package pipeline

import (
	"fmt"
	"net/url"

	"github.com/iterum-provenance/cli/pipeline"
	"github.com/iterum-provenance/cli/util"
	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(submitCmd)
	submitCmd.PersistentFlags().StringVarP(&ManagerURL, "url", "u", "http://localhost:3001", "URL at which the manager can be reached")
}

var submitCmd = &cobra.Command{
	Use:   "submit [pipeline.json]",
	Short: "Submit a pipeline configuration to the manager for deployment",
	Long:  `Instruct the manager to deploy a pipeline matching the specifications in 'pipeline.json'`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errWrongAmountOfArgs(1)
		}
		if !util.FileExists(args[0]) {
			return fmt.Errorf("'%v' is not an existing file", args[0])
		}
		return nil
	},
	Run: submitRun,
}

func submitRun(cmd *cobra.Command, args []string) {
	url, err := url.Parse(ManagerURL)
	if err != nil {
		log.Fatalln(err)
	}
	err = pipeline.SubmitPipeline(args[0], url)
	if err != nil {
		log.Fatalln(err)
	}
}
