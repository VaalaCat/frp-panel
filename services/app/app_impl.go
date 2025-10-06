package app

import (
	"context"
	"sync"

	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/utils/logger"
	"github.com/casbin/casbin/v2"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/credentials"
)

type application struct {
	streamLogHookMgr StreamLogHookMgr
	masterCli        MasterClient

	shellPTYMgr          ShellPTYMgr
	clientLogManager     ClientLogManager
	clientRPCHandler     ClientRPCHandler
	dbManager            DBManager
	clientController     ClientController
	clientRecvMap        *sync.Map
	clientsManager       ClientsManager
	serverHandler        ServerHandler
	serverController     ServerController
	rpcCred              credentials.TransportCredentials
	conf                 conf.Config
	currentRole          string
	permManager          PermissionManager
	enforcer             *casbin.Enforcer
	workerExecManager    WorkerExecManager
	workersManager       WorkersManager
	wireGuardManager     WireGuardManager
	networkTopologyCache NetworkTopologyCache

	loggerInstance *logrus.Logger
}

func (a *application) GetClientBase() *pb.ClientBase {
	return &pb.ClientBase{
		ClientId:     a.GetConfig().Client.ID,
		ClientSecret: a.GetConfig().Client.Secret,
	}
}

func (a *application) GetServerBase() *pb.ServerBase {
	return &pb.ServerBase{
		ServerId:     a.GetConfig().Client.ID,
		ServerSecret: a.GetConfig().Client.Secret,
	}
}

func (a *application) SetLogger(l *logrus.Logger) {
	a.loggerInstance = l
}

func (a *application) Logger(ctx context.Context) *logrus.Entry {
	if a.loggerInstance == nil {
		return logger.Logger(ctx)
	}
	return a.loggerInstance.WithContext(ctx)
}

// GetWorkersManager implements Application.
func (a *application) GetWorkersManager() WorkersManager {
	return a.workersManager
}

// SetWorkersManager implements Application.
func (a *application) SetWorkersManager(w WorkersManager) {
	a.workersManager = w
}

// GetWorkerExecManager implements Application.
func (a *application) GetWorkerExecManager() WorkerExecManager {
	return a.workerExecManager
}

// SetWorkerExecManager implements Application.
func (a *application) SetWorkerExecManager(w WorkerExecManager) {
	a.workerExecManager = w
}

// GetEnforcer implements Application.
func (a *application) GetEnforcer() *casbin.Enforcer {
	return a.enforcer
}

// SetEnforcer implements Application.
func (a *application) SetEnforcer(c *casbin.Enforcer) {
	a.enforcer = c
}

// GetPermManager implements Application.
func (a *application) GetPermManager() PermissionManager {
	return a.permManager
}

// SetPermManager implements Application.
func (a *application) SetPermManager(p PermissionManager) {
	a.permManager = p
}

// GetCurrentRole implements Application.
func (a *application) GetCurrentRole() string {
	return a.currentRole
}

// SetCurrentRole implements Application.
func (a *application) SetCurrentRole(r string) {
	a.currentRole = r
}

// GetClientCred implements Application.
func (a *application) GetRPCCred() credentials.TransportCredentials {
	return a.rpcCred
}

// SetClientCred implements Application.
func (a *application) SetRPCCred(cred credentials.TransportCredentials) {
	a.rpcCred = cred
}

// GetConfig implements Application.
func (a *application) GetConfig() conf.Config {
	return a.conf
}

// SetConfig implements Application.
func (a *application) SetConfig(c conf.Config) {
	a.conf = c
}

// GetServerController implements Application.
func (a *application) GetServerController() ServerController {
	return a.serverController
}

// SetServerController implements Application.
func (a *application) SetServerController(serverController ServerController) {
	a.serverController = serverController
}

// GetClientController implements Application.
func (a *application) GetClientController() ClientController {
	return a.clientController
}

// SetClientController implements Application.
func (a *application) SetClientController(clientController ClientController) {
	a.clientController = clientController
}

// GetServerHandler implements Application.
func (a *application) GetServerHandler() ServerHandler {
	return a.serverHandler
}

// SetServerHandler implements Application.
func (a *application) SetServerHandler(serverHandler ServerHandler) {
	a.serverHandler = serverHandler
}

// GetClientRPCHandler implements Application.
func (a *application) GetClientRPCHandler() ClientRPCHandler {
	return a.clientRPCHandler
}

// SetClientRPCHandler implements Application.
func (a *application) SetClientRPCHandler(clientRPCHandler ClientRPCHandler) {
	a.clientRPCHandler = clientRPCHandler
}

// GetMasterCli implements Application.
func (a *application) GetMasterCli() MasterClient {
	return a.masterCli
}

// SetMasterCli implements Application.
func (a *application) SetMasterCli(masterCli MasterClient) {
	a.masterCli = masterCli
}

// GetClientsManager implements Application.
func (a *application) GetClientsManager() ClientsManager {
	return a.clientsManager
}

// SetClientsManager implements Application.
func (a *application) SetClientsManager(clientsManager ClientsManager) {
	a.clientsManager = clientsManager
}

// GetClientRecvMap implements Application.
func (a *application) GetClientRecvMap() *sync.Map {
	return a.clientRecvMap
}

// SetClientRecvMap implements Application.
func (a *application) SetClientRecvMap(clientRecvMap *sync.Map) {
	a.clientRecvMap = clientRecvMap
}

// GetDBManager implements Application.
func (a *application) GetDBManager() DBManager {
	return a.dbManager
}

// SetDBManager implements Application.
func (a *application) SetDBManager(dbManager DBManager) {
	a.dbManager = dbManager
}

// GetClientLogManager implements Application.
func (a *application) GetClientLogManager() ClientLogManager {
	return a.clientLogManager
}

// SetClientLogManager implements Application.
func (a *application) SetClientLogManager(clientLogManager ClientLogManager) {
	a.clientLogManager = clientLogManager
}

// GetShellPTYMgr implements Application.
func (a *application) GetShellPTYMgr() ShellPTYMgr {
	return a.shellPTYMgr
}

// GetStreamLogHookMgr implements Application.
func (a *application) GetStreamLogHookMgr() StreamLogHookMgr {
	return a.streamLogHookMgr
}

// GetWireGuardManager implements Application.
func (a *application) GetWireGuardManager() WireGuardManager {
	return a.wireGuardManager
}

// SetWireGuardManager implements Application.
func (a *application) SetWireGuardManager(wireGuardManager WireGuardManager) {
	a.wireGuardManager = wireGuardManager
}

// SetShellPTYMgr implements Application.
func (a *application) SetShellPTYMgr(shellPTYMgr ShellPTYMgr) {
	a.shellPTYMgr = shellPTYMgr
}

// SetStreamLogHookMgr implements Application.
func (a *application) SetStreamLogHookMgr(streamLogHookMgr StreamLogHookMgr) {
	a.streamLogHookMgr = streamLogHookMgr
}

// GetNetworkTopologyCache implements Application.
func (a *application) GetNetworkTopologyCache() NetworkTopologyCache {
	return a.networkTopologyCache
}

// SetNetworkTopologyCache implements Application.
func (a *application) SetNetworkTopologyCache(networkTopologyCache NetworkTopologyCache) {
	a.networkTopologyCache = networkTopologyCache
}
