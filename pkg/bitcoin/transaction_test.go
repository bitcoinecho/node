package bitcoin

import (
	"testing"
)

// TestNewTransaction tests creating new transactions
func TestNewTransaction(t *testing.T) {
	// Create sample inputs and outputs
	hash, _ := NewHash256FromString("0000000000000000000000000000000000000000000000000000000000000001")
	outpoint := OutPoint{Hash: hash, Index: 0}

	input := TxInput{
		PreviousOutput: outpoint,
		ScriptSig:      []byte{0x76, 0xa9}, // Sample script
		Sequence:       0xffffffff,
	}

	output := TxOutput{
		Value:        5000000000, // 50 BTC
		ScriptPubKey: []byte{0x76, 0xa9, 0x14}, // Sample P2PKH script
	}

	tx := NewTransaction(1, []TxInput{input}, []TxOutput{output}, 0)

	if tx.Version != 1 {
		t.Errorf("expected version 1, got %d", tx.Version)
	}

	if len(tx.Inputs) != 1 {
		t.Errorf("expected 1 input, got %d", len(tx.Inputs))
	}

	if len(tx.Outputs) != 1 {
		t.Errorf("expected 1 output, got %d", len(tx.Outputs))
	}

	if tx.LockTime != 0 {
		t.Errorf("expected locktime 0, got %d", tx.LockTime)
	}
}

// TestTransaction_IsCoinbase tests coinbase transaction detection
func TestTransaction_IsCoinbase(t *testing.T) {
	tests := []struct {
		name     string
		tx       *Transaction
		expected bool
	}{
		{
			name: "valid coinbase transaction",
			tx: &Transaction{
				Inputs: []TxInput{{
					PreviousOutput: OutPoint{Hash: ZeroHash, Index: 0xffffffff},
					ScriptSig:      []byte("coinbase data"),
					Sequence:       0xffffffff,
				}},
				Outputs: []TxOutput{{Value: 5000000000, ScriptPubKey: []byte{0x76, 0xa9}}},
			},
			expected: true,
		},
		{
			name: "non-coinbase transaction",
			tx: &Transaction{
				Inputs: []TxInput{{
					PreviousOutput: OutPoint{Hash: Hash256{0x01}, Index: 0},
					ScriptSig:      []byte{},
					Sequence:       0xffffffff,
				}},
				Outputs: []TxOutput{{Value: 1000000, ScriptPubKey: []byte{0x76, 0xa9}}},
			},
			expected: false,
		},
		{
			name: "multiple inputs (not coinbase)",
			tx: &Transaction{
				Inputs: []TxInput{
					{PreviousOutput: OutPoint{Hash: ZeroHash, Index: 0xffffffff}},
					{PreviousOutput: OutPoint{Hash: ZeroHash, Index: 0}},
				},
				Outputs: []TxOutput{{Value: 1000000, ScriptPubKey: []byte{0x76, 0xa9}}},
			},
			expected: false,
		},
		{
			name: "wrong index for coinbase",
			tx: &Transaction{
				Inputs: []TxInput{{
					PreviousOutput: OutPoint{Hash: ZeroHash, Index: 0},
					ScriptSig:      []byte("coinbase data"),
					Sequence:       0xffffffff,
				}},
				Outputs: []TxOutput{{Value: 5000000000, ScriptPubKey: []byte{0x76, 0xa9}}},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.tx.IsCoinbase()
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestTransaction_TotalOutput tests total output calculation
func TestTransaction_TotalOutput(t *testing.T) {
	tests := []struct {
		name     string
		outputs  []TxOutput
		expected uint64
	}{
		{
			name:     "single output",
			outputs:  []TxOutput{{Value: 5000000000}},
			expected: 5000000000,
		},
		{
			name: "multiple outputs",
			outputs: []TxOutput{
				{Value: 1000000000},
				{Value: 2000000000},
				{Value: 500000000},
			},
			expected: 3500000000,
		},
		{
			name:     "zero outputs",
			outputs:  []TxOutput{},
			expected: 0,
		},
		{
			name: "outputs with zero value",
			outputs: []TxOutput{
				{Value: 1000000000},
				{Value: 0},
				{Value: 2000000000},
			},
			expected: 3000000000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &Transaction{Outputs: tt.outputs}
			result := tx.TotalOutput()
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

// TestTransaction_HasWitness tests witness data detection
func TestTransaction_HasWitness(t *testing.T) {
	tests := []struct {
		name      string
		witnesses []TxWitness
		expected  bool
	}{
		{
			name:      "no witness data",
			witnesses: []TxWitness{},
			expected:  false,
		},
		{
			name:      "nil witness data",
			witnesses: nil,
			expected:  false,
		},
		{
			name: "has witness data",
			witnesses: []TxWitness{
				{Stack: [][]byte{[]byte("witness data")}},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &Transaction{Witnesses: tt.witnesses}
			result := tx.HasWitness()
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestTransaction_Validate tests transaction validation
func TestTransaction_Validate(t *testing.T) {
	// Helper to create a valid outpoint
	validOutpoint := OutPoint{
		Hash:  Hash256{0x01},
		Index: 0,
	}

	tests := []struct {
		name        string
		tx          *Transaction
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid transaction",
			tx: &Transaction{
				Version: 1,
				Inputs: []TxInput{{
					PreviousOutput: validOutpoint,
					ScriptSig:      []byte{0x76, 0xa9},
					Sequence:       0xffffffff,
				}},
				Outputs: []TxOutput{{
					Value:        1000000000, // 10 BTC
					ScriptPubKey: []byte{0x76, 0xa9, 0x14},
				}},
				LockTime: 0,
			},
			expectError: false,
		},
		{
			name: "no inputs",
			tx: &Transaction{
				Version:  1,
				Inputs:   []TxInput{},
				Outputs:  []TxOutput{{Value: 1000000000, ScriptPubKey: []byte{0x76}}},
				LockTime: 0,
			},
			expectError: true,
			errorMsg:    "transaction has no inputs",
		},
		{
			name: "no outputs",
			tx: &Transaction{
				Version: 1,
				Inputs: []TxInput{{
					PreviousOutput: validOutpoint,
					ScriptSig:      []byte{0x76, 0xa9},
					Sequence:       0xffffffff,
				}},
				Outputs:  []TxOutput{},
				LockTime: 0,
			},
			expectError: true,
			errorMsg:    "transaction has no outputs",
		},
		{
			name: "duplicate inputs",
			tx: &Transaction{
				Version: 1,
				Inputs: []TxInput{
					{PreviousOutput: validOutpoint, ScriptSig: []byte{0x76}, Sequence: 0xffffffff},
					{PreviousOutput: validOutpoint, ScriptSig: []byte{0xa9}, Sequence: 0xffffffff}, // Same outpoint
				},
				Outputs: []TxOutput{{Value: 1000000000, ScriptPubKey: []byte{0x76}}},
			},
			expectError: true,
			errorMsg:    "transaction has duplicate inputs",
		},
		{
			name: "output value exceeds maximum",
			tx: &Transaction{
				Version: 1,
				Inputs: []TxInput{{
					PreviousOutput: validOutpoint,
					ScriptSig:      []byte{0x76},
					Sequence:       0xffffffff,
				}},
				Outputs: []TxOutput{{
					Value:        MaxMoney + 1,
					ScriptPubKey: []byte{0x76},
				}},
			},
			expectError: true,
			errorMsg:    "output 0 value exceeds maximum",
		},
		{
			name: "total output exceeds maximum",
			tx: &Transaction{
				Version: 1,
				Inputs: []TxInput{{
					PreviousOutput: validOutpoint,
					ScriptSig:      []byte{0x76},
					Sequence:       0xffffffff,
				}},
				Outputs: []TxOutput{
					{Value: MaxMoney / 2 + 1, ScriptPubKey: []byte{0x76}},
					{Value: MaxMoney / 2 + 1, ScriptPubKey: []byte{0xa9}},
				},
			},
			expectError: true,
			errorMsg:    "total output value exceeds maximum",
		},
		{
			name: "valid coinbase transaction",
			tx: &Transaction{
				Version: 1,
				Inputs: []TxInput{{
					PreviousOutput: OutPoint{Hash: ZeroHash, Index: 0xffffffff},
					ScriptSig:      []byte("The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"),
					Sequence:       0xffffffff,
				}},
				Outputs: []TxOutput{{
					Value:        5000000000, // 50 BTC
					ScriptPubKey: []byte{0x76, 0xa9, 0x14}, // P2PKH
				}},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.tx.Validate()

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

// TestOutPoint_String tests OutPoint string representation
func TestOutPoint_String(t *testing.T) {
	hash, _ := NewHash256FromString("000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f")
	outpoint := OutPoint{Hash: hash, Index: 0}

	expected := "000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f:0"
	result := outpoint.String()

	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

// TestOutPoint_IsNull tests null outpoint detection
func TestOutPoint_IsNull(t *testing.T) {
	tests := []struct {
		name     string
		outpoint OutPoint
		expected bool
	}{
		{
			name:     "null outpoint (coinbase)",
			outpoint: OutPoint{Hash: ZeroHash, Index: 0xffffffff},
			expected: true,
		},
		{
			name:     "non-null outpoint",
			outpoint: OutPoint{Hash: Hash256{0x01}, Index: 0},
			expected: false,
		},
		{
			name:     "zero hash but wrong index",
			outpoint: OutPoint{Hash: ZeroHash, Index: 0},
			expected: false,
		},
		{
			name:     "correct index but non-zero hash",
			outpoint: OutPoint{Hash: Hash256{0x01}, Index: 0xffffffff},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.outpoint.IsNull()
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestTransaction_Hash tests transaction hashing (currently returns zero)
func TestTransaction_Hash(t *testing.T) {
	tx := &Transaction{
		Version: 1,
		Inputs: []TxInput{{
			PreviousOutput: OutPoint{Hash: Hash256{0x01}, Index: 0},
			ScriptSig:      []byte{0x76, 0xa9},
			Sequence:       0xffffffff,
		}},
		Outputs: []TxOutput{{
			Value:        1000000000,
			ScriptPubKey: []byte{0x76, 0xa9, 0x14},
		}},
		LockTime: 0,
	}

	hash1 := tx.Hash()
	hash2 := tx.Hash()

	// Hash should be cached and consistent
	if hash1 != hash2 {
		t.Errorf("hash not consistent: %s != %s", hash1.String(), hash2.String())
	}

	// Currently returns zero hash (TODO: implement actual hashing)
	if !hash1.IsZero() {
		t.Logf("Hash implementation complete: %s", hash1.String())
	}
}

// TestTransaction_WitnessHash tests witness transaction hashing (currently returns zero)
func TestTransaction_WitnessHash(t *testing.T) {
	tx := &Transaction{
		Version: 2,
		Inputs: []TxInput{{
			PreviousOutput: OutPoint{Hash: Hash256{0x01}, Index: 0},
			ScriptSig:      []byte{},
			Sequence:       0xffffffff,
		}},
		Outputs: []TxOutput{{
			Value:        1000000000,
			ScriptPubKey: []byte{0x00, 0x14}, // P2WPKH
		}},
		Witnesses: []TxWitness{{
			Stack: [][]byte{[]byte("signature"), []byte("pubkey")},
		}},
		LockTime: 0,
	}

	wHash1 := tx.WitnessHash()
	wHash2 := tx.WitnessHash()

	// Hash should be cached and consistent
	if wHash1 != wHash2 {
		t.Errorf("witness hash not consistent: %s != %s", wHash1.String(), wHash2.String())
	}

	// Currently returns zero hash (TODO: implement actual hashing)
	if !wHash1.IsZero() {
		t.Logf("Witness hash implementation complete: %s", wHash1.String())
	}
}

// TestConstants tests Bitcoin constants
func TestConstants(t *testing.T) {
	// Test MaxMoney constant
	expectedMaxMoney := uint64(21000000 * 100000000) // 21 million BTC in satoshis
	if MaxMoney != expectedMaxMoney {
		t.Errorf("MaxMoney constant incorrect: expected %d, got %d", expectedMaxMoney, MaxMoney)
	}

	// Verify MaxMoney is exactly 21 million BTC
	maxBTC := float64(MaxMoney) / 100000000.0
	if maxBTC != 21000000.0 {
		t.Errorf("MaxMoney should equal 21 million BTC, got %.8f", maxBTC)
	}
}

// BenchmarkTransaction_TotalOutput benchmarks output summation
func BenchmarkTransaction_TotalOutput(b *testing.B) {
	// Create transaction with many outputs
	outputs := make([]TxOutput, 1000)
	for i := range outputs {
		outputs[i] = TxOutput{Value: uint64(i + 1000000)}
	}

	tx := &Transaction{Outputs: outputs}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tx.TotalOutput()
	}
}

// BenchmarkTransaction_Validate benchmarks transaction validation
func BenchmarkTransaction_Validate(b *testing.B) {
	// Create a valid transaction
	validOutpoint := OutPoint{Hash: Hash256{0x01}, Index: 0}
	tx := &Transaction{
		Version: 1,
		Inputs: []TxInput{{
			PreviousOutput: validOutpoint,
			ScriptSig:      []byte{0x76, 0xa9},
			Sequence:       0xffffffff,
		}},
		Outputs: []TxOutput{{
			Value:        1000000000,
			ScriptPubKey: []byte{0x76, 0xa9, 0x14},
		}},
		LockTime: 0,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tx.Validate()
	}
}