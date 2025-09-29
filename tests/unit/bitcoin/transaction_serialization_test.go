package bitcoin_test

import (
	"bytes"
	"encoding/hex"
	"testing"
	"bitcoinecho.org/node/pkg/bitcoin"
)

// TestTransaction_Serialize tests Bitcoin transaction serialization to wire format
func TestTransaction_Serialize(t *testing.T) {
	tests := []struct {
		name        string
		transaction *bitcoin.Transaction
		expectedHex string
		isWitness   bool
		description string
	}{
		{
			name: "Genesis block coinbase transaction",
			transaction: &bitcoin.Transaction{
				Version: 1,
				Inputs: []bitcoin.TxInput{
					{
						PreviousOutput: bitcoin.OutPoint{
							Hash:  bitcoin.Hash256{}, // Zero hash
							Index: 0xFFFFFFFF, // Coinbase marker
						},
						ScriptSig: mustDecodeHex("4d04ffff001d0104455468652054696d65732030332f4a616e2f32303039204368616e63656c6c6f72206f6e206272696e6b206f66207365636f6e64206261696c6f757420666f722062616e6b73"),
						Sequence:  0xFFFFFFFF,
					},
				},
				Outputs: []bitcoin.TxOutput{
					{
						Value:        5000000000, // 50 BTC in satoshis
						ScriptPubKey: mustDecodeHex("4104678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5fac"),
					},
				},
				LockTime: 0,
			},
			expectedHex: "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff4d04ffff001d0104455468652054696d65732030332f4a616e2f32303039204368616e63656c6c6f72206f6e206272696e6b206f66207365636f6e64206261696c6f757420666f722062616e6b73ffffffff0100f2052a01000000434104678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5fac00000000",
			isWitness:   false,
			description: "First Bitcoin transaction ever created",
		},
		{
			name: "Simple P2PKH transaction (legacy format)",
			transaction: &bitcoin.Transaction{
				Version: 1,
				Inputs: []bitcoin.TxInput{
					{
						PreviousOutput: bitcoin.OutPoint{
							Hash:  mustDecodeHash256("7957a35fe64f80d234d76d83a2a8f1a0d8149a41d81de548f0a65a8a999f6f18"),
							Index: 0,
						},
						ScriptSig: mustDecodeHex("483045022100884d142d86652a3f47ba4746ec719bbfbd040a570b1deccbb6498c75c4ae24cb02204b9f039ff08df09cbe9f6addac960298cad530a863ea8f53982c09db8f6e381301"),
						Sequence:  0xFFFFFFFF,
					},
				},
				Outputs: []bitcoin.TxOutput{
					{
						Value:    5000000000, // 50 BTC
						ScriptPubKey: mustDecodeHex("76a914389ffce9cd9ae88dcc0631e88a821ffdbe9bfe2688ac"), // P2PKH
					},
					{
						Value:    4900000000, // 49 BTC (change)
						ScriptPubKey: mustDecodeHex("76a914389ffce9cd9ae88dcc0631e88a821ffdbe9bfe2688ac"), // P2PKH
					},
				},
				LockTime: 0,
			},
			expectedHex: "0100000001186f9f998a5aa6f048e51dd8419a14d8a0f1a8a2836dd734d2804fe65fa35779000000006b483045022100884d142d86652a3f47ba4746ec719bbfbd040a570b1deccbb6498c75c4ae24cb02204b9f039ff08df09cbe9f6addac960298cad530a863ea8f53982c09db8f6e3813014104278485b64e6c5eb9e97e51e6e7b5c0b99e5b4e6e35e9e8fb8e4e8e4e8e4e8e4e8e4e8e4e8e4e8e4e8e4e8e4e8e4e8e4e8e4e8e4e8e4e8e4e8e4e8e4e8e4e8e4e8ffffffff0200f2052a010000001976a914389ffce9cd9ae88dcc0631e88a821ffdbe9bfe2688ac00286bee010000001976a914389ffce9cd9ae88dcc0631e88a821ffdbe9bfe2688ac00000000",
			isWitness:   false,
			description: "Standard P2PKH transaction with two outputs",
		},
		{
			name: "SegWit transaction with witness data",
			transaction: &bitcoin.Transaction{
				Version: 2,
				Inputs: []bitcoin.TxInput{
					{
						PreviousOutput: bitcoin.OutPoint{
							Hash:  mustDecodeHash256("fff7f7881a8dd6f5d734f3bc3f8f03e93c56d3f63f7d3a5a6f8b4b4b4b4b4b4b"),
							Index: 0,
						},
						ScriptSig: mustDecodeHex(""), // Empty scriptSig for SegWit
						Sequence:  0xEEEEEEEE,
						Witness: [][]byte{
							mustDecodeHex("304402203ad1cc746a3cb70dd42ce5d7b6b0c7b8f9e5e0c2d1a2b3c4d5e6f708090a0b0c02200b1c2d3e4f506172839405a6b7c8d9e0f1a2b3c4d5e6f708090a0b0c0d0e0f01"),
							mustDecodeHex("0279be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798"),
						},
					},
				},
				Outputs: []bitcoin.TxOutput{
					{
						Value:    199996600, // ~2 BTC
						ScriptPubKey: mustDecodeHex("00140279be667ef9dcbbac55a06295ce870b07029bf0"), // P2WPKH
					},
				},
				LockTime: 17,
			},
			expectedHex: "020000000001014b4b4b4b4b4b4b8b6f5a3a7d3ff6d3563ce9038f3fbc4f73d7f5d68d1a88f7f7ff0000000000eeeeeeee01d8eaae0b000000001600140279be667ef9dcbbac55a06295ce870b07029bf0247304402203ad1cc746a3cb70dd42ce5d7b6b0c7b8f9e5e0c2d1a2b3c4d5e6f708090a0b0c02200b1c2d3e4f506172839405a6b7c8d9e0f1a2b3c4d5e6f708090a0b0c0d0e0f01210279be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f8179811000000",
			isWitness:   true,
			description: "SegWit P2WPKH transaction with witness data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test serialization
			serialized, err := tt.transaction.Serialize()
			if err != nil {
				t.Fatalf("Failed to serialize transaction: %v", err)
			}

			actualHex := hex.EncodeToString(serialized)
			if actualHex != tt.expectedHex {
				t.Errorf("Serialized transaction hex mismatch:\nExpected: %s\nActual:   %s", tt.expectedHex, actualHex)
			}

			// Test witness detection
			hasWitness := tt.transaction.HasWitness()
			if hasWitness != tt.isWitness {
				t.Errorf("Witness detection failed: expected %v, got %v", tt.isWitness, hasWitness)
			}

			t.Logf("✓ %s", tt.description)
		})
	}
}

// TestTransaction_Deserialize tests Bitcoin transaction deserialization from wire format
func TestTransaction_Deserialize(t *testing.T) {
	tests := []struct {
		name         string
		inputHex     string
		expectedTx   *bitcoin.Transaction
		isWitness    bool
		description  string
	}{
		{
			name:     "Genesis block coinbase transaction",
			inputHex: "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff4d04ffff001d0104455468652054696d65732030332f4a616e2f32303039204368616e63656c6c6f72206f6e206272696e6b206f66207365636f6e64206261696c6f757420666f722062616e6b73ffffffff0100f2052a01000000434104678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5fac00000000",
			expectedTx: &bitcoin.Transaction{
				Version: 1,
				Inputs: []bitcoin.TxInput{
					{
						PreviousOutput: bitcoin.OutPoint{
							Hash:  bitcoin.Hash256{}, // Zero hash
							Index: 0xFFFFFFFF,
						},
						ScriptSig: mustDecodeHex("4d04ffff001d0104455468652054696d65732030332f4a616e2f32303039204368616e63656c6c6f72206f6e206272696e6b206f66207365636f6e64206261696c6f757420666f722062616e6b73"),
						Sequence:  0xFFFFFFFF,
					},
				},
				Outputs: []bitcoin.TxOutput{
					{
						Value:    5000000000,
						ScriptPubKey: mustDecodeHex("4104678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5fac"),
					},
				},
				LockTime: 0,
			},
			isWitness:   false,
			description: "Deserialize genesis coinbase transaction",
		},
		{
			name:     "Simple legacy transaction",
			inputHex: "0100000001c997a5e56e104102fa209c6a852dd90660a20b2d9c352423edce25857fcd3704000000004847304402204e45e16932b8af514961a1d3a1a25fdf3f4f7732e9d624c6c61548ab5fb8cd410220181522ec8eca07de4860a4acdd12909d831cc56cbbac4622082221a8768d1d0901ffffffff0200ca9a3b00000000434104ae1a62fe09c5f51b13905f07f06b99a2f7159b2225f374cd378d71302fa28414e7aab37397f554a7df5f142c21c1b7303b8a0626f1baded5c72a704f7e6cd84cac00286bee0000000043410411db93e1dcdb8a016b49840f8c53bc1eb68a382e97b1482ecad7b148a6909a5cb2e0eaddfb84ccf9744464f82e160bfa9b8b64f9d4c03f999b8643f656b412a3ac00000000",
			expectedTx: &bitcoin.Transaction{
				Version: 1,
				Inputs: []bitcoin.TxInput{
					{
						PreviousOutput: bitcoin.OutPoint{
							Hash:  mustDecodeHash256("0437cd7f8525ceed2324359c2d0ba26006d92d856a9c20fa0241106ee5a597c9"),
							Index: 0,
						},
						ScriptSig: mustDecodeHex("47304402204e45e16932b8af514961a1d3a1a25fdf3f4f7732e9d624c6c61548ab5fb8cd410220181522ec8eca07de4860a4acdd12909d831cc56cbbac4622082221a8768d1d0901"),
						Sequence:  0xFFFFFFFF,
					},
				},
				Outputs: []bitcoin.TxOutput{
					{
						Value:    1000000000,
						ScriptPubKey: mustDecodeHex("4104ae1a62fe09c5f51b13905f07f06b99a2f7159b2225f374cd378d71302fa28414e7aab37397f554a7df5f142c21c1b7303b8a0626f1baded5c72a704f7e6cd84cac"),
					},
					{
						Value:    4000000000,
						ScriptPubKey: mustDecodeHex("410411db93e1dcdb8a016b49840f8c53bc1eb68a382e97b1482ecad7b148a6909a5cb2e0eaddfb84ccf9744464f82e160bfa9b8b64f9d4c03f999b8643f656b412a3ac"),
					},
				},
				LockTime: 0,
			},
			isWitness:   false,
			description: "Deserialize early Bitcoin P2PK transaction",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Decode hex input
			txBytes, err := hex.DecodeString(tt.inputHex)
			if err != nil {
				t.Fatalf("Failed to decode hex: %v", err)
			}

			// Test deserialization
			tx, err := bitcoin.DeserializeTransaction(txBytes)
			if err != nil {
				t.Fatalf("Failed to deserialize transaction: %v", err)
			}

			// Compare version
			if tx.Version != tt.expectedTx.Version {
				t.Errorf("Version mismatch: expected %d, got %d", tt.expectedTx.Version, tx.Version)
			}

			// Compare inputs
			if len(tx.Inputs) != len(tt.expectedTx.Inputs) {
				t.Errorf("Input count mismatch: expected %d, got %d", len(tt.expectedTx.Inputs), len(tx.Inputs))
			}

			for i, expectedInput := range tt.expectedTx.Inputs {
				if i >= len(tx.Inputs) {
					break
				}
				actualInput := tx.Inputs[i]

				if actualInput.PreviousOutput.Hash != expectedInput.PreviousOutput.Hash {
					t.Errorf("Input %d hash mismatch", i)
				}
				if actualInput.PreviousOutput.Index != expectedInput.PreviousOutput.Index {
					t.Errorf("Input %d index mismatch: expected %d, got %d", i, expectedInput.PreviousOutput.Index, actualInput.PreviousOutput.Index)
				}
				if !bytes.Equal(actualInput.ScriptSig, expectedInput.ScriptSig) {
					t.Errorf("Input %d scriptSig mismatch", i)
				}
				if actualInput.Sequence != expectedInput.Sequence {
					t.Errorf("Input %d sequence mismatch: expected %d, got %d", i, expectedInput.Sequence, actualInput.Sequence)
				}
			}

			// Compare outputs
			if len(tx.Outputs) != len(tt.expectedTx.Outputs) {
				t.Errorf("Output count mismatch: expected %d, got %d", len(tt.expectedTx.Outputs), len(tx.Outputs))
			}

			for i, expectedOutput := range tt.expectedTx.Outputs {
				if i >= len(tx.Outputs) {
					break
				}
				actualOutput := tx.Outputs[i]

				if actualOutput.Value != expectedOutput.Value {
					t.Errorf("Output %d value mismatch: expected %d, got %d", i, expectedOutput.Value, actualOutput.Value)
				}
				if !bytes.Equal(actualOutput.ScriptPubKey, expectedOutput.ScriptPubKey) {
					t.Errorf("Output %d scriptPK mismatch", i)
				}
			}

			// Compare locktime
			if tx.LockTime != tt.expectedTx.LockTime {
				t.Errorf("LockTime mismatch: expected %d, got %d", tt.expectedTx.LockTime, tx.LockTime)
			}

			t.Logf("✓ %s", tt.description)
		})
	}
}

// TestTransaction_RoundTrip tests serialize -> deserialize -> serialize consistency
func TestTransaction_RoundTrip(t *testing.T) {
	tests := []struct {
		name        string
		inputHex    string
		description string
	}{
		{
			name:        "Genesis coinbase transaction",
			inputHex:    "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff4d04ffff001d0104455468652054696d65732030332f4a616e2f32303039204368616e63656c6c6f72206f6e206272696e6b206f66207365636f6e64206261696c6f757420666f722062616e6b73ffffffff0100f2052a01000000434104678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5fac00000000",
			description: "Round-trip genesis transaction",
		},
		{
			name:        "Early Bitcoin P2PK transaction",
			inputHex:    "0100000001c997a5e56e104102fa209c6a852dd90660a20b2d9c352423edce25857fcd3704000000004847304402204e45e16932b8af514961a1d3a1a25fdf3f4f7732e9d624c6c61548ab5fb8cd410220181522ec8eca07de4860a4acdd12909d831cc56cbbac4622082221a8768d1d0901ffffffff0200ca9a3b00000000434104ae1a62fe09c5f51b13905f07f06b99a2f7159b2225f374cd378d71302fa28414e7aab37397f554a7df5f142c21c1b7303b8a0626f1baded5c72a704f7e6cd84cac00286bee0000000043410411db93e1dcdb8a016b49840f8c53bc1eb68a382e97b1482ecad7b148a6909a5cb2e0eaddfb84ccf9744464f82e160bfa9b8b64f9d4c03f999b8643f656b412a3ac00000000",
			description: "Round-trip P2PK transaction",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Step 1: Decode original hex
			originalBytes, err := hex.DecodeString(tt.inputHex)
			if err != nil {
				t.Fatalf("Failed to decode original hex: %v", err)
			}

			// Step 2: Deserialize to Transaction
			tx, err := bitcoin.DeserializeTransaction(originalBytes)
			if err != nil {
				t.Fatalf("Failed to deserialize transaction: %v", err)
			}

			// Step 3: Serialize back to bytes
			serializedBytes, err := tx.Serialize()
			if err != nil {
				t.Fatalf("Failed to serialize transaction: %v", err)
			}

			// Step 4: Compare original and round-trip bytes
			if !bytes.Equal(originalBytes, serializedBytes) {
				t.Errorf("Round-trip failed for %s", tt.description)
				t.Errorf("Original:  %s", hex.EncodeToString(originalBytes))
				t.Errorf("Round-trip: %s", hex.EncodeToString(serializedBytes))
			}

			t.Logf("✓ %s", tt.description)
		})
	}
}

// TestVarInt tests variable integer encoding/decoding
func TestVarInt(t *testing.T) {
	tests := []struct {
		name        string
		value       uint64
		expectedHex string
		description string
	}{
		{
			name:        "Single byte (0-252)",
			value:       42,
			expectedHex: "2a",
			description: "Values under 253 use 1 byte",
		},
		{
			name:        "Two bytes (253-65535)",
			value:       1000,
			expectedHex: "fd03e8",
			description: "Values 253-65535 use 3 bytes (fd + 2 bytes LE)",
		},
		{
			name:        "Four bytes (65536-4294967295)",
			value:       100000,
			expectedHex: "fe00000186a0",
			description: "Values 65536+ use 5 bytes (fe + 4 bytes LE)",
		},
		{
			name:        "Eight bytes (4294967296+)",
			value:       5000000000,
			expectedHex: "ff00000001000000002af31dc4",
			description: "Values over 4GB use 9 bytes (ff + 8 bytes LE)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test encoding
			encoded := bitcoin.EncodeVarInt(tt.value)
			actualHex := hex.EncodeToString(encoded)
			if actualHex != tt.expectedHex {
				t.Errorf("VarInt encode failed: expected %s, got %s", tt.expectedHex, actualHex)
			}

			// Test decoding
			encodedBytes, _ := hex.DecodeString(tt.expectedHex)
			decoded, bytesRead := bitcoin.DecodeVarInt(encodedBytes)
			if decoded != tt.value {
				t.Errorf("VarInt decode failed: expected %d, got %d", tt.value, decoded)
			}
			if bytesRead != len(encodedBytes) {
				t.Errorf("VarInt bytes read mismatch: expected %d, got %d", len(encodedBytes), bytesRead)
			}

			t.Logf("✓ %s", tt.description)
		})
	}
}

// Helper functions for test data
func mustDecodeHex(hexStr string) []byte {
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		panic("Invalid hex in test data: " + err.Error())
	}
	return data
}

func mustDecodeHash256(hexStr string) bitcoin.Hash256 {
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		panic("Invalid hash hex in test data: " + err.Error())
	}
	if len(data) != 32 {
		panic("Hash256 must be 32 bytes")
	}
	var hash bitcoin.Hash256
	copy(hash[:], data)
	return hash
}