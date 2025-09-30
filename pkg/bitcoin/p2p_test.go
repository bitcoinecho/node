package bitcoin

import (
	"net"
	"testing"
	"time"
)

// TestP2P_MessageSerialization tests Bitcoin P2P message serialization (TDD RED phase)
func TestP2P_MessageSerialization(t *testing.T) {
	tests := []struct {
		name        string
		command     string
		payload     []byte
		expectedLen int
		description string
	}{
		{
			name:        "Version message",
			command:     "version",
			payload:     []byte{0x01, 0x02, 0x03, 0x04}, // Dummy payload
			expectedLen: 24 + 4,                         // Header + payload
			description: "Version message should serialize correctly",
		},
		{
			name:        "Ping message",
			command:     "ping",
			payload:     []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}, // 8-byte nonce
			expectedLen: 24 + 8,                                                 // Header + nonce
			description: "Ping message should serialize with 8-byte nonce",
		},
		{
			name:        "Empty message (verack)",
			command:     "verack",
			payload:     []byte{},
			expectedLen: 24, // Header only
			description: "Empty message should serialize header only",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("TDD RED: %s - %s", tt.name, tt.description)

			// This should fail since we haven't implemented P2P message serialization yet
			msg := NewP2PMessage(tt.command, tt.payload)
			serialized := msg.Serialize()

			if len(serialized) != tt.expectedLen {
				t.Errorf("Expected serialized length %d, got %d", tt.expectedLen, len(serialized))
			}

			t.Logf("Message: %s, Payload length: %d, Serialized length: %d",
				tt.command, len(tt.payload), len(serialized))
		})
	}
}

// TestP2P_MessageDeserialization tests Bitcoin P2P message deserialization (TDD RED phase)
func TestP2P_MessageDeserialization(t *testing.T) {
	tests := []struct {
		name        string
		messageData []byte
		expectedCmd string
		expectedLen int
		description string
	}{
		{
			name: "Version message deserialization",
			// Bitcoin mainnet magic + "version" command + length + checksum + payload
			messageData: []byte{
				0xf9, 0xbe, 0xb4, 0xd9, // Magic bytes (mainnet)
				0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x00, 0x00, 0x00, 0x00, 0x00, // "version" command
				0x04, 0x00, 0x00, 0x00, // Payload length (4 bytes)
				0x8d, 0xe4, 0x72, 0xe2, // Correct checksum for [1,2,3,4]
				0x01, 0x02, 0x03, 0x04, // Payload
			},
			expectedCmd: "version",
			expectedLen: 4,
			description: "Version message should deserialize correctly",
		},
		{
			name: "Ping message deserialization",
			messageData: []byte{
				0xf9, 0xbe, 0xb4, 0xd9, // Magic bytes
				0x70, 0x69, 0x6e, 0x67, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // "ping" command
				0x08, 0x00, 0x00, 0x00, // Payload length (8 bytes)
				0x25, 0x02, 0xfa, 0x94, // Correct checksum for 8-byte nonce
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, // 8-byte nonce
			},
			expectedCmd: "ping",
			expectedLen: 8,
			description: "Ping message should deserialize with correct nonce",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("TDD RED: %s - %s", tt.name, tt.description)

			// This should fail since we haven't implemented P2P message deserialization yet
			msg, err := DeserializeP2PMessage(tt.messageData)
			if err != nil {
				t.Fatalf("Failed to deserialize message: %v", err)
			}

			if msg.Command() != tt.expectedCmd {
				t.Errorf("Expected command %s, got %s", tt.expectedCmd, msg.Command())
			}

			if len(msg.Payload()) != tt.expectedLen {
				t.Errorf("Expected payload length %d, got %d", tt.expectedLen, len(msg.Payload()))
			}

			t.Logf("Deserialized command: %s, payload length: %d", msg.Command(), len(msg.Payload()))
		})
	}
}

// TestP2P_MessageValidation tests P2P message validation (TDD RED phase)
func TestP2P_MessageValidation(t *testing.T) {
	tests := []struct {
		name        string
		messageData []byte
		shouldFail  bool
		description string
	}{
		{
			name: "Valid message with correct checksum",
			messageData: []byte{
				0xf9, 0xbe, 0xb4, 0xd9, // Magic bytes
				0x70, 0x69, 0x6e, 0x67, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // "ping"
				0x08, 0x00, 0x00, 0x00, // Length
				0x25, 0x02, 0xfa, 0x94, // Correct checksum for payload below
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, // Payload
			},
			shouldFail:  false,
			description: "Valid message should pass validation",
		},
		{
			name: "Invalid magic bytes",
			messageData: []byte{
				0x00, 0x00, 0x00, 0x00, // Wrong magic bytes
				0x70, 0x69, 0x6e, 0x67, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x08, 0x00, 0x00, 0x00,
				0x9c, 0x12, 0xcf, 0xdc,
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
			},
			shouldFail:  true,
			description: "Invalid magic bytes should fail validation",
		},
		{
			name: "Invalid checksum",
			messageData: []byte{
				0xf9, 0xbe, 0xb4, 0xd9,
				0x70, 0x69, 0x6e, 0x67, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x08, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, // Wrong checksum
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
			},
			shouldFail:  true,
			description: "Invalid checksum should fail validation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("TDD RED: %s - %s", tt.name, tt.description)

			// This should fail since we haven't implemented message validation yet
			isValid := ValidateP2PMessage(tt.messageData)

			if tt.shouldFail && isValid {
				t.Errorf("Expected message to fail validation, but it passed")
			}

			if !tt.shouldFail && !isValid {
				t.Errorf("Expected message to pass validation, but it failed")
			}

			t.Logf("Message validation result: %v (expected fail: %v)", isValid, tt.shouldFail)
		})
	}
}

// TestP2P_PeerHandshake tests Bitcoin peer handshake process (TDD RED phase)
func TestP2P_PeerHandshake(t *testing.T) {
	tests := []struct {
		name        string
		peerAddr    string
		ourVersion  uint32
		expected    bool
		description string
	}{
		{
			name:        "Successful handshake",
			peerAddr:    "peer.example.com:8333", // Mock peer that will work
			ourVersion:  70015,                   // Protocol version
			expected:    true,
			description: "Valid peer should complete handshake",
		},
		{
			name:        "Version too old",
			peerAddr:    "127.0.0.1:8333",
			ourVersion:  60000, // Old version
			expected:    false,
			description: "Old protocol version should fail handshake",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("TDD RED: %s - %s", tt.name, tt.description)

			// This should fail since we haven't implemented peer handshake yet
			peer := NewPeer(tt.peerAddr)
			success := peer.PerformHandshake(tt.ourVersion)

			if success != tt.expected {
				t.Errorf("Expected handshake result %v, got %v", tt.expected, success)
			}

			t.Logf("Handshake with %s (version %d): %v", tt.peerAddr, tt.ourVersion, success)
		})
	}
}

// TestP2P_PeerConnection tests peer connection management (TDD RED phase)
func TestP2P_PeerConnection(t *testing.T) {
	tests := []struct {
		name          string
		peerAddr      string
		shouldConnect bool
		description   string
	}{
		{
			name:          "Connect to localhost",
			peerAddr:      "127.0.0.1:8333",
			shouldConnect: false, // Will fail since no actual node running
			description:   "Connection attempt should be handled gracefully",
		},
		{
			name:          "Invalid address",
			peerAddr:      "invalid:address",
			shouldConnect: false,
			description:   "Invalid address should fail connection",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("TDD RED: %s - %s", tt.name, tt.description)

			// This should fail since we haven't implemented peer connection yet
			peer := NewPeer(tt.peerAddr)
			err := peer.Connect()

			if tt.shouldConnect && err != nil {
				t.Errorf("Expected successful connection, got error: %v", err)
			}

			if !tt.shouldConnect && err == nil {
				t.Errorf("Expected connection to fail, but it succeeded")
			}

			t.Logf("Connection to %s: %v", tt.peerAddr, err)
		})
	}
}

// TestPeer_Close tests peer connection closing
func TestPeer_Close(t *testing.T) {
	peer := NewPeer("127.0.0.1:8333")

	// Should not error on closing unconnected peer
	err := peer.Close()
	if err != nil {
		t.Errorf("Expected no error closing unconnected peer, got: %v", err)
	}

	// Test closing after attempting connection
	peer.Connect() // Try to connect first
	err = peer.Close()
	if err != nil {
		t.Errorf("Expected no error closing peer, got: %v", err)
	}
}

// TestDeserializeP2PMessage_EdgeCases tests P2P message deserialization edge cases
func TestDeserializeP2PMessage_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		shouldError bool
		errorMsg    string
	}{
		{
			name:        "Message too short",
			data:        []byte{0x01, 0x02},
			shouldError: true,
			errorMsg:    "message too short",
		},
		{
			name:        "Invalid magic bytes",
			data:        append([]byte{0x00, 0x00, 0x00, 0x00}, make([]byte, 20)...),
			shouldError: true,
			errorMsg:    "invalid magic bytes",
		},
		{
			name:        "Invalid checksum",
			data:        append([]byte{0xF9, 0xBE, 0xB4, 0xD9}, append(make([]byte, 16), []byte{0x00, 0x00, 0x00, 0x01, 0xFF, 0xFF, 0xFF, 0xFF, 0x00}...)...),
			shouldError: true,
			errorMsg:    "invalid checksum",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DeserializeP2PMessage(tt.data)
			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error '%s', got none", tt.errorMsg)
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Expected error '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// TestPeer_Connect_EdgeCases tests more edge cases for peer connection
func TestPeer_Connect_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		address     string
		shouldError bool
		description string
	}{
		{
			name:        "Connect to nonexistent host",
			address:     "nonexistent.example.com:8333",
			shouldError: true,
			description: "Connection to nonexistent host should fail",
		},
		{
			name:        "Connect with invalid port",
			address:     "127.0.0.1:99999",
			shouldError: true,
			description: "Connection with invalid port should fail",
		},
		{
			name:        "Connect with malformed address",
			address:     ":::invalid:::8333",
			shouldError: true,
			description: "Connection with malformed address should fail",
		},
		{
			name:        "Connect to localhost unreachable port",
			address:     "127.0.0.1:1",
			shouldError: true,
			description: "Connection to unreachable port should fail",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			peer := NewPeer(tt.address)
			err := peer.Connect()

			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected connection to fail for %s, but it succeeded", tt.address)
				}
			} else {
				if err != nil {
					t.Errorf("Expected connection to succeed for %s, got error: %v", tt.address, err)
				}
			}

			t.Logf("Connection attempt to %s: %v", tt.address, err)
		})
	}
}

// TestPeer_Connect_Success tests successful connection paths
func TestPeer_Connect_Success(t *testing.T) {
	// This test uses a mock successful connection by creating our own test listener
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create test listener: %v", err)
	}
	defer listener.Close()

	// Get the actual address the listener is bound to
	addr := listener.Addr().String()

	// Accept one connection in a goroutine
	go func() {
		conn, err := listener.Accept()
		if err == nil {
			conn.Close() // Close immediately after accepting
		}
	}()

	// Test connecting to our listener
	peer := NewPeer(addr)
	err = peer.Connect()
	if err != nil {
		t.Errorf("Expected successful connection to test listener, got error: %v", err)
	}

	// Test that connection was established (peer.conn should be set)
	if peer.conn == nil {
		t.Error("Expected peer.conn to be set after successful connection")
	}

	// Clean up
	peer.Close()
}

// TestPeer_Close_States tests peer closing in different states
func TestPeer_Close_States(t *testing.T) {
	// Test closing peer that was never connected
	peer1 := NewPeer("127.0.0.1:8333")
	err := peer1.Close()
	if err != nil {
		t.Errorf("Expected no error closing unconnected peer, got: %v", err)
	}

	// Test closing peer after failed connection attempt
	peer2 := NewPeer("invalid:address")
	peer2.Connect() // This will fail
	err = peer2.Close()
	if err != nil {
		t.Errorf("Expected no error closing peer after failed connection, got: %v", err)
	}

	// Test closing peer with active connection
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to create test listener: %v", err)
	}
	defer listener.Close()

	addr := listener.Addr().String()

	// Accept and hold connection
	go func() {
		conn, err := listener.Accept()
		if err == nil {
			time.Sleep(100 * time.Millisecond) // Keep connection alive briefly
			conn.Close()
		}
	}()

	peer3 := NewPeer(addr)
	err = peer3.Connect()
	if err == nil {
		// Now close the connected peer
		err = peer3.Close()
		if err != nil {
			t.Errorf("Expected no error closing connected peer, got: %v", err)
		}
	}

	// Test double close
	peer4 := NewPeer("127.0.0.1:8333")
	peer4.Connect()     // Try to connect
	peer4.Close()       // First close
	err = peer4.Close() // Second close
	if err != nil {
		t.Errorf("Expected no error on double close, got: %v", err)
	}
}

// TestPeer_PerformHandshake_EdgeCases tests handshake edge cases
func TestPeer_PerformHandshake_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		version     uint32
		peerAddr    string
		expected    bool
		description string
	}{
		{
			name:        "Very old version",
			version:     1,
			peerAddr:    "127.0.0.1:8333",
			expected:    false,
			description: "Very old protocol version should fail",
		},
		{
			name:        "Future version",
			version:     999999,
			peerAddr:    "peer.example.com:8333",
			expected:    true,
			description: "Future version should work",
		},
		{
			name:        "Boundary version",
			version:     70002,
			peerAddr:    "peer.example.com:8333",
			expected:    true,
			description: "Boundary version should work",
		},
		{
			name:        "Zero version",
			version:     0,
			peerAddr:    "127.0.0.1:8333",
			expected:    false,
			description: "Zero version should fail",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			peer := NewPeer(tt.peerAddr)
			success := peer.PerformHandshake(tt.version)

			if success != tt.expected {
				t.Errorf("Expected handshake result %v for version %d, got %v", tt.expected, tt.version, success)
			}

			t.Logf("Handshake with version %d: %v (expected: %v)", tt.version, success, tt.expected)
		})
	}
}
