package cmd

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/Mantsje/iterum-cli/config"
	"github.com/Mantsje/iterum-cli/config/project"
	"github.com/Mantsje/iterum-cli/constants"
	"github.com/Mantsje/iterum-cli/util"
)

func init() {
	rootCmd.AddCommand(registerCmd)
	rootCmd.AddCommand(deregisterCmd)
	deregisterCmd.PersistentFlags().BoolVarP(&RemoveFiles, "rm-files", "", false, "Remove the folder called ./*deregistered-component*")
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Registers an untracked unit/flow to this project",
	Long:  `Registers a unit/flow that is untracked within this project. When a repo is copied into this folder for example`,
	Args:  cobra.ExactArgs(1),
	Run:   registerRun,
}

var deregisterCmd = &cobra.Command{
	Use:   "deregister",
	Short: "Deregisters a tracked unit/flow of this project",
	Long:  `Deregisters a unit/flow that is tracked within this project. When you've removed a folder for example`,
	Args:  cobra.ExactArgs(1),
	Run:   deregisterRun,
}

func register(name string, repo config.RepoType, project project.ProjectConf) error {
	if _, ok := project.Registered[name]; ok {
		return errRegistrationClash
	}
	project.Registered[name] = repo
	err := util.WriteJSONFile(constants.ConfigFilePath, project)
	if err != nil {
		return errRegistrationFailed
	}
	return nil
}

func deregister(name string, project project.ProjectConf) error {
	if _, ok := project.Registered[name]; ok {
		delete(project.Registered, name)
		err := util.WriteJSONFile(constants.ConfigFilePath, project)
		if err != nil {
			return errRegistrationFailed
		}
	} else {
		return errNotDeregisterable
	}
	return nil
}

func registerRun(cmd *cobra.Command, args []string) {
	proj, err := ensureRootLocation()
	if err != nil {
		log.Fatal(err)
	}
	configPath := "./" + args[0] + "/" + constants.ConfigFilePath
	if util.FileExists(configPath) {
		repo := struct {
			RepoType config.RepoType
		}{}
		err := util.ReadJSONFile(configPath, &repo)
		if err != nil {
			log.Fatal(err)
		}
		if repo.RepoType == config.Project {
			log.Fatal(errProjectNesting)
		}
		err = register(args[0], repo.RepoType, proj)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal(errConfigNotFound)
	}
}

func deregisterRun(cmd *cobra.Command, args []string) {
	proj, err := ensureRootLocation()
	if err != nil {
		log.Fatal(err)
	}
	err = deregister(args[0], proj)
	if err != nil {
		log.Fatal(err)
	}
	configPath := "./" + args[0] + "/" + constants.ConfigFilePath
	if util.FileExists(configPath) {
		if RemoveFiles {
			rmFiles := exec.Command("rm", "-rf", "./"+args[0])
			util.RunCommand(rmFiles, "./", false)
		} else {
			fmt.Println("Iterum does not remove the deregistered component's folder, so do this yourself or use --rm-files")
		}
	}
	fmt.Println("Component deregistered")
}
