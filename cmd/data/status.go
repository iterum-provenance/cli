package data

import (
	"fmt"
	"log"

	"github.com/Mantsje/iterum-cli/idv"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "List information about currently stages updates, removes and additions",
	Long:  `Lists all files and the type of update staged for the upcoming commit and their possible conflicts`,
	Run:   statusRun,
}

func statusRun(cmd *cobra.Command, args []string) {
	report, err := idv.Status()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Staged file changes:")
	fmt.Println(report)
}
