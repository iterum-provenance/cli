// Package pipeline contains all the cobra subcommands found in `iterum pipeline [subcommand]` of the CLI.
// It deals with pipeline deployment and results analysis/interaction. The cobra commands mostly function as
// a simple interface calling API endpoints in the main /pipeline package
package pipeline

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "The subset of commands corresponding to pipeline deployment and management for iterum",
	Long:  "The `pipeline` sub-command gives access to creating, updating and managing (remotely) deployed pipelines",
}

// GetRootCmd returns the root of the data subcommand
func GetRootCmd() *cobra.Command {
	return rootCmd
}
