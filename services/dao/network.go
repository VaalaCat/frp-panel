package dao

import (
	"fmt"

	"github.com/VaalaCat/frp-panel/models"
	"gorm.io/gorm"
)

func (q *queryImpl) CreateNetwork(userInfo models.UserInfo, network *models.NetworkEntity) error {
	if network == nil {
		return fmt.Errorf("invalid network entity")
	}
	if len(network.Name) == 0 || len(network.CIDR) == 0 {
		return fmt.Errorf("invalid network name or cidr")
	}
	// scope
	network.UserId = uint32(userInfo.GetUserID())
	network.TenantId = uint32(userInfo.GetTenantID())

	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	return db.Create(&models.Network{NetworkEntity: network}).Error
}

func (q *queryImpl) UpdateNetwork(userInfo models.UserInfo, id uint, network *models.NetworkEntity) error {
	if id == 0 || network == nil {
		return fmt.Errorf("invalid network id or entity")
	}
	// scope
	network.UserId = uint32(userInfo.GetUserID())
	network.TenantId = uint32(userInfo.GetTenantID())

	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	return db.Where(&models.Network{
		Model: gorm.Model{ID: id},
		NetworkEntity: &models.NetworkEntity{
			UserId:   uint32(userInfo.GetUserID()),
			TenantId: uint32(userInfo.GetTenantID()),
		},
	}).Save(&models.Network{
		Model:         gorm.Model{ID: id},
		NetworkEntity: network,
	}).Error
}

func (q *queryImpl) DeleteNetwork(userInfo models.UserInfo, id uint) error {
	if id == 0 {
		return fmt.Errorf("invalid network id")
	}
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	return db.Unscoped().Where(&models.Network{
		Model: gorm.Model{ID: id},
		NetworkEntity: &models.NetworkEntity{
			UserId:   uint32(userInfo.GetUserID()),
			TenantId: uint32(userInfo.GetTenantID()),
		},
	}).Delete(&models.Network{}).Error
}

func (q *queryImpl) GetNetworkByID(userInfo models.UserInfo, id uint) (*models.Network, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid network id")
	}
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	var n models.Network
	if err := db.Where(&models.Network{
		Model: gorm.Model{ID: id},
		NetworkEntity: &models.NetworkEntity{
			UserId:   uint32(userInfo.GetUserID()),
			TenantId: uint32(userInfo.GetTenantID()),
		},
	}).First(&n).Error; err != nil {
		return nil, err
	}
	return &n, nil
}

func (q *queryImpl) ListNetworks(userInfo models.UserInfo, page, pageSize int) ([]*models.Network, error) {
	if page < 1 || pageSize < 1 {
		return nil, fmt.Errorf("invalid page or page size")
	}
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	var list []*models.Network
	offset := (page - 1) * pageSize
	if err := db.Where(&models.Network{NetworkEntity: &models.NetworkEntity{
		UserId:   uint32(userInfo.GetUserID()),
		TenantId: uint32(userInfo.GetTenantID()),
	}}).Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (q *queryImpl) ListNetworksWithKeyword(userInfo models.UserInfo, page, pageSize int, keyword string) ([]*models.Network, error) {
	if page < 1 || pageSize < 1 || len(keyword) == 0 {
		return nil, fmt.Errorf("invalid page or page size or keyword")
	}
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	var list []*models.Network
	offset := (page - 1) * pageSize
	if err := db.Where(&models.Network{NetworkEntity: &models.NetworkEntity{
		UserId:   uint32(userInfo.GetUserID()),
		TenantId: uint32(userInfo.GetTenantID()),
	}}).Where("name like ? OR cidr like ?", "%"+keyword+"%", "%"+keyword+"%").Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (q *queryImpl) CountNetworks(userInfo models.UserInfo) (int64, error) {
	var count int64
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	if err := db.Model(&models.Network{}).Where(&models.Network{NetworkEntity: &models.NetworkEntity{
		UserId:   uint32(userInfo.GetUserID()),
		TenantId: uint32(userInfo.GetTenantID()),
	}}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (q *queryImpl) CountNetworksWithKeyword(userInfo models.UserInfo, keyword string) (int64, error) {
	var count int64
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	if err := db.Model(&models.Network{}).Where(&models.Network{NetworkEntity: &models.NetworkEntity{
		UserId:   uint32(userInfo.GetUserID()),
		TenantId: uint32(userInfo.GetTenantID()),
	}}).Where("name like ? OR cidr like ?", "%"+keyword+"%", "%"+keyword+"%").Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
