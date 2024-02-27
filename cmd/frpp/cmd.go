package main

import (
	"github.com/VaalaCat/frp-panel/conf"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	clientCmd *cobra.Command
	serverCmd *cobra.Command
	masterCmd *cobra.Command
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

	clientCmd = &cobra.Command{
		Use:   "client [-s client secret] [-i client id] [-a app secret] [-r rpc host] [-c rpc port] [-p api port]",
		Short: "run managed frpc",
		Run: func(cmd *cobra.Command, args []string) {
			patchConfig(rpcHost, appSecret,
				clientID, clientSecret,
				apiScheme, rpcPort, apiPort)
			runClient()
		},
	}
	serverCmd = &cobra.Command{
		Use:   "server [-s client secret] [-i client id] [-a app secret] [-r rpc host] [-c rpc port] [-p api port]",
		Short: "run managed frps",
		Run: func(cmd *cobra.Command, args []string) {
			patchConfig(rpcHost, appSecret,
				clientID, clientSecret,
				apiScheme, rpcPort, apiPort)
			runServer()
		},
	}
	masterCmd = &cobra.Command{
		Use:   "master",
		Short: "run frp-panel manager",
		Run: func(cmd *cobra.Command, args []string) {
			runMaster()
		},
	}
	rootCmd = &cobra.Command{
		Use:   "frp-panel",
		Short: "frp-panel is a frp panel QwQ",
	}
	rootCmd.AddCommand(clientCmd, serverCmd, masterCmd)
	clientCmd.Flags().StringVarP(&clientSecret, "secret", "s", "", "client secret")
	serverCmd.Flags().StringVarP(&clientSecret, "secret", "s", "", "client secret")
	clientCmd.Flags().StringVarP(&clientID, "id", "i", "", "client id")
	serverCmd.Flags().StringVarP(&clientID, "id", "i", "", "client id")
	clientCmd.Flags().StringVarP(&rpcHost, "rpc", "r", "", "rpc host")
	serverCmd.Flags().StringVarP(&rpcHost, "rpc", "r", "", "rpc host")
	clientCmd.Flags().StringVarP(&appSecret, "app", "a", "", "app secret")
	serverCmd.Flags().StringVarP(&appSecret, "app", "a", "", "app secret")

	clientCmd.Flags().IntVarP(&rpcPort, "rpc-port", "c", 0, "rpc port")
	serverCmd.Flags().IntVarP(&rpcPort, "rpc-port", "c", 0, "rpc port")
	clientCmd.Flags().IntVarP(&apiPort, "api-port", "p", 0, "api port")
	serverCmd.Flags().IntVarP(&apiPort, "api-port", "p", 0, "api port")

	clientCmd.Flags().StringVarP(&apiScheme, "api-scheme", "e", "", "api scheme")
	serverCmd.Flags().StringVarP(&apiScheme, "api-scheme", "e", "", "api scheme")
}

func initLogger() {
	logrus.SetReportCaller(true)
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
	logrus.Infof("env config rpc host: %s, rpc port: %d, api host: %s, api port: %d, api scheme: %s",
		conf.Get().Master.RPCHost, conf.Get().Master.RPCPort,
		conf.Get().Master.APIHost, conf.Get().Master.APIPort,
		conf.Get().Master.APIScheme)
}
