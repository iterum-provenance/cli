package data

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(setupCmd)
}

var setupCmd = &cobra.Command{
	Use:   "setup [config-file]",
	Short: "Sets up a tracked data repository via the iterum daemon",
	Long:  `Setting up a dataset will hand the current idv config to the daemon which will try to create the necesarry files and setup remotely`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errWrongAmountOfArgs(1)
		}
		if isValidLocation(args[0]) {
			return nil
		}
		return errInvalidLocation
	},
	Run: setupRun,
}

func setupRun(cmd *cobra.Command, args []string) {
	log.Println("`iterum data setup`")
}
