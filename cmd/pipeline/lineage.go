package pipeline

import (
	"net/url"
	"os"

	"github.com/iterum-provenance/cli/pipeline"
	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(lineageCmd)
	lineageCmd.PersistentFlags().StringVarP(&DaemonURL, "url", "u", "http://localhost:3000", "URL at which the daemon can be reached")
}

// Gather more detailed status information for a specific pipeline
var lineageCmd = &cobra.Command{
	Use:   "lineage [pipeline-hash] [dirpath]",
	Short: "Download lineage information of a pipeline",
	Long:  `Downloads all know lineage information of a pipeline to the passed dir`,
	Args:  cobra.ExactArgs(2),
	Run:   lineageRun,
}

func lineageRun(cmd *cobra.Command, args []string) {
	url, err := url.Parse(DaemonURL)
	if err != nil {
		log.Fatalln(err)
	}
	info, err := os.Stat(args[1])
	if err != nil {
		log.Fatalln(err)
	}
	if !info.IsDir() {
		log.Fatalln("Passed path is not a directory")
	}
	err = pipeline.Lineage(args[0], args[1], url)
	if err != nil {
		log.Fatalln(err)
	}
}
