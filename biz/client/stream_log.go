package client

import (
	"context"

	"github.com/VaalaCat/frp-panel/biz/common"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/VaalaCat/frp-panel/utils/logger"
)

func StartSteamLogHandler(ctx *app.Context, req *pb.StartSteamLogRequest) (*pb.CommonResponse, error) {
	return common.StartSteamLogHandler(ctx, req, initStreamLog)
}

func StopSteamLogHandler(ctx *app.Context, req *pb.CommonRequest) (*pb.CommonResponse, error) {
	return common.StopSteamLogHandler(ctx, req)
}

func initStreamLog(ctx *app.Context, h app.StreamLogHookMgr) {
	clientID := ctx.GetApp().GetConfig().Client.ID
	clientSecret := ctx.GetApp().GetConfig().Client.Secret

	handler, err := ctx.GetApp().GetClientRPCHandler().GetCli().Call().PushClientStreamLog(
		context.Background())
	if err != nil {
		logger.Logger(ctx).Error(err)
	}

	h.AddStream(func(msg string) {
		handler.Send(&pb.PushClientStreamLogReq{
			Log: []byte(utils.EncodeBase64(msg)),
			Base: &pb.ClientBase{
				ClientId:     clientID,
				ClientSecret: clientSecret,
			},
		})
	}, func() { handler.CloseSend() })
}
