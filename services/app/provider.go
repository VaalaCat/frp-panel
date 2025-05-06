package app

import (
	"sync"
	"time"

	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/casbin/casbin/v2"

	"github.com/fatedier/frp/client/proxy"
	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/fatedier/frp/pkg/metrics/mem"
	"gorm.io/gorm"
)

// biz/common/stream_log.go
type StreamLogHookMgr interface {
	AddStream(send func(msg string), closeSend func())
	Close()
	Lock()
	TryLock() bool
	Unlock()
}

// utils/sync.go
type SyncMap[K comparable, V any] interface {
	Clone() utils.SyncMap[K, V]
	Delete(key K)
	Grow(size int)
	Keys() []K
	Len() (l int)
	Load(key K) (value V, loaded bool)
	LoadAndDelete(key K) (value V, loaded bool)
	LoadOrStore(key K, value V) (actual V, loaded bool)
	Range(f func(key K, value V) (shouldContinue bool))
	Store(key K, value V)
	Values() []V
}

type GoSyncMap interface {
	Clear()
	CompareAndDelete(key any, old any) (deleted bool)
	CompareAndSwap(key any, old any, new any) (swapped bool)
	Delete(key any)
	Load(key any) (value any, ok bool)
	LoadAndDelete(key any) (value any, loaded bool)
	LoadOrStore(key any, value any) (actual any, loaded bool)
	Range(f func(key any, value any) bool)
	Store(key any, value any)
	Swap(key any, value any) (previous any, loaded bool)
}

// biz/master/shell/mgr.go
type ShellPTYMgr interface {
	SyncMap[string, pb.Master_PTYConnectServer]
	Add(sessionID string, conn pb.Master_PTYConnectServer)
	IsSessionDone(sessionID string) bool
	SetSessionDone(sessionID string)
}

// biz/master/streamlog/collect_log.go
type ClientLogManager interface {
	SyncMap[string, chan string]
	GetClientLock(clientId string) *sync.Mutex
}

// models/db.go
type DBManager interface {
	GetDB(dbType string, dbRole string) *gorm.DB
	GetDefaultDB() *gorm.DB
	SetDB(dbType string, dbRole string, db *gorm.DB)
	RemoveDB(dbType string, dbRole string)
	SetDebug(bool)
	Init()
}

type ClientsManager interface {
	Get(cliID string) *Connector
	Set(cliID, clientType string, sender pb.Master_ServerSendServer)
	Remove(cliID string)
	ClientAddr(cliID string) string
	ConnectTime(cliID string) (time.Time, bool)
}

type Connector struct {
	CliID   string
	Conn    pb.Master_ServerSendServer
	CliType string
}

type Service interface {
	Run()
	Stop()
}

// services/client/frpc_service.go
type ClientHandler interface {
	Run()
	Stop()
	Wait()
	Running() bool
	Update([]v1.ProxyConfigurer, []v1.VisitorConfigurer)
	AddProxy(v1.ProxyConfigurer)
	AddVisitor(v1.VisitorConfigurer)
	RemoveProxy(v1.ProxyConfigurer)
	RemoveVisitor(v1.VisitorConfigurer)
	GetProxyStatus(string) (*proxy.WorkingStatus, bool)
	GetCommonCfg() *v1.ClientCommonConfig
	GetProxyCfgs() map[string]v1.ProxyConfigurer
	GetVisitorCfgs() map[string]v1.VisitorConfigurer
}

// services/rpcclient/rpc_service.go
type ClientRPCHandler interface {
	Run()
	Stop()
	GetCli() MasterClient
}

type ClientController interface {
	Add(clientID string, serverID string, clientHandler ClientHandler)
	Get(clientID string, serverID string) ClientHandler
	Delete(clientID string, serverID string)
	Set(clientID string, serverID string, clientHandler ClientHandler)
	Run(clientID string, serverID string) // 不阻塞
	Stop(clientID string, serverID string)
	GetByClient(clientID string) *utils.SyncMap[string, ClientHandler]
	DeleteByClient(clientID string)
	RunByClient(clientID string) // 不阻塞
	StopByClient(clientID string)
	StopAll()
	DeleteAll()
	RunAll()
	List() []string
}

type ServerController interface {
	Add(serverID string, serverHandler ServerHandler)
	Get(serverID string) ServerHandler
	Delete(serverID string)
	Set(serverID string, serverHandler ServerHandler)
	Run(serverID string) // 不阻塞
	Stop(serverID string)
	List() []string
}

type ServerHandler interface {
	Run()
	Stop()
	IsFirstSync() bool
	GetCommonCfg() *v1.ServerConfig
	GetMem() *mem.ServerStats
	GetProxyStatsByType(v1.ProxyType) []*mem.ProxyStats
}

// rpc/master.go
type MasterClient interface {
	Call() pb.MasterClient
}

// services/rbac/perm_manager.go
type PermissionManager interface {
	AddUserToGroup(userID int, groupID string, tenantID int) (bool, error)
	CheckPermission(userID int, objType defs.RBACObj, objID string, action defs.RBACAction, tenantID int) (bool, error)
	Enforcer() *casbin.Enforcer
	GrantGroupPermission(groupID string, objType defs.RBACObj, objID string, action defs.RBACAction, tenantID int) (bool, error)
	GrantUserPermission(userID int, objType defs.RBACObj, objID string, action defs.RBACAction, tenantID int) (bool, error)
	RemoveUserFromGroup(userID int, groupID string, tenantID int) (bool, error)
	RevokeGroupPermission(groupID string, objType defs.RBACObj, objID string, action defs.RBACAction, tenantID int) (bool, error)
	RevokeUserPermission(userID int, objType defs.RBACObj, objID string, action defs.RBACAction, tenantID int) (bool, error)
}

// services/workerd/exec_manager.go
type WorkerExecManager interface {
	RunCmd(workerId string, cwd string, argv []string)
	ExitCmd(workerId string)
	ExitAllCmd()
	UpdateBinaryPath(path string)
}

// services/workerd/workerd.go
type WorkerController interface {
	RunWorker(c *Context)
	StopWorker(c *Context)
	// GetWorkerStatus(c *Context) defs.WorkerStatus
	GarbageCollect()
	Init(c *Context) error
}

// services/workerd/workers_manager.go
type WorkersManager interface {
	GetWorker(ctx *Context, id string) (WorkerController, bool)
	RunWorker(ctx *Context, id string, worker WorkerController) error
	StopWorker(ctx *Context, id string) error
	GetWorkerStatus(ctx *Context, id string) (defs.WorkerStatus, error)
	// install workerd bin to workerd bin path, if not specified, use default path /usr/local/bin/workerd
	InstallWorkerd(ctx *Context, url string, path string) (string, error)
}
