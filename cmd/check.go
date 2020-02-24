package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(checkCmd)
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Validates Iterum's dependencies",
	Long:  `Checks whether dependency applications are installed and usable by Iterum`,
	Run:   checkRun,
}

func checkRun(cmd *cobra.Command, args []string) {
	fmt.Println("'Iterum check' command")
}
