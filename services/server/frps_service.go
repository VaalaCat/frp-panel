package server

import (
	"context"
	"sync"

	"github.com/VaalaCat/frp-panel/app"
	"github.com/VaalaCat/frp-panel/logger"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/fatedier/frp/pkg/config/v1/validation"
	"github.com/fatedier/frp/pkg/metrics/mem"
	"github.com/fatedier/frp/pkg/util/log"
	"github.com/fatedier/frp/server"
	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/conc"
)

type serverImpl struct {
	srv       *server.Service
	Common    *v1.ServerConfig
	firstSync sync.Once
}

func NewServerHandler(svrCfg *v1.ServerConfig) app.ServerHandler {
	svrCfg.Complete()

	warning, err := validation.ValidateServerConfig(svrCfg)
	if warning != nil {
		logger.Logger(context.Background()).WithError(err).Warnf("validate server config warning: %+v", warning)
	}
	if err != nil {
		logrus.Panic(err)
	}

	log.InitLogger(svrCfg.Log.To, svrCfg.Log.Level, int(svrCfg.Log.MaxDays), svrCfg.Log.DisablePrintColor)

	var svr *server.Service

	if svr, err = server.NewService(svrCfg); err != nil {
		logger.Logger(context.Background()).WithError(err).Panic("cannot create server, exit and restart")
	}

	return &serverImpl{
		srv:       svr,
		Common:    svrCfg,
		firstSync: sync.Once{},
	}
}

func (s *serverImpl) Run() {
	wg := conc.NewWaitGroup()
	wg.Go(func() { s.srv.Run(context.Background()) })
	wg.Wait()
}

func (s *serverImpl) Stop() {
	c := context.Background()
	wg := conc.NewWaitGroup()
	wg.Go(func() {
		err := s.srv.Close()
		if err != nil {
			logger.Logger(c).Errorf("close server error: %v", err)
		}
		logger.Logger(c).Infof("server closed")
	})
	wg.Wait()
}

func (s *serverImpl) GetCommonCfg() *v1.ServerConfig {
	return s.Common
}

func (s *serverImpl) GetMem() *mem.ServerStats {
	return mem.StatsCollector.GetServer()
}

func (s *serverImpl) GetProxyStatsByType(proxyType v1.ProxyType) []*mem.ProxyStats {
	return mem.StatsCollector.GetProxiesByType(string(proxyType))
}

func (s *serverImpl) IsFirstSync() bool {
	result := false
	s.firstSync.Do(func() {
		result = true
	})
	return result
}
