package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/Mantsje/iterum-cli/cmd/prompter"
	"github.com/Mantsje/iterum-cli/config"
	"github.com/Mantsje/iterum-cli/config/project"
	"github.com/Mantsje/iterum-cli/git"
	"github.com/Mantsje/iterum-cli/util"
)

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.PersistentFlags().BoolVarP(&NoRemote, "no-remote", "", false, "flag to use if git should not be initialized remotely")
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Iterum project",
	Long:  `This creates a new Iterum project and asks some basic questions to set up that process`,
	Run:   initRun,
}

func initRun(cmd *cobra.Command, args []string) {
	var name string = prompter.Name()
	if util.FileExists(config.ConfigFilePath) || util.FileExists(name+"/"+config.ConfigFilePath) {
		log.Fatal(errProjectNesting)
	}
	// Guaranteed to be correct, so no checking needed
	var projectType, _ = project.InferProjectType(prompter.ProjectType())
	var gitPlatform, _ = git.NewPlatform(prompter.Platform())
	var gitProtocol, _ = git.NewProtocol(prompter.Protocol())

	var projectConfig = project.NewProjectConf(name)
	projectConfig.ProjectType = projectType
	projectConfig.Git.Platform = gitPlatform
	projectConfig.Git.Protocol = gitProtocol

	createComponentFolder(name)
	err := util.WriteJSONFile(name+"/"+config.ConfigFilePath, projectConfig)
	if err != nil {
		log.Fatal(errConfigWriteFailed)
	}

	initVersionTracking(&projectConfig)
}
