package bitcoin

import (
	"math/big"
)

// ValidateProofOfWork checks if a block hash meets the difficulty target
// TDD GREEN: Minimal implementation to make tests pass
func ValidateProofOfWork(blockHash Hash256, targetBits uint32) bool {
	// Convert compact target to full target
	target := CompactToBigTarget(targetBits)

	// Convert block hash to big.Int for comparison
	hashInt := new(big.Int)
	hashInt.SetBytes(blockHash[:])

	// Hash must be less than or equal to target
	return hashInt.Cmp(target) <= 0
}

// CompactToBigTarget converts compact target representation to big.Int
// TDD REFACTOR: Enhanced implementation with proper Bitcoin format handling
func CompactToBigTarget(compactBits uint32) *big.Int {
	// Handle zero/invalid input
	if compactBits == 0 {
		return big.NewInt(0)
	}

	// Extract exponent and mantissa from compact representation
	// Bitcoin format: 0xEEMMMMNN where EE is exponent, MMMMNN is mantissa
	exponent := compactBits >> 24
	mantissa := compactBits & 0x00ffffff

	// Handle invalid cases
	if exponent > 32 {
		return big.NewInt(0) // Invalid exponent
	}

	// Handle special cases for small exponents
	if exponent <= 3 {
		// For small exponents, mantissa is shifted right
		target := big.NewInt(int64(mantissa))
		if exponent < 3 {
			target.Rsh(target, uint((3-exponent)*8))
		}
		return target
	}

	// Normal case: mantissa * 256^(exponent-3)
	target := big.NewInt(int64(mantissa))
	target.Lsh(target, uint((exponent-3)*8))

	return target
}

// BigTargetToHash256 converts big.Int target to Hash256 for comparison
// TDD GREEN: Basic conversion for testing
func BigTargetToHash256(target *big.Int) Hash256 {
	var hash Hash256

	// Convert big.Int to bytes (big-endian)
	targetBytes := target.Bytes()

	// Copy to hash (right-aligned, like Bitcoin does)
	if len(targetBytes) <= 32 {
		copy(hash[32-len(targetBytes):], targetBytes)
	}

	return hash
}

// AdjustDifficulty calculates new difficulty target based on time taken
// TDD REFACTOR: Complete Bitcoin difficulty adjustment algorithm
func AdjustDifficulty(currentTargetBits, actualTimeSeconds uint32) uint32 {
	// Bitcoin difficulty adjustment constants
	const targetTimespan = 14 * 24 * 60 * 60 // 2 weeks in seconds
	const maxAdjustment = 4                  // Max 4x adjustment up or down

	// Handle edge case
	if actualTimeSeconds == 0 {
		return currentTargetBits
	}

	// If time is exactly 2 weeks, no adjustment needed
	if actualTimeSeconds == targetTimespan {
		return currentTargetBits
	}

	// Convert current target to big.Int for calculation
	currentTarget := CompactToBigTarget(currentTargetBits)

	// Calculate adjustment ratio: actualTime / targetTime
	// newTarget = currentTarget * actualTime / targetTime
	actualTime := big.NewInt(int64(actualTimeSeconds))
	targetTime := big.NewInt(targetTimespan)

	// Apply adjustment limits (max 4x up or down)
	maxTime := big.NewInt(targetTimespan * maxAdjustment)
	minTime := big.NewInt(targetTimespan / maxAdjustment)

	if actualTime.Cmp(maxTime) > 0 {
		actualTime = maxTime
	}
	if actualTime.Cmp(minTime) < 0 {
		actualTime = minTime
	}

	// Calculate new target: currentTarget * actualTime / targetTime
	newTarget := new(big.Int)
	newTarget.Mul(currentTarget, actualTime)
	newTarget.Div(newTarget, targetTime)

	// Convert back to compact representation
	return BigTargetToCompact(newTarget)
}

// BigTargetToCompact converts big.Int target back to compact representation
// TDD REFACTOR: Added for complete difficulty adjustment
func BigTargetToCompact(target *big.Int) uint32 {
	if target.Sign() <= 0 {
		return 0
	}

	// Convert to bytes
	targetBytes := target.Bytes()
	if len(targetBytes) == 0 {
		return 0
	}

	// Find the exponent (number of bytes)
	exponent := len(targetBytes)

	// Extract mantissa (first 3 bytes, big-endian)
	var mantissa uint32
	if exponent >= 3 {
		mantissa = uint32(targetBytes[0])<<16 | uint32(targetBytes[1])<<8 | uint32(targetBytes[2])
	} else if exponent == 2 {
		mantissa = uint32(targetBytes[0])<<16 | uint32(targetBytes[1])<<8
	} else {
		mantissa = uint32(targetBytes[0]) << 16
	}

	// Handle the case where the high bit is set (negative in Bitcoin's format)
	if mantissa&0x800000 != 0 {
		mantissa >>= 8
		exponent++
	}

	// Combine exponent and mantissa
	// Ensure exponent is within valid range before conversion
	if exponent < 0 || exponent > 255 {
		return 0
	}
	compact := uint32(exponent)<<24 | (mantissa & 0x00ffffff)

	return compact
}
