package cmd

import (
	"fmt"

	"github.com/Mantsje/iterum-cli/config/parser"
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
	conf, repoType, err := parser.ParseConfigFile("iterum.conf")
	fmt.Println(err)
	fmt.Println(conf)
	fmt.Println(repoType)
}
