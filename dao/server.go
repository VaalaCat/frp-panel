package dao

import (
	"fmt"

	"github.com/VaalaCat/frp-panel/models"
	"github.com/samber/lo"
)

func ValidateServerSecret(serverID string, secret string) (*models.ServerEntity, error) {
	if serverID == "" || secret == "" {
		return nil, fmt.Errorf("invalid request")
	}
	db := models.GetDBManager().GetDefaultDB()
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

func AdminGetServerByServerID(serverID string) (*models.ServerEntity, error) {
	if serverID == "" {
		return nil, fmt.Errorf("invalid server id")
	}
	db := models.GetDBManager().GetDefaultDB()
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

func GetServerByServerID(userInfo models.UserInfo, serverID string) (*models.ServerEntity, error) {
	if serverID == "" {
		return nil, fmt.Errorf("invalid server id")
	}
	db := models.GetDBManager().GetDefaultDB()
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

func CreateServer(userInfo models.UserInfo, server *models.ServerEntity) error {
	server.UserID = userInfo.GetUserID()
	server.TenantID = userInfo.GetTenantID()
	c := &models.Server{
		ServerEntity: server,
	}
	db := models.GetDBManager().GetDefaultDB()
	return db.Create(c).Error
}

func DeleteServer(userInfo models.UserInfo, serverID string) error {
	if serverID == "" {
		return fmt.Errorf("invalid server id")
	}
	db := models.GetDBManager().GetDefaultDB()
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

func UpdateServer(userInfo models.UserInfo, server *models.ServerEntity) error {
	c := &models.Server{
		ServerEntity: server,
	}
	db := models.GetDBManager().GetDefaultDB()
	return db.Where(
		&models.Server{
			ServerEntity: &models.ServerEntity{
				UserID:   userInfo.GetUserID(),
				TenantID: userInfo.GetTenantID(),
			},
		},
	).Save(c).Error
}

func ListServers(userInfo models.UserInfo, page, pageSize int) ([]*models.ServerEntity, error) {
	if page < 1 || pageSize < 1 {
		return nil, fmt.Errorf("invalid page or page size")
	}

	db := models.GetDBManager().GetDefaultDB()
	offset := (page - 1) * pageSize

	var servers []*models.Server
	err := db.Where(
		&models.Server{
			ServerEntity: &models.ServerEntity{
				UserID:   userInfo.GetUserID(),
				TenantID: userInfo.GetTenantID(),
			},
		},
	).Offset(offset).Limit(pageSize).Find(&servers).Error
	if err != nil {
		return nil, err
	}

	return lo.Map(servers, func(c *models.Server, _ int) *models.ServerEntity {
		return c.ServerEntity
	}), nil
}

func CountServers(userInfo models.UserInfo) (int64, error) {
	db := models.GetDBManager().GetDefaultDB()
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

func CountUnconfiguredServers(userInfo models.UserInfo) (int64, error) {
	db := models.GetDBManager().GetDefaultDB()
	var count int64
	err := db.Model(&models.Server{}).Where(
		&models.Server{
			ServerEntity: &models.ServerEntity{
				UserID:        userInfo.GetUserID(),
				TenantID:      userInfo.GetTenantID(),
				ConfigContent: []byte{},
			},
		},
	).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func CountConfiguredServers(userInfo models.UserInfo) (int64, error) {
	db := models.GetDBManager().GetDefaultDB()
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
