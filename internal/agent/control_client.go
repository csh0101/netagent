package agent

import (
	"errors"
	"net"

	"github.com/csh0101/netagent.git/internal/protocol"
	"github.com/google/uuid"
)

func RegisterSelf(conn net.Conn, msg *protocol.RegisterReqMessage) error {
	if msg.RequestID.String() == "" {
		msg.RequestID = uuid.New()
	}
	if msg.AgentName == "" {
		return errors.New("AgentName is required")
	}

	if data, err := msg.Encode(); err != nil {
		return err
	} else {
		if _, err := conn.Write(data); err != nil {
			return err
		}
	}
	return nil
}

func SendHearbeat(conn net.Conn, msg *protocol.HeartbeatMessage) error {
	if msg.AgentName == "" || msg.LanAddress == 0 || msg.RequestID.String() == "" {
		return errors.New("err when send heartbeat params is required")
	}

	if data, err := msg.Encode(); err != nil {
		return err
	} else {
		if _, err := conn.Write(data); err != nil {
			return err
		}
	}
	return nil
}
