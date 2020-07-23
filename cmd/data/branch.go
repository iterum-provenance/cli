package data

import (
	"github.com/prometheus/common/log"

	"github.com/iterum-provenance/cli/idv"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(branchCmd)
	branchCmd.PersistentFlags().StringVarP(&CommitHashOrName, "commit", "c", "", "A commit name or hash (needs --hash) to a specific commit that we want to branch off of")
	branchCmd.PersistentFlags().BoolVarP(&IsHash, "hash", "#", false, "If the value passed in [--commit -c] flag is a hash rather than a name (default) ")
}

var branchCmd = &cobra.Command{
	Use:   "branch [branch-name]",
	Short: "Branch the current (local) commit onto a new branch",
	Long:  `Create a new branch using the current commit as the root of this branch`,
	Args:  cobra.ExactArgs(1),
	Run:   branchRun,
}

func branchRun(cmd *cobra.Command, args []string) {
	err := idv.BranchFromCommit(args[0], CommitHashOrName, IsHash)
	if err != nil {
		log.Fatalln(err)
	}
}
