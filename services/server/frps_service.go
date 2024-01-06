package server

import (
	"context"

	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/fatedier/frp/pkg/config/v1/validation"
	"github.com/fatedier/frp/pkg/util/log"
	"github.com/fatedier/frp/server"
	"github.com/sirupsen/logrus"
)

type ServerHandler interface {
	Run()
	Stop()
	GetCommonCfg() *v1.ServerConfig
}

type Server struct {
	srv    *server.Service
	Common *v1.ServerConfig
}

var (
	srv *Server
)

func InitGlobalServerService(svrCfg *v1.ServerConfig) {
	if srv != nil {
		logrus.Warn("server has been initialized")
		return
	}

	svrCfg.Complete()
	srv = NewServerHandler(svrCfg)
}

func GetGlobalServerSerivce() ServerHandler {
	if srv == nil {
		logrus.Panic("server has not been initialized")
	}
	return srv
}

func GetServerSerivce(svrCfg *v1.ServerConfig) ServerHandler {
	svrCfg.Complete()
	return NewServerHandler(svrCfg)
}

func NewServerHandler(svrCfg *v1.ServerConfig) *Server {
	warning, err := validation.ValidateServerConfig(svrCfg)
	if warning != nil {
		logrus.WithError(err).Warnf("validate server config warning: %+v", warning)
	}
	if err != nil {
		logrus.Panic(err)
	}

	log.InitLog(svrCfg.Log.To, svrCfg.Log.Level, svrCfg.Log.MaxDays, svrCfg.Log.DisablePrintColor)

	svr, err := server.NewService(svrCfg)
	if err != nil {
		logrus.Panic(err)
	}

	return &Server{
		srv:    svr,
		Common: svrCfg,
	}
}

func (s *Server) Run() {
	s.srv.Run(context.Background())
}

func (s *Server) Stop() {
	s.srv.Close()
}

func (s *Server) GetCommonCfg() *v1.ServerConfig {
	return s.Common
}
