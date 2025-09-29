package bitcoin

import (
	"fmt"
)

// UTXO represents an Unspent Transaction Output
// TDD GREEN: Basic implementation to make tests pass
type UTXO struct {
	txHash      Hash256
	outputIndex uint32
	amount      uint64
	scriptPubKey []byte
}

// NewUTXO creates a new UTXO
func NewUTXO(txHash Hash256, outputIndex uint32, amount uint64, scriptPubKey []byte) *UTXO {
	script := make([]byte, len(scriptPubKey))
	copy(script, scriptPubKey)
	return &UTXO{
		txHash:      txHash,
		outputIndex: outputIndex,
		amount:      amount,
		scriptPubKey: script,
	}
}

// TxHash returns the transaction hash
func (u *UTXO) TxHash() Hash256 {
	return u.txHash
}

// OutputIndex returns the output index
func (u *UTXO) OutputIndex() uint32 {
	return u.outputIndex
}

// Amount returns the amount in satoshis
func (u *UTXO) Amount() uint64 {
	return u.amount
}

// ScriptPubKey returns the script public key
func (u *UTXO) ScriptPubKey() []byte {
	return u.scriptPubKey
}

// UTXOSet represents a set of unspent transaction outputs
// TDD GREEN: Basic implementation using map for fast lookups
type UTXOSet struct {
	utxos map[string]*UTXO // key: txHash:outputIndex
}

// NewUTXOSet creates a new UTXO set
func NewUTXOSet() *UTXOSet {
	return &UTXOSet{
		utxos: make(map[string]*UTXO),
	}
}

// makeKey creates a unique key for UTXO indexing
func (s *UTXOSet) makeKey(txHash Hash256, outputIndex uint32) string {
	return fmt.Sprintf("%s:%d", txHash.String(), outputIndex)
}

// Add adds a UTXO to the set
func (s *UTXOSet) Add(utxo *UTXO) {
	key := s.makeKey(utxo.txHash, utxo.outputIndex)
	s.utxos[key] = utxo
}

// Remove removes a UTXO from the set
func (s *UTXOSet) Remove(txHash Hash256, outputIndex uint32) bool {
	key := s.makeKey(txHash, outputIndex)
	if _, exists := s.utxos[key]; exists {
		delete(s.utxos, key)
		return true
	}
	return false
}

// Find finds a UTXO in the set
func (s *UTXOSet) Find(txHash Hash256, outputIndex uint32) (*UTXO, bool) {
	key := s.makeKey(txHash, outputIndex)
	utxo, exists := s.utxos[key]
	return utxo, exists
}

// Size returns the number of UTXOs in the set
func (s *UTXOSet) Size() int {
	return len(s.utxos)
}

// ValidateSpend validates if a UTXO can be spent
// TDD GREEN: Basic validation to check if UTXO exists and has sufficient amount
func (s *UTXOSet) ValidateSpend(txHash Hash256, outputIndex uint32, amount uint64) bool {
	utxo, exists := s.Find(txHash, outputIndex)
	if !exists {
		return false // UTXO doesn't exist
	}

	// For basic validation, just check if UTXO exists
	// In full implementation, would also check script validation, etc.
	return utxo.amount >= amount
}

// TotalValue calculates the total value of all UTXOs in the set
func (s *UTXOSet) TotalValue() uint64 {
	total := uint64(0)
	for _, utxo := range s.utxos {
		total += utxo.amount
	}
	return total
}

// GetAllUTXOs returns all UTXOs in the set
func (s *UTXOSet) GetAllUTXOs() []*UTXO {
	utxos := make([]*UTXO, 0, len(s.utxos))
	for _, utxo := range s.utxos {
		utxos = append(utxos, utxo)
	}
	return utxos
}

// Clear removes all UTXOs from the set
func (s *UTXOSet) Clear() {
	s.utxos = make(map[string]*UTXO)
}