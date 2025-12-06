package dao

import (
	"fmt"

	"github.com/VaalaCat/frp-panel/models"
	"gorm.io/gorm"
)

type EndpointQuery interface {
	GetEndpointByID(userInfo models.UserInfo, id uint) (*models.Endpoint, error)
	ListEndpoints(userInfo models.UserInfo, page, pageSize int) ([]*models.Endpoint, error)
	CountEndpoints(userInfo models.UserInfo) (int64, error)
	ListEndpointsWithFilters(userInfo models.UserInfo, page, pageSize int, clientID string, wireguardID uint, keyword string) ([]*models.Endpoint, error)
	CountEndpointsWithFilters(userInfo models.UserInfo, clientID string, wireguardID uint, keyword string) (int64, error)
}

type EndpointMutation interface {
	CreateEndpoint(userInfo models.UserInfo, endpoint *models.EndpointEntity) error
	UpdateEndpoint(userInfo models.UserInfo, id uint, endpoint *models.EndpointEntity) error
	DeleteEndpoint(userInfo models.UserInfo, id uint) error
}

type endpointQuery struct{ *queryImpl }
type endpointMutation struct{ *mutationImpl }

func newEndpointQuery(base *queryImpl) EndpointQuery          { return &endpointQuery{base} }
func newEndpointMutation(base *mutationImpl) EndpointMutation { return &endpointMutation{base} }

func (m *endpointMutation) CreateEndpoint(userInfo models.UserInfo, endpoint *models.EndpointEntity) error {
	if endpoint == nil {
		return fmt.Errorf("invalid endpoint entity")
	}
	if len(endpoint.Host) == 0 || endpoint.Port == 0 {
		return fmt.Errorf("invalid endpoint host or port")
	}
	// scope via parent wireguard/client
	db := m.ctx.GetApp().GetDBManager().GetDefaultDB()
	return db.Create(&models.Endpoint{EndpointEntity: endpoint}).Error
}

func (m *endpointMutation) UpdateEndpoint(userInfo models.UserInfo, id uint, endpoint *models.EndpointEntity) error {
	if id == 0 || endpoint == nil {
		return fmt.Errorf("invalid endpoint id or entity")
	}
	db := m.ctx.GetApp().GetDBManager().GetDefaultDB()
	return db.Where(&models.Endpoint{
		Model: gorm.Model{ID: id},
	}).Save(&models.Endpoint{Model: gorm.Model{ID: id}, EndpointEntity: endpoint}).Error
}

func (m *endpointMutation) DeleteEndpoint(userInfo models.UserInfo, id uint) error {
	if id == 0 {
		return fmt.Errorf("invalid endpoint id")
	}
	db := m.ctx.GetApp().GetDBManager().GetDefaultDB()
	return db.Unscoped().Where(&models.Endpoint{Model: gorm.Model{ID: id}}).Delete(&models.Endpoint{}).Error
}

func (q *endpointQuery) GetEndpointByID(userInfo models.UserInfo, id uint) (*models.Endpoint, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid endpoint id")
	}
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	var e models.Endpoint
	if err := db.Where(&models.Endpoint{Model: gorm.Model{ID: id}}).First(&e).Error; err != nil {
		return nil, err
	}
	return &e, nil
}

func (q *endpointQuery) ListEndpoints(userInfo models.UserInfo, page, pageSize int) ([]*models.Endpoint, error) {
	if page < 1 || pageSize < 1 {
		return nil, fmt.Errorf("invalid page or page size")
	}
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	var list []*models.Endpoint
	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (q *endpointQuery) CountEndpoints(userInfo models.UserInfo) (int64, error) {
	var count int64
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	if err := db.Model(&models.Endpoint{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ListEndpointsWithFilters 根据 clientID / wireguardID / keyword 过滤端点
func (q *endpointQuery) ListEndpointsWithFilters(userInfo models.UserInfo, page, pageSize int, clientID string, wireguardID uint, keyword string) ([]*models.Endpoint, error) {
	if page < 1 || pageSize < 1 {
		return nil, fmt.Errorf("invalid page or page size")
	}
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()

	// 若指定 clientID，先校验归属
	if len(clientID) > 0 {
		if _, err := newClientQuery(q.queryImpl).GetClientByClientID(userInfo, clientID); err != nil {
			return nil, err
		}
	}

	var list []*models.Endpoint
	offset := (page - 1) * pageSize
	query := db.Model(&models.Endpoint{})
	if len(clientID) > 0 {
		query = query.Where(&models.Endpoint{EndpointEntity: &models.EndpointEntity{ClientID: clientID}})
	}
	if wireguardID > 0 {
		query = query.Where(&models.Endpoint{EndpointEntity: &models.EndpointEntity{WireGuardID: wireguardID}})
	}
	if len(keyword) > 0 {
		query = query.Where("host like ?", "%"+keyword+"%")
	}
	if err := query.Offset(offset).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (q *endpointQuery) CountEndpointsWithFilters(userInfo models.UserInfo, clientID string, wireguardID uint, keyword string) (int64, error) {
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()

	if len(clientID) > 0 {
		if _, err := newClientQuery(q.queryImpl).GetClientByClientID(userInfo, clientID); err != nil {
			return 0, err
		}
	}

	var count int64
	query := db.Model(&models.Endpoint{})
	if len(clientID) > 0 {
		query = query.Where(&models.Endpoint{EndpointEntity: &models.EndpointEntity{ClientID: clientID}})
	}
	if wireguardID > 0 {
		query = query.Where(&models.Endpoint{EndpointEntity: &models.EndpointEntity{WireGuardID: wireguardID}})
	}
	if len(keyword) > 0 {
		query = query.Where("host like ?", "%"+keyword+"%")
	}
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
