package main

import (
	"context"
	"embed"
	"path/filepath"

	bizmaster "github.com/VaalaCat/frp-panel/biz/master"
	"github.com/VaalaCat/frp-panel/biz/master/auth"
	"github.com/VaalaCat/frp-panel/biz/master/proxy"
	"github.com/VaalaCat/frp-panel/biz/master/streamlog"
	bizserver "github.com/VaalaCat/frp-panel/biz/server"
	"github.com/VaalaCat/frp-panel/cache"
	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/dao"
	"github.com/VaalaCat/frp-panel/logger"
	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/rpc"
	"github.com/VaalaCat/frp-panel/services/master"
	"github.com/VaalaCat/frp-panel/services/mux"
	"github.com/VaalaCat/frp-panel/services/rpcclient"
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/VaalaCat/frp-panel/watcher"
	"github.com/fatedier/golib/crypto"
	"github.com/glebarez/sqlite"
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

func runMaster() {
	c := context.Background()
	crypto.DefaultSalt = conf.MasterDefaultSalt()

	streamlog.Init()

	initDatabase(c)
	cache.InitCache()
	auth.InitAuth()
	creds := dao.InitCert(conf.GetCertTemplate())

	router := bizmaster.NewRouter(fs)

	lisOpt := conf.GetListener(c)

	masterService := master.NewMasterService(credentials.NewTLS(creds))
	server := masterService.GetServer()
	muxServer := mux.NewMux(server, router, lisOpt.MuxLis, creds)

	masterH2CService := master.NewMasterService(insecure.NewCredentials())
	serverH2c := masterH2CService.GetServer()
	httpMuxServer := mux.NewMux(serverH2c, router, lisOpt.ApiLis, nil)

	tasks := watcher.NewClient()
	tasks.AddCronTask("0 0 3 * * *", proxy.CollectDailyStats)
	defer tasks.Stop()

	var wg conc.WaitGroup

	logger.Logger(c).Infof("start to run master")

	wg.Go(runDefaultInternalServer)
	wg.Go(muxServer.Run)
	wg.Go(httpMuxServer.Run)
	wg.Go(tasks.Run)

	wg.Wait()
}

func initDatabase(c context.Context) {
	logger.Logger(c).Infof("start to init database, type: %s", conf.Get().DB.Type)
	models.MustInitDBManager(nil, conf.Get().DB.Type)

	if conf.Get().IsDebug {
		models.GetDBManager().SetDebug(true)
	}

	switch conf.Get().DB.Type {
	case "sqlite3":
		if err := utils.EnsureDirectoryExists(conf.Get().DB.DSN); err != nil {
			logrus.WithError(err).Warnf("ensure directory failed, data location: [%s], keep data in current directory", conf.Get().DB.DSN)
			conf.Get().DB.DSN = filepath.Base(conf.Get().DB.DSN)
			logrus.Infof("new data location: [%s]", conf.Get().DB.DSN)
		}

		if sqlitedb, err := gorm.Open(sqlite.Open(conf.Get().DB.DSN), &gorm.Config{}); err != nil {
			logrus.Panic(err)
		} else {
			models.GetDBManager().SetDB("sqlite3", sqlitedb)
			logger.Logger(c).Infof("init database success, data location: [%s]", conf.Get().DB.DSN)
		}
	case "mysql":
		if mysqlDB, err := gorm.Open(mysql.Open(conf.Get().DB.DSN), &gorm.Config{}); err != nil {
			logrus.Panic(err)
		} else {
			models.GetDBManager().SetDB("mysql", mysqlDB)
			logger.Logger(c).Infof("init database success, data type: [%s]", "mysql")
		}
	case "postgres":
		if postgresDB, err := gorm.Open(postgres.Open(conf.Get().DB.DSN), &gorm.Config{}); err != nil {
			logrus.Panic(err)
		} else {
			models.GetDBManager().SetDB("postgres", postgresDB)
			logger.Logger(c).Infof("init database success, data type: [%s]", "postgres")
		}
	default:
		logrus.Panicf("currently unsupported database type: %s", conf.Get().DB.Type)
	}

	models.GetDBManager().Init()
}

func runDefaultInternalServer() {
	dao.InitDefaultServer(conf.Get().Master.APIHost)
	defaultServer, err := dao.GetDefaultServer()
	if err != nil {
		logrus.Fatal(err)
	}

	cred, err := utils.TLSClientCertNoValidate(rpc.GetClientCert(
		defaultServer.ServerID, defaultServer.ConnectSecret, pb.ClientType_CLIENT_TYPE_FRPS))
	if err != nil {
		logrus.Fatal(err)
	}
	conf.ClientCred = cred
	rpcclient.MustInitClientRPCSerivce(
		defaultServer.ServerID, defaultServer.ConnectSecret,
		pb.Event_EVENT_REGISTER_SERVER,
		bizserver.HandleServerMessage,
	)

	conf.Get().Client.ID = defaultServer.ServerID
	conf.Get().Client.Secret = defaultServer.ConnectSecret

	r := rpcclient.GetClientRPCSerivce()

	w := watcher.NewClient()
	w.AddDurationTask(common.PullConfigDuration, bizserver.PullConfig, defaultServer.ServerID, defaultServer.ConnectSecret)
	w.AddDurationTask(common.PushProxyInfoDuration, bizserver.PushProxyInfo, defaultServer.ServerID, defaultServer.ConnectSecret)

	go initServerOnce(defaultServer.ServerID, defaultServer.ConnectSecret)
	var wg conc.WaitGroup

	defer w.Stop()
	defer r.Stop()

	wg.Go(w.Run)
	wg.Go(r.Run)

	wg.Wait()
}
