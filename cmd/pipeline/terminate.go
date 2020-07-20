package pipeline

import (
	"log"
	"net/url"

	"github.com/spf13/cobra"

	"github.com/iterum-provenance/cli/manager"
)

func init() {
	rootCmd.AddCommand(terminateCmd)
}

// Terminate and delete all or specific pipelines
var terminateCmd = &cobra.Command{
	Use:   "terminate [pipeline-hash?]",
	Short: "Terminate and delete either 1 or all pipelines",
	Long:  `Removes pipelines by killing all running jobs, does not wipe MinIO or RabbitMQ`,
	Args:  cobra.MaximumNArgs(1),
	Run:   terminateRun,
}

func terminateRun(cmd *cobra.Command, args []string) {
	url, err := url.Parse(ManagerURL)
	if err != nil {
		log.Fatalln(err)
	}
	if len(args) == 0 {
		err := manager.TerminateAllPipelines(url)
		if err != nil {
			log.Fatal(err)
		}
	} else if len(args) == 1 {
		err := manager.TerminatePipeline(args[0], url)
		if err != nil {
			log.Fatal(err)
		}
	}
}
