package api

import (
	"net"

	"github.com/gin-gonic/gin"
)

type ApiService interface {
	Run()
	Stop()
}

type server struct {
	srv    *gin.Engine
	addr   net.Listener
	enable bool
}

var (
	_ ApiService = (*server)(nil)
)

func NewApiService(listen net.Listener, router *gin.Engine, enable bool) *server {
	return &server{
		srv:    router,
		addr:   listen,
		enable: enable,
	}
}

func (s *server) Run() {
	// 如果完全使用mux，可以不启动
	if !s.enable {
		return
	}
	s.srv.RunListener(s.addr)
}

func (s *server) Stop() {
}
