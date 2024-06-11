package main

import (
	"io"
	"log"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Println("Read error:", err)
			}
			break
		}
		_, err = conn.Write(buf[:n])
		if err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":5555")
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
	defer listener.Close()
	log.Println("Server started on port 5555")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}
