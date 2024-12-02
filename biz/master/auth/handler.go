package auth

import (
	"net/http"

	"github.com/VaalaCat/frp-panel/cache"
	"github.com/VaalaCat/frp-panel/dao"
	"github.com/VaalaCat/frp-panel/logger"
	plugin "github.com/fatedier/frp/pkg/plugin/server"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Msg string `json:"msg"`
}

type HTTPError struct {
	Code int
	Err  error
}

func (e *HTTPError) Error() string {
	return e.Err.Error()
}

type HandlerFunc func(ctx *gin.Context) (interface{}, error)

func MakeGinHandlerFunc(handler HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res, err := handler(ctx)
		if err != nil {
			logger.Logger(ctx).Infof("handle %s error: %v", ctx.Request.URL.Path, err)
			switch e := err.(type) {
			case *HTTPError:
				ctx.JSON(e.Code, &Response{Msg: e.Err.Error()})
			default:
				ctx.JSON(500, &Response{Msg: err.Error()})
			}
			return
		}
		ctx.JSON(http.StatusOK, res)
	}
}

func HandleLogin(ctx *gin.Context) (interface{}, error) {
	var r plugin.Request
	var content plugin.LoginContent
	r.Content = &content
	if err := ctx.BindJSON(&r); err != nil {
		return nil, &HTTPError{
			Code: http.StatusBadRequest,
			Err:  err,
		}
	}

	var res plugin.Response
	token := content.Metas["token"]
	if len(content.User) == 0 || len(token) == 0 {
		res.Reject = true
		res.RejectReason = "user or meta token can not be empty"
		return res, nil
	}

	userToken, err := cache.Get().Get([]byte(content.User))
	if err != nil {
		u, err := dao.GetUserByUserName(content.User)
		if err != nil || u == nil {
			res.Reject = true
			res.RejectReason = "invalid frp auth"
			return res, nil
		}
		cache.Get().Set([]byte(u.GetUserName()), []byte(u.GetToken()), 0)
		userToken = []byte(u.GetToken())
	}

	if string(userToken) == token {
		res.Unchange = true
		return res, nil
	}

	res.Reject = true
	res.RejectReason = "invalid meta token"
	return res, nil
}
