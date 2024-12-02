package rpcclient

import (
	"context"

	"github.com/VaalaCat/frp-panel/logger"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/sirupsen/logrus"
)

type ClientRPCHandler interface {
	Run()
	Stop()
	GetCli() pb.MasterClient
}

type ClientRPC struct {
	rpcClient    pb.MasterClient
	done         chan bool
	handerFunc   func(req *pb.ServerMessage) *pb.ClientMessage
	clientID     string
	clientSecret string
	event        pb.Event
}

var (
	cliRpc *ClientRPC
)

func MustInitClientRPCSerivce(
	clientID,
	clientSecret string,
	event pb.Event,
	handerFunc func(req *pb.ServerMessage) *pb.ClientMessage,
) {
	ctx := context.Background()
	if cliRpc != nil {
		logger.Logger(ctx).Warn("rpc client has been initialized")
		return
	}
	cliRpc = NewClientRPCHandler(clientID, clientSecret, event, handerFunc)
}

func GetClientRPCSerivce() ClientRPCHandler {
	if cliRpc == nil {
		logrus.Panic("rpc client has not been initialized")
	}
	return cliRpc
}

func NewClientRPCHandler(
	clientID,
	clientSecret string,
	event pb.Event,
	handerFunc func(req *pb.ServerMessage) *pb.ClientMessage,
) *ClientRPC {
	rpcCli, err := NewMasterCli()
	if err != nil {
		logrus.Fatalf("new rpc client failed: %v", err)
	}
	done := make(chan bool)
	return &ClientRPC{
		rpcClient:    rpcCli,
		done:         done,
		handerFunc:   handerFunc,
		clientID:     clientID,
		clientSecret: clientSecret,
		event:        event,
	}
}

func (s *ClientRPC) Run() {
	StartRPCClient(s.rpcClient, s.done, s.clientID, s.clientSecret, s.event, s.handerFunc)
}

func (s *ClientRPC) Stop() {
	close(s.done)
}

func (s *ClientRPC) GetCli() pb.MasterClient {
	return s.rpcClient
}
