package client

import (
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
)

func DeleteWireGuard(ctx *app.Context, req *pb.DeleteWireGuardRequest) (*pb.DeleteWireGuardResponse, error) {
	log := ctx.Logger().WithField("op", "DeleteWireGuard")

	log.Debugf("delete wireguard service, client_id: %s, interface_name: %s", req.GetClientId(), req.GetInterfaceName())

	err := ctx.GetApp().GetWireGuardManager().RemoveService(req.GetInterfaceName())
	if err != nil {
		log.WithError(err).Errorf("remove wireguard service failed")
		return nil, err
	}

	log.Debugf("remove wireguard service success")

	return &pb.DeleteWireGuardResponse{Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "success"}}, nil
}
