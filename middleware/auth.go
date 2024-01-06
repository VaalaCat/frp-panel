package middleware

import (
	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/dao"
	"github.com/VaalaCat/frp-panel/models"
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func AuthCtx(c *gin.Context) {
	var uid int64
	var err error
	var u *models.UserEntity

	defer func() {
		logrus.WithContext(c).Info("finish auth user middleware")
	}()

	userID, ok := utils.GetValue[int](c, common.UserIDKey)
	if !ok {
		logrus.WithContext(c).Errorf("invalid user id")
		common.ErrUnAuthorized(c, "token invalid")
		c.Abort()
		return
	}

	u, err = dao.GetUserByUserID(userID)
	if err != nil {
		logrus.WithContext(c).Errorf("get user by user id failed: %v", err)
		common.ErrUnAuthorized(c, "token invalid")
		c.Abort()
		return
	}

	logrus.WithContext(c).Infof("auth middleware authed user is: [%+v]", u)

	if u.Valid() {
		logrus.WithContext(c).Infof("set auth user to context, login success")
		c.Set(common.UserInfoKey, u)
		c.Next()
		return
	} else {
		if uid == 1 {
			logrus.WithContext(c).Infof("seems to be admin assign token login")
			c.Next()
			return
		}
		logrus.WithContext(c).Errorf("invalid authorization, auth ctx middleware login failed")
		common.ErrUnAuthorized(c, "token invalid")
		c.Abort()
		return
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
