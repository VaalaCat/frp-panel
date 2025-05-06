package client

import (
	"context"

	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/workerd"
	"github.com/VaalaCat/frp-panel/utils/logger"
)

func PullWorkers(appInstance app.Application, clientID, clientSecret string) error {
	ctx := app.NewContext(context.Background(), appInstance)

	if !ctx.GetApp().GetConfig().Client.Features.EnableFunctions {
		logger.Logger(ctx).Infof("function features are not enabled")
		return nil
	}

	logger.Logger(ctx).Infof("start to pull workers belong to client, clientID: [%s]", clientID)

	cli := ctx.GetApp().GetMasterCli()

	resp, err := cli.Call().ListClientWorkers(ctx, &pb.ListClientWorkersRequest{
		Base: &pb.ClientBase{
			ClientId:     clientID,
			ClientSecret: clientSecret,
		},
	})
	if err != nil {
		logger.Logger(ctx).WithError(err).Error("cannot list client workers")
		return err
	}

	if len(resp.GetWorkers()) == 0 {
		logger.Logger(ctx).Infof("client [%s] has no workers", clientID)
		return nil
	}

	ctrl := ctx.GetApp().GetWorkersManager()
	for _, worker := range resp.GetWorkers() {
		ctrl.RunWorker(ctx, worker.GetWorkerId(), workerd.NewWorkerdController(worker, ctx.GetApp().GetConfig().Client.Worker.WorkerdWorkDir))
	}

	logger.Logger(ctx).Infof("pull workers belong to client success, clientID: [%s], will run [%d] workers", clientID, len(resp.GetWorkers()))

	return nil
}
