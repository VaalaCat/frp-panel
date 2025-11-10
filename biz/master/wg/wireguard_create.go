package wg

import (
	"errors"
	"net/netip"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
	"github.com/VaalaCat/frp-panel/services/rpc"
	wgsvc "github.com/VaalaCat/frp-panel/services/wg"
	"github.com/VaalaCat/frp-panel/utils"
)

// Create/Update/Get/List WireGuard 基于 pb.WireGuardConfig
// 将 pb 映射到 models.WireGuard + models.Endpoint 列表

func CreateWireGuard(ctx *app.Context, req *pb.CreateWireGuardRequest) (*pb.CreateWireGuardResponse, error) {
	log := ctx.Logger().WithField("op", "CreateWireGuard")

	userInfo := common.GetUserInfo(ctx)
	if !userInfo.Valid() {
		return nil, errors.New("invalid user")
	}
	cfg := req.GetWireguardConfig()
	if cfg == nil || len(cfg.GetClientId()) == 0 || len(cfg.GetInterfaceName()) == 0 || len(cfg.GetLocalAddress()) == 0 {
		return nil, errors.New("invalid wireguard config")
	}

	ips, err := dao.NewQuery(ctx).GetWireGuardLocalAddressesByNetworkID(userInfo, uint(cfg.GetNetworkId()))
	if err != nil {
		log.WithError(err).Errorf("get wireguard local addresses by network id failed")
		return nil, err
	}

	network, err := dao.NewQuery(ctx).GetNetworkByID(userInfo, uint(cfg.GetNetworkId()))
	if err != nil {
		log.WithError(err).Errorf("get network by id failed")
		return nil, err
	}

	newIpStr, err := utils.AllocateIP(network.CIDR, ips, cfg.GetLocalAddress())
	if err != nil {
		log.WithError(err).Errorf("allocate ip failed")
		return nil, err
	}

	newIp, err := netip.ParseAddr(newIpStr)
	if err != nil {
		log.WithError(err).Errorf("parse ip failed")
		return nil, err
	}

	networkCidr, err := netip.ParsePrefix(network.CIDR)
	if err != nil {
		log.WithError(err).Errorf("parse network cidr failed")
		return nil, err
	}

	newIpCidr := netip.PrefixFrom(newIp, networkCidr.Bits())

	keys := wgsvc.GenerateKeys()

	wgModel := &models.WireGuard{}
	wgModel.FromPB(cfg)
	wgModel.UserId = uint32(userInfo.GetUserID())
	wgModel.TenantId = uint32(userInfo.GetTenantID())
	wgModel.PrivateKey = keys.PrivateKeyBase64
	wgModel.LocalAddress = newIpCidr.String()

	log.Debugf("create wireguard with config: %+v", wgModel)

	if err := dao.NewQuery(ctx).CreateWireGuard(userInfo, wgModel); err != nil {
		return nil, err
	}

	// 处理端点：优先复用已存在的 endpoint（通过 id），否则创建新端点并绑定到该 WireGuard
	for _, ep := range cfg.GetAdvertisedEndpoints() {
		if ep == nil {
			continue
		}
		if ep.GetId() > 0 {
			// 复用现有 endpoint，要求归属同一 client
			exist, err := dao.NewQuery(ctx).GetEndpointByID(userInfo, uint(ep.GetId()))
			if err != nil {
				return nil, err
			}
			if exist.ClientID != cfg.GetClientId() {
				return nil, errors.New("endpoint client mismatch")
			}
			exist.WireGuardID = wgModel.ID

			if err := dao.NewQuery(ctx).UpdateEndpoint(userInfo, uint(exist.ID), exist.EndpointEntity); err != nil {
				return nil, err
			}
		} else {
			// 创建并绑定新端点
			newEp := &models.Endpoint{}
			newEp.FromPB(ep)
			newEp.ClientID = cfg.GetClientId()
			newEp.WireGuardID = wgModel.ID
			if err := dao.NewQuery(ctx).CreateEndpoint(userInfo, newEp.EndpointEntity); err != nil {
				return nil, err
			}
		}
	}

	go func() {
		peers, err := dao.NewQuery(ctx).GetWireGuardsByNetworkID(userInfo, uint(cfg.GetNetworkId()))
		if err != nil {
			log.WithError(err).Errorf("get wireguards by network id failed")
			return
		}
		links, err := dao.NewQuery(ctx).ListWireGuardLinksByNetwork(userInfo, uint(cfg.GetNetworkId()))
		if err != nil {
			log.WithError(err).Errorf("get wireguard links by network id failed")
			return
		}

		peerConfigs, err := wgsvc.PlanAllowedIPs(
			peers,
			links,
			wgsvc.DefaultRoutingPolicy(
				wgsvc.NewACL().LoadFromPB(network.ACL.Data),
				ctx.GetApp().GetNetworkTopologyCache(),
				ctx.GetApp().GetClientsManager(),
			))
		if err != nil {
			log.WithError(err).Errorf("build peer configs for network failed")
			return
		}

		for _, peer := range peers {
			if peer.ClientID == cfg.GetClientId() {
				if err := emitCreateWireGuardEventToClient(ctx, peer, peerConfigs[peer.ID]); err != nil {
					log.WithError(err).Errorf("update config to client failed")
				}
				continue
			}

			if err := emitPatchWireGuardEventToClient(ctx, peer, peerConfigs[peer.ID]); err != nil {
				log.WithError(err).Errorf("add wireguard event send to client error")
				continue
			}

			log.Debugf("update config to client success, client id: [%s], wireguard interface: [%s]", peer.ClientID, peer.Name)
		}
	}()

	return &pb.CreateWireGuardResponse{Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "success"}, WireguardConfig: cfg}, nil
}

func emitCreateWireGuardEventToClient(ctx *app.Context, peer *models.WireGuard, peerConfigs []*pb.WireGuardPeerConfig) error {
	log := ctx.Logger().WithField("op", "updateConfigToClient")
	userInfo := common.GetUserInfo(ctx)
	if !userInfo.Valid() {
		return errors.New("invalid user")
	}

	cfg := peer.ToPB()
	cfg.Peers = peerConfigs

	resp := &pb.CreateWireGuardResponse{}

	req := &pb.CreateWireGuardRequest{
		WireguardConfig: cfg,
	}

	err := rpc.CallClientWrapper(ctx, peer.ClientID, pb.Event_EVENT_CREATE_WIREGUARD, req, resp)
	if err != nil {
		log.WithError(err).Errorf("create wireguard event send to client error")
		return err
	}

	log.Infof("create wireguard event send to client success, client id: [%s], wireguard interface: [%s]",
		peer.ClientID, peer.Name)
	return nil
}

func emitPatchWireGuardEventToClient(ctx *app.Context, peer *models.WireGuard, peerConfigs []*pb.WireGuardPeerConfig) error {
	log := ctx.Logger().WithField("op", "patchWireGuardToClient")
	userInfo := common.GetUserInfo(ctx)
	if !userInfo.Valid() {
		return errors.New("invalid user")
	}

	cfg := peer.ToPB()
	cfg.Peers = peerConfigs

	resp := &pb.UpdateWireGuardResponse{}
	req := &pb.UpdateWireGuardRequest{
		WireguardConfig: cfg,
		UpdateType:      pb.UpdateWireGuardRequest_UPDATE_TYPE_PATCH_PEERS.Enum(),
	}

	err := rpc.CallClientWrapper(ctx, peer.ClientID, pb.Event_EVENT_UPDATE_WIREGUARD, req, resp)
	if err != nil {
		log.WithError(err).Errorf("add wireguard event send to client error")
		return err
	}

	log.Infof("add wireguard event send to client success, client id: [%s], wireguard interface: [%s]",
		peer.ClientID, peer.Name)
	return nil
}
