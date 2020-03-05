package data

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(branchCmd)
}

var branchCmd = &cobra.Command{
	Use:   "branch [branch-name]",
	Short: "Branch the current (local) commit onto a new branch",
	Long:  `Create a new branch using the current commit as the root of this branch`,
	Args:  cobra.ExactArgs(1),
	Run:   branchRun,
}

func branchRun(cmd *cobra.Command, args []string) {
	log.Println("`iterum data branch`")
}
