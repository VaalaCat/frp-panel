package main

import (
	"fmt"
	"os"

	"github.com/VaalaCat/frp-panel/conf"
	"github.com/sirupsen/logrus"
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
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	clientCmd = &cobra.Command{
		Use:   "client [-s client secret] [-i client id] [-a app secret] [-r rpc host] [-c rpc port] [-p api port]",
		Short: "run managed frpc",
		Run: func(cmd *cobra.Command, args []string) {
			patchConfig(rpcHost, appSecret, rpcPort, apiPort)
			runClient(clientID, clientSecret)
		},
	}
	rootCmd = &cobra.Command{
		Use:   "frp-panel-client",
		Short: "frp-panel-client is a frp panel client QwQ",
	}
	rootCmd.AddCommand(clientCmd)
	clientCmd.Flags().StringVarP(&clientSecret, "secret", "s", "", "client secret")
	clientCmd.Flags().StringVarP(&clientID, "id", "i", hostname, "client id")
	clientCmd.Flags().StringVarP(&rpcHost, "rpc", "r", "", "rpc host")
	clientCmd.Flags().StringVarP(&appSecret, "app", "a", "", "app secret")

	clientCmd.Flags().IntVarP(&rpcPort, "rpc-port", "c", 0, "rpc port")
	clientCmd.Flags().IntVarP(&apiPort, "api-port", "p", 0, "api port")

	clientCmd.Flags().StringVarP(&apiScheme, "api-scheme", "e", "", "api scheme")
}

func initLogger() {
	logrus.SetReportCaller(true)
}

func patchConfig(host, secret string, rpcPort, apiPort int) {
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
}
