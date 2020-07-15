package data

import (
	"fmt"
	"log"

	"github.com/iterum-provenance/cli/idv"
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
	err := idv.Initialize()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("IDV repo initialized")
	}
}
