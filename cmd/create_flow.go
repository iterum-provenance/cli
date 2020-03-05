package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/Mantsje/iterum-cli/config/flow"
)

func init() {
	createCmd.AddCommand(createFlowCmd)
}

var createFlowCmd = &cobra.Command{
	Use:   "flow",
	Short: "Create a new flow for this Iterum project",
	Long:  `Create or pull a new flow and add it to this iterum project`,
	Run:   createFlowRun,
}

func createFlowRun(cmd *cobra.Command, args []string) {
	proj, name, gitConf, err := initCreate()
	if err != nil {
		log.Fatal(err.Error())
	}

	var flowConfig = flow.NewFlowConf(name)
	flowConfig.Git = gitConf

	finalizeCreate(&flowConfig, proj)
}
