package main

import (
	"context"
	"net"

	"github.com/VaalaCat/frp-panel/app"
	bizserver "github.com/VaalaCat/frp-panel/biz/server"
	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/logger"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/rpc"
	"github.com/VaalaCat/frp-panel/services/api"
	"github.com/VaalaCat/frp-panel/services/rpcclient"
	"github.com/VaalaCat/frp-panel/tunnel"
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/VaalaCat/frp-panel/watcher"
	"github.com/fatedier/golib/crypto"
	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/conc"
)

func runServer(appInstance app.Application) {
	var (
		c            = context.Background()
		clientID     = appInstance.GetConfig().Client.ID
		clientSecret = appInstance.GetConfig().Client.Secret
		cfg          = appInstance.GetConfig()
	)
	crypto.DefaultSalt = cfg.App.Secret
	logger.Logger(c).Infof("start to run server")

	if len(clientID) == 0 {
		logrus.Fatal("client id cannot be empty")
	}

	l, err := net.Listen("tcp", conf.ServerAPIListenAddr(cfg))
	if err != nil {
		logger.Logger(c).WithError(err).Fatalf("failed to listen addr: %v", conf.ServerAPIListenAddr(cfg))
		return
	}

	a := api.NewApiService(l, bizserver.NewRouter(appInstance), true)
	defer a.Stop()

	cred, err := utils.TLSClientCertNoValidate(rpc.GetClientCert(appInstance, clientID, clientSecret, pb.ClientType_CLIENT_TYPE_FRPS))
	if err != nil {
		logrus.Fatal(err)
	}
	appInstance.SetClientCred(cred)

	appInstance.SetMasterCli(rpc.NewMasterCli(appInstance))
	appInstance.SetServerController(tunnel.NewServerController())

	cliHandler := rpcclient.NewClientRPCHandler(
		appInstance,
		clientID,
		clientSecret,
		pb.Event_EVENT_REGISTER_SERVER,
		bizserver.HandleServerMessage,
	)

	appInstance.SetClientRPCHandler(cliHandler)

	r := cliHandler
	defer r.Stop()

	w := watcher.NewClient()
	w.AddDurationTask(defs.PullConfigDuration, bizserver.PullConfig, clientID, clientSecret)
	w.AddDurationTask(defs.PushProxyInfoDuration, bizserver.PushProxyInfo, clientID, clientSecret)
	defer w.Stop()

	initServerOnce(appInstance, clientID, clientSecret)

	var wg conc.WaitGroup
	wg.Go(r.Run)
	wg.Go(w.Run)
	wg.Go(a.Run)
	wg.Wait()
}

func initServerOnce(appInstance app.Application, clientID, clientSecret string) {
	err := bizserver.PullConfig(appInstance, clientID, clientSecret)
	if err != nil {
		logger.Logger(context.Background()).WithError(err).Errorf("cannot pull server config, wait for retry")
	}
}
