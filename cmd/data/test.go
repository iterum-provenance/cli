package data

import (
	"fmt"
	"log"

	"github.com/Mantsje/iterum-cli/idv"
	"github.com/Mantsje/iterum-cli/idv/ctl"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(testCmd)
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test random things",
	Long:  `Test command, does nothing`,
	Run:   testRun,
}

func testRun(cmd *cobra.Command, args []string) {
	fmt.Println("'Iterum testing' command")
	var ctl ctl.DataCTL
	err := ctl.ParseFromFile("idv-config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(idv.PostDataset(ctl))
	fmt.Println(idv.PullVTree())
}
