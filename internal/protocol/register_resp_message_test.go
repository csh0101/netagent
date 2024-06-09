package protocol

import (
	"bytes"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRegisterRespMessage_EncodeDecode(t *testing.T) {
	// Original RegisterRespMessage
	originalMsg := &RegisterRespMessage{
		RequestID:  uuid.New(),
		ClientID:   12,
		Success:    true,
		LanAddress: 2885812509,
		Msg:        "Registration successful",
	}

	// Encode the original message
	encoded, err := originalMsg.Encode()
	if err != nil {
		t.Fatalf("Failed to encode message: %v", err)
	}

	// Decode the encoded message
	_, buf, err := UnPackMessage(bytes.NewBuffer(encoded))
	assert.Nil(t, err)
	decodedMsg := &RegisterRespMessage{}
	err = decodedMsg.Decode(buf)
	if err != nil {
		t.Fatalf("Failed to decode message: %v", err)
	}

	// Check if the decoded message matches the original message
	if originalMsg.RequestID != decodedMsg.RequestID {
		t.Errorf("RequestID mismatch: got %v, want %v", decodedMsg.RequestID, originalMsg.RequestID)
	}
	if originalMsg.ClientID != decodedMsg.ClientID {
		t.Errorf("ClientID mismatch: got %v, want %v", decodedMsg.ClientID, originalMsg.ClientID)
	}
	if originalMsg.Success != decodedMsg.Success {
		t.Errorf("Success mismatch: got %v, want %v", decodedMsg.Success, originalMsg.Success)
	}
	if originalMsg.LanAddress != decodedMsg.LanAddress {
		t.Errorf("LanAddress mismatch: got %v, want %v", decodedMsg.LanAddress, originalMsg.LanAddress)
	}
	if originalMsg.Msg != decodedMsg.Msg {
		t.Errorf("Msg mismatch: got %v, want %v", decodedMsg.Msg, originalMsg.Msg)
	}
}
