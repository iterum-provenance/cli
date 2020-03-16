package data

import (
	"errors"
	"log"
	"regexp"
	"strings"

	"github.com/Mantsje/iterum-cli/idv"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(commitCmd)
}

var commitCmd = &cobra.Command{
	Use:   "commit [name/tag] [description]",
	Short: "Commit applies staged changes to the (remote) store.",
	Long:  `Commit applies locally staged changes to the (remote) store of data via the iterum daemon. The name is the shorthand naming of the commit, where the description is a longer message describing the type of changes`,
	Run:   commitRun,
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if len(args) != 2 {
			return errNotEnoughArgs
		}
		input := strings.TrimSpace(args[0])
		rexp, _ := regexp.Compile("[ \t\n\r]")
		if rexp.ReplaceAllString(input, "") != input {
			return errors.New("Error: name of commit cannot contain whitespace")
		}
		if len(args[0]) > 32 {
			return errors.New("Error: name of commit should be <= 32 characters")
		}
		return nil
	},
}

func commitRun(cmd *cobra.Command, args []string) {
	err := idv.ApplyCommit(args[0], args[1])
	if err != nil {
		log.Fatal(err)
	}
}
