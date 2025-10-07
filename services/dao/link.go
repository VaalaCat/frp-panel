package dao

import (
	"fmt"
	"strconv"

	"github.com/VaalaCat/frp-panel/models"
	"gorm.io/gorm"
)

func (q *queryImpl) CreateWireGuardLink(userInfo models.UserInfo, link *models.WireGuardLink) error {
	if link == nil {
		return fmt.Errorf("invalid wg link")
	}
	if link.WireGuardLinkEntity == nil {
		link.WireGuardLinkEntity = &models.WireGuardLinkEntity{}
	}
	link.UserId = uint32(userInfo.GetUserID())
	link.TenantId = uint32(userInfo.GetTenantID())
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	return db.Create(link).Error
}

func (q *queryImpl) CreateWireGuardLinks(userInfo models.UserInfo, links ...*models.WireGuardLink) error {
	if len(links) == 0 {
		return fmt.Errorf("invalid wg links")
	}
	for _, link := range links {
		link.UserId = uint32(userInfo.GetUserID())
		link.TenantId = uint32(userInfo.GetTenantID())
	}
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	return db.Create(links).Error
}

func (q *queryImpl) UpdateWireGuardLink(userInfo models.UserInfo, id uint, link *models.WireGuardLink) error {
	if id == 0 || link == nil {
		return fmt.Errorf("invalid wg link id or entity")
	}
	link.Model.ID = id
	if link.WireGuardLinkEntity == nil {
		link.WireGuardLinkEntity = &models.WireGuardLinkEntity{}
	}
	link.UserId = uint32(userInfo.GetUserID())
	link.TenantId = uint32(userInfo.GetTenantID())
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	return db.Where(&models.WireGuardLink{
		Model:               link.Model,
		WireGuardLinkEntity: &models.WireGuardLinkEntity{UserId: link.UserId, TenantId: link.TenantId},
	}).Save(link).Error
}

func (q *queryImpl) DeleteWireGuardLink(userInfo models.UserInfo, id uint) error {
	if id == 0 {
		return fmt.Errorf("invalid wg link id")
	}
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	return db.Unscoped().Where(&models.WireGuardLink{
		Model: gorm.Model{ID: id},
		WireGuardLinkEntity: &models.WireGuardLinkEntity{
			UserId:   uint32(userInfo.GetUserID()),
			TenantId: uint32(userInfo.GetTenantID()),
		}}).Delete(&models.WireGuardLink{}).Error
}

func (q *queryImpl) ListWireGuardLinksByNetwork(userInfo models.UserInfo, networkID uint) ([]*models.WireGuardLink, error) {
	if networkID == 0 {
		return nil, fmt.Errorf("invalid network id")
	}
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	var list []*models.WireGuardLink
	if err := db.Preload("ToEndpoint").Where(&models.WireGuardLink{
		WireGuardLinkEntity: &models.WireGuardLinkEntity{
			NetworkID: networkID,
			UserId:    uint32(userInfo.GetUserID()),
			TenantId:  uint32(userInfo.GetTenantID()),
		}}).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// GetWireGuardLinkByID 根据 ID 查询 Link（按租户隔离）
func (q *queryImpl) GetWireGuardLinkByID(userInfo models.UserInfo, id uint) (*models.WireGuardLink, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid wg link id")
	}
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	var m models.WireGuardLink
	if err := db.Preload("ToEndpoint").Where(&models.WireGuardLink{
		Model: gorm.Model{ID: id},
		WireGuardLinkEntity: &models.WireGuardLinkEntity{
			UserId:   uint32(userInfo.GetUserID()),
			TenantId: uint32(userInfo.GetTenantID()),
		},
	}).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (q *queryImpl) GetWireGuardLinkByClientIDs(userInfo models.UserInfo, fromClientId, toClientId uint) (*models.WireGuardLink, error) {
	if fromClientId == 0 || toClientId == 0 {
		return nil, fmt.Errorf("invalid from client id or to client id")
	}
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	var link *models.WireGuardLink
	if err := db.Preload("ToEndpoint").Where(&models.WireGuardLink{WireGuardLinkEntity: &models.WireGuardLinkEntity{
		UserId:          uint32(userInfo.GetUserID()),
		TenantId:        uint32(userInfo.GetTenantID()),
		FromWireGuardID: fromClientId,
		ToWireGuardID:   toClientId,
	}}).First(&link).Error; err != nil {
		return nil, err
	}
	return link, nil
}

// ListWireGuardLinksWithFilters 分页查询 Link，支持按 networkID 过滤与关键字（数字时匹配 from/to id）
func (q *queryImpl) ListWireGuardLinksWithFilters(userInfo models.UserInfo, page, pageSize int, networkID uint, keyword string) ([]*models.WireGuardLink, error) {
	if page < 1 || pageSize < 1 {
		return nil, fmt.Errorf("invalid page or page size")
	}
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	var list []*models.WireGuardLink
	offset := (page - 1) * pageSize

	base := db.Preload("ToEndpoint").Where(&models.WireGuardLink{WireGuardLinkEntity: &models.WireGuardLinkEntity{
		UserId:   uint32(userInfo.GetUserID()),
		TenantId: uint32(userInfo.GetTenantID()),
	}})

	if networkID > 0 {
		base = base.Where(&models.WireGuardLink{WireGuardLinkEntity: &models.WireGuardLinkEntity{NetworkID: networkID}})
	}
	if len(keyword) > 0 {
		if v, err := strconv.ParseUint(keyword, 10, 64); err == nil {
			base = base.Where("from_wire_guard_id = ? OR to_wire_guard_id = ?", uint(v), uint(v))
		}
	}

	if err := base.Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// CountWireGuardLinksWithFilters 统计分页条件下的总数
func (q *queryImpl) CountWireGuardLinksWithFilters(userInfo models.UserInfo, networkID uint, keyword string) (int64, error) {
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	var count int64

	base := db.Model(&models.WireGuardLink{}).Where(&models.WireGuardLink{WireGuardLinkEntity: &models.WireGuardLinkEntity{
		UserId:   uint32(userInfo.GetUserID()),
		TenantId: uint32(userInfo.GetTenantID()),
	}})

	if networkID > 0 {
		base = base.Where(&models.WireGuardLink{WireGuardLinkEntity: &models.WireGuardLinkEntity{NetworkID: networkID}})
	}
	if len(keyword) > 0 {
		if v, err := strconv.ParseUint(keyword, 10, 64); err == nil {
			base = base.Where("from_wire_guard_id = ? OR to_wire_guard_id = ?", uint(v), uint(v))
		}
	}

	if err := base.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (q *queryImpl) AdminListWireGuardLinksWithNetworkIDs(networkIDs []uint) ([]*models.WireGuardLink, error) {
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	var list []*models.WireGuardLink
	if err := db.Where("network_id IN ?", networkIDs).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
