package main

import (
	"github.com/VaalaCat/frp-panel/cmd/frpp/shared"
	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/utils/logger"
	"github.com/fatedier/golib/crypto"
	"github.com/spf13/cobra"
)

func main() {
	crypto.DefaultSalt = "frp"
	logger.InitLogger()
	cobra.MousetrapHelpText = ""

	rootCmd := shared.NewRootCmd(
		shared.NewClientCmd(conf.NewConfig()),
		shared.NewJoinCmd(),
		shared.NewInstallServiceCmd(),
		shared.NewUninstallServiceCmd(),
		shared.NewStartServiceCmd(),
		shared.NewStopServiceCmd(),
		shared.NewRestartServiceCmd(),
		shared.NewVersionCmd(),
	)

	shared.SetClientCommandIfNonePresent(rootCmd)
	rootCmd.Execute()
}
