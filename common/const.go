package common

import "time"

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

const (
	DefaultServerID    = "default"
	DefaultAdminUserID = 1
)

const (
	PullConfigDuration    = 30 * time.Second
	PushProxyInfoDuration = 30 * time.Second
)
