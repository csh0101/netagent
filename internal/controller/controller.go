package controller

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/csh0101/netagent.git/internal/protocol"
	"github.com/csh0101/netagent.git/internal/tcp"
	"github.com/csh0101/netagent.git/internal/util"
)

type Config struct {
	ControlPort uint16
	DataPort    int
	Cidr        string
}

type Controller struct {
	nodeManager *NodeManager
}

// GetConnection implements common.ConnectionRetriever.
func (c *Controller) GetConnection(address string) net.Conn {
	ip := strings.Split(address, ":")[0]
	conn := c.nodeManager.GetNodeDataTunnelByIp(ip)
	return conn
}

func (c *Controller) Run(config *Config) error {

	// api sever run
	go RunApiServer(c)

	if ipManager, err := NewIPManager(config.Cidr); err != nil {
		return err
	} else {
		c.nodeManager = NewNodeManager(ipManager)
		go c.nodeManager.RunAliveCheckBackground()
	}

	if err := RunDataSrv(c.DataFunc(), &tcp.ServerConfig{
		Port: config.DataPort,
	}); err != nil {
		return err
	}

	controlSrv := tcp.NewServer(tcp.NameServerOption("control-server"))

	controlSrv.RegisterHandler(new(protocol.HeartbeatMessage).MsgType(), c.HeartbeatFunc())
	controlSrv.RegisterHandler(new(protocol.RegisterReqMessage).MsgType(), c.RegisterFunc())
	// todo register handler
	fmt.Println(config.ControlPort)
	go controlSrv.Run(&tcp.ServerConfig{
		Port:    int(config.ControlPort),
		Address: "0.0.0.0",
		TlsKey:  "",
		TlsPem:  "",
	})

	return nil
}

func (c *Controller) DataFunc() tcp.HandleFunc {

	return func(ctx context.Context, buf []byte, src net.Conn) ([]byte, net.Conn, error) {
		msg := &protocol.DataMessage{}
		if err := msg.Decode(buf); err != nil {
			return nil, nil, err
		}
		downstream := c.GetConnection(msg.DstAddr)
		if downstream == nil {
			ip := strings.Split(msg.SrcAddr, ":")[0]
			// reverse addr for register
			if msg.DstAddr == "0.0.0.1:65535" {
				c.nodeManager.SetNodeDataTunnelByIp(src, ip)
				return nil, nil, nil
			} else {
				resp := protocol.DataMessage{
					RequestID: msg.RequestID,
					ClientID:  msg.ClientID,
					TunnelID:  msg.TunnelID,
					Error:     "dst not out live",
				}
				if buf, err := resp.Encode(); err != nil {
					return nil, nil, err
				} else {
					if _, err := src.Write(buf); err != nil {
						return nil, nil, err
					} else {
						return buf, src, nil
					}
				}
			}
		}
		return buf, downstream, nil
	}
}

func (c *Controller) RegisterFunc() tcp.HandleFunc {
	return func(ctx context.Context, buf []byte, src net.Conn) ([]byte, net.Conn, error) {

		errResp := &protocol.RegisterRespMessage{}
		req := &protocol.RegisterReqMessage{}
		resp := &protocol.RegisterRespMessage{}
		{
			if err := req.Decode(buf); err != nil {
				errResp.Success = false
				errResp.Msg = "data can't decode"
				goto ERROR
			}

			node, err := c.nodeManager.Add(req.AgentName, src)
			if err != nil {
				errResp.Success = false
				errResp.Msg = "agent name is not unique"
				goto ERROR
			}
			resp.ClientID = node.id
			resp.LanAddress, _ = util.IPToUint32(node.ip)
		}

		{
			resp.Success = true
			resp.RequestID = req.RequestID
			buf, err := resp.Encode()
			if err != nil {
				return nil, nil, err
			}
			msg := &protocol.RegisterRespMessage{}

			if err := msg.Decode(buf); err != nil {
				panic(err)
			}
			return buf, src, nil
		}
	ERROR:
		buf, err := errResp.Encode()
		if err != nil {
			return nil, nil, err
		}
		buf, err = protocol.PackMessage(new(protocol.RegisterRespMessage).MsgType(), buf)
		if err != nil {
			return nil, nil, err
		}
		return buf, src, nil
	}
}

func (c *Controller) HeartbeatFunc() tcp.HandleFunc {
	return func(ctx context.Context, data []byte, src net.Conn) ([]byte, net.Conn, error) {
		req := &protocol.HeartbeatMessage{}
		if err := req.Decode(data); err != nil {
			//
			return nil, nil, err
		}

		c.nodeManager.Update(req.AgentName)

		fmt.Println("receive heart..", req.AgentName)

		return nil, nil, nil

	}
}

// RetrieveAllNode implements apiserver.DataSource.
func (c *Controller) RetrieveAllNode() ([]*NodeResp, error) {
	nodes := make([]*NodeResp, 0)
	c.nodeManager.Foreach(func(n *Node) {
		resp := &NodeResp{
			Name:     n.name,
			IP:       n.ip.String(),
			LastSeen: n.lastSeen.String(),
		}
		nodes = append(nodes, resp)
	})
	return nodes, nil
}

func RunDataSrv(f tcp.HandleFunc, conf *tcp.ServerConfig) error {
	s := tcp.NewServer(tcp.NameServerOption("data-server"))
	s.RegisterHandler(new(protocol.DataMessage).MsgType(), f)
	go s.Run(&tcp.ServerConfig{
		Port:    int(conf.Port),
		Address: "0.0.0.0",
		TlsKey:  conf.TlsKey,
		TlsPem:  conf.TlsPem,
	})
	return nil
}

type NodeResp struct {
	IP       string `json:"ip"`
	Name     string `json:"name"`
	LastSeen string `json:"last_seen"`
}
