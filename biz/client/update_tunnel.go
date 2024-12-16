package client

import (
	"context"
	"reflect"

	"github.com/VaalaCat/frp-panel/logger"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/client"
	"github.com/VaalaCat/frp-panel/tunnel"
	"github.com/VaalaCat/frp-panel/utils"
)

func UpdateFrpcHander(ctx context.Context, req *pb.UpdateFRPCRequest) (*pb.UpdateFRPCResponse, error) {
	logger.Logger(ctx).Infof("update frpc, req: [%+v]", req)
	content := req.GetConfig()
	c, p, v, err := utils.LoadClientConfig(content, false)
	if err != nil {
		logger.Logger(ctx).WithError(err).Errorf("cannot load config")
		return &pb.UpdateFRPCResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: err.Error()},
		}, err
	}

	cli := tunnel.GetClientController().Get(req.GetClientId(), req.GetServerId())
	if cli != nil {
		if reflect.DeepEqual(c, cli.GetCommonCfg()) {
			logger.Logger(ctx).Warnf("client common config not changed")
			cli.Update(p, v)
		} else {
			cli.Stop()
			tunnel.GetClientController().Delete(req.GetClientId(), req.GetServerId())
			tunnel.GetClientController().Add(req.GetClientId(), req.GetServerId(), client.NewClientHandler(c, p, v))
			tunnel.GetClientController().Run(req.GetClientId(), req.GetServerId())
		}
		logger.Logger(ctx).Infof("update client, id: [%s] success, running", req.GetClientId())
	} else {
		tunnel.GetClientController().Add(req.GetClientId(), req.GetServerId(), client.NewClientHandler(c, p, v))
		tunnel.GetClientController().Run(req.GetClientId(), req.GetServerId())
		logger.Logger(ctx).Infof("add new client, id: [%s], running", req.GetClientId())
	}

	return &pb.UpdateFRPCResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
	}, nil
}
