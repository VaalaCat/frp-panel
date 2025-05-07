package worker

import (
	"fmt"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
	"github.com/VaalaCat/frp-panel/services/rpc"
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/VaalaCat/frp-panel/utils/logger"
	"github.com/samber/lo"
)

func RedeployWorker(ctx *app.Context, req *pb.RedeployWorkerRequest) (*pb.RedeployWorkerResponse, error) {
	var (
		clientIds    = req.GetClientIds()
		workerId     = req.GetWorkerId()
		userInfo     = common.GetUserInfo(ctx)
		oldClientIds []string
	)

	if len(workerId) == 0 {
		logger.Logger(ctx).Errorf("redeploy worker, worker id req: [%s]", req.String())
		return nil, fmt.Errorf("worker id is empty")

	}

	workerToUpdate, err := dao.NewQuery(ctx).GetWorkerByWorkerID(userInfo, workerId)
	if err != nil {
		logger.Logger(ctx).WithError(err).Errorf("redeploy worker cannot get worker, id: [%s]", workerId)
		return nil, fmt.Errorf("cannot get worker, id: [%s]", workerId)
	}

	clis := []*models.Client{}
	if len(clientIds) != 0 {
		clis, err = dao.NewQuery(ctx).GetClientsByClientIDs(userInfo, clientIds)
		if err != nil {
			logger.Logger(ctx).WithError(err).Errorf("redeploy worker cannot get client, id: [%s]", utils.MarshalForJson(clientIds))
			return nil, fmt.Errorf("cannot get client, id: [%s]", utils.MarshalForJson(clientIds))
		}
	} else {
		logger.Logger(ctx).Infof("redeploy worker, no clientId, redeploy on all clients")
	}

	var clisToRedeploy []string
	allCliIds := lo.Map(workerToUpdate.Clients, func(c models.Client, _ int) string { return c.ClientID })
	reqCliIds := lo.SliceToMap(clis, func(c *models.Client) (string, struct{}) { return c.ClientID, struct{}{} })
	// 如果大于0，需要过滤
	if len(clis) > 0 {
		clisToRedeploy = lo.Filter(allCliIds, func(c string, _ int) bool {
			_, ok := reqCliIds[c]
			return ok
		})
	} else {
		clisToRedeploy = allCliIds
	}

	go func() {
		bgCtx := ctx.Background()

		for _, cliId := range clisToRedeploy {
			removeResp := &pb.RemoveWorkerResponse{}
			err := rpc.CallClientWrapper(bgCtx, cliId, pb.Event_EVENT_REMOVE_WORKER, &pb.RemoveWorkerRequest{
				ClientId: &cliId,
				WorkerId: &workerToUpdate.ID,
			}, removeResp)
			if err != nil {
				logger.Logger(bgCtx).WithError(err).Errorf("remove old worker event send to client error, clients: [%s], worker name: [%s]", cliId, workerToUpdate.Name)
			}

			createResp := &pb.CreateWorkerResponse{}
			err = rpc.CallClientWrapper(bgCtx, cliId, pb.Event_EVENT_CREATE_WORKER, &pb.CreateWorkerRequest{
				ClientId: &cliId,
				Worker:   workerToUpdate.ToPB(),
			}, createResp)
			if err != nil {
				logger.Logger(bgCtx).WithError(err).Errorf("update new worker event send to client error, client id: [%s], worker name: [%s]", cliId, workerToUpdate.Name)
			}
		}

		logger.Logger(ctx).Infof("redeploy worker event send to client success, clients: [%s], worker name: [%s], remove old worker send to those clients: %s",
			utils.MarshalForJson(clis), workerToUpdate.Name, utils.MarshalForJson(oldClientIds))
	}()

	logger.Logger(ctx).Infof("redeploy worker success, id: [%s], clients: %s", workerId, utils.MarshalForJson(clientIds))

	return &pb.RedeployWorkerResponse{
		Status: &pb.Status{
			Code:    pb.RespCode_RESP_CODE_SUCCESS,
			Message: "ok",
		},
	}, nil
}
