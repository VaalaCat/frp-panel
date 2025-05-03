package client

import (
	"fmt"
	"os"

	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/utils/logger"
)

func InstallWorkerd(ctx *app.Context, req *pb.InstallWorkerdRequest) (*pb.InstallWorkerdResponse, error) {
	if !ctx.GetApp().GetConfig().Client.Features.EnableFunctions {
		logger.Logger(ctx).Errorf("function features are not enabled")
		return nil, fmt.Errorf("function features are not enabled")
	}

	workersMgr := ctx.GetApp().GetWorkersManager()

	cwd, err := os.Getwd()
	if err != nil {
		logger.Logger(ctx).Errorf("failed to get current working directory: %v, will install workerd in /usr/local/bin", err)
	}

	binPath, err := workersMgr.InstallWorkerd(ctx, req.GetDownloadUrl(), cwd)
	if err != nil {
		logger.Logger(ctx).Errorf("failed to install workerd: %v", err)
		return nil, fmt.Errorf("failed to install workerd: %v", err)
	}

	execMgr := ctx.GetApp().GetWorkerExecManager()
	execMgr.UpdateBinaryPath(binPath)

	return &pb.InstallWorkerdResponse{
		Status: &pb.Status{
			Code:    pb.RespCode_RESP_CODE_SUCCESS,
			Message: "ok",
		},
	}, nil
}
