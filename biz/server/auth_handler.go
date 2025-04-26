package server

import (
	"net/http"

	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/utils/logger"
	plugin "github.com/fatedier/frp/pkg/plugin/server"
	"github.com/gin-gonic/gin"
)

func NewRouter(appInstance app.Application) *gin.Engine {
	router := gin.Default()
	router.POST("/auth", MakeGinHandlerFunc(appInstance, HandleLogin))
	return router
}

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

type HandlerFunc func(ctx *app.Context) (interface{}, error)

func MakeGinHandlerFunc(appInstance app.Application, handler HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		res, err := handler(app.NewContext(ctx, appInstance))
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

func HandleLogin(ctx *app.Context) (interface{}, error) {
	var r plugin.Request
	var content plugin.LoginContent
	r.Content = &content
	if err := ctx.GetGinCtx().BindJSON(&r); err != nil {
		return nil, &HTTPError{
			Code: http.StatusBadRequest,
			Err:  err,
		}
	}

	var res plugin.Response
	token := content.Metas[defs.FRPAuthTokenKey]
	if len(content.User) == 0 || len(token) == 0 {
		res.Reject = true
		res.RejectReason = "user or meta token can not be empty"
		return res, nil
	}
	cli := ctx.GetApp().GetMasterCli()
	authResponse, err := cli.Call().FRPCAuth(ctx, &pb.FRPAuthRequest{User: content.User, Token: token})
	if err != nil {
		res.Reject = true
		res.RejectReason = "invalid meta token"
		return res, nil
	}

	if authResponse.GetOk() {
		res.Unchange = true
		return res, nil
	}

	res.Reject = true
	res.RejectReason = "invalid meta token"
	return res, nil
}
