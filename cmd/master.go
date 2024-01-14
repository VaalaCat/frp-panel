package main

import (
	"embed"

	bizmaster "github.com/VaalaCat/frp-panel/biz/master"
	"github.com/VaalaCat/frp-panel/biz/master/auth"
	"github.com/VaalaCat/frp-panel/cache"
	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/services/api"
	"github.com/VaalaCat/frp-panel/services/master"
	"github.com/VaalaCat/frp-panel/services/server"
	"github.com/VaalaCat/frp-panel/utils"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/fatedier/golib/crypto"
	"github.com/glebarez/sqlite"
	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/conc"
	"gorm.io/gorm"
)

//go:embed all:out
var fs embed.FS

func runMaster() {
	crypto.DefaultSalt = conf.MasterDefaultSalt()
	master.MustInitMasterService()

	router := bizmaster.NewRouter(fs)
	api.MustInitApiService(conf.MasterAPIListenAddr(), router)

	initDatabase()
	cache.InitCache()
	auth.InitAuth()

	logrus.Infof("start to run master")
	m := master.GetMasterSerivce()
	opt := utils.NewBaseFRPServerUserAuthConfig(
		conf.Get().Master.InternalFRPServerPort,
		[]v1.HTTPPluginOptions{conf.FRPsAuthOption()},
	)

	s := server.GetServerSerivce(opt)
	a := api.GetAPIService()

	var wg conc.WaitGroup
	wg.Go(s.Run)
	wg.Go(m.Run)
	wg.Go(a.Run)
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
	default:
		logrus.Panicf("currently unsupported database type: %s", conf.Get().DB.Type)
	}

	models.GetDBManager().Init()
}
