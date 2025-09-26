# Bitcoin Echo Node

A Pure Bitcoin Node Implementation

## Overview

Bitcoin Echo aims to provide:
- **Protocol Fidelity**: Exact implementation of Bitcoin consensus rules
- **User Sovereignty**: All policy decisions are user-configurable
- **Universal Compatibility**: Seamless interaction with the Bitcoin network

## Status

**Current Development Phase**: Test-Driven Implementation (September 2025)

✅ **Foundation Complete:**
- Core Bitcoin types (Hash256, Transaction, Block, Script)
- Comprehensive test suite with TDD approach
- Cryptographic functions (SHA-256 double hashing)
- Transaction validation with Bitcoin consensus rules
- **Block header hashing with Genesis Block verification**

🧪 **Test Coverage:**
- Hash functions: 100% coverage
- Transaction validation: 100% coverage
- Block operations: ~97% coverage
- **Overall project: 60.2% coverage achieved**

🚧 **Next Implementation Priorities:**
- Script execution engine (Bitcoin Script interpreter)
- Transaction serialization and witness data
- P2P networking layer

See [WHITEPAPER.md](./WHITEPAPER.md) for the complete architectural specification and philosophy.

## Quick Start

```bash
# Build the node
go build ./cmd/bitcoin-echo

# Run all tests with coverage
go test -v ./...
go test -cover ./pkg/bitcoin

# Run specific test suites
go test -v ./pkg/bitcoin -run TestHash
go test -v ./pkg/bitcoin -run TestTransaction
go test -v ./pkg/bitcoin -run TestBlock

# Clean dependencies
go mod tidy
```

### Development Workflow

Bitcoin Echo follows **Test-Driven Development (TDD)**:

1. **Red**: Write failing tests for desired functionality
2. **Green**: Implement minimal code to make tests pass
3. **Refactor**: Improve code while keeping tests green

All consensus-critical code requires 100% test coverage before implementation.

## Architecture

Bitcoin Echo is structured with clear separation between consensus rules (immutable) and policy settings (user-configurable):

- **Consensus Engine** - Block/transaction validation, script execution
- **Network Layer** - P2P protocol and peer management
- **Storage Engine** - Blockchain data and UTXO set
- **Mempool** - Transaction validation with configurable policies
- **RPC Server** - Standard Bitcoin RPC API
- **Web UI** - Optional dashboard and configuration

## Implementation Status

### ✅ Completed Components

**Cryptographic Foundation:**
- `Hash256` - Bitcoin's 256-bit hash type with full validation
- `DoubleHashSHA256` - Bitcoin's standard double SHA-256 function
- `Hash160` - 160-bit hash type for Bitcoin addresses

**Transaction System:**
- Complete transaction structure (inputs, outputs, witness data)
- Coinbase transaction detection and validation
- Input/output validation with Bitcoin money limits (21M BTC)
- Duplicate input detection and basic consensus checks

**Block System:**
- Bitcoin-compliant block header serialization (80 bytes)
- Double SHA-256 block hashing with correct byte ordering
- ✅ **Genesis Block verification**: `000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f`
- Block validation with coinbase and transaction checks
- Block size/weight limit enforcement (1MB/4M weight units)

### 🚧 In Development

**Script Engine:**
- Bitcoin Script interpreter and execution environment
- Standard script types (P2PKH, P2SH, P2WPKH, P2WSH, Taproot)
- Signature verification (ECDSA, Schnorr)

**Enhanced Features:**
- Merkle tree construction and validation
- Proof-of-work verification
- Transaction serialization with witness data

## Documentation

- [WHITEPAPER.md](./WHITEPAPER.md) - Complete project specification
- [Website](https://bitcoinecho.org) - Project website

## License

MIT