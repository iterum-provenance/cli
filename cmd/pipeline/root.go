package pipeline

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "The subset of commands corresponding to pipeline deployment and management for iterum",
	Long:  "The `pipeline` sub-command gives access to creating, updating and managing (remotely) deployed pipelines",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&ManagerURL, "url", "u", "http://localhost:3001", "URL at which the manager can be reached")
}

// GetRootCmd returns the root of the data subcommand
func GetRootCmd() *cobra.Command {
	return rootCmd
}
