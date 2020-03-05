package data

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(lsCmd)
}

var lsCmd = &cobra.Command{
	Use:   "ls [selector]",
	Short: "List files in the current commit (with selector)",
	Long:  `List all files in the current commit filtered using specified regex selector`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errNotEnoughArgs
		}
		if len(args) == 0 || isValidSelector(args[0]) {
			return nil
		}
		return errInvalidArgs(args[0])
	},
	Run: lsRun,
}

func lsRun(cmd *cobra.Command, args []string) {
	log.Println("`iterum data ls`")
}
