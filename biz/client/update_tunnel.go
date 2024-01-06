package client

import (
	"context"
	"fmt"
	"reflect"

	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/client"
	"github.com/VaalaCat/frp-panel/tunnel"
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/sirupsen/logrus"
)

func UpdateFrpcHander(ctx context.Context, req *pb.UpdateFRPCRequest) (*pb.UpdateFRPCResponse, error) {
	logrus.Infof("update frpc, req: [%+v]", req)
	content := req.GetConfig()
	c, p, v, err := utils.LoadClientConfig(content, false)
	if err != nil {
		logrus.WithError(err).Errorf("cannot load config")
		return &pb.UpdateFRPCResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: err.Error()},
		}, err
	}

	cli := tunnel.GetClientController().Get(req.GetClientId())
	if cli == nil {
		logrus.Errorf("cannot get client, id: [%s]", req.GetClientId())
		return &pb.UpdateFRPCResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "cannot get client"},
		}, fmt.Errorf("cannot get client")
	}

	if reflect.DeepEqual(c, cli.GetCommonCfg()) {
		logrus.Warnf("client config not changed")
		cli.Update(p, v)
	} else {
		tcli := tunnel.GetClientController().Get(req.GetClientId())
		if tcli != nil {
			tcli.Stop()
			tunnel.GetClientController().Delete(req.GetClientId())
		}
		tunnel.GetClientController().Add(req.GetClientId(), client.NewClientHandler(c, p, v))
		tunnel.GetClientController().Run(req.GetClientId())
	}

	return &pb.UpdateFRPCResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "ok"},
	}, nil
}
