package pipeline

import (
	"net/url"

	"github.com/iterum-provenance/cli/pipeline"
	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(resultsCmd)
	resultsCmd.PersistentFlags().StringVarP(&DaemonURL, "url", "u", "http://localhost:3000", "URL at which the daemon can be reached")
}

// Gather more detailed status information for a specific pipeline
var resultsCmd = &cobra.Command{
	Use:   "results [pipeline-hash]",
	Short: "Returns a list of file names found in the result of a pipeline",
	Long:  `Returns a the files found in the results of a pipeline execution`,
	Args:  cobra.ExactArgs(1),
	Run:   resultsRun,
}

func resultsRun(cmd *cobra.Command, args []string) {
	url, err := url.Parse(DaemonURL)
	if err != nil {
		log.Fatalln(err)
	}
	err = pipeline.Results(args[0], url)
	if err != nil {
		log.Fatalln(err)
	}
}
