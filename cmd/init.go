package cmd

import (
	"net/url"
	"os"
	"os/exec"

	"github.com/Mantsje/iterum-cli/cmd/prompter"
	"github.com/Mantsje/iterum-cli/config"
	"github.com/Mantsje/iterum-cli/consts"
	"github.com/Mantsje/iterum-cli/git"
	"github.com/Mantsje/iterum-cli/util"
	"github.com/prometheus/common/log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.PersistentFlags().BoolVarP(&FromURL, "from-url", "u", false, "Pull from existing repo rather than creating")
	initCmd.PersistentFlags().BoolVarP(&NoRemote, "no-remote", "", false, "flag to use if git should not be initialized remotely")
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Iterum component",
	Long:  `This creates or pulls a new Iterum component and initializes it based on some basic questions to set up`,
	Run:   initRun,
}

// Write the new config to disk
func writeConfig(conf config.Configurable) error {
	base := conf.GetBaseConf()

	err := util.WriteJSONFile(base.Name+"/"+consts.ConfigFilePath, conf)
	if err != nil {
		return errConfigWriteFailed
	}

	if err != nil { // Erase work if we failed
		os.RemoveAll("./" + base.Name)
	}
	return err
}

func initCreate() (name string, gitConf config.GitConf, err error) {
	_, _, errAlreadyAComponent := ensureIterumComponent("")
	if errAlreadyAComponent == nil { // Meaning this is already an iterum component
		err = errComponentNesting
		return
	}

	name = prompter.Name()

	// if to be created component folder already exists
	if _, errExist := os.Stat(name); errExist == nil {
		_, _, errAlreadyAComponent = ensureIterumComponent(name)
		if errAlreadyAComponent == nil { // Meaning THIS is already an iterum component
			err = errComponentNesting
			return
		}
	}
	createComponentFolder(name)

	var gitPlatform git.Platform
	var gitProtocol git.Protocol
	if NoRemote {
		gitPlatform = git.None
		gitProtocol = git.HTTPS
	} else {
		gitPlatform, _ = git.NewPlatform(prompter.Platform())
		gitProtocol, _ = git.NewProtocol(prompter.Protocol())
	}

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
		os.Exit(0)
	}

	return
}

func finalizeCreate(conf config.Configurable) {
	err := writeConfig(conf)
	if err != nil {
		log.Fatal(err.Error())
	}
	initVersionTracking(conf)
}

func initRun(cmd *cobra.Command, args []string) {
	_, _, err := ensureIterumComponent("")
	if err == nil { // Meaning this is already an iterum component
		log.Fatal(errComponentNesting)
	}

	componentType, _ := config.NewRepoType(prompter.RepoType())
	switch componentType {
	case config.Data:
		initDataRun(cmd, args)
	case config.Unit:
		initUnitRun(cmd, args)
	case config.Flow:
		initFlowRun(cmd, args)
	}
}
