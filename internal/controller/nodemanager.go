package controller

import (
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type NodeManager struct {
	nodes     map[string]*Node
	ipmap     map[string]*Node
	ipManager *IPManager
	acc       uint32
	mu        *sync.RWMutex
}

type Node struct {
	ip       net.IP
	id       uint32
	name     string
	tunnel   net.Conn
	lastSeen time.Time
}

func NewNodeManager(ipManager *IPManager) *NodeManager {
	return &NodeManager{
		nodes:     make(map[string]*Node),
		ipmap:     make(map[string]*Node),
		ipManager: ipManager,
		mu:        &sync.RWMutex{},
	}
}

// Add 添加一个新的节点
// name 应该是唯一的
func (m *NodeManager) Add(name string, conn net.Conn) (*Node, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.nodes[name]; exists {
		return nil, errors.New("node name already exists")
	}

	ip, err := m.ipManager.GenerateUniqueIP()
	if err != nil {
		return nil, err
	}

	id := atomic.AddUint32(&m.acc, 1)

	node := &Node{
		ip:       ip,
		id:       id,
		name:     name,
		lastSeen: time.Now(),
	}

	m.nodes[name] = node
	m.ipmap[node.ip.String()] = node
	return node, nil
}

func (m *NodeManager) GetNodeDataTunnelByIp(ip string) net.Conn {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if _, ok := m.ipmap[ip]; !ok {
		return nil
	}
	// todo double check ?
	return m.ipmap[ip].tunnel
}

func (m *NodeManager) SetNodeDataTunnelByIp(dataTunnel net.Conn, ip string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exist := m.ipmap[ip]; exist {
		node := m.ipmap[ip]
		node.tunnel = dataTunnel
	}
}

// Remove 移除一个节点
func (m *NodeManager) Remove(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if node, exists := m.nodes[name]; exists {
		delete(m.nodes, name)
		// 释放已分配的IP地址
		m.ipManager.ReleaseIP(node.ip)
	}
}

func (m *NodeManager) Update(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.nodes[name].lastSeen = time.Now()
}

func (m *NodeManager) Foreach(f func(*Node)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, node := range m.nodes {
		f(node)
	}
}

func (m *NodeManager) RunAliveCheckBackground() error {
	for {
		now := time.Now()
		m.Foreach(func(node *Node) {
			if now.Sub(node.lastSeen).Seconds() > 60 {
				// 这里假设有人在往这个tunnel写会不会panic?
				if node.tunnel != nil {
					node.tunnel.Close()
				}
				ip := m.ipManager.ReleaseIP(node.ip)
				delete(m.ipmap, ip)
				delete(m.nodes, node.name)
			}
		})
		time.Sleep(time.Second * 60)
	}
}
