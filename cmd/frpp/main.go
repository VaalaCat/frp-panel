package main

import (
	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/rpc"
)

func main() {
	initLogger()
	initCommand()
	conf.InitConfig()
	rpc.InitRPCClients()

	rootCmd.Execute()
}
