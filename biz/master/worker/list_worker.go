package worker

import (
	"fmt"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
	"github.com/VaalaCat/frp-panel/utils/logger"
	"github.com/samber/lo"
)

func ListWorkers(ctx *app.Context, req *pb.ListWorkersRequest) (*pb.ListWorkersResponse, error) {
	var (
		userInfo     = common.GetUserInfo(ctx)
		page         = int(req.GetPage())
		pageSize     = int(req.GetPageSize())
		keyword      = req.GetKeyword()
		err          error
		workers      []*models.Worker
		workerCounts int64
		hasKeyword   = len(keyword) > 0
	)

	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	if hasKeyword {
		workers, err = dao.NewQuery(ctx).ListWorkersWithKeyword(userInfo, page, pageSize, keyword)
	} else {
		workers, err = dao.NewQuery(ctx).ListWorkers(userInfo, page, pageSize)
	}
	if err != nil {
		logger.Logger(ctx).WithError(err).Errorf("cannot list workers, page: [%d], pageSize: [%d], keyword: [%s]", page, pageSize, keyword)
		return &pb.ListWorkersResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_NOT_FOUND, Message: err.Error()},
		}, fmt.Errorf("cannot list workers, page: [%d], pageSize: [%d], keyword: [%s]", page, pageSize, keyword)
	}

	if hasKeyword {
		workerCounts, err = dao.NewQuery(ctx).CountWorkersWithKeyword(userInfo, keyword)
	} else {
		workerCounts, err = dao.NewQuery(ctx).CountWorkers(userInfo)
	}
	if err != nil {
		logger.Logger(ctx).WithError(err).Errorf("cannot count workers, keyword: [%s]", keyword)
		return &pb.ListWorkersResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_NOT_FOUND, Message: err.Error()},
		}, fmt.Errorf("cannot count workers, keyword: [%s]", keyword)
	}

	return &pb.ListWorkersResponse{
		Status: &pb.Status{
			Code:    pb.RespCode_RESP_CODE_SUCCESS,
			Message: "success",
		},
		Total: lo.ToPtr(int32(workerCounts)),
		Workers: lo.Map(workers, func(w *models.Worker, _ int) *pb.Worker {
			k := w.ToPB()
			k.Code = nil
			k.ConfigTemplate = nil
			return k
		}),
	}, nil
}

func ListClientWorkers(ctx *app.Context, req *pb.ListClientWorkersRequest) (*pb.ListClientWorkersResponse, error) {
	var (
		err      error
		workers  []*models.Worker
		clientId = req.GetBase().GetClientId()
	)

	workers, err = dao.NewQuery(ctx).AdminListWorkersByClientID(clientId)
	if err != nil {
		logger.Logger(ctx).WithError(err).Errorf("cannot list workers, clientId: [%s]", clientId)
		return &pb.ListClientWorkersResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_NOT_FOUND, Message: err.Error()},
		}, fmt.Errorf("cannot list workers, clientId: [%s]", clientId)
	}

	logger.Logger(ctx).Infof("list workers, clientId: [%s], worker len: [%d]", clientId, len(workers))

	return &pb.ListClientWorkersResponse{
		Status: &pb.Status{
			Code:    pb.RespCode_RESP_CODE_SUCCESS,
			Message: "success",
		},
		Workers: lo.Map(workers, func(w *models.Worker, _ int) *pb.Worker {
			k := w.ToPB()
			return k
		}),
	}, nil
}
