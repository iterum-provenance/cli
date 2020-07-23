package data

import (
	"github.com/prometheus/common/log"

	"github.com/iterum-provenance/cli/idv"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Make the daemon mount the data set using the idv-config",
	Long:  `Applies the idv-config at the daemon when already initialized locally. Useful after rebooting or wiping the Daemon`,
	Run:   applyRun,
}

func applyRun(cmd *cobra.Command, args []string) {
	err := idv.Apply()
	if err != nil {
		log.Fatalln(err)
	}
}
