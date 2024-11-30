package main

import (
	"embed"

	bizmaster "github.com/VaalaCat/frp-panel/biz/master"
	"github.com/VaalaCat/frp-panel/biz/master/auth"
	"github.com/VaalaCat/frp-panel/biz/master/proxy"
	bizserver "github.com/VaalaCat/frp-panel/biz/server"
	"github.com/VaalaCat/frp-panel/cache"
	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/dao"
	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/rpc"
	"github.com/VaalaCat/frp-panel/services/api"
	"github.com/VaalaCat/frp-panel/services/master"
	"github.com/VaalaCat/frp-panel/services/rpcclient"
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/VaalaCat/frp-panel/watcher"
	"github.com/fatedier/golib/crypto"
	"github.com/glebarez/sqlite"
	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/conc"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//go:embed all:out
var fs embed.FS

func runMaster() {
	crypto.DefaultSalt = conf.MasterDefaultSalt()

	initDatabase()
	cache.InitCache()
	auth.InitAuth()
	creds := dao.InitCert(conf.GetCertTemplate())

	master.MustInitMasterService(creds)
	router := bizmaster.NewRouter(fs)
	api.MustInitApiService(conf.MasterAPIListenAddr(), router)

	logrus.Infof("start to run master")
	m := master.GetMasterSerivce()
	a := api.GetAPIService()

	r, w := initDefaultInternalServer()
	defer w.Stop()
	defer r.Stop()

	tasks := watcher.NewClient()
	tasks.AddCronTask("0 0 3 * * *", proxy.CollectDailyStats)
	defer tasks.Stop()

	var wg conc.WaitGroup
	wg.Go(w.Run)
	wg.Go(r.Run)
	wg.Go(m.Run)
	wg.Go(a.Run)
	wg.Go(tasks.Run)
	wg.Wait()
}

func initDatabase() {
	logrus.Infof("start to init database, type: %s", conf.Get().DB.Type)
	models.MustInitDBManager(nil, conf.Get().DB.Type)

	switch conf.Get().DB.Type {
	case "sqlite3":
		if sqlitedb, err := gorm.Open(sqlite.Open(conf.Get().DB.DSN), &gorm.Config{}); err != nil {
			logrus.Panic(err)
		} else {
			models.GetDBManager().SetDB("sqlite3", sqlitedb)
			logrus.Infof("init database success, data location: [%s]", conf.Get().DB.DSN)
		}
	case "mysql":
		if mysqlDB, err := gorm.Open(mysql.Open(conf.Get().DB.DSN), &gorm.Config{}); err != nil {
			logrus.Panic(err)
		} else {
			models.GetDBManager().SetDB("mysql", mysqlDB)
			logrus.Infof("init database success, data type: [%s]", "mysql")
		}
	case "postgres":
		if postgresDB, err := gorm.Open(postgres.Open(conf.Get().DB.DSN), &gorm.Config{}); err != nil {
			logrus.Panic(err)
		} else {
			models.GetDBManager().SetDB("postgres", postgresDB)
			logrus.Infof("init database success, data type: [%s]", "postgres")
		}
	default:
		logrus.Panicf("currently unsupported database type: %s", conf.Get().DB.Type)
	}

	models.GetDBManager().Init()
}

func initDefaultInternalServer() (rpcclient.ClientRPCHandler, watcher.Client) {
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

	r := rpcclient.GetClientRPCSerivce()

	w := watcher.NewClient()
	w.AddDurationTask(common.PullConfigDuration, bizserver.PullConfig, defaultServer.ServerID, defaultServer.ConnectSecret)
	w.AddDurationTask(common.PushProxyInfoDuration, bizserver.PushProxyInfo, defaultServer.ServerID, defaultServer.ConnectSecret)

	go initServerOnce(defaultServer.ServerID, defaultServer.ConnectSecret)
	return r, w
}
