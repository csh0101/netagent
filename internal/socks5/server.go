package socks5

import (
	"context"
	"fmt"
	"net"
)

// RunkaSocks5 should exist for restart extension
func RunSocks5(ctx context.Context, f ProxyProcessFunc) error {
	server, err := net.Listen("tcp", ":1080")
	if err != nil {
		fmt.Printf("Listen failed: %v\n", err)
		return nil
	}
	defer server.Close()

	for {
		conn, err := server.Accept()
		if err != nil {
			continue
		}
		go process(conn, f)
	}
}

func process(conn net.Conn, f ProxyProcessFunc) {
	if err := Socks5Auth(conn); err != nil {
		// todo add log
		fmt.Println("socks5 auth error: ", err.Error())
		return
	}
	addr, port, err := Socks5Connect(conn)
	if err != nil {
		return
	}
	f(conn, addr, port)
}
