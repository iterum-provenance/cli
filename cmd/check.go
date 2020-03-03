package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Mantsje/iterum-cli/deps"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(checkCmd)
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Validates Iterum's dependencies",
	Long:  `Checks whether dependency applications are installed and usable by Iterum`,
	Run:   checkRun,
}

func formatDependency(dep deps.Dep) string {
	cmd := dep.Cmd
	name := dep.Name
	return name + strings.Repeat(" ", 20-len(name)) + cmd + strings.Repeat(" ", 20-len(cmd)) + strconv.FormatBool(dep.IsUsable()) + "\n"
}

func checkRun(cmd *cobra.Command, args []string) {
	fmt.Printf("Checking whether Iterum has access to all the commands used by this tool\n\n")
	fmt.Println("Name" + strings.Repeat(" ", 20-len("Name")) + "Command" + strings.Repeat(" ", 20-len("Command")) + "Success")
	fmt.Println("--------------------------------------------------")
	for _, dep := range deps.Dependencies {
		fmt.Printf(formatDependency(dep))
	}
	fmt.Println()
}
