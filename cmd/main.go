package main

import (
	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/rpc"
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/fatedier/golib/crypto"
)

func main() {
	initLogger()
	initCommand()
	conf.InitConfig()
	rpc.InitRPCClients()

	crypto.DefaultSalt = utils.MD5(conf.Get().App.GlobalSecret)
	rootCmd.Execute()
}
