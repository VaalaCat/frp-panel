package dao

import (
	"fmt"

	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/models"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type ServerQuery interface {
	GetDefaultServer() (*models.ServerEntity, error)
	ValidateServerSecret(serverID string, secret string) (*models.ServerEntity, error)
	AdminGetServerByServerID(serverID string) (*models.ServerEntity, error)
	GetServerByServerID(userInfo models.UserInfo, serverID string) (*models.ServerEntity, error)
	ListServers(userInfo models.UserInfo, page, pageSize int) ([]*models.ServerEntity, error)
	ListServersWithKeyword(userInfo models.UserInfo, page, pageSize int, keyword string) ([]*models.ServerEntity, error)
	CountServers(userInfo models.UserInfo) (int64, error)
	CountServersWithKeyword(userInfo models.UserInfo, keyword string) (int64, error)
	CountConfiguredServers(userInfo models.UserInfo) (int64, error)
}

type ServerMutation interface {
	InitDefaultServer(serverIP string)
	UpdateDefaultServer(c *models.Server) error
	CreateServer(userInfo models.UserInfo, server *models.ServerEntity) error
	DeleteServer(userInfo models.UserInfo, serverID string) error
	UpdateServer(userInfo models.UserInfo, server *models.ServerEntity) error
}

type serverQuery struct{ *queryImpl }
type serverMutation struct{ *mutationImpl }

func newServerQuery(base *queryImpl) ServerQuery          { return &serverQuery{base} }
func newServerMutation(base *mutationImpl) ServerMutation { return &serverMutation{base} }

func (m *serverMutation) InitDefaultServer(serverIP string) {
	db := m.ctx.GetApp().GetDBManager().GetDefaultDB()
	db.Where(&models.Server{
		ServerEntity: &models.ServerEntity{
			ServerID: defs.DefaultServerID,
		},
	}).Attrs(&models.Server{
		ServerEntity: &models.ServerEntity{
			ServerID:      defs.DefaultServerID,
			ServerIP:      serverIP,
			ConnectSecret: uuid.New().String(),
		},
	}).FirstOrCreate(&models.Server{})
}

func (q *serverQuery) GetDefaultServer() (*models.ServerEntity, error) {
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	c := &models.Server{}
	err := db.
		Where(&models.Server{ServerEntity: &models.ServerEntity{
			ServerID: defs.DefaultServerID,
		}}).
		First(c).Error
	if err != nil {
		return nil, err
	}
	return c.ServerEntity, nil
}

func (m *serverMutation) UpdateDefaultServer(c *models.Server) error {
	db := m.ctx.GetApp().GetDBManager().GetDefaultDB()
	c.ServerID = defs.DefaultServerID
	err := db.Where(&models.Server{
		ServerEntity: &models.ServerEntity{
			ServerID: defs.DefaultServerID,
		}}).Save(c).Error
	if err != nil {
		return err
	}
	return nil
}

func (q *serverQuery) ValidateServerSecret(serverID string, secret string) (*models.ServerEntity, error) {
	if serverID == "" || secret == "" {
		return nil, fmt.Errorf("invalid request")
	}
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	c := &models.Server{}
	err := db.
		Where(&models.Server{ServerEntity: &models.ServerEntity{
			ServerID: serverID,
		}}).
		First(c).Error
	if err != nil {
		return nil, err
	}
	if c.ConnectSecret != secret {
		return nil, fmt.Errorf("invalid secret")
	}
	return c.ServerEntity, nil
}

func (q *serverQuery) AdminGetServerByServerID(serverID string) (*models.ServerEntity, error) {
	if serverID == "" {
		return nil, fmt.Errorf("invalid server id")
	}
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	c := &models.Server{}
	err := db.
		Where(&models.Server{ServerEntity: &models.ServerEntity{
			ServerID: serverID,
		}}).
		First(c).Error
	if err != nil {
		return nil, err
	}
	return c.ServerEntity, nil
}

func (q *serverQuery) GetServerByServerID(userInfo models.UserInfo, serverID string) (*models.ServerEntity, error) {
	if serverID == "" {
		return nil, fmt.Errorf("invalid server id")
	}
	if userInfo.GetUserID() == defs.DefaultAdminUserID && serverID == defs.DefaultServerID {
		return q.GetDefaultServer()
	}
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	c := &models.Server{}
	err := db.
		Where(&models.Server{ServerEntity: &models.ServerEntity{
			TenantID: userInfo.GetTenantID(),
			UserID:   userInfo.GetUserID(),
			ServerID: serverID,
		}}).
		First(c).Error
	if err != nil {
		return nil, err
	}
	return c.ServerEntity, nil
}

func (m *serverMutation) CreateServer(userInfo models.UserInfo, server *models.ServerEntity) error {
	server.UserID = userInfo.GetUserID()
	server.TenantID = userInfo.GetTenantID()
	c := &models.Server{
		ServerEntity: server,
	}
	db := m.ctx.GetApp().GetDBManager().GetDefaultDB()
	return db.Create(c).Error
}

func (m *serverMutation) DeleteServer(userInfo models.UserInfo, serverID string) error {
	if serverID == "" {
		return fmt.Errorf("invalid server id")
	}
	db := m.ctx.GetApp().GetDBManager().GetDefaultDB()
	return db.Unscoped().Where(
		&models.Server{
			ServerEntity: &models.ServerEntity{
				TenantID: userInfo.GetTenantID(),
				UserID:   userInfo.GetUserID(),
			},
		},
	).Delete(&models.Server{
		ServerEntity: &models.ServerEntity{
			ServerID: serverID,
		},
	}).Error
}

func (m *serverMutation) UpdateServer(userInfo models.UserInfo, server *models.ServerEntity) error {
	c := &models.Server{
		ServerEntity: server,
	}
	if userInfo.GetUserID() == defs.DefaultAdminUserID && server.ServerID == defs.DefaultServerID {
		return m.UpdateDefaultServer(c)
	}
	db := m.ctx.GetApp().GetDBManager().GetDefaultDB()
	return db.Where(
		&models.Server{
			ServerEntity: &models.ServerEntity{
				UserID:   userInfo.GetUserID(),
				TenantID: userInfo.GetTenantID(),
			},
		},
	).Save(c).Error
}

func (q *serverQuery) ListServers(userInfo models.UserInfo, page, pageSize int) ([]*models.ServerEntity, error) {
	if page < 1 || pageSize < 1 {
		return nil, fmt.Errorf("invalid page or page size")
	}

	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	offset := (page - 1) * pageSize

	var servers []*models.Server
	err := db.Where(
		&models.Server{
			ServerEntity: &models.ServerEntity{
				UserID:   userInfo.GetUserID(),
				TenantID: userInfo.GetTenantID(),
			},
		},
	).Or(&models.Server{
		ServerEntity: &models.ServerEntity{
			ServerID: defs.DefaultServerID,
		},
	}).Offset(offset).Limit(pageSize).Find(&servers).Error
	if err != nil {
		return nil, err
	}

	return lo.Map(servers, func(c *models.Server, _ int) *models.ServerEntity {
		return c.ServerEntity
	}), nil
}

func (q *serverQuery) ListServersWithKeyword(userInfo models.UserInfo, page, pageSize int, keyword string) ([]*models.ServerEntity, error) {
	if page < 1 || pageSize < 1 || len(keyword) == 0 {
		return nil, fmt.Errorf("invalid page or page size or keyword")
	}

	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	offset := (page - 1) * pageSize

	var servers []*models.Server
	err := db.Where(
		&models.Server{
			ServerEntity: &models.ServerEntity{
				UserID:   userInfo.GetUserID(),
				TenantID: userInfo.GetTenantID(),
			},
		},
	).Where("server_id like ?", "%"+keyword+"%").
		Offset(offset).Limit(pageSize).Find(&servers).Error
	if err != nil {
		return nil, err
	}

	return lo.Map(servers, func(c *models.Server, _ int) *models.ServerEntity {
		return c.ServerEntity
	}), nil
}

func (q *serverQuery) CountServers(userInfo models.UserInfo) (int64, error) {
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	var count int64
	err := db.Model(&models.Server{}).Where(
		&models.Server{
			ServerEntity: &models.ServerEntity{
				UserID:   userInfo.GetUserID(),
				TenantID: userInfo.GetTenantID(),
			},
		},
	).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (q *serverQuery) CountServersWithKeyword(userInfo models.UserInfo, keyword string) (int64, error) {
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	var count int64
	err := db.Model(&models.Server{}).Where(
		&models.Server{
			ServerEntity: &models.ServerEntity{
				UserID:   userInfo.GetUserID(),
				TenantID: userInfo.GetTenantID(),
			},
		},
	).Where("server_id like ?", "%"+keyword+"%").Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (q *serverQuery) CountConfiguredServers(userInfo models.UserInfo) (int64, error) {
	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	var count int64
	err := db.Model(&models.Server{}).Where(
		&models.Server{
			ServerEntity: &models.ServerEntity{
				UserID:   userInfo.GetUserID(),
				TenantID: userInfo.GetTenantID(),
			},
		},
	).Not(
		&models.Server{
			ServerEntity: &models.ServerEntity{
				ConfigContent: []byte{},
			},
		},
	).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
