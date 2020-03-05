package data

import (
	"log"

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
	log.Println("`iterum data status`")
}
