package bitcoin

import (
	"testing"
	"time"
)

// TestNewBlock tests block creation
func TestNewBlock(t *testing.T) {
	// Create a sample block header
	header := NewBlockHeader(
		1,                // version
		ZeroHash, // previous block hash (genesis)
		ZeroHash, // merkle root (placeholder)
		1640995200,       // timestamp (Jan 1, 2022)
		0x1d00ffff,       // difficulty bits
		12345,            // nonce
	)

	// Create a sample coinbase transaction
	coinbaseTx := &Transaction{
		Version: 1,
		Inputs: []TxInput{{
			PreviousOutput: OutPoint{Hash: ZeroHash, Index: 0xffffffff},
			ScriptSig:      []byte("Genesis block coinbase"),
			Sequence:       0xffffffff,
		}},
		Outputs: []TxOutput{{
			Value:        5000000000,               // 50 BTC
			ScriptPubKey: []byte{0x76, 0xa9, 0x14}, // P2PKH placeholder
		}},
		LockTime: 0,
	}

	block := NewBlock(header, []Transaction{*coinbaseTx})

	if block.Header.Version != 1 {
		t.Errorf("expected version 1, got %d", block.Header.Version)
	}

	if len(block.Transactions) != 1 {
		t.Errorf("expected 1 transaction, got %d", len(block.Transactions))
	}

	if !block.HasCoinbase() {
		t.Errorf("expected block to have coinbase transaction")
	}
}

// TestNewBlockHeader tests block header creation
func TestNewBlockHeader(t *testing.T) {
	prevHash := Hash256{0x01, 0x02, 0x03}
	merkleRoot := Hash256{0x04, 0x05, 0x06}

	header := NewBlockHeader(2, prevHash, merkleRoot, 1640995200, 0x1d00ffff, 54321)

	if header.Version != 2 {
		t.Errorf("expected version 2, got %d", header.Version)
	}

	if header.PrevBlockHash != prevHash {
		t.Errorf("expected prev hash %v, got %v", prevHash, header.PrevBlockHash)
	}

	if header.MerkleRoot != merkleRoot {
		t.Errorf("expected merkle root %v, got %v", merkleRoot, header.MerkleRoot)
	}

	if header.Timestamp != 1640995200 {
		t.Errorf("expected timestamp 1640995200, got %d", header.Timestamp)
	}

	if header.Bits != 0x1d00ffff {
		t.Errorf("expected bits 0x1d00ffff, got 0x%x", header.Bits)
	}

	if header.Nonce != 54321 {
		t.Errorf("expected nonce 54321, got %d", header.Nonce)
	}
}

// TestBlock_IsGenesis tests genesis block detection
func TestBlock_IsGenesis(t *testing.T) {
	tests := []struct {
		name         string
		prevHash     Hash256
		expectedBool bool
	}{
		{
			name:         "genesis block (zero prev hash)",
			prevHash:     ZeroHash,
			expectedBool: true,
		},
		{
			name:         "non-genesis block",
			prevHash:     Hash256{0x01},
			expectedBool: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := NewBlockHeader(1, tt.prevHash, ZeroHash, 1640995200, 0x1d00ffff, 0)
			block := NewBlock(header, []Transaction{})

			result := block.IsGenesis()
			if result != tt.expectedBool {
				t.Errorf("expected %v, got %v", tt.expectedBool, result)
			}
		})
	}
}

// TestBlock_HasCoinbase tests coinbase transaction detection
func TestBlock_HasCoinbase(t *testing.T) {
	tests := []struct {
		name         string
		transactions []Transaction
		expected     bool
	}{
		{
			name:         "no transactions",
			transactions: []Transaction{},
			expected:     false,
		},
		{
			name: "has coinbase transaction",
			transactions: []Transaction{{
				Version: 1,
				Inputs: []TxInput{{
					PreviousOutput: OutPoint{Hash: ZeroHash, Index: 0xffffffff},
					ScriptSig:      []byte("coinbase"),
					Sequence:       0xffffffff,
				}},
				Outputs: []TxOutput{{Value: 5000000000, ScriptPubKey: []byte{0x76}}},
			}},
			expected: true,
		},
		{
			name: "first transaction not coinbase",
			transactions: []Transaction{{
				Version: 1,
				Inputs: []TxInput{{
					PreviousOutput: OutPoint{Hash: Hash256{0x01}, Index: 0},
					ScriptSig:      []byte{},
					Sequence:       0xffffffff,
				}},
				Outputs: []TxOutput{{Value: 1000000, ScriptPubKey: []byte{0x76}}},
			}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := NewBlockHeader(1, ZeroHash, ZeroHash, 1640995200, 0x1d00ffff, 0)
			block := NewBlock(header, tt.transactions)

			result := block.HasCoinbase()
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestBlock_CoinbaseTransaction tests coinbase transaction retrieval
func TestBlock_CoinbaseTransaction(t *testing.T) {
	// Create coinbase transaction
	coinbaseTx := Transaction{
		Version: 1,
		Inputs: []TxInput{{
			PreviousOutput: OutPoint{Hash: ZeroHash, Index: 0xffffffff},
			ScriptSig:      []byte("Genesis coinbase"),
			Sequence:       0xffffffff,
		}},
		Outputs: []TxOutput{{Value: 5000000000, ScriptPubKey: []byte{0x76, 0xa9}}},
	}

	// Create regular transaction
	regularTx := Transaction{
		Version: 1,
		Inputs: []TxInput{{
			PreviousOutput: OutPoint{Hash: Hash256{0x01}, Index: 0},
			ScriptSig:      []byte{0x76, 0xa9},
			Sequence:       0xffffffff,
		}},
		Outputs: []TxOutput{{Value: 1000000, ScriptPubKey: []byte{0x76}}},
	}

	tests := []struct {
		name         string
		transactions []Transaction
		expectNil    bool
	}{
		{
			name:         "no transactions",
			transactions: []Transaction{},
			expectNil:    true,
		},
		{
			name:         "has coinbase transaction",
			transactions: []Transaction{coinbaseTx, regularTx},
			expectNil:    false,
		},
		{
			name:         "first transaction not coinbase",
			transactions: []Transaction{regularTx},
			expectNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := NewBlockHeader(1, ZeroHash, ZeroHash, 1640995200, 0x1d00ffff, 0)
			block := NewBlock(header, tt.transactions)

			coinbase := block.CoinbaseTransaction()

			if tt.expectNil {
				if coinbase != nil {
					t.Errorf("expected nil coinbase transaction, got %v", coinbase)
				}
			} else {
				if coinbase == nil {
					t.Errorf("expected coinbase transaction, got nil")
				} else if !coinbase.IsCoinbase() {
					t.Errorf("returned transaction is not coinbase")
				}
			}
		})
	}
}

// TestBlock_TransactionCount tests transaction counting
func TestBlock_TransactionCount(t *testing.T) {
	tests := []struct {
		name     string
		txCount  int
		expected int
	}{
		{"empty block", 0, 0},
		{"single transaction", 1, 1},
		{"multiple transactions", 5, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transactions := make([]Transaction, tt.txCount)
			for i := 0; i < tt.txCount; i++ {
				transactions[i] = Transaction{Version: 1}
			}

			header := NewBlockHeader(1, ZeroHash, ZeroHash, 1640995200, 0x1d00ffff, 0)
			block := NewBlock(header, transactions)

			count := block.TransactionCount()
			if count != tt.expected {
				t.Errorf("expected %d transactions, got %d", tt.expected, count)
			}
		})
	}
}

// TestBlock_Height tests block height management
func TestBlock_Height(t *testing.T) {
	header := NewBlockHeader(1, ZeroHash, ZeroHash, 1640995200, 0x1d00ffff, 0)
	block := NewBlock(header, []Transaction{})

	// Initially height should be nil
	if block.Height() != nil {
		t.Errorf("expected nil height initially, got %v", *block.Height())
	}

	// Set height
	block.SetHeight(100)

	if block.Height() == nil {
		t.Errorf("expected height to be set, got nil")
	} else if *block.Height() != 100 {
		t.Errorf("expected height 100, got %d", *block.Height())
	}

	// Update height
	block.SetHeight(200)
	if *block.Height() != 200 {
		t.Errorf("expected height 200, got %d", *block.Height())
	}
}

// TestBlockHeader_Time tests timestamp conversion
func TestBlockHeader_Time(t *testing.T) {
	timestamp := uint32(1640995200) // Jan 1, 2022 00:00:00 UTC
	header := NewBlockHeader(1, ZeroHash, ZeroHash, timestamp, 0x1d00ffff, 0)

	expectedTime := time.Unix(1640995200, 0)
	actualTime := header.Time()

	if !actualTime.Equal(expectedTime) {
		t.Errorf("expected time %v, got %v", expectedTime, actualTime)
	}
}

// TestBlockHeader_Difficulty tests difficulty access
func TestBlockHeader_Difficulty(t *testing.T) {
	bits := uint32(0x1d00ffff)
	header := NewBlockHeader(1, ZeroHash, ZeroHash, 1640995200, bits, 0)

	difficulty := header.Difficulty()
	if difficulty != bits {
		t.Errorf("expected difficulty 0x%x, got 0x%x", bits, difficulty)
	}
}

// TestBlock_Hash tests block hashing (currently returns zero)
func TestBlock_Hash(t *testing.T) {
	header := NewBlockHeader(1, ZeroHash, ZeroHash, 1640995200, 0x1d00ffff, 0)
	block := NewBlock(header, []Transaction{})

	hash1 := block.Hash()
	hash2 := block.Hash()

	// Hash should be consistent (cached)
	if hash1 != hash2 {
		t.Errorf("block hash not consistent: %s != %s", hash1.String(), hash2.String())
	}

	// Currently returns zero hash (TODO: implement actual hashing)
	if !hash1.IsZero() {
		t.Logf("Block hash implementation complete: %s", hash1.String())
	}
}

// TestBlockHeader_Hash tests header hashing with real implementation
func TestBlockHeader_Hash(t *testing.T) {
	header := NewBlockHeader(1, ZeroHash, ZeroHash, 1640995200, 0x1d00ffff, 12345)

	hash1 := header.Hash()
	hash2 := header.Hash()

	// Hash should be consistent (cached)
	if hash1 != hash2 {
		t.Errorf("header hash not consistent: %s != %s", hash1.String(), hash2.String())
	}

	// Hash should not be zero (real implementation now)
	if hash1.IsZero() {
		t.Errorf("header hash should not be zero with real implementation")
	}

	t.Logf("Header hash: %s", hash1.String())
}

// TestBlockHeader_Hash_GenesisBlock tests with Bitcoin Genesis Block data
func TestBlockHeader_Hash_GenesisBlock(t *testing.T) {
	// Bitcoin Genesis Block header data
	// https://blockstream.info/block/000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f

	merkleRoot, err := NewHash256FromString("4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b")
	if err != nil {
		t.Fatalf("failed to create merkle root hash: %v", err)
	}

	// Genesis block header
	genesisHeader := NewBlockHeader(
		1,                // version
		ZeroHash, // previous block hash (genesis has no previous)
		merkleRoot,       // merkle root
		1231006505,       // timestamp (Jan 3, 2009 18:15:05 UTC)
		0x1d00ffff,       // bits (difficulty)
		2083236893,       // nonce
	)

	// Expected genesis block hash
	expectedHash := "000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"

	actualHash := genesisHeader.Hash()
	actualHashStr := actualHash.String()

	if actualHashStr != expectedHash {
		t.Errorf("Genesis block hash mismatch:\n  expected: %s\n  actual:   %s", expectedHash, actualHashStr)
	}

	t.Logf("✅ Genesis block hash verified: %s", actualHashStr)
}

// TestBlock_Validate tests basic block validation
func TestBlock_Validate(t *testing.T) {
	// Create valid coinbase transaction
	validCoinbase := Transaction{
		Version: 1,
		Inputs: []TxInput{{
			PreviousOutput: OutPoint{Hash: ZeroHash, Index: 0xffffffff},
			ScriptSig:      []byte("Block validation test"),
			Sequence:       0xffffffff,
		}},
		Outputs: []TxOutput{{
			Value:        5000000000,
			ScriptPubKey: []byte{0x76, 0xa9, 0x14},
		}},
		LockTime: 0,
	}

	// Create valid regular transaction
	validRegular := Transaction{
		Version: 1,
		Inputs: []TxInput{{
			PreviousOutput: OutPoint{Hash: Hash256{0x01}, Index: 0},
			ScriptSig:      []byte{0x76, 0xa9},
			Sequence:       0xffffffff,
		}},
		Outputs: []TxOutput{{
			Value:        1000000000,
			ScriptPubKey: []byte{0x76, 0xa9},
		}},
		LockTime: 0,
	}

	tests := []struct {
		name         string
		transactions []Transaction
		expectError  bool
		errorMsg     string
	}{
		{
			name:         "valid block with coinbase only",
			transactions: []Transaction{validCoinbase},
			expectError:  false,
		},
		{
			name:         "valid block with coinbase and regular tx",
			transactions: []Transaction{validCoinbase, validRegular},
			expectError:  false,
		},
		{
			name:         "no transactions",
			transactions: []Transaction{},
			expectError:  true,
			errorMsg:     "block has no transactions",
		},
		{
			name:         "first transaction not coinbase",
			transactions: []Transaction{validRegular},
			expectError:  true,
			errorMsg:     "first transaction is not coinbase",
		},
		{
			name:         "multiple coinbase transactions",
			transactions: []Transaction{validCoinbase, validCoinbase},
			expectError:  true,
			errorMsg:     "transaction 1 is coinbase (only first can be)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := NewBlockHeader(1, ZeroHash, ZeroHash, 1640995200, 0x1d00ffff, 0)
			block := NewBlock(header, tt.transactions)

			err := block.Validate()

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if err.Error() != tt.errorMsg {
					t.Errorf("expected error '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

// TestBlockHeader_Validate tests block header validation
func TestBlockHeader_Validate(t *testing.T) {
	tests := []struct {
		name        string
		timestamp   uint32
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid timestamp (current)",
			timestamp:   uint32(time.Now().Unix()),
			expectError: false,
		},
		{
			name:        "valid timestamp (1 hour in future)",
			timestamp:   uint32(time.Now().Add(1 * time.Hour).Unix()),
			expectError: false,
		},
		{
			name:        "invalid timestamp (3 hours in future)",
			timestamp:   uint32(time.Now().Add(3 * time.Hour).Unix()),
			expectError: true,
			errorMsg:    "block timestamp too far in future",
		},
		{
			name:        "valid timestamp (past)",
			timestamp:   uint32(time.Now().Add(-24 * time.Hour).Unix()),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := NewBlockHeader(1, ZeroHash, ZeroHash, tt.timestamp, 0x1d00ffff, 0)

			err := header.Validate()

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if err.Error() != tt.errorMsg {
					t.Errorf("expected error '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

// TestBlock_Size tests block size estimation
func TestBlock_Size(t *testing.T) {
	// Create simple transaction for size testing
	simpleTx := Transaction{
		Version: 1,
		Inputs: []TxInput{{
			PreviousOutput: OutPoint{Hash: ZeroHash, Index: 0xffffffff},
			ScriptSig:      []byte("simple"),
			Sequence:       0xffffffff,
		}},
		Outputs: []TxOutput{{
			Value:        5000000000,
			ScriptPubKey: []byte{0x76, 0xa9},
		}},
		LockTime: 0,
	}

	header := NewBlockHeader(1, ZeroHash, ZeroHash, 1640995200, 0x1d00ffff, 0)
	block := NewBlock(header, []Transaction{simpleTx})

	size := block.Size()

	// Should be greater than just header size (80 bytes)
	if size <= 80 {
		t.Errorf("expected block size > 80 bytes, got %d", size)
	}

	// Should be reasonable for a single simple transaction
	if size > 1000 {
		t.Errorf("block size seems too large for simple transaction: %d bytes", size)
	}
}

// TestBlock_Weight tests block weight calculation
func TestBlock_Weight(t *testing.T) {
	simpleTx := Transaction{
		Version: 1,
		Inputs: []TxInput{{
			PreviousOutput: OutPoint{Hash: ZeroHash, Index: 0xffffffff},
			ScriptSig:      []byte("weight test"),
			Sequence:       0xffffffff,
		}},
		Outputs: []TxOutput{{
			Value:        5000000000,
			ScriptPubKey: []byte{0x76, 0xa9},
		}},
		LockTime: 0,
	}

	header := NewBlockHeader(1, ZeroHash, ZeroHash, 1640995200, 0x1d00ffff, 0)
	block := NewBlock(header, []Transaction{simpleTx})

	weight := block.Weight()
	size := block.Size()

	// Currently weight = size * 4 (placeholder implementation)
	expectedWeight := size * 4
	if weight != expectedWeight {
		t.Errorf("expected weight %d, got %d", expectedWeight, weight)
	}

	// Weight should be within limits for a simple block
	if weight > MaxBlockWeight {
		t.Errorf("block weight %d exceeds maximum %d", weight, MaxBlockWeight)
	}
}

// TestBlockConstants tests Bitcoin block constants
func TestBlockConstants(t *testing.T) {
	// Test MaxBlockSize constant
	expectedMaxSize := 1000000 // 1MB
	if MaxBlockSize != expectedMaxSize {
		t.Errorf("MaxBlockSize should be %d, got %d", expectedMaxSize, MaxBlockSize)
	}

	// Test MaxBlockWeight constant
	expectedMaxWeight := 4000000 // 4M weight units
	if MaxBlockWeight != expectedMaxWeight {
		t.Errorf("MaxBlockWeight should be %d, got %d", expectedMaxWeight, MaxBlockWeight)
	}
}

// BenchmarkBlock_Hash benchmarks block hashing
func BenchmarkBlock_Hash(b *testing.B) {
	header := NewBlockHeader(1, ZeroHash, ZeroHash, 1640995200, 0x1d00ffff, 0)
	block := NewBlock(header, []Transaction{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = block.Hash()
	}
}

// BenchmarkBlock_Validate benchmarks block validation
func BenchmarkBlock_Validate(b *testing.B) {
	// Create valid block with coinbase
	coinbase := Transaction{
		Version: 1,
		Inputs: []TxInput{{
			PreviousOutput: OutPoint{Hash: ZeroHash, Index: 0xffffffff},
			ScriptSig:      []byte("Benchmark coinbase"),
			Sequence:       0xffffffff,
		}},
		Outputs: []TxOutput{{
			Value:        5000000000,
			ScriptPubKey: []byte{0x76, 0xa9, 0x14},
		}},
		LockTime: 0,
	}

	header := NewBlockHeader(1, ZeroHash, ZeroHash, 1640995200, 0x1d00ffff, 0)
	block := NewBlock(header, []Transaction{coinbase})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = block.Validate()
	}
}

// TestBlockHeader_Serialize tests block header serialization edge cases
func TestBlockHeader_Serialize(t *testing.T) {
	tests := []struct {
		name        string
		header      BlockHeader
		shouldError bool
		description string
	}{
		{
			name: "Normal block header",
			header: BlockHeader{
				Version:       1,
				PrevBlockHash: mustParseHash("000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"),
				MerkleRoot:    mustParseHash("4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b"),
				Timestamp:     1231006505,
				Bits:          0x1d00ffff,
				Nonce:         2083236893,
			},
			shouldError: false,
			description: "Genesis block header should serialize correctly",
		},
		{
			name: "Block header with zero hashes",
			header: BlockHeader{
				Version:       1,
				PrevBlockHash: Hash256{},
				MerkleRoot:    Hash256{},
				Timestamp:     1234567890,
				Bits:          0x207fffff,
				Nonce:         0,
			},
			shouldError: false,
			description: "Block header with zero hashes should serialize",
		},
		{
			name: "Block header with maximum values",
			header: BlockHeader{
				Version:       0xffffffff,
				PrevBlockHash: Hash256{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
				MerkleRoot:    Hash256{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
				Timestamp:     0xffffffff,
				Bits:          0xffffffff,
				Nonce:         0xffffffff,
			},
			shouldError: false,
			description: "Block header with maximum values should serialize",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := tt.header.serialize()

			if tt.shouldError {
				if err == nil {
					t.Error("Expected serialization to fail, but it succeeded")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected serialization error: %v", err)
				}

				// Verify data length (80 bytes for block header)
				expectedLen := 80
				if len(data) != expectedLen {
					t.Errorf("Expected serialized length %d, got %d", expectedLen, len(data))
				}
			}

			t.Logf("Serialization result for %s: length=%d, error=%v", tt.description, len(data), err)
		})
	}
}

// TestBlock_Validate_EdgeCases tests block validation edge cases
func TestBlock_Validate_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		block       *Block
		shouldError bool
		description string
	}{
		{
			name: "Block with duplicate transactions",
			block: &Block{
				Header: BlockHeader{
					Version:       1,
					PrevBlockHash: Hash256{},
					MerkleRoot:    Hash256{},
					Timestamp:     1234567890,
					Bits:          0x207fffff,
					Nonce:         1,
				},
				Transactions: []Transaction{
					// Coinbase transaction
					{
						Version: 1,
						Inputs: []TxInput{{
							PreviousOutput: OutPoint{Hash: Hash256{}, Index: 0xffffffff},
							ScriptSig:      []byte{0x01, 0x42},
							Sequence:       0xffffffff,
						}},
						Outputs: []TxOutput{{
							Value:        5000000000,
							ScriptPubKey: []byte{0x76, 0xa9, 0x14},
						}},
						LockTime: 0,
					},
					// Same transaction again (duplicate)
					{
						Version: 1,
						Inputs: []TxInput{{
							PreviousOutput: OutPoint{Hash: Hash256{}, Index: 0xffffffff},
							ScriptSig:      []byte{0x01, 0x42},
							Sequence:       0xffffffff,
						}},
						Outputs: []TxOutput{{
							Value:        5000000000,
							ScriptPubKey: []byte{0x76, 0xa9, 0x14},
						}},
						LockTime: 0,
					},
				},
			},
			shouldError: true,
			description: "Block with duplicate transactions should fail validation",
		},
		{
			name: "Block with multiple coinbase transactions",
			block: &Block{
				Header: BlockHeader{
					Version:       1,
					PrevBlockHash: Hash256{},
					MerkleRoot:    Hash256{},
					Timestamp:     1234567890,
					Bits:          0x207fffff,
					Nonce:         1,
				},
				Transactions: []Transaction{
					// First coinbase transaction
					{
						Version: 1,
						Inputs: []TxInput{{
							PreviousOutput: OutPoint{Hash: Hash256{}, Index: 0xffffffff},
							ScriptSig:      []byte{0x01, 0x42},
							Sequence:       0xffffffff,
						}},
						Outputs: []TxOutput{{
							Value:        5000000000,
							ScriptPubKey: []byte{0x76, 0xa9, 0x14},
						}},
						LockTime: 0,
					},
					// Second coinbase transaction (invalid)
					{
						Version: 1,
						Inputs: []TxInput{{
							PreviousOutput: OutPoint{Hash: Hash256{}, Index: 0xffffffff},
							ScriptSig:      []byte{0x01, 0x43},
							Sequence:       0xffffffff,
						}},
						Outputs: []TxOutput{{
							Value:        2500000000,
							ScriptPubKey: []byte{0x76, 0xa9, 0x14},
						}},
						LockTime: 0,
					},
				},
			},
			shouldError: true,
			description: "Block with multiple coinbase transactions should fail validation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.block.Validate()

			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected validation to fail for %s, but it succeeded", tt.description)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected validation error for %s: %v", tt.description, err)
				}
			}

			t.Logf("Validation result for %s: %v", tt.description, err)
		})
	}
}

func mustParseHash(hashStr string) Hash256 {
	hash, err := NewHash256FromString(hashStr)
	if err != nil {
		panic("Failed to parse hash: " + err.Error())
	}
	return hash
}
