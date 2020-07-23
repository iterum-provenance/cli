package data

import (
	"fmt"
	"github.com/prometheus/common/log"

	"github.com/iterum-provenance/cli/cmd/prompter"
	"github.com/iterum-provenance/cli/idv"
	"github.com/iterum-provenance/cli/idv/ctl/storage"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init creates a local empty, unbound IDV repository",
	Long:  `Init initializes a local empty, unbound IDV repository, if there is an idv-config.yaml, this can be used afterwards for that data set`,
	Run:   initRun,
}

func initRun(cmd *cobra.Command, args []string) {
	fmt.Println("Please input the name of the data set")
	name := prompter.Name()
	fmt.Println("Please input a short description for the data set")
	description := prompter.Description()
	backend := storage.Backend(prompter.StorageType())
	if err := backend.IsValid(); err != nil {
		log.Fatalln(err)
	}

	err := idv.Initialize(name, description, backend)
	if err != nil {
		log.Fatalln(err)
	} else {
		fmt.Println("IDV repo initialized")
	}
}
