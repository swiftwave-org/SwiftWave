package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/swiftwave-org/swiftwave/system_config"
	"os"
)

var systemConfig *system_config.Config

var configFilePath = "/etc/swiftwave/config.yml"

var rootCmd = &cobra.Command{
	Use:   "swiftwave",
	Short: "SwiftWave is a self-hosted lightweight PaaS solution",
	Long:  `SwiftWave is a self-hosted lightweight PaaS solution to deploy and manage your applications on any VPS without any hassle of managing servers.`,
	Run: func(cmd *cobra.Command, args []string) {
		// print help
		err := cmd.Help()
		if err != nil {
			return
		}
	},
}

func Execute(config *system_config.Config) {
	systemConfig = config
	// set config and manager
	cobra.EnableCommandSorting = false
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
