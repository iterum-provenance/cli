package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/iterum-provenance/cli/config/flow"
)

func init() {
	initCmd.AddCommand(initFlowCmd)
}

var initFlowCmd = &cobra.Command{
	Use:   "flow",
	Short: "Initialize a new Iterum flow component",
	Long:  `Initialize or pull a new Iterum flow component`,
	Run:   initFlowRun,
}

func initFlowRun(cmd *cobra.Command, args []string) {
	name, gitConf, err := initCreate()
	if err != nil {
		log.Fatal(err.Error())
	}

	var flowConfig = flow.NewFlowConf(name)
	flowConfig.Git = gitConf

	finalizeCreate(&flowConfig)
}
