package bitcoin

import (
	"fmt"
)

// Transaction represents a Bitcoin transaction
type Transaction struct {
	Version    uint32     `json:"version"`
	Inputs     []TxInput  `json:"inputs"`
	Outputs    []TxOutput `json:"outputs"`
	LockTime   uint32     `json:"locktime"`

	// Witness data for SegWit transactions
	Witnesses []TxWitness `json:"witnesses,omitempty"`

	// Cached values
	hash   *Hash256 // Transaction ID (excludes witness data)
	wthash *Hash256 // Witness Transaction ID (includes witness data)
}

// TxInput represents a transaction input
type TxInput struct {
	PreviousOutput OutPoint `json:"previous_output"`
	ScriptSig      []byte   `json:"script_sig"`
	Sequence       uint32   `json:"sequence"`
}

// TxOutput represents a transaction output
type TxOutput struct {
	Value        uint64 `json:"value"`        // Amount in satoshis
	ScriptPubKey []byte `json:"script_pubkey"`
}

// TxWitness represents witness data for a SegWit input
type TxWitness struct {
	Stack [][]byte `json:"stack"`
}

// OutPoint represents a reference to a transaction output
type OutPoint struct {
	Hash  Hash256 `json:"hash"`  // Transaction hash
	Index uint32  `json:"index"` // Output index
}

// NewTransaction creates a new transaction
func NewTransaction(version uint32, inputs []TxInput, outputs []TxOutput, lockTime uint32) *Transaction {
	return &Transaction{
		Version:  version,
		Inputs:   inputs,
		Outputs:  outputs,
		LockTime: lockTime,
	}
}

// Hash returns the transaction ID (excludes witness data)
func (tx *Transaction) Hash() Hash256 {
	if tx.hash == nil {
		// TODO: Implement transaction serialization and hashing
		// For now, return zero hash
		hash := ZeroHash
		tx.hash = &hash
	}
	return *tx.hash
}

// WitnessHash returns the witness transaction ID (includes witness data)
func (tx *Transaction) WitnessHash() Hash256 {
	if tx.wthash == nil {
		// TODO: Implement witness transaction serialization and hashing
		// For now, return zero hash
		hash := ZeroHash
		tx.wthash = &hash
	}
	return *tx.wthash
}

// IsCoinbase returns true if this is a coinbase transaction
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 &&
		   tx.Inputs[0].PreviousOutput.Hash.IsZero() &&
		   tx.Inputs[0].PreviousOutput.Index == 0xffffffff
}

// HasWitness returns true if the transaction has witness data
func (tx *Transaction) HasWitness() bool {
	return len(tx.Witnesses) > 0
}

// TotalOutput calculates the total value of all outputs
func (tx *Transaction) TotalOutput() uint64 {
	var total uint64
	for _, output := range tx.Outputs {
		total += output.Value
	}
	return total
}

// IsStandard checks if the transaction follows standard rules
func (tx *Transaction) IsStandard() bool {
	// TODO: Implement standard transaction checks
	// - Version check
	// - Size limits
	// - Standard script types
	// - Dust threshold
	return true // Placeholder
}

// Validate performs basic validation checks
func (tx *Transaction) Validate() error {
	// Basic sanity checks
	if len(tx.Inputs) == 0 {
		return fmt.Errorf("transaction has no inputs")
	}

	if len(tx.Outputs) == 0 {
		return fmt.Errorf("transaction has no outputs")
	}

	// Check for duplicate inputs
	seen := make(map[OutPoint]bool)
	for _, input := range tx.Inputs {
		if seen[input.PreviousOutput] {
			return fmt.Errorf("transaction has duplicate inputs")
		}
		seen[input.PreviousOutput] = true
	}

	// Check output values
	for i, output := range tx.Outputs {
		if output.Value > MaxMoney {
			return fmt.Errorf("output %d value exceeds maximum", i)
		}
	}

	// Check total output value
	if tx.TotalOutput() > MaxMoney {
		return fmt.Errorf("total output value exceeds maximum")
	}

	return nil
}

// Constants
const (
	MaxMoney = 21000000 * 100000000 // 21 million BTC in satoshis
)

// String returns a string representation of the OutPoint
func (op OutPoint) String() string {
	return fmt.Sprintf("%s:%d", op.Hash.String(), op.Index)
}

// IsNull returns true if the outpoint is null (coinbase)
func (op OutPoint) IsNull() bool {
	return op.Hash.IsZero() && op.Index == 0xffffffff
}