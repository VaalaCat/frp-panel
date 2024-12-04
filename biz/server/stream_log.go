package server

import (
	"context"

	"github.com/VaalaCat/frp-panel/biz/common"
	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/rpcclient"
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/sirupsen/logrus"
)

func StartSteamLogHandler(ctx context.Context, req *pb.CommonRequest) (*pb.CommonResponse, error) {
	return common.StartSteamLogHandler(ctx, req, initStreamLog)
}

func StopSteamLogHandler(ctx context.Context, req *pb.CommonRequest) (*pb.CommonResponse, error) {
	return common.StopSteamLogHandler(ctx, req)
}

func initStreamLog(h *common.HookMgr) {
	clientID := conf.Get().Client.ID
	clientSecret := conf.Get().Client.Secret

	handler, err := rpcclient.GetClientRPCSerivce().GetCli().PushServerStreamLog(
		context.Background())
	if err != nil {
		logrus.Error(err)
	}

	h.AddStream(func(msg string) {
		handler.Send(&pb.PushServerStreamLogReq{
			Log: []byte(utils.EncodeBase64(msg)),
			Base: &pb.ServerBase{
				ServerId:     clientID,
				ServerSecret: clientSecret,
			},
		})
	}, func() { handler.CloseSend() })
}
