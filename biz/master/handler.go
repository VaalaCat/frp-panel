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
	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(fs embed.FS) *gin.Engine {
	router := gin.Default()
	HandleStaticFile(fs, router)
	ConfigureRouter(router)
	return router
}

func ConfigureRouter(router *gin.Engine) {
	router.POST("/auth", auth.MakeGinHandlerFunc(auth.HandleLogin))

	api := router.Group("/api")
	v1 := api.Group("/v1")
	{
		authRouter := v1.Group("/auth")
		{
			authRouter.POST("/login", common.Wrapper(auth.LoginHandler))
			authRouter.POST("/register", common.Wrapper(auth.RegisterHandler))
			authRouter.GET("/logout", auth.RemoveJWTHandler)
			authRouter.POST("/cert", common.Wrapper(auth.GetClientCert))
		}
		userRouter := v1.Group("/user", middleware.JWTAuth, middleware.AuthCtx)
		{
			userRouter.POST("/get", common.Wrapper(user.GetUserInfoHandler))
			userRouter.POST("/update", common.Wrapper(user.UpdateUserInfoHander))
		}
		platformRouter := v1.Group("/platform", middleware.JWTAuth, middleware.AuthCtx)
		{
			platformRouter.GET("/baseinfo", platform.GetPlatformInfo)
			platformRouter.POST("/clientsstatus", common.Wrapper(platform.GetClientsStatus))
		}
		clientRouter := v1.Group("/client", middleware.JWTAuth, middleware.AuthCtx)
		{
			clientRouter.POST("/get", common.Wrapper(client.GetClientHandler))
			clientRouter.POST("/init", common.Wrapper(client.InitClientHandler))
			clientRouter.POST("/delete", common.Wrapper(client.DeleteClientHandler))
			clientRouter.POST("/list", common.Wrapper(client.ListClientsHandler))
		}
		serverRouter := v1.Group("/server", middleware.JWTAuth, middleware.AuthCtx)
		{
			serverRouter.POST("/get", common.Wrapper(server.GetServerHandler))
			serverRouter.POST("/init", common.Wrapper(server.InitServerHandler))
			serverRouter.POST("/delete", common.Wrapper(server.DeleteServerHandler))
			serverRouter.POST("/list", common.Wrapper(server.ListServersHandler))
		}
		frpcRouter := v1.Group("/frpc", middleware.JWTAuth, middleware.AuthCtx)
		{
			frpcRouter.POST("/update", common.Wrapper(client.UpdateFrpcHander))
			frpcRouter.POST("/delete", common.Wrapper(client.RemoveFrpcHandler))
			frpcRouter.POST("/stop", common.Wrapper(client.StopFRPCHandler))
			frpcRouter.POST("/start", common.Wrapper(client.StartFRPCHandler))
		}
		frpsRouter := v1.Group("/frps", middleware.JWTAuth, middleware.AuthCtx)
		{
			frpsRouter.POST("/update", common.Wrapper(server.UpdateFrpsHander))
			frpsRouter.POST("/delete", common.Wrapper(server.RemoveFrpsHandler))
		}
		proxyRouter := v1.Group("/proxy", middleware.JWTAuth, middleware.AuthCtx)
		{
			proxyRouter.POST("/get_by_cid", common.Wrapper(proxy.GetProxyStatsByClientID))
			proxyRouter.POST("/get_by_sid", common.Wrapper(proxy.GetProxyStatsByServerID))
			proxyRouter.POST("/list_configs", common.Wrapper(proxy.ListProxyConfigs))
			proxyRouter.POST("/create_config", common.Wrapper(proxy.CreateProxyConfig))
			proxyRouter.POST("/update_config", common.Wrapper(proxy.UpdateProxyConfig))
			proxyRouter.POST("/delete_config", common.Wrapper(proxy.DeleteProxyConfig))
			proxyRouter.POST("/get_config", common.Wrapper(proxy.GetProxyConfig))
		}
		v1.GET("/pty/:clientID", middleware.JWTAuth, middleware.AuthCtx, shell.PTYHandler)
		v1.GET("/log", middleware.JWTAuth, middleware.AuthCtx, streamlog.GetLogHander)
	}
}
