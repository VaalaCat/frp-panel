package server

import (
	"context"
	"reflect"

	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/rpc"
	"github.com/VaalaCat/frp-panel/services/server"
	"github.com/VaalaCat/frp-panel/tunnel"
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/sirupsen/logrus"
)

func PullConfig(serverID, serverSecret string) error {
	logrus.Infof("start to pull server config, serverID: [%s]", serverID)
	ctx := context.Background()
	cli, err := rpc.MasterCli(ctx)
	if err != nil {
		logrus.WithError(err).Error("cannot get master server")
		return err
	}
	resp, err := cli.PullServerConfig(ctx, &pb.PullServerConfigReq{
		Base: &pb.ServerBase{
			ServerId:     serverID,
			ServerSecret: serverSecret,
		},
	})
	if err != nil {
		logrus.WithError(err).Error("cannot pull server config")
		return err
	}

	if len(resp.GetServer().GetConfig()) == 0 {
		logrus.Infof("server [%s] config is empty, wait for server init", serverID)
		return nil
	}

	s, err := utils.LoadServerConfig([]byte(resp.GetServer().GetConfig()), true)
	if err != nil {
		logrus.WithError(err).Error("cannot load server config")
		return err
	}

	ctrl := tunnel.GetServerController()

	if t := ctrl.Get(serverID); t != nil {
		if !reflect.DeepEqual(t.GetCommonCfg(), s) {
			t.Stop()
			ctrl.Delete(serverID)
			logrus.Infof("server %s config changed, will recreate it", serverID)
		} else {
			logrus.Infof("server %s config not changed", serverID)
			return nil
		}
	}
	ctrl.Add(serverID, server.NewServerHandler(s))
	ctrl.Run(serverID)

	logrus.Infof("pull server config success, serverID: [%s]", serverID)
	return nil
}
