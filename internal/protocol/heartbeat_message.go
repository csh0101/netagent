package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/google/uuid"
)

type HeartbeatMessage struct {
	RequestID  uuid.UUID
	ClientID   uint32
	AgentName  string
	LanAddress uint32
}

// Decode implements tcp.MessageProtocol.
func (h *HeartbeatMessage) Decode(data []byte) error {
	if len(data) < 24 { // UUID is 16 bytes, ClientID is 4 bytes, and LanAddress is 4 bytes
		return errors.New("insufficient data to decode HeartbeatMessage")
	}

	h.RequestID = uuid.UUID{}
	copy(h.RequestID[:], data[:16])
	h.ClientID = binary.BigEndian.Uint32(data[16:20])
	h.LanAddress = binary.BigEndian.Uint32(data[20:24])
	agentNameLength := len(data[24:])
	h.AgentName = string(data[24 : 24+agentNameLength])

	return nil
}

// Encode implements tcp.MessageProtocol.
func (h *HeartbeatMessage) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.Write(h.RequestID[:])
	ClientIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(ClientIDBytes, h.ClientID)
	buf.Write(ClientIDBytes)
	LanAddressBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(LanAddressBytes, h.LanAddress)
	buf.Write(LanAddressBytes)
	buf.Write([]byte(h.AgentName))
	return PackMessage(h.MsgType(), buf.Bytes())
}

// MsgType implements tcp.MessageProtocol.
func (h *HeartbeatMessage) MsgType() uint32 {
	return 10
}
