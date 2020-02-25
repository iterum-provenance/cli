package cmd

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
)

// RawURL contains the raw url found in the optional url flag
var RawURL string

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.AddCommand(createUnitCmd)
	createCmd.AddCommand(createFlowCmd)
	createCmd.PersistentFlags().StringVarP(&RawURL, "url", "u", "", "url to pull from")
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create elements to the Iterum project",
	Long:  `Create or pull new elements to add to this iterum project. Create units and flows`,
}

var createUnitCmd = &cobra.Command{
	Use:   "unit",
	Short: "Create a new unit to the Iterum project",
	Long:  `Create or pull a new unit and add it to this iterum project`,
	Args:  argsValidator,
	Run:   createUnitRun,
}

var createFlowCmd = &cobra.Command{
	Use:   "flow",
	Short: "Create a new flow to the Iterum project",
	Long:  `Create or pull a new flow and add it to this iterum project`,
	Args:  argsValidator,
	Run:   createFlowRun,
}

func createUnitRun(cmd *cobra.Command, args []string) {
	fmt.Println("'Iterum create unit' command")
}

func createFlowRun(cmd *cobra.Command, args []string) {
	fmt.Println("'Iterum create flow' command")
}

func urlFlagValidator(cmd *cobra.Command, args []string) error {
	if RawURL != "" {
		_, err := url.Parse(RawURL)
		if err != nil {
			return errors.New("Error: url flag could not be parsed")
		}
	}
	return nil
}

func argsValidator(cmd *cobra.Command, args []string) error {
	err := urlFlagValidator(cmd, args)
	if err == nil && len(args) != 1 {
		return errors.New("Error: Not enough arguments passed. Missing [name]")
	}
	return err
}
