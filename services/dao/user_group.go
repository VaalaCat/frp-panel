package dao

import (
	"fmt"

	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/models"
)

func (q *queryImpl) CreateGroup(userInfo models.UserInfo, groupID, groupName, comment string) (*models.UserGroup, error) {
	if groupID == "" || groupName == "" {
		return nil, fmt.Errorf("invalid group id or group name")
	}

	if userInfo.GetRole() != defs.UserRole_Admin {
		return nil, fmt.Errorf("only admin can create group")
	}

	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()

	g := &models.UserGroup{
		TenantID:  userInfo.GetTenantID(),
		GroupID:   groupID,
		GroupName: groupName,
		Comment:   comment,
	}
	err := db.Create(g).Error
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (q *queryImpl) DeleteGroup(userInfo models.UserInfo, groupID string) error {
	if userInfo.GetRole() != defs.UserRole_Admin {
		return fmt.Errorf("only admin can delete group")
	}

	db := q.ctx.GetApp().GetDBManager().GetDefaultDB()
	return db.Unscoped().Where(&models.UserGroup{
		TenantID: userInfo.GetTenantID(),
		GroupID:  groupID,
	}).Delete(&models.UserGroup{}).Error
}
