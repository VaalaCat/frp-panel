package rpcclient

import (
	"github.com/VaalaCat/frp-panel/app"
	"github.com/VaalaCat/frp-panel/pb"
)

type ClientRPCHandler interface {
	Run()
	Stop()
	GetCli() pb.MasterClient
}

type clientRPC struct {
	appInstance  app.Application
	rpcClient    pb.MasterClient
	done         chan bool
	handerFunc   func(appInstance app.Application, req *pb.ServerMessage) *pb.ClientMessage
	clientID     string
	clientSecret string
	event        pb.Event
}

func NewClientRPCHandler(
	appInstance app.Application,
	clientID,
	clientSecret string,
	event pb.Event,
	handerFunc func(appInstance app.Application, req *pb.ServerMessage) *pb.ClientMessage,
) app.ClientRPCHandler {
	rpcCli := appInstance.GetMasterCli()
	done := make(chan bool)
	return &clientRPC{
		appInstance:  appInstance,
		rpcClient:    rpcCli,
		done:         done,
		handerFunc:   handerFunc,
		clientID:     clientID,
		clientSecret: clientSecret,
		event:        event,
	}
}

func (s *clientRPC) Run() {
	StartRPCClient(s.appInstance, s.rpcClient, s.done, s.clientID, s.clientSecret, s.event, s.handerFunc)
}

func (s *clientRPC) Stop() {
	close(s.done)
}

func (s *clientRPC) GetCli() pb.MasterClient {
	return s.rpcClient
}
