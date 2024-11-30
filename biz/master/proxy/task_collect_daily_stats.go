package proxy

import (
	"github.com/VaalaCat/frp-panel/dao"
	"github.com/VaalaCat/frp-panel/models"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

func CollectDailyStats() error {
	tx := models.GetDBManager().GetDefaultDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	proxies, err := dao.AdminGetAllProxies(tx)
	if err != nil {
		logrus.WithError(err).Error("CollectDailyStats cannot get proxies")
		return err
	}

	proxyDailyStats := lo.Map(proxies, func(item *models.ProxyEntity, _ int) *models.HistoryProxyStats {
		return &models.HistoryProxyStats{
			ProxyID:    item.ProxyID,
			ServerID:   item.ServerID,
			ClientID:   item.ClientID,
			Name:       item.Name,
			Type:       item.Type,
			UserID:     item.UserID,
			TenantID:   item.TenantID,
			TrafficIn:  item.HistoryTrafficIn,
			TrafficOut: item.HistoryTrafficOut,
		}
	})

	if err := dao.AdminMSaveTodyStats(tx, proxyDailyStats); err != nil {
		logrus.WithError(err).Error("CollectDailyStats cannot save stats")
		return err
	}

	logrus.Infof("CollectDailyStats success")

	return nil
}
