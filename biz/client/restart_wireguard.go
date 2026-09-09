package client

import (
	"errors"

	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
)

func RestartWireGuard(ctx *app.Context, req *pb.RestartWireGuardRequest) (*pb.RestartWireGuardResponse, error) {
	log := ctx.Logger().WithField("op", "RestartWireGuard")

	if req == nil || len(req.GetInterfaceName()) == 0 {
		return nil, errors.New("invalid interface name")
	}

	if err := ctx.GetApp().GetWireGuardManager().RestartService(req.GetInterfaceName()); err != nil {
		log.WithError(err).Errorf("restart wireguard service failed, interface: %s", req.GetInterfaceName())
		return nil, err
	}

	log.Infof("restart wireguard service success, interface: %s", req.GetInterfaceName())

	return &pb.RestartWireGuardResponse{
		Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "success"},
	}, nil
}
