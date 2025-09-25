package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bitcoinecho/node/pkg/bitcoin"
)

const (
	Name    = "bitcoin-echo"
	Version = "0.1.0-dev"
)

func main() {
	fmt.Printf("%s v%s\n", Name, Version)
	fmt.Println("A Pure Bitcoin Node Implementation")
	fmt.Println("")

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version":
			printVersion()
		case "help":
			printHelp()
		case "test":
			runTests()
		default:
			fmt.Printf("Unknown command: %s\n", os.Args[1])
			printHelp()
			os.Exit(1)
		}
	} else {
		// Default: start the node
		startNode()
	}
}

func printVersion() {
	fmt.Printf("%s version %s\n", Name, Version)
	fmt.Println("Built with Go")
	fmt.Println("")
	fmt.Println("Bitcoin Echo: Faithfully reflecting the Bitcoin protocol since 2025")
}

func printHelp() {
	fmt.Printf("Usage: %s [command]\n", Name)
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  help        Show this help message")
	fmt.Println("  version     Show version information")
	fmt.Println("  test        Run basic functionality tests")
	fmt.Println("  (no args)   Start the Bitcoin Echo node")
	fmt.Println("")
	fmt.Println("For more information, visit: https://bitcoinecho.org")
}

func startNode() {
	fmt.Println("🚀 Starting Bitcoin Echo node...")
	fmt.Println("")

	// TODO: Implement full node startup
	fmt.Println("⚠️  Node implementation in progress")
	fmt.Println("📋 Current status: Core types defined")
	fmt.Println("")

	// For now, demonstrate that our types work
	demonstrateTypes()

	fmt.Println("Node would continue running here...")
	fmt.Println("Use Ctrl+C to stop")
}

func runTests() {
	fmt.Println("🧪 Running basic functionality tests...")
	fmt.Println("")

	demonstrateTypes()

	fmt.Println("✅ Basic tests completed")
}

func demonstrateTypes() {
	// Create a sample transaction
	fmt.Println("📦 Creating sample transaction...")

	// Create some dummy inputs and outputs
	prevHash, err := bitcoin.NewHash256FromString("0000000000000000000000000000000000000000000000000000000000000000")
	if err != nil {
		log.Printf("Error creating hash: %v", err)
		return
	}

	outpoint := bitcoin.OutPoint{
		Hash:  prevHash,
		Index: 0,
	}

	input := bitcoin.TxInput{
		PreviousOutput: outpoint,
		ScriptSig:      []byte{0x76, 0xa9, 0x14}, // Dummy script
		Sequence:       0xffffffff,
	}

	output := bitcoin.TxOutput{
		Value:        5000000000, // 50 BTC in satoshis
		ScriptPubKey: []byte{0x76, 0xa9, 0x14}, // Dummy P2PKH script
	}

	tx := bitcoin.NewTransaction(1, []bitcoin.TxInput{input}, []bitcoin.TxOutput{output}, 0)

	fmt.Printf("   Transaction ID: %s\n", tx.Hash().String())
	fmt.Printf("   Is Coinbase: %t\n", tx.IsCoinbase())
	fmt.Printf("   Output Value: %d satoshis\n", tx.TotalOutput())

	// Validate the transaction
	if err := tx.Validate(); err != nil {
		fmt.Printf("   ⚠️ Transaction validation failed: %v\n", err)
	} else {
		fmt.Printf("   ✅ Transaction validation passed\n")
	}

	fmt.Println("")

	// Create a sample block
	fmt.Println("🧱 Creating sample block...")

	header := bitcoin.NewBlockHeader(
		1,                    // Version
		bitcoin.ZeroHash,     // Previous block hash (genesis)
		bitcoin.ZeroHash,     // Merkle root (placeholder)
		1640995200,          // Timestamp (Jan 1, 2022)
		0x1d00ffff,          // Bits (difficulty)
		12345,               // Nonce
	)

	block := bitcoin.NewBlock(header, []bitcoin.Transaction{*tx})

	fmt.Printf("   Block Hash: %s\n", block.Hash().String())
	fmt.Printf("   Is Genesis: %t\n", block.IsGenesis())
	fmt.Printf("   Transaction Count: %d\n", block.TransactionCount())
	fmt.Printf("   Has Coinbase: %t\n", block.HasCoinbase())

	// Validate the block
	if err := block.Validate(); err != nil {
		fmt.Printf("   ⚠️ Block validation failed: %v\n", err)
	} else {
		fmt.Printf("   ✅ Block validation passed\n")
	}

	fmt.Println("")

	// Demonstrate script analysis
	fmt.Println("📜 Analyzing sample scripts...")

	// P2PKH script
	p2pkhScript := bitcoin.Script{0x76, 0xa9, 0x14, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x88, 0xac}
	fmt.Printf("   P2PKH Script Type: %v\n", p2pkhScript.AnalyzeScript())
	fmt.Printf("   P2PKH Is Standard: %t\n", p2pkhScript.IsStandard())

	// OP_RETURN script
	opReturnScript := bitcoin.Script{0x6a, 0x0b, 'H', 'e', 'l', 'l', 'o', ' ', 'W', 'o', 'r', 'l', 'd'}
	fmt.Printf("   OP_RETURN Script Type: %v\n", opReturnScript.AnalyzeScript())
	fmt.Printf("   OP_RETURN Is Standard: %t\n", opReturnScript.IsStandard())

	fmt.Println("")
}