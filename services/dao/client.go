package dao

import (
	"fmt"
	"time"

	"github.com/VaalaCat/frp-panel/models"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type ClientQuery interface {
	ValidateClientSecret(clientID, clientSecret string) (*models.ClientEntity, error)
	AdminGetClientByClientID(clientID string) (*models.Client, error)
	GetClientByClientID(userInfo models.UserInfo, clientID string) (*models.Client, error)
	GetClientsByClientIDs(userInfo models.UserInfo, clientIDs []string) ([]*models.Client, error)
	GetClientByFilter(userInfo models.UserInfo, client *models.ClientEntity, shadow *bool) (*models.ClientEntity, error)
	GetClientByOriginClientID(originClientID string) (*models.ClientEntity, error)
	ListClients(userInfo models.UserInfo, page, pageSize int) ([]*models.ClientEntity, error)
	ListClientsWithKeyword(userInfo models.UserInfo, page, pageSize int, keyword string) ([]*models.ClientEntity, error)
	GetAllClients(userInfo models.UserInfo) ([]*models.ClientEntity, error)
	CountClients(userInfo models.UserInfo) (int64, error)
	CountClientsWithKeyword(userInfo models.UserInfo, keyword string) (int64, error)
	CountConfiguredClients(userInfo models.UserInfo) (int64, error)
	CountClientsInShadow(userInfo models.UserInfo, clientID string) (int64, error)
	GetClientIDsInShadowByClientID(userInfo models.UserInfo, clientID string) ([]string, error)
	AdminGetClientIDsInShadowByClientID(clientID string) ([]string, error)
}

type ClientMutation interface {
	CreateClient(userInfo models.UserInfo, client *models.ClientEntity) error
	DeleteClient(userInfo models.UserInfo, clientID string) error
	UpdateClient(userInfo models.UserInfo, client *models.ClientEntity) error
	AdminUpdateClientLastSeen(clientID string) error
}

type clientQuery struct{ *queryImpl }
type clientMutation struct{ *mutationImpl }

func newClientQuery(base *queryImpl) ClientQuery          { return &clientQuery{base} }
func newClientMutation(base *mutationImpl) ClientMutation { return &clientMutation{base} }

func (q *clientQuery) ValidateClientSecret(clientID, clientSecret string) (*models.ClientEntity, error) {
	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("invalid client id or client secret")
	}
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	c := &models.Client{}
	err := db.
		Where(&models.Client{ClientEntity: &models.ClientEntity{
			ClientID: clientID,
		}}).
		First(c).Error
	if err != nil {
		return nil, err
	}
	if c.ConnectSecret != clientSecret {
		return nil, fmt.Errorf("invalid client secret")
	}
	return c.ClientEntity, nil
}

func (q *clientQuery) AdminGetClientByClientID(clientID string) (*models.Client, error) {
	if clientID == "" {
		return nil, fmt.Errorf("invalid client id")
	}
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	c := &models.Client{}
	err := db.
		Where(&models.Client{ClientEntity: &models.ClientEntity{
			ClientID: clientID,
		}}).
		First(c).Error
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (q *clientQuery) GetClientByClientID(userInfo models.UserInfo, clientID string) (*models.Client, error) {
	if clientID == "" {
		return nil, fmt.Errorf("invalid client id")
	}
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	c := &models.Client{}
	err := db.
		Where(&models.Client{ClientEntity: &models.ClientEntity{
			UserID:   userInfo.GetUserID(),
			TenantID: userInfo.GetTenantID(),
			ClientID: clientID,
		}}).
		First(c).Error
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (q *clientQuery) GetClientsByClientIDs(userInfo models.UserInfo, clientIDs []string) ([]*models.Client, error) {
	if len(clientIDs) == 0 {
		return nil, fmt.Errorf("invalid client ids")
	}

	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	cs := []*models.Client{}
	err := db.Where("client_id IN ?", clientIDs).Find(&cs).Error
	if err != nil {
		return nil, err
	}

	return cs, nil
}

func (q *clientQuery) GetClientByFilter(userInfo models.UserInfo, client *models.ClientEntity, shadow *bool) (*models.ClientEntity, error) {
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	filter := &models.ClientEntity{}
	if len(client.ClientID) != 0 {
		filter.ClientID = client.ClientID
	}
	if len(client.OriginClientID) != 0 {
		filter.OriginClientID = client.OriginClientID
	}
	if len(client.ConnectSecret) != 0 {
		filter.ConnectSecret = client.ConnectSecret
	}
	if len(client.ServerID) != 0 {
		filter.ServerID = client.ServerID
	}
	if shadow != nil {
		filter.IsShadow = *shadow
	}
	c := &models.Client{}

	err := db.
		Where(&models.Client{ClientEntity: &models.ClientEntity{
			UserID:   userInfo.GetUserID(),
			TenantID: userInfo.GetTenantID(),
		}}).
		Where(filter).
		First(c).Error
	if err != nil {
		return nil, err
	}
	return c.ClientEntity, nil
}

func (q *clientQuery) GetClientByOriginClientID(originClientID string) (*models.ClientEntity, error) {
	if originClientID == "" {
		return nil, fmt.Errorf("invalid origin client id")
	}
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	c := &models.Client{}
	err := db.
		Where(&models.Client{ClientEntity: &models.ClientEntity{
			OriginClientID: originClientID,
		}}).
		First(c).Error
	if err != nil {
		return nil, err
	}
	return c.ClientEntity, nil
}

func (m *clientMutation) CreateClient(userInfo models.UserInfo, client *models.ClientEntity) error {
	client.UserID = userInfo.GetUserID()
	client.TenantID = userInfo.GetTenantID()
	c := &models.Client{
		ClientEntity: client,
	}
	db := m.ctx.GetApp().GetDBManager().GetDefaultDB()
	return db.Create(c).Error
}

func (m *clientMutation) DeleteClient(userInfo models.UserInfo, clientID string) error {
	if clientID == "" {
		return fmt.Errorf("invalid client id")
	}
	db := m.ctx.GetApp().GetDBManager().GetDefaultDB()
	return db.Unscoped().Where(&models.Client{
		ClientEntity: &models.ClientEntity{
			ClientID: clientID,
			UserID:   userInfo.GetUserID(),
			TenantID: userInfo.GetTenantID(),
		},
	}).Or(&models.Client{
		ClientEntity: &models.ClientEntity{
			OriginClientID: clientID,
			UserID:         userInfo.GetUserID(),
			TenantID:       userInfo.GetTenantID(),
		},
	}).Delete(&models.Client{}).Error
}

func (m *clientMutation) UpdateClient(userInfo models.UserInfo, client *models.ClientEntity) error {
	c := &models.Client{
		ClientEntity: client,
	}
	db := m.ctx.GetApp().GetDBManager().GetDefaultDB()
	return db.Where(&models.Client{
		ClientEntity: &models.ClientEntity{
			UserID:   userInfo.GetUserID(),
			TenantID: userInfo.GetTenantID(),
		},
	}).Save(c).Error
}

func (q *clientQuery) ListClients(userInfo models.UserInfo, page, pageSize int) ([]*models.ClientEntity, error) {
	if page < 1 || pageSize < 1 {
		return nil, fmt.Errorf("invalid page or page size")
	}

	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	offset := (page - 1) * pageSize

	var clients []*models.Client
	err := db.Where(&models.Client{
		ClientEntity: &models.ClientEntity{
			UserID:   userInfo.GetUserID(),
			TenantID: userInfo.GetTenantID(),
		},
	}).
		Where(
			db.Where(
				normalClientFilter(db),
			),
		).Offset(offset).Limit(pageSize).Find(&clients).Error
	if err != nil {
		return nil, err
	}

	return lo.Map(clients, func(c *models.Client, _ int) *models.ClientEntity {
		return c.ClientEntity
	}), nil
}

func (q *clientQuery) ListClientsWithKeyword(userInfo models.UserInfo, page, pageSize int, keyword string) ([]*models.ClientEntity, error) {
	// 只获取没shadow且config有东西
	// 或isShadow的client
	if page < 1 || pageSize < 1 || len(keyword) == 0 {
		return nil, fmt.Errorf("invalid page or page size or keyword")
	}

	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	offset := (page - 1) * pageSize

	var clients []*models.Client
	err := db.Where("client_id like ?", "%"+keyword+"%").
		Where(&models.Client{ClientEntity: &models.ClientEntity{
			UserID:   userInfo.GetUserID(),
			TenantID: userInfo.GetTenantID(),
		}}).
		Where(normalClientFilter(db)).
		Offset(offset).Limit(pageSize).Find(&clients).Error
	if err != nil {
		return nil, err
	}

	return lo.Map(clients, func(c *models.Client, _ int) *models.ClientEntity {
		return c.ClientEntity
	}), nil
}

func (q *clientQuery) GetAllClients(userInfo models.UserInfo) ([]*models.ClientEntity, error) {
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	var clients []*models.Client
	err := db.Where(&models.Client{
		ClientEntity: &models.ClientEntity{
			UserID:   userInfo.GetUserID(),
			TenantID: userInfo.GetTenantID(),
		},
	}).Find(&clients).Error
	if err != nil {
		return nil, err
	}

	return lo.Map(clients, func(c *models.Client, _ int) *models.ClientEntity {
		return c.ClientEntity
	}), nil
}

func (q *clientQuery) CountClients(userInfo models.UserInfo) (int64, error) {
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	var count int64
	err := db.Model(&models.Client{}).Where(&models.Client{
		ClientEntity: &models.ClientEntity{
			UserID:   userInfo.GetUserID(),
			TenantID: userInfo.GetTenantID(),
		},
	}).
		Where(normalClientFilter(db)).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (q *clientQuery) CountClientsWithKeyword(userInfo models.UserInfo, keyword string) (int64, error) {
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	var count int64
	err := db.Model(&models.Client{}).Where(&models.Client{
		ClientEntity: &models.ClientEntity{
			UserID:   userInfo.GetUserID(),
			TenantID: userInfo.GetTenantID(),
		},
	}).
		Where(normalClientFilter(db)).Where("client_id like ?", "%"+keyword+"%").Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (q *clientQuery) CountConfiguredClients(userInfo models.UserInfo) (int64, error) {
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	var count int64
	err := db.Model(&models.Client{}).
		Where(&models.Client{
			ClientEntity: &models.ClientEntity{
				UserID:   userInfo.GetUserID(),
				TenantID: userInfo.GetTenantID(),
			}}).
		Where(normalClientFilter(db)).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (q *clientQuery) CountClientsInShadow(userInfo models.UserInfo, clientID string) (int64, error) {
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	var count int64
	err := db.Model(&models.Client{}).
		Where(&models.Client{
			ClientEntity: &models.ClientEntity{
				UserID:         userInfo.GetUserID(),
				TenantID:       userInfo.GetTenantID(),
				OriginClientID: clientID,
			}}).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (q *clientQuery) GetClientIDsInShadowByClientID(userInfo models.UserInfo, clientID string) ([]string, error) {
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	var clients []*models.Client
	err := db.Where(&models.Client{
		ClientEntity: &models.ClientEntity{
			UserID:         userInfo.GetUserID(),
			TenantID:       userInfo.GetTenantID(),
			OriginClientID: clientID,
		}}).Find(&clients).Error
	if err != nil {
		return nil, err
	}
	return lo.Map(clients, func(c *models.Client, _ int) string {
		return c.ClientID
	}), nil
}

func (q *clientQuery) AdminGetClientIDsInShadowByClientID(clientID string) ([]string, error) {
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	var clients []*models.Client
	err := db.Where(&models.Client{
		ClientEntity: &models.ClientEntity{
			OriginClientID: clientID,
		}}).Find(&clients).Error
	if err != nil {
		return nil, err
	}
	return lo.Map(clients, func(c *models.Client, _ int) string {
		return c.ClientID
	}), nil
}

func (m *clientMutation) AdminUpdateClientLastSeen(clientID string) error {
	db := m.ctx.GetApp().GetDBManager().GetDefaultDB()
	return db.Model(&models.Client{
		ClientEntity: &models.ClientEntity{
			ClientID: clientID,
		}}).Update("last_seen_at", time.Now()).Error
}

func normalClientFilter(db *gorm.DB) *gorm.DB {
	// 1. 没shadow过的老client
	// 2. shadow过的shadow client
	// 3. 非临时节点
	return db.Where(
		db.Where("origin_client_id is NULL").
			Or("is_shadow = ?", true).
			Or("LENGTH(origin_client_id) = ?", 0),
	).
		Where(
			db.Where(
				db.Where("ephemeral is NULL").
					Or("ephemeral = ?", false),
			).Or(
				db.Where("ephemeral = ?", true).
					Where("last_seen_at > ?", time.Now().Add(-5*time.Minute)),
			),
		)
}
