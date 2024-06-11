package socks5

import (
	"context"
	"fmt"
	"net"
)

// RunkaSocks5 should exist for restart extension
func RunSocks5(ctx context.Context, port int, f ProxyProcessFunc) error {
	// todo logger
	fmt.Println("socks5 proxy server running on port: ", port)
	server, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
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
