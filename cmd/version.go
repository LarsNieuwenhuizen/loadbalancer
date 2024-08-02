package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of loadbalancer",
	Long:  `All software has versions. This is loadbalancer's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("##VERSION##")
	},
}
