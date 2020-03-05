package cmd

import (
	"log"
	"net/url"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/Mantsje/iterum-cli/cmd/prompter"
	"github.com/Mantsje/iterum-cli/config"
	"github.com/Mantsje/iterum-cli/config/project"
	"github.com/Mantsje/iterum-cli/git"
	"github.com/Mantsje/iterum-cli/util"
)

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.PersistentFlags().BoolVarP(&FromURL, "from-url", "u", false, "Pull from existing repo rather than creating")
	createCmd.PersistentFlags().BoolVarP(&NoRemote, "no-remote", "", false, "Skip initializing and pushing to remote repo")
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create elements for this Iterum project",
	Long:  `Create or pull new elements to add to this iterum project. Create units and flows`,
}

// Write the new config to disk and update the registered elements of the project config
func writeConfigAndUpdate(conf config.Configurable, project project.ProjectConf) error {
	c := conf.GetBaseConf()
	// Write config
	err := util.WriteJSONFile(c.Name+"/"+config.ConfigFilePath, conf)
	if err != nil {
		return errConfigWriteFailed
	}

	err = register(c.Name, c.RepoType, project)
	if err != nil { // Erase work if we failed
		os.RemoveAll("./" + c.Name)
	}
	return err
}

func initCreate() (proj project.ProjectConf, name string, gitConf config.GitConf, err error) {
	proj, err = ensureRootLocation()
	if err != nil {
		return
	}

	name = prompter.Name()
	if _, ok := proj.Registered[name]; ok {
		err = errAlreadyExists
		return
	}

	createComponentFolder(name)

	// Guaranteed to be correct, so no checking needed
	var gitPlatform, _ = git.NewPlatform(prompter.Platform())
	var gitProtocol, _ = git.NewProtocol(prompter.Protocol())

	gitConf = config.GitConf{
		Platform: gitPlatform,
		Protocol: gitProtocol,
	}

	if FromURL {
		var gitPath = prompter.GitRepoPath()
		git.CloneRepo(gitProtocol, gitPlatform, gitPath, name)
		uri, _ := url.Parse("https://" + gitPlatform.String() + ".com/" + gitPath)
		gitConf.URI = *uri
		rename := exec.Command("iterum", "set", "Name", name)
		util.RunCommand(rename, "./"+name, false)
		register := exec.Command("iterum", "register", name)
		util.RunCommand(register, "./", false)
		os.Exit(0)
	}

	return
}

func finalizeCreate(conf config.Configurable, project project.ProjectConf) {
	err := writeConfigAndUpdate(conf, project)
	if err != nil {
		log.Fatal(err.Error())
	}
	initVersionTracking(conf)
}
