package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/net/proxy"
)

func main() {

	dst := os.Getenv("VPN_TARGET")

	// SOCKS5 代理地址
	socks5Addr := "localhost:1080" // 替换为实际的 SOCKS5 代理地址

	// 代理认证信息（如果需要）
	auth := proxy.Auth{
		// User:     "username", // 替换为实际的用户名
		// Password: "password", // 替换为实际的密码
	}

	// 创建一个 SOCKS5 代理拨号器
	dialer, err := proxy.SOCKS5("tcp", socks5Addr, &auth, proxy.Direct)
	if err != nil {
		log.Fatal("Error creating SOCKS5 dialer:", err)
	}

	target := fmt.Sprintf("%s:5555", dst)

	// 使用代理拨号器连接到服务端
	conn, err := dialer.Dial("tcp", target)
	if err != nil {
		log.Fatal("Error connecting to server via SOCKS5 proxy:", err)
	}
	defer conn.Close()

	for {

		message := "Hello Wrold"

		// 发送消息到服务端
		_, err = conn.Write([]byte(message))
		if err != nil {
			log.Fatal("Error writing to server:", err)
		}

		// 接收服务端回显的消息
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			log.Fatal("Error reading from server:", err)
		}

		// 打印服务端回显的消息
		fmt.Println("Received from server:", string(buf[:n]))
	}

}
