package protocol

import (
	"bytes"
	"errors"

	"github.com/google/uuid"
)

type RegisterReqMessage struct {
	RequestID uuid.UUID
	AgentName string
}

// Decode implements tcp.MessageProtocol.
func (r *RegisterReqMessage) Decode(data []byte) error {
	if len(data) < 16 { // UUID is 16 bytes
		return errors.New("insufficient data to decode RegisterReqMessage")
	}

	r.RequestID = uuid.UUID{}
	copy(r.RequestID[:], data[:16])

	agentNameLength := len(data[16:])
	r.AgentName = string(data[16 : 16+agentNameLength])

	return nil
}

// Encode implements tcp.MessageProtocol.
func (r *RegisterReqMessage) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.Write(r.RequestID[:])
	buf.Write([]byte(r.AgentName))
	return PackMessage(r.MsgType(), buf.Bytes())
}

// MsgType implements tcp.MessageProtocol.
func (r *RegisterReqMessage) MsgType() uint32 {
	return 1 // Assuming 1 is the message type for RegisterReqMessage
}
