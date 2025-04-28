package middleware

import (
	"regexp"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/VaalaCat/frp-panel/utils/logger"
	"github.com/gin-gonic/gin"
)

func RBAC(appInstance app.Application) func(*gin.Context) {
	return func(c *gin.Context) {
		// appCtx := app.NewContext(c, appInstance)
		perms, err := common.GetTokenPermission(c)
		userInfo := common.GetUserInfo(c)
		token := common.GetTokenString(c)
		path := c.Request.URL.Path
		method := c.Request.Method

		if err != nil {
			logger.Logger(c).WithError(err).Errorf("get token permission error, token: [%s], userInfo:[%s]", token, utils.MarshalForJson(userInfo.GetSafeUserInfo()))
			common.ErrResp(c, &pb.CommonResponse{Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: err.Error()}}, err.Error())
			c.Abort()
			return
		}

		if len(perms) == 0 {
			logger.Logger(c).WithError(err).Errorf("user has no permission, token: [%s], userInfo:[%s]", token, utils.MarshalForJson(userInfo.GetSafeUserInfo()))
			common.ErrResp(c, &pb.CommonResponse{Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "user has no permission"}}, "user has no permission")
			c.Abort()
			return
		}

		for _, perm := range perms {
			if ruleMatched(ruleMatchParam{
				RuleMethod:    perm.Method,
				RulePath:      perm.Path,
				RequestPath:   path,
				RequestMethod: method,
			}) {
				logger.Logger(c).Infof("user has api permission, continue")
				c.Next()
				return
			}
		}

		logger.Logger(c).Errorf("user has no permission, perms: %s, userInfo: [%s], ", utils.MarshalForJson(perms), utils.MarshalForJson(userInfo.GetSafeUserInfo()))
		common.ErrResp(c, &pb.CommonResponse{Status: &pb.Status{Code: pb.RespCode_RESP_CODE_INVALID, Message: "user has no permission"}}, "user has no permission")
		c.Abort()
		return
	}
}

type ruleMatchParam struct {
	RuleMethod    string
	RulePath      string
	RequestPath   string
	RequestMethod string
}

func ruleMatched(param ruleMatchParam) bool {
	methodMatch := false
	if param.RuleMethod == param.RequestMethod || param.RuleMethod == "*" {
		methodMatch = true
	}

	if !methodMatch {
		return false
	}

	pathMatch := false
	if param.RulePath == "*" {
		pathMatch = true
	} else {
		pathMatch = regexp.MustCompile(param.RulePath).MatchString(param.RequestPath)
	}

	return pathMatch
}
