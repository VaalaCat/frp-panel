package server

import (
	"context"
	"reflect"

	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/server"
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/VaalaCat/frp-panel/utils/logger"
)

func PullConfig(appInstance app.Application, serverID, serverSecret string) error {
	ctx := context.Background()
	logger.Logger(ctx).Infof("start to pull server config, serverID: [%s]", serverID)

	cli := appInstance.GetMasterCli()
	resp, err := cli.Call().PullServerConfig(ctx, &pb.PullServerConfigReq{
		Base: &pb.ServerBase{
			ServerId:     serverID,
			ServerSecret: serverSecret,
		},
	})
	if err != nil {
		logger.Logger(context.Background()).WithError(err).Error("cannot pull server config")
		return err
	}

	if len(resp.GetServer().GetConfig()) == 0 {
		logger.Logger(ctx).Infof("server [%s] config is empty, wait for server init", serverID)
		return nil
	}

	s, err := utils.LoadServerConfig([]byte(resp.GetServer().GetConfig()), true)
	if err != nil {
		logger.Logger(context.Background()).WithError(err).Error("cannot load server config")
		return err
	}

	ctrl := appInstance.GetServerController()

	if t := ctrl.Get(serverID); t != nil {
		if !reflect.DeepEqual(t.GetCommonCfg(), s) {
			t.Stop()
			ctrl.Delete(serverID)
			logger.Logger(ctx).Infof("server %s config changed, will recreate it", serverID)
		} else {
			logger.Logger(ctx).Infof("server %s config not changed", serverID)
			return nil
		}
	}
	ctrl.Add(serverID, server.NewServerHandler(s))
	ctrl.Run(serverID)

	logger.Logger(ctx).Infof("pull server config success, serverID: [%s]", serverID)
	return nil
}
