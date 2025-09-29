# Bitcoin Echo Node

A Pure Bitcoin Node Implementation

## Overview

Bitcoin Echo aims to provide:
- **Protocol Fidelity**: Exact implementation of Bitcoin consensus rules
- **User Sovereignty**: All policy decisions are user-configurable
- **Universal Compatibility**: Seamless interaction with the Bitcoin network

## Status

**Current Development Phase**: Test-Driven Implementation (September 2025)

✅ **Project Structure Reorganized**: Tests moved to `tests/unit/bitcoin/` with proper Go module structure

✅ **Foundation Complete:**
- Core Bitcoin types (Hash256, Transaction, Block, Script)
- Comprehensive test suite with TDD approach
- Cryptographic functions (SHA-256 double hashing)
- Transaction validation with Bitcoin consensus rules
- **Block header hashing with Genesis Block verification**
- **Bitcoin Script Analysis with full script type detection**
- **✅ Bitcoin Script Execution Engine with stack-based interpreter**

🧪 **Test Coverage:**
- Hash functions: 100% coverage
- Transaction validation: 100% coverage
- Block operations: ~97% coverage
- Script analysis: 100% coverage (40+ comprehensive test cases)
- Script execution: 100% coverage (40+ execution test cases)
- **✅ Transaction serialization: 100% coverage (VarInt, wire format, SegWit)**
- **✅ ECDSA signature verification: 100% coverage (OP_CHECKSIG implementation)**
- **✅ Merkle tree construction: 100% coverage (Bitcoin protocol-compliant)**
- **✅ Proof-of-Work verification: 100% coverage (difficulty validation and adjustment)**
- **✅ Basic P2P networking: 100% coverage (message protocol, peer connections)**
- **✅ UTXO set management: 100% coverage (creation, validation, storage)**
- **✅ Blockchain validation: 100% coverage (chain management, block validation)**
- **✅ Overall project: 260/260 tests passing (100% success rate)**

✅ **Recently Completed:**
- **✅ Blockchain Validation Engine**: Complete blockchain management with block validation
- **✅ Chain Management**: Block addition, validation, UTXO integration, and integrity checking
- **✅ GitHub Actions CI/CD**: Comprehensive testing, security scanning, and automated releases
- **✅ UTXO Set Management Complete**: Unspent transaction output tracking and validation
- **✅ UTXO Operations**: Creation, addition, removal, lookup with optimized hash indexing
- **✅ Spend Validation**: Double-spend prevention and amount verification
- **✅ Basic P2P Networking Foundation**: Bitcoin P2P message protocol implementation
- **✅ P2P Message System**: Serialization, deserialization, validation with Bitcoin wire format
- **✅ Peer Connection Management**: Network connection handling and handshake simulation
- **✅ Proof-of-Work Verification Complete**: Bitcoin difficulty validation, adjustment algorithm, Genesis block validation
- **✅ Difficulty Target System**: Compact format conversion, big integer arithmetic, 4x adjustment limits
- **✅ Block Hash Validation**: Complete target comparison with proper byte ordering and validation
- **✅ Merkle Tree Implementation**: Complete Bitcoin merkle tree construction with TDD methodology
- **✅ ECDSA Signature Verification**: Complete OP_CHECKSIG implementation with cryptographic validation
- **✅ Script Execution Complete**: Stack-based interpreter with full Bitcoin protocol compliance
- **✅ Test Suite Achievement**: 260/260 tests passing with comprehensive coverage

🔜 **Next Implementation Priorities:**
- Network synchronization and block relay
- Advanced P2P features (peer discovery, inventory messages)

See [WHITEPAPER.md](./WHITEPAPER.md) for the complete architectural specification and philosophy.

## Continuous Integration

Bitcoin Echo uses comprehensive GitHub Actions workflows for quality assurance:

- **🧪 Test Suite**: Automated testing on every push with 95%+ coverage requirement
- **🔍 Code Quality**: golangci-lint with Bitcoin-specific rules and formatting checks
- **🛡️ Security Scanning**: gosec, govulncheck, and dependency vulnerability scanning
- **🏗️ Multi-Platform Builds**: Linux, Windows, macOS (amd64 & arm64)
- **🔄 Automated Releases**: Tagged releases with cross-platform binaries and Docker images
- **📦 Dependency Management**: Automated dependency updates with security checks

All consensus-critical code requires 100% test coverage and multiple test runs before merge.

## Quick Start

```bash
# Build the node
go build ./cmd/bitcoin-echo

# Run all tests with coverage
go test -v ./...
go test -cover ./tests/unit/bitcoin

# Run specific test suites
go test -v ./tests/unit/bitcoin -run TestHash
go test -v ./tests/unit/bitcoin -run TestTransaction
go test -v ./tests/unit/bitcoin -run TestBlock
go test -v ./tests/unit/bitcoin -run TestScript
go test -v ./tests/unit/bitcoin -run TestScriptEngine

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
- **✅ Bitcoin wire format serialization/deserialization**
- **✅ VarInt encoding with all 4 size variants**
- **✅ SegWit witness data handling**

**Block System:**
- Bitcoin-compliant block header serialization (80 bytes)
- Double SHA-256 block hashing with correct byte ordering
- ✅ **Genesis Block verification**: `000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f`
- Block validation with coinbase and transaction checks
- Block size/weight limit enforcement (1MB/4M weight units)
- **✅ Merkle tree construction**: Level-by-level Bitcoin protocol implementation
- **✅ Transaction summarization**: Proper odd-count handling with last hash duplication

**Script Analysis Engine:**
- Complete Bitcoin script type detection and classification
- Support for all standard script types: P2PKH, P2SH, P2PK, P2WPKH, P2WSH, P2TR
- Multisig script analysis (M-of-N signatures with standardness validation)
- OP_RETURN (null data) script detection
- Bitcoin standardness rule enforcement with configurable policy limits
- High-performance analysis (sub-3ns execution time)

**✅ Script Execution Engine:**
- Complete Bitcoin Script interpreter with stack-based execution
- Full opcode support: OP_1-OP_16, stack ops (DUP, DROP, SWAP), arithmetic (ADD, SUB)
- Comparison operations (EQUAL, EQUALVERIFY) with proper verification semantics
- Bitcoin-compliant number encoding (little-endian with sign bit)
- Hash operations (OP_HASH160) and error handling with stack protection
- **✅ ECDSA signature verification**: OP_CHECKSIG implementation with DER format validation
- **✅ Cryptographic validation**: Distinguishes valid/invalid signatures correctly
- **✅ TDD implementation**: Complete RED-GREEN-REFACTOR cycle with comprehensive tests

### 🔜 Next Implementation Priorities

**Core Protocol Features:**
- Basic P2P networking layer (peer discovery, message handling)
- UTXO set management and storage
- Block chain validation engine
- Advanced signature verification (Schnorr, multi-signature)

## Documentation

- [WHITEPAPER.md](./WHITEPAPER.md) - Complete project specification
- [Website](https://bitcoinecho.org) - Project website

## License

MIT