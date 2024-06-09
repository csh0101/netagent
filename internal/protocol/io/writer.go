package io

import (
	"io"
	"net"

	"github.com/csh0101/netagent.git/internal/protocol"
	"github.com/google/uuid"
)

type WriterAlias func(p []byte) (n int, err error)

var _ io.Writer = new(WriterAlias)

// 实现 io.Writer 接口的 Write 方法
func (w WriterAlias) Write(p []byte) (n int, err error) {
	return w(p)
}

func DataMessageWriter(clientID, TunnelID uint32, srcAddr string, dstAddr string, target net.Conn) io.Writer {
	var f WriterAlias = func(p []byte) (n int, err error) {
		msg := &protocol.DataMessage{
			RequestID: uuid.New(),
			ClientID:  clientID,
			TunnelID:  TunnelID,
			SrcAddr:   srcAddr,
			DstAddr:   dstAddr,
			Data:      p,
		}

		buf, err := msg.Encode()
		if err != nil {
			return 0, err
		}

		if buf, err = protocol.PackMessage(new(protocol.DataMessage).MsgType(), buf); err != nil {
			return 0, err
		}

		if _, err = target.Write(buf); err != nil {
			return 0, err
		}

		return len(p), nil
	}
	return f
}
