package main

import (
	"context"
	"net"

	bizserver "github.com/VaalaCat/frp-panel/biz/server"
	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/logger"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/rpc"
	"github.com/VaalaCat/frp-panel/services/api"
	"github.com/VaalaCat/frp-panel/services/rpcclient"
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/VaalaCat/frp-panel/watcher"
	"github.com/fatedier/golib/crypto"
	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/conc"
)

func runServer() {
	var (
		c            = context.Background()
		clientID     = conf.Get().Client.ID
		clientSecret = conf.Get().Client.Secret
	)
	crypto.DefaultSalt = conf.Get().App.Secret
	logger.Logger(c).Infof("start to run server")

	if len(clientID) == 0 {
		logrus.Fatal("client id cannot be empty")
	}

	l, err := net.Listen("tcp", conf.ServerAPIListenAddr())
	if err != nil {
		logger.Logger(c).WithError(err).Fatalf("failed to listen addr: %v", conf.ServerAPIListenAddr())
		return
	}

	a := api.NewApiService(l, bizserver.NewRouter(), true)
	defer a.Stop()

	cred, err := utils.TLSClientCertNoValidate(rpc.GetClientCert(clientID, clientSecret, pb.ClientType_CLIENT_TYPE_FRPS))
	if err != nil {
		logrus.Fatal(err)
	}
	conf.ClientCred = cred
	rpcclient.MustInitClientRPCSerivce(
		clientID,
		clientSecret,
		pb.Event_EVENT_REGISTER_SERVER,
		bizserver.HandleServerMessage,
	)

	r := rpcclient.GetClientRPCSerivce()
	defer r.Stop()

	w := watcher.NewClient()
	w.AddDurationTask(common.PullConfigDuration, bizserver.PullConfig, clientID, clientSecret)
	w.AddDurationTask(common.PushProxyInfoDuration, bizserver.PushProxyInfo, clientID, clientSecret)
	defer w.Stop()

	initServerOnce(clientID, clientSecret)

	var wg conc.WaitGroup
	wg.Go(r.Run)
	wg.Go(w.Run)
	wg.Go(a.Run)
	wg.Wait()
}

func initServerOnce(clientID, clientSecret string) {
	err := bizserver.PullConfig(clientID, clientSecret)
	if err != nil {
		logger.Logger(context.Background()).WithError(err).Errorf("cannot pull server config, wait for retry")
	}
}
