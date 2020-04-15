package data

import (
	"fmt"
	"log"

	"github.com/iterum-provenance/cli/idv"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rmCmd)
	rmCmd.PersistentFlags().BoolVarP(&Recursive, "recursive", "r", false, "Descend recursively into passed folders")
	rmCmd.PersistentFlags().StringSliceVarP(&Exclusions, "exclude", "x", []string{}, "Exclude files and folders from being removed using -x selector1 -x selector2,selector3. So if you want to remove all files from a folder recursively except for a few files in those folders")
	rmCmd.PersistentFlags().BoolVarP(&ShowExcluded, "show-excluded", "s", false, "Show list of excluded files which are NOT removed")
	rmCmd.PersistentFlags().BoolVarP(&AsSelector, "as-selector", "a", false, "Use the passed argument(s) as a regex to match over committed files rather than paths or exact names. All files matching any of the arguments are staged for removal")
	rmCmd.PersistentFlags().BoolVarP(&Unstage, "unstage", "u", false, "Don't only stage committed files for removal, but also unstage any matching staged adds/updates (does not undo staged removes!)")
}

var rmCmd = &cobra.Command{
	Use:     "rm [idvname/file/folder]...",
	Aliases: []string{"remove"},
	Short:   "Stage committed files for removal",
	Long:    `Stages files to be removed from the dataset. If existing paths are given, these files will be converted to internal formatting and removed from commit if they match any. if locally non-existent paths are passed, iterum uses them as literal names and tries to remove them from the dataset.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errNotEnoughArgs
		}
		if AsSelector || Unstage {
			var invalids []string
			format := "Invalid selector args:\n"
			for _, arg := range args {
				if !isValidSelector(arg) {
					invalids = append(invalids, arg)
					format += fmt.Sprintf("%v\n", arg)
				}
			}
			if len(invalids) > 0 {
				return fmt.Errorf(format)
			}
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
	var err error
	var removals, unstaged int
	if AsSelector {
		selector := buildSelector(args)
		removals, unstaged, err = idv.RemoveWithSelector(selector, Unstage)
	} else {
		paths, names := getPaths(args)
		allFiles := getAllFiles(paths)
		whitelisted := exclude(allFiles)
		removals, unstaged, err = idv.RemoveFiles(whitelisted, names, Unstage)
		fmt.Printf("GOT %v potential files\n", len(names)+len(whitelisted))
	}
	if err != nil {
		log.Fatal(err)
	}

	if Unstage {
		fmt.Printf("UNSTAGED %v file(s)\n", unstaged)
	}
	fmt.Printf("REMOVED %v file(s)\n", removals)
}
