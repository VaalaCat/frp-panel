package master

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type MasterHandler interface {
	Run()
	Stop()
}

type Master struct {
	grpcServer *grpc.Server
}

var (
	cli *Master
)

func MustInitMasterService() {
	if cli != nil {
		logrus.Warn("server has been initialized")
		return
	}
	cli = NewMasterHandler()
}

func GetMasterSerivce() MasterHandler {
	if cli == nil {
		logrus.Panic("server has not been initialized")
	}
	return cli
}

func NewMasterHandler() *Master {
	s := NewRpcServer()
	return &Master{
		grpcServer: s,
	}
}

func (s *Master) Run() {
	RunRpcServer(s.grpcServer)
}

func (s *Master) Stop() {
	s.grpcServer.Stop()
}
