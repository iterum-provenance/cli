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
	var selector *regexp.Regexp
	if len(args) == 0 {
		selector, _ = regexp.Compile("")
	} else {
		selector, _ = regexp.Compile(args[0])
	}
	report, err := idv.Ls(selector, ShowFullPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Files in data set:")
	fmt.Println(report)
}
