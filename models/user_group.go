package models

import "time"

type UserGroup struct {
	GroupID   string `json:"group_id" gorm:"primaryKey"`
	GroupName string `json:"group_name" gorm:"type:varchar(255);uniqueIndex:idx_group_tenant_name;not null"`
	TenantID  int    `json:"tenant_id" gorm:"uniqueIndex:idx_group_tenant_name;not null"`
	Comment   string `json:"comment"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Users []*User `json:"users,omitempty" gorm:"many2many:user_group_memberships;"`
}

func (u *UserGroup) TableName() string {
	return "user_groups"
}
