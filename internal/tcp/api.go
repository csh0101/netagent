package tcp

import (
	"context"
	"net"
)

type HandleFunc func(context.Context, []byte, net.Conn) ([]byte, net.Conn, error)
