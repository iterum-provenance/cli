package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/Mantsje/iterum-cli/config"
	"github.com/Mantsje/iterum-cli/config/git"
	"github.com/Mantsje/iterum-cli/config/project"
	"github.com/Mantsje/iterum-cli/util"
)

const configName = config.ConfigFileName

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Iterum project",
	Long:  `This creates a new Iterum project and asks some basic questions to set up that process`,
	Run:   initRun,
}

func initRun(cmd *cobra.Command, args []string) {
	var name string = runPrompt(namePrompt)
	if util.FileExists(configName) || util.FileExists(name+"/"+configName) {
		fmt.Println(errors.New("Error: current or ./*project-name* is already (part of) an iterum project, you cannot start another here"))
	} else {
		// Guaranteed to be correct, so no checking needed
		var projectType, _ = project.InferProjectType(runSelect(projectTypePrompt))
		var gitPlatform, _ = git.NewGitPlatform(runSelect(gitPlatformPrompt))
		var gitProtocol, _ = git.NewGitProtocol(runSelect(gitProtocolPrompt))

		var projectConfig = project.NewProjectConf(name)
		projectConfig.ProjectType = projectType
		projectConfig.Git.Platform = gitPlatform
		projectConfig.Git.Protocol = gitProtocol

		os.Mkdir("./"+name, 0755)
		err := util.JSONWriteFile(name+"/"+configName, projectConfig)
		if err != nil {
			fmt.Println("Error: Writing config to file failed, project creation failed")
		}
	}
}
