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
   _       U  ___ u    _      ____    ____      _       _         _      _   _      ____ U _____ u   ____
  |"|       \/"_ \/U  /"\  u |  _"\U | __")uU  /"\  u  |"|    U  /"\  u | \ |"|  U /"___|\| ___"|/U |  _"\ u
U | | u     | | | | \/ _ \/ /| | | |\|  _ \/ \/ _ \/ U | | u   \/ _ \/ <|  \| |> \| | u   |  _|"   \| |_) |/
 \| |/__.-,_| |_| | / ___ \ U| |_| |\| |_) | / ___ \  \| |/__  / ___ \ U| |\  |u  | |/__  | |___    |  _ <
  |_____|\_)-\___/ /_/   \_\ |____/ u|____/ /_/   \_\  |_____|/_/   \_\ |_| \_|    \____| |_____|   |_| \_\
  //  \\      \\    \\    >>  |||_  _|| \\_  \\    >>  //  \\  \\    >> ||   \\,-._// \\  <<   >>   //   \\_
 (_")("_)    (__)  (__)  (__)(__)_)(__) (__)(__)  (__)(_")("_)(__)  (__)(_")  (_/(__)(__)(__) (__) (__)  (__)
		`)
		cmd.Help()
	},
}
