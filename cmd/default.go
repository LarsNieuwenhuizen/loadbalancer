package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "help",
	Short: "A loadblancer written in go",
	Long:  `This loadbalancer passes requests to configured backend servers.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`
__    __   __  ____ ____  __  __    __  __ _  ___ ____ ____
(  )  /  \ / _\(    (  _ \/ _\(  )  / _\(  ( \/ __(  __(  _ \
/ (_/(  O /    \) D () _ /    / (_//    /    ( (__ ) _) )   /
\____/\__/\_/\_(____(____\_/\_\____\_/\_\_)__)\___(____(__\_)
				`)
		cmd.Help()
	},
}
