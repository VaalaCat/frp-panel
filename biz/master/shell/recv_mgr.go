package shell

import (
	"fmt"
	"io"

	"github.com/VaalaCat/frp-panel/biz/master/client"
	"github.com/VaalaCat/frp-panel/biz/master/server"
	"github.com/VaalaCat/frp-panel/logger"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/samber/lo"
)

func PTYConnect(sender pb.Master_PTYConnectServer) error {
	msg, err := sender.Recv()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}

	clientID := ""

	if msg.GetClientBase() != nil {
		_, err = client.ValidateClientRequest(msg.GetClientBase())
		clientID = msg.GetClientBase().GetClientId()
	}
	if msg.GetServerBase() != nil {
		_, err = server.ValidateServerRequest(msg.GetServerBase())
		clientID = msg.GetServerBase().GetServerId()
	}
	if err != nil {
		return err
	}

	if len(clientID) == 0 {
		return fmt.Errorf("invalid client connect")
	}

	logger.Logger(sender.Context()).Infof("start pty connect, client id: [%s], session id: [%s]", clientID, msg.GetSessionId())

	Mgr().Add(msg.GetSessionId(), sender)

	if err := sender.Send(&pb.PTYServerMessage{Data: lo.ToPtr("ok")}); err != nil {
		return err
	}

	Mgr().IsSessionDone(msg.GetSessionId())
	return nil
}
