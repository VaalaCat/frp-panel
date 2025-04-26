package shell

import (
	"fmt"
	"io"

	"github.com/VaalaCat/frp-panel/biz/master/client"
	"github.com/VaalaCat/frp-panel/biz/master/server"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/utils/logger"
)

func PTYConnect(ctx *app.Context, sender pb.Master_PTYConnectServer) error {
	msg, err := sender.Recv()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}

	clientID := ""

	if msg.GetClientBase() != nil {
		_, err = client.ValidateClientRequest(ctx, msg.GetClientBase())
		clientID = msg.GetClientBase().GetClientId()
	}
	if msg.GetServerBase() != nil {
		_, err = server.ValidateServerRequest(ctx, msg.GetServerBase())
		clientID = msg.GetServerBase().GetServerId()
	}
	if err != nil {
		return err
	}

	if len(clientID) == 0 {
		return fmt.Errorf("invalid client connect")
	}

	logger.Logger(sender.Context()).Infof("start pty connect, client id: [%s], session id: [%s]", clientID, msg.GetSessionId())

	ctx.GetApp().GetShellPTYMgr().Add(msg.GetSessionId(), sender)

	if err := sender.Send(&pb.PTYServerMessage{Data: []byte("ok")}); err != nil {
		return err
	}

	ctx.GetApp().GetShellPTYMgr().IsSessionDone(msg.GetSessionId())
	return nil
}
