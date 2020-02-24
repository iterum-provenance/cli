package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(setCmd)
	setCmd.SetUsageFunc(setUsage)
}

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Update values of variables",
	Long:  `Update the values of project/unit/flow variables based on the current value of $PWD`,
	Args:  cobra.ExactArgs(2),
	Run:   setRun,
}

func setRun(cmd *cobra.Command, args []string) {
	fmt.Println("'Iterum set' command")
	fmt.Println(args[0])
	fmt.Println(args[1])
}

func setUsage(cmd *cobra.Command) error {
	fmt.Println(`Usage:
	iterum set [variable] [value]

Flags: 
	-h, --help	help for set

Examples:
	iterum set git.protocol https`)
	return nil
}
