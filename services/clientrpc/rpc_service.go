package clientrpc

import (
	"context"

	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/utils/logger"
)

type ClientRPCHandler interface {
	Run()
	Stop()
	GetCli() pb.MasterClient
}

type clientRPCHandler struct {
	appInstance  app.Application
	rpcClient    app.MasterClient
	ctx          context.Context
	cancel       context.CancelFunc
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
	ctx, cancel := context.WithCancel(context.Background())
	return &clientRPCHandler{
		appInstance:  appInstance,
		rpcClient:    rpcCli,
		ctx:          ctx,
		cancel:       cancel,
		handerFunc:   handerFunc,
		clientID:     clientID,
		clientSecret: clientSecret,
		event:        event,
	}
}

func (s *clientRPCHandler) Run() {
	defer func() {
		if err := recover(); err != nil {
			logger.Logger(context.Background()).Fatalf("client rpc handler panic: %v", err)
		}
	}()

	startClientRpcHandler(s.appInstance, s.rpcClient, s.ctx, s.clientID, s.clientSecret, s.event, s.handerFunc)
}

func (s *clientRPCHandler) Stop() {
	s.cancel()
}

func (s *clientRPCHandler) GetCli() app.MasterClient {
	return s.rpcClient
}
