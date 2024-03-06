package main

import (
	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/rpc"
	"github.com/spf13/cobra"
)

func main() {
	cobra.MousetrapHelpText = ""

	initLogger()
	initCommand()
	conf.InitConfig()
	rpc.InitRPCClients()

	rootCmd.Execute()
}
