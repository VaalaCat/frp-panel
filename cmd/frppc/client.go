package main

import (
	"context"

	bizclient "github.com/VaalaCat/frp-panel/biz/client"
	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/rpc"
	"github.com/VaalaCat/frp-panel/services/rpcclient"
	"github.com/VaalaCat/frp-panel/services/tunnel"
	"github.com/VaalaCat/frp-panel/services/watcher"
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/VaalaCat/frp-panel/utils/logger"
	"github.com/sourcegraph/conc"
)

func runClient(appInstance app.Application) {
	var (
		c            = context.Background()
		clientID     = appInstance.GetConfig().Client.ID
		clientSecret = appInstance.GetConfig().Client.Secret
		ctx          = context.Background()
	)
	logger.Logger(c).Infof("start to run client")
	if len(clientSecret) == 0 {
		logger.Logger(ctx).Fatal("client secret cannot be empty")
	}

	if len(clientID) == 0 {
		logger.Logger(ctx).Fatal("client id cannot be empty")
	}

	cred, err := utils.TLSClientCertNoValidate(rpc.GetClientCert(appInstance, clientID, clientSecret, pb.ClientType_CLIENT_TYPE_FRPC))
	if err != nil {
		logger.Logger(ctx).Fatal(err)
	}

	appInstance.SetRPCCred(cred)
	appInstance.SetMasterCli(rpc.NewMasterCli(appInstance))
	appInstance.SetClientController(tunnel.NewClientController())

	r := rpcclient.NewClientRPCHandler(
		appInstance,
		clientID,
		clientSecret,
		pb.Event_EVENT_REGISTER_CLIENT,
		bizclient.HandleServerMessage,
	)
	appInstance.SetClientRPCHandler(r)

	w := watcher.NewClient()
	w.AddDurationTask(defs.PullConfigDuration, bizclient.PullConfig, clientID, clientSecret)

	initClientOnce(appInstance, clientID, clientSecret)

	defer w.Stop()
	defer r.Stop()

	var wg conc.WaitGroup
	wg.Go(r.Run)
	wg.Go(w.Run)
	wg.Wait()
}

func initClientOnce(appInstance app.Application, clientID, clientSecret string) {
	err := bizclient.PullConfig(appInstance, clientID, clientSecret)
	if err != nil {
		logger.Logger(context.Background()).WithError(err).Errorf("cannot pull client config, wait for retry")
	}
}
