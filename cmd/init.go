package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/Mantsje/iterum-cli/config"
	"github.com/Mantsje/iterum-cli/config/git"
	"github.com/Mantsje/iterum-cli/config/project"
	"github.com/Mantsje/iterum-cli/util"
)

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.PersistentFlags().BoolVarP(&NoRemote, "skip-remote", "", false, "flag to use if git should not be initialized remotely")
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Iterum project",
	Long:  `This creates a new Iterum project and asks some basic questions to set up that process`,
	Run:   initRun,
}

func initRun(cmd *cobra.Command, args []string) {
	var name string = runPrompt(namePrompt)
	if util.FileExists(config.ConfigFileName) || util.FileExists(name+"/"+config.ConfigFileName) {
		log.Fatal(errProjectNesting)
	}
	// Guaranteed to be correct, so no checking needed
	var projectType, _ = project.InferProjectType(runSelect(projectTypePrompt))
	var gitPlatform, _ = git.NewGitPlatform(runSelect(gitPlatformPrompt))
	var gitProtocol, _ = git.NewGitProtocol(runSelect(gitProtocolPrompt))

	var projectConfig = project.NewProjectConf(name)
	projectConfig.ProjectType = projectType
	projectConfig.Git.Platform = gitPlatform
	projectConfig.Git.Protocol = gitProtocol

	os.Mkdir("./"+name, 0755)
	err := util.JSONWriteFile(name+"/"+config.ConfigFileName, projectConfig)
	if err != nil {
		log.Fatal(errConfigWriteFailed)
	}
	if NoRemote {
		createRepo("Creation of Iterum project `"+name+"`", git.None, "./"+name)
	} else {
		uri := createRepo("Creation of Iterum project `"+name+"`", gitPlatform, "./"+name)
		projectConfig.Git.URI = uri
		err := util.JSONWriteFile(name+"/"+config.ConfigFileName, projectConfig)
		if err != nil {
			log.Fatal(errConfigWriteFailed)
		}
	}
}
