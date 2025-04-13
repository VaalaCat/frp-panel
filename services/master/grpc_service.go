package master

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type MasterService interface {
	Run()
	Stop()
	GetServer() *grpc.Server
}

type master struct {
	grpcServer *grpc.Server
}

func NewMasterService(creds credentials.TransportCredentials) MasterService {
	s := newRpcServer(creds)
	return &master{
		grpcServer: s,
	}
}

func (s *master) Run() {
	runRpcServer(s.grpcServer)
}

func (s *master) Stop() {
	s.grpcServer.Stop()
}

func (s *master) GetServer() *grpc.Server {
	return s.grpcServer
}
