package proxy

import (
	"errors"

	"github.com/VaalaCat/frp-panel/biz/master/client"
	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
	"github.com/VaalaCat/frp-panel/utils/logger"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

func convertProxyStatsList(proxyList []*models.ProxyStatsEntity) []*pb.ProxyInfo {
	return lo.Map(proxyList, func(item *models.ProxyStatsEntity, index int) *pb.ProxyInfo {
		return &pb.ProxyInfo{
			Name:              lo.ToPtr(item.Name),
			Type:              lo.ToPtr(item.Type),
			ClientId:          lo.ToPtr(item.ClientID),
			ServerId:          lo.ToPtr(item.ServerID),
			TodayTrafficIn:    lo.ToPtr(item.TodayTrafficIn),
			TodayTrafficOut:   lo.ToPtr(item.TodayTrafficOut),
			HistoryTrafficIn:  lo.ToPtr(item.HistoryTrafficIn),
			HistoryTrafficOut: lo.ToPtr(item.HistoryTrafficOut),
		}
	})
}

// getClientWithMakeShadow
// 1. 检查是否有已连接该服务端的客户端
// 2. 检查是否有Shadow客户端
// 3. 如果没有，则新建Shadow客户端和子客户端
func getClientWithMakeShadow(c *app.Context, clientID, serverID string) (*models.ClientEntity, error) {
	userInfo := common.GetUserInfo(c)
	clientEntity, err := dao.NewQuery(c).GetClientByFilter(userInfo, &models.ClientEntity{OriginClientID: clientID, ServerID: serverID}, lo.ToPtr(false))
	if errors.Is(err, gorm.ErrRecordNotFound) {
		clientEntity, err = dao.NewQuery(c).GetClientByFilter(userInfo, &models.ClientEntity{ClientID: clientID}, nil)
		if err != nil {
			logger.Logger(c).WithError(err).Errorf("cannot get client, id: [%s]", clientID)
			return nil, err
		}
		if (!clientEntity.IsShadow || len(clientEntity.ConfigContent) != 0) && len(clientEntity.OriginClientID) == 0 {
			// 没shadow过，需要shadow
			_, err = client.MakeClientShadowed(c, serverID, clientEntity)
			if err != nil {
				logger.Logger(c).WithError(err).Errorf("cannot make client shadow, id: [%s]", clientID)
				return nil, err
			}
		}
		// shadow过，但没找到子客户端，需要新建
		clientEntity, err = client.ChildClientForServer(c, serverID, clientEntity)
		if err != nil {
			logger.Logger(c).WithError(err).Errorf("cannot create child client, id: [%s]", clientID)
			return nil, err
		}
	}
	// 有任何失败，返回
	if err != nil {
		logger.Logger(c).WithError(err).Errorf("cannot get client, id: [%s]", clientID)
		return nil, err
	}

	return clientEntity, nil
}
