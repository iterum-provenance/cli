package data

import (
	"log"

	"github.com/iterum-provenance/cli/idv"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Applies the idv-config at the daemon when already initialized locally",
	Long:  `Applies the idv-config at the daemon when already initialized locally. Useful after rebooting or wiping the Daemon`,
	Run:   applyRun,
}

func applyRun(cmd *cobra.Command, args []string) {
	err := idv.Apply()
	if err != nil {
		log.Fatal(err)
	}
}
