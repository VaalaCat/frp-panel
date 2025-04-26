package proxy

import (
	"github.com/VaalaCat/frp-panel/app"
	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/dao"
	"github.com/VaalaCat/frp-panel/logger"
	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/samber/lo"
)

func ListProxyConfigs(ctx *app.Context, req *pb.ListProxyConfigsRequest) (*pb.ListProxyConfigsResponse, error) {
	logger.Logger(ctx).Infof("list proxy configs, req: [%+v]", req)

	var (
		userInfo = common.GetUserInfo(ctx)
	)

	if !userInfo.Valid() {
		return &pb.ListProxyConfigsResponse{
			Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "invalid user"},
		}, nil
	}

	var (
		page         = int(req.GetPage())
		pageSize     = int(req.GetPageSize())
		keyword      = req.GetKeyword()
		clientID     = req.GetClientId()
		serverID     = req.GetServerId()
		hasKeyword   = len(keyword) > 0
		hasClientID  = len(clientID) > 0
		hasServerID  = len(serverID) > 0
		proxyConfigs []*models.ProxyConfig
		err          error
		proxyCounts  int64
		filter       = &models.ProxyConfigEntity{}
	)

	if hasClientID {
		filter.OriginClientID = clientID
	}
	if hasServerID {
		filter.ServerID = serverID
	}

	if hasKeyword {
		proxyConfigs, err = dao.NewQuery(ctx).ListProxyConfigsWithFiltersAndKeyword(userInfo, page, pageSize, filter, keyword)
	} else {
		proxyConfigs, err = dao.NewQuery(ctx).ListProxyConfigsWithFilters(userInfo, page, pageSize, filter)
	}

	if err != nil {
		return nil, err
	}

	if hasKeyword {
		proxyCounts, err = dao.NewQuery(ctx).CountProxyConfigsWithFiltersAndKeyword(userInfo, filter, keyword)
	} else {
		proxyCounts, err = dao.NewQuery(ctx).CountProxyConfigsWithFilters(userInfo, filter)
	}

	if err != nil {
		return nil, err
	}

	respProxyConfigs := lo.Map(proxyConfigs, func(item *models.ProxyConfig, _ int) *pb.ProxyConfig {
		return &pb.ProxyConfig{
			Id:             lo.ToPtr(uint32(item.ID)),
			Name:           lo.ToPtr(item.Name),
			Type:           lo.ToPtr(item.Type),
			ClientId:       lo.ToPtr(item.ClientID),
			ServerId:       lo.ToPtr(item.ServerID),
			Config:         lo.ToPtr(string(item.Content)),
			OriginClientId: lo.ToPtr(item.OriginClientID),
		}
	})

	return &pb.ListProxyConfigsResponse{
		Status:       &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS, Message: "success"},
		ProxyConfigs: respProxyConfigs,
		Total:        lo.ToPtr(int32(proxyCounts)),
	}, nil
}
