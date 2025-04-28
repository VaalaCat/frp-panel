package shared

import (
	"context"

	bizclient "github.com/VaalaCat/frp-panel/biz/client"
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

type runClientParam struct {
	fx.In

	Lc fx.Lifecycle

	Ctx         *app.Context
	AppInstance app.Application
	TaskManager watcher.Client `name:"clientTaskManager"`
	Cfg         conf.Config
}

func runClient(param runClientParam) {
	var (
		ctx          = param.Ctx
		clientID     = param.AppInstance.GetConfig().Client.ID
		clientSecret = param.AppInstance.GetConfig().Client.Secret
		appInstance  = param.AppInstance
	)
	logger.Logger(ctx).Infof("start to run client")
	if len(clientSecret) == 0 {
		logger.Logger(ctx).Fatal("client secret cannot be empty")
	}

	if len(clientID) == 0 {
		logger.Logger(ctx).Fatal("client id cannot be empty")
	}

	param.TaskManager.AddDurationTask(defs.PullConfigDuration,
		bizclient.PullConfig, appInstance, clientID, clientSecret)

	var wg conc.WaitGroup
	param.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			appInstance.SetRPCCred(NewClientCred(appInstance))
			appInstance.SetMasterCli(NewClientMasterCli(appInstance))
			appInstance.SetClientController(tunnel.NewClientController())

			cliRpcHandler := rpcclient.NewClientRPCHandler(
				appInstance,
				clientID,
				clientSecret,
				pb.Event_EVENT_REGISTER_CLIENT,
				bizclient.HandleServerMessage,
			)
			appInstance.SetClientRPCHandler(cliRpcHandler)

			initClientOnce(appInstance, clientID, clientSecret)

			wg.Go(cliRpcHandler.Run)
			wg.Go(param.TaskManager.Run)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			param.TaskManager.Stop()
			appInstance.GetClientRPCHandler().Stop()

			wg.Wait()
			return nil
		},
	})
}

func initClientOnce(appInstance app.Application, clientID, clientSecret string) {
	err := bizclient.PullConfig(appInstance, clientID, clientSecret)
	if err != nil {
		logger.Logger(context.Background()).WithError(err).Errorf("cannot pull client config, wait for retry")
	}
}
