package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Iterum",
	Long:  `This prints the version of the Iterum CLI`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Iterum provenance tracker v0.1")
	},
}
