package pipeline

import (
	"net/url"

	"github.com/prometheus/common/log"

	"github.com/iterum-provenance/cli/pipeline"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(describeCmd)
	describeCmd.PersistentFlags().StringVarP(&ManagerURL, "url", "u", "http://localhost:3001", "URL at which the manager can be reached")
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
	url, err := url.Parse(ManagerURL)
	if err != nil {
		log.Fatalln(err)
	}
	err = pipeline.Describe(args[0], url)
	if err != nil {
		log.Fatalln(err)
	}
}
