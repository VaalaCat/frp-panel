package auth

import (
	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/gin-gonic/gin"
)

func RemoveJWTHandler(c *gin.Context) {
	c.SetCookie(conf.Get().App.CookieName,
		"", -1,
		conf.Get().App.CookiePath,
		conf.Get().App.CookieDomain,
		conf.Get().App.CookieSecure,
		conf.Get().App.CookieHTTPOnly)
	common.OKResp(c, &pb.CommonResponse{Status: &pb.Status{Code: pb.RespCode_RESP_CODE_SUCCESS,
		Message: "ok"}})
}
