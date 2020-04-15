package data

import (
	"fmt"
	"log"

	"github.com/iterum-provenance/cli/idv"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.PersistentFlags().BoolVarP(&Recursive, "recursive", "r", false, "Descend recursively into passed folders")
	addCmd.PersistentFlags().StringSliceVarP(&Exclusions, "exclude", "x", []string{}, "Exclude files and folders -x selector1 -x selector2,selector3")
	addCmd.PersistentFlags().BoolVarP(&ShowExcluded, "show-excluded", "s", false, "Show list of excluded files")
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
			return errInvalidArgs("Not a file or directory", invalids...)
		}
		return nil
	},
	Run: addRun,
}

func addRun(cmd *cobra.Command, args []string) {
	allFiles := getAllFiles(args)
	whitelisted := exclude(allFiles)
	adds, updates, err := idv.AddFiles(whitelisted)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("PASSED %v file(s)\nADDED %v file(s)\nUPDATED %v file(s)\n", len(whitelisted), adds, updates)
	if len(whitelisted) > adds+updates && len(whitelisted) > 0 {
		fmt.Println("To see which files were uploaded exactly use `iterum data status --local-path`")
	}
}
