package proxy

import (
	"context"

	"github.com/VaalaCat/frp-panel/dao"
	"github.com/VaalaCat/frp-panel/logger"
	"github.com/VaalaCat/frp-panel/models"
	"github.com/samber/lo"
)

func CollectDailyStats() error {
	ctx := context.Background()

	tx := models.GetDBManager().GetDefaultDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	proxies, err := dao.AdminGetAllProxyStats(tx)
	if err != nil {
		logger.Logger(context.Background()).WithError(err).Error("CollectDailyStats cannot get proxies")
		return err
	}

	proxyDailyStats := lo.Map(proxies, func(item *models.ProxyStatsEntity, _ int) *models.HistoryProxyStats {
		return &models.HistoryProxyStats{
			ProxyID:        item.ProxyID,
			ServerID:       item.ServerID,
			ClientID:       item.ClientID,
			OriginClientID: item.OriginClientID,
			Name:           item.Name,
			Type:           item.Type,
			UserID:         item.UserID,
			TenantID:       item.TenantID,
			TrafficIn:      item.HistoryTrafficIn,
			TrafficOut:     item.HistoryTrafficOut,
		}
	})

	if err := dao.AdminMSaveTodyStats(tx, proxyDailyStats); err != nil {
		logger.Logger(context.Background()).WithError(err).Error("CollectDailyStats cannot save stats")
		return err
	}

	logger.Logger(ctx).Infof("CollectDailyStats success")

	return nil
}
