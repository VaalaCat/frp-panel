package middleware

import (
	"errors"
	"strings"
	"time"

	"github.com/VaalaCat/frp-panel/app"
	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/logger"
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
				if err = resignAndPatchCtxJWT(c, appInstance, t, tokenStr); err != nil {
					logger.Logger(c).WithError(err).Errorf("resign jwt error")
					common.ErrUnAuthorized(c, "resign jwt error")
					c.Abort()
					return
				}
				c.Next()
				SetToken(c, appInstance, utils.ToStr(t[defs.UserIDKey]))
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
				if err = resignAndPatchCtxJWT(c, appInstance, t, cookieToken); err != nil {
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
			if err = resignAndPatchCtxJWT(c, appInstance, t, tokenStr); err != nil {
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

func resignAndPatchCtxJWT(c *gin.Context, appInstance app.Application, t jwt.MapClaims, tokenStr string) error {
	tokenExpire, _ := t.GetExpirationTime()
	now := time.Now().Add(time.Duration(appInstance.GetConfig().App.CookieAge/2) * time.Second)
	if now.Before(tokenExpire.Time) {
		logger.Logger(c).Infof("jwt not going to expire, continue to use old one")
		c.Set(defs.TokenKey, tokenStr)
		return nil
	}

	token, err := utils.GetJwtTokenFromMap(conf.JWTSecret(appInstance.GetConfig()),
		time.Now().Unix(),
		int64(appInstance.GetConfig().App.CookieAge),
		map[string]string{defs.UserIDKey: utils.ToStr(t[defs.UserIDKey])})
	if err != nil {
		c.Set(defs.TokenKey, tokenStr)
		logger.Logger(c).WithError(err).Errorf("resign jwt error")
		return err
	}

	logger.Logger(c).Infof("jwt going to expire, resigning token")
	c.Header(defs.SetAuthorizationKey, token)
	c.SetCookie(appInstance.GetConfig().App.CookieName,
		token,
		appInstance.GetConfig().App.CookieAge,
		appInstance.GetConfig().App.CookiePath,
		appInstance.GetConfig().App.CookieDomain,
		appInstance.GetConfig().App.CookieSecure,
		appInstance.GetConfig().App.CookieHTTPOnly)
	c.Set(defs.TokenKey, token)
	return nil
}

func SetToken(c *gin.Context, appInstance app.Application, uid string) (string, error) {
	logger.Logger(c).Infof("set token for uid:[%s]", uid)
	token, err := utils.GetJwtTokenFromMap(conf.JWTSecret(appInstance.GetConfig()),
		time.Now().Unix(),
		int64(appInstance.GetConfig().App.CookieAge),
		map[string]string{defs.UserIDKey: uid})
	if err != nil {
		return "", err
	}
	c.SetCookie(appInstance.GetConfig().App.CookieName,
		token,
		appInstance.GetConfig().App.CookieAge,
		appInstance.GetConfig().App.CookiePath,
		appInstance.GetConfig().App.CookieDomain,
		appInstance.GetConfig().App.CookieSecure,
		appInstance.GetConfig().App.CookieHTTPOnly)
	c.Set(defs.TokenKey, token)
	c.Header(defs.SetAuthorizationKey, token)
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
