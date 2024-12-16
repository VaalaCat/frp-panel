package proxy

import (
	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/samber/lo"
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
