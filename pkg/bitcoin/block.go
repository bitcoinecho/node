package bitcoin

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
)

// Block represents a Bitcoin block
type Block struct {
	Header       BlockHeader   `json:"header"`
	Transactions []Transaction `json:"transactions"`

	// Cached values
	hash   *Hash256 // Block hash
	height *int32   // Block height (set when connected to chain)
}

// BlockHeader represents a Bitcoin block header
type BlockHeader struct {
	Version       uint32    `json:"version"`
	PrevBlockHash Hash256   `json:"prev_block_hash"`
	MerkleRoot    Hash256   `json:"merkle_root"`
	Timestamp     uint32    `json:"timestamp"`
	Bits          uint32    `json:"bits"`          // Difficulty target
	Nonce         uint32    `json:"nonce"`

	// Cached values
	hash *Hash256 // Header hash
}

// NewBlock creates a new block
func NewBlock(header BlockHeader, transactions []Transaction) *Block { //nolint:gocritic // header copied intentionally for immutability
	return &Block{
		Header:       header,
		Transactions: transactions,
	}
}

// NewBlockHeader creates a new block header
func NewBlockHeader(version uint32, prevHash, merkleRoot Hash256, timestamp, bits, nonce uint32) BlockHeader {
	return BlockHeader{
		Version:       version,
		PrevBlockHash: prevHash,
		MerkleRoot:    merkleRoot,
		Timestamp:     timestamp,
		Bits:          bits,
		Nonce:         nonce,
	}
}

// Hash returns the block hash
func (b *Block) Hash() Hash256 {
	if b.hash == nil {
		hash := b.Header.Hash()
		b.hash = &hash
	}
	return *b.hash
}

// Height returns the block height if known
func (b *Block) Height() *int32 {
	return b.height
}

// SetHeight sets the block height
func (b *Block) SetHeight(height int32) {
	b.height = &height
}

// Hash returns the header hash
func (bh *BlockHeader) Hash() Hash256 {
	if bh.hash == nil {
		// Serialize block header and hash with double SHA-256
		serialized, err := bh.serialize()
		if err != nil {
			// In case of serialization error, return zero hash
			return ZeroHash
		}
		rawHash := DoubleHashSHA256(serialized)

		// Bitcoin displays block hashes in reverse byte order
		var hash Hash256
		rawBytes := rawHash.Bytes()
		for i := 0; i < 32; i++ {
			hash[i] = rawBytes[31-i]
		}

		bh.hash = &hash
	}
	return *bh.hash
}

// serialize serializes the block header for hashing
func (bh *BlockHeader) serialize() ([]byte, error) {
	buf := new(bytes.Buffer)

	// Version (4 bytes, little-endian)
	if err := binary.Write(buf, binary.LittleEndian, bh.Version); err != nil {
		return nil, fmt.Errorf("failed to write version: %w", err)
	}

	// Previous block hash (32 bytes) - Bitcoin stores hashes in reverse order
	prevHashBytes := bh.PrevBlockHash.Bytes()
	for i := len(prevHashBytes) - 1; i >= 0; i-- {
		buf.WriteByte(prevHashBytes[i])
	}

	// Merkle root (32 bytes) - Bitcoin stores hashes in reverse order
	merkleBytes := bh.MerkleRoot.Bytes()
	for i := len(merkleBytes) - 1; i >= 0; i-- {
		buf.WriteByte(merkleBytes[i])
	}

	// Timestamp (4 bytes, little-endian)
	if err := binary.Write(buf, binary.LittleEndian, bh.Timestamp); err != nil {
		return nil, fmt.Errorf("failed to write timestamp: %w", err)
	}

	// Bits (4 bytes, little-endian)
	if err := binary.Write(buf, binary.LittleEndian, bh.Bits); err != nil {
		return nil, fmt.Errorf("failed to write bits: %w", err)
	}

	// Nonce (4 bytes, little-endian)
	if err := binary.Write(buf, binary.LittleEndian, bh.Nonce); err != nil {
		return nil, fmt.Errorf("failed to write nonce: %w", err)
	}

	return buf.Bytes(), nil
}

// Time returns the block timestamp as a time.Time
func (bh *BlockHeader) Time() time.Time {
	return time.Unix(int64(bh.Timestamp), 0)
}

// Difficulty returns the difficulty target as a compact representation
func (bh *BlockHeader) Difficulty() uint32 {
	return bh.Bits
}

// IsGenesis returns true if this is the genesis block
func (b *Block) IsGenesis() bool {
	return b.Header.PrevBlockHash.IsZero()
}

// TransactionCount returns the number of transactions in the block
func (b *Block) TransactionCount() int {
	return len(b.Transactions)
}

// HasCoinbase returns true if the block has a coinbase transaction
func (b *Block) HasCoinbase() bool {
	return len(b.Transactions) > 0 && b.Transactions[0].IsCoinbase()
}

// CoinbaseTransaction returns the coinbase transaction if present
func (b *Block) CoinbaseTransaction() *Transaction {
	if b.HasCoinbase() {
		return &b.Transactions[0]
	}
	return nil
}

// Size returns the serialized size of the block in bytes
func (b *Block) Size() int {
	// TODO: Implement actual block serialization to get accurate size
	// This is a placeholder estimation
	size := 80 // Block header size
	size += 9  // Approximate varint for transaction count

	for i := range b.Transactions {
		size += b.estimateTransactionSize(&b.Transactions[i])
	}

	return size
}

// Weight returns the block weight as defined by BIP141
func (b *Block) Weight() int {
	// TODO: Implement BIP141 weight calculation
	// Weight = (base_size * 3) + total_size
	// where base_size excludes witness data and total_size includes it
	return b.Size() * 4 // Placeholder: assume no witness data
}

// estimateTransactionSize estimates the serialized size of a transaction
func (b *Block) estimateTransactionSize(tx *Transaction) int {
	// Rough estimation - in practice this should serialize the transaction
	size := 4  // Version
	size += 1  // Input count varint (simplified)
	size += len(tx.Inputs) * 36  // Each input: 32-byte hash + 4-byte index + script + sequence
	size += 1  // Output count varint (simplified)
	size += len(tx.Outputs) * 9 // Each output: 8-byte value + script (simplified)
	size += 4  // Lock time

	// Add script sizes
	for _, input := range tx.Inputs {
		size += len(input.ScriptSig) + 1 // Script + length varint
	}
	for _, output := range tx.Outputs {
		size += len(output.ScriptPubKey) + 1 // Script + length varint
	}

	return size
}

// Validate performs basic block validation
func (b *Block) Validate() error {
	// Check if block has transactions
	if len(b.Transactions) == 0 {
		return fmt.Errorf("block has no transactions")
	}

	// Check if first transaction is coinbase
	if !b.Transactions[0].IsCoinbase() {
		return fmt.Errorf("first transaction is not coinbase")
	}

	// Check that only the first transaction is coinbase
	for i, tx := range b.Transactions[1:] {
		if tx.IsCoinbase() {
			return fmt.Errorf("transaction %d is coinbase (only first can be)", i+1)
		}
	}

	// Validate each transaction
	for i, tx := range b.Transactions {
		if err := tx.Validate(); err != nil {
			return fmt.Errorf("transaction %d validation failed: %v", i, err)
		}
	}

	// Check block size/weight limits
	if b.Size() > MaxBlockSize {
		return fmt.Errorf("block size %d exceeds maximum %d", b.Size(), MaxBlockSize)
	}

	if b.Weight() > MaxBlockWeight {
		return fmt.Errorf("block weight %d exceeds maximum %d", b.Weight(), MaxBlockWeight)
	}

	// TODO: Additional validations:
	// - Merkle root validation
	// - Proof of work validation
	// - Timestamp validation
	// - Difficulty target validation

	return nil
}

// ValidateHeader performs block header validation
func (bh *BlockHeader) Validate() error {
	// Check timestamp is not too far in future (2 hours)
	maxTime := time.Now().Add(2 * time.Hour)
	if bh.Time().After(maxTime) {
		return fmt.Errorf("block timestamp too far in future")
	}

	// TODO: Additional header validations:
	// - Version validation
	// - Proof of work validation
	// - Difficulty target validation

	return nil
}

// Constants
const (
	MaxBlockSize   = 1000000  // 1MB (legacy limit)
	MaxBlockWeight = 4000000  // 4M weight units (BIP141)
)