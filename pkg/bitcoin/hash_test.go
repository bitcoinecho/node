package bitcoin

import (
	"testing"
)

// TestHash256_NewHash256FromBytes tests creating Hash256 from byte slices
func TestHash256_NewHash256FromBytes(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		expectError bool
		expected    string // hex representation
	}{
		{
			name:     "valid 32-byte input",
			input:    make([]byte, 32),
			expected: "0000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			name:        "too short input",
			input:       make([]byte, 31),
			expectError: true,
		},
		{
			name:        "too long input",
			input:       make([]byte, 33),
			expectError: true,
		},
		{
			name:     "all 0xFF bytes",
			input:    []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			expected: "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		},
		{
			name:        "empty input",
			input:       []byte{},
			expectError: true,
		},
		{
			name:        "nil input",
			input:       nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := NewHash256FromBytes(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if hash.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, hash.String())
			}
		})
	}
}

// TestHash256_NewHash256FromString tests creating Hash256 from hex strings
func TestHash256_NewHash256FromString(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		expected    string
	}{
		{
			name:     "valid zero hash",
			input:    "0000000000000000000000000000000000000000000000000000000000000000",
			expected: "0000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			name:     "valid all FF hash",
			input:    "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			expected: "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		},
		{
			name:     "valid mixed hash",
			input:    "000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f",
			expected: "000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f",
		},
		{
			name:        "too short hex string",
			input:       "00112233445566778899aabbccddeeff00112233445566778899aabbccdde",
			expectError: true,
		},
		{
			name:        "too long hex string",
			input:       "000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f00",
			expectError: true,
		},
		{
			name:        "invalid hex characters",
			input:       "000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce2zz",
			expectError: true,
		},
		{
			name:        "odd length hex string",
			input:       "000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26",
			expectError: true,
		},
		{
			name:        "empty string",
			input:       "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := NewHash256FromString(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if hash.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, hash.String())
			}
		})
	}
}

// TestHash256_IsZero tests the IsZero method
func TestHash256_IsZero(t *testing.T) {
	tests := []struct {
		name     string
		hash     Hash256
		expected bool
	}{
		{
			name:     "zero hash",
			hash:     ZeroHash,
			expected: true,
		},
		{
			name:     "non-zero hash",
			hash:     Hash256{0x01},
			expected: false,
		},
		{
			name:     "all FF hash",
			hash:     Hash256{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.hash.IsZero()
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestHash256_Bytes tests the Bytes method
func TestHash256_Bytes(t *testing.T) {
	testHash := Hash256{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20}

	bytes := testHash.Bytes()

	if len(bytes) != 32 {
		t.Errorf("expected 32 bytes, got %d", len(bytes))
	}

	for i, b := range bytes {
		if b != testHash[i] {
			t.Errorf("byte mismatch at index %d: expected %02x, got %02x", i, testHash[i], b)
		}
	}
}

// TestDoubleHashSHA256 tests the Bitcoin double SHA-256 hashing function
func TestDoubleHashSHA256(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "empty input",
			input:    []byte{},
			expected: "5df6e0e2761359d30a8275058e299fcc0381534545f55cf43e41983f5d4c9456", // SHA256(SHA256(""))
		},
		{
			name:     "single byte",
			input:    []byte{0x00},
			expected: "1406e05881e299367766d313e26c05564ec91bf721d31726bd6e46e60689539a", // SHA256(SHA256([0x00]))
		},
		{
			name:     "hello world",
			input:    []byte("hello world"),
			expected: "bc62d4b80d9e36da29c16c5d4d9f11731f36052c72401a76c23c0fb5a9b74423", // SHA256(SHA256("hello world"))
		},
		{
			name:     "bitcoin genesis block coinbase",
			input:    []byte("The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"),
			expected: "687c09c2b4c2392a47717f58c468698b998fef0eed2ec9c8f8736d42a1b8c26a", // Verified double SHA256
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DoubleHashSHA256(tt.input)
			resultHex := result.String()

			if resultHex != tt.expected {
				t.Errorf("DoubleHashSHA256(%q) = %s; expected %s", tt.input, resultHex, tt.expected)
			}
		})
	}
}

// TestDoubleHashSHA256_KnownBitcoinValues tests against known Bitcoin values
func TestDoubleHashSHA256_KnownBitcoinValues(t *testing.T) {
	// Test the empty input case - this is a well-known value in Bitcoin
	// SHA256(SHA256("")) = 5df6e0e2761359d30a8275058e299fcc0381534545f55cf43e41983f5d4c9456
	emptyHash := DoubleHashSHA256([]byte{})
	expectedEmpty := "5df6e0e2761359d30a8275058e299fcc0381534545f55cf43e41983f5d4c9456"

	if emptyHash.String() != expectedEmpty {
		t.Errorf("Double SHA256 of empty input failed: got %s, expected %s", emptyHash.String(), expectedEmpty)
	}
}

// TestHash160_NewHash160FromBytes tests creating Hash160 from byte slices
func TestHash160_NewHash160FromBytes(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		expectError bool
		expected    string
	}{
		{
			name:     "valid 20-byte input",
			input:    make([]byte, 20),
			expected: "0000000000000000000000000000000000000000",
		},
		{
			name:        "too short input",
			input:       make([]byte, 19),
			expectError: true,
		},
		{
			name:        "too long input",
			input:       make([]byte, 21),
			expectError: true,
		},
		{
			name:     "all 0xFF bytes",
			input:    []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			expected: "ffffffffffffffffffffffffffffffffffffffff",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := NewHash160FromBytes(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if hash.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, hash.String())
			}
		})
	}
}

// TestHash160_Bytes tests the Hash160 Bytes method
func TestHash160_Bytes(t *testing.T) {
	testHash := Hash160{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14}

	bytes := testHash.Bytes()

	if len(bytes) != 20 {
		t.Errorf("expected 20 bytes, got %d", len(bytes))
	}

	for i, b := range bytes {
		if b != testHash[i] {
			t.Errorf("byte mismatch at index %d: expected %02x, got %02x", i, testHash[i], b)
		}
	}
}

// BenchmarkDoubleHashSHA256 benchmarks the double SHA256 function
func BenchmarkDoubleHashSHA256(b *testing.B) {
	data := []byte("The Times 03/Jan/2009 Chancellor on brink of second bailout for banks")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = DoubleHashSHA256(data)
	}
}

// BenchmarkNewHash256FromString benchmarks creating Hash256 from string
func BenchmarkNewHash256FromString(b *testing.B) {
	hashStr := "000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = NewHash256FromString(hashStr)
	}
}