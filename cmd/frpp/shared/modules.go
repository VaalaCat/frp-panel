package shared

import (
	bizmaster "github.com/VaalaCat/frp-panel/biz/master"
	"github.com/VaalaCat/frp-panel/biz/master/shell"
	"github.com/VaalaCat/frp-panel/biz/master/streamlog"
	bizserver "github.com/VaalaCat/frp-panel/biz/server"
	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/services/rpc"
	"github.com/VaalaCat/frp-panel/utils/logger"
	"go.uber.org/fx"
)

var (
	clientMod = fx.Module("cmd.client",
		fx.Provide(
			NewWorkerExecManager,
			NewWorkersManager,
			fx.Annotate(NewWatcher, fx.ResultTags(`name:"clientTaskManager"`)),
		))

	serverMod = fx.Module("cmd.server", fx.Provide(
		fx.Annotate(NewServerAPI, fx.ResultTags(`name:"serverApiService"`)),
		fx.Annotate(bizserver.NewRouter, fx.ResultTags(`name:"serverRouter"`)),
		fx.Annotate(NewWatcher, fx.ResultTags(`name:"serverTaskManager"`)),
	))

	masterMod = fx.Module("cmd.master", fx.Provide(
		NewPermissionManager,
		NewEnforcer,
		conf.GetListener,
		NewDBManager,
		NewWSListener,
		NewMasterTLSConfig,
		NewWSUpgrader,
		streamlog.NewClientLogManager,
		// wireguard.NewWireGuardManager,
		fx.Annotate(NewWatcher, fx.ResultTags(`name:"masterTaskManager"`)),
		fx.Annotate(bizmaster.NewRouter, fx.ResultTags(`name:"masterRouter"`)),
		fx.Annotate(NewHTTPMasterService, fx.ResultTags(`name:"httpMasterService"`)),
		fx.Annotate(NewHTTPMasterService, fx.ResultTags(`name:"wsMasterService"`)),
		fx.Annotate(NewTLSMasterService, fx.ResultTags(`name:"tlsMasterService"`)),
		fx.Annotate(NewMux, fx.ResultTags(`name:"tlsMux"`)),
		fx.Annotate(NewHTTPMux, fx.ResultTags(`name:"httpMux"`)),
		fx.Annotate(NewWSGrpcHandler, fx.ResultTags(`name:"wsGrpcHandler"`)),
	))

	commonMod = fx.Module("common", fx.Provide(
		logger.Logger,
		logger.Instance,
		NewLogHookManager,
		shell.NewPTYMgr,
		NewBaseApp,
		NewContext,
		NewAndFinishNormalContext,
		rpc.NewClientsManager,
		NewAutoJoin, // provide final config
		fx.Annotate(NewPatchedConfig, fx.ResultTags(`name:"argsPatchedConfig"`)),
	))
)
