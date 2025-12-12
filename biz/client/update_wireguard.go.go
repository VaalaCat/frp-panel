package client

import (
	"errors"

	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/samber/lo"
)

func UpdateWireGuard(ctx *app.Context, req *pb.UpdateWireGuardRequest) (*pb.UpdateWireGuardResponse, error) {
	log := ctx.Logger().WithField("op", "UpdateWireGuard")

	wgSvc, ok := ctx.GetApp().GetWireGuardManager().GetService(req.GetWireguardConfig().GetInterfaceName())
	if !ok || wgSvc == nil {
		log.Errorf("get wireguard service failed")
		return nil, errors.New("wireguard service not found")
	}

	switch req.GetUpdateType() {
	case pb.UpdateWireGuardRequest_UPDATE_TYPE_ADD_PEER:
		return AddPeer(ctx, wgSvc, req)
	case pb.UpdateWireGuardRequest_UPDATE_TYPE_REMOVE_PEER:
		return RemovePeer(ctx, wgSvc, req)
	case pb.UpdateWireGuardRequest_UPDATE_TYPE_UPDATE_PEER:
		return UpdatePeer(ctx, wgSvc, req)
	case pb.UpdateWireGuardRequest_UPDATE_TYPE_PATCH_PEERS:
		return PatchPeers(ctx, wgSvc, req)
	default:
	}

	return nil, errors.New("update type not found, please check the update type in the request")
}

func AddPeer(ctx *app.Context, wgSvc app.WireGuard, req *pb.UpdateWireGuardRequest) (*pb.UpdateWireGuardResponse, error) {
	log := ctx.Logger().WithField("op", "AddPeer")

	log.Debugf("add peer, peer_config: %+v", req.GetWireguardConfig().GetPeers())

	for _, peer := range req.GetWireguardConfig().GetPeers() {
		err := wgSvc.AddPeer(&defs.WireGuardPeerConfig{WireGuardPeerConfig: peer})
		if err != nil {
			log.WithError(err).Errorf("add peer failed")
			continue
		}
	}

	if err := wgSvc.UpdateAdjs(req.GetWireguardConfig().GetAdjs()); err != nil {
		log.WithError(err).Errorf("update adjs failed, adjs: %+v", req.GetWireguardConfig().GetAdjs())
		return nil, err
	}

	log.Infof("add peer done")

	return &pb.UpdateWireGuardResponse{Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "success"}}, nil
}

func RemovePeer(ctx *app.Context, wgSvc app.WireGuard, req *pb.UpdateWireGuardRequest) (*pb.UpdateWireGuardResponse, error) {
	log := ctx.Logger().WithField("op", "RemovePeer")

	log.Debugf("remove peer, peer_config: %+v", req.GetWireguardConfig().GetPeers())

	for _, peer := range req.GetWireguardConfig().GetPeers() {
		err := wgSvc.RemovePeer(peer.GetPublicKey())
		if err != nil {
			log.WithError(err).Errorf("remove peer failed")
			continue
		}
	}

	if err := wgSvc.UpdateAdjs(req.GetWireguardConfig().GetAdjs()); err != nil {
		log.WithError(err).Errorf("update adjs failed, adjs: %+v", req.GetWireguardConfig().GetAdjs())
		return nil, err
	}

	log.Infof("remove peer done")

	return &pb.UpdateWireGuardResponse{Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "success"}}, nil
}

func UpdatePeer(ctx *app.Context, wgSvc app.WireGuard, req *pb.UpdateWireGuardRequest) (*pb.UpdateWireGuardResponse, error) {
	log := ctx.Logger().WithField("op", "UpdatePeer")

	log.Debugf("update peer, peer_config: %+v", req.GetWireguardConfig().GetPeers())

	for _, peer := range req.GetWireguardConfig().GetPeers() {
		err := wgSvc.UpdatePeer(&defs.WireGuardPeerConfig{WireGuardPeerConfig: peer})
		if err != nil {
			log.WithError(err).Errorf("update peer failed")
			continue
		}
	}

	if err := wgSvc.UpdateAdjs(req.GetWireguardConfig().GetAdjs()); err != nil {
		log.WithError(err).Errorf("update adjs failed, adjs: %+v", req.GetWireguardConfig().GetAdjs())
		return nil, err
	}

	log.Infof("update peer done")

	return &pb.UpdateWireGuardResponse{Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "success"}}, nil
}

func PatchPeers(ctx *app.Context, wgSvc app.WireGuard, req *pb.UpdateWireGuardRequest) (*pb.UpdateWireGuardResponse, error) {
	log := ctx.Logger().WithField("op", "PatchPeers")

	log.Debugf("patch peers, peer_config: %+v", req.GetWireguardConfig().GetPeers())

	wgCfg := &defs.WireGuardConfig{WireGuardConfig: req.GetWireguardConfig()}

	diffResp, err := wgSvc.PatchPeers(wgCfg.GetParsedPeers())
	if err != nil {
		log.WithError(err).Errorf("patch peers failed")
		return nil, err
	}

	if err = wgSvc.UpdateAdjs(req.GetWireguardConfig().GetAdjs()); err != nil {
		log.WithError(err).Errorf("update adjs failed, adjs: %+v", req.GetWireguardConfig().GetAdjs())
		return nil, err
	}

	log.Debugf("patch peers done, add_peers: %+v, remove_peers: %+v",
		lo.Map(diffResp.AddPeers, func(item *defs.WireGuardPeerConfig, _ int) string { return item.GetClientId() }),
		lo.Map(diffResp.RemovePeers, func(item *defs.WireGuardPeerConfig, _ int) string { return item.GetClientId() }))

	return &pb.UpdateWireGuardResponse{Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "success"}}, nil
}
