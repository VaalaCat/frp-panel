package app

import (
	"sync"

	"github.com/VaalaCat/frp-panel/conf"
	"google.golang.org/grpc/credentials"
)

type application struct {
	streamLogHookMgr StreamLogHookMgr
	masterCli        MasterClient

	shellPTYMgr      ShellPTYMgr
	clientLogManager ClientLogManager
	clientRPCHandler ClientRPCHandler
	dbManager        DBManager
	clientController ClientController
	clientRecvMap    *sync.Map
	clientsManager   ClientsManager
	serverHandler    ServerHandler
	serverController ServerController
	rpcCred          credentials.TransportCredentials
	conf             conf.Config
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

// SetShellPTYMgr implements Application.
func (a *application) SetShellPTYMgr(shellPTYMgr ShellPTYMgr) {
	a.shellPTYMgr = shellPTYMgr
}

// SetStreamLogHookMgr implements Application.
func (a *application) SetStreamLogHookMgr(streamLogHookMgr StreamLogHookMgr) {
	a.streamLogHookMgr = streamLogHookMgr
}
