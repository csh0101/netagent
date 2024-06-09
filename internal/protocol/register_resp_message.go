package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/google/uuid"
)

type RegisterRespMessage struct {
	RequestID  uuid.UUID
	ClientID   uint32
	Success    bool
	LanAddress uint32
	Msg        string
}

// Decode implements tcp.MessageProtocol.
func (r *RegisterRespMessage) Decode(data []byte) error {
	if len(data) < 25 { // UUID (16 bytes) + ClientID (4 bytes) + Success (1 byte) + LanAddress (4 bytes)
		return errors.New("insufficient data to decode RegisterRespMessage")
	}

	r.RequestID = uuid.UUID{}
	copy(r.RequestID[:], data[:16])

	r.ClientID = binary.BigEndian.Uint32(data[16:20])

	r.Success = data[20] == 1

	r.LanAddress = binary.BigEndian.Uint32(data[21:25])

	if len(data) > 25 {
		r.Msg = string(data[25:])
	} else {
		r.Msg = ""
	}

	return nil
}

// Encode implements tcp.MessageProtocol.
func (r *RegisterRespMessage) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.Write(r.RequestID[:])
	clientIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(clientIDBytes, r.ClientID)
	buf.Write(clientIDBytes)
	successByte := byte(0)
	if r.Success {
		successByte = 1
	}
	buf.WriteByte(successByte)
	lanAddressBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(lanAddressBytes, r.LanAddress)
	buf.Write(lanAddressBytes)
	buf.Write([]byte(r.Msg))
	return PackMessage(r.MsgType(), buf.Bytes())
}

// MsgType implements tcp.MessageProtocol.
func (r *RegisterRespMessage) MsgType() uint32 {
	return 2 // Assuming 2 is the message type for RegisterRespMessage
}
