package worker

import (
	"github.com/VaalaCat/frp-panel/biz/master/proxy"
	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
	"github.com/VaalaCat/frp-panel/services/rpc"
	"github.com/VaalaCat/frp-panel/utils/logger"
	"github.com/samber/lo"
)

func RemoveWorker(ctx *app.Context, req *pb.RemoveWorkerRequest) (*pb.RemoveWorkerResponse, error) {
	var (
		userInfo = common.GetUserInfo(ctx)
		workerId = req.GetWorkerId()
	)

	logger.Logger(ctx).Infof("start remove worker, id: [%s]", workerId)

	workerToDelete, err := dao.NewQuery(ctx).GetWorkerByWorkerID(userInfo, workerId)
	if err != nil {
		logger.Logger(ctx).WithError(err).Errorf("cannot get worker, id: [%s]", workerId)
		return nil, err
	}

	if ingressesToDelete, err := dao.NewQuery(ctx).GetProxyConfigsByWorkerId(userInfo, workerId); err == nil {
		for _, ingressToDelete := range ingressesToDelete {
			logger.Logger(ctx).Infof("start to remove worker ingress on server: [%s] client: [%s], name: [%s]", ingressToDelete.ServerID, ingressToDelete.ClientID, ingressToDelete.Name)

			resp, err := proxy.DeleteProxyConfig(ctx, &pb.DeleteProxyConfigRequest{
				ServerId: lo.ToPtr(ingressToDelete.ServerID),
				ClientId: lo.ToPtr(ingressToDelete.ClientID),
				Name:     lo.ToPtr(ingressToDelete.Name),
			})

			if err != nil {
				logger.Logger(ctx).WithError(err).Errorf("cannot remove worker ingress on server: [%s] client: [%s], name: [%s], resp is: [%s]",
					ingressToDelete.ServerID, ingressToDelete.ClientID, ingressToDelete.Name, resp.String())
				return nil, err
			}
			logger.Logger(ctx).Infof("remove worker ingress success on server: [%s] client: [%s], name: [%s]", ingressToDelete.ServerID, ingressToDelete.ClientID, ingressToDelete.Name)
		}
	}

	if err := dao.NewQuery(ctx).DeleteWorker(userInfo, workerId); err != nil {
		logger.Logger(ctx).WithError(err).Errorf("cannot remove worker, id: [%s]", workerId)
		return nil, err
	}

	go func() {
		bgCtx := ctx.Background()
		hasErr := false

		for _, cli := range workerToDelete.Clients {
			resp := &pb.RemoveWorkerResponse{}
			if err := rpc.CallClientWrapper(bgCtx, cli.ClientID, pb.Event_EVENT_REMOVE_WORKER, req, resp); err != nil {
				logger.Logger(bgCtx).WithError(err).Errorf("remove event send to client error, client id: [%s]", cli.ClientID)
				hasErr = true
				continue
			}
		}

		if hasErr {
			logger.Logger(bgCtx).Errorf("remove event send to client error")
		}
	}()

	logger.Logger(ctx).Infof("remove worker success, id: [%s]", workerId)

	return &pb.RemoveWorkerResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
	}, nil
}
