package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/LarsNieuwenhuizen/loadbalancer"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "loadbalancer",
	Short: "A loadblancer written in go",
	Long:  `This loadbalancer passes requests to configured backend servers.`,
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	version := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of loadbalancer",
		Long:  `All software has versions. This is loadbalancer's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("##VERSION##")
		},
	}
	rootCmd.AddCommand(version)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	homeDir, _ := os.UserHomeDir()
	cfgFile := filepath.Join(homeDir, ".loadbalancer.yaml")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", cfgFile, "config yaml file")
}
