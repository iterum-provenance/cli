package data

import (
	"fmt"
	"log"

	"github.com/iterum-provenance/cli/idv"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(setupCmd)
}

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Sets up a tracked data repository via the iterum daemon",
	Long:  `Setting up a dataset will hand the current idv config to the daemon which will try to create the necesarry files and setup remotely`,
	Run:   setupRun,
}

func setupRun(cmd *cobra.Command, args []string) {
	err := idv.Setup()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Daemon and local IDV repo initialized")
		fmt.Println("Ready to version your data...")
	}
}
