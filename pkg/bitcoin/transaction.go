package bitcoin

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Transaction represents a Bitcoin transaction
type Transaction struct {
	Version    uint32     `json:"version"`
	Inputs     []TxInput  `json:"inputs"`
	Outputs    []TxOutput `json:"outputs"`
	LockTime   uint32     `json:"locktime"`

	// Witness data for SegWit transactions
	Witnesses []TxWitness `json:"witnesses,omitempty"`

	// Cached values
	hash   *Hash256 // Transaction ID (excludes witness data)
	wthash *Hash256 // Witness Transaction ID (includes witness data)
}

// TxInput represents a transaction input
type TxInput struct {
	PreviousOutput OutPoint  `json:"previous_output"`
	ScriptSig      []byte    `json:"script_sig"`
	Sequence       uint32    `json:"sequence"`
	Witness        [][]byte  `json:"witness,omitempty"` // SegWit witness data
}

// TxOutput represents a transaction output
type TxOutput struct {
	Value        uint64 `json:"value"`        // Amount in satoshis
	ScriptPubKey []byte `json:"script_pubkey"`
}

// TxWitness represents witness data for a SegWit input
type TxWitness struct {
	Stack [][]byte `json:"stack"`
}

// OutPoint represents a reference to a transaction output
type OutPoint struct {
	Hash  Hash256 `json:"hash"`  // Transaction hash
	Index uint32  `json:"index"` // Output index
}

// NewTransaction creates a new transaction
func NewTransaction(version uint32, inputs []TxInput, outputs []TxOutput, lockTime uint32) *Transaction {
	return &Transaction{
		Version:  version,
		Inputs:   inputs,
		Outputs:  outputs,
		LockTime: lockTime,
	}
}

// EncodeVarInt encodes an integer as a Bitcoin variable-length integer
func EncodeVarInt(value uint64) []byte {
	if value < 0xfd {
		return []byte{byte(value)}
	} else if value <= 0xffff {
		buf := make([]byte, 3)
		buf[0] = 0xfd
		binary.LittleEndian.PutUint16(buf[1:], uint16(value))
		return buf
	} else if value <= 0xffffffff {
		buf := make([]byte, 5)
		buf[0] = 0xfe
		binary.LittleEndian.PutUint32(buf[1:], uint32(value))
		return buf
	} else {
		buf := make([]byte, 9)
		buf[0] = 0xff
		binary.LittleEndian.PutUint64(buf[1:], value)
		return buf
	}
}

// DecodeVarInt decodes a Bitcoin variable-length integer
func DecodeVarInt(data []byte) (value uint64, bytesRead int, err error) {
	if len(data) == 0 {
		return 0, 0, fmt.Errorf("empty data")
	}

	first := data[0]
	if first < 0xfd {
		return uint64(first), 1, nil
	} else if first == 0xfd {
		if len(data) < 3 {
			return 0, 0, fmt.Errorf("insufficient data for fd varint")
		}
		return uint64(binary.LittleEndian.Uint16(data[1:3])), 3, nil
	} else if first == 0xfe {
		if len(data) < 5 {
			return 0, 0, fmt.Errorf("insufficient data for fe varint")
		}
		return uint64(binary.LittleEndian.Uint32(data[1:5])), 5, nil
	} else { // first == 0xff
		if len(data) < 9 {
			return 0, 0, fmt.Errorf("insufficient data for ff varint")
		}
		return binary.LittleEndian.Uint64(data[1:9]), 9, nil
	}
}

// Serialize converts the transaction to Bitcoin wire format
func (tx *Transaction) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	hasWitness := tx.HasWitness()

	// Version (4 bytes, little-endian)
	if err := binary.Write(&buf, binary.LittleEndian, tx.Version); err != nil {
		return nil, fmt.Errorf("failed to write version: %w", err)
	}

	// SegWit flag and marker (if witness data present)
	if hasWitness {
		buf.WriteByte(0x00) // marker
		buf.WriteByte(0x01) // flag
	}

	// Input count (varint)
	buf.Write(EncodeVarInt(uint64(len(tx.Inputs))))

	// Inputs
	for _, input := range tx.Inputs {
		// Previous output hash (32 bytes, reversed for wire format)
		hashBytes := input.PreviousOutput.Hash.Bytes()
		for i := len(hashBytes) - 1; i >= 0; i-- {
			buf.WriteByte(hashBytes[i])
		}
		// Previous output index (4 bytes, little-endian)
		if err := binary.Write(&buf, binary.LittleEndian, input.PreviousOutput.Index); err != nil {
			return nil, fmt.Errorf("failed to write previous output index: %w", err)
		}
		// Script length (varint)
		buf.Write(EncodeVarInt(uint64(len(input.ScriptSig))))
		// Script
		buf.Write(input.ScriptSig)
		// Sequence (4 bytes, little-endian)
		if err := binary.Write(&buf, binary.LittleEndian, input.Sequence); err != nil {
			return nil, fmt.Errorf("failed to write sequence: %w", err)
		}
	}

	// Output count (varint)
	buf.Write(EncodeVarInt(uint64(len(tx.Outputs))))

	// Outputs
	for _, output := range tx.Outputs {
		// Value (8 bytes, little-endian)
		if err := binary.Write(&buf, binary.LittleEndian, output.Value); err != nil {
			return nil, fmt.Errorf("failed to write output value: %w", err)
		}
		// Script length (varint)
		buf.Write(EncodeVarInt(uint64(len(output.ScriptPubKey))))
		// Script
		buf.Write(output.ScriptPubKey)
	}

	// Witness data (if present)
	if hasWitness {
		// Witness data for each input
		for _, input := range tx.Inputs {
			// Number of witness elements (varint)
			buf.Write(EncodeVarInt(uint64(len(input.Witness))))
			// Witness elements
			for _, element := range input.Witness {
				// Element length (varint)
				buf.Write(EncodeVarInt(uint64(len(element))))
				// Element data
				buf.Write(element)
			}
		}

		// Also handle transaction-level witnesses if present
		for _, witness := range tx.Witnesses {
			// Number of witness elements (varint)
			buf.Write(EncodeVarInt(uint64(len(witness.Stack))))
			// Witness elements
			for _, element := range witness.Stack {
				// Element length (varint)
				buf.Write(EncodeVarInt(uint64(len(element))))
				// Element data
				buf.Write(element)
			}
		}
	}

	// Locktime (4 bytes, little-endian)
	if err := binary.Write(&buf, binary.LittleEndian, tx.LockTime); err != nil {
		return nil, fmt.Errorf("failed to write locktime: %w", err)
	}

	return buf.Bytes(), nil
}

// DeserializeTransaction deserializes a transaction from Bitcoin wire format
func DeserializeTransaction(data []byte) (*Transaction, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty transaction data")
	}

	tx := &Transaction{}
	offset := 0

	// Version (4 bytes)
	if len(data[offset:]) < 4 {
		return nil, fmt.Errorf("insufficient data for version")
	}
	tx.Version = binary.LittleEndian.Uint32(data[offset : offset+4])
	offset += 4

	// Input count
	inputCount, bytesRead, err := DecodeVarInt(data[offset:])
	if err != nil {
		return nil, fmt.Errorf("failed to decode input count: %v", err)
	}
	offset += bytesRead

	// Validate input count before conversion
	if inputCount > 0x7fffffff { // Max int value
		return nil, fmt.Errorf("input count too large: %d", inputCount)
	}

	// Inputs
	tx.Inputs = make([]TxInput, int(inputCount))
	for i := uint64(0); i < inputCount; i++ {
		// Previous output hash (32 bytes, reversed from wire format)
		if len(data[offset:]) < 32 {
			return nil, fmt.Errorf("insufficient data for input %d hash", i)
		}
		// Reverse the hash bytes from wire format
		for j := 0; j < 32; j++ {
			tx.Inputs[i].PreviousOutput.Hash[j] = data[offset+31-j]
		}
		offset += 32

		// Previous output index (4 bytes)
		if len(data[offset:]) < 4 {
			return nil, fmt.Errorf("insufficient data for input %d index", i)
		}
		tx.Inputs[i].PreviousOutput.Index = binary.LittleEndian.Uint32(data[offset : offset+4])
		offset += 4

		// Script length
		scriptLen, bytesRead, err := DecodeVarInt(data[offset:])
		if err != nil {
			return nil, fmt.Errorf("failed to decode input %d script length: %v", i, err)
		}
		offset += bytesRead

		// Script
		// Validate script length before conversion
		if scriptLen > 0x7fffffff { // Max int value
			return nil, fmt.Errorf("input %d script length too large: %d", i, scriptLen)
		}
		scriptLenInt := int(scriptLen)
		if len(data[offset:]) < scriptLenInt {
			return nil, fmt.Errorf("insufficient data for input %d script", i)
		}
		tx.Inputs[i].ScriptSig = make([]byte, scriptLen)
		copy(tx.Inputs[i].ScriptSig, data[offset:offset+scriptLenInt])
		offset += scriptLenInt

		// Sequence (4 bytes)
		if len(data[offset:]) < 4 {
			return nil, fmt.Errorf("insufficient data for input %d sequence", i)
		}
		tx.Inputs[i].Sequence = binary.LittleEndian.Uint32(data[offset : offset+4])
		offset += 4
	}

	// Output count
	outputCount, bytesRead, err := DecodeVarInt(data[offset:])
	if err != nil {
		return nil, fmt.Errorf("failed to decode output count: %v", err)
	}
	offset += bytesRead

	// Validate output count before conversion
	if outputCount > 0x7fffffff { // Max int value
		return nil, fmt.Errorf("output count too large: %d", outputCount)
	}

	// Outputs
	tx.Outputs = make([]TxOutput, int(outputCount))
	for i := uint64(0); i < outputCount; i++ {
		// Value (8 bytes)
		if len(data[offset:]) < 8 {
			return nil, fmt.Errorf("insufficient data for output %d value", i)
		}
		tx.Outputs[i].Value = binary.LittleEndian.Uint64(data[offset : offset+8])
		offset += 8

		// Script length
		scriptLen, bytesRead, err := DecodeVarInt(data[offset:])
		if err != nil {
			return nil, fmt.Errorf("failed to decode output %d script length: %v", i, err)
		}
		offset += bytesRead

		// Script
		// Validate script length before conversion
		if scriptLen > 0x7fffffff { // Max int value
			return nil, fmt.Errorf("output %d script length too large: %d", i, scriptLen)
		}
		scriptLenInt := int(scriptLen)
		if len(data[offset:]) < scriptLenInt {
			return nil, fmt.Errorf("insufficient data for output %d script", i)
		}
		tx.Outputs[i].ScriptPubKey = make([]byte, scriptLen)
		copy(tx.Outputs[i].ScriptPubKey, data[offset:offset+scriptLenInt])
		offset += scriptLenInt
	}

	// Locktime (4 bytes)
	if len(data[offset:]) < 4 {
		return nil, fmt.Errorf("insufficient data for locktime")
	}
	tx.LockTime = binary.LittleEndian.Uint32(data[offset : offset+4])
	offset += 4

	// TODO: Handle witness data for SegWit transactions
	// This would require detecting the witness flag and parsing witness data

	return tx, nil
}

// Hash returns the transaction ID (excludes witness data)
func (tx *Transaction) Hash() Hash256 {
	if tx.hash == nil {
		// TODO: Implement transaction serialization and hashing
		// For now, return zero hash
		hash := ZeroHash
		tx.hash = &hash
	}
	return *tx.hash
}

// WitnessHash returns the witness transaction ID (includes witness data)
func (tx *Transaction) WitnessHash() Hash256 {
	if tx.wthash == nil {
		// TODO: Implement witness transaction serialization and hashing
		// For now, return zero hash
		hash := ZeroHash
		tx.wthash = &hash
	}
	return *tx.wthash
}

// IsCoinbase returns true if this is a coinbase transaction
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 &&
		   tx.Inputs[0].PreviousOutput.Hash.IsZero() &&
		   tx.Inputs[0].PreviousOutput.Index == 0xffffffff
}

// HasWitness returns true if the transaction has witness data
func (tx *Transaction) HasWitness() bool {
	// Check witness data in transaction-level field
	if len(tx.Witnesses) > 0 {
		return true
	}

	// Check witness data in individual inputs
	for _, input := range tx.Inputs {
		if len(input.Witness) > 0 {
			return true
		}
	}

	return false
}

// TotalOutput calculates the total value of all outputs
func (tx *Transaction) TotalOutput() uint64 {
	var total uint64
	for _, output := range tx.Outputs {
		total += output.Value
	}
	return total
}

// IsStandard checks if the transaction follows standard rules
func (tx *Transaction) IsStandard() bool {
	// TODO: Implement standard transaction checks
	// - Version check
	// - Size limits
	// - Standard script types
	// - Dust threshold
	return true // Placeholder
}

// Validate performs basic validation checks
func (tx *Transaction) Validate() error {
	// Basic sanity checks
	if len(tx.Inputs) == 0 {
		return fmt.Errorf("transaction has no inputs")
	}

	if len(tx.Outputs) == 0 {
		return fmt.Errorf("transaction has no outputs")
	}

	// Check for duplicate inputs
	seen := make(map[OutPoint]bool)
	for _, input := range tx.Inputs {
		if seen[input.PreviousOutput] {
			return fmt.Errorf("transaction has duplicate inputs")
		}
		seen[input.PreviousOutput] = true
	}

	// Check output values
	for i, output := range tx.Outputs {
		if output.Value > MaxMoney {
			return fmt.Errorf("output %d value exceeds maximum", i)
		}
	}

	// Check total output value
	if tx.TotalOutput() > MaxMoney {
		return fmt.Errorf("total output value exceeds maximum")
	}

	return nil
}

// Constants
const (
	MaxMoney = 21000000 * 100000000 // 21 million BTC in satoshis
)

// String returns a string representation of the OutPoint
func (op OutPoint) String() string {
	return fmt.Sprintf("%s:%d", op.Hash.String(), op.Index)
}

// IsNull returns true if the outpoint is null (coinbase)
func (op OutPoint) IsNull() bool {
	return op.Hash.IsZero() && op.Index == 0xffffffff
}

