package common

import (
	"context"
	"encoding/json"

	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/models"
)

func GetUserInfo(c context.Context) models.UserInfo {
	val := c.Value(defs.UserInfoKey)
	if val == nil {
		return nil
	}

	u, ok := val.(*models.UserEntity)
	if !ok {
		return nil
	}

	return u
}

func GetTokenPermission(c context.Context) ([]defs.APIPermission, error) {
	val := c.Value(defs.TokenPayloadKey_Permissions)
	if val == nil {
		return nil, nil
	}

	raw, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}

	perms := []defs.APIPermission{}
	err = json.Unmarshal(raw, &perms)
	if err != nil {
		return nil, err
	}

	return perms, nil
}

func GetTokenString(c context.Context) string {
	return c.Value(defs.TokenKey).(string)
}
