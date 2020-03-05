package data

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(inspectCmd)
}

var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Probe the iterum daemon for which data config is currently active",
	Long:  `Probes the iterum daemon to return the current dataset configurarion and prints this`,
	Run:   inspectRun,
}

func inspectRun(cmd *cobra.Command, args []string) {
	log.Println("`iterum data inspect`")
}
