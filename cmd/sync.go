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
	syncCmd.PersistentFlags().BoolVarP(&CurrentComponentOnly, "current-only", "c", false, "Only sync the current working dir instead of entire project")
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Pulls the main repo and all its submodules",
	Long:  `Updates each of the proejct's units, flows, etc`,
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
	if !CurrentComponentOnly {
		project, err := ensureRootLocation()
		if err != nil {
			log.Fatal("Not in root of Iterum project")
		}
		// Pull current repo and all its registered components
		for key := range project.Registered {
			fmt.Print("Syncing `" + key + "`..." + strings.Repeat(" ", 20-len(key)))
			git.PullRepo(key)
			fmt.Println("Done")
		}
	}
}
