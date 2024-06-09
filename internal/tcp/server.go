package tcp

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"net"

	"github.com/csh0101/netagent.git/internal/util"
)

type ServerConfig struct {
	Port    int
	Address string
	TlsKey  string
	TlsPem  string
}

func (config *ServerConfig) tcpAddressTuple() string {
	return fmt.Sprintf("%s:%d", config.Address, config.Port)
}

func (config *ServerConfig) EnableTLS() bool {
	return config.TlsKey != "" && config.TlsPem != ""
}

type HandlersMap map[uint32]HandleFunc

func (handlers HandlersMap) registerHandler(route uint32, handler HandleFunc) error {

	if _, ok := handlers[route]; ok {
		// todo add log
		return errors.New("tcp handler route is not unique, please check")
	}
	// handler
	handlers[route] = handler

	return nil
}

func (handlers HandlersMap) doHandler(ctx context.Context, conn net.Conn) error {

	// todo , for feature extension by ctx
	buf := util.GetBuf()
	n, err := conn.Read(buf[:4])
	if err != nil {
		return err
	}

	// todo because of the msg type is a int32 value, parse to []byte 4
	if n != 4 {
		return errors.New(" ")
	}

	// handlers
	msgType := binary.BigEndian.Uint32(buf[:4])

	handler := handlers[msgType]

	if handler == nil {
		// todo handler is nil , log warning and return
		return errors.New("unmap handler! please check msgType")
	}

	n, err = conn.Read(buf[:4])
	if err != nil {
		// todo add logger
		return err
	}
	if n != 4 {
		return errors.New("unexcepted data length")
	}

	length := binary.BigEndian.Uint32(buf[:4])
	if _, err := conn.Read(buf[:length]); err != nil {
		return err
	}

	if buf, target, err := handler(ctx, buf[:length], conn); err != nil {
		// add log
		return err
	} else if target != nil && buf != nil {
		// "\x00\x00\x00\x02\x00\x00\x00\x19 \xa5\v\xe6<\xe4I\x95\x8e\x12{\x0f\x8be\xf9\xf5\x00\x00\x00\x01\x01\xac\x02\x01\x0f"
		_, err := target.Write(buf)

		if err != nil {
			// add log
			fmt.Println("target write buf err:", err.Error())
		}
	}

	buf = buf[:]
	util.PutBuf(buf)
	return nil
}

type ServerOption func(s *Server)

func NameServerOption(name string) ServerOption {
	return func(s *Server) {
		s.Name = name
	}
}

type Server struct {
	Name     string
	ctx      context.Context
	handlers HandlersMap
}

func NewServer(opts ...ServerOption) *Server {
	s := &Server{
		ctx:      context.Background(),
		handlers: make(HandlersMap),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Server) Run(c *ServerConfig) {
	s.ctx = context.Background()
	// todo if enable tls, deal it
	if c.EnableTLS() {

	} else {
		listener, err := net.Listen("tcp", c.tcpAddressTuple())
		if err != nil {
			fmt.Println("Error starting server: ", err)
			return
		}
		defer listener.Close()
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Error accepting connection:", err)
				panic(err)
			}
			go s.dealConnection(context.Background(), conn)
		}
	}

}

func (s *Server) RegisterHandler(route uint32, handler HandleFunc) {
	s.handlers.registerHandler(route, handler)
}

func (s *Server) dealConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	for {
		err := s.handlers.doHandler(ctx, conn)
		if err != nil {
			fmt.Printf("%s occur err %s,then close conn\n", s.Name, err.Error())
			// todo are there to split
			conn.Close()
			return
		}
	}
}
