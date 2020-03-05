package data

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(checkoutCmd)
	checkoutCmd.PersistentFlags().BoolVarP(&IsCommit, "commit", "c", false, "Passed arg is a commit refference rather than a branch")
	checkoutCmd.PersistentFlags().BoolVarP(&IsHash, "hash", "h", false, "Passed arg is a hash rather than a name")
}

var checkoutCmd = &cobra.Command{
	Use:   "checkout [target]",
	Short: "Checkout of current commit onto a new one",
	Long:  `Change the current commit to another one based on either the hash or name of a branch/commit`,
	Args:  cobra.ExactArgs(1),
	Run:   checkoutRun,
}

func checkoutRun(cmd *cobra.Command, args []string) {
	log.Println("`iterum data checkout`")
}
