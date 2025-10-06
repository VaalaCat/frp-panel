package wg

import (
	"errors"

	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
)

func ReportWireGuardRuntimeInfo(ctx *app.Context, req *pb.ReportWireGuardRuntimeInfoReq) (*pb.ReportWireGuardRuntimeInfoResp, error) {
	var (
		interfaceName = req.GetInterfaceName()
		clientId      = req.GetBase().GetClientId()
		log           = ctx.Logger().WithField("op", "ReportWireGuardRuntimeInfo")
	)

	wgIfce, err := dao.NewQuery(ctx).AdminGetWireGuardByClientIDAndInterfaceName(clientId, interfaceName)
	if err != nil {
		log.WithError(err).Errorf("failed to get wireguard by client id and interface name, clientId: [%s], interfaceName: [%s]", clientId, interfaceName)
		return nil, errors.New("failed to get wireguard by client id and interface name")
	}

	networkTopologyCache := ctx.GetApp().GetNetworkTopologyCache()
	networkTopologyCache.SetRuntimeInfo(uint(wgIfce.ID), req.GetRuntimeInfo())

	return &pb.ReportWireGuardRuntimeInfoResp{
		Status: &pb.Status{
			Code:    pb.RespCode_RESP_CODE_SUCCESS,
			Message: "success",
		},
	}, nil
}
