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
	Short: "Print the version number of Iterum CLI",
	Long:  `This prints the version of the Iterum CLI`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Iterum CLI v0.1")
	},
}
