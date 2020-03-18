package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/Mantsje/iterum-cli/git"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(syncCmd)
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Pulls the iterum component repo",
	Long:  `Updates the current iterum repo, basically an alias for git pull`,
	Run:   syncRun,
}

func syncRun(cmd *cobra.Command, args []string) {
	_, _, err := ensureIterumComponent()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("Syncing..." + strings.Repeat(" ", 20+3))
	git.PullRepo("")
	fmt.Println("Done")
}
