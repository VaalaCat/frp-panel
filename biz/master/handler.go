package master

import (
	"embed"

	"github.com/VaalaCat/frp-panel/app"
	"github.com/VaalaCat/frp-panel/biz/master/auth"
	"github.com/VaalaCat/frp-panel/biz/master/client"
	"github.com/VaalaCat/frp-panel/biz/master/platform"
	"github.com/VaalaCat/frp-panel/biz/master/proxy"
	"github.com/VaalaCat/frp-panel/biz/master/server"
	"github.com/VaalaCat/frp-panel/biz/master/shell"
	"github.com/VaalaCat/frp-panel/biz/master/streamlog"
	"github.com/VaalaCat/frp-panel/biz/master/user"
	"github.com/VaalaCat/frp-panel/middleware"
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
	v1 := api.Group("/v1")
	{
		authRouter := v1.Group("/auth")
		{
			authRouter.POST("/login", app.Wrapper(appInstance, auth.LoginHandler))
			authRouter.POST("/register", app.Wrapper(appInstance, auth.RegisterHandler))
			authRouter.GET("/logout", auth.RemoveJWTHandler(appInstance))
			authRouter.POST("/cert", app.Wrapper(appInstance, auth.GetClientCert))
		}
		userRouter := v1.Group("/user", middleware.JWTAuth(appInstance), middleware.AuthCtx(appInstance))
		{
			userRouter.POST("/get", app.Wrapper(appInstance, user.GetUserInfoHandler))
			userRouter.POST("/update", app.Wrapper(appInstance, user.UpdateUserInfoHander))
		}
		platformRouter := v1.Group("/platform", middleware.JWTAuth(appInstance), middleware.AuthCtx(appInstance))
		{
			platformRouter.GET("/baseinfo", platform.GetPlatformInfo(appInstance))
			platformRouter.POST("/clientsstatus", app.Wrapper(appInstance, platform.GetClientsStatus))
		}
		clientRouter := v1.Group("/client", middleware.JWTAuth(appInstance), middleware.AuthCtx(appInstance))
		{
			clientRouter.POST("/get", app.Wrapper(appInstance, client.GetClientHandler))
			clientRouter.POST("/init", app.Wrapper(appInstance, client.InitClientHandler))
			clientRouter.POST("/delete", app.Wrapper(appInstance, client.DeleteClientHandler))
			clientRouter.POST("/list", app.Wrapper(appInstance, client.ListClientsHandler))
		}
		serverRouter := v1.Group("/server", middleware.JWTAuth(appInstance), middleware.AuthCtx(appInstance))
		{
			serverRouter.POST("/get", app.Wrapper(appInstance, server.GetServerHandler))
			serverRouter.POST("/init", app.Wrapper(appInstance, server.InitServerHandler))
			serverRouter.POST("/delete", app.Wrapper(appInstance, server.DeleteServerHandler))
			serverRouter.POST("/list", app.Wrapper(appInstance, server.ListServersHandler))
		}
		frpcRouter := v1.Group("/frpc", middleware.JWTAuth(appInstance), middleware.AuthCtx(appInstance))
		{
			frpcRouter.POST("/update", app.Wrapper(appInstance, client.UpdateFrpcHander))
			frpcRouter.POST("/delete", app.Wrapper(appInstance, client.RemoveFrpcHandler))
			frpcRouter.POST("/stop", app.Wrapper(appInstance, client.StopFRPCHandler))
			frpcRouter.POST("/start", app.Wrapper(appInstance, client.StartFRPCHandler))
		}
		frpsRouter := v1.Group("/frps", middleware.JWTAuth(appInstance), middleware.AuthCtx(appInstance))
		{
			frpsRouter.POST("/update", app.Wrapper(appInstance, server.UpdateFrpsHander))
			frpsRouter.POST("/delete", app.Wrapper(appInstance, server.RemoveFrpsHandler))
		}
		proxyRouter := v1.Group("/proxy", middleware.JWTAuth(appInstance), middleware.AuthCtx(appInstance))
		{
			proxyRouter.POST("/get_by_cid", app.Wrapper(appInstance, proxy.GetProxyStatsByClientID))
			proxyRouter.POST("/get_by_sid", app.Wrapper(appInstance, proxy.GetProxyStatsByServerID))
			proxyRouter.POST("/list_configs", app.Wrapper(appInstance, proxy.ListProxyConfigs))
			proxyRouter.POST("/create_config", app.Wrapper(appInstance, proxy.CreateProxyConfig))
			proxyRouter.POST("/update_config", app.Wrapper(appInstance, proxy.UpdateProxyConfig))
			proxyRouter.POST("/delete_config", app.Wrapper(appInstance, proxy.DeleteProxyConfig))
			proxyRouter.POST("/get_config", app.Wrapper(appInstance, proxy.GetProxyConfig))
		}
		v1.GET("/pty/:clientID", middleware.JWTAuth(appInstance), middleware.AuthCtx(appInstance), shell.PTYHandler(appInstance))
		v1.GET("/log", middleware.JWTAuth(appInstance), middleware.AuthCtx(appInstance), streamlog.GetLogHandler(appInstance))
	}
}
