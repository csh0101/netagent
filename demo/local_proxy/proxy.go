package main

import (
	"fmt"
	"net"
)

func main() {
	// 监听端口
	listener, err := net.Listen("tcp", ":18888")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server listening on port 18888")

	for {
		// 接受客户端连接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Client connected")

	// 读取客户端发送的数据
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading from connection:", err)
			return
		}
		if n == 0 {
			break
		}

		fmt.Printf("Received: %s\n", string(buf[:n]))

		// 响应客户端
		response := "Message received"
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing to connection:", err)
			return
		}
		fmt.Printf("Sent: %s\n", response)
	}
}
