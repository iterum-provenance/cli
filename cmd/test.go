package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Mantsje/iterum-cli/config"
	"github.com/Mantsje/iterum-cli/git"
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
	url := git.CreateRepo("\"Initial commit for unit\"", config.Github, "./uniterum")
	fmt.Println(url)
}
