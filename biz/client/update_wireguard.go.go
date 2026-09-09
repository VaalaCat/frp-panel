package client

import (
	"errors"

	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
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

	// 主链路：先更新 adjs（保证后续 wg 内部的预连接/清理逻辑使用最新拓扑）
	if err := updateAdjsFirst(log, wgSvc, req); err != nil {
		return nil, err
	}

	applyPeerOps(log, req.GetWireguardConfig().GetPeers(), "add peer", func(peer *pb.WireGuardPeerConfig) error {
		return wgSvc.AddPeer(&defs.WireGuardPeerConfig{WireGuardPeerConfig: peer})
	})

	log.Infof("add peer done")

	return &pb.UpdateWireGuardResponse{Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "success"}}, nil
}

func RemovePeer(ctx *app.Context, wgSvc app.WireGuard, req *pb.UpdateWireGuardRequest) (*pb.UpdateWireGuardResponse, error) {
	log := ctx.Logger().WithField("op", "RemovePeer")

	log.Debugf("remove peer, peer_config: %+v", req.GetWireguardConfig().GetPeers())

	// 主链路：先更新 adjs（保证后续 wg 内部的预连接/清理逻辑使用最新拓扑）
	if err := updateAdjsFirst(log, wgSvc, req); err != nil {
		return nil, err
	}

	applyPeerOps(log, req.GetWireguardConfig().GetPeers(), "remove peer routes", func(peer *pb.WireGuardPeerConfig) error {
		return wgSvc.RemovePeer(peer.GetPublicKey())
	})

	log.Infof("remove peer done")

	return &pb.UpdateWireGuardResponse{Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "success"}}, nil
}

func UpdatePeer(ctx *app.Context, wgSvc app.WireGuard, req *pb.UpdateWireGuardRequest) (*pb.UpdateWireGuardResponse, error) {
	log := ctx.Logger().WithField("op", "UpdatePeer")

	log.Debugf("update peer, peer_config: %+v", req.GetWireguardConfig().GetPeers())

	// 主链路：先更新 adjs（保证后续 wg 内部的预连接/清理逻辑使用最新拓扑）
	if err := updateAdjsFirst(log, wgSvc, req); err != nil {
		return nil, err
	}

	applyPeerOps(log, req.GetWireguardConfig().GetPeers(), "update peer", func(peer *pb.WireGuardPeerConfig) error {
		return wgSvc.UpdatePeer(&defs.WireGuardPeerConfig{WireGuardPeerConfig: peer})
	})

	log.Infof("update peer done")

	return &pb.UpdateWireGuardResponse{Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "success"}}, nil
}

func PatchPeers(ctx *app.Context, wgSvc app.WireGuard, req *pb.UpdateWireGuardRequest) (*pb.UpdateWireGuardResponse, error) {
	log := ctx.Logger().WithField("op", "PatchPeers")

	log.Debugf("patch peers, peer_config: %+v", req.GetWireguardConfig().GetPeers())

	// 主链路：先更新 adjs（保证后续 wg 内部的预连接/清理逻辑使用最新拓扑）
	if err := updateAdjsFirst(log, wgSvc, req); err != nil {
		return nil, err
	}

	wgCfg := &defs.WireGuardConfig{WireGuardConfig: req.GetWireguardConfig()}

	diffResp, err := wgSvc.PatchPeers(wgCfg.GetParsedPeers())
	if err != nil {
		log.WithError(err).Errorf("patch peers failed")
		return nil, err
	}

	log.Debugf("patch peers done, add_peers: %+v, remove_peers: %+v",
		lo.Map(diffResp.AddPeers, func(item *defs.WireGuardPeerConfig, _ int) string { return item.GetClientId() }),
		lo.Map(diffResp.RemovePeers, func(item *defs.WireGuardPeerConfig, _ int) string { return item.GetClientId() }))

	return &pb.UpdateWireGuardResponse{Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "success"}}, nil
}

func updateAdjsFirst(log *logrus.Entry, wgSvc app.WireGuard, req *pb.UpdateWireGuardRequest) error {
	if req == nil || req.GetWireguardConfig() == nil {
		return nil
	}
	if err := wgSvc.UpdateAdjs(req.GetWireguardConfig().GetAdjs()); err != nil {
		log.WithError(err).Errorf("update adjs failed, adjs: %+v", req.GetWireguardConfig().GetAdjs())
		return err
	}
	return nil
}

func applyPeerOps(log *logrus.Entry, peers []*pb.WireGuardPeerConfig, op string, fn func(peer *pb.WireGuardPeerConfig) error) {
	for _, peer := range peers {
		if peer == nil {
			continue
		}
		if err := fn(peer); err != nil {
			log.WithError(err).Errorf("%s failed", op)
			continue
		}
	}
}
