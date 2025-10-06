package wg

import (
	"errors"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
	"github.com/VaalaCat/frp-panel/services/rpc"
	wgsvc "github.com/VaalaCat/frp-panel/services/wg"
)

func DeleteWireGuard(ctx *app.Context, req *pb.DeleteWireGuardRequest) (*pb.DeleteWireGuardResponse, error) {
	userInfo := common.GetUserInfo(ctx)
	log := ctx.Logger().WithField("op", "DeleteWireGuard")
	if !userInfo.Valid() {
		return &pb.DeleteWireGuardResponse{Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "invalid user"}}, nil
	}
	id := uint(req.GetId())
	if id == 0 {
		return &pb.DeleteWireGuardResponse{Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "invalid id"}}, nil
	}

	wgToDelete, err := dao.NewQuery(ctx).GetWireGuardByID(userInfo, id)
	if err != nil {
		log.WithError(err).Errorf("get wireguard by id failed")
		return nil, err
	}

	if err := dao.NewQuery(ctx).DeleteWireGuard(userInfo, id); err != nil {
		log.WithError(err).Errorf("delete wireguard failed")
		return nil, err
	}

	log.Debugf("delete wireguard success, id: %d", id)

	ctxBg := ctx.Background()

	go func() {
		resp, err := rpc.CallClient(ctxBg, wgToDelete.ClientID, pb.Event_EVENT_DELETE_WIREGUARD, &pb.DeleteWireGuardRequest{
			ClientId:      &wgToDelete.ClientID,
			InterfaceName: &wgToDelete.Name,
		})
		if err != nil {
			log.WithError(err).Errorf("delete wireguard event send to client failed")
		}

		if resp == nil {
			log.Errorf("cannot get response, client id: [%s]", wgToDelete.ClientID)
		}

		peers, err := dao.NewQuery(ctx).GetWireGuardsByNetworkID(userInfo, uint(wgToDelete.NetworkID))
		if err != nil {
			log.WithError(err).Errorf("get wireguards by network id failed")
			return
		}

		if len(peers) == 0 {
			log.Infof("no wireguards in network, network id: [%d]", wgToDelete.NetworkID)
			return
		}

		links, err := dao.NewQuery(ctx).ListWireGuardLinksByNetwork(userInfo, uint(wgToDelete.NetworkID))
		if err != nil {
			log.WithError(err).Errorf("get wireguard links by network id failed")
			return
		}

		peerConfigs, err := wgsvc.PlanAllowedIPs(peers, links,
			wgsvc.DefaultRoutingPolicy(
				wgsvc.NewACL().LoadFromPB(wgToDelete.Network.ACL.Data),
				ctx.GetApp().GetNetworkTopologyCache(),
			))
		if err != nil {
			log.WithError(err).Errorf("build peer configs for network failed")
			return
		}

		for _, peer := range peers {
			if err := emitPatchWireGuardEventToClient(ctx, peer, peerConfigs[peer.ID]); err != nil {
				log.WithError(err).Errorf("patch wireguard event send to client error")
				continue
			}

			log.Debugf("update config to client success, client id: [%s], wireguard interface: [%s]", peer.ClientID, peer.Name)
		}
	}()

	return &pb.DeleteWireGuardResponse{Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "success"}}, nil
}

func emitDeleteWireGuardEventToClient(ctx *app.Context, peerNeedRemoveWg *models.WireGuard, wgToDelete *models.WireGuard) error {
	log := ctx.Logger().WithField("op", "emitDeleteWireGuardEventToClient")
	userInfo := common.GetUserInfo(ctx)
	if !userInfo.Valid() {
		return errors.New("invalid user")
	}

	resp := &pb.UpdateWireGuardResponse{}

	err := rpc.CallClientWrapper(ctx, peerNeedRemoveWg.ClientID, pb.Event_EVENT_UPDATE_WIREGUARD, &pb.UpdateWireGuardRequest{
		WireguardConfig: &pb.WireGuardConfig{
			InterfaceName: peerNeedRemoveWg.Name,
			Peers:         []*pb.WireGuardPeerConfig{{ClientId: peerNeedRemoveWg.ClientID}},
		},
		UpdateType: pb.UpdateWireGuardRequest_UPDATE_TYPE_REMOVE_PEER.Enum(),
	}, resp)
	if err != nil {
		log.WithError(err).Errorf("delete wireguard event send to client error")
		return err
	}

	log.Infof("delete wireguard event send to client success, client id: [%s], wireguard interface: [%s]",
		wgToDelete.ClientID, wgToDelete.Name)
	return nil
}
