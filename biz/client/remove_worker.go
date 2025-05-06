package client

import (
	"fmt"

	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/utils/logger"
)

func RemoveWorker(ctx *app.Context, req *pb.RemoveWorkerRequest) (*pb.RemoveWorkerResponse, error) {
	if !ctx.GetApp().GetConfig().Client.Features.EnableFunctions {
		logger.Logger(ctx).Errorf("function features are not enabled")
		return nil, fmt.Errorf("function features are not enabled")
	}

	mgr := ctx.GetApp().GetWorkersManager()

	workerId := req.GetWorkerId()
	logger.Logger(ctx).Infof("start remove worker, id: [%s]", workerId)

	if err := mgr.StopWorker(ctx, workerId); err != nil {
		logger.Logger(ctx).WithError(err).Errorf("cannot remove worker, id: [%s]", workerId)
		return nil, err
	}

	logger.Logger(ctx).Infof("remove worker success, id: [%s]", workerId)
	return &pb.RemoveWorkerResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
	}, nil
}
