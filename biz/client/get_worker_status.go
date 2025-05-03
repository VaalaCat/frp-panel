package client

import (
	"fmt"

	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/utils/logger"
)

func GetWorkerStatus(ctx *app.Context, req *pb.GetWorkerStatusRequest) (*pb.GetWorkerStatusResponse, error) {
	if !ctx.GetApp().GetConfig().Client.Features.EnableFunctions {
		logger.Logger(ctx).Errorf("function features are not enabled")
		return nil, fmt.Errorf("function features are not enabled")
	}

	clientId := ctx.GetApp().GetConfig().Client.ID

	workersMgr := ctx.GetApp().GetWorkersManager()

	status, err := workersMgr.GetWorkerStatus(ctx, req.GetWorkerId())
	if err != nil {
		logger.Logger(ctx).Errorf("failed to get worker status: %v", err)
		return nil, fmt.Errorf("failed to get worker status: %v", err)
	}
	logger.Logger(ctx).Infof("get worker status for worker [%s], status: [%s]", req.GetWorkerId(), status)
	return &pb.GetWorkerStatusResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
		WorkerStatus: map[string]string{
			clientId: string(status),
		},
	}, nil
}
