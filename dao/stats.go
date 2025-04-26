package dao

import (
	"fmt"

	"github.com/VaalaCat/frp-panel/models"
	"gorm.io/gorm"
)

const (
	MSetBatchSize = 100
)

func (q *queryImpl) AdminSaveTodyStats(s *models.HistoryProxyStats) error {
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	return db.Save(s).Error
}

func (q *queryImpl) AdminMSaveTodyStats(tx *gorm.DB, s []*models.HistoryProxyStats) error {
	if len(s) == 0 {
		return nil
	}

	if err := tx.CreateInBatches(s, MSetBatchSize).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("AdminMSaveTodyStats failed to save history proxy stats: %v", err)
	}
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("AdminMSaveTodyStats failed to commit transaction: %v", err)
	}
	return nil
}

func (q *queryImpl) GetHistoryStatsByProxyID(userInfo models.UserInfo, proxyID int) ([]*models.HistoryProxyStats, error) {
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	var stats []*models.HistoryProxyStats
	err := db.Where(&models.HistoryProxyStats{
		ProxyID:  proxyID,
		UserID:   userInfo.GetUserID(),
		TenantID: userInfo.GetTenantID(),
	}).Find(&stats).Error
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (q *queryImpl) GetHistoryStatsByClientID(userInfo models.UserInfo, clientID string) ([]*models.HistoryProxyStats, error) {
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	var stats []*models.HistoryProxyStats
	err := db.Where(&models.HistoryProxyStats{
		ClientID: clientID,
		UserID:   userInfo.GetUserID(),
		TenantID: userInfo.GetTenantID(),
	}).Find(&stats).Error
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (q *queryImpl) GetHistoryStatsByServerID(userInfo models.UserInfo, serverID string) ([]*models.HistoryProxyStats, error) {
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	var stats []*models.HistoryProxyStats
	err := db.Where(&models.HistoryProxyStats{
		ServerID: serverID,
		UserID:   userInfo.GetUserID(),
		TenantID: userInfo.GetTenantID(),
	}).Find(&stats).Error
	if err != nil {
		return nil, err
	}
	return stats, nil
}
