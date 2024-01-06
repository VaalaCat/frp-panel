package api

import (
	"github.com/gin-gonic/gin"
)

type ApiService interface {
	Run()
	Stop()
}

type server struct {
	srv  *gin.Engine
	addr string
}

var (
	_          ApiService = (*server)(nil)
	apiService *server
)

func NewApiService(listenAddr string, router *gin.Engine) *server {
	return &server{
		srv:  router,
		addr: listenAddr,
	}
}

func MustInitApiService(listenAddr string, router *gin.Engine) {
	apiService = NewApiService(listenAddr, router)
}

func GetAPIService() ApiService {
	return apiService
}

func (s *server) Run() {
	s.srv.Run(s.addr)
}

func (s *server) Stop() {
}
