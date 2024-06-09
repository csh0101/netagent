package protocol

import (
	"bytes"
	"testing"

	"github.com/google/uuid"
)

func TestDataMessage_EncodeDecode(t *testing.T) {
	// Original DataMessage
	originalMsg := &DataMessage{
		RequestID: uuid.New(),
		ClientID:  12345,
		TunnelID:  67890,
		DstAddr:   "192.168.0.1:8080",
		SrcAddr:   "10.0.0.1:9090",
		Error:     "",
		Data:      []byte("test data"),
	}

	// Encode the original message
	encoded, err := originalMsg.Encode()
	if err != nil {
		t.Fatalf("Failed to encode message: %v", err)
	}

	// Decode the encoded message
	decodedMsg := &DataMessage{}
	err = decodedMsg.Decode(encoded)
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
	if originalMsg.TunnelID != decodedMsg.TunnelID {
		t.Errorf("TunnelID mismatch: got %v, want %v", decodedMsg.TunnelID, originalMsg.TunnelID)
	}
	if originalMsg.DstAddr != decodedMsg.DstAddr {
		t.Errorf("DstAddr mismatch: got %v, want %v", decodedMsg.DstAddr, originalMsg.DstAddr)
	}
	if originalMsg.SrcAddr != decodedMsg.SrcAddr {
		t.Errorf("SrcAddr mismatch: got %v, want %v", decodedMsg.SrcAddr, originalMsg.SrcAddr)
	}
	if originalMsg.Error != decodedMsg.Error {
		t.Errorf("Error mismatch: got %v, want %v", decodedMsg.Error, originalMsg.Error)
	}
	if !bytes.Equal(originalMsg.Data, decodedMsg.Data) {
		t.Errorf("Data mismatch: got %v, want %v", decodedMsg.Data, originalMsg.Data)
	}
}
