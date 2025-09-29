package bitcoin

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// Hash256 represents a 256-bit hash (32 bytes)
type Hash256 [32]byte

// ZeroHash represents an all-zero hash
var ZeroHash = Hash256{}

// NewHash256FromBytes creates a Hash256 from a byte slice
func NewHash256FromBytes(b []byte) (Hash256, error) {
	if len(b) != 32 {
		return ZeroHash, fmt.Errorf("invalid hash length: expected 32 bytes, got %d", len(b))
	}
	var hash Hash256
	copy(hash[:], b)
	return hash, nil
}

// NewHash256FromString creates a Hash256 from a hex string
func NewHash256FromString(s string) (Hash256, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return ZeroHash, fmt.Errorf("invalid hex string: %v", err)
	}
	return NewHash256FromBytes(b)
}

// String returns the hash as a hex string
func (h Hash256) String() string {
	return hex.EncodeToString(h[:])
}

// Bytes returns the hash as a byte slice
func (h Hash256) Bytes() []byte {
	return h[:]
}

// IsZero returns true if the hash is all zeros
func (h Hash256) IsZero() bool {
	return h == ZeroHash
}

// DoubleHashSHA256 performs double SHA256 hashing (SHA256(SHA256(data)))
func DoubleHashSHA256(data []byte) Hash256 {
	first := sha256.Sum256(data)
	second := sha256.Sum256(first[:])
	return Hash256(second)
}

// Hash160 represents a 160-bit hash (20 bytes) used for addresses
type Hash160 [20]byte

// ZeroHash160 represents an all-zero hash160
var ZeroHash160 = Hash160{}

// NewHash160FromBytes creates a Hash160 from a byte slice
func NewHash160FromBytes(b []byte) (Hash160, error) {
	if len(b) != 20 {
		return ZeroHash160, fmt.Errorf("invalid hash160 length: expected 20 bytes, got %d", len(b))
	}
	var hash Hash160
	copy(hash[:], b)
	return hash, nil
}

// String returns the hash160 as a hex string
func (h Hash160) String() string {
	return hex.EncodeToString(h[:])
}

// Bytes returns the hash160 as a byte slice
func (h Hash160) Bytes() []byte {
	return h[:]
}
