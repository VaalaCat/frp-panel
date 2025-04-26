package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/VaalaCat/frp-panel/conf"
	"github.com/VaalaCat/frp-panel/defs"
	"github.com/VaalaCat/frp-panel/pb"
	"github.com/VaalaCat/frp-panel/services/app"
	"github.com/VaalaCat/frp-panel/services/rpc"
	"github.com/VaalaCat/frp-panel/utils"
	"github.com/VaalaCat/frp-panel/utils/logger"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	clientCmd *cobra.Command
	rootCmd   *cobra.Command
)

func initCommand(appInstance app.Application) {
	rootCmd = &cobra.Command{
		Use:   "frp-panel",
		Short: "frp-panel is a frp panel QwQ",
	}
	CmdListWithFlag := initCmdWithFlag(appInstance)
	CmdListWithoutFlag := initCmdWithoutFlag(appInstance)
	rootCmd.AddCommand(CmdListWithFlag...)
	rootCmd.AddCommand(CmdListWithoutFlag...)
}

func initCmdWithFlag(appInstance app.Application) []*cobra.Command {
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
		rpcUrl       string
		apiUrl       string
	)

	clientCmd = &cobra.Command{
		Use:   "client [-s client secret] [-i client id] [-a app secret] [-t api host] [-r rpc host] [-c rpc port] [-p api port]",
		Short: "run managed frpc",
		Run: func(cmd *cobra.Command, args []string) {
			run := func() {
				patchConfig(appInstance,
					apiHost, rpcHost, appSecret,
					clientID, clientSecret,
					apiScheme, rpcPort, apiPort,
					apiUrl, rpcUrl)
				runClient(appInstance)
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
			pullRunConfig(appInstance,
				joinToken, appSecret, rpcHost, apiScheme, rpcPort, apiPort, clientID, apiHost, apiUrl, rpcUrl)
		},
	}

	clientCmd.Flags().StringVarP(&clientSecret, "secret", "s", "", "client secret")
	clientCmd.Flags().StringVarP(&clientID, "id", "i", "", "client id")

	clientCmd.Flags().StringVarP(&appSecret, "app", "a", "", "app secret")

	clientCmd.Flags().StringVar(&rpcUrl, "rpc-url", "", "rpc url, master rpc url, scheme can be grpc/ws/wss://hostname:port")

	clientCmd.Flags().StringVar(&apiUrl, "api-url", "", "api url, master api url, scheme can be http/https://hostname:port")

	// deprecated start
	clientCmd.Flags().StringVarP(&rpcHost, "rpc", "r", "", "deprecated, use --rpc-url instead, rpc host, canbe ip or domain")
	clientCmd.Flags().StringVarP(&apiHost, "api", "t", "", "deprecated, use --api-url instead, api host, canbe ip or domain")
	clientCmd.Flags().IntVarP(&rpcPort, "rpc-port", "c", 0, "deprecated, use --rpc-url instead, rpc port, master rpc port, scheme is grpc")
	clientCmd.Flags().IntVarP(&apiPort, "api-port", "p", 0, "deprecated, use --api-url instead, api port, master api port, scheme is http/https")
	clientCmd.Flags().StringVarP(&apiScheme, "api-scheme", "e", "", "deprecated, use --api-url instead, api scheme, master api scheme, scheme is http/https")
	joinCmd.Flags().IntVarP(&rpcPort, "rpc-port", "c", 0, "deprecated, use --rpc-url instead, rpc port, master rpc port, scheme is grpc")
	joinCmd.Flags().IntVarP(&apiPort, "api-port", "p", 0, "deprecated, use --api-url instead, api port, master api port, scheme is http/https")
	joinCmd.Flags().StringVarP(&rpcHost, "rpc", "r", "", "deprecated, use --rpc-url instead, rpc host, canbe ip or domain")
	joinCmd.Flags().StringVarP(&apiHost, "api", "t", "", "deprecated, use --api-url instead, api host, canbe ip or domain")
	joinCmd.Flags().StringVarP(&apiScheme, "api-scheme", "e", "", "deprecated, use --api-url instead, api scheme, master api scheme, scheme is http/https")
	// deprecated end

	joinCmd.Flags().StringVarP(&appSecret, "app", "a", "", "app secret")
	joinCmd.Flags().StringVarP(&joinToken, "join-token", "j", "", "your token from master")
	joinCmd.Flags().StringVarP(&clientID, "id", "i", "", "client id")
	joinCmd.Flags().StringVar(&rpcUrl, "rpc-url", "", "rpc url, master rpc url, scheme can be grpc/ws/wss://hostname:port")
	joinCmd.Flags().StringVar(&apiUrl, "api-url", "", "api url, master api url, scheme can be http/https://hostname:port")

	return []*cobra.Command{clientCmd, joinCmd}
}

func initCmdWithoutFlag(appInstance app.Application) []*cobra.Command {

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
		installServiceCmd,
		uninstallServiceCmd,
		startServiceCmd,
		stopServiceCmd,
		restartServiceCmd,
		versionCmd,
	}
}

func patchConfig(appInstance app.Application, apiHost, rpcHost, secret, clientID, clientSecret, apiScheme string, rpcPort, apiPort int, apiUrl, rpcUrl string) {
	c := context.Background()
	tmpCfg := appInstance.GetConfig()
	if len(rpcHost) != 0 {
		tmpCfg.Master.RPCHost = rpcHost
		tmpCfg.Master.APIHost = rpcHost
	}

	if len(apiHost) != 0 {
		tmpCfg.Master.APIHost = apiHost
	}

	if len(secret) != 0 {
		tmpCfg.App.Secret = secret
	}
	if rpcPort != 0 {
		tmpCfg.Master.RPCPort = rpcPort
	}
	if apiPort != 0 {
		tmpCfg.Master.APIPort = apiPort
	}
	if len(apiScheme) != 0 {
		tmpCfg.Master.APIScheme = apiScheme
	}
	if len(clientID) != 0 {
		tmpCfg.Client.ID = clientID
	}
	if len(clientSecret) != 0 {
		tmpCfg.Client.Secret = clientSecret
	}

	if len(apiUrl) != 0 {
		tmpCfg.Client.APIUrl = apiUrl
	}
	if len(rpcUrl) != 0 {
		tmpCfg.Client.RPCUrl = rpcUrl
	}

	if rpcPort != 0 || apiPort != 0 || len(apiScheme) != 0 || len(rpcHost) != 0 || len(apiHost) != 0 {
		logger.Logger(c).Warnf("deprecatedenv configs !!! rpc host: %s, rpc port: %d, api host: %s, api port: %d, api scheme: %s",
			tmpCfg.Master.RPCHost, tmpCfg.Master.RPCPort,
			tmpCfg.Master.APIHost, tmpCfg.Master.APIPort,
			tmpCfg.Master.APIScheme)
	}
	logger.Logger(c).Infof("env config, api url: %s, rpc url: %s", tmpCfg.Client.APIUrl, tmpCfg.Client.RPCUrl)
}

func setMasterCommandIfNonePresent() {
	cmd, _, err := rootCmd.Find(os.Args[1:])
	if err == nil && cmd.Use == rootCmd.Use && cmd.Flags().Parse(os.Args[1:]) != pflag.ErrHelp {
		args := append([]string{"master"}, os.Args[1:]...)
		rootCmd.SetArgs(args)
	}
}

func pullRunConfig(appInstance app.Application, joinToken, appSecret, rpcHost, apiScheme string, rpcPort, apiPort int, clientID, apiHost string, apiUrl, rpcUrl string) {
	c := context.Background()
	if err := checkPullParams(joinToken, apiHost, apiScheme, apiPort, apiUrl); err != nil {
		logger.Logger(c).Errorf("check pull params failed: %s", err.Error())
		return
	}

	if err := utils.EnsureDirectoryExists(defs.SysEnvPath); err != nil {
		logger.Logger(c).Errorf("ensure directory failed: %s", err.Error())
		return
	}

	if len(clientID) == 0 {
		clientID = utils.GetHostnameWithIP()
	}

	clientID = utils.MakeClientIDPermited(clientID)
	patchConfig(appInstance, apiHost, rpcHost, appSecret, "", "", apiScheme, rpcPort, apiPort, apiUrl, rpcUrl)

	initResp, err := rpc.InitClient(appInstance, clientID, joinToken)
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
	clientResp, err := rpc.GetClient(appInstance, clientID, joinToken)
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

	envMap, err := godotenv.Read(defs.SysEnvPath)
	if err != nil {
		envMap = make(map[string]string)
		logger.Logger(c).Warnf("read env file failed, try to create: %s", err.Error())
	}

	envMap[defs.EnvAppSecret] = appSecret
	envMap[defs.EnvClientID] = clientID
	envMap[defs.EnvClientSecret] = client.GetSecret()
	envMap[defs.EnvClientAPIUrl] = apiUrl
	envMap[defs.EnvClientRPCUrl] = rpcUrl

	if err = godotenv.Write(envMap, defs.SysEnvPath); err != nil {
		logger.Logger(c).Errorf("write env file failed: %s", err.Error())
		return
	}
	logger.Logger(c).Infof("config saved to env file: %s, you can use `frp-panel client` without args to run client,\n\nconfig is: [%v]", defs.SysEnvPath, envMap)
}

func checkPullParams(joinToken, apiHost, apiScheme string, apiPort int, apiUrl string) error {
	if len(joinToken) == 0 {
		return errors.New("join token is empty")
	}
	if len(apiUrl) == 0 {
		if len(apiHost) == 0 {
			return errors.New("api host is empty")
		}
		if len(apiScheme) == 0 {
			return errors.New("api scheme is empty")
		}
	}

	if apiPort == 0 {
		return errors.New("api port is empty")
	}
	return nil
}
