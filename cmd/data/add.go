package data

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.PersistentFlags().BoolVarP(&Recursive, "recursive", "r", false, "Descend recursively into passed folders")
}

var addCmd = &cobra.Command{
	Use:   "add [file/folder]...",
	Short: "Add or Update files to the current commit",
	Long:  `Stages files to be added to the dataset, or in case of name clash to update those files`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errNotEnoughArgs
		}
		var invalids []string
		for _, arg := range args {
			if !isValidLocation(arg) {
				invalids = append(invalids, arg)
			}
		}
		if len(invalids) != 0 {
			return errInvalidArgs(invalids...)
		}
		return nil
	},
	Run: addRun,
}

func addRun(cmd *cobra.Command, args []string) {
	log.Println("`iterum data add`")
}
