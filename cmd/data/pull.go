package data

import (
	"log"

	"github.com/Mantsje/iterum-cli/idv"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pullCmd)
}

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "If nonexistent pulls HEAD of `master` branch, else pulls HEAD of current branch, always pulls vtree",
	Long:  "Pulls the HEAD version information of the current commit or in case of a first pull, it clones the HEAD of the `master`",
	Run:   pullRun,
}

func pullRun(cmd *cobra.Command, args []string) {
	err := idv.Pull()
	if err != nil {
		log.Fatal(err)
	}
}
