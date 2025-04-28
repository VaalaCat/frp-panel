package defs

type RBACAction string

const (
	RBACActionCreate RBACAction = "create"
	RBACActionRead   RBACAction = "read"
	RBACActionUpdate RBACAction = "update"
	RBACActionDelete RBACAction = "delete"
)

type RBACObj string

const (
	RBACObjServer RBACObj = "server"
	RBACObjClient RBACObj = "client"
	RBACObjUser   RBACObj = "user"
	RBACObjGroup  RBACObj = "group"
	RBACObjAPI    RBACObj = "api"
)

type RBACSubject string

const (
	RBACSubjectUser  RBACSubject = "user"
	RBACSubjectGroup RBACSubject = "group"
	RBACSubjectToken RBACSubject = "token"
)

type RBACDomain string

const (
	RBACDomainTenant RBACDomain = "tenant"
)

type APIPermission struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

const (
	TokenPayloadKey_Permissions = "permissions"
)
