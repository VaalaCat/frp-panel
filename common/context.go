package common

import (
	"context"

	"github.com/VaalaCat/frp-panel/models"
)

func GetUserInfo(c context.Context) models.UserInfo {
	val := c.Value(UserInfoKey)
	if val == nil {
		return nil
	}

	u, ok := val.(*models.UserEntity)
	if !ok {
		return nil
	}

	return u
}
