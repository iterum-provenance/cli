package pipeline

import (
	"log"

	"github.com/iterum-provenance/cli/manager"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

// Gather status information for deployed pipelines
var statusCmd = &cobra.Command{
	Use:   "status [pipeline-hash?]",
	Short: "List pipeline deployments known to manager. If argument, gather more detailed info for that pipeline",
	Long:  `List all pipelines that are deployed in either completed, errored of running status. If an argument is passed, it is treated as a pipeline hash and information is gathered for that specific pipeline on a more detailed level.`,
	Args:  cobra.MaximumNArgs(1),
	Run:   statusRun,
}

func statusRun(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		err := manager.Status()
		if err != nil {
			log.Fatal(err)
		}
	} else if len(args) == 1 {
		err := manager.PipelineStatus(args[0])
		if err != nil {
			log.Fatal(err)
		}
	}
}
