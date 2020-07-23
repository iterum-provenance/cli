package pipeline

import (
	"net/url"

	"github.com/prometheus/common/log"

	"github.com/spf13/cobra"

	"github.com/iterum-provenance/cli/pipeline"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.PersistentFlags().StringVarP(&ManagerURL, "url", "u", "http://localhost:3001", "URL at which the manager can be reached")
}

// delete and delete all or specific pipelines
var deleteCmd = &cobra.Command{
	Use:   "delete [pipeline-hash?]",
	Short: "delete either a specific or all pipelines, keeps results",
	Long:  `Removes pipelines by killing all running jobs, does not wipe MinIO or RabbitMQ. Keeps data and results`,
	Args:  cobra.MaximumNArgs(1),
	Run:   deleteRun,
}

func deleteRun(cmd *cobra.Command, args []string) {
	url, err := url.Parse(ManagerURL)
	if err != nil {
		log.Fatalln(err)
	}
	if len(args) == 0 {
		err := pipeline.DeleteAllPipelines(url)
		if err != nil {
			log.Fatalln(err)
		}
	} else if len(args) == 1 {
		err := pipeline.DeletePipeline(args[0], url)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
