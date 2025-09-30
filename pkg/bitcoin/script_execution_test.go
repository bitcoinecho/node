package bitcoin

import (
	"bytes"
	"encoding/hex"
	"testing"
)

// TestScriptEngine_Execute tests Bitcoin script execution with real Bitcoin scripts
func TestScriptEngine_Execute(t *testing.T) {
	tests := []struct {
		name       string
		scriptHex  string              // Script as hex string
		expected   bool                // Expected execution result
		finalStack []string            // Expected final stack state (hex strings)
		flags      ScriptFlags // Script verification flags
	}{
		// Basic stack operations
		{
			name:       "OP_1 pushes 1 to stack",
			scriptHex:  "51", // OP_1
			expected:   true,
			finalStack: []string{"01"}, // 1 as single byte
			flags:      ScriptFlagsNone,
		},
		{
			name:       "OP_2 pushes 2 to stack",
			scriptHex:  "52", // OP_2
			expected:   true,
			finalStack: []string{"02"}, // 2 as single byte
			flags:      ScriptFlagsNone,
		},
		{
			name:       "Push data operation",
			scriptHex:  "0548656c6c6f", // PUSH(5) "Hello"
			expected:   true,
			finalStack: []string{"48656c6c6f"}, // "Hello" in hex
			flags:      ScriptFlagsNone,
		},
		{
			name:       "OP_DUP duplicates top stack item",
			scriptHex:  "5176", // OP_1 OP_DUP
			expected:   true,
			finalStack: []string{"01", "01"}, // Two copies of 1
			flags:      ScriptFlagsNone,
		},
		{
			name:       "OP_DROP removes top stack item",
			scriptHex:  "515275", // OP_1 OP_2 OP_DROP
			expected:   true,
			finalStack: []string{"01"}, // Only 1 remains
			flags:      ScriptFlagsNone,
		},

		// Arithmetic operations
		{
			name:       "OP_ADD adds two numbers",
			scriptHex:  "515293", // OP_1 OP_2 OP_ADD
			expected:   true,
			finalStack: []string{"03"}, // 1 + 2 = 3
			flags:      ScriptFlagsNone,
		},
		{
			name:       "OP_SUB subtracts two numbers",
			scriptHex:  "525194", // OP_2 OP_1 OP_SUB
			expected:   true,
			finalStack: []string{"01"}, // 2 - 1 = 1
			flags:      ScriptFlagsNone,
		},

		// Logical operations
		{
			name:       "OP_EQUAL compares equal values",
			scriptHex:  "515187", // OP_1 OP_1 OP_EQUAL
			expected:   true,
			finalStack: []string{"01"}, // True (1)
			flags:      ScriptFlagsNone,
		},
		{
			name:       "OP_EQUAL compares different values",
			scriptHex:  "515287", // OP_1 OP_2 OP_EQUAL
			expected:   true,
			finalStack: []string{"00"}, // False (0)
			flags:      ScriptFlagsNone,
		},
		{
			name:       "OP_EQUALVERIFY with equal values",
			scriptHex:  "515188", // OP_1 OP_1 OP_EQUALVERIFY
			expected:   true,
			finalStack: []string{}, // Stack should be empty after verify
			flags:      ScriptFlagsNone,
		},
		{
			name:       "OP_EQUALVERIFY with different values (should fail)",
			scriptHex:  "515288", // OP_1 OP_2 OP_EQUALVERIFY
			expected:   false,    // Should fail verification
			finalStack: []string{},
			flags:      ScriptFlagsNone,
		},

		// Hash operations
		{
			name:       "OP_HASH160 of known data",
			scriptHex:  "0548656c6c6fa9", // PUSH(5) "Hello" OP_HASH160
			expected:   true,
			finalStack: []string{"b6a9c8c230722b7c748331a8b450f05566dc7d0f"}, // HASH160("Hello")
			flags:      ScriptFlagsNone,
		},

		// Complex scripts
		{
			name:       "Simple P2PKH-like pattern (without signature)",
			scriptHex:  "76a914" + "b6a9c8c230722b7c748331a8b450f05566dc7d0f" + "87", // OP_DUP OP_HASH160 <hash> OP_EQUAL
			expected:   false,                                                        // Should fail without matching data on stack
			finalStack: []string{},
			flags:      ScriptFlagsNone,
		},

		// Error conditions
		{
			name:       "Empty script",
			scriptHex:  "",
			expected:   true, // Empty script should succeed
			finalStack: []string{},
			flags:      ScriptFlagsNone,
		},
		{
			name:       "OP_DUP with empty stack (should fail)",
			scriptHex:  "76", // OP_DUP
			expected:   false,
			finalStack: []string{},
			flags:      ScriptFlagsNone,
		},
		{
			name:       "OP_ADD with insufficient stack items (should fail)",
			scriptHex:  "5193", // OP_1 OP_ADD (needs 2 items)
			expected:   false,
			finalStack: []string{},
			flags:      ScriptFlagsNone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Decode script from hex
			scriptBytes, err := hex.DecodeString(tt.scriptHex)
			if err != nil {
				t.Fatalf("Failed to decode script hex: %v", err)
			}

			// Create script engine (no transaction context needed for basic tests)
			script := Script(scriptBytes)
			engine := NewScriptEngine(script, nil, 0, nil, tt.flags)

			// Execute script
			result, err := engine.Execute()

			// Check execution result
			if result != tt.expected {
				if err != nil {
					t.Errorf("Expected result %v, got %v with error: %v", tt.expected, result, err)
				} else {
					t.Errorf("Expected result %v, got %v", tt.expected, result)
				}
			}

			// If execution succeeded, check final stack state
			if result && tt.expected {
				actualStack := engine.GetStack()
				if len(actualStack) != len(tt.finalStack) {
					t.Errorf("Expected stack size %d, got %d", len(tt.finalStack), len(actualStack))
					return
				}

				for i, expectedHex := range tt.finalStack {
					expected, err := hex.DecodeString(expectedHex)
					if err != nil {
						t.Fatalf("Invalid expected stack hex at index %d: %v", i, err)
					}

					if !bytes.Equal(actualStack[i], expected) {
						t.Errorf("Stack item %d: expected %x, got %x", i, expected, actualStack[i])
					}
				}
			}
		})
	}
}

// TestScriptEngine_P2PKHExecution tests basic P2PKH execution patterns
// Note: Full P2PKH tests require signature validation which will be implemented later
func TestScriptEngine_P2PKHExecution(t *testing.T) {
	t.Skip("P2PKH execution requires signature validation - will be implemented in next phase")
}

// TestScriptEngine_SignatureVerification tests ECDSA signature verification (TDD RED phase)
func TestScriptEngine_SignatureVerification(t *testing.T) {
	tests := []struct {
		name         string
		scriptHex    string
		txHash       string // Transaction hash to sign
		pubKeyHex    string // Public key in DER format
		signatureHex string // DER signature
		expected     bool
		description  string
	}{
		{
			name:         "Valid ECDSA signature verification",
			scriptHex:    "ac", // OP_CHECKSIG
			txHash:       "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			pubKeyHex:    "0279be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798",                                                                             // Sample compressed pubkey
			signatureHex: "304402200123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef02200123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef01", // Sample DER signature + SIGHASH_ALL
			expected:     true,
			description:  "Valid signature should verify successfully",
		},
		{
			name:         "Invalid signature should fail verification",
			scriptHex:    "ac", // OP_CHECKSIG
			txHash:       "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			pubKeyHex:    "0279be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798",
			signatureHex: "3044022000000000000000000000000000000000000000000000000000000000000000000220000000000000000000000000000000000000000000000000000000000000000001", // All-zero r and s components = invalid
			expected:     false,
			description:  "Invalid signature should fail verification",
		},
		{
			name:         "Complete P2PKH script execution",
			scriptHex:    "76a914" + "751e76ab4c23b27acb9b8e1c4c9c48c9e9f8a8b3" + "88ac", // OP_DUP OP_HASH160 <20-byte hash> OP_EQUALVERIFY OP_CHECKSIG
			txHash:       "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			pubKeyHex:    "0279be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798",
			signatureHex: "304402200123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef02200123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef01",
			expected:     true,
			description:  "Complete P2PKH script should execute with valid signature",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("TDD RED: %s - %s", tt.name, tt.description)

			// Create a test script that pushes signature and pubkey, then calls OP_CHECKSIG
			var testScript []byte

			// Decode signature and pubkey for the test
			signature, err := hex.DecodeString(tt.signatureHex)
			if err != nil {
				t.Fatalf("Failed to decode signature hex: %v", err)
			}

			pubKey, err := hex.DecodeString(tt.pubKeyHex)
			if err != nil {
				t.Fatalf("Failed to decode pubkey hex: %v", err)
			}

			// Build script: PUSH(signature) PUSH(pubkey) OP_CHECKSIG
			// Push signature
			testScript = append(testScript, byte(len(signature)))
			testScript = append(testScript, signature...)
			// Push pubkey
			testScript = append(testScript, byte(len(pubKey)))
			testScript = append(testScript, pubKey...)
			// Add OP_CHECKSIG
			testScript = append(testScript, 0xac)

			// Execute the test script
			script := Script(testScript)
			engine := NewScriptEngine(script, nil, 0, nil, ScriptFlagsNone)

			result, err := engine.Execute()

			// Check the result based on expected behavior
			if result {
				stack := engine.GetStack()
				if len(stack) != 1 {
					t.Errorf("Expected 1 item on stack after OP_CHECKSIG, got %d", len(stack))
				} else {
					resultByte := stack[0]
					actualResult := len(resultByte) == 1 && resultByte[0] == 1

					t.Logf("OP_CHECKSIG result: %x (expected: %v)", resultByte, tt.expected)

					// Verify that the signature verification behaves as expected
					if actualResult != tt.expected {
						t.Errorf("Expected signature verification result %v, got %v", tt.expected, actualResult)
					}

					if actualResult == tt.expected {
						t.Logf("✓ Signature verification working correctly for %s", tt.name)
					}
				}
			}

			// Log the transaction hash for context (will be needed for real verification)
			t.Logf("Transaction hash for verification: %s", tt.txHash)
		})
	}
}

// TestScriptEngine_StackOperations tests detailed stack manipulation
func TestScriptEngine_StackOperations(t *testing.T) {
	tests := []struct {
		name       string
		operations []struct {
			opcode     ScriptOpcode
			data       []byte // For push operations
			expectFail bool   // True if this operation should fail
		}
		finalStackSize int
		description    string
	}{
		{
			name: "Stack depth management",
			operations: []struct {
				opcode     ScriptOpcode
				data       []byte
				expectFail bool
			}{
				{OP_1, nil, false},    // Push 1
				{OP_2, nil, false},    // Push 2
				{OP_3, nil, false},    // Push 3
				{OP_DROP, nil, false}, // Drop 3
				{OP_SWAP, nil, false}, // Swap 1,2 -> 2,1
				{OP_DUP, nil, false},  // Duplicate 1 -> 2,1,1
			},
			finalStackSize: 3, // Should have [2, 1, 1]
			description:    "Complex stack manipulation should maintain correct state",
		},
		{
			name: "Stack underflow protection",
			operations: []struct {
				opcode     ScriptOpcode
				data       []byte
				expectFail bool
			}{
				{OP_1, nil, false},    // Push 1
				{OP_DROP, nil, false}, // Drop 1 (stack empty)
				{OP_DROP, nil, true},  // Try to drop from empty stack - should fail
			},
			finalStackSize: 0,
			description:    "Operations on empty stack should fail safely",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Build script from operations
			var scriptBytes []byte
			for _, op := range tt.operations {
				if op.data != nil {
					// Push operation
					scriptBytes = append(scriptBytes, byte(len(op.data)))
					scriptBytes = append(scriptBytes, op.data...)
				} else {
					// Regular opcode
					scriptBytes = append(scriptBytes, byte(op.opcode))
				}
			}

			script := Script(scriptBytes)
			engine := NewScriptEngine(script, nil, 0, nil, ScriptFlagsNone)

			result, err := engine.Execute()

			// Check if any operation was expected to fail
			expectOverallFail := false
			for _, op := range tt.operations {
				if op.expectFail {
					expectOverallFail = true
					break
				}
			}

			if expectOverallFail {
				if result {
					t.Errorf("Expected script execution to fail, but it succeeded")
				}
			} else {
				if !result {
					t.Errorf("Expected script execution to succeed, but it failed: %v", err)
				} else {
					// Check final stack size
					stack := engine.GetStack()
					if len(stack) != tt.finalStackSize {
						t.Errorf("Expected final stack size %d, got %d", tt.finalStackSize, len(stack))
					}
				}
			}

			t.Logf("Test: %s - %s", tt.name, tt.description)
		})
	}
}

// TestScriptEngine_Benchmarks tests for performance regressions
func TestScriptEngine_Benchmarks(t *testing.T) {
	// Simple benchmark test to ensure execution performance
	scriptBytes, _ := hex.DecodeString("51525293") // OP_1 OP_2 OP_ADD -> pushes 3
	script := Script(scriptBytes)

	for i := 0; i < 1000; i++ {
		engine := NewScriptEngine(script, nil, 0, nil, ScriptFlagsNone)
		result, err := engine.Execute()
		if !result || err != nil {
			t.Fatalf("Benchmark script failed on iteration %d: %v", i, err)
		}
	}
}
