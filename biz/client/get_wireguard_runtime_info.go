package client

import (
	"fmt"

	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
)

func GetWireGuardRuntimeInfo(ctx *app.Context, req *pb.GetWireGuardRuntimeInfoRequest) (*pb.GetWireGuardRuntimeInfoResponse, error) {
	var (
		interfaceName = req.GetInterfaceName()
		log           = ctx.Logger().WithField("op", "GetWireGuardRuntimeInfo")
	)

	if interfaceName == "" {
		log.Errorf("interface_name is required")
		return nil, fmt.Errorf("interface_name is required")
	}

	wgSvc, ok := ctx.GetApp().GetWireGuardManager().GetService(interfaceName)
	if !ok {
		log.Errorf("wireguard service not found, interface_name: %s", interfaceName)
		return nil, fmt.Errorf("wireguard service not found, interface_name: %s", interfaceName)
	}

	runtimeInfo, err := wgSvc.GetWGRuntimeInfo()
	if err != nil {
		log.WithError(err).Errorf("get wireguard runtime info failed")
		return nil, fmt.Errorf("get wireguard runtime info failed: %v", err)
	}

	log.Debugf("get wireguard runtime info with interface_name: %s, runtimeInfo: %s", interfaceName, runtimeInfo.String())

	return &pb.GetWireGuardRuntimeInfoResponse{
		Status: &pb.Status{
			Code:    pb.RespCode_RESP_CODE_SUCCESS,
			Message: "success",
		},
		WgDeviceRuntimeInfo: runtimeInfo,
	}, nil
}
