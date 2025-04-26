package middleware

import (
	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/dao"
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/VaalaCat/frp-panel/utils/logger"
	"github.com/gin-gonic/gin"
)

func AuthCtx(appInstance app.Application) func(*gin.Context) {
	return func(c *gin.Context) {
		var uid int64
		var err error
		var u *models.UserEntity
		appCtx := app.NewContext(c, appInstance)

		defer func() {
			logger.Logger(c).Info("finish auth user middleware")
		}()

		userID, ok := utils.GetValue[int](c, defs.UserIDKey)
		if !ok {
			logger.Logger(c).Errorf("invalid user id")
			common.ErrUnAuthorized(c, "token invalid")
			c.Abort()
			return
		}

		u, err = dao.NewQuery(appCtx).GetUserByUserID(userID)
		if err != nil {
			logger.Logger(c).Errorf("get user by user id failed: %v", err)
			common.ErrUnAuthorized(c, "token invalid")
			c.Abort()
			return
		}

		logger.Logger(c).Infof("auth middleware authed user is: [%+v]", u)

		if u.Valid() {
			logger.Logger(c).Infof("set auth user to context, login success")
			c.Set(defs.UserInfoKey, u)
			c.Next()
			return
		} else {
			if uid == 1 {
				logger.Logger(c).Infof("seems to be admin assign token login")
				c.Next()
				return
			}
			logger.Logger(c).Errorf("invalid authorization, auth ctx middleware login failed")
			common.ErrUnAuthorized(c, "token invalid")
			c.Abort()
			return
		}
	}
}

func AuthAdmin(c *gin.Context) {
	u := common.GetUserInfo(c)
	if u != nil && u.GetRole() == models.ROLE_ADMIN {
		common.ErrUnAuthorized(c, "permission denied")
		c.Abort()
		return
	}
	c.Next()
}
