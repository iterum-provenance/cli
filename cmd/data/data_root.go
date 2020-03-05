package data

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "data",
	Short: "The subset of commands corresponding to data versioning for iterum",
	Long:  "The `data` sub-command gives access to updating and adjusting remotely stored datasets in a git-like fashion",
}

// GetRootCmd returns the root of the data subcommand
func GetRootCmd() *cobra.Command {
	return rootCmd
}
