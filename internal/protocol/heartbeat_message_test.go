package protocol

import (
	"testing"

	"github.com/google/uuid"
)

func TestHeartbeatMessage_EncodeDecode(t *testing.T) {
	// Original HeartbeatMessage
	originalMsg := &HeartbeatMessage{
		RequestID:  uuid.New(),
		ClientID:   12345,
		AgentName:  "test",
		LanAddress: 3030303030,
	}

	// Encode the original message
	encoded, err := originalMsg.Encode()
	if err != nil {
		t.Fatalf("Failed to encode message: %v", err)
	}

	// Decode the encoded message
	decodedMsg := &HeartbeatMessage{}
	err = decodedMsg.Decode(encoded)
	if err != nil {
		t.Fatalf("Failed to decode message: %v", err)
	}

	// Check if the decoded message matches the original message
	if originalMsg.RequestID != decodedMsg.RequestID {
		t.Errorf("RequestID mismatch: got %v, want %v", decodedMsg.RequestID, originalMsg.RequestID)
	}

	if originalMsg.AgentName != decodedMsg.AgentName {
		t.Errorf("AgentName mismatch: got %v, want %v", decodedMsg.AgentName, originalMsg.AgentName)
	}

	if originalMsg.ClientID != decodedMsg.ClientID {
		t.Errorf("ClientID mismatch: got %v, want %v", decodedMsg.ClientID, originalMsg.ClientID)
	}
	if originalMsg.LanAddress != decodedMsg.LanAddress {
		t.Errorf("LanAddress mismatch: got %v, want %v", decodedMsg.ClientID, originalMsg.ClientID)
	}
}
