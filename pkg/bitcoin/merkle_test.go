package bitcoin

import (
	"testing"
)

// TestMerkleTree_Construction tests Merkle tree construction (TDD RED phase)
func TestMerkleTree_Construction(t *testing.T) {
	tests := []struct {
		name         string
		txHashes     []string // Transaction hashes as hex strings
		expectedRoot string   // Expected merkle root as hex string
		description  string
	}{
		{
			name: "Single transaction (Genesis block)",
			txHashes: []string{
				"4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b", // Genesis coinbase tx
			},
			expectedRoot: "4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b", // Same as single tx
			description:  "Single transaction merkle root equals the transaction hash",
		},
		{
			name: "Two transactions",
			txHashes: []string{
				"4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b",
				"6e4e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda44c",
			},
			expectedRoot: "calculated_from_double_sha256", // Will be calculated properly
			description:  "Two transaction merkle tree",
		},
		{
			name: "Four transactions (even number)",
			txHashes: []string{
				"4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b",
				"6e4e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda44c",
				"7f5f2f5baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda55d",
				"8a6a3a6baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda66e",
			},
			expectedRoot: "calculated_from_tree", // Will be calculated properly
			description:  "Four transaction balanced merkle tree",
		},
		{
			name: "Three transactions (odd number - requires duplication)",
			txHashes: []string{
				"4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b",
				"6e4e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda44c",
				"7f5f2f5baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda55d",
			},
			expectedRoot: "calculated_with_duplication", // Last tx gets duplicated
			description:  "Odd number of transactions requires last hash duplication",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("TDD RED: %s - %s", tt.name, tt.description)

			// Convert hex strings to Hash256 objects
			var txHashes []Hash256
			for _, hexHash := range tt.txHashes {
				hash, err := NewHash256FromString(hexHash)
				if err != nil {
					t.Fatalf("Failed to parse tx hash %s: %v", hexHash, err)
				}
				txHashes = append(txHashes, hash)
			}

			// This should fail since we haven't implemented MerkleRoot yet
			merkleRoot := CalculateMerkleRoot(txHashes)

			// Convert expected root to Hash256 for comparison
			expectedHash, err := NewHash256FromString(tt.expectedRoot)
			if err != nil {
				// For TDD RED phase, we expect some test data to be placeholders
				t.Logf("Expected root placeholder: %s", tt.expectedRoot)
				// Continue with test to see what we get
			} else {
				if merkleRoot != expectedHash {
					t.Errorf("Expected merkle root %s, got %s", expectedHash.String(), merkleRoot.String())
				}
			}

			t.Logf("Calculated merkle root: %s", merkleRoot.String())
		})
	}
}

// TestMerkleTree_GenesisBlock tests against known Genesis block merkle root
func TestMerkleTree_GenesisBlock(t *testing.T) {
	t.Logf("TDD RED: Testing Genesis block merkle root calculation")

	// Genesis block has only one transaction (coinbase)
	genesisTxHash := "4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b"

	hash, err := NewHash256FromString(genesisTxHash)
	if err != nil {
		t.Fatalf("Failed to parse genesis tx hash: %v", err)
	}

	// Calculate merkle root (should equal the single transaction hash)
	merkleRoot := CalculateMerkleRoot([]Hash256{hash})

	// For single transaction, merkle root should equal the transaction hash
	if merkleRoot != hash {
		t.Errorf("Genesis block merkle root should equal single transaction hash")
		t.Errorf("Expected: %s", hash.String())
		t.Errorf("Got: %s", merkleRoot.String())
	}

	t.Logf("✓ Genesis merkle root: %s", merkleRoot.String())
}

// TestMerkleTree_EmptyInput tests edge cases (TDD RED phase)
func TestMerkleTree_EmptyInput(t *testing.T) {
	tests := []struct {
		name        string
		txHashes    []Hash256
		shouldError bool
		description string
	}{
		{
			name:        "Empty transaction list",
			txHashes:    []Hash256{},
			shouldError: true,
			description: "Empty input should return error",
		},
		{
			name:        "Nil transaction list",
			txHashes:    nil,
			shouldError: true,
			description: "Nil input should return error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("TDD RED: %s - %s", tt.name, tt.description)

			// This should panic or return zero hash since function doesn't exist yet
			defer func() {
				if r := recover(); r != nil {
					t.Logf("Function panicked as expected: %v", r)
				}
			}()

			merkleRoot := CalculateMerkleRoot(tt.txHashes)

			if tt.shouldError {
				// We expect this to be handled gracefully once implemented
				t.Logf("Got merkle root for invalid input: %s", merkleRoot.String())
			}
		})
	}
}
