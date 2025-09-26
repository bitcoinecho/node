package bitcoin

import (
	"fmt"
)

// Script represents a Bitcoin script
type Script []byte

// ScriptOpcode represents a script operation code
type ScriptOpcode byte

// Script operation codes
const (
	// Constants
	OP_0                   ScriptOpcode = 0x00
	OP_FALSE               ScriptOpcode = OP_0
	OP_PUSHDATA1           ScriptOpcode = 0x4c
	OP_PUSHDATA2           ScriptOpcode = 0x4d
	OP_PUSHDATA4           ScriptOpcode = 0x4e
	OP_1NEGATE             ScriptOpcode = 0x4f
	OP_RESERVED            ScriptOpcode = 0x50
	OP_1                   ScriptOpcode = 0x51
	OP_TRUE                ScriptOpcode = OP_1

	// Flow control
	OP_NOP                 ScriptOpcode = 0x61
	OP_VER                 ScriptOpcode = 0x62
	OP_IF                  ScriptOpcode = 0x63
	OP_NOTIF               ScriptOpcode = 0x64
	OP_VERIF               ScriptOpcode = 0x65
	OP_VERNOTIF            ScriptOpcode = 0x66
	OP_ELSE                ScriptOpcode = 0x67
	OP_ENDIF               ScriptOpcode = 0x68
	OP_VERIFY              ScriptOpcode = 0x69
	OP_RETURN              ScriptOpcode = 0x6a

	// Stack ops
	OP_TOALTSTACK          ScriptOpcode = 0x6b
	OP_FROMALTSTACK        ScriptOpcode = 0x6c
	OP_2DROP               ScriptOpcode = 0x6d
	OP_2DUP                ScriptOpcode = 0x6e
	OP_3DUP                ScriptOpcode = 0x6f
	OP_2OVER               ScriptOpcode = 0x70
	OP_2ROT                ScriptOpcode = 0x71
	OP_2SWAP               ScriptOpcode = 0x72
	OP_IFDUP               ScriptOpcode = 0x73
	OP_DEPTH               ScriptOpcode = 0x74
	OP_DROP                ScriptOpcode = 0x75
	OP_DUP                 ScriptOpcode = 0x76
	OP_NIP                 ScriptOpcode = 0x77
	OP_OVER                ScriptOpcode = 0x78
	OP_PICK                ScriptOpcode = 0x79
	OP_ROLL                ScriptOpcode = 0x7a
	OP_ROT                 ScriptOpcode = 0x7b
	OP_SWAP                ScriptOpcode = 0x7c
	OP_TUCK                ScriptOpcode = 0x7d

	// String ops
	OP_SIZE                ScriptOpcode = 0x82

	// Bitwise logic
	OP_EQUAL               ScriptOpcode = 0x87
	OP_EQUALVERIFY         ScriptOpcode = 0x88

	// Arithmetic
	OP_1ADD                ScriptOpcode = 0x8b
	OP_1SUB                ScriptOpcode = 0x8c
	OP_NEGATE              ScriptOpcode = 0x8f
	OP_ABS                 ScriptOpcode = 0x90
	OP_NOT                 ScriptOpcode = 0x91
	OP_0NOTEQUAL           ScriptOpcode = 0x92
	OP_ADD                 ScriptOpcode = 0x93
	OP_SUB                 ScriptOpcode = 0x94
	OP_BOOLAND             ScriptOpcode = 0x9a
	OP_BOOLOR              ScriptOpcode = 0x9b
	OP_NUMEQUAL            ScriptOpcode = 0x9c
	OP_NUMEQUALVERIFY      ScriptOpcode = 0x9d
	OP_NUMNOTEQUAL         ScriptOpcode = 0x9e
	OP_LESSTHAN            ScriptOpcode = 0x9f
	OP_GREATERTHAN         ScriptOpcode = 0xa0
	OP_LESSTHANOREQUAL     ScriptOpcode = 0xa1
	OP_GREATERTHANOREQUAL  ScriptOpcode = 0xa2
	OP_MIN                 ScriptOpcode = 0xa3
	OP_MAX                 ScriptOpcode = 0xa4
	OP_WITHIN              ScriptOpcode = 0xa5

	// Crypto
	OP_RIPEMD160           ScriptOpcode = 0xa6
	OP_SHA1                ScriptOpcode = 0xa7
	OP_SHA256              ScriptOpcode = 0xa8
	OP_HASH160             ScriptOpcode = 0xa9
	OP_HASH256             ScriptOpcode = 0xaa
	OP_CODESEPARATOR       ScriptOpcode = 0xab
	OP_CHECKSIG            ScriptOpcode = 0xac
	OP_CHECKSIGVERIFY      ScriptOpcode = 0xad
	OP_CHECKMULTISIG       ScriptOpcode = 0xae
	OP_CHECKMULTISIGVERIFY ScriptOpcode = 0xaf

	// Expansion
	OP_NOP1                ScriptOpcode = 0xb0
	OP_CHECKLOCKTIMEVERIFY ScriptOpcode = 0xb1 // BIP65
	OP_CHECKSEQUENCEVERIFY ScriptOpcode = 0xb2 // BIP112
	OP_NOP4                ScriptOpcode = 0xb3
	OP_NOP5                ScriptOpcode = 0xb4
	OP_NOP6                ScriptOpcode = 0xb5
	OP_NOP7                ScriptOpcode = 0xb6
	OP_NOP8                ScriptOpcode = 0xb7
	OP_NOP9                ScriptOpcode = 0xb8
	OP_NOP10               ScriptOpcode = 0xb9

	// Invalid opcodes
	OP_INVALIDOPCODE       ScriptOpcode = 0xff
)

// ScriptType represents the type of a script
type ScriptType int

const (
	ScriptTypeUnknown ScriptType = iota
	ScriptTypeP2PK    // Pay-to-Public-Key
	ScriptTypeP2PKH   // Pay-to-Public-Key-Hash
	ScriptTypeP2SH    // Pay-to-Script-Hash
	ScriptTypeP2WPKH  // Pay-to-Witness-Public-Key-Hash
	ScriptTypeP2WSH   // Pay-to-Witness-Script-Hash
	ScriptTypeP2TR    // Pay-to-Taproot
	ScriptTypeMultisig
	ScriptTypeNullData // OP_RETURN
)

// ScriptEngine executes Bitcoin scripts
type ScriptEngine struct {
	stack    [][]byte
	altStack [][]byte
	script   Script
	pc       int

	// Execution flags
	flags ScriptFlags

	// Transaction context for signature verification
	tx       *Transaction
	txIdx    int
	prevOuts []TxOutput
}

// ScriptFlags control script execution behavior
type ScriptFlags uint32

const (
	ScriptFlagsNone                    ScriptFlags = 0
	ScriptVerifyP2SH                  ScriptFlags = 1 << 0  // BIP16
	ScriptVerifyStrictEnc             ScriptFlags = 1 << 1  // Strict DER encoding
	ScriptVerifyDERSig                ScriptFlags = 1 << 2  // Strict DER signatures
	ScriptVerifyLowS                  ScriptFlags = 1 << 3  // Low S values
	ScriptVerifyNullDummy             ScriptFlags = 1 << 4  // Null dummy for multisig
	ScriptVerifySigPushOnly           ScriptFlags = 1 << 5  // Only push operations in scriptSig
	ScriptVerifyMinimalData           ScriptFlags = 1 << 6  // Minimal pushdata operations
	ScriptVerifyDiscourageUpgradableNops ScriptFlags = 1 << 7
	ScriptVerifyCleanStack            ScriptFlags = 1 << 8  // Clean stack after execution
	ScriptVerifyCheckLockTimeVerify   ScriptFlags = 1 << 9  // BIP65
	ScriptVerifyCheckSequenceVerify   ScriptFlags = 1 << 10 // BIP112
	ScriptVerifyWitness              ScriptFlags = 1 << 11 // BIP141
	ScriptVerifyDiscourageUpgradableWitnessProgram ScriptFlags = 1 << 12
	ScriptVerifyMinimalIf            ScriptFlags = 1 << 13
	ScriptVerifyNullFail             ScriptFlags = 1 << 14
	ScriptVerifyWitnessPubkeyType    ScriptFlags = 1 << 15
	ScriptVerifyConstScriptCode      ScriptFlags = 1 << 16 // BIP342
	ScriptVerifyTaproot              ScriptFlags = 1 << 17 // BIP340/341/342
)

// NewScriptEngine creates a new script execution engine
func NewScriptEngine(script Script, tx *Transaction, txIdx int, prevOuts []TxOutput, flags ScriptFlags) *ScriptEngine {
	return &ScriptEngine{
		stack:    make([][]byte, 0, 100),
		altStack: make([][]byte, 0, 100),
		script:   script,
		pc:       0,
		flags:    flags,
		tx:       tx,
		txIdx:    txIdx,
		prevOuts: prevOuts,
	}
}

// Execute runs the script and returns true if successful
func (se *ScriptEngine) Execute() (bool, error) {
	// TODO: Implement full script execution
	// This is a placeholder that will need comprehensive implementation

	for se.pc < len(se.script) {
		opcode := ScriptOpcode(se.script[se.pc])
		se.pc++

		if err := se.executeOpcode(opcode); err != nil {
			return false, err
		}
	}

	// Script succeeds if stack is not empty and top element is true
	if len(se.stack) == 0 {
		return false, nil
	}

	return se.isTrue(se.stack[len(se.stack)-1]), nil
}

// executeOpcode executes a single opcode
func (se *ScriptEngine) executeOpcode(opcode ScriptOpcode) error {
	switch opcode {
	case OP_DUP:
		if len(se.stack) < 1 {
			return fmt.Errorf("OP_DUP: insufficient stack items")
		}
		top := se.stack[len(se.stack)-1]
		se.stack = append(se.stack, append([]byte{}, top...))

	case OP_HASH160:
		if len(se.stack) < 1 {
			return fmt.Errorf("OP_HASH160: insufficient stack items")
		}
		data := se.stack[len(se.stack)-1]
		se.stack = se.stack[:len(se.stack)-1]
		hash := hash160(data)
		se.stack = append(se.stack, hash[:])

	case OP_EQUAL:
		if len(se.stack) < 2 {
			return fmt.Errorf("OP_EQUAL: insufficient stack items")
		}
		a := se.stack[len(se.stack)-2]
		b := se.stack[len(se.stack)-1]
		se.stack = se.stack[:len(se.stack)-2]

		if bytesEqual(a, b) {
			se.stack = append(se.stack, []byte{1})
		} else {
			se.stack = append(se.stack, []byte{0})
		}

	case OP_EQUALVERIFY:
		if err := se.executeOpcode(OP_EQUAL); err != nil {
			return err
		}
		return se.executeOpcode(OP_VERIFY)

	case OP_VERIFY:
		if len(se.stack) < 1 {
			return fmt.Errorf("OP_VERIFY: insufficient stack items")
		}
		top := se.stack[len(se.stack)-1]
		se.stack = se.stack[:len(se.stack)-1]

		if !se.isTrue(top) {
			return fmt.Errorf("OP_VERIFY: failed")
		}

	case OP_CHECKSIG:
		// TODO: Implement signature verification
		// This is complex and requires proper implementation
		// For now, always return true (placeholder)
		se.stack = append(se.stack, []byte{1})

	default:
		// Handle push operations
		if opcode >= 1 && opcode <= 75 {
			// Direct push of N bytes
			n := int(opcode)
			if se.pc+n > len(se.script) {
				return fmt.Errorf("push operation exceeds script bounds")
			}
			data := se.script[se.pc : se.pc+n]
			se.pc += n
			se.stack = append(se.stack, data)
		} else {
			return fmt.Errorf("unimplemented opcode: %02x", opcode)
		}
	}

	return nil
}

// isTrue returns true if the byte slice represents a true value
func (se *ScriptEngine) isTrue(data []byte) bool {
	if len(data) == 0 {
		return false
	}

	// All bytes must be zero except possibly the last byte which can be 0x80
	for i := 0; i < len(data)-1; i++ {
		if data[i] != 0 {
			return true
		}
	}

	// Last byte can be 0x00 or 0x80 (negative zero) for false
	last := data[len(data)-1]
	return last != 0 && last != 0x80
}

// Script size constants
const (
	P2PKHScriptSize         = 25 // OP_DUP OP_HASH160 <20-byte hash> OP_EQUALVERIFY OP_CHECKSIG
	P2SHScriptSize          = 23 // OP_HASH160 <20-byte hash> OP_EQUAL
	P2WPKHScriptSize        = 22 // OP_0 <20-byte hash>
	P2WSHScriptSize         = 34 // OP_0 <32-byte hash>
	P2TRScriptSize          = 34 // OP_1 <32-byte key>
	CompressedPubKeySize    = 33 // 0x02/0x03 + 32 bytes
	UncompressedPubKeySize  = 65 // 0x04 + 64 bytes
	Hash160Size             = 20 // RIPEMD160 output
	Hash256Size             = 32 // SHA256 output
)

// AnalyzeScript determines the type of a script
func (s Script) AnalyzeScript() ScriptType {
	if len(s) == 0 {
		return ScriptTypeUnknown
	}

	// P2PKH: OP_DUP OP_HASH160 <20-byte hash> OP_EQUALVERIFY OP_CHECKSIG
	if len(s) == P2PKHScriptSize &&
		s[0] == byte(OP_DUP) &&
		s[1] == byte(OP_HASH160) &&
		s[2] == Hash160Size &&
		s[23] == byte(OP_EQUALVERIFY) &&
		s[24] == byte(OP_CHECKSIG) {
		return ScriptTypeP2PKH
	}

	// P2SH: OP_HASH160 <20-byte hash> OP_EQUAL
	if len(s) == P2SHScriptSize &&
		s[0] == byte(OP_HASH160) &&
		s[1] == Hash160Size &&
		s[22] == byte(OP_EQUAL) {
		return ScriptTypeP2SH
	}

	// P2PK: <pubkey> OP_CHECKSIG
	if len(s) >= 35 && s[len(s)-1] == byte(OP_CHECKSIG) {
		// Check for compressed pubkey (push 33 + 33-byte key + OP_CHECKSIG)
		if len(s) >= 35 && s[0] == CompressedPubKeySize && (s[1] == 0x02 || s[1] == 0x03) {
			return ScriptTypeP2PK
		}
		// Check for uncompressed pubkey (push 65 + 65-byte key + OP_CHECKSIG)
		if len(s) >= 67 && s[0] == UncompressedPubKeySize && s[1] == 0x04 {
			return ScriptTypeP2PK
		}
	}

	// P2WPKH: OP_0 <20-byte hash>
	if len(s) == P2WPKHScriptSize && s[0] == byte(OP_0) && s[1] == Hash160Size {
		return ScriptTypeP2WPKH
	}

	// P2WSH: OP_0 <32-byte hash>
	if len(s) == P2WSHScriptSize && s[0] == byte(OP_0) && s[1] == Hash256Size {
		return ScriptTypeP2WSH
	}

	// P2TR: OP_1 <32-byte key>
	if len(s) == P2TRScriptSize && s[0] == byte(OP_1) && s[1] == Hash256Size {
		return ScriptTypeP2TR
	}

	// Multisig: OP_M <pubkey1> ... <pubkeyN> OP_N OP_CHECKMULTISIG
	if len(s) >= 4 && s[len(s)-1] == byte(OP_CHECKMULTISIG) {
		// Check if starts with OP_1 through OP_16 (0x51-0x60)
		if s[0] >= 0x51 && s[0] <= 0x60 {
			// Check if second-to-last byte is also OP_1 through OP_16
			if s[len(s)-2] >= 0x51 && s[len(s)-2] <= 0x60 {
				return ScriptTypeMultisig
			}
		}
	}

	// OP_RETURN (null data)
	if len(s) > 0 && s[0] == byte(OP_RETURN) {
		return ScriptTypeNullData
	}

	return ScriptTypeUnknown
}

// IsStandard returns true if the script is considered standard
func (s Script) IsStandard() bool {
	scriptType := s.AnalyzeScript()
	switch scriptType {
	case ScriptTypeP2PKH, ScriptTypeP2SH, ScriptTypeP2WPKH, ScriptTypeP2WSH, ScriptTypeP2TR, ScriptTypeP2PK:
		return true
	case ScriptTypeNullData:
		// OP_RETURN scripts are standard if they're not too large
		return len(s) <= 80
	case ScriptTypeMultisig:
		// Validate multisig constraints (M-of-N limits)
		return s.isStandardMultisig()
	default:
		return false
	}
}

// isStandardMultisig checks if a multisig script meets standardness rules
func (s Script) isStandardMultisig() bool {
	if len(s) < 4 || s[len(s)-1] != byte(OP_CHECKMULTISIG) {
		return false
	}

	// Check if starts with OP_1 through OP_3 (standard M values)
	if s[0] < 0x51 || s[0] > 0x53 {
		return false // Only 1-of-N, 2-of-N, 3-of-N are standard
	}

	// Check if second-to-last byte is OP_1 through OP_3 (standard N values)
	if s[len(s)-2] < 0x51 || s[len(s)-2] > 0x53 {
		return false // Only M-of-1, M-of-2, M-of-3 are standard
	}

	// M should be <= N
	m := s[0] - 0x50  // OP_1 = 0x51, so M = s[0] - 0x50
	n := s[len(s)-2] - 0x50 // N = s[len(s)-2] - 0x50

	return m <= n && n <= 3 // Standard multisig is limited to 3 keys max
}

// Helper functions
func hash160(data []byte) Hash160 {
	// TODO: Implement proper HASH160 (RIPEMD160(SHA256(data)))
	// For now, return zero hash
	return ZeroHash160
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}