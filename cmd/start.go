package cmd

import (
	"fmt"

	"github.com/Mantsje/iterum-cli/container"
	"github.com/Mantsje/iterum-cli/deps"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the Iterum daemon",
	Long:  `Runs a (docker) process in the background that starts the Iterum daemon in a docker container`,
	Run:   startRun,
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops the Iterum daemon",
	Long:  `Kills the (docker) process in the background that is the Iterum daemon`,
	Run:   startRun,
}

func startRun(cmd *cobra.Command, args []string) {
	deps.EnsureDep(container.DockerDep)
	fmt.Println("`iterum start` not fully implemented yet")
}

func stopRun(cmd *cobra.Command, args []string) {
	deps.EnsureDep(container.DockerDep)
	fmt.Println("`iterum stop` not fully implemented yet")
}
