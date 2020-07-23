package pipeline

import (
	"net/url"

	"github.com/prometheus/common/log"

	"github.com/spf13/cobra"

	"github.com/iterum-provenance/cli/manager"
)

func init() {
	rootCmd.AddCommand(purgeCmd)
	purgeCmd.PersistentFlags().StringVarP(&ManagerURL, "url", "u", "http://localhost:3001", "URL at which the manager can be reached")
}

// purge and purge all or specific pipelines
var purgeCmd = &cobra.Command{
	Use:   "purge [pipeline-hash]",
	Short: "purge a specific pipeline, including Kubernetes jobs and results",
	Long:  `Purges a pipeline by killing all running jobs and wiping results from daemon, does not wipe MinIO or RabbitMQ`,
	Args:  cobra.ExactArgs(1),
	Run:   purgeRun,
}

func purgeRun(cmd *cobra.Command, args []string) {
	url, err := url.Parse(ManagerURL)
	if err != nil {
		log.Fatalln(err)
	}

	err = manager.PurgePipeline(args[0], url)
	if err != nil {
		log.Fatalln(err)
	}
}
