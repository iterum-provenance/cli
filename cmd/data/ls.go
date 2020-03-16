package data

import (
	"fmt"
	"log"
	"regexp"

	"github.com/Mantsje/iterum-cli/idv"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(lsCmd)
	lsCmd.PersistentFlags().BoolVarP(&ShowFullPath, "full-path", "f", false, "Show entire path rather than just filename")
	lsCmd.PersistentFlags().BoolVarP(&Branches, "branches", "b", false, "Show a list of branches instead of files, selector is then ignored")
	lsCmd.PersistentFlags().BoolVarP(&Commits, "commits", "c", false, "Show a list of commits instead of files, selector is then ignored")
}

var lsCmd = &cobra.Command{
	Use:   "ls [selector]",
	Short: "List files in the current commit (that match with selector)",
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
	var report string
	var err error
	if Branches {
		report, err = idv.LsBranches()
	} else if Commits {
		report, err = idv.LsCommits()
	} else {
		var selector *regexp.Regexp
		if len(args) == 0 {
			selector, _ = regexp.Compile("")
		} else {
			selector, _ = regexp.Compile(args[0])
		}
		report, err = idv.Ls(selector, ShowFullPath)
		report = "Files in data set:\n" + report
	}
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(report)
}
