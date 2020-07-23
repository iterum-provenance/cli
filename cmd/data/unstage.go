package data

import (
	"fmt"
	"github.com/prometheus/common/log"

	"github.com/iterum-provenance/cli/idv"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(unstageCmd)
}

var unstageCmd = &cobra.Command{
	Use:   "unstage [selector]...",
	Short: "Unstage uncommitted files that match the passed selector(s)",
	Long:  `Unstages files that were staged for add/update/remove earlier that match the passed selector(s)`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errNotEnoughArgs
		}
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
		return nil
	},
	Run: unstageRun,
}

func unstageRun(cmd *cobra.Command, args []string) {
	var err error
	var unstaged int
	selector := buildSelector(args)
	unstaged, err = idv.Unstage(selector)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("UNSTAGED %v file(s)\n", unstaged)
}
