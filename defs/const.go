package defs

import (
	"time"
)

const (
	AuthorizationKey    = "authorization"
	SetAuthorizationKey = "x-set-authorization"
	MsgKey              = "msg"
	UserIDKey           = "x-vaala-userid"
	EndpointKey         = "endpoint"
	IpAddrKey           = "ipaddr"
	HeaderKey           = "header"
	MethodKey           = "method"
	UAKey               = "User-Agent"
	ContentTypeKey      = "Content-Type"
	TraceIDKey          = "TraceID"
	TokenKey            = "token"
	FRPAuthTokenKey     = "token"
	ErrKey              = "err"
	UserInfoKey         = "x-vaala-userinfo"
	FRPClientIDKey      = "x-vaala-frp-client-id"
)

const (
	ErrInvalidRequest   = "invalid request"
	ErrUserInfoNotValid = "user info not valid"
	ErrInternalError    = "internal error"
	ErrParamNotValid    = "param not valid"
	ErrDB               = "database error"
	ErrNotFound         = "data not found"
	ErrCodeNotFound     = "code not found"
	ErrCodeAlreadyUsed  = "code already used"
)

const (
	ReqSuccess = "success"
)

const (
	TimeLayout = time.RFC3339
)

const (
	StatusPending = "pending"
	StatusSuccess = "success"
	StatusFailed  = "failed"
	StatusDone    = "done"
)

const (
	CliTypeServer = "server"
	CliTypeClient = "client"
)

type AppRole string

const (
	AppRole_Client AppRole = CliTypeClient
	AppRole_Server AppRole = CliTypeServer
	AppRole_Master AppRole = "master"
)

const (
	DefaultServerID    = "default"
	DefaultAdminUserID = 1
)

const (
	PullConfigDuration        = 30 * time.Second
	PushProxyInfoDuration     = 30 * time.Second
	PullClientWorkersDuration = 30 * time.Second
)

const (
	CurEnvPath         = ".env"
	SysEnvPath         = "/etc/frpp/.env"
	EnvClientID        = "CLIENT_ID"
	EnvClientSecret    = "CLIENT_SECRET"
	EnvMasterRPCHost   = "MASTER_RPC_HOST"
	EnvMasterAPIHost   = "MASTER_API_HOST"
	EnvMasterRPCPort   = "MASTER_RPC_PORT"
	EnvMasterAPIPort   = "MASTER_API_PORT"
	EnvMasterAPIScheme = "MASTER_API_SCHEME"
	EnvClientAPIUrl    = "CLIENT_API_URL"
	EnvClientRPCUrl    = "CLIENT_RPC_URL"
)

const (
	DBRoleDefault = "default"
	DBRoleRam     = "ram"
)

const (
	DBTypeSQLite3  = "sqlite3"
	DBTypeMysql    = "mysql"
	DBTypePostgres = "postgres"
)

const (
	UserRole_Admin  = "admin"
	UserRole_Normal = "normal"
	CapFileName     = "workerd.capnp"
	WorkerInfoPath  = "workers"
	WorkerCodePath  = "src"
	DBTypeSqlite    = "sqlite"

	DefaultHostName       = "127.0.0.1"
	DefaultNodeName       = "default"
	DefaultExternalPath   = "/"
	DefaultEntry          = "entry.js"
	DefaultSocketTemplate = "unix-abstract:/tmp/frpp-worker-%s.sock"
	DefaultCode           = `export default {
  async fetch(req, env) {
    try {
		let resp = new Response("worker: " + req.url + " is online! -- " + new Date())
		return resp
	} catch(e) {
		return new Response(e.stack, { status: 500 })
	}
  }
};`

	DefaultConfigTemplate = `using Workerd = import "/workerd/workerd.capnp";

const config :Workerd.Config = (
  services = [
    (name = "{{.WorkerId}}", worker = .v{{.WorkerId}}Worker),
  ],

  sockets = [
    (
      name = "{{.WorkerId}}",
      address = "{{.Socket.Address}}",
      http=(),
      service="{{.WorkerId}}"
    ),
  ]
);

const v{{.WorkerId}}Worker :Workerd.Worker = (
  modules = [
    (name = "{{.CodeEntry}}", esModule = embed "src/{{.CodeEntry}}"),
  ],
  compatibilityDate = "2023-04-03",
);`
)

type TokenStatus string

const (
	TokenStatusActive   TokenStatus = "active"
	TokenStatusInactive TokenStatus = "inactive"
	TokenStatusRevoked  TokenStatus = "revoked"
)

const (
	KeyNodeName    = "node_name"
	KeyNodeSecret  = "node_secret"
	KeyNodeProto   = "node_proto"
	KeyWorkerProto = "worker_proto"
)

type WorkerStatus string

const (
	WorkerStatus_Unknown  WorkerStatus = "unknown"
	WorkerStatus_Running  WorkerStatus = "running"
	WorkerStatus_Inactive WorkerStatus = "inactive"
)

const (
	FrpProxyAnnotationsKey_Ingress  = "ingress"
	FrpProxyAnnotationsKey_WorkerId = "worker_id"
)
