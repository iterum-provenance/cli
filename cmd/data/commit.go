package data

import (
	"log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(commitCmd)
}

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Commit changes to the (remote) store",
	Long:  `Commit locally staged changes to the (remote) store of data via the iterum daemon`,
	Run:   commitRun,
}

func commitRun(cmd *cobra.Command, args []string) {
	log.Println("`iterum data commit`")
}
