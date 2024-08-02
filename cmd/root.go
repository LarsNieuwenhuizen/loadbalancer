package cmd

import (
	"os"
	"path/filepath"
)

func Execute() {
	rootCmd.AddCommand(startCommand)
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
