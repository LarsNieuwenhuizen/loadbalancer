package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/LarsNieuwenhuizen/loadbalancer"
	"github.com/spf13/cobra"
)

var startCommand = &cobra.Command{
	Use:     "start",
	Short:   "Start the loadbalancer",
	Long:    `This command starts the loadbalancer with the given configuration.`,
	Example: `loadbalancer start --config /path/to/config.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		configFlag, _ := cmd.Flags().GetString("config")
		_, err := os.Stat(configFlag)
		if errors.Is(err, fs.ErrNotExist) {
			fmt.Println("Config file " + configFlag + " does not exist, check the path or change it with the --config flag")
			os.Exit(1)
		}

		lb := loadbalancer.LoadBalancer{}
		lb.ConfigureFromYaml(configFlag)
		fmt.Println("Loadbalancer started, configured with", len(lb.Configuration.BackendServers), "backend servers")
		lb.Start()
	},
}
