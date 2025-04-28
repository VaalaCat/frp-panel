package middleware

import (
	"errors"
	"strings"
	"time"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/VaalaCat/frp-panel/utils/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/cast"
)

// JWTMAuth check if authed and set uid to context
func JWTAuth(appInstance app.Application) func(c *gin.Context) {
	return func(c *gin.Context) {
		defer func() {
			logger.Logger(c).Info("finish jwt middleware")
		}()

		var tokenStr string

		if tokenStr = c.Copy().Query(defs.TokenKey); len(tokenStr) != 0 {
			if t, err := utils.ParseToken(conf.JWTSecret(appInstance.GetConfig()), tokenStr); err == nil {
				for k, v := range t {
					c.Set(k, v)
				}
				logger.Logger(c).Infof("query auth success")
				if err = resignAndPatchCtxJWT(c, appInstance, cast.ToInt(t[defs.UserIDKey]), t, tokenStr); err != nil {
					logger.Logger(c).WithError(err).Errorf("resign jwt error")
					common.ErrUnAuthorized(c, "resign jwt error")
					c.Abort()
					return
				}
				c.Next()
				SetToken(c, appInstance, cast.ToInt(t[defs.UserIDKey]), t)
				return
			}
			logger.Logger(c).Infof("query auth failed")
		}

		cookieToken, err := c.Cookie(appInstance.GetConfig().App.CookieName)
		if err == nil {
			if t, err := utils.ParseToken(conf.JWTSecret(appInstance.GetConfig()), cookieToken); err == nil {
				for k, v := range t {
					c.Set(k, v)
				}
				logger.Logger(c).Infof("cookie auth success")
				if err = resignAndPatchCtxJWT(c, appInstance, cast.ToInt(t[defs.UserIDKey]), t, cookieToken); err != nil {
					logger.Logger(c).WithError(err).Errorf("resign jwt error")
					common.ErrUnAuthorized(c, "resign jwt error")
					c.Abort()
					return
				}
				c.Next()
				return
			} else {
				logger.Logger(c).WithError(err).Errorf("jwt middleware parse token error")
				common.ErrUnAuthorized(c, "invalid authorization")
				c.Abort()
				return
			}
		}

		tokenStr = c.Request.Header.Get(defs.AuthorizationKey)
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
		if tokenStr == "" || tokenStr == "null" {
			logger.Logger(c).WithError(errors.New("authorization is empty")).Infof("authorization is empty")
			common.ErrUnAuthorized(c, "invalid authorization")
			c.Abort()
			return
		}

		if t, err := utils.ParseToken(conf.JWTSecret(appInstance.GetConfig()), tokenStr); err == nil {
			for k, v := range t {
				c.Set(k, v)
			}
			logger.Logger(c).Infof("header auth success")
			if err = resignAndPatchCtxJWT(c, appInstance, cast.ToInt(t[defs.UserIDKey]), t, tokenStr); err != nil {
				logger.Logger(c).WithError(err).Errorf("resign jwt error")
				common.ErrUnAuthorized(c, "resign jwt error")
				c.Abort()
				return
			}
			c.Next()
			return
		} else {
			logger.Logger(c).WithError(err).Errorf("jwt middleware parse token error")
		}
	}
}

func resignAndPatchCtxJWT(c *gin.Context, appInstance app.Application, userID int, t jwt.MapClaims, tokenStr string) error {
	tokenExpire, _ := t.GetExpirationTime()
	now := time.Now().Add(time.Duration(appInstance.GetConfig().App.CookieAge/2) * time.Second)
	if now.Before(tokenExpire.Time) {
		logger.Logger(c).Infof("jwt not going to expire, continue to use old one")
		c.Set(defs.TokenKey, tokenStr)
		return nil
	}

	tokenStr, err := SetToken(c, appInstance, userID, t)
	if err != nil {
		logger.Logger(c).WithError(err).Errorf("resign jwt error")
		return err
	}

	PushTokenStr(c, appInstance, tokenStr)

	logger.Logger(c).Infof("jwt going to expire, resigning token")
	return nil
}

// SetToken 设置新token并写入ctx
func SetToken(c *gin.Context, appInstance app.Application, userID int, payload jwt.MapClaims) (string, error) {
	logger.Logger(c).Debugf("set token for userID:[%d]", userID)
	if payload == nil {
		payload = make(map[string]interface{})
	}

	payload[defs.UserIDKey] = userID

	token, err := conf.GetJWTWithPayload(appInstance.GetConfig(), userID, payload)
	if err != nil {
		return "", err
	}
	c.Set(defs.TokenKey, token)
	return token, nil
}

// PushTokenStr 推送token到客户端
func PushTokenStr(c *gin.Context, appInstance app.Application, tokenStr string) {
	logger.Logger(c).Infof("push new token to client")
	c.SetCookie(appInstance.GetConfig().App.CookieName,
		tokenStr,
		appInstance.GetConfig().App.CookieAge,
		appInstance.GetConfig().App.CookiePath,
		appInstance.GetConfig().App.CookieDomain,
		appInstance.GetConfig().App.CookieSecure,
		appInstance.GetConfig().App.CookieHTTPOnly)
	c.Header(defs.SetAuthorizationKey, tokenStr)
}
