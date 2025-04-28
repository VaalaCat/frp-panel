package shared

import (
	"context"

	bizserver "github.com/VaalaCat/frp-panel/biz/server"
	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/rpcclient"
	"github.com/VaalaCat/frp-panel/services/tunnel"
	"github.com/VaalaCat/frp-panel/services/watcher"
	"github.com/VaalaCat/frp-panel/utils/logger"
	"github.com/sourcegraph/conc"
	"go.uber.org/fx"
)

type runServerParam struct {
	fx.In

	Lc fx.Lifecycle

	AppInstance      app.Application
	AppCtx           *app.Context
	ServerApiService app.Service    `name:"serverApiService"`
	TaskManager      watcher.Client `name:"serverTaskManager"`
	Cfg              conf.Config
}

func runServer(param runServerParam) {
	var (
		c            = context.Background()
		clientID     = param.AppInstance.GetConfig().Client.ID
		clientSecret = param.AppInstance.GetConfig().Client.Secret
		appInstance  = param.AppInstance
		ctx          = param.AppCtx
	)

	logger.Logger(c).Infof("start to init server")

	if len(clientID) == 0 {
		logger.Logger(ctx).Fatal("client id cannot be empty")
	}

	param.TaskManager.AddDurationTask(defs.PullConfigDuration, bizserver.PullConfig, appInstance, clientID, clientSecret)
	param.TaskManager.AddDurationTask(defs.PushProxyInfoDuration, bizserver.PushProxyInfo, appInstance, clientID, clientSecret)

	var wg conc.WaitGroup

	param.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Logger(ctx).Infof("start to run server, serverID: [%s]", clientID)
			appInstance.SetRPCCred(NewServerCred(appInstance))
			appInstance.SetMasterCli(NewServerMasterCli(appInstance))

			cliHandler := rpcclient.NewClientRPCHandler(
				appInstance,
				clientID,
				clientSecret,
				pb.Event_EVENT_REGISTER_SERVER,
				bizserver.HandleServerMessage,
			)

			appInstance.SetClientRPCHandler(cliHandler)
			appInstance.SetServerController(tunnel.NewServerController())

			go initServerOnce(appInstance, clientID, clientSecret)
			wg.Go(cliHandler.Run)
			wg.Go(param.TaskManager.Run)
			wg.Go(param.ServerApiService.Run)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			param.TaskManager.Stop()
			appInstance.GetClientRPCHandler().Stop()
			param.ServerApiService.Stop()
			wg.Wait()
			return nil
		},
	})

	logger.Logger(ctx).Infof("server started successfully, serverID: [%s]", clientID)
}

func initServerOnce(appInstance app.Application, clientID, clientSecret string) {
	err := bizserver.PullConfig(appInstance, clientID, clientSecret)
	if err != nil {
		logger.Logger(context.Background()).WithError(err).Errorf("cannot pull server config, wait for retry")
	}
}
