package cmd

import (
	"log"
	"os"

	"github.com/iterum-provenance/cli/config/data"
	"github.com/iterum-provenance/cli/idv"
	"github.com/spf13/cobra"
)

func init() {
	initCmd.AddCommand(initDataCmd)
}

var initDataCmd = &cobra.Command{
	Use:   "data",
	Short: "Initialize a new iterum data repository component",
	Long:  `Initialize or pull a new data repository component. A data component is used to manage and version control an iterum enabled data set`,
	Run:   initDataRun,
}

func initDataRun(cmd *cobra.Command, args []string) {
	name, gitConf, err := initCreate()
	if err != nil {
		log.Fatal(err.Error())
	}

	var dataConfig = data.NewDataConf(name)
	dataConfig.Git = gitConf

	finalizeCreate(&dataConfig)

	os.Chdir("./" + name)
	idv.Initialize()
}
