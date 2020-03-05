package cmd

import (
	"log"

	"github.com/Mantsje/iterum-cli/config/data"
	"github.com/spf13/cobra"
)

func init() {
	createCmd.AddCommand(createDataCmd)
}

var createDataCmd = &cobra.Command{
	Use:   "data",
	Short: "Create a new data components for this Iterum project",
	Long:  `Create or pull a new data component and add it to this iterum project. A data component is used to manage and version control an iterum enabled data set`,
	Run:   createDataRun,
}

func createDataRun(cmd *cobra.Command, args []string) {
	proj, name, gitConf, err := initCreate()
	if err != nil {
		log.Fatal(err.Error())
	}

	var dataConfig = data.NewDataConf(name)
	dataConfig.Git = gitConf

	finalizeCreate(&dataConfig, proj)
}
