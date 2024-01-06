package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	clientSecret string
	clientID     string
	clientCmd    *cobra.Command
	serverCmd    *cobra.Command
	masterCmd    *cobra.Command
	rootCmd      *cobra.Command
)

func initCommand() {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	clientCmd = &cobra.Command{
		Use:   "client [-s client secret] [-i client id]",
		Short: "run managed frpc",
		Run: func(cmd *cobra.Command, args []string) {
			runClient(clientID, clientSecret)
		},
	}
	serverCmd = &cobra.Command{
		Use:   "server [-s client secret] [-i client id]",
		Short: "run managed frps",
		Run: func(cmd *cobra.Command, args []string) {
			runServer(clientID, clientSecret)
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
	clientCmd.Flags().StringVarP(&clientID, "id", "i", hostname, "client id")
	serverCmd.Flags().StringVarP(&clientID, "id", "i", hostname, "client id")
}

func initLogger() {
	logrus.SetReportCaller(true)
}
