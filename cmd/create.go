package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/Mantsje/iterum-cli/config"
	"github.com/Mantsje/iterum-cli/config/git"
	"github.com/Mantsje/iterum-cli/config/parser"
	"github.com/Mantsje/iterum-cli/config/unit"
	"github.com/Mantsje/iterum-cli/util"
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
	Short: "Create elements for this Iterum project",
	Long:  `Create or pull new elements to add to this iterum project. Create units and flows`,
}

var createUnitCmd = &cobra.Command{
	Use:   "unit",
	Short: "Create a new unit for this Iterum project",
	Long:  `Create or pull a new unit and add it to this iterum project`,
	Args:  urlFlagValidator,
	Run:   createUnitRun,
}

var createFlowCmd = &cobra.Command{
	Use:   "flow",
	Short: "Create a new flow for this Iterum project",
	Long:  `Create or pull a new flow and add it to this iterum project`,
	Args:  urlFlagValidator,
	Run:   createFlowRun,
}

func ensureLocation() error {
	_, repo, err := parser.ParseConfigFile(config.ConfigFileName)
	if err != nil {
		return errors.New("Error: either this folder is not an iterum project or the .conf file is corrupted")
	}
	if repo != config.Project {
		return errors.New("Error: current folder is not root of an Iterum project")
	}
	return nil
}

func createUnitRun(cmd *cobra.Command, args []string) {
	if err := ensureLocation(); err != nil {
		fmt.Println(err.Error())
		return
	}
	var name string = runPrompt(namePrompt)
	os.Mkdir("./"+name, 0755)

	// Guaranteed to be correct, so no checking needed
	var unitType, _ = unit.NewUnitType(runSelect(unitTypePrompt))
	var gitPlatform, _ = git.NewGitPlatform(runSelect(gitPlatformPrompt))
	var gitProtocol, _ = git.NewGitProtocol(runSelect(gitProtocolPrompt))

	var unitConfig = unit.NewUnitConf(name)
	unitConfig.UnitType = unitType
	unitConfig.Git.Platform = gitPlatform
	unitConfig.Git.Protocol = gitProtocol

	err := util.JSONWriteFile(name+"/"+config.ConfigFileName, unitConfig)
	if err != nil {
		fmt.Println("Error: Writing config to file failed, unit creation failed")
	}
}

func createFlowRun(cmd *cobra.Command, args []string) {
	if err := ensureLocation(); err != nil {
		fmt.Println(err.Error())
		return
	}
	var name string = runPrompt(namePrompt)
	os.Mkdir("./"+name, 0755)

	// Guaranteed to be correct, so no checking needed
	var gitPlatform, _ = git.NewGitPlatform(runSelect(gitPlatformPrompt))
	var gitProtocol, _ = git.NewGitProtocol(runSelect(gitProtocolPrompt))

	var flowConfig = unit.NewUnitConf(name)
	flowConfig.Git.Platform = gitPlatform
	flowConfig.Git.Protocol = gitProtocol

	err := util.JSONWriteFile(name+"/"+config.ConfigFileName, flowConfig)
	if err != nil {
		fmt.Println("Error: Writing config to file failed, flow creation failed")
	}
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
