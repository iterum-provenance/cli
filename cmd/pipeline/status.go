package pipeline

import (
	"net/url"

	"github.com/prometheus/common/log"

	"github.com/iterum-provenance/cli/pipeline"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.PersistentFlags().StringVarP(&ManagerURL, "url", "u", "http://localhost:3001", "URL at which the manager can be reached")
}

// Gather status information for deployed pipelines
var statusCmd = &cobra.Command{
	Use:   "status [pipeline-hash?]",
	Short: "List pipeline deployments known to pipeline. If argument, gather more detailed info for that pipeline",
	Long:  `List all pipelines that are deployed in either completed, errored of running status. If an argument is passed, it is treated as a pipeline hash and information is gathered for that specific pipeline on a more detailed level.`,
	Args:  cobra.MaximumNArgs(1),
	Run:   statusRun,
}

func statusRun(cmd *cobra.Command, args []string) {
	url, err := url.Parse(ManagerURL)
	if err != nil {
		log.Fatalln(err)
	}
	if len(args) == 0 {
		err := pipeline.Status(url)
		if err != nil {
			log.Fatalln(err)
		}
	} else if len(args) == 1 {
		err := pipeline.PipelineStatus(args[0], url)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
