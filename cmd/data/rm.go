package data

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rmCmd)
	rmCmd.PersistentFlags().BoolVarP(&Recursive, "recursive", "r", false, "Descend recursively into passed folders")
}

var rmCmd = &cobra.Command{
	Use:   "rm [file/folder]...",
	Short: "Stages the removal of specified files from the dataset",
	Long:  `Locally stages the removal of files from the dataset. Will be reflected in next commit. Folders will pick all files in said folder.`,
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
	Run: rmRun,
}

func rmRun(cmd *cobra.Command, args []string) {
	log.Println("`iterum data rm`")
}
