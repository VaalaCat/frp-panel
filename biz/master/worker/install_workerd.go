package worker

import (
	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
	"github.com/VaalaCat/frp-panel/services/rpc"
	"github.com/VaalaCat/frp-panel/utils/logger"
)

func InstallWorkerd(ctx *app.Context, req *pb.InstallWorkerdRequest) (*pb.InstallWorkerdResponse, error) {
	var (
		userInfo = common.GetUserInfo(ctx)
		clientId = req.GetClientId()
	)
	logger.Logger(ctx).Infof("installw orkerd called with userInfo: %v, clientId: %s", userInfo, clientId)

	_, err := dao.NewQuery(ctx).GetClientByClientID(userInfo, clientId)
	if err != nil {
		logger.Logger(ctx).WithError(err).Errorf("failed to get client by clientID: %s", clientId)
		return nil, err
	}

	resp := &pb.InstallWorkerdResponse{}
	if err := rpc.CallClientWrapper(ctx, clientId, pb.Event_EVENT_INSTALL_WORKERD, req, resp); err != nil {
		logger.Logger(ctx).WithError(err).Errorf("failed to call install workerd with clientId: %s", clientId)
		return nil, err
	}
	logger.Logger(ctx).Infof("install workerd success with clientId: %s", clientId)

	return &pb.InstallWorkerdResponse{
		Status: &pb.Status{
			Code:    pb.RespCode_RESP_CODE_SUCCESS,
			Message: "ok",
		},
	}, nil
}
