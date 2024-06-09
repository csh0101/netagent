package tcp

type MessageType int32

type MessageProtocol interface {
	Encode() ([]byte, error)
	Decode(data []byte) error
	MsgType() uint32
}
