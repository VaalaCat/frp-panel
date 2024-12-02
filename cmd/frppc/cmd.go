package main

import (
	"context"

	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/logger"
	"github.com/spf13/cobra"
)

var (
	clientCmd *cobra.Command
	rootCmd   *cobra.Command
)

func initCommand() {
	var (
		clientSecret string
		clientID     string
		rpcHost      string
		appSecret    string
		rpcPort      int
		apiPort      int
		apiScheme    string
	)

	rootCmd = &cobra.Command{
		Use:   "frp-panel-client [-s client secret] [-i client id] [-a app secret] [-r rpc host] [-c rpc port] [-p api port]",
		Short: "frp-panel-client is a frp panel client QwQ",
		Run: func(cmd *cobra.Command, args []string) {
			patchConfig(rpcHost, appSecret,
				clientID, clientSecret,
				apiScheme, rpcPort, apiPort)
			runClient()
		},
	}
	rootCmd.Flags().StringVarP(&clientSecret, "secret", "s", "", "client secret")
	rootCmd.Flags().StringVarP(&clientID, "id", "i", "", "client id")
	rootCmd.Flags().StringVarP(&rpcHost, "rpc", "r", "", "rpc host")
	rootCmd.Flags().StringVarP(&appSecret, "app", "a", "", "app secret")

	rootCmd.Flags().IntVarP(&rpcPort, "rpc-port", "c", 0, "rpc port")
	rootCmd.Flags().IntVarP(&apiPort, "api-port", "p", 0, "api port")

	rootCmd.Flags().StringVarP(&apiScheme, "api-scheme", "e", "", "api scheme")
}

func initLogger() {
	logger.Instance().SetReportCaller(true)
}

func patchConfig(host, secret, clientID, clientSecret, apiScheme string, rpcPort, apiPort int) {
	if len(host) != 0 {
		conf.Get().Master.RPCHost = host
		conf.Get().Master.APIHost = host
	}
	if len(secret) != 0 {
		conf.Get().App.Secret = secret
	}
	if rpcPort != 0 {
		conf.Get().Master.RPCPort = rpcPort
	}
	if apiPort != 0 {
		conf.Get().Master.APIPort = apiPort
	}
	if len(apiScheme) != 0 {
		conf.Get().Master.APIScheme = apiScheme
	}
	if len(clientID) != 0 {
		conf.Get().Client.ID = clientID
	}
	if len(clientSecret) != 0 {
		conf.Get().Client.Secret = clientSecret
	}
	c := context.Background()
	logger.Logger(c).Infof("env config rpc host: %s, rpc port: %d, api host: %s, api port: %d, api scheme: %s",
		conf.Get().Master.RPCHost, conf.Get().Master.RPCPort,
		conf.Get().Master.APIHost, conf.Get().Master.APIPort,
		conf.Get().Master.APIScheme)
}
