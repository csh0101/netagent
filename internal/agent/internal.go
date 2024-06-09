package agent

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/csh0101/netagent.git/internal/protocol"

	protocl_io "github.com/csh0101/netagent.git/internal/protocol/io"
	"github.com/csh0101/netagent.git/internal/socks5"
	"github.com/csh0101/netagent.git/internal/tcp"
	"github.com/csh0101/netagent.git/internal/util"
	"github.com/google/uuid"
)

type Config struct {
	DataTunnelAddress    string
	DataTunnelPort       int
	ControlTunnelAddress string
	ControlTunnelPort    int
	Name                 string
	MaxRetries           int
}

type Agent struct {
	ctx       context.Context
	name      string
	privateIP uint32
	// todo! server need return clientID
	clientID uint32
	local    map[string]net.Conn
	control  net.Conn
	remote   net.Conn
	mutex    *sync.Mutex
	register util.AtomicBool
	conf     *Config
}

func (a *Agent) Run(conf *Config) error {

	// ext for restart , abstract object
	a.ctx = context.Background()
	a.local = make(map[string]net.Conn)
	a.mutex = &sync.Mutex{}
	a.name = conf.Name
	a.conf = conf

	var pipeline util.Pipeline = make(chan struct{})

	// todo 设计一下重连方案..

	if err := pipeline.Run(
		util.NewCallback("initControllerTunnel", func() (chan struct{}, error) {
			// todo initControlTunnel now is very simple,but for feature ext, it can hold on Pipeline Async logic.
			return a.initControlTunnel()
		}),
		util.NewCallback("RegisterAgent", func() (chan struct{}, error) {
			return a.registerAgent()
		}),
		util.NewCallback("InitDataTunnel", func() (chan struct{}, error) {
			return a.initDataTunnel()
		}),
	); err != nil {
		return err
	}
	// todo! this should't be a DataSrv, which is a real tcp server
	// todo! this should be a loop listern to a tcp socket..

	// start hearbeat goroutine
	go func() {
		for {
			err := SendHearbeat(a.control, &protocol.HeartbeatMessage{
				RequestID:  uuid.New(),
				ClientID:   a.clientID,
				AgentName:  a.name,
				LanAddress: a.privateIP,
			})
			if err != nil {
				// add log warnning
				fmt.Println("err when send heart beat", err.Error())
				continue
			}
			time.Sleep(time.Second * 10)
		}
	}()
	go runDataLoop(a.DataFunc(), a.remote)
	go socks5.RunSocks5(a.ctx, func(src net.Conn, addr string, port uint16) error {
		// get local port
		sourcePort := strings.Split(src.LocalAddr().String(), ":")[1]
		localAddr := fmt.Sprintf("%s:%s", util.Uint32ToIP(a.privateIP).String(), sourcePort)
		writeBack := protocl_io.DataMessageWriter(a.clientID, 0, localAddr, fmt.Sprintf("%s:%d", addr, port), a.remote)
		a.mutex.Lock()
		a.local[localAddr] = src
		a.mutex.Unlock()
		go func() {
			defer src.Close()
			io.Copy(writeBack, src)
		}()
		return nil
	})
	return nil
}

func (a *Agent) getConnection(k string) net.Conn {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	// 一种是这次请求本来就是由我发起，所以这个连接已经维持在本地
	// 对于这种情况，只需要拿到原来得了连接，把数据写进去就可以了
	if _, exist := a.local[k]; exist {
		return a.local[k]
	}
	return nil
}

func (a *Agent) DataFunc() tcp.HandleFunc {
	return func(ctx context.Context, buf []byte, src net.Conn) ([]byte, net.Conn, error) {
		msg := &protocol.DataMessage{}
		if err := msg.Decode(buf); err != nil {
			return nil, nil, err
		}
		local := a.getConnection(msg.DstAddr)
		dst := protocl_io.DataMessageWriter(a.clientID, msg.TunnelID, msg.DstAddr, msg.SrcAddr, src)
		if local == nil {
			ip := "127.0.0.1"
			port := strings.Split(msg.DstAddr, ":")[1]
			address := fmt.Sprintf("%s:%s", ip, port)
			src, err := net.Dial("tcp", address)
			if err != nil {
				// todo log error
				resp := protocol.DataMessage{
					RequestID: msg.RequestID,
					ClientID:  msg.ClientID,
					TunnelID:  msg.TunnelID,
					Error:     "local process is not exist" + err.Error(),
				}
				if buf, err := resp.Encode(); err != nil {
					return nil, nil, err
				} else {
					return buf, src, nil
				}
			}
			// assign src to local
			local = src
			a.mutex.Lock()
			a.local[address] = src
			a.mutex.Unlock()
			go func(src net.Conn, dst io.Writer) {
				defer src.Close()
				io.Copy(dst, src)
			}(src, dst)
		}
		if _, err := local.Write(msg.Data); err != nil {
			return nil, nil, err
		}
		return nil, nil, nil
	}
}

func (a *Agent) initControlTunnel() (chan struct{}, error) {
	finished, err := util.GetTaskExecutor().RunTask(context.Background(), "initControluTunnel", func() error {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", a.conf.ControlTunnelAddress, a.conf.ControlTunnelPort))
		if err != nil {
			fmt.Println("err when dial control tcp", err.Error())
		}
		a.control = conn
		go a.runControlLoop(a.control)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return finished, nil
}

func (a *Agent) initDataTunnel() (chan struct{}, error) {
	finished, err := util.GetTaskExecutor().RunTask(context.Background(), "initDataTunnel", func() error {
		var err error
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", a.conf.DataTunnelAddress, a.conf.DataTunnelPort))
		if err != nil {
			return err
		}
		a.remote = conn
		// send special msg for register self's datatunnel connection to controller
		// 通过发送一个传递给保留地址的数据消息，服务端特殊处理，注册进服务端Server
		srcAddr := fmt.Sprintf("%s:%d", util.Uint32ToIP(a.privateIP).String(), 0)
		req := protocol.DataMessage{
			RequestID: uuid.New(),
			ClientID:  a.clientID,
			TunnelID:  0,
			// useless ip, 只是为了让agent注册自己的数据面连接。本质上DstAddr永远不会在控制器的节点注册管理器中找到，
			// 就会触发旁路逻辑，对应的旁路逻辑就是注册数据本身
			DstAddr: "0.0.0.1:65535",
			SrcAddr: srcAddr,
			Data:    []byte{},
		}

		buf, err := req.Encode()
		if err != nil {
			return err
		}
		if _, err := a.remote.Write(buf); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return finished, nil
}

func (a *Agent) registerAgent() (chan struct{}, error) {
	finished, err := util.GetTaskExecutor().RunTask(context.Background(), "registerAgent", func() error {
		err := RegisterSelf(a.control, &protocol.RegisterReqMessage{
			RequestID: uuid.New(),
			AgentName: a.name,
		})
		for !a.register.Get() {
			time.Sleep(time.Microsecond * 501)
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return finished, nil
}

// todo raw net.Conn should be protect. after use conn to send initial msg , agent can't write any thing into conn
func runDataLoop(f tcp.HandleFunc, conn net.Conn) {
	defer conn.Close()
	for {
		_, data, err := protocol.UnPackMessage(conn)
		if err != nil {
			// add warning
			continue
		}
		buf, io, err := f(context.Background(), data, conn)
		if err != nil {
			// add warning
			fmt.Println("err when handle data msg", err.Error())
			continue
		}

		if io != nil && buf != nil {
			_, err := io.Write(buf)
			if err != nil {
				// add warning
				continue
			}
		}
	}
}

// todo raw net.Conn should be protect. after use conn to send initial msg , agent can't write any thing into conn
func (a *Agent) runControlLoop(conn net.Conn) {
	for {
		msgType, data, err := protocol.UnPackMessage(conn)
		if err != nil {
			if err == io.EOF {
				// todo for other reconnect method..
				fmt.Println("peer has close, so receive io.EOF")
				return
			} else {
				// todo should there directly return ?
				fmt.Println("unpackerror err", err.Error())
				return
			}
		}
		switch msgType {
		case new(protocol.HeartbeatMessage).MsgType():
			// todo! maybe agent also need heart message
		case new(protocol.RegisterRespMessage).MsgType():
			resp := &protocol.RegisterRespMessage{}
			if err := resp.Decode(data); err != nil {
				fmt.Println("decode err", err.Error())
				return
			}
			if resp.Success {
				a.privateIP = resp.LanAddress
				a.register.Set(true)
			} else {
				fmt.Println(resp.Msg)
			}
		}
	}
}
