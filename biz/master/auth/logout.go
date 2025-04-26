package auth

import (
	"github.com/VaalaCat/frp-panel/app"
	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/gin-gonic/gin"
)

func RemoveJWTHandler(appInstance app.Application) func(c *gin.Context) {
	return func(ctx *gin.Context) {
		cfg := appInstance.GetConfig()
		ctx.SetCookie(cfg.App.CookieName,
			"", -1,
			cfg.App.CookiePath,
			cfg.App.CookieDomain,
			cfg.App.CookieSecure,
			cfg.App.CookieHTTPOnly)
		common.OKResp(ctx, &pb.CommonResponse{Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS,
			Message: "ok"}})
	}
}
