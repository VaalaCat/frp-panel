package master

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

func MustInitMasterService(creds credentials.TransportCredentials) {
	if cli != nil {
		logrus.Warn("server has been initialized")
		return
	}
	cli = NewMasterHandler(creds)
}

func GetMasterSerivce() MasterHandler {
	if cli == nil {
		logrus.Panic("server has not been initialized")
	}
	return cli
}

func NewMasterHandler(creds credentials.TransportCredentials) *Master {
	s := NewRpcServer(creds)
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
