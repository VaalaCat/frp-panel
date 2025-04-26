package main

import (
	"github.com/VaalaCat/frp-panel/utils/logger"
	"github.com/spf13/cobra"
)

func main() {
	logger.InitLogger()
	cobra.MousetrapHelpText = ""
	rootCmd := buildCommand()
	setMasterCommandIfNonePresent(rootCmd)
	rootCmd.Execute()
}
