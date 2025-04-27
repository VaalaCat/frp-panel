package main

import (
	"github.com/VaalaCat/frp-panel/utils/logger"
	"github.com/fatedier/golib/crypto"
	"github.com/spf13/cobra"
)

func main() {
	crypto.DefaultSalt = "frp"
	logger.InitLogger()
	cobra.MousetrapHelpText = ""
	rootCmd := buildCommand()
	setMasterCommandIfNonePresent(rootCmd)
	rootCmd.Execute()
}
