package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "iterum",
	Short: "Iterum is a provenance tracking pipeline deployment tool",
	Long:  "The ideal tool for tracking any type of provenance for (distributed) (data science) pipelines aimed at supporting academic research",
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
