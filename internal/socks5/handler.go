package socks5

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
)

type ProxyProcessFunc func(src net.Conn, addr string, port uint16) error

var TransParentProxyProcessFunc ProxyProcessFunc = func(src net.Conn, addr string, port uint16) error {

	if addr == "" || port == 0 {
		// todo add wraning log
		fmt.Println("addr is empty or port is not correct")
		return errors.New("addr")
	}
	destAddrPort := fmt.Sprintf("%s:%d", addr, port)
	// Conn
	dest, err := net.Dial("tcp", destAddrPort)
	if err != nil {
		fmt.Println("dial dst: " + err.Error())
		return err
	}
	forward := func(src, dest net.Conn) {
		defer src.Close()
		defer dest.Close()
		io.Copy(src, dest)
	}
	go forward(src, dest)
	go forward(dest, src)
	return nil
}

func Socks5Auth(conn net.Conn) (err error) {
	buf := make([]byte, 256)

	n, err := io.ReadFull(conn, buf[:2])

	if err != nil {
		return errors.New("reading header: " + err.Error())
	}

	if n != 2 {
		return errors.New("unexcepted data length")
	}

	version, nMethods := int(buf[0]), int(buf[1])

	if version != 5 {
		return errors.New("invaild version")
	}

	n, err = io.ReadFull(conn, buf[:nMethods])
	if n != nMethods {
		return errors.New("reading methods: " + err.Error())
	}

	//    The SOCKS request information is sent by the client as soon as it has
	//    established a connection to the SOCKS server, and completed the
	//    authentication negotiations.  The server evaluates the request, and
	//    returns a reply formed as follows
	//  +----+-----+-------+------+----------+----------+
	// |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
	// +----+-----+-------+------+----------+----------+
	// | 1  |  1  | X'00' |  1   | Variable |    2     |
	// +----+-----+-------+------+----------+----------+

	// todo specify it, version = 0x05, 0x00 means reprensent
	n, err = conn.Write([]byte{0x05, 0x00})

	if n != 2 || err != nil {
		return errors.New("write rsp: " + err.Error())
	}

	return nil
}

func Socks5Connect(conn net.Conn) (dst string, port uint16, err error) {

	buf := make([]byte, 256)

	n, err := io.ReadFull(conn, buf[:4])

	if n != 4 {
		return "", 0, errors.New("read header: " + err.Error())
	}

	ver, cmd, _, atyp := buf[0], buf[1], buf[2], buf[3]

	if ver != 5 || cmd != 1 {
		return "", 0, errors.New("invalid ver/cmd")
	}

	addr := ""

	switch atyp {
	case 1:
		n, err := io.ReadFull(conn, buf[:4])
		if n != 4 {
			return "", 0, errors.New("invalid IPv4: " + err.Error())
		}
		addr = fmt.Sprintf("%d.%d.%d.%d", buf[0], buf[1], buf[2], buf[3])
	case 3:
		n, err := io.ReadFull(conn, buf[:1])
		if n != 1 {
			return "", 0, errors.New("invalid hostname: " + err.Error())
		}
		addrLen := int(buf[0])

		n, err = io.ReadFull(conn, buf[:addrLen])
		if n != addrLen {
			return "", 0, errors.New("invalid hostname: " + err.Error())
		}
		addr = string(buf[:addrLen])
	case 4:
		return "", 0, errors.New("IPv6: no supported yet")
	default:
		return "", 0, errors.New("invalid atyp")
	}

	n, err = io.ReadFull(conn, buf[:2])

	if n != 2 {
		return "", 0, errors.New("read port:  " + err.Error())
	}

	port = binary.BigEndian.Uint16((buf[:2]))

	// write rsp
	_, err = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	if err != nil {
		return "", 0, errors.New("write rsp: " + err.Error())
	}

	return addr, port, nil
}
