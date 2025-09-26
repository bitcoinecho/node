# Bitcoin Echo Node

A Pure Bitcoin Node Implementation

## Overview

Bitcoin Echo aims to provide:
- **Protocol Fidelity**: Exact implementation of Bitcoin consensus rules
- **User Sovereignty**: All policy decisions are user-configurable
- **Universal Compatibility**: Seamless interaction with the Bitcoin network

## Status

This project is currently in early development. See [WHITEPAPER.md](./WHITEPAPER.md) for the complete architectural specification and philosophy.

## Quick Start

```bash
# Build the node
go build ./cmd/bitcoin-echo

# Run tests
go test ./...

# Clean dependencies
go mod tidy
```

## Architecture

Bitcoin Echo is structured with clear separation between consensus rules (immutable) and policy settings (user-configurable):

- **Consensus Engine** - Block/transaction validation, script execution
- **Network Layer** - P2P protocol and peer management
- **Storage Engine** - Blockchain data and UTXO set
- **Mempool** - Transaction validation with configurable policies
- **RPC Server** - Standard Bitcoin RPC API
- **Web UI** - Optional dashboard and configuration

## Documentation

- [WHITEPAPER.md](./WHITEPAPER.md) - Complete project specification
- [Website](https://bitcoinecho.org) - Project website

## License

MIT