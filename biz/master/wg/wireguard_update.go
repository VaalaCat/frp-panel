package wg

import (
	"errors"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
	"github.com/samber/lo"
)

func UpdateWireGuard(ctx *app.Context, req *pb.UpdateWireGuardRequest) (*pb.UpdateWireGuardResponse, error) {
	userInfo := common.GetUserInfo(ctx)
	if !userInfo.Valid() {
		return nil, errors.New("invalid user")
	}
	cfg := req.GetWireguardConfig()
	if cfg == nil || cfg.GetId() == 0 || len(cfg.GetClientId()) == 0 || len(cfg.GetInterfaceName()) == 0 || len(cfg.GetPrivateKey()) == 0 || len(cfg.GetLocalAddress()) == 0 {
		return nil, errors.New("invalid wireguard config")
	}
	model := &models.WireGuard{
		WireGuardEntity: &models.WireGuardEntity{
			Name:         cfg.GetInterfaceName(),
			UserId:       uint32(userInfo.GetUserID()),
			TenantId:     uint32(userInfo.GetTenantID()),
			PrivateKey:   cfg.GetPrivateKey(),
			LocalAddress: cfg.GetLocalAddress(),
			ListenPort:   cfg.GetListenPort(),
			InterfaceMtu: cfg.GetInterfaceMtu(),
			DnsServers:   models.GormArray[string](cfg.GetDnsServers()),
			ClientID:     cfg.GetClientId(),
			NetworkID:    uint(cfg.GetNetworkId()),
			Tags:         models.GormArray[string](cfg.GetTags()),
		},
	}
	if err := dao.NewQuery(ctx).UpdateWireGuard(userInfo, uint(cfg.GetId()), model); err != nil {
		return nil, err
	}

	// 端点赋值与解绑
	// 1) 读取当前绑定的端点
	currentList, err := dao.NewQuery(ctx).ListEndpointsWithFilters(userInfo, 1, 1000, "", uint(cfg.GetId()), "")
	if err != nil {
		return nil, err
	}
	currentSet := lo.SliceToMap(currentList, func(e *models.Endpoint) (uint, struct{}) { return uint(e.ID), struct{}{} })

	// 2) 处理本次配置中的端点：存在则更新并绑定；不存在则创建绑定
	newSet := map[uint]struct{}{}
	for _, ep := range cfg.GetAdvertisedEndpoints() {
		if ep == nil {
			continue
		}
		if ep.GetId() > 0 {
			exist, err := dao.NewQuery(ctx).GetEndpointByID(userInfo, uint(ep.GetId()))
			if err != nil {
				return nil, err
			}
			// 必须属于相同 client
			if exist.ClientID != cfg.GetClientId() {
				return nil, errors.New("endpoint client mismatch")
			}
			entity := &models.EndpointEntity{Host: ep.GetHost(), Port: ep.GetPort(), ClientID: exist.ClientID, WireGuardID: uint(cfg.GetId())}
			if err := dao.NewQuery(ctx).UpdateEndpoint(userInfo, uint(exist.ID), entity); err != nil {
				return nil, err
			}
			newSet[uint(exist.ID)] = struct{}{}
		} else {
			entity := &models.EndpointEntity{Host: ep.GetHost(), Port: ep.GetPort(), ClientID: cfg.GetClientId(), WireGuardID: uint(cfg.GetId())}
			if err := dao.NewQuery(ctx).CreateEndpoint(userInfo, entity); err != nil {
				return nil, err
			}
			// 无法获取新建 id，这里不加入 newSet；不影响后续解绑逻辑（仅解绑 current - new）
		}
	}

	// 3) 解绑本次未包含的端点（将 WireGuardID 置 0）
	for id := range currentSet {
		if _, ok := newSet[id]; ok {
			continue
		}
		exist, err := dao.NewQuery(ctx).GetEndpointByID(userInfo, id)
		if err != nil {
			return nil, err
		}
		entity := &models.EndpointEntity{Host: exist.Host, Port: exist.Port, ClientID: exist.ClientID, WireGuardID: 0}
		if err := dao.NewQuery(ctx).UpdateEndpoint(userInfo, id, entity); err != nil {
			return nil, err
		}
	}
	return &pb.UpdateWireGuardResponse{Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "success"}, WireguardConfig: cfg}, nil
}
