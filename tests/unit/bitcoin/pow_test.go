package bitcoin_test

import (
	"testing"
	"bitcoinecho.org/node/pkg/bitcoin"
)

// TestProofOfWork_TargetValidation tests PoW target validation (TDD RED phase)
func TestProofOfWork_TargetValidation(t *testing.T) {
	tests := []struct {
		name          string
		blockHash     string // Block hash as hex string
		targetBits    uint32 // Compact target representation
		expectedValid bool
		description   string
	}{
		{
			name:          "Genesis block valid PoW",
			blockHash:     "000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f", // Genesis block hash
			targetBits:    0x1d00ffff,                                                           // Genesis target
			expectedValid: true,
			description:   "Genesis block should pass PoW validation",
		},
		{
			name:          "Hash above target (invalid)",
			blockHash:     "100000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f", // Higher hash
			targetBits:    0x1d00ffff,                                                           // Same target
			expectedValid: false,
			description:   "Hash above target should fail validation",
		},
		{
			name:          "Maximum difficulty target",
			blockHash:     "0000000000000000000000000000000000000000000000000000000000000001", // Very low hash
			targetBits:    0x207fffff,                                                           // Maximum difficulty
			expectedValid: true,
			description:   "Hash below maximum difficulty should pass",
		},
		{
			name:          "Minimum difficulty target",
			blockHash:     "00000000ffff0000000000000000000000000000000000000000000000000000", // Medium hash
			targetBits:    0x1d00ffff,                                                           // Minimum difficulty (genesis)
			expectedValid: true,
			description:   "Hash below minimum difficulty should pass",
		},
		{
			name:          "Border case - exactly at target",
			blockHash:     "00000000ffff0000000000000000000000000000000000000000000000000000",
			targetBits:    0x1d00ffff,
			expectedValid: true, // Exactly at target should be valid
			description:   "Hash exactly at target should pass validation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("TDD RED: %s - %s", tt.name, tt.description)

			// Parse the block hash
			blockHash, err := bitcoin.NewHash256FromString(tt.blockHash)
			if err != nil {
				t.Fatalf("Failed to parse block hash: %v", err)
			}

			// This should fail since we haven't implemented PoW validation yet
			isValid := bitcoin.ValidateProofOfWork(blockHash, tt.targetBits)

			if isValid != tt.expectedValid {
				t.Errorf("Expected PoW validation %v, got %v", tt.expectedValid, isValid)
			}

			t.Logf("PoW validation result: %v (expected: %v)", isValid, tt.expectedValid)
		})
	}
}

// TestProofOfWork_CompactTargetConversion tests compact target conversion (TDD RED phase)
func TestProofOfWork_CompactTargetConversion(t *testing.T) {
	tests := []struct {
		name         string
		compactBits  uint32
		expectedTarget string // Expected full target as hex
		description  string
	}{
		{
			name:         "Genesis block target",
			compactBits:  0x1d00ffff,
			expectedTarget: "00000000ffff0000000000000000000000000000000000000000000000000000", // Genesis target
			description:  "Genesis block compact target conversion",
		},
		{
			name:         "Maximum difficulty",
			compactBits:  0x207fffff,
			expectedTarget: "7fffff0000000000000000000000000000000000000000000000000000000000", // Max difficulty (correct format)
			description:  "Maximum difficulty target conversion",
		},
		{
			name:         "Early Bitcoin target",
			compactBits:  0x1d00ffff,
			expectedTarget: "00000000ffff0000000000000000000000000000000000000000000000000000",
			description:  "Early Bitcoin difficulty target",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("TDD RED: %s - %s", tt.name, tt.description)

			// This should fail since we haven't implemented target conversion yet
			target := bitcoin.CompactToBigTarget(tt.compactBits)

			expectedTarget, err := bitcoin.NewHash256FromString(tt.expectedTarget)
			if err != nil {
				t.Fatalf("Failed to parse expected target: %v", err)
			}

			// Convert target to Hash256 for comparison
			targetHash := bitcoin.BigTargetToHash256(target)

			if targetHash != expectedTarget {
				t.Errorf("Expected target %s, got %s", expectedTarget.String(), targetHash.String())
			}

			t.Logf("Compact 0x%08x -> Target: %s", tt.compactBits, targetHash.String())
		})
	}
}

// TestProofOfWork_DifficultyAdjustment tests difficulty adjustment algorithm (TDD RED phase)
func TestProofOfWork_DifficultyAdjustment(t *testing.T) {
	tests := []struct {
		name           string
		currentTarget  uint32
		actualTime     uint32 // Time taken for last 2016 blocks (seconds)
		expectedTarget uint32
		description    string
	}{
		{
			name:           "No adjustment needed",
			currentTarget:  0x1d00ffff,
			actualTime:     1209600, // Exactly 2 weeks
			expectedTarget: 0x1d00ffff,
			description:    "Perfect timing should not adjust difficulty",
		},
		{
			name:           "Increase difficulty (blocks too fast)",
			currentTarget:  0x1d00ffff,
			actualTime:     604800, // 1 week (half time)
			expectedTarget: 0x1c7fff80, // Should increase difficulty (calculated value)
			description:    "Fast blocks should increase difficulty",
		},
		{
			name:           "Decrease difficulty (blocks too slow)",
			currentTarget:  0x1d00ffff,
			actualTime:     2419200, // 4 weeks (double time)
			expectedTarget: 0x1d01fffe, // Should decrease difficulty (calculated value)
			description:    "Slow blocks should decrease difficulty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("TDD RED: %s - %s", tt.name, tt.description)

			// This should fail since we haven't implemented difficulty adjustment yet
			newTarget := bitcoin.AdjustDifficulty(tt.currentTarget, tt.actualTime)

			if newTarget != tt.expectedTarget {
				t.Errorf("Expected new target 0x%08x, got 0x%08x", tt.expectedTarget, newTarget)
			}

			t.Logf("Difficulty adjustment: 0x%08x -> 0x%08x (time: %d seconds)",
				tt.currentTarget, newTarget, tt.actualTime)
		})
	}
}

// TestProofOfWork_GenesisBlock tests against known Genesis block (TDD RED phase)
func TestProofOfWork_GenesisBlock(t *testing.T) {
	t.Logf("TDD RED: Testing Genesis block PoW validation")

	// Known Genesis block values
	genesisHash := "000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"
	genesisTarget := uint32(0x1d00ffff)

	blockHash, err := bitcoin.NewHash256FromString(genesisHash)
	if err != nil {
		t.Fatalf("Failed to parse Genesis hash: %v", err)
	}

	// Validate Genesis block PoW
	isValid := bitcoin.ValidateProofOfWork(blockHash, genesisTarget)

	if !isValid {
		t.Errorf("Genesis block should pass PoW validation")
	}

	t.Logf("✓ Genesis block PoW validation: %v", isValid)
	t.Logf("Genesis hash: %s", blockHash.String())
	t.Logf("Genesis target: 0x%08x", genesisTarget)
}

// TestProofOfWork_EdgeCases tests edge cases and error conditions (TDD RED phase)
func TestProofOfWork_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		blockHash   string
		targetBits  uint32
		shouldPanic bool
		description string
	}{
		{
			name:        "Invalid target bits (overflow)",
			blockHash:   "000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f",
			targetBits:  0xff000000, // Invalid exponent
			shouldPanic: false,       // Should handle gracefully
			description: "Invalid target bits should be handled gracefully",
		},
		{
			name:        "Zero target bits",
			blockHash:   "000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f",
			targetBits:  0x00000000,
			shouldPanic: false,
			description: "Zero target should be handled gracefully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("TDD RED: %s - %s", tt.name, tt.description)

			blockHash, err := bitcoin.NewHash256FromString(tt.blockHash)
			if err != nil {
				t.Fatalf("Failed to parse block hash: %v", err)
			}

			if tt.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("Expected function to panic, but it didn't")
					}
				}()
			}

			// This should handle edge cases gracefully once implemented
			isValid := bitcoin.ValidateProofOfWork(blockHash, tt.targetBits)
			t.Logf("Edge case result: %v", isValid)
		})
	}
}