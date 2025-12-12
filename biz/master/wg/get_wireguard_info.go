package wg

import (
	"errors"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
	"github.com/VaalaCat/frp-panel/services/rpc"
)

func GetWireGuardRuntimeInfo(ctx *app.Context, req *pb.GetWireGuardRuntimeInfoRequest) (*pb.GetWireGuardRuntimeInfoResponse, error) {
	log := ctx.Logger().WithField("op", "GetWireGuardRuntimeInfo")

	userInfo := common.GetUserInfo(ctx)
	if !userInfo.Valid() {
		log.Errorf("invalid user")
		return nil, errors.New("invalid user")
	}

	wgRecord, err := dao.NewQuery(ctx).GetWireGuardByID(userInfo, uint(req.GetId()))
	if err != nil {
		log.WithError(err).Errorf("get wireguard by id failed, clientId: [%s], id: [%d]", req.GetClientId(), req.GetId())
		return nil, errors.New("get wireguard by id failed")
	}

	resp := &pb.GetWireGuardRuntimeInfoResponse{}
	if err := rpc.CallClientWrapper(ctx, wgRecord.ClientID, pb.Event_EVENT_GET_WIREGUARD_RUNTIME_INFO, &pb.GetWireGuardRuntimeInfoRequest{
		InterfaceName: &wgRecord.Name,
	}, resp); err != nil {
		log.WithError(err).Errorf("failed to call get wireguard runtime info with clientId: [%s], id: [%d]", wgRecord.ClientID, req.GetId())
		return nil, errors.New("failed to call get wireguard runtime info")
	}
	resp.WgDeviceRuntimeInfo.ClientId = wgRecord.ClientID
	resp.WgDeviceRuntimeInfo.VirtualIp = wgRecord.LocalAddress

	log.Debugf("get wireguard runtime info success with clientId: [%s], interfaceName: [%s], runtimeInfo: [%s]",
		wgRecord.ClientID, wgRecord.Name, resp.GetWgDeviceRuntimeInfo().String())

	return resp, nil
}
