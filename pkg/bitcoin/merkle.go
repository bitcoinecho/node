package bitcoin

import (
	"crypto/sha256"
)

// CalculateMerkleRoot calculates the merkle root for a list of transaction hashes
// This implements the Bitcoin merkle tree algorithm with proper duplication for odd counts
func CalculateMerkleRoot(txHashes []Hash256) Hash256 {
	// Handle edge cases
	if len(txHashes) == 0 {
		return ZeroHash // Return zero hash for empty input
	}

	// For single transaction, merkle root equals the transaction hash
	if len(txHashes) == 1 {
		return txHashes[0]
	}

	// TDD REFACTOR: Implement proper merkle tree algorithm
	// Make a copy to avoid modifying the original slice
	hashes := make([]Hash256, len(txHashes))
	copy(hashes, txHashes)

	// Build the merkle tree level by level
	for len(hashes) > 1 {
		var nextLevel []Hash256

		// Process pairs of hashes
		for i := 0; i < len(hashes); i += 2 {
			var left, right Hash256
			left = hashes[i]

			// If odd number of hashes, duplicate the last one (Bitcoin rule)
			if i+1 < len(hashes) {
				right = hashes[i+1]
			} else {
				right = hashes[i] // Duplicate last hash
			}

			// Combine the pair using double SHA-256
			combined := doubleSHA256(left, right)
			nextLevel = append(nextLevel, combined)
		}

		hashes = nextLevel
	}

	// Return the final root hash
	return hashes[0]
}

// doubleSHA256 performs Bitcoin's double SHA-256 hash on two concatenated hashes
func doubleSHA256(left, right Hash256) Hash256 {
	// Concatenate the two hashes
	combined := make([]byte, 64) // 32 + 32 bytes
	copy(combined[0:32], left[:])
	copy(combined[32:64], right[:])

	// First SHA-256
	first := sha256.Sum256(combined)
	// Second SHA-256
	second := sha256.Sum256(first[:])

	var result Hash256
	copy(result[:], second[:])
	return result
}
