package dao

import "github.com/VaalaCat/frp-panel/services/app"

type Query interface {
	CertQuery
	ClientQuery
	EndpointQuery
	LinkQuery
	NetworkQuery
	ProxyQuery
	ServerQuery
	StatsQuery
	UserQuery
	WireGuardQuery
	WorkerQuery
}

type Mutation interface {
	CertMutation
	ClientMutation
	EndpointMutation
	LinkMutation
	NetworkMutation
	ProxyMutation
	ServerMutation
	StatsMutation
	UserMutation
	WireGuardMutation
	WorkerMutation
	UserGroupMutation
}

// queryImpl/mutationImpl 是具体表级实现的基础结构（持有 ctx）。
type queryImpl struct {
	ctx *app.Context
}

type mutationImpl struct {
	ctx *app.Context
}

// compositeQuery / compositeMutation 组合各子领域实现，对外暴露统一入口。
type compositeQuery struct {
	CertQuery
	ClientQuery
	EndpointQuery
	LinkQuery
	NetworkQuery
	ProxyQuery
	ServerQuery
	StatsQuery
	UserQuery
	WireGuardQuery
	WorkerQuery
}

type compositeMutation struct {
	CertMutation
	ClientMutation
	EndpointMutation
	LinkMutation
	NetworkMutation
	ProxyMutation
	ServerMutation
	StatsMutation
	UserMutation
	WireGuardMutation
	WorkerMutation
	UserGroupMutation
}

func NewQuery(ctx *app.Context) Query {
	base := &queryImpl{ctx: ctx}
	return &compositeQuery{
		CertQuery:      newCertQuery(base),
		ClientQuery:    newClientQuery(base),
		EndpointQuery:  newEndpointQuery(base),
		LinkQuery:      newLinkQuery(base),
		NetworkQuery:   newNetworkQuery(base),
		ProxyQuery:     newProxyQuery(base),
		ServerQuery:    newServerQuery(base),
		StatsQuery:     newStatsQuery(base),
		UserQuery:      newUserQuery(base),
		WireGuardQuery: newWireGuardQuery(base),
		WorkerQuery:    newWorkerQuery(base),
	}
}

func NewMutation(ctx *app.Context) Mutation {
	base := &mutationImpl{ctx: ctx}
	return &compositeMutation{
		CertMutation:      newCertMutation(base),
		ClientMutation:    newClientMutation(base),
		EndpointMutation:  newEndpointMutation(base),
		LinkMutation:      newLinkMutation(base),
		NetworkMutation:   newNetworkMutation(base),
		ProxyMutation:     newProxyMutation(base),
		ServerMutation:    newServerMutation(base),
		StatsMutation:     newStatsMutation(base),
		UserMutation:      newUserMutation(base),
		WireGuardMutation: newWireGuardMutation(base),
		WorkerMutation:    newWorkerMutation(base),
		UserGroupMutation: newUserGroupMutation(base),
	}
}

var (
	_ Query    = (*compositeQuery)(nil)
	_ Mutation = (*compositeMutation)(nil)
)
