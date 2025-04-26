package main

import (
	"context"
	"embed"
	"net/http"
	"path/filepath"

	"github.com/VaalaCat/frp-panel/app"
	bizmaster "github.com/VaalaCat/frp-panel/biz/master"
	"github.com/VaalaCat/frp-panel/biz/master/auth"
	"github.com/VaalaCat/frp-panel/biz/master/proxy"
	"github.com/VaalaCat/frp-panel/biz/master/streamlog"
	bizserver "github.com/VaalaCat/frp-panel/biz/server"
	"github.com/VaalaCat/frp-panel/cache"
	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/dao"
	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/logger"
	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/rpc"
	"github.com/VaalaCat/frp-panel/services/master"
	"github.com/VaalaCat/frp-panel/services/mux"
	"github.com/VaalaCat/frp-panel/services/rpcclient"
	"github.com/VaalaCat/frp-panel/tunnel"
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/VaalaCat/frp-panel/utils/wsgrpc"
	"github.com/VaalaCat/frp-panel/watcher"
	"github.com/fatedier/golib/crypto"
	"github.com/glebarez/sqlite"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/conc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//go:embed all:out
var fs embed.FS

func runMaster(appInstance app.Application) {
	c := app.NewContext(context.Background(), appInstance)
	cfg := appInstance.GetConfig()
	crypto.DefaultSalt = conf.MasterDefaultSalt(cfg)

	appInstance.SetClientLogManager(streamlog.NewClientLogManager())

	initDatabase(c, appInstance)
	cache.InitCache(cfg)
	auth.InitAuth(appInstance)
	creds := dao.NewQuery(c).InitCert(conf.GetCertTemplate(cfg))

	router := bizmaster.NewRouter(fs, appInstance)

	lisOpt := conf.GetListener(c, cfg)

	// ---- ws grpc start -----
	wsListener := wsgrpc.NewWSListener("ws-listener", "wsgrpc", 100)

	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	router.GET("/wsgrpc", wsgrpc.GinWSHandler(wsListener, upgrader))
	// ---- ws grpc end -----

	masterService := master.NewMasterService(appInstance, credentials.NewTLS(creds))
	server := masterService.GetServer()
	muxServer := mux.NewMux(server, router, lisOpt.MuxLis, creds)

	masterH2CService := master.NewMasterService(appInstance, insecure.NewCredentials())
	serverH2c := masterH2CService.GetServer()
	httpMuxServer := mux.NewMux(serverH2c, router, lisOpt.ApiLis, nil)

	tasks := watcher.NewClient()
	tasks.AddCronTask("0 0 3 * * *", proxy.CollectDailyStats, appInstance)
	defer tasks.Stop()

	var wg conc.WaitGroup

	logger.Logger(c).Infof("start to run master")

	go func() {
		wsGrpcServer := master.NewMasterService(appInstance, insecure.NewCredentials()).GetServer()
		if err := wsGrpcServer.Serve(wsListener); err != nil {
			logrus.Fatalf("gRPC server error: %v", err)
		}
	}()
	wg.Go(func() { runDefaultInternalServer(appInstance) })
	wg.Go(muxServer.Run)
	wg.Go(httpMuxServer.Run)
	wg.Go(tasks.Run)

	wg.Wait()
}

func initDatabase(c context.Context, appInstance app.Application) {
	logger.Logger(c).Infof("start to init database, type: %s", appInstance.GetConfig().DB.Type)
	mgr := models.NewDBManager(nil, appInstance.GetConfig().DB.Type)
	appInstance.SetDBManager(mgr)

	if appInstance.GetConfig().IsDebug {
		appInstance.GetDBManager().SetDebug(true)
	}

	switch appInstance.GetConfig().DB.Type {
	case "sqlite3":
		if err := utils.EnsureDirectoryExists(appInstance.GetConfig().DB.DSN); err != nil {
			logrus.WithError(err).Warnf("ensure directory failed, data location: [%s], keep data in current directory",
				appInstance.GetConfig().DB.DSN)
			tmpCfg := appInstance.GetConfig()
			tmpCfg.DB.DSN = filepath.Base(appInstance.GetConfig().DB.DSN)
			appInstance.SetConfig(tmpCfg)
			logrus.Infof("new data location: [%s]", appInstance.GetConfig().DB.DSN)
		}

		if sqlitedb, err := gorm.Open(sqlite.Open(appInstance.GetConfig().DB.DSN), &gorm.Config{}); err != nil {
			logrus.Panic(err)
		} else {
			appInstance.GetDBManager().SetDB("sqlite3", sqlitedb)
			logger.Logger(c).Infof("init database success, data location: [%s]", appInstance.GetConfig().DB.DSN)
		}
	case "mysql":
		if mysqlDB, err := gorm.Open(mysql.Open(appInstance.GetConfig().DB.DSN), &gorm.Config{}); err != nil {
			logrus.Panic(err)
		} else {
			appInstance.GetDBManager().SetDB("mysql", mysqlDB)
			logger.Logger(c).Infof("init database success, data type: [%s]", "mysql")
		}
	case "postgres":
		if postgresDB, err := gorm.Open(postgres.Open(appInstance.GetConfig().DB.DSN), &gorm.Config{}); err != nil {
			logrus.Panic(err)
		} else {
			appInstance.GetDBManager().SetDB("postgres", postgresDB)
			logger.Logger(c).Infof("init database success, data type: [%s]", "postgres")
		}
	default:
		logrus.Panicf("currently unsupported database type: %s", appInstance.GetConfig().DB.Type)
	}

	appInstance.GetDBManager().Init()
}

func runDefaultInternalServer(appInstance app.Application) {
	logger.Logger(context.Background()).Infof("init default internal server")

	appCtx := app.NewContext(clientCmd.Context(), appInstance)

	dao.NewQuery(appCtx).InitDefaultServer(appInstance.GetConfig().Master.APIHost)
	defaultServer, err := dao.NewQuery(appCtx).GetDefaultServer()
	if err != nil {
		logrus.Fatal(err)
	}

	cred, err := utils.TLSClientCertNoValidate(rpc.GetClientCert(
		appInstance,
		defaultServer.ServerID, defaultServer.ConnectSecret, pb.ClientType_CLIENT_TYPE_FRPS))
	if err != nil {
		logrus.Fatal(err)
	}
	appInstance.SetClientCred(cred)

	appInstance.SetMasterCli(rpc.NewMasterCli(appInstance))
	appInstance.SetServerController(tunnel.NewServerController())

	cliHandler := rpcclient.NewClientRPCHandler(
		appInstance,
		defaultServer.ServerID, defaultServer.ConnectSecret,
		pb.Event_EVENT_REGISTER_SERVER,
		bizserver.HandleServerMessage,
	)

	appInstance.SetClientRPCHandler(cliHandler)

	tmpCfg := appInstance.GetConfig()
	tmpCfg.Client.ID = defaultServer.ServerID
	tmpCfg.Client.Secret = defaultServer.ConnectSecret
	appInstance.SetConfig(tmpCfg)

	r := cliHandler

	w := watcher.NewClient()
	w.AddDurationTask(defs.PullConfigDuration, bizserver.PullConfig, appInstance, defaultServer.ServerID, defaultServer.ConnectSecret)
	w.AddDurationTask(defs.PushProxyInfoDuration, bizserver.PushProxyInfo, appInstance, defaultServer.ServerID, defaultServer.ConnectSecret)

	go initServerOnce(appInstance, defaultServer.ServerID, defaultServer.ConnectSecret)
	var wg conc.WaitGroup

	defer w.Stop()
	defer r.Stop()

	wg.Go(w.Run)
	wg.Go(r.Run)

	wg.Wait()
}
