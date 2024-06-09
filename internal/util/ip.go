package util

import (
	"encoding/binary"
	"fmt"
	"net"
)

func IPToUint32(ip net.IP) (uint32, error) {
	ip = ip.To4()
	if ip == nil {
		return 0, fmt.Errorf("not an IPv4 address")
	}
	return binary.BigEndian.Uint32(ip), nil
}

// Uint32ToIP 将 uint32 转换为 net.IP
func Uint32ToIP(n uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, n)
	return ip
}
