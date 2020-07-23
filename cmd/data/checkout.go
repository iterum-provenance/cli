package data

import (
	"github.com/prometheus/common/log"

	"github.com/iterum-provenance/cli/idv"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(checkoutCmd)
	checkoutCmd.PersistentFlags().BoolVarP(&IsCommit, "commit", "c", false, "Passed argument is a commit reference rather than a branch")
	checkoutCmd.PersistentFlags().BoolVarP(&IsHash, "hash", "#", false, "Passed argument is a hash rather than a name")
}

var checkoutCmd = &cobra.Command{
	Use:   "checkout [target]",
	Short: "Checkout of current commit onto a new one",
	Long:  `Change the current commit to another one based on either the hash or name of a branch/commit`,
	Args:  cobra.ExactArgs(1),
	Run:   checkoutRun,
}

func checkoutRun(cmd *cobra.Command, args []string) {
	err := idv.Checkout(args[0], IsCommit, IsHash)
	if err != nil {
		log.Fatalln(err)
	}
}
