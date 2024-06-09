// package main

// import (
// 	"fmt"
// 	"io"
// 	"net"
// 	"os"
// 	"syscall"
// 	"unsafe"
// )

// func main() {
// 	listener, err := net.Listen("tcp", ":9999")
// 	if err != nil {
// 		fmt.Println("Error listening:", err)
// 		os.Exit(1)
// 	}
// 	defer listener.Close()

// 	fmt.Println("Listening on :9999")

// 	for {
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			fmt.Println("Error accepting connection:", err)
// 			continue
// 		}

// 		go handleClient(conn)
// 	}
// }

// func handleClient(clientConn net.Conn) {
// 	defer clientConn.Close()

// 	// 获取原始目的地址和端口
// 	rawConn, err := clientConn.(*net.TCPConn).SyscallConn()
// 	if err != nil {
// 		fmt.Println("Error getting raw connection:", err)
// 		return
// 	}

// 	var SO_ORIGINAL_DST = 80
// 	var originalDst syscall.Sockaddr
// 	err = rawConn.Control(func(fd uintptr) {
// 		originalDst, err = getsockopt(fd, syscall.IPPROTO_IP, SO_ORIGINAL_DST)
// 	})
// 	if err != nil {
// 		fmt.Println("Error getting original destination:", err)
// 		return
// 	}

// 	// 转发连接
// 	switch addr := originalDst.(type) {
// 	case *syscall.SockaddrInet4:
// 		targetAddr := fmt.Sprintf("%s:%d", net.IP(addr.Addr[:]).String(), addr.Port)
// 		fmt.Println("xyz")
// 		fmt.Println(targetAddr)
// 		// forward(clientConn, targetAddr)
// 	case *syscall.SockaddrInet6:
// 		targetAddr := fmt.Sprintf("[%s]:%d", net.IP(addr.Addr[:]).String(), addr.Port)
// 		fmt.Println(targetAddr)
// 		// forward(clientConn, targetAddr)
// 	default:
// 		fmt.Println("Unsupported address type")
// 	}
// }

// func getsockopt(fd uintptr, level, optname int) (syscall.Sockaddr, error) {
// 	var addr syscall.RawSockaddrAny
// 	var size uintptr = unsafe.Sizeof(addr)
// 	_, _, errno := syscall.Syscall6(syscall.SYS_GETSOCKOPT, fd, uintptr(level), uintptr(optname), uintptr(unsafe.Pointer(&addr)), uintptr(unsafe.Pointer(&size)), 0)
// 	if errno != 0 {
// 		return nil, errno
// 	}

// 	return sockaddrFromAny(&addr), nil
// }

// func sockaddrFromAny(rsa *syscall.RawSockaddrAny) syscall.Sockaddr {
// 	switch rsa.Addr.Family {
// 	case syscall.AF_INET:
// 		pp := (*syscall.RawSockaddrInet4)(unsafe.Pointer(rsa))
// 		addr := new(syscall.SockaddrInet4)
// 		addr.Port = int(pp.Port)
// 		addr.Addr = pp.Addr
// 		return addr
// 	case syscall.AF_INET6:
// 		pp := (*syscall.RawSockaddrInet6)(unsafe.Pointer(rsa))
// 		addr := new(syscall.SockaddrInet6)
// 		addr.Port = int(pp.Port)
// 		addr.Addr = pp.Addr
// 		addr.ZoneId = pp.Scope_id
// 		return addr
// 	}
// 	return nil
// }

// func forward(src net.Conn, targetAddr string) {
// 	dst, err := net.Dial("tcp", targetAddr)
// 	if err != nil {
// 		fmt.Println("Error connecting to target:", err)
// 		return
// 	}
// 	defer dst.Close()

// 	go func() {
// 		io.Copy(dst, src)
// 	}()
// 	io.Copy(src, dst)
// }

package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"syscall"
	"unsafe"
)

func main() {
	listener, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println("Error listening:", err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Println("Listening on :9999")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleClient(conn)
	}
}

func handleClient(clientConn net.Conn) {
	defer clientConn.Close()

	// 获取原始目的地址和端口
	rawConn, err := clientConn.(*net.TCPConn).SyscallConn()
	if err != nil {
		fmt.Println("Error getting raw connection:", err)
		return
	}

	var SO_ORIGINAL_DST = 80
	var originalDst syscall.Sockaddr
	err = rawConn.Control(func(fd uintptr) {
		originalDst, err = getsockopt(fd, syscall.IPPROTO_IP, SO_ORIGINAL_DST)
	})
	if err != nil {
		fmt.Println("Error getting original destination:", err)
		return
	}

	// 转发连接
	switch addr := originalDst.(type) {
	case *syscall.SockaddrInet4:
		port := ntohs(uint16(addr.Port))
		targetAddr := fmt.Sprintf("%s:%d", net.IP(addr.Addr[:]).String(), port)
		fmt.Println(targetAddr)
		// forward(clientConn, targetAddr)
	case *syscall.SockaddrInet6:
		port := ntohs(uint16(addr.Port))
		targetAddr := fmt.Sprintf("[%s]:%d", net.IP(addr.Addr[:]).String(), port)
		fmt.Println(targetAddr)
		// forward(clientConn, targetAddr)
	default:
		fmt.Println("Unsupported address type")
	}
}

func getsockopt(fd uintptr, level, optname int) (syscall.Sockaddr, error) {
	var addr syscall.RawSockaddrAny
	var size uintptr = unsafe.Sizeof(addr)
	_, _, errno := syscall.Syscall6(syscall.SYS_GETSOCKOPT, fd, uintptr(level), uintptr(optname), uintptr(unsafe.Pointer(&addr)), uintptr(unsafe.Pointer(&size)), 0)
	if errno != 0 {
		return nil, errno
	}

	return sockaddrFromAny(&addr), nil
}

func sockaddrFromAny(rsa *syscall.RawSockaddrAny) syscall.Sockaddr {
	switch rsa.Addr.Family {
	case syscall.AF_INET:
		pp := (*syscall.RawSockaddrInet4)(unsafe.Pointer(rsa))
		addr := new(syscall.SockaddrInet4)
		addr.Port = int(pp.Port)
		addr.Addr = pp.Addr
		return addr
	case syscall.AF_INET6:
		pp := (*syscall.RawSockaddrInet6)(unsafe.Pointer(rsa))
		addr := new(syscall.SockaddrInet6)
		addr.Port = int(pp.Port)
		addr.Addr = pp.Addr
		addr.ZoneId = pp.Scope_id
		return addr
	}
	return nil
}

func ntohs(n uint16) uint16 {
	return (n>>8)&0xff | (n&0xff)<<8
}

func forward(src net.Conn, targetAddr string) {
	dst, err := net.Dial("tcp", targetAddr)
	if err != nil {
		fmt.Println("Error connecting to target:", err)
		return
	}
	defer dst.Close()

	go func() {
		io.Copy(dst, src)
	}()
	io.Copy(src, dst)
}
