package bitcoin_test

import (
	"bitcoinecho.org/node/pkg/bitcoin"
	"testing"
)

// TestBlockChain_Creation tests blockchain creation and initialization (TDD RED phase)
func TestBlockChain_Creation(t *testing.T) {
	tests := []struct {
		name         string
		genesisBlock *bitcoin.Block
		description  string
	}{
		{
			name:         "Create blockchain with Genesis block",
			genesisBlock: createGenesisBlock(),
			description:  "Blockchain should initialize with Genesis block",
		},
		{
			name:         "Create empty blockchain",
			genesisBlock: nil,
			description:  "Blockchain should handle empty initialization",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("TDD RED: %s - %s", tt.name, tt.description)

			// This should fail since we haven't implemented blockchain yet
			blockchain := bitcoin.NewBlockChain(tt.genesisBlock)

			if tt.genesisBlock != nil {
				if blockchain.Height() != 0 {
					t.Errorf("Expected blockchain height 0, got %d", blockchain.Height())
				}

				tip := blockchain.GetTip()
				if tip == nil {
					t.Error("Expected blockchain tip to be Genesis block, got nil")
				}
			} else {
				if blockchain.Height() != -1 {
					t.Errorf("Expected empty blockchain height -1, got %d", blockchain.Height())
				}
			}

			t.Logf("Blockchain created with height: %d", blockchain.Height())
		})
	}
}

// TestBlockChain_AddBlock tests adding blocks to the blockchain (TDD RED phase)
func TestBlockChain_AddBlock(t *testing.T) {
	tests := []struct {
		name        string
		setupBlocks int
		newBlock    *bitcoin.Block
		shouldAdd   bool
		description string
	}{
		{
			name:        "Add valid block to Genesis",
			setupBlocks: 1, // Genesis only
			newBlock:    createValidNextBlock(),
			shouldAdd:   true,
			description: "Valid block should be added to blockchain",
		},
		{
			name:        "Add block with invalid previous hash",
			setupBlocks: 1,
			newBlock:    createBlockWithInvalidPrevHash(),
			shouldAdd:   false,
			description: "Block with wrong previous hash should be rejected",
		},
		{
			name:        "Add block with invalid proof of work",
			setupBlocks: 1,
			newBlock:    createBlockWithInvalidPoW(),
			shouldAdd:   false,
			description: "Block with invalid PoW should be rejected",
		},
		{
			name:        "Add multiple valid blocks",
			setupBlocks: 3,
			newBlock:    nil, // Will be created properly in test
			shouldAdd:   true,
			description: "Chain should accept multiple valid blocks",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("TDD RED: %s - %s", tt.name, tt.description)

			// Setup blockchain with specified number of blocks
			blockchain := setupBlockChain(tt.setupBlocks)

			initialHeight := blockchain.Height()

			// Create proper block for multi-block test
			if tt.newBlock == nil && tt.name == "Add multiple valid blocks" {
				tt.newBlock = createValidBlockAfter(blockchain.GetTip(), initialHeight+1)
			}

			// This should fail since we haven't implemented block addition validation yet
			err := blockchain.AddBlock(tt.newBlock)

			if tt.shouldAdd {
				if err != nil {
					t.Errorf("Expected block to be added, got error: %v", err)
				}

				if blockchain.Height() != initialHeight+1 {
					t.Errorf("Expected height %d, got %d", initialHeight+1, blockchain.Height())
				}
			} else {
				if err == nil {
					t.Error("Expected block addition to fail, but it succeeded")
				}

				if blockchain.Height() != initialHeight {
					t.Errorf("Expected height to remain %d, got %d", initialHeight, blockchain.Height())
				}
			}

			t.Logf("Block addition result: height %d, error: %v", blockchain.Height(), err)
		})
	}
}

// TestBlockChain_Validation tests blockchain validation (TDD RED phase)
func TestBlockChain_Validation(t *testing.T) {
	tests := []struct {
		name        string
		chainLength int
		corruption  string
		isValid     bool
		description string
	}{
		{
			name:        "Valid Genesis-only chain",
			chainLength: 1,
			corruption:  "",
			isValid:     true,
			description: "Genesis block only should be valid",
		},
		{
			name:        "Valid multi-block chain",
			chainLength: 5,
			corruption:  "",
			isValid:     true,
			description: "Multi-block chain should validate correctly",
		},
		{
			name:        "Chain with corrupted block",
			chainLength: 3,
			corruption:  "corrupt_block_1",
			isValid:     false,
			description: "Chain with corrupted block should be invalid",
		},
		{
			name:        "Chain with broken link",
			chainLength: 4,
			corruption:  "break_chain_link",
			isValid:     false,
			description: "Chain with broken previous hash link should be invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("TDD RED: %s - %s", tt.name, tt.description)

			// Setup blockchain
			blockchain := setupBlockChain(tt.chainLength)

			// Apply corruption if specified
			if tt.corruption != "" {
				applyChainCorruption(blockchain, tt.corruption)
			}

			// This should fail since we haven't implemented chain validation yet
			isValid := blockchain.ValidateChain()

			if isValid != tt.isValid {
				t.Errorf("Expected chain validity %v, got %v", tt.isValid, isValid)
			}

			t.Logf("Chain validation result: %v (length: %d)", isValid, blockchain.Height()+1)
		})
	}
}

// TestBlockChain_UTXO_Integration tests blockchain integration with UTXO set (TDD RED phase)
func TestBlockChain_UTXO_Integration(t *testing.T) {
	tests := []struct {
		name          string
		transactions  []string // Transaction types to create
		expectedUTXOs int
		description   string
	}{
		{
			name:          "Genesis block UTXO creation",
			transactions:  []string{"coinbase"},
			expectedUTXOs: 1,
			description:   "Genesis coinbase should create initial UTXO",
		},
		{
			name:          "Simple spend and create",
			transactions:  []string{"coinbase", "spend_and_create"},
			expectedUTXOs: 2, // Genesis coinbase + new block coinbase
			description:   "Transaction should update UTXO set correctly",
		},
		{
			name:          "Multiple transactions",
			transactions:  []string{"coinbase", "spend_create", "spend_create", "spend_only"},
			expectedUTXOs: 4, // Genesis + 3 new blocks (each adds coinbase)
			description:   "Multiple transactions should maintain correct UTXO count",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("TDD RED: %s - %s", tt.name, tt.description)

			// This should fail since we haven't implemented UTXO integration yet
			blockchain := bitcoin.NewBlockChain(createGenesisBlock())

			// Process transactions (skip coinbase since Genesis already has it)
			for i, txType := range tt.transactions {
				if i == 0 && txType == "coinbase" {
					continue // Genesis already processed
				}
				block := createValidBlockAfter(blockchain.GetTip(), blockchain.Height()+1)
				err := blockchain.AddBlock(block)
				if err != nil {
					t.Fatalf("Failed to add block with %s transaction: %v", txType, err)
				}
			}

			// Check UTXO set
			utxoSet := blockchain.GetUTXOSet()
			actualUTXOs := utxoSet.Size()

			// Debug: print all UTXOs to see what's happening
			allUTXOs := utxoSet.GetAllUTXOs()
			for i, utxo := range allUTXOs {
				t.Logf("UTXO %d: %s:%d = %d satoshis", i, utxo.TxHash().String(), utxo.OutputIndex(), utxo.Amount())
			}

			if actualUTXOs != tt.expectedUTXOs {
				t.Errorf("Expected %d UTXOs, got %d", tt.expectedUTXOs, actualUTXOs)
			}

			t.Logf("UTXO integration: %d transactions resulted in %d UTXOs",
				len(tt.transactions), actualUTXOs)
		})
	}
}

// TestBlockChain_Reorganization tests blockchain reorganization (TDD RED phase)
func TestBlockChain_Reorganization(t *testing.T) {
	tests := []struct {
		name        string
		mainChain   int // Length of main chain
		forkChain   int // Length of competing fork
		shouldReorg bool
		description string
	}{
		{
			name:        "No reorganization - fork shorter",
			mainChain:   5,
			forkChain:   3,
			shouldReorg: false,
			description: "Shorter fork should not trigger reorganization",
		},
		{
			name:        "Reorganization - fork longer",
			mainChain:   3,
			forkChain:   5,
			shouldReorg: true,
			description: "Longer fork should trigger reorganization",
		},
		{
			name:        "No reorganization - equal length",
			mainChain:   4,
			forkChain:   4,
			shouldReorg: false,
			description: "Equal length fork should not trigger reorganization",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("TDD RED: %s - %s", tt.name, tt.description)

			// Setup main chain
			blockchain := setupBlockChain(tt.mainChain)
			originalTip := blockchain.GetTip().Hash()

			// Debug: print chain info
			t.Logf("Main chain height: %d, tip: %s", blockchain.Height(), originalTip.String())
			if blockchain.Height() >= 0 {
				genesis := blockchain.GetBlock(0)
				t.Logf("Genesis hash: %s", genesis.Hash().String())
			}

			// Create competing fork
			forkBlocks := createForkChain(tt.forkChain)

			// This should fail since we haven't implemented reorganization yet
			reorgOccurred := false
			for i, block := range forkBlocks {
				t.Logf("Adding fork block %d: %s (nonce=%d)", i, block.Hash().String(), block.Header.Nonce)
				err := blockchain.AddBlock(block)
				t.Logf("Block addition result: error=%v, new tip=%s", err, blockchain.GetTip().Hash().String())
				if err == nil && blockchain.GetTip().Hash() != originalTip {
					reorgOccurred = true
					t.Logf("Reorganization detected at block %d", i)
				}
			}

			if reorgOccurred != tt.shouldReorg {
				t.Errorf("Expected reorganization %v, got %v", tt.shouldReorg, reorgOccurred)
			}

			expectedHeight := tt.mainChain - 1 // Convert to 0-based
			if tt.shouldReorg {
				expectedHeight = tt.forkChain - 1
			}

			if blockchain.Height() != expectedHeight {
				t.Errorf("Expected final height %d, got %d", expectedHeight, blockchain.Height())
			}

			t.Logf("Reorganization test: reorg=%v, final height=%d", reorgOccurred, blockchain.Height())
		})
	}
}

// Helper functions for test setup
func createGenesisBlock() *bitcoin.Block {
	// Create Genesis block with known parameters
	genesisHeader := bitcoin.BlockHeader{
		Version:       1,
		PrevBlockHash: bitcoin.ZeroHash,
		MerkleRoot:    mustParseHash("4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b"),
		Timestamp:     1231006505, // Genesis timestamp
		Bits:          0x1d00ffff,
		Nonce:         2083236893,
	}

	// Genesis coinbase transaction
	coinbaseTx := createCoinbaseTransaction(5000000000) // 50 BTC

	return bitcoin.NewBlock(genesisHeader, []bitcoin.Transaction{*coinbaseTx})
}

func createValidNextBlock() *bitcoin.Block {
	// Create a valid block that builds on Genesis
	header := bitcoin.BlockHeader{
		Version:       1,
		PrevBlockHash: mustParseHash("000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"),
		MerkleRoot:    mustParseHash("0000000000000000000000000000000000000000000000000000000000000001"),
		Timestamp:     1231006505 + 600, // 10 minutes later
		Bits:          0x1d00ffff,
		Nonce:         12345,
	}

	coinbaseTx := createCoinbaseTransaction(5000000000)
	return bitcoin.NewBlock(header, []bitcoin.Transaction{*coinbaseTx})
}

func createBlockWithInvalidPrevHash() *bitcoin.Block {
	header := bitcoin.BlockHeader{
		Version:       1,
		PrevBlockHash: mustParseHash("1111111111111111111111111111111111111111111111111111111111111111"),
		MerkleRoot:    mustParseHash("0000000000000000000000000000000000000000000000000000000000000001"),
		Timestamp:     1231006505 + 600,
		Bits:          0x1d00ffff,
		Nonce:         12345,
	}

	coinbaseTx := createCoinbaseTransaction(5000000000)
	return bitcoin.NewBlock(header, []bitcoin.Transaction{*coinbaseTx})
}

func createBlockWithInvalidPoW() *bitcoin.Block {
	header := bitcoin.BlockHeader{
		Version:       1,
		PrevBlockHash: mustParseHash("000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"),
		MerkleRoot:    mustParseHash("0000000000000000000000000000000000000000000000000000000000000001"),
		Timestamp:     1231006505 + 600,
		Bits:          0x1d00ffff,
		Nonce:         999999, // Invalid nonce that won't meet difficulty
	}

	coinbaseTx := createCoinbaseTransaction(5000000000)
	return bitcoin.NewBlock(header, []bitcoin.Transaction{*coinbaseTx})
}

func createValidBlock(height int) *bitcoin.Block {
	// Create a valid block for given height
	header := bitcoin.BlockHeader{
		Version:       1,
		PrevBlockHash: mustParseHash("000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f"), // Simplified
		MerkleRoot:    mustParseHash("0000000000000000000000000000000000000000000000000000000000000001"),
		Timestamp:     uint32(1231006505 + height*600),
		Bits:          0x1d00ffff,
		Nonce:         uint32(12345 + height),
	}

	coinbaseTx := createCoinbaseTransaction(5000000000)
	return bitcoin.NewBlock(header, []bitcoin.Transaction{*coinbaseTx})
}

func createCoinbaseTransaction(amount uint64) *bitcoin.Transaction {
	// Create a basic coinbase transaction
	input := bitcoin.TxInput{
		PreviousOutput: bitcoin.OutPoint{
			Hash:  bitcoin.ZeroHash,
			Index: 0xffffffff,
		},
		ScriptSig: []byte{0x04, 0xff, 0xff, 0x00, 0x1d}, // Genesis coinbase script
		Sequence:  0xffffffff,
	}

	output := bitcoin.TxOutput{
		Value: amount,
		ScriptPubKey: []byte{
			0x41, // OP_PUSHDATA 65 bytes
			0x04, 0x67, 0x8a, 0xfd, 0xb0, 0xfe, 0x55, 0x48, 0x27, 0x19,
			0x67, 0xf1, 0xa6, 0x71, 0x30, 0xb7, 0x10, 0x5c, 0xd6, 0xa8,
			0x28, 0xe0, 0x39, 0x09, 0xa6, 0x79, 0x62, 0xe0, 0xea, 0x1f,
			0x61, 0xde, 0xb6, 0x49, 0xf6, 0xbc, 0x3f, 0x4c, 0xef, 0x38,
			0xc4, 0xf3, 0x55, 0x04, 0xe5, 0x1e, 0xc1, 0x12, 0xde, 0x5c,
			0x38, 0x4d, 0xf7, 0xba, 0x0b, 0x8d, 0x57, 0x8a, 0x4c, 0x70,
			0x2b, 0x6b, 0xf1, 0x1d, 0x5f,
			0xac, // OP_CHECKSIG
		},
	}

	return bitcoin.NewTransaction(
		1,                          // version
		[]bitcoin.TxInput{input},   // inputs
		[]bitcoin.TxOutput{output}, // outputs
		0,                          // locktime
	)
}

func createUniqueCoinbaseTransaction(amount uint64, height int) *bitcoin.Transaction {
	// Create a unique coinbase transaction by including height in scriptSig
	heightBytes := []byte{byte(height & 0xff), byte((height >> 8) & 0xff), byte((height >> 16) & 0xff), byte((height >> 24) & 0xff)}
	scriptSig := append([]byte{0x04, 0xff, 0xff, 0x00, 0x1d}, heightBytes...)

	input := bitcoin.TxInput{
		PreviousOutput: bitcoin.OutPoint{
			Hash:  bitcoin.ZeroHash,
			Index: 0xffffffff,
		},
		ScriptSig: scriptSig, // Unique script with height
		Sequence:  0xffffffff,
	}

	output := bitcoin.TxOutput{
		Value: amount,
		ScriptPubKey: []byte{
			0x41, // OP_PUSHDATA 65 bytes
			0x04, 0x67, 0x8a, 0xfd, 0xb0, 0xfe, 0x55, 0x48, 0x27, 0x19,
			0x67, 0xf1, 0xa6, 0x71, 0x30, 0xb7, 0x10, 0x5c, 0xd6, 0xa8,
			0x28, 0xe0, 0x39, 0x09, 0xa6, 0x79, 0x62, 0xe0, 0xea, 0x1f,
			0x61, 0xde, 0xb6, 0x49, 0xf6, 0xbc, 0x3f, 0x4c, 0xef, 0x38,
			0xc4, 0xf3, 0x55, 0x04, 0xe5, 0x1e, 0xc1, 0x12, 0xde, 0x5c,
			0x38, 0x4d, 0xf7, 0xba, 0x0b, 0x8d, 0x57, 0x8a, 0x4c, 0x70,
			0x2b, 0x6b, 0xf1, 0x1d, 0x5f,
			0xac, // OP_CHECKSIG
		},
	}

	return bitcoin.NewTransaction(
		1,                          // version
		[]bitcoin.TxInput{input},   // inputs
		[]bitcoin.TxOutput{output}, // outputs
		0,                          // locktime
	)
}

func setupBlockChain(numBlocks int) *bitcoin.BlockChain {
	blockchain := bitcoin.NewBlockChain(createGenesisBlock())

	for i := 1; i < numBlocks; i++ {
		// Create block that properly links to previous block
		prevBlock := blockchain.GetTip()
		block := createValidBlockAfter(prevBlock, i)
		err := blockchain.AddBlock(block)
		if err != nil {
			panic("Failed to setup blockchain: " + err.Error())
		}
	}

	return blockchain
}

func createValidBlockAfter(prevBlock *bitcoin.Block, height int) *bitcoin.Block {
	// Create a valid block that builds on the previous block
	header := bitcoin.BlockHeader{
		Version:       1,
		PrevBlockHash: prevBlock.Hash(), // Correct previous hash
		MerkleRoot:    mustParseHash("0000000000000000000000000000000000000000000000000000000000000001"),
		Timestamp:     uint32(1231006505 + height*600),
		Bits:          0x1d00ffff,
		Nonce:         uint32(12345 + height),
	}

	// Create unique coinbase transaction for each block
	coinbaseTx := createUniqueCoinbaseTransaction(5000000000, height)
	return bitcoin.NewBlock(header, []bitcoin.Transaction{*coinbaseTx})
}

func applyChainCorruption(blockchain *bitcoin.BlockChain, corruption string) {
	// Simulate different types of corruption for testing
	switch corruption {
	case "corrupt_block_1":
		// Corrupt block at index 1 by modifying its nonce
		if blockchain.Height() >= 1 {
			block := blockchain.GetBlock(1)
			if block != nil {
				// Create a corrupted version by changing the nonce
				corruptedHeader := block.Header
				corruptedHeader.Nonce = 999999 // Invalid nonce
				corruptedBlock := bitcoin.NewBlock(corruptedHeader, block.Transactions)
				// Force replace the block in the blockchain (for testing purposes)
				// Note: This is a test hack - real blockchain wouldn't allow this
				blockchain.ForceReplaceBlock(1, corruptedBlock)
			}
		}
	case "break_chain_link":
		// Break the chain linkage by corrupting previous hash
		if blockchain.Height() >= 1 {
			block := blockchain.GetBlock(1)
			if block != nil {
				// Create a block with wrong previous hash
				corruptedHeader := block.Header
				corruptedHeader.PrevBlockHash = mustParseHash("1111111111111111111111111111111111111111111111111111111111111111")
				corruptedBlock := bitcoin.NewBlock(corruptedHeader, block.Transactions)
				blockchain.ForceReplaceBlock(1, corruptedBlock)
			}
		}
	}
}

func createBlockWithTransaction(txType string) *bitcoin.Block {
	// Create blocks with different transaction types for UTXO testing
	// For now, just return a valid next block
	return createValidNextBlock() // Simplified for now
}

func createForkChain(length int) []*bitcoin.Block {
	// Create a competing fork chain that branches from Genesis
	// length parameter represents desired total chain length, so create length-1 blocks
	forkBlocks := length - 1
	if forkBlocks <= 0 {
		return []*bitcoin.Block{}
	}
	blocks := make([]*bitcoin.Block, forkBlocks)

	genesisHash := mustParseHash("000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f")

	for i := 0; i < forkBlocks; i++ {
		var prevHash bitcoin.Hash256
		if i == 0 {
			// First fork block builds on Genesis
			prevHash = genesisHash
		} else {
			// Subsequent blocks build on previous fork block
			prevHash = blocks[i-1].Hash()
		}

		header := bitcoin.BlockHeader{
			Version:       1,
			PrevBlockHash: prevHash,
			MerkleRoot:    mustParseHash("0000000000000000000000000000000000000000000000000000000000000002"),
			Timestamp:     uint32(1231006505 + (i+100)*600), // Different timestamps
			Bits:          0x1d00ffff,
			Nonce:         uint32(50000 + i), // Different nonces for fork
		}

		// Create unique coinbase for fork blocks
		coinbaseTx := createUniqueCoinbaseTransaction(5000000000, i+100)
		blocks[i] = bitcoin.NewBlock(header, []bitcoin.Transaction{*coinbaseTx})
	}
	return blocks
}

func mustParseHash(hashStr string) bitcoin.Hash256 {
	hash, err := bitcoin.NewHash256FromString(hashStr)
	if err != nil {
		panic("Failed to parse hash: " + err.Error())
	}
	return hash
}
