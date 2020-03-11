package data

import (
	"fmt"
	"log"

	"github.com/Mantsje/iterum-cli/idv"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rmCmd)
	rmCmd.PersistentFlags().BoolVarP(&Recursive, "recursive", "r", false, "Descend recursively into passed folders")
	rmCmd.PersistentFlags().StringSliceVarP(&Exclusions, "exclude", "x", []string{}, "Exclude files and folders from removal using -x selector1 -x selector2,selector3")
	rmCmd.PersistentFlags().BoolVarP(&ShowExcluded, "show-excluded", "s", false, "Show list of excluded files which are NOT removed")
}

var rmCmd = &cobra.Command{
	Use:   "rm [idvname/file/folder]...",
	Short: "Remove files or (un)stage them for the current commit",
	Long:  `Stages files to be removed from the dataset. If existing paths are given, these files will be converted to internal formatting and removed from commit, if non-existent paths are passed, iterum tries to find files that match with the given value. The exclusions are used to exclude files from removal. So if you want to remove all files from a folder recursively from the commit except for a few`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errNotEnoughArgs
		}
		return nil
	},
	Run: rmRun,
}

func getPaths(args []string) (paths, names []string) {
	for _, arg := range args {
		if isValidLocation(arg) {
			paths = append(paths, arg)
		} else {
			names = append(names, arg)
		}
	}
	return
}

func rmRun(cmd *cobra.Command, args []string) {
	paths, names := getPaths(args)
	allFiles := getAllFiles(paths)
	whitelisted := exclude(allFiles)
	removals, err := idv.RemoveFiles(whitelisted, names)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("GOT %v potential files\nREMOVED %v file(s)\n", len(names)+len(whitelisted), removals)
}
