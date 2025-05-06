package worker

import (
	"fmt"
	"maps"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
	"github.com/VaalaCat/frp-panel/services/rpc"
	"github.com/VaalaCat/frp-panel/utils/logger"
	"github.com/samber/lo"
	"github.com/sourcegraph/conc/pool"
)

func GetWorkerStatus(ctx *app.Context, req *pb.GetWorkerStatusRequest) (*pb.GetWorkerStatusResponse, error) {
	var (
		workerID = req.GetWorkerId()
		userInfo = common.GetUserInfo(ctx)
	)

	if len(workerID) == 0 {
		logger.Logger(ctx).Errorf("worker id is empty")
		return nil, fmt.Errorf("worker id is empty")
	}

	workerRecord, err := dao.NewQuery(ctx).GetWorkerByWorkerID(userInfo, workerID)
	if err != nil {
		logger.Logger(ctx).WithError(err).Errorf("get worker by id failed")
		return nil, err
	}

	clientIds := lo.Map(workerRecord.Clients, func(cli models.Client, _ int) string {
		return cli.ClientID
	})

	var pool pool.ResultErrorPool[*pb.GetWorkerStatusResponse]
	for _, clientID := range clientIds {
		pool.Go(func() (*pb.GetWorkerStatusResponse, error) {
			bgCtx := ctx.Background()
			cliResp := &pb.GetWorkerStatusResponse{}
			err := rpc.CallClientWrapper(bgCtx, clientID, pb.Event_EVENT_GET_WORKER_STATUS, &pb.GetWorkerStatusRequest{}, cliResp)
			return cliResp, err
		})
	}

	resps, err := pool.Wait()
	if err != nil {
		logger.Logger(ctx).WithError(err).Warnf("get worker status failed")
	}

	statusMap := map[string]string{}

	for _, r := range resps {
		s := r.GetWorkerStatus()
		maps.Copy(statusMap, s)
	}

	return &pb.GetWorkerStatusResponse{
		Status:       &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
		WorkerStatus: statusMap,
	}, nil
}
