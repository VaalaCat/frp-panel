package shared

import (
	"go.uber.org/fx"
)

var (
	clientMod = fx.Module("cmd.client",
		fx.Provide(
			fx.Annotate(NewWatcher, fx.ResultTags(`name:"clientTaskManager"`)),
		))

	serverMod = fx.Module("cmd.server", fx.Provide(
		fx.Annotate(NewServerAPI, fx.ResultTags(`name:"serverApiService"`)),
		fx.Annotate(NewServerRouter, fx.ResultTags(`name:"serverRouter"`)),
		fx.Annotate(NewWatcher, fx.ResultTags(`name:"serverTaskManager"`)),
	))

	masterMod = fx.Module("cmd.master", fx.Provide(
		NewPermissionManager,
		NewEnforcer,
		NewListenerOptions,
		NewDBManager,
		NewWSListener,
		NewMasterTLSConfig,
		NewWSUpgrader,
		NewClientLogManager,
		fx.Annotate(NewWatcher, fx.ResultTags(`name:"masterTaskManager"`)),
		fx.Annotate(NewMasterRouter, fx.ResultTags(`name:"masterRouter"`)),
		fx.Annotate(NewHTTPMasterService, fx.ResultTags(`name:"httpMasterService"`)),
		fx.Annotate(NewHTTPMasterService, fx.ResultTags(`name:"wsMasterService"`)),
		fx.Annotate(NewTLSMasterService, fx.ResultTags(`name:"tlsMasterService"`)),
		fx.Annotate(NewMux, fx.ResultTags(`name:"tlsMux"`)),
		fx.Annotate(NewHTTPMux, fx.ResultTags(`name:"httpMux"`)),
		fx.Annotate(NewWSGrpcHandler, fx.ResultTags(`name:"wsGrpcHandler"`)),
	))

	commonMod = fx.Module("common", fx.Provide(
		NewLogHookManager,
		NewPTYManager,
		NewBaseApp,
		NewContext,
		NewClientsManager,
		NewAutoJoin,
		fx.Annotate(NewPatchedConfig, fx.ResultTags(`name:"argsPatchedConfig"`)),
	))
)
