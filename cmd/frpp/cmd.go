package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/VaalaCat/frp-panel/common"
	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/logger"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/rpc"
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/joho/godotenv"
	"github.com/spf13/cast"
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
	rootCmd = &cobra.Command{
		Use:   "frp-panel",
		Short: "frp-panel is a frp panel QwQ",
	}
	CmdListWithFlag := initCmdWithFlag()
	CmdListWithoutFlag := initCmdWithoutFlag()
	rootCmd.AddCommand(CmdListWithFlag...)
	rootCmd.AddCommand(CmdListWithoutFlag...)
}

func initCmdWithFlag() []*cobra.Command {
	var (
		clientSecret string
		clientID     string
		rpcHost      string
		apiHost      string
		appSecret    string
		rpcPort      int
		apiPort      int
		apiScheme    string
		joinToken    string
	)

	clientCmd = &cobra.Command{
		Use:   "client [-s client secret] [-i client id] [-a app secret] [-t api host] [-r rpc host] [-c rpc port] [-p api port]",
		Short: "run managed frpc",
		Run: func(cmd *cobra.Command, args []string) {
			run := func() {
				patchConfig(apiHost, rpcHost, appSecret,
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
				patchConfig(apiHost, rpcHost, appSecret,
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

	joinCmd := &cobra.Command{
		Use:   "join [-j join token] [-r rpc host] [-p api port] [-e api scheme]",
		Short: "join to master with token, save param to config",
		Run: func(cmd *cobra.Command, args []string) {
			pullRunConfig(joinToken, appSecret, rpcHost, apiScheme, rpcPort, apiPort, clientID, apiHost)
		},
	}

	clientCmd.Flags().StringVarP(&clientSecret, "secret", "s", "", "client secret")
	serverCmd.Flags().StringVarP(&clientSecret, "secret", "s", "", "client secret")
	clientCmd.Flags().StringVarP(&clientID, "id", "i", "", "client id")
	serverCmd.Flags().StringVarP(&clientID, "id", "i", "", "client id")
	clientCmd.Flags().StringVarP(&rpcHost, "rpc", "r", "", "rpc host, canbe ip or domain")
	serverCmd.Flags().StringVarP(&rpcHost, "rpc", "r", "", "rpc host, canbe ip or domain")

	clientCmd.Flags().StringVarP(&apiHost, "api", "t", "", "api host, canbe ip or domain")
	serverCmd.Flags().StringVarP(&apiHost, "api", "t", "", "api host, canbe ip or domain")

	clientCmd.Flags().StringVarP(&appSecret, "app", "a", "", "app secret")
	serverCmd.Flags().StringVarP(&appSecret, "app", "a", "", "app secret")

	clientCmd.Flags().IntVarP(&rpcPort, "rpc-port", "c", 0, "rpc port, master rpc port, scheme is grpc")
	serverCmd.Flags().IntVarP(&rpcPort, "rpc-port", "c", 0, "rpc port, master rpc port, scheme is grpc")
	clientCmd.Flags().IntVarP(&apiPort, "api-port", "p", 0, "api port, master api port, scheme is http/https")
	serverCmd.Flags().IntVarP(&apiPort, "api-port", "p", 0, "api port, master api port, scheme is http/https")

	clientCmd.Flags().StringVarP(&apiScheme, "api-scheme", "e", "", "api scheme, master api scheme, scheme is http/https")
	serverCmd.Flags().StringVarP(&apiScheme, "api-scheme", "e", "", "api scheme, master api scheme, scheme is http/https")

	joinCmd.Flags().IntVarP(&rpcPort, "rpc-port", "c", 0, "rpc port, master rpc port, scheme is grpc")
	joinCmd.Flags().IntVarP(&apiPort, "api-port", "p", 0, "api port, master api port, scheme is http/https")
	joinCmd.Flags().StringVarP(&appSecret, "app", "a", "", "app secret")
	joinCmd.Flags().StringVarP(&joinToken, "join-token", "j", "", "your token from master")
	joinCmd.Flags().StringVarP(&rpcHost, "rpc", "r", "", "rpc host, canbe ip or domain")
	joinCmd.Flags().StringVarP(&apiHost, "api", "t", "", "api host, canbe ip or domain")
	joinCmd.Flags().StringVarP(&apiScheme, "api-scheme", "e", "", "api scheme, master api scheme, scheme is http/https")
	joinCmd.Flags().StringVarP(&clientID, "id", "i", "", "client id")

	return []*cobra.Command{clientCmd, serverCmd, joinCmd}
}

func initCmdWithoutFlag() []*cobra.Command {
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
	return []*cobra.Command{
		masterCmd,
		installServiceCmd,
		uninstallServiceCmd,
		startServiceCmd,
		stopServiceCmd,
		restartServiceCmd,
		versionCmd,
	}
}

func initLogger() {
	logger.Instance().SetReportCaller(true)
}

func patchConfig(apiHost, rpcHost, secret, clientID, clientSecret, apiScheme string, rpcPort, apiPort int) {
	c := context.Background()
	if len(rpcHost) != 0 {
		conf.Get().Master.RPCHost = rpcHost
		conf.Get().Master.APIHost = rpcHost
	}

	if len(apiHost) != 0 {
		conf.Get().Master.APIHost = apiHost
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

func pullRunConfig(joinToken, appSecret, rpcHost, apiScheme string, rpcPort, apiPort int, clientID, apiHost string) {
	c := context.Background()
	if err := checkPullParams(joinToken, apiHost, apiScheme, apiPort); err != nil {
		logger.Logger(c).Errorf("check pull params failed: %s", err.Error())
		return
	}

	if err := utils.EnsureDirectoryExists(common.SysEnvPath); err != nil {
		logger.Logger(c).Errorf("ensure directory failed: %s", err.Error())
		return
	}

	if len(clientID) == 0 {
		clientID = utils.GetHostnameWithIP()
	}

	clientID = utils.MakeClientIDPermited(clientID)
	patchConfig(apiHost, rpcHost, appSecret, "", "", apiScheme, rpcPort, apiPort)

	initResp, err := rpc.InitClient(clientID, joinToken)
	if err != nil {
		logger.Logger(c).Errorf("init client failed: %s", err.Error())
		return
	}
	if initResp == nil {
		logger.Logger(c).Errorf("init resp is nil")
		return
	}
	if initResp.GetStatus().GetCode() != pb.RespCode_RESP_CODE_SUCCESS {
		logger.Logger(c).Errorf("init client failed with status: %s", initResp.GetStatus().GetMessage())
		return
	}

	clientID = initResp.GetClientId()
	clientResp, err := rpc.GetClient(clientID, joinToken)
	if err != nil {
		logger.Logger(c).Errorf("get client failed: %s", err.Error())
		return
	}
	if clientResp == nil {
		logger.Logger(c).Errorf("client resp is nil")
		return
	}
	if clientResp.GetStatus().GetCode() != pb.RespCode_RESP_CODE_SUCCESS {
		logger.Logger(c).Errorf("client resp code is not success: %s", clientResp.GetStatus().GetMessage())
		return
	}

	client := clientResp.GetClient()
	if client == nil {
		logger.Logger(c).Errorf("client is nil")
		return
	}

	envMap, err := godotenv.Read(common.SysEnvPath)
	if err != nil {
		envMap = make(map[string]string)
		logger.Logger(c).Warnf("read env file failed, try to create: %s", err.Error())
	}

	envMap[common.EnvAppSecret] = appSecret
	envMap[common.EnvClientID] = clientID
	envMap[common.EnvClientSecret] = client.GetSecret()
	envMap[common.EnvMasterRPCHost] = rpcHost
	envMap[common.EnvMasterRPCPort] = cast.ToString(rpcPort)
	envMap[common.EnvMasterAPIHost] = apiHost
	envMap[common.EnvMasterAPIPort] = cast.ToString(apiPort)
	envMap[common.EnvMasterAPIScheme] = apiScheme

	if err = godotenv.Write(envMap, common.SysEnvPath); err != nil {
		logger.Logger(c).Errorf("write env file failed: %s", err.Error())
		return
	}
	logger.Logger(c).Infof("config saved to env file: %s, you can use `frp-panel client` without args to run client,\n\nconfig is: [%v]", common.SysEnvPath, envMap)
}

func checkPullParams(joinToken, apiHost, apiScheme string, apiPort int) error {
	if len(joinToken) == 0 {
		return errors.New("join token is empty")
	}
	if len(apiHost) == 0 {
		return errors.New("api host is empty")
	}
	if len(apiScheme) == 0 {
		return errors.New("api scheme is empty")
	}
	if apiPort == 0 {
		return errors.New("api port is empty")
	}
	return nil
}
