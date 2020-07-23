// Package data contains all the subcommands found in `iterum data [subcommand]` of the CLI.
// It deals with data versioning. Creating and updating commits. Creating branches, inspecting
// data sets. It prepares the arguments by validating file paths and filtering certain sets of files.
// This is especially useful for staging changes to commits. Many of these helper functions can be found
// in `fileutils.go`
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
