package app

import (
	"context"
	"sync"

	"github.com/VaalaCat/frp-panel/conf"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/credentials"
)

type Application interface {
	GetStreamLogHookMgr() StreamLogHookMgr
	SetStreamLogHookMgr(StreamLogHookMgr)
	GetShellPTYMgr() ShellPTYMgr
	SetShellPTYMgr(ShellPTYMgr)
	GetClientLogManager() ClientLogManager
	SetClientLogManager(ClientLogManager)
	GetDBManager() DBManager
	SetDBManager(DBManager)
	GetClientRecvMap() *sync.Map
	SetClientRecvMap(*sync.Map)
	GetClientsManager() ClientsManager
	SetClientsManager(ClientsManager)
	GetMasterCli() MasterClient
	SetMasterCli(MasterClient)
	GetClientRPCHandler() ClientRPCHandler
	SetClientRPCHandler(ClientRPCHandler)
	GetServerHandler() ServerHandler
	SetServerHandler(ServerHandler)
	GetClientController() ClientController
	SetClientController(ClientController)
	GetServerController() ServerController
	SetServerController(ServerController)
	GetConfig() conf.Config
	SetConfig(conf.Config)
	GetRPCCred() credentials.TransportCredentials
	SetRPCCred(credentials.TransportCredentials)
}

type Context struct {
	context.Context
	appInstance Application
}

func (c *Context) GetApp() Application {
	return c.appInstance
}

func (c *Context) GetGinCtx() *gin.Context {
	return c.Context.(*gin.Context)
}

func (c *Context) GetCtx() context.Context {
	return c.Context
}

func NewContext(c context.Context, appInstance Application) *Context {
	return &Context{
		Context:     c,
		appInstance: appInstance,
	}
}

func NewApp() Application {
	return &application{}
}

// var app *application

// func GetApp() Application {
// 	if app == nil {
// 		app = NewApp().(*application)
// 	}
// 	return app
// }

// func SetAppInstance(a Application) {
// 	app = a.(*application)
// }
