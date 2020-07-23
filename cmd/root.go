// Package cmd is the root of all CLI command specified using the excellent spf13/cobra package.
// Its elegant structure is used to create an intuitive command line interface.
package cmd

import (
	"github.com/prometheus/common/log"

	"github.com/iterum-provenance/cli/cmd/data"
	"github.com/iterum-provenance/cli/cmd/pipeline"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "iterum",
	Short: "Iterum is a provenance tracking pipeline deployment tool",
	Long:  "The ideal tool for tracking any type of provenance for (distributed) (data science) pipelines aimed at supporting academic research",
}

func init() {
	rootCmd.AddCommand(data.GetRootCmd())
	rootCmd.AddCommand(pipeline.GetRootCmd())
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
