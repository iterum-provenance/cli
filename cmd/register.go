package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/Mantsje/iterum-cli/config"
	"github.com/Mantsje/iterum-cli/config/project"
	"github.com/Mantsje/iterum-cli/util"
)

func init() {
	rootCmd.AddCommand(registerCmd)
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "register an untracked unit/flow to this project",
	Long:  `register a unit/flow that is untracked within this project. When a repo is copied into this folder for example`,
	Args:  cobra.ExactArgs(1),
	Run:   registerRun,
}

func register(name string, repo config.RepoType, project project.ProjectConf) error {
	if _, ok := project.Registered[name]; ok {
		return errRegistrationClash
	}
	project.Registered[name] = repo
	err := util.JSONWriteFile(config.ConfigFileName, project)
	if err != nil {
		return errRegistrationFailed
	}
	return nil
}

func registerRun(cmd *cobra.Command, args []string) {
	proj, err := ensureRootLocation()
	if err != nil {
		log.Fatal(err)
	}
	configPath := "./" + args[0] + "/" + config.ConfigFileName
	if util.FileExists(configPath) {
		repo := struct {
			RepoType config.RepoType
		}{}
		err := util.JSONReadFile(configPath, &repo)
		if err != nil {
			log.Fatal(err)
		}
		if repo.RepoType == config.Project {
			log.Fatal(errProjectNesting)
		}
		register(args[0], repo.RepoType, proj)
	} else {
		log.Fatal(errConfigNotFound)
	}
}
