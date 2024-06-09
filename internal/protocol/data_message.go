package protocol

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/google/uuid"
)

type DataMessage struct {
	RequestID uuid.UUID
	ClientID  uint32
	TunnelID  uint32
	DstAddr   string
	SrcAddr   string
	Error     string
	Data      []byte
}

// Decode implements tcp.MessageProtocol.
func (s *DataMessage) Decode(data []byte) error {
	buf := bytes.NewReader(data)

	// Read RequestID (16 bytes)
	uuidBytes := make([]byte, 16)
	if _, err := io.ReadFull(buf, uuidBytes); err != nil {
		return err
	}
	requestID, err := uuid.FromBytes(uuidBytes)
	if err != nil {
		return err
	}
	s.RequestID = requestID

	// Read ClientID and TunnelID (4 bytes each)
	if err := binary.Read(buf, binary.BigEndian, &s.ClientID); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &s.TunnelID); err != nil {
		return err
	}

	// Read DstAddr
	dstAddrLen, err := buf.ReadByte()
	if err != nil {
		return err
	}
	dstAddr := make([]byte, dstAddrLen)
	if _, err := io.ReadFull(buf, dstAddr); err != nil {
		return err
	}
	s.DstAddr = string(dstAddr)

	// Read SrcAddr
	srcAddrLen, err := buf.ReadByte()
	if err != nil {
		return err
	}
	srcAddr := make([]byte, srcAddrLen)
	if _, err := io.ReadFull(buf, srcAddr); err != nil {
		return err
	}
	s.SrcAddr = string(srcAddr)

	// Read Error
	errorLen, err := buf.ReadByte()
	if err != nil {
		return err
	}
	errorMsg := make([]byte, errorLen)
	if _, err := io.ReadFull(buf, errorMsg); err != nil {
		return err
	}
	s.Error = string(errorMsg)

	// Read Data
	dataLen := buf.Len()
	s.Data = make([]byte, dataLen)
	if _, err := io.ReadFull(buf, s.Data); err != nil {
		return err
	}

	return nil
}

// Encode implements tcp.MessageProtocol.
func (s *DataMessage) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)

	uuidBytes, err := s.RequestID.MarshalBinary()
	if err != nil {
		return nil, err
	}
	if _, err := buf.Write(uuidBytes); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, s.ClientID); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, s.TunnelID); err != nil {
		return nil, err
	}

	if err := buf.WriteByte(byte(len(s.DstAddr))); err != nil {
		return nil, err
	}
	if _, err := buf.Write([]byte(s.DstAddr)); err != nil {
		return nil, err
	}

	if err := buf.WriteByte(byte(len(s.SrcAddr))); err != nil {
		return nil, err
	}
	if _, err := buf.Write([]byte(s.SrcAddr)); err != nil {
		return nil, err
	}

	if err := buf.WriteByte(byte(len(s.Error))); err != nil {
		return nil, err
	}
	if _, err := buf.Write([]byte(s.Error)); err != nil {
		return nil, err
	}

	if _, err := buf.Write(s.Data); err != nil {
		return nil, err
	}
	return PackMessage(new(DataMessage).MsgType(), buf.Bytes())
}

// MsgType implements tcp.MessageProtocol.
func (s *DataMessage) MsgType() uint32 {
	return 1 // 假设DataMessage的消息类型为1
}
