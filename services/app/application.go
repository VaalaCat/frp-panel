package app

import (
	"context"
	"sync"

	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
	GetCurrentRole() string
	SetCurrentRole(string)
	GetEnforcer() *casbin.Enforcer
	SetEnforcer(*casbin.Enforcer)
	GetPermManager() PermissionManager
	SetPermManager(PermissionManager)
	GetWorkerExecManager() WorkerExecManager
	SetWorkerExecManager(WorkerExecManager)
	GetWorkersManager() WorkersManager
	SetWorkersManager(WorkersManager)
	SetLogger(*logrus.Logger)
	Logger(ctx context.Context) *logrus.Entry
	GetClientBase() *pb.ClientBase
	GetServerBase() *pb.ServerBase
}

type Context struct {
	context.Context
	appInstance    Application
	loggerInstance *logrus.Logger
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

func (c *Context) Background() *Context {
	return NewContext(context.Background(), c.appInstance)
}

func (c *Context) Copy() *Context {
	return NewContext(c.Context, c.appInstance)
}

func (c *Context) CopyWithCancel() (*Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(c.Context)
	return NewContext(ctx, c.appInstance), cancel
}

func (c *Context) BackgroundWithCancel() (*Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	return NewContext(ctx, c.appInstance), cancel
}

func (c *Context) Logger() *logrus.Entry {
	if c.loggerInstance != nil {
		return c.loggerInstance.WithContext(c)
	}
	return c.GetApp().Logger(c)
}

func (c *Context) SetLogger(logger *logrus.Logger) {
	c.loggerInstance = logger
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
