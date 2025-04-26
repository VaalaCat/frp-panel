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
	"go.uber.org/fx"
)

type CommonArgs struct {
	ClientSecret *string
	ClientID     *string
	AppSecret    *string
	RpcUrl       *string
	ApiUrl       *string

	RpcHost   *string
	ApiHost   *string
	RpcPort   *int
	ApiPort   *int
	ApiScheme *string
}

func buildCommand() *cobra.Command {
	cfg := conf.NewConfig()

	return NewRootCmd(
		NewMasterCmd(cfg),
		NewClientCmd(cfg),
		NewServerCmd(cfg),
		NewJoinCmd(),
		NewInstallServiceCmd(),
		NewUninstallServiceCmd(),
		NewStartServiceCmd(),
		NewStopServiceCmd(),
		NewRestartServiceCmd(),
		NewVersionCmd(),
	)
}

func AddCommonFlags(commonCmd *cobra.Command) {
	commonCmd.Flags().StringP("secret", "s", "", "client secret")
	commonCmd.Flags().StringP("id", "i", "", "client id")
	commonCmd.Flags().StringP("app", "a", "", "app secret")
	commonCmd.Flags().String("rpc-url", "", "rpc url, master rpc url, scheme can be grpc/ws/wss://hostname:port")
	commonCmd.Flags().String("api-url", "", "api url, master api url, scheme can be http/https://hostname:port")

	// deprecated start
	commonCmd.Flags().StringP("rpc-host", "r", "", "deprecated, use --rpc-url instead, rpc host, canbe ip or domain")
	commonCmd.Flags().StringP("api-host", "t", "", "deprecated, use --api-url instead, api host, canbe ip or domain")
	commonCmd.Flags().IntP("rpc-port", "c", 0, "deprecated, use --rpc-url instead, rpc port, master rpc port, scheme is grpc")
	commonCmd.Flags().IntP("api-port", "p", 0, "deprecated, use --api-url instead, api port, master api port, scheme is http/https")
	commonCmd.Flags().StringP("api-scheme", "e", "", "deprecated, use --api-url instead, api scheme, master api scheme, scheme is http/https")
	// deprecated end
}

func GetCommonArgs(cmd *cobra.Command) CommonArgs {
	var commonArgs CommonArgs

	if clientSecret, err := cmd.Flags().GetString("secret"); err == nil {
		commonArgs.ClientSecret = &clientSecret
	}

	if clientID, err := cmd.Flags().GetString("id"); err == nil {
		commonArgs.ClientID = &clientID
	}

	if appSecret, err := cmd.Flags().GetString("app"); err == nil {
		commonArgs.AppSecret = &appSecret
	}

	if rpcURL, err := cmd.Flags().GetString("rpc-url"); err == nil {
		commonArgs.RpcUrl = &rpcURL
	}

	if apiURL, err := cmd.Flags().GetString("api-url"); err == nil {
		commonArgs.ApiUrl = &apiURL
	}

	if rpcHost, err := cmd.Flags().GetString("rpc-host"); err == nil {
		commonArgs.RpcHost = &rpcHost
	}

	if apiHost, err := cmd.Flags().GetString("api-host"); err == nil {
		commonArgs.ApiHost = &apiHost
	}

	if rpcPort, err := cmd.Flags().GetInt("rpc-port"); err == nil {
		commonArgs.RpcPort = &rpcPort
	}

	if apiPort, err := cmd.Flags().GetInt("api-port"); err == nil {
		commonArgs.ApiPort = &apiPort
	}

	if apiScheme, err := cmd.Flags().GetString("api-scheme"); err == nil {
		commonArgs.ApiScheme = &apiScheme
	}

	return commonArgs
}

type JoinArgs struct {
	CommonArgs
	JoinToken *string
}

func NewJoinCmd() *cobra.Command {
	joinCmd := &cobra.Command{
		Use:   "join [-j join token] [-r rpc host] [-p api port] [-e api scheme]",
		Short: "join to master with token, save param to config",
		Run: func(cmd *cobra.Command, args []string) {
			commonArgs := GetCommonArgs(cmd)
			joinArgs := &JoinArgs{
				CommonArgs: commonArgs,
			}
			if joinToken, err := cmd.Flags().GetString("join-token"); err == nil {
				joinArgs.JoinToken = &joinToken
			}

			appInstance := app.NewApp()
			pullRunConfig(appInstance, joinArgs)
		},
	}

	joinCmd.Flags().StringP("join-token", "j", "", "your token from master")
	AddCommonFlags(joinCmd)

	return joinCmd
}

func NewMasterCmd(cfg conf.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "master",
		Short: "run frp-panel manager",
		Run: func(cmd *cobra.Command, args []string) {

			opts := []fx.Option{
				commonMod,
				masterMod,
				serverMod,
				fx.Supply(
					CommonArgs{},
					fx.Annotate(cfg, fx.ResultTags(`name:"originConfig"`)),
				),
				fx.Provide(fx.Annotate(NewDefaultServerConfig, fx.ResultTags(`name:"defaultServerConfig"`))),
				fx.Invoke(runMaster),
				fx.Invoke(runServer),
			}

			if !cfg.IsDebug {
				opts = append(opts, fx.NopLogger)
			}

			run := func() {
				masterApp := fx.New(opts...)
				masterApp.Run()
				if err := masterApp.Err(); err != nil {
					logger.Logger(context.Background()).Fatalf("masterApp FX Application Error: %v", err)
				}
			}

			if srv, err := utils.CreateSystemService(args, run); err != nil {
				run()
			} else {
				srv.Run()
			}
		},
	}
}

func NewClientCmd(cfg conf.Config) *cobra.Command {
	clientCmd := &cobra.Command{
		Use:   "client [-s client secret] [-i client id] [-a app secret] [-t api host] [-r rpc host] [-c rpc port] [-p api port]",
		Short: "run managed frpc",
		Run: func(cmd *cobra.Command, args []string) {
			commonArgs := GetCommonArgs(cmd)

			opts := []fx.Option{
				clientMod,
				commonMod,
				fx.Supply(
					commonArgs,
					fx.Annotate(cfg, fx.ResultTags(`name:"originConfig"`)),
				),
				fx.Invoke(runClient),
			}

			if !cfg.IsDebug {
				opts = append(opts, fx.NopLogger)
			}

			run := func() {
				clientApp := fx.New(opts...)
				clientApp.Run()
				if err := clientApp.Err(); err != nil {
					logger.Logger(context.Background()).Fatalf("clientApp FX Application Error: %v", err)
				}
			}
			if srv, err := utils.CreateSystemService(args, run); err != nil {
				run()
			} else {
				srv.Run()
			}
		},
	}

	AddCommonFlags(clientCmd)

	return clientCmd
}

func NewServerCmd(cfg conf.Config) *cobra.Command {
	serverCmd := &cobra.Command{
		Use:   "server [-s client secret] [-i client id] [-a app secret] [-r rpc host] [-c rpc port] [-p api port]",
		Short: "run managed frps",
		Run: func(cmd *cobra.Command, args []string) {
			commonArgs := GetCommonArgs(cmd)
			opts := []fx.Option{
				serverMod,
				commonMod,
				fx.Supply(
					commonArgs,
					fx.Annotate(cfg, fx.ResultTags(`name:"originConfig"`)),
				),
				fx.Invoke(runServer),
			}

			if !cfg.IsDebug {
				opts = append(opts, fx.NopLogger)
			}

			run := func() {
				serverApp := fx.New(opts...)
				serverApp.Run()
				if err := serverApp.Err(); err != nil {
					logger.Logger(context.Background()).Fatalf("serverApp FX Application Error: %v", err)
				}
			}
			if srv, err := utils.CreateSystemService(args, run); err != nil {
				run()
			} else {
				srv.Run()
			}
		},
	}

	AddCommonFlags(serverCmd)

	return serverCmd
}

func NewInstallServiceCmd() *cobra.Command {
	return &cobra.Command{
		Use:                   "install",
		Short:                 "install frp-panel as service",
		DisableFlagParsing:    true,
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			utils.ControlSystemService(args, "install", func() {})
		},
	}
}

func NewUninstallServiceCmd() *cobra.Command {
	return &cobra.Command{
		Use:                   "uninstall",
		Short:                 "uninstall frp-panel service",
		DisableFlagParsing:    true,
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			utils.ControlSystemService(args, "uninstall", func() {})
		},
	}
}

func NewStartServiceCmd() *cobra.Command {
	return &cobra.Command{
		Use:                   "start",
		Short:                 "start frp-panel service",
		DisableFlagParsing:    true,
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			utils.ControlSystemService(args, "start", func() {})
		},
	}
}

func NewStopServiceCmd() *cobra.Command {
	return &cobra.Command{
		Use:                   "stop",
		Short:                 "stop frp-panel service",
		DisableFlagParsing:    true,
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			utils.ControlSystemService(args, "stop", func() {})
		},
	}
}

func NewRestartServiceCmd() *cobra.Command {
	return &cobra.Command{
		Use:                   "restart",
		Short:                 "restart frp-panel service",
		DisableFlagParsing:    true,
		DisableFlagsInUseLine: true,
		Run: func(cmd *cobra.Command, args []string) {
			utils.ControlSystemService(args, "restart", func() {})
		},
	}
}

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version info of frp-panel",
		Long:  `All software has versions. This is frp-panel's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(conf.GetVersion().String())
		},
	}
}

func patchConfig(appInstance app.Application, commonArgs CommonArgs) conf.Config {
	c := context.Background()
	tmpCfg := appInstance.GetConfig()

	if commonArgs.RpcHost != nil {
		tmpCfg.Master.RPCHost = *commonArgs.RpcHost
		tmpCfg.Master.APIHost = *commonArgs.RpcHost
	}

	if commonArgs.ApiHost != nil {
		tmpCfg.Master.APIHost = *commonArgs.ApiHost
	}

	if commonArgs.AppSecret != nil {
		tmpCfg.App.Secret = *commonArgs.AppSecret
	}
	if commonArgs.RpcPort != nil {
		tmpCfg.Master.RPCPort = *commonArgs.RpcPort
	}
	if commonArgs.ApiPort != nil {
		tmpCfg.Master.APIPort = *commonArgs.ApiPort
	}
	if commonArgs.ApiScheme != nil {
		tmpCfg.Master.APIScheme = *commonArgs.ApiScheme
	}
	if commonArgs.ClientID != nil {
		tmpCfg.Client.ID = *commonArgs.ClientID
	}
	if commonArgs.ClientSecret != nil {
		tmpCfg.Client.Secret = *commonArgs.ClientSecret
	}

	if commonArgs.ApiUrl != nil {
		tmpCfg.Client.APIUrl = *commonArgs.ApiUrl
	}
	if commonArgs.RpcUrl != nil {
		tmpCfg.Client.RPCUrl = *commonArgs.RpcUrl
	}

	if commonArgs.RpcPort != nil || commonArgs.ApiPort != nil ||
		commonArgs.ApiScheme != nil ||
		commonArgs.RpcHost != nil || commonArgs.ApiHost != nil {
		logger.Logger(c).Warnf("deprecatedenv configs !!! rpc host: %s, rpc port: %d, api host: %s, api port: %d, api scheme: %s",
			tmpCfg.Master.RPCHost, tmpCfg.Master.RPCPort,
			tmpCfg.Master.APIHost, tmpCfg.Master.APIPort,
			tmpCfg.Master.APIScheme)
	}
	logger.Logger(c).Infof("env config, api url: %s, rpc url: %s", tmpCfg.Client.APIUrl, tmpCfg.Client.RPCUrl)
	return tmpCfg
}

func setMasterCommandIfNonePresent(rootCmd *cobra.Command) {
	cmd, _, err := rootCmd.Find(os.Args[1:])
	if err == nil && cmd.Use == rootCmd.Use && cmd.Flags().Parse(os.Args[1:]) != pflag.ErrHelp {
		args := append([]string{"master"}, os.Args[1:]...)
		rootCmd.SetArgs(args)
	}
}

func pullRunConfig(appInstance app.Application, joinArgs *JoinArgs) {
	c := context.Background()
	if err := checkPullParams(joinArgs); err != nil {
		logger.Logger(c).Errorf("check pull params failed: %s", err.Error())
		return
	}

	if err := utils.EnsureDirectoryExists(defs.SysEnvPath); err != nil {
		logger.Logger(c).Errorf("ensure directory failed: %s", err.Error())
		return
	}

	var clientID string

	if cliID := joinArgs.ClientID; cliID == nil || len(*cliID) == 0 {
		clientID = utils.GetHostnameWithIP()
	}

	clientID = utils.MakeClientIDPermited(clientID)
	patchConfig(appInstance, joinArgs.CommonArgs)

	initResp, err := rpc.InitClient(appInstance, clientID, *joinArgs.JoinToken)
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
	clientResp, err := rpc.GetClient(appInstance, clientID, *joinArgs.JoinToken)
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

	envMap[defs.EnvAppSecret] = *joinArgs.AppSecret
	envMap[defs.EnvClientID] = clientID
	envMap[defs.EnvClientSecret] = client.GetSecret()
	envMap[defs.EnvClientAPIUrl] = *joinArgs.ApiUrl
	envMap[defs.EnvClientRPCUrl] = *joinArgs.RpcUrl

	if err = godotenv.Write(envMap, defs.SysEnvPath); err != nil {
		logger.Logger(c).Errorf("write env file failed: %s", err.Error())
		return
	}
	logger.Logger(c).Infof("config saved to env file: %s, you can use `frp-panel client` without args to run client,\n\nconfig is: [%v]", defs.SysEnvPath, envMap)
}

func checkPullParams(joinArgs *JoinArgs) error {
	if joinToken := joinArgs.JoinToken; joinToken != nil && len(*joinToken) == 0 {
		return errors.New("join token is empty")
	}

	if apiUrl := joinArgs.ApiUrl; apiUrl == nil || len(*apiUrl) == 0 {
		if apiHost := joinArgs.ApiHost; apiHost == nil || len(*apiHost) == 0 {
			return errors.New("api host is empty")
		}
		if apiScheme := joinArgs.ApiScheme; apiScheme == nil || len(*apiScheme) == 0 {
			return errors.New("api scheme is empty")
		}
	}

	if apiPort := joinArgs.ApiPort; apiPort == nil || *apiPort == 0 {
		return errors.New("api port is empty")
	}

	return nil
}

func NewRootCmd(cmds ...*cobra.Command) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "frp-panel",
		Short: "frp-panel is a frp panel QwQ",
	}

	rootCmd.AddCommand(cmds...)

	return rootCmd
}
