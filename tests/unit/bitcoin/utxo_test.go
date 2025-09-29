package bitcoin_test

import (
	"bitcoinecho.org/node/pkg/bitcoin"
	"testing"
)

// TestUTXO_Creation tests UTXO creation and basic operations (TDD RED phase)
func TestUTXO_Creation(t *testing.T) {
	tests := []struct {
		name         string
		txHash       string
		outputIndex  uint32
		amount       uint64
		scriptPubKey []byte
		description  string
	}{
		{
			name:        "Create standard P2PKH UTXO",
			txHash:      "abcd1234567890abcdef1234567890abcdef1234567890abcdef1234567890ab",
			outputIndex: 0,
			amount:      5000000000, // 50 BTC in satoshis
			scriptPubKey: []byte{
				0x76, 0xa9, 0x14, // OP_DUP OP_HASH160 <20 bytes>
				0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0, 0x12, 0x34,
				0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0, 0x12, 0x34, 0x56, 0x78,
				0x88, 0xac, // OP_EQUALVERIFY OP_CHECKSIG
			},
			description: "Standard P2PKH UTXO should be created correctly",
		},
		{
			name:        "Create coinbase UTXO",
			txHash:      "0000000000000000000000000000000000000000000000000000000000000000",
			outputIndex: 0,
			amount:      5000000000, // Genesis block coinbase reward
			scriptPubKey: []byte{
				0x41, // OP_PUSHDATA 65 bytes
				0x04, 0x67, 0x8a, 0xfd, 0xb0, 0xfe, 0x55, 0x48, 0x27, 0x19,
				0x67, 0xf1, 0xa6, 0x71, 0x30, 0xb7, 0x10, 0x5c, 0xd6, 0xa8,
				0x28, 0xe0, 0x39, 0x09, 0xa6, 0x79, 0x62, 0xe0, 0xea, 0x1f,
				0x61, 0xde, 0xb6, 0x49, 0xf6, 0xbc, 0x3f, 0x4c, 0xef, 0x38,
				0xc4, 0xf3, 0x55, 0x04, 0xe5, 0x1e, 0xc1, 0x12, 0xde, 0x5c,
				0x38, 0x4d, 0xf7, 0xba, 0x0b, 0x8d, 0x57, 0x8a, 0x4c, 0x70,
				0x2b, 0x6b, 0xf1, 0x1d, 0x5f,
				0xac, // OP_CHECKSIG
			},
			description: "Genesis block coinbase UTXO should be created",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("TDD RED: %s - %s", tt.name, tt.description)

			// This should fail since we haven't implemented UTXO creation yet
			txHash, err := bitcoin.NewHash256FromString(tt.txHash)
			if err != nil {
				t.Fatalf("Failed to parse transaction hash: %v", err)
			}

			utxo := bitcoin.NewUTXO(txHash, tt.outputIndex, tt.amount, tt.scriptPubKey)

			if utxo.TxHash() != txHash {
				t.Errorf("Expected tx hash %s, got %s", txHash.String(), utxo.TxHash().String())
			}

			if utxo.OutputIndex() != tt.outputIndex {
				t.Errorf("Expected output index %d, got %d", tt.outputIndex, utxo.OutputIndex())
			}

			if utxo.Amount() != tt.amount {
				t.Errorf("Expected amount %d, got %d", tt.amount, utxo.Amount())
			}

			t.Logf("UTXO created: %s:%d, amount: %d satoshis",
				utxo.TxHash().String()[:8], utxo.OutputIndex(), utxo.Amount())
		})
	}
}

// TestUTXOSet_AddRemove tests UTXO set operations (TDD RED phase)
func TestUTXOSet_AddRemove(t *testing.T) {
	tests := []struct {
		name        string
		operations  []string // "add" or "remove"
		txHashes    []string
		outputIndex []uint32
		amounts     []uint64
		description string
	}{
		{
			name:        "Add single UTXO",
			operations:  []string{"add"},
			txHashes:    []string{"abcd1234567890abcdef1234567890abcdef1234567890abcdef1234567890ab"},
			outputIndex: []uint32{0},
			amounts:     []uint64{5000000000},
			description: "Adding single UTXO should increase set size",
		},
		{
			name:       "Add multiple UTXOs",
			operations: []string{"add", "add", "add"},
			txHashes: []string{
				"abcd1234567890abcdef1234567890abcdef1234567890abcdef1234567890ab",
				"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
				"567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234",
			},
			outputIndex: []uint32{0, 1, 0},
			amounts:     []uint64{5000000000, 2500000000, 1000000000},
			description: "Adding multiple UTXOs should track all correctly",
		},
		{
			name:       "Add and remove UTXO",
			operations: []string{"add", "remove"},
			txHashes: []string{
				"abcd1234567890abcdef1234567890abcdef1234567890abcdef1234567890ab",
				"abcd1234567890abcdef1234567890abcdef1234567890abcdef1234567890ab",
			},
			outputIndex: []uint32{0, 0},
			amounts:     []uint64{5000000000, 5000000000},
			description: "Adding then removing UTXO should result in empty set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("TDD RED: %s - %s", tt.name, tt.description)

			// This should fail since we haven't implemented UTXO set yet
			utxoSet := bitcoin.NewUTXOSet()

			for i, op := range tt.operations {
				txHash, err := bitcoin.NewHash256FromString(tt.txHashes[i])
				if err != nil {
					t.Fatalf("Failed to parse transaction hash: %v", err)
				}

				if op == "add" {
					utxo := bitcoin.NewUTXO(txHash, tt.outputIndex[i], tt.amounts[i], []byte{0x76, 0xa9})
					utxoSet.Add(utxo)
				} else if op == "remove" {
					utxoSet.Remove(txHash, tt.outputIndex[i])
				}
			}

			expectedSize := 0
			for _, op := range tt.operations {
				if op == "add" {
					expectedSize++
				} else if op == "remove" {
					expectedSize--
				}
			}

			if utxoSet.Size() != expectedSize {
				t.Errorf("Expected UTXO set size %d, got %d", expectedSize, utxoSet.Size())
			}

			t.Logf("UTXO set operations completed, final size: %d", utxoSet.Size())
		})
	}
}

// TestUTXOSet_Find tests UTXO lookup operations (TDD RED phase)
func TestUTXOSet_Find(t *testing.T) {
	tests := []struct {
		name           string
		setupUTXOs     int
		searchTxHash   string
		searchIndex    uint32
		shouldFind     bool
		expectedAmount uint64
		description    string
	}{
		{
			name:           "Find existing UTXO",
			setupUTXOs:     3,
			searchTxHash:   "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			searchIndex:    1,
			shouldFind:     true,
			expectedAmount: 2500000000,
			description:    "Should find existing UTXO in set",
		},
		{
			name:           "UTXO not found - wrong hash",
			setupUTXOs:     3,
			searchTxHash:   "9999999999999999999999999999999999999999999999999999999999999999",
			searchIndex:    0,
			shouldFind:     false,
			expectedAmount: 0,
			description:    "Should not find non-existent UTXO",
		},
		{
			name:           "UTXO not found - wrong index",
			setupUTXOs:     3,
			searchTxHash:   "abcd1234567890abcdef1234567890abcdef1234567890abcdef1234567890ab",
			searchIndex:    99,
			shouldFind:     false,
			expectedAmount: 0,
			description:    "Should not find UTXO with wrong output index",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("TDD RED: %s - %s", tt.name, tt.description)

			// Setup UTXO set with test data
			utxoSet := bitcoin.NewUTXOSet()

			// Add some test UTXOs
			testHashes := []string{
				"abcd1234567890abcdef1234567890abcdef1234567890abcdef1234567890ab",
				"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
				"567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234",
			}
			testAmounts := []uint64{5000000000, 2500000000, 1000000000}

			for i := 0; i < tt.setupUTXOs && i < len(testHashes); i++ {
				txHash, _ := bitcoin.NewHash256FromString(testHashes[i])
				utxo := bitcoin.NewUTXO(txHash, uint32(i), testAmounts[i], []byte{0x76, 0xa9})
				utxoSet.Add(utxo)
			}

			// Search for UTXO
			searchHash, err := bitcoin.NewHash256FromString(tt.searchTxHash)
			if err != nil {
				t.Fatalf("Failed to parse search hash: %v", err)
			}

			utxo, found := utxoSet.Find(searchHash, tt.searchIndex)

			if found != tt.shouldFind {
				t.Errorf("Expected found=%v, got found=%v", tt.shouldFind, found)
			}

			if tt.shouldFind {
				if utxo.Amount() != tt.expectedAmount {
					t.Errorf("Expected amount %d, got %d", tt.expectedAmount, utxo.Amount())
				}
			}

			t.Logf("UTXO search for %s:%d - found: %v",
				searchHash.String()[:8], tt.searchIndex, found)
		})
	}
}

// TestUTXOSet_Validation tests UTXO validation operations (TDD RED phase)
func TestUTXOSet_Validation(t *testing.T) {
	tests := []struct {
		name        string
		inputTxHash string
		inputIndex  uint32
		amount      uint64
		isValid     bool
		description string
	}{
		{
			name:        "Valid UTXO spend",
			inputTxHash: "abcd1234567890abcdef1234567890abcdef1234567890abcdef1234567890ab",
			inputIndex:  0,
			amount:      5000000000,
			isValid:     true,
			description: "Spending existing UTXO should be valid",
		},
		{
			name:        "Invalid UTXO spend - not found",
			inputTxHash: "9999999999999999999999999999999999999999999999999999999999999999",
			inputIndex:  0,
			amount:      5000000000,
			isValid:     false,
			description: "Spending non-existent UTXO should be invalid",
		},
		{
			name:        "Invalid UTXO spend - double spend",
			inputTxHash: "abcd1234567890abcdef1234567890abcdef1234567890abcdef1234567890ab",
			inputIndex:  0,
			amount:      5000000000,
			isValid:     false,
			description: "Double spending UTXO should be invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("TDD RED: %s - %s", tt.name, tt.description)

			// Setup UTXO set
			utxoSet := bitcoin.NewUTXOSet()

			// Add a test UTXO
			txHash, _ := bitcoin.NewHash256FromString("abcd1234567890abcdef1234567890abcdef1234567890abcdef1234567890ab")
			utxo := bitcoin.NewUTXO(txHash, 0, 5000000000, []byte{0x76, 0xa9})
			utxoSet.Add(utxo)

			// For double spend test, remove the UTXO first
			if tt.name == "Invalid UTXO spend - double spend" {
				utxoSet.Remove(txHash, 0)
			}

			// Validate the spend
			inputHash, _ := bitcoin.NewHash256FromString(tt.inputTxHash)
			isValid := utxoSet.ValidateSpend(inputHash, tt.inputIndex, tt.amount)

			if isValid != tt.isValid {
				t.Errorf("Expected validation result %v, got %v", tt.isValid, isValid)
			}

			t.Logf("UTXO spend validation for %s:%d - valid: %v",
				inputHash.String()[:8], tt.inputIndex, isValid)
		})
	}
}

// TestUTXOSet_TotalValue tests UTXO set value calculations (TDD RED phase)
func TestUTXOSet_TotalValue(t *testing.T) {
	tests := []struct {
		name          string
		utxoAmounts   []uint64
		expectedTotal uint64
		description   string
	}{
		{
			name:          "Empty UTXO set",
			utxoAmounts:   []uint64{},
			expectedTotal: 0,
			description:   "Empty set should have zero total value",
		},
		{
			name:          "Single UTXO",
			utxoAmounts:   []uint64{5000000000},
			expectedTotal: 5000000000,
			description:   "Single UTXO set total should equal the UTXO amount",
		},
		{
			name:          "Multiple UTXOs",
			utxoAmounts:   []uint64{5000000000, 2500000000, 1000000000, 500000000},
			expectedTotal: 9000000000,
			description:   "Multiple UTXO set should sum all amounts correctly",
		},
		{
			name:          "Genesis block coinbase",
			utxoAmounts:   []uint64{5000000000}, // 50 BTC original coinbase
			expectedTotal: 5000000000,
			description:   "Genesis coinbase should calculate correctly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("TDD RED: %s - %s", tt.name, tt.description)

			// This should fail since we haven't implemented UTXO set value calculation yet
			utxoSet := bitcoin.NewUTXOSet()

			// Add UTXOs with specified amounts
			for i, amount := range tt.utxoAmounts {
				txHashStr := "abcd1234567890abcdef1234567890abcdef1234567890abcdef1234567890" +
					string(rune('a'+i%10)) + string(rune('b'+i%10))
				txHash, _ := bitcoin.NewHash256FromString(txHashStr)
				utxo := bitcoin.NewUTXO(txHash, uint32(i), amount, []byte{0x76, 0xa9})
				utxoSet.Add(utxo)
			}

			totalValue := utxoSet.TotalValue()

			if totalValue != tt.expectedTotal {
				t.Errorf("Expected total value %d, got %d", tt.expectedTotal, totalValue)
			}

			t.Logf("UTXO set total value: %d satoshis (%.8f BTC)",
				totalValue, float64(totalValue)/100000000)
		})
	}
}
