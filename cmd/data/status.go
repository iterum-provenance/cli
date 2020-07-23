package data

import (
	"fmt"
	"github.com/prometheus/common/log"

	"github.com/iterum-provenance/cli/idv"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.PersistentFlags().BoolVarP(&ShowFullPath, "full-path", "f", false, "Show entire path rather than just filename")
	statusCmd.PersistentFlags().BoolVarP(&ShowLocalPath, "local-path", "l", false, "Show full paths to locally staged files")
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "List information about currently stages updates, removes and additions",
	Long:  `Lists all files and the type of update staged for the upcoming commit and their possible conflicts`,
	Run:   statusRun,
}

func statusRun(cmd *cobra.Command, args []string) {
	report, err := idv.Status(ShowFullPath, ShowLocalPath)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Staged file changes:")
	fmt.Println(report)
}
