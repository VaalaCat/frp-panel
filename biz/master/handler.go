package master

import (
	"embed"

	"github.com/VaalaCat/frp-panel/biz/master/auth"
	"github.com/VaalaCat/frp-panel/biz/master/client"
	"github.com/VaalaCat/frp-panel/biz/master/platform"
	"github.com/VaalaCat/frp-panel/biz/master/proxy"
	"github.com/VaalaCat/frp-panel/biz/master/server"
	"github.com/VaalaCat/frp-panel/biz/master/shell"
	"github.com/VaalaCat/frp-panel/biz/master/streamlog"
	"github.com/VaalaCat/frp-panel/biz/master/user"
	"github.com/VaalaCat/frp-panel/middleware"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/gin-gonic/gin"
)

func NewRouter(fs embed.FS, appInstance app.Application) *gin.Engine {
	router := gin.Default()
	HandleStaticFile(fs, router)
	ConfigureRouter(appInstance, router)
	return router
}

func ConfigureRouter(appInstance app.Application, router *gin.Engine) {
	router.POST("/auth", auth.MakeGinHandlerFunc(appInstance, auth.HandleLogin))

	api := router.Group("/api")
	api.POST("/v1/auth/cert", app.Wrapper(appInstance, auth.GetClientCert))
	api.POST("/v1/auth/login", app.Wrapper(appInstance, auth.LoginHandler))
	api.POST("/v1/auth/register", app.Wrapper(appInstance, auth.RegisterHandler))
	api.GET("/v1/auth/logout", auth.RemoveJWTHandler(appInstance))

	v1 := api.Group("/v1", middleware.JWTAuth(appInstance), middleware.AuthCtx(appInstance), middleware.RBAC(appInstance))
	{
		userRouter := v1.Group("/user")
		{
			userRouter.POST("/get", app.Wrapper(appInstance, user.GetUserInfoHandler))
			userRouter.POST("/update", app.Wrapper(appInstance, user.UpdateUserInfoHander))
			userRouter.POST("/sign-token", app.Wrapper(appInstance, user.SignTokenHandler))
		}
		platformRouter := v1.Group("/platform")
		{
			platformRouter.GET("/baseinfo", platform.GetPlatformInfo(appInstance))
			platformRouter.POST("/clientsstatus", app.Wrapper(appInstance, platform.GetClientsStatus))
		}
		clientRouter := v1.Group("/client")
		{
			clientRouter.POST("/get", app.Wrapper(appInstance, client.GetClientHandler))
			clientRouter.POST("/init", app.Wrapper(appInstance, client.InitClientHandler))
			clientRouter.POST("/delete", app.Wrapper(appInstance, client.DeleteClientHandler))
			clientRouter.POST("/list", app.Wrapper(appInstance, client.ListClientsHandler))
		}
		serverRouter := v1.Group("/server")
		{
			serverRouter.POST("/get", app.Wrapper(appInstance, server.GetServerHandler))
			serverRouter.POST("/init", app.Wrapper(appInstance, server.InitServerHandler))
			serverRouter.POST("/delete", app.Wrapper(appInstance, server.DeleteServerHandler))
			serverRouter.POST("/list", app.Wrapper(appInstance, server.ListServersHandler))
		}
		frpcRouter := v1.Group("/frpc")
		{
			frpcRouter.POST("/update", app.Wrapper(appInstance, client.UpdateFrpcHander))
			frpcRouter.POST("/delete", app.Wrapper(appInstance, client.RemoveFrpcHandler))
			frpcRouter.POST("/stop", app.Wrapper(appInstance, client.StopFRPCHandler))
			frpcRouter.POST("/start", app.Wrapper(appInstance, client.StartFRPCHandler))
		}
		frpsRouter := v1.Group("/frps")
		{
			frpsRouter.POST("/update", app.Wrapper(appInstance, server.UpdateFrpsHander))
			frpsRouter.POST("/delete", app.Wrapper(appInstance, server.RemoveFrpsHandler))
		}
		proxyRouter := v1.Group("/proxy")
		{
			proxyRouter.POST("/get_by_cid", app.Wrapper(appInstance, proxy.GetProxyStatsByClientID))
			proxyRouter.POST("/get_by_sid", app.Wrapper(appInstance, proxy.GetProxyStatsByServerID))
			proxyRouter.POST("/list_configs", app.Wrapper(appInstance, proxy.ListProxyConfigs))
			proxyRouter.POST("/create_config", app.Wrapper(appInstance, proxy.CreateProxyConfig))
			proxyRouter.POST("/update_config", app.Wrapper(appInstance, proxy.UpdateProxyConfig))
			proxyRouter.POST("/delete_config", app.Wrapper(appInstance, proxy.DeleteProxyConfig))
			proxyRouter.POST("/get_config", app.Wrapper(appInstance, proxy.GetProxyConfig))
		}
		v1.GET("/pty/:clientID", shell.PTYHandler(appInstance))
		v1.GET("/log", streamlog.GetLogHandler(appInstance))
	}
}
