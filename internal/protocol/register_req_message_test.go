package protocol

import (
	"testing"

	"github.com/google/uuid"
)

func TestRegisterReqMessage_EncodeDecode(t *testing.T) {
	// Original RegisterReqMessage
	originalMsg := &RegisterReqMessage{
		RequestID: uuid.New(),
		AgentName: "TestAgent",
	}

	// Encode the original message
	encoded, err := originalMsg.Encode()
	if err != nil {
		t.Fatalf("Failed to encode message: %v", err)
	}

	// Decode the encoded message
	decodedMsg := &RegisterReqMessage{}
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
}
