package main

import (
	"context"
	"fmt"
	"os"

	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/logger"
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
			run := func() {
				patchConfig(rpcHost, appSecret,
					clientID, clientSecret,
					apiScheme, rpcPort, apiPort)
				runClient()
			}
			if srv, err := utils.CreateSystemService(args, run); err != nil {
				run()
			} else {
				srv.Run()
			}
		},
	}
	serverCmd = &cobra.Command{
		Use:   "server [-s client secret] [-i client id] [-a app secret] [-r rpc host] [-c rpc port] [-p api port]",
		Short: "run managed frps",
		Run: func(cmd *cobra.Command, args []string) {
			run := func() {
				patchConfig(rpcHost, appSecret,
					clientID, clientSecret,
					apiScheme, rpcPort, apiPort)
				runServer()
			}
			if srv, err := utils.CreateSystemService(args, run); err != nil {
				run()
			} else {
				srv.Run()
			}
		},
	}
	masterCmd = &cobra.Command{
		Use:   "master",
		Short: "run frp-panel manager",
		Run: func(cmd *cobra.Command, args []string) {
			if srv, err := utils.CreateSystemService(args, runMaster); err != nil {
				runMaster()
			} else {
				srv.Run()
			}
		},
	}
	rootCmd = &cobra.Command{
		Use:   "frp-panel",
		Short: "frp-panel is a frp panel QwQ",
	}

	installServiceCmd := &cobra.Command{
		Use:                   "install",
		Short:                 "install frp-panel as service",
		DisableFlagParsing:    true,
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			utils.ControlSystemService(args, "install", func() {})
		},
	}

	uninstallServiceCmd := &cobra.Command{
		Use:                   "uninstall",
		Short:                 "uninstall frp-panel service",
		DisableFlagParsing:    true,
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			utils.ControlSystemService(args, "uninstall", func() {})
		},
	}

	startServiceCmd := &cobra.Command{
		Use:                   "start",
		Short:                 "start frp-panel service",
		DisableFlagParsing:    true,
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			utils.ControlSystemService(args, "start", func() {})
		},
	}

	stopServiceCmd := &cobra.Command{
		Use:                   "stop",
		Short:                 "stop frp-panel service",
		DisableFlagParsing:    true,
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			utils.ControlSystemService(args, "stop", func() {})
		},
	}

	restartServiceCmd := &cobra.Command{
		Use:                   "restart",
		Short:                 "restart frp-panel service",
		DisableFlagParsing:    true,
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			utils.ControlSystemService(args, "restart", func() {})
		},
	}

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version info of frp-panel",
		Long:  `All software has versions. This is frp-panel's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(conf.GetVersion().String())
		},
	}

	rootCmd.AddCommand(clientCmd, serverCmd, masterCmd, versionCmd,
		installServiceCmd, uninstallServiceCmd,
		startServiceCmd, stopServiceCmd, restartServiceCmd)

	clientCmd.Flags().StringVarP(&clientSecret, "secret", "s", "", "client secret")
	serverCmd.Flags().StringVarP(&clientSecret, "secret", "s", "", "client secret")
	clientCmd.Flags().StringVarP(&clientID, "id", "i", "", "client id")
	serverCmd.Flags().StringVarP(&clientID, "id", "i", "", "client id")
	clientCmd.Flags().StringVarP(&rpcHost, "rpc", "r", "", "rpc host, canbe ip or domain")
	serverCmd.Flags().StringVarP(&rpcHost, "rpc", "r", "", "rpc host, canbe ip or domain")
	clientCmd.Flags().StringVarP(&appSecret, "app", "a", "", "app secret")
	serverCmd.Flags().StringVarP(&appSecret, "app", "a", "", "app secret")

	clientCmd.Flags().IntVarP(&rpcPort, "rpc-port", "c", 0, "rpc port, master rpc port, scheme is grpc")
	serverCmd.Flags().IntVarP(&rpcPort, "rpc-port", "c", 0, "rpc port, master rpc port, scheme is grpc")
	clientCmd.Flags().IntVarP(&apiPort, "api-port", "p", 0, "api port, master api port, scheme is http/https")
	serverCmd.Flags().IntVarP(&apiPort, "api-port", "p", 0, "api port, master api port, scheme is http/https")

	clientCmd.Flags().StringVarP(&apiScheme, "api-scheme", "e", "", "api scheme, master api scheme, scheme is http/https")
	serverCmd.Flags().StringVarP(&apiScheme, "api-scheme", "e", "", "api scheme, master api scheme, scheme is http/https")
}

func initLogger() {
	logger.Instance().SetReportCaller(true)
}

func patchConfig(host, secret, clientID, clientSecret, apiScheme string, rpcPort, apiPort int) {
	c := context.Background()
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
	logger.Logger(c).Infof("env config rpc host: %s, rpc port: %d, api host: %s, api port: %d, api scheme: %s",
		conf.Get().Master.RPCHost, conf.Get().Master.RPCPort,
		conf.Get().Master.APIHost, conf.Get().Master.APIPort,
		conf.Get().Master.APIScheme)
}

func setMasterCommandIfNonePresent() {
	cmd, _, err := rootCmd.Find(os.Args[1:])
	if err == nil && cmd.Use == rootCmd.Use && cmd.Flags().Parse(os.Args[1:]) != pflag.ErrHelp {
		args := append([]string{"master"}, os.Args[1:]...)
		rootCmd.SetArgs(args)
	}
}
