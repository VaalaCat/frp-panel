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

func UpdateWorker(ctx *app.Context, req *pb.UpdateWorkerRequest) (*pb.UpdateWorkerResponse, error) {
	var (
		clientIds    = req.GetClientIds()
		wrokerReq    = req.GetWorker()
		userInfo     = common.GetUserInfo(ctx)
		oldClientIds []string
	)

	workerToUpdate, err := dao.NewQuery(ctx).GetWorkerByWorkerID(userInfo, wrokerReq.GetWorkerId())
	if err != nil {
		logger.Logger(ctx).WithError(err).Errorf("cannot get worker, id: [%s]", wrokerReq.GetWorkerId())
		return nil, fmt.Errorf("cannot get worker, id: [%s]", wrokerReq.GetWorkerId())
	}

	clis, err := dao.NewQuery(ctx).GetClientsByClientIDs(userInfo, clientIds)
	if err != nil {
		logger.Logger(ctx).WithError(err).Errorf("cannot get client, id: [%s]", utils.MarshalForJson(clientIds))
		return nil, fmt.Errorf("cannot get client, id: [%s]", utils.MarshalForJson(clientIds))
	}

	updatedFields := []string{}

	if len(clientIds) != 0 {
		oldClientIds = lo.Map(workerToUpdate.Clients, func(c models.Client, _ int) string { return c.ClientID })
		workerToUpdate.Clients = lo.Map(clis, func(c *models.Client, _ int) models.Client { return *c })
		updatedFields = append(updatedFields, "client_id")
	} else {
		oldClientIds = lo.Map(workerToUpdate.Clients, func(c models.Client, _ int) string { return c.ClientID })
	}

	if len(wrokerReq.GetName()) != 0 {
		workerToUpdate.Name = wrokerReq.GetName()
		updatedFields = append(updatedFields, "name")
	}

	if len(wrokerReq.GetCode()) != 0 {
		workerToUpdate.Code = wrokerReq.GetCode()
		updatedFields = append(updatedFields, "code")
	}

	if len(wrokerReq.GetConfigTemplate()) != 0 {
		workerToUpdate.ConfigTemplate = wrokerReq.GetConfigTemplate()
		updatedFields = append(updatedFields, "config_template")
	}

	if err := dao.NewQuery(ctx).UpdateWorker(userInfo, workerToUpdate); err != nil {
		logger.Logger(ctx).WithError(err).Errorf("cannot update worker, id: [%s]", wrokerReq.GetWorkerId())
		return nil, fmt.Errorf("cannot update worker, id: [%s]", wrokerReq.GetWorkerId())
	}

	go func() {
		bgCtx := ctx.Background()

		for _, oldClientId := range oldClientIds {
			removeResp := &pb.RemoveWorkerResponse{}
			err := rpc.CallClientWrapper(bgCtx, oldClientId, pb.Event_EVENT_REMOVE_WORKER, &pb.RemoveWorkerRequest{
				ClientId: &oldClientId,
				WorkerId: &workerToUpdate.ID,
			}, removeResp)
			if err != nil {
				logger.Logger(bgCtx).WithError(err).Errorf("remove old worker event send to client error, clients: [%s], worker name: [%s]", oldClientId, workerToUpdate.Name)
			}
		}

		for _, newClient := range clis {
			createResp := &pb.CreateWorkerResponse{}
			err = rpc.CallClientWrapper(bgCtx, newClient.ClientID, pb.Event_EVENT_CREATE_WORKER, &pb.CreateWorkerRequest{
				ClientId: &newClient.ClientID,
				Worker:   workerToUpdate.ToPB(),
			}, createResp)
			if err != nil {
				logger.Logger(bgCtx).WithError(err).Errorf("update new worker event send to client error, client id: [%s], worker name: [%s]", newClient.ClientID, workerToUpdate.Name)
			}
		}

		logger.Logger(ctx).Infof("update worker event send to client success, clients: [%s], worker name: [%s], remove old worker send to those clients: %s",
			utils.MarshalForJson(clis), workerToUpdate.Name, utils.MarshalForJson(oldClientIds))
	}()

	logger.Logger(ctx).Infof("update worker success, id: [%s], updated fields: %s", wrokerReq.GetWorkerId(), utils.MarshalForJson(updatedFields))

	return &pb.UpdateWorkerResponse{
		Status: &pb.Status{
			Code:    pb.RespCode_RESP_CODE_SUCCESS,
			Message: "ok",
		},
	}, nil
}
