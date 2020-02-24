package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(syncCmd)
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Pulls the main repo and all its submodules",
	Long:  `Updates each of the proejct's units, flows, etc`,
	Run:   syncRun,
}

func syncRun(cmd *cobra.Command, args []string) {
	fmt.Println("'Iterum sync' command")
}
