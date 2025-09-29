package bitcoin_test

import (
	"encoding/hex"
	"testing"
	"bitcoinecho.org/node/pkg/bitcoin"
)

// TestScript_AnalyzeScript tests script type detection with real Bitcoin scripts
func TestScript_AnalyzeScript(t *testing.T) {
	tests := []struct {
		name     string
		script   string // hex representation
		expected bitcoin.ScriptType
	}{
		// P2PKH (Pay-to-Public-Key-Hash) scripts
		{
			name:     "P2PKH standard script",
			script:   "76a914389ffce9cd9ae88dcc0631e88a821ffdbe9bfe2688ac", // Real P2PKH (25 bytes)
			expected: bitcoin.ScriptTypeP2PKH,
		},
		{
			name:     "P2PKH Genesis Block coinbase output",
			script:   "4104678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5fac", // Genesis P2PK (not P2PKH)
			expected: bitcoin.ScriptTypeP2PK,
		},
		{
			name:     "P2PKH another real example",
			script:   "76a9141b72503639a13f190bf79acf6d76255d772360b088ac",
			expected: bitcoin.ScriptTypeP2PKH,
		},

		// P2SH (Pay-to-Script-Hash) scripts
		{
			name:     "P2SH standard script",
			script:   "a91487916d4c8984d29dc696c7c9e14c9c9ad44b1e5987", // Real P2SH
			expected: bitcoin.ScriptTypeP2SH,
		},
		{
			name:     "P2SH multisig wrapper",
			script:   "a914b7fcfa3c16db5d7c17cd2db5e4e6b5cd12b5b47287",
			expected: bitcoin.ScriptTypeP2SH,
		},

		// P2PK (Pay-to-Public-Key) scripts - legacy format
		{
			name:     "P2PK compressed pubkey",
			script:   "21034f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa5288ac",
			expected: bitcoin.ScriptTypeP2PK,
		},
		{
			name:     "P2PK uncompressed pubkey",
			script:   "4104678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5fac",
			expected: bitcoin.ScriptTypeP2PK,
		},

		// P2WPKH (Pay-to-Witness-Public-Key-Hash) - Native SegWit
		{
			name:     "P2WPKH native SegWit",
			script:   "0014751e76ab4c23b27acb9b8e1c4c9c48c9e9f8a8b3", // OP_0 + 20-byte hash
			expected: bitcoin.ScriptTypeP2WPKH,
		},
		{
			name:     "P2WPKH another example",
			script:   "00141234567890abcdef1234567890abcdef12345678",
			expected: bitcoin.ScriptTypeP2WPKH,
		},

		// P2WSH (Pay-to-Witness-Script-Hash) - Native SegWit
		{
			name:     "P2WSH native SegWit",
			script:   "0020701a8d401c84fb13e6baf169d59684e17abd9fa216c8cc5b9fc63d622ff8c58d", // OP_0 + 32-byte hash
			expected: bitcoin.ScriptTypeP2WSH,
		},
		{
			name:     "P2WSH multisig script",
			script:   "00201234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			expected: bitcoin.ScriptTypeP2WSH,
		},

		// P2TR (Pay-to-Taproot) - Taproot
		{
			name:     "P2TR Taproot script",
			script:   "5120751e76ab4c23b27acb9b8e1c4c9c48c9e9f8a8b3751e76ab4c23b27acb9b8e1c", // OP_1 + 32-byte key
			expected: bitcoin.ScriptTypeP2TR,
		},

		// Multisig scripts
		{
			name:     "Multisig 2-of-3",
			script:   "5221034f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa52103459b6b64b9ce7f2e3ad0a9c60b8ddf6d87f1ad95e00db2e41d79e7c71d6be9021037e6d1d1a05f7c4e2b0f9c7b2e2c4b7a6d0a4a2c3e7b8d7a0d1f2c9b4e6a8d0a653ae",
			expected: bitcoin.ScriptTypeMultisig,
		},
		{
			name:     "Multisig 1-of-2",
			script:   "5121034f355bdcb7cc0af728ef3cceb9615d90684bb5b2ca5f859ab0f0b704075871aa521037e6d1d1a05f7c4e2b0f9c7b2e2c4b7a6d0a4a2c3e7b8d7a0d1f2c9b4e6a8d0a52ae",
			expected: bitcoin.ScriptTypeMultisig,
		},

		// OP_RETURN (Null Data) scripts
		{
			name:     "OP_RETURN with data",
			script:   "6a0b48656c6c6f20576f726c64", // OP_RETURN "Hello World"
			expected: bitcoin.ScriptTypeNullData,
		},
		{
			name:     "OP_RETURN empty",
			script:   "6a", // Just OP_RETURN
			expected: bitcoin.ScriptTypeNullData,
		},
		{
			name:     "OP_RETURN with longer data",
			script:   "6a4c50546869732069732061206c6f6e6720444154414441544144415441444154414441544144415441444154414441544144415441444154414441544144415441444154414441544144415441",
			expected: bitcoin.ScriptTypeNullData,
		},

		// Edge cases and invalid scripts
		{
			name:     "Empty script",
			script:   "",
			expected: bitcoin.ScriptTypeUnknown,
		},
		{
			name:     "Random bytes",
			script:   "deadbeef",
			expected: bitcoin.ScriptTypeUnknown,
		},
		{
			name:     "Almost P2PKH but wrong length",
			script:   "76a914389ffce9cd9ae88dcc0631e88a821ffdbe9bfe261558", // Missing final bytes (fixed hex)
			expected: bitcoin.ScriptTypeUnknown,
		},
		{
			name:     "Almost P2SH but wrong opcode",
			script:   "a81487916d4c8984d29dc696c7c9e14c9c9ad44b1e59c087", // Wrong first opcode
			expected: bitcoin.ScriptTypeUnknown,
		},
		{
			name:     "P2WPKH wrong hash length",
			script:   "0015751e76ab4c23b27acb9b8e1c4c9c48c9e9f8a8b3ff", // 21 bytes instead of 20
			expected: bitcoin.ScriptTypeUnknown,
		},
		{
			name:     "P2WSH wrong hash length",
			script:   "0021701a8d401c84fb13e6baf169d59684e17abd9fa216c8cc5b9fc63d622ff8c58dff", // 33 bytes instead of 32
			expected: bitcoin.ScriptTypeUnknown,
		},
		{
			name:     "Non-standard script with valid opcodes",
			script:   "6351676351676351676351",
			expected: bitcoin.ScriptTypeUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scriptBytes, err := hex.DecodeString(tt.script)
			if err != nil && tt.script != "" {
				t.Fatalf("Failed to decode hex script: %v", err)
			}

			script := bitcoin.Script(scriptBytes)
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
			name:     "Very large multisig is not standard (>3 pubkeys)",
			script:   "54" + // OP_4 (4-of-N)
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

			script := bitcoin.Script(scriptBytes)
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
			script := bitcoin.Script(tt.script)
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
	script := bitcoin.Script(scriptBytes)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = script.AnalyzeScript()
	}
}

// BenchmarkScript_IsStandard benchmarks standardness checking
func BenchmarkScript_IsStandard(b *testing.B) {
	// P2PKH script for benchmarking
	scriptBytes, _ := hex.DecodeString("76a914389ffce9cd9ae88dcc0631e88a821ffdbe9bfe2615588ac")
	script := bitcoin.Script(scriptBytes)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = script.IsStandard()
	}
}