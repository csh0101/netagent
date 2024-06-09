package controller

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net"
	"sync"
)

type IPManager struct {
	cidr       *net.IPNet
	baseIP     net.IP
	allocated  map[string]struct{}
	mu         sync.RWMutex
	totalAddrs *big.Int
}

// NewIPManager 创建一个新的IPManager
func NewIPManager(cidrStr string) (*IPManager, error) {
	_, cidr, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return nil, fmt.Errorf("invalid CIDR: %v", err)
	}

	// 计算CIDR段内的总地址数
	ones, bits := cidr.Mask.Size()
	totalAddrs := big.NewInt(1)
	totalAddrs.Lsh(totalAddrs, uint(bits-ones))

	return &IPManager{
		cidr:       cidr,
		baseIP:     cidr.IP,
		allocated:  make(map[string]struct{}),
		totalAddrs: totalAddrs,
	}, nil
}

// GenerateUniqueIP 生成唯一的IP地址
func (im *IPManager) GenerateUniqueIP() (net.IP, error) {
	im.mu.Lock()
	defer im.mu.Unlock()

	if len(im.allocated) >= int(im.totalAddrs.Int64()-2) { // 减去网络地址和广播地址
		return nil, errors.New("no available IP addresses")
	}

	for {
		ip := make(net.IP, len(im.baseIP))
		copy(ip, im.baseIP)

		// 生成随机偏移
		offset, err := rand.Int(rand.Reader, im.totalAddrs)
		if err != nil {
			return nil, fmt.Errorf("failed to generate random offset: %v", err)
		}

		// 跳过网络地址和广播地址
		offset.Add(offset, big.NewInt(1))
		if offset.Cmp(im.totalAddrs) >= 0 {
			offset.Sub(offset, big.NewInt(2))
		}

		for i := len(ip) - 1; i >= 0; i-- {
			ip[i] += byte(offset.Uint64() & 0xFF)
			offset.Rsh(offset, 8)
		}

		ipStr := ip.String()
		if _, exists := im.allocated[ipStr]; !exists {
			im.allocated[ipStr] = struct{}{}
			return ip, nil
		}
	}
}

func (im *IPManager) ReleaseIP(ip net.IP) string {
	im.mu.Lock()
	defer im.mu.Unlock()

	ipStr := ip.String()
	delete(im.allocated, ipStr)
	return ipStr
}
