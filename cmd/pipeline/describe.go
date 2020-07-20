package pipeline

import (
	"log"

	"github.com/iterum-provenance/cli/manager"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(describeCmd)
}

// Gather more detailed status information for a specific pipeline
var describeCmd = &cobra.Command{
	Use:   "describe [pipeline-hash]",
	Short: "List the pipeline run configuration of a specific pipeline",
	Long:  `List all the known information about a pipeline run configuration, echoing the deployed json structure`,
	Args:  cobra.ExactArgs(1),
	Run:   describeRun,
}

func describeRun(cmd *cobra.Command, args []string) {
	err := manager.Describe(args[0])
	if err != nil {
		log.Fatal(err)
	}
}
