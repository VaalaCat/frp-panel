package main

import (
	"sync"

	"github.com/VaalaCat/frp-panel/biz/common"
	"github.com/VaalaCat/frp-panel/biz/master/shell"
	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/rpc"
	"github.com/VaalaCat/frp-panel/utils/logger"
	"github.com/fatedier/golib/crypto"
	"github.com/spf13/cobra"
)

func main() {
	crypto.DefaultSalt = "frp"

	logger.InitLogger()
	cobra.MousetrapHelpText = ""
	cfg := conf.NewConfig()

	appInstance := app.NewApp()
	appInstance.SetConfig(cfg)
	appInstance.SetClientsManager(rpc.NewClientsManager())
	appInstance.SetStreamLogHookMgr(&common.HookMgr{})
	appInstance.SetShellPTYMgr(shell.NewPTYMgr())
	appInstance.SetClientRecvMap(&sync.Map{})

	initCommand(appInstance)

	setMasterCommandIfNonePresent()
	rootCmd.Execute()
}
