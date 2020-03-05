package data

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:   "apply [config-file]",
	Short: "Apply a data config file at the iterum daemon",
	Long:  `Updates the data config at the iterum daemon. Used for setting/switching current datasets`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errWrongAmountOfArgs(1)
		}
		if isValidLocation(args[0]) {
			return nil
		}
		return errInvalidLocation
	},
	Run: applyRun,
}

func applyRun(cmd *cobra.Command, args []string) {
	log.Println("`iterum data apply`")
}
