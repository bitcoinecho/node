package bitcoin

import (
	"bytes"
	"encoding/hex"
	"testing"
)

// TestScript_AnalyzeScript tests script type detection with real Bitcoin scripts
func TestScript_AnalyzeScript(t *testing.T) {
	tests := []struct {
		name     string
		script   string // hex representation
		expected ScriptType
	}{
		// P2PKH (Pay-to-Public-Key-Hash) scripts
		{
			name:     "P2PKH standard script",
			script:   "76a914389ffce9cd9ae88dcc0631e88a821ffdbe9bfe2688ac", // Real P2PKH (25 bytes)
			expected: ScriptTypeP2PKH,
		},
		{
			name:     "P2PKH Genesis Block coinbase output",
			script:   "4104678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5fac", // Genesis P2PK (not P2PKH)
			expected: ScriptTypeP2PK,
		},
		{
			name:     "P2PKH another real example",
			script:   "76a9141b72503639a13f190bf79acf6d76255d772360b088ac",
			expected: ScriptTypeP2PKH,
		},

		// P2SH (Pay-to-Script-Hash) scripts
		{
			name:     "P2SH standard script",
			script:   "a91487916d4c8984d29dc696c7c9e14c9c9ad44b1e5987", // Real P2SH
			expected: ScriptTypeP2SH,
		},
		{
			name:     "P2SH multisig wrapper",
			script:   "a914b7fcfa3c16db5d7c17cd2db5e4e6b5cd12b5b47287",
			expected: ScriptTypeP2SH,
		},

		// P2PK (Pay-to-Public-Key) scripts - legacy format
		{
			name:     "P2PK compressed pubkey",
			script:   "21034f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa5288ac",
			expected: ScriptTypeP2PK,
		},
		{
			name:     "P2PK uncompressed pubkey",
			script:   "4104678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5fac",
			expected: ScriptTypeP2PK,
		},

		// P2WPKH (Pay-to-Witness-Public-Key-Hash) - Native SegWit
		{
			name:     "P2WPKH native SegWit",
			script:   "0014751e76ab4c23b27acb9b8e1c4c9c48c9e9f8a8b3", // OP_0 + 20-byte hash
			expected: ScriptTypeP2WPKH,
		},
		{
			name:     "P2WPKH another example",
			script:   "00141234567890abcdef1234567890abcdef12345678",
			expected: ScriptTypeP2WPKH,
		},

		// P2WSH (Pay-to-Witness-Script-Hash) - Native SegWit
		{
			name:     "P2WSH native SegWit",
			script:   "0020701a8d401c84fb13e6baf169d59684e17abd9fa216c8cc5b9fc63d622ff8c58d", // OP_0 + 32-byte hash
			expected: ScriptTypeP2WSH,
		},
		{
			name:     "P2WSH multisig script",
			script:   "00201234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			expected: ScriptTypeP2WSH,
		},

		// P2TR (Pay-to-Taproot) - Taproot
		{
			name:     "P2TR Taproot script",
			script:   "5120751e76ab4c23b27acb9b8e1c4c9c48c9e9f8a8b3751e76ab4c23b27acb9b8e1c", // OP_1 + 32-byte key
			expected: ScriptTypeP2TR,
		},

		// Multisig scripts
		{
			name:     "Multisig 2-of-3",
			script:   "5221034f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa52103459b6b64b9ce7f2e3ad0a9c60b8ddf6d87f1ad95e00db2e41d79e7c71d6be9021037e6d1d1a05f7c4e2b0f9c7b2e2c4b7a6d0a4a2c3e7b8d7a0d1f2c9b4e6a8d0a653ae",
			expected: ScriptTypeMultisig,
		},
		{
			name:     "Multisig 1-of-2",
			script:   "5121034f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa521037e6d1d1a05f7c4e2b0f9c7b2e2c4b7a6d0a4a2c3e7b8d7a0d1f2c9b4e6a8d0a52ae",
			expected: ScriptTypeMultisig,
		},

		// OP_RETURN (Null Data) scripts
		{
			name:     "OP_RETURN with data",
			script:   "6a0b48656c6c6f20576f726c64", // OP_RETURN "Hello World"
			expected: ScriptTypeNullData,
		},
		{
			name:     "OP_RETURN empty",
			script:   "6a", // Just OP_RETURN
			expected: ScriptTypeNullData,
		},
		{
			name:     "OP_RETURN with longer data",
			script:   "6a4c50546869732069732061206c6f6e6720444154414441544144415441444154414441544144415441444154414441544144415441444154414441544144415441444154414441544144415441",
			expected: ScriptTypeNullData,
		},

		// Edge cases and invalid scripts
		{
			name:     "Empty script",
			script:   "",
			expected: ScriptTypeUnknown,
		},
		{
			name:     "Random bytes",
			script:   "deadbeef",
			expected: ScriptTypeUnknown,
		},
		{
			name:     "Almost P2PKH but wrong length",
			script:   "76a914389ffce9cd9ae88dcc0631e88a821ffdbe9bfe261558", // Missing final bytes (fixed hex)
			expected: ScriptTypeUnknown,
		},
		{
			name:     "Almost P2SH but wrong opcode",
			script:   "a81487916d4c8984d29dc696c7c9e14c9c9ad44b1e59c087", // Wrong first opcode
			expected: ScriptTypeUnknown,
		},
		{
			name:     "P2WPKH wrong hash length",
			script:   "0015751e76ab4c23b27acb9b8e1c4c9c48c9e9f8a8b3ff", // 21 bytes instead of 20
			expected: ScriptTypeUnknown,
		},
		{
			name:     "P2WSH wrong hash length",
			script:   "0021701a8d401c84fb13e6baf169d59684e17abd9fa216c8cc5b9fc63d622ff8c58dff", // 33 bytes instead of 32
			expected: ScriptTypeUnknown,
		},
		{
			name:     "Non-standard script with valid opcodes",
			script:   "6351676351676351676351",
			expected: ScriptTypeUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scriptBytes, err := hex.DecodeString(tt.script)
			if err != nil && tt.script != "" {
				t.Fatalf("Failed to decode hex script: %v", err)
			}

			script := Script(scriptBytes)
			result := script.AnalyzeScript()

			if result != tt.expected {
				t.Errorf("Expected script type %v, got %v\nScript: %s",
					tt.expected, result, tt.script)
			}
		})
	}
}

// TestScript_IsStandard tests script standardness rules
func TestScript_IsStandard(t *testing.T) {
	tests := []struct {
		name     string
		script   string // hex representation
		expected bool
	}{
		// Standard scripts should return true
		{
			name:     "P2PKH is standard",
			script:   "76a914389ffce9cd9ae88dcc0631e88a821ffdbe9bfe2688ac",
			expected: true,
		},
		{
			name:     "P2SH is standard",
			script:   "a91487916d4c8984d29dc696c7c9e14c9c9ad44b1e5987",
			expected: true,
		},
		{
			name:     "P2WPKH is standard",
			script:   "0014751e76ab4c23b27acb9b8e1c4c9c48c9e9f8a8b3",
			expected: true,
		},
		{
			name:     "P2WSH is standard",
			script:   "0020701a8d401c84fb13e6baf169d59684e17abd9fa216c8cc5b9fc63d622ff8c58d",
			expected: true,
		},
		{
			name:     "P2TR is standard",
			script:   "5120751e76ab4c23b27acb9b8e1c4c9c48c9e9f8a8b3751e76ab4c23b27acb9b8e1c",
			expected: true,
		},
		{
			name:     "P2PK compressed is standard",
			script:   "21034f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa5288ac",
			expected: true,
		},
		{
			name:     "P2PK uncompressed is standard",
			script:   "4104678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5fac",
			expected: true,
		},
		{
			name:     "Small multisig is standard",
			script:   "5121034f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa521037e6d1d1a05f7c4e2b0f9c7b2e2c4b7a6d0a4a2c3e7b8d7a0d1f2c9b4e6a8d0a52ae",
			expected: true,
		},
		{
			name:     "OP_RETURN with reasonable data is standard",
			script:   "6a0b48656c6c6f20576f726c64", // "Hello World"
			expected: true,
		},

		// Non-standard scripts should return false
		{
			name:     "Empty script is not standard",
			script:   "",
			expected: false,
		},
		{
			name:     "Unknown script type is not standard",
			script:   "deadbeef",
			expected: false,
		},
		{
			name:     "Large OP_RETURN is not standard (over 80 bytes)",
			script:   "6a4c50546869732069732061206c6f6e6720444154414441544144415441444154414441544144415441444154414441544144415441444154414441544144415441444154414441544144415441444154414441544144415441",
			expected: false, // This is >80 bytes of data
		},
		{
			name:     "Malformed P2PKH is not standard",
			script:   "76a914389ffce9cd9ae88dcc0631e88a821ffdbe9bfe2615", // Incomplete
			expected: false,
		},
		{
			name: "Very large multisig is not standard (>3 pubkeys)",
			script: "54" + // OP_4 (4-of-N)
				"21034f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa52" +
				"21037e6d1d1a05f7c4e2b0f9c7b2e2c4b7a6d0a4a2c3e7b8d7a0d1f2c9b4e6a8d0a6" +
				"21034f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa52" +
				"21037e6d1d1a05f7c4e2b0f9c7b2e2c4b7a6d0a4a2c3e7b8d7a0d1f2c9b4e6a8d0a6" +
				"21034f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa52" + // 5th pubkey
				"55ae", // OP_5 OP_CHECKMULTISIG
			expected: false,
		},
		{
			name:     "Non-standard script with weird opcodes",
			script:   "6351676351676351676351", // Random valid opcodes
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scriptBytes, err := hex.DecodeString(tt.script)
			if err != nil && tt.script != "" {
				t.Fatalf("Failed to decode hex script: %v", err)
			}

			script := Script(scriptBytes)
			result := script.IsStandard()

			if result != tt.expected {
				t.Errorf("Expected standard=%v, got %v\nScript: %s\nType: %v",
					tt.expected, result, tt.script, script.AnalyzeScript())
			}
		})
	}
}

// TestScript_String tests script string representation
// TODO: Implement String() method in future TDD iteration
/*
func TestScript_String(t *testing.T) {
	tests := []struct {
		name     string
		script   []byte
		expected string
	}{
		{
			name:     "Empty script",
			script:   []byte{},
			expected: "",
		},
		{
			name:     "Simple P2PKH script representation",
			script:   []byte{0x76, 0xa9, 0x14}, // OP_DUP OP_HASH160 PUSH(20)
			expected: "OP_DUP OP_HASH160 OP_PUSHDATA(20)",
		},
		{
			name:     "OP_RETURN script",
			script:   []byte{0x6a, 0x0b}, // OP_RETURN PUSH(11)
			expected: "OP_RETURN OP_PUSHDATA(11)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			script := Script(tt.script)
			result := script.String()

			if result != tt.expected {
				t.Errorf("Expected string %q, got %q", tt.expected, result)
			}
		})
	}
}
*/

// BenchmarkScript_AnalyzeScript benchmarks script analysis
func BenchmarkScript_AnalyzeScript(b *testing.B) {
	// P2PKH script for benchmarking
	scriptBytes, _ := hex.DecodeString("76a914389ffce9cd9ae88dcc0631e88a821ffdbe9bfe2615588ac")
	script := Script(scriptBytes)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = script.AnalyzeScript()
	}
}

// BenchmarkScript_IsStandard benchmarks standardness checking
func BenchmarkScript_IsStandard(b *testing.B) {
	// P2PKH script for benchmarking
	scriptBytes, _ := hex.DecodeString("76a914389ffce9cd9ae88dcc0631e88a821ffdbe9bfe2615588ac")
	script := Script(scriptBytes)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = script.IsStandard()
	}
}

// TestScriptEngine_SetScript tests setting a new script
func TestScriptEngine_SetScript(t *testing.T) {
	// Create script engine with initial script
	initialScript := Script([]byte{0x51}) // OP_1
	tx := &Transaction{}
	prevOuts := []TxOutput{}
	flags := ScriptFlags(0)

	engine := NewScriptEngine(initialScript, tx, 0, prevOuts, flags)

	// Set new script
	newScript := Script([]byte{0x52, 0x53}) // OP_2 OP_3
	engine.SetScript(newScript)

	// Execute and verify the new script was set
	result, err := engine.Execute()
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	// Should have executed OP_2 OP_3, so stack should have [2, 3]
	if !result {
		t.Error("Script execution should have succeeded")
	}

	stack := engine.GetStack()
	if len(stack) != 2 {
		t.Errorf("Expected stack length 2, got %d", len(stack))
	}
}

// TestScriptEngine_BytesToNum tests byte array to number conversion
func TestScriptEngine_BytesToNum(t *testing.T) {
	script := Script([]byte{0x51}) // OP_1
	tx := &Transaction{}
	prevOuts := []TxOutput{}
	flags := ScriptFlags(0)

	engine := NewScriptEngine(script, tx, 0, prevOuts, flags)

	tests := []struct {
		name     string
		input    []byte
		expected int64
	}{
		{
			name:     "empty bytes",
			input:    []byte{},
			expected: 0,
		},
		{
			name:     "single byte positive",
			input:    []byte{0x01},
			expected: 1,
		},
		{
			name:     "single byte negative",
			input:    []byte{0x81},
			expected: -1,
		},
		{
			name:     "multi-byte positive",
			input:    []byte{0x01, 0x02},
			expected: 513, // 0x0201 in little-endian
		},
		{
			name:     "multi-byte negative",
			input:    []byte{0x01, 0x82},
			expected: -513, // 0x0201 with sign bit
		},
		{
			name:     "too many bytes (overflow protection)",
			input:    make([]byte, 10), // More than 8 bytes
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := engine.bytesToNum(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

// TestScriptEngine_NumToBytes tests number to byte array conversion
func TestScriptEngine_NumToBytes(t *testing.T) {
	script := Script([]byte{0x51}) // OP_1
	tx := &Transaction{}
	prevOuts := []TxOutput{}
	flags := ScriptFlags(0)

	engine := NewScriptEngine(script, tx, 0, prevOuts, flags)

	tests := []struct {
		name     string
		input    int64
		expected []byte
	}{
		{
			name:     "zero",
			input:    0,
			expected: []byte{},
		},
		{
			name:     "positive single byte",
			input:    1,
			expected: []byte{0x01},
		},
		{
			name:     "negative single byte",
			input:    -1,
			expected: []byte{0x81},
		},
		{
			name:     "multi-byte positive",
			input:    256,
			expected: []byte{0x00, 0x01},
		},
		{
			name:     "multi-byte negative",
			input:    -256,
			expected: []byte{0x00, 0x81},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := engine.numToBytes(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("Expected length %d, got %d", len(tt.expected), len(result))
				return
			}
			for i, b := range tt.expected {
				if result[i] != b {
					t.Errorf("Expected byte %d to be %x, got %x", i, b, result[i])
				}
			}
		})
	}
}

// TestScriptEngine_IsTrue tests truth value evaluation
func TestScriptEngine_IsTrue(t *testing.T) {
	script := Script([]byte{0x51}) // OP_1
	tx := &Transaction{}
	prevOuts := []TxOutput{}
	flags := ScriptFlags(0)

	engine := NewScriptEngine(script, tx, 0, prevOuts, flags)

	tests := []struct {
		name     string
		input    []byte
		expected bool
	}{
		{
			name:     "empty bytes (false)",
			input:    []byte{},
			expected: false,
		},
		{
			name:     "zero byte (false)",
			input:    []byte{0x00},
			expected: false,
		},
		{
			name:     "negative zero (false)",
			input:    []byte{0x80},
			expected: false,
		},
		{
			name:     "positive number (true)",
			input:    []byte{0x01},
			expected: true,
		},
		{
			name:     "negative number (true)",
			input:    []byte{0x81},
			expected: true,
		},
		{
			name:     "multiple zeros (false)",
			input:    []byte{0x00, 0x00},
			expected: false,
		},
		{
			name:     "zero with negative sign (false)",
			input:    []byte{0x00, 0x80},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := engine.isTrue(tt.input)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestScriptEngine_NumericOpcodes tests numeric opcode execution
func TestScriptEngine_NumericOpcodes(t *testing.T) {
	tests := []struct {
		name     string
		opcode   ScriptOpcode
		expected []byte
	}{
		{name: "OP_0", opcode: OP_0, expected: []byte{}},
		{name: "OP_1", opcode: OP_1, expected: []byte{1}},
		{name: "OP_2", opcode: OP_2, expected: []byte{2}},
		{name: "OP_3", opcode: OP_3, expected: []byte{3}},
		{name: "OP_4", opcode: OP_4, expected: []byte{4}},
		{name: "OP_5", opcode: OP_5, expected: []byte{5}},
		{name: "OP_6", opcode: OP_6, expected: []byte{6}},
		{name: "OP_7", opcode: OP_7, expected: []byte{7}},
		{name: "OP_8", opcode: OP_8, expected: []byte{8}},
		{name: "OP_9", opcode: OP_9, expected: []byte{9}},
		{name: "OP_10", opcode: OP_10, expected: []byte{10}},
		{name: "OP_11", opcode: OP_11, expected: []byte{11}},
		{name: "OP_12", opcode: OP_12, expected: []byte{12}},
		{name: "OP_13", opcode: OP_13, expected: []byte{13}},
		{name: "OP_14", opcode: OP_14, expected: []byte{14}},
		{name: "OP_15", opcode: OP_15, expected: []byte{15}},
		{name: "OP_16", opcode: OP_16, expected: []byte{16}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			script := Script([]byte{byte(tt.opcode)})
			tx := &Transaction{}
			prevOuts := []TxOutput{}
			flags := ScriptFlags(0)
			engine := NewScriptEngine(script, tx, 0, prevOuts, flags)

			err := engine.executeOpcode(tt.opcode)
			if err != nil {
				t.Errorf("Unexpected error executing %s: %v", tt.name, err)
				return
			}

			if len(engine.stack) != 1 {
				t.Errorf("Expected stack size 1, got %d", len(engine.stack))
				return
			}

			if len(tt.expected) == 0 && len(engine.stack[0]) != 0 {
				t.Errorf("Expected empty stack item for %s, got %v", tt.name, engine.stack[0])
			} else if len(tt.expected) > 0 && !bytes.Equal(engine.stack[0], tt.expected) {
				t.Errorf("Expected %v on stack for %s, got %v", tt.expected, tt.name, engine.stack[0])
			}
		})
	}
}

// TestScriptEngine_StackOpcodes tests stack manipulation opcodes
func TestScriptEngine_StackOpcodes(t *testing.T) {
	t.Run("OP_DUP", func(t *testing.T) {
		script := Script([]byte{})
		tx := &Transaction{}
		prevOuts := []TxOutput{}
		flags := ScriptFlags(0)
		engine := NewScriptEngine(script, tx, 0, prevOuts, flags)

		// Push test data and duplicate
		engine.stack = append(engine.stack, []byte{0x42})
		err := engine.executeOpcode(OP_DUP)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if len(engine.stack) != 2 {
			t.Errorf("Expected stack size 2, got %d", len(engine.stack))
		}

		if !bytes.Equal(engine.stack[0], []byte{0x42}) || !bytes.Equal(engine.stack[1], []byte{0x42}) {
			t.Errorf("OP_DUP failed to duplicate top stack item")
		}
	})

	t.Run("OP_DUP_insufficient_stack", func(t *testing.T) {
		script := Script([]byte{})
		tx := &Transaction{}
		prevOuts := []TxOutput{}
		flags := ScriptFlags(0)
		engine := NewScriptEngine(script, tx, 0, prevOuts, flags)

		err := engine.executeOpcode(OP_DUP)
		if err == nil {
			t.Error("Expected error for OP_DUP with empty stack")
		}
	})

	t.Run("OP_DROP", func(t *testing.T) {
		script := Script([]byte{})
		tx := &Transaction{}
		prevOuts := []TxOutput{}
		flags := ScriptFlags(0)
		engine := NewScriptEngine(script, tx, 0, prevOuts, flags)

		// Push test data and drop
		engine.stack = append(engine.stack, []byte{0x42})
		err := engine.executeOpcode(OP_DROP)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if len(engine.stack) != 0 {
			t.Errorf("Expected empty stack after OP_DROP, got size %d", len(engine.stack))
		}
	})

	t.Run("OP_SWAP", func(t *testing.T) {
		script := Script([]byte{})
		tx := &Transaction{}
		prevOuts := []TxOutput{}
		flags := ScriptFlags(0)
		engine := NewScriptEngine(script, tx, 0, prevOuts, flags)

		// Push test data and swap
		engine.stack = append(engine.stack, []byte{0x11}, []byte{0x22})
		err := engine.executeOpcode(OP_SWAP)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if len(engine.stack) != 2 {
			t.Errorf("Expected stack size 2, got %d", len(engine.stack))
		}

		if !bytes.Equal(engine.stack[0], []byte{0x22}) || !bytes.Equal(engine.stack[1], []byte{0x11}) {
			t.Errorf("OP_SWAP failed to swap stack items")
		}
	})
}

// TestScriptEngine_ArithmeticOpcodes tests arithmetic opcodes
func TestScriptEngine_ArithmeticOpcodes(t *testing.T) {
	t.Run("OP_ADD", func(t *testing.T) {
		script := Script([]byte{})
		tx := &Transaction{}
		prevOuts := []TxOutput{}
		flags := ScriptFlags(0)
		engine := NewScriptEngine(script, tx, 0, prevOuts, flags)

		// Push test numbers (3 + 2 = 5)
		engine.stack = append(engine.stack, []byte{3}, []byte{2})
		err := engine.executeOpcode(OP_ADD)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if len(engine.stack) != 1 {
			t.Errorf("Expected stack size 1, got %d", len(engine.stack))
		}

		if !bytes.Equal(engine.stack[0], []byte{5}) {
			t.Errorf("Expected result 5, got %v", engine.stack[0])
		}
	})

	t.Run("OP_ADD_insufficient_stack", func(t *testing.T) {
		script := Script([]byte{})
		tx := &Transaction{}
		prevOuts := []TxOutput{}
		flags := ScriptFlags(0)
		engine := NewScriptEngine(script, tx, 0, prevOuts, flags)

		// Only push one number
		engine.stack = append(engine.stack, []byte{3})
		err := engine.executeOpcode(OP_ADD)
		if err == nil {
			t.Error("Expected error for OP_ADD with insufficient stack")
		}
	})
}
