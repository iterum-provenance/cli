package cmd

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/Mantsje/iterum-cli/cmd/prompter"
	"github.com/Mantsje/iterum-cli/config"
	"github.com/Mantsje/iterum-cli/config/flow"
	"github.com/Mantsje/iterum-cli/config/project"
	"github.com/Mantsje/iterum-cli/config/unit"
	"github.com/Mantsje/iterum-cli/git"
	"github.com/Mantsje/iterum-cli/util"
)

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.AddCommand(createUnitCmd)
	createCmd.AddCommand(createFlowCmd)
	createCmd.PersistentFlags().StringVarP(&RawURL, "url", "u", "", "Valid existing git clone url to be used isntead of making a new component")
	createCmd.PersistentFlags().BoolVarP(&NoRemote, "no-remote", "", false, "Skip initializing and pushing to remote repo")
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

// Ensures that the url flag is a valid url that can be parsed
func urlFlagValidator(cmd *cobra.Command, args []string) error {
	if RawURL != "" {
		_, err := url.Parse(RawURL)
		if err != nil {
			return errMalformedURL
		}
	}
	return nil
}

// Write the new config to disk and update the registered elements of the project config
func writeConfigAndUpdate(conf config.Configurable, project project.ProjectConf) error {
	c := conf.GetBaseConf()
	fmt.Println(c)
	// Write config
	err := util.JSONWriteFile(c.Name+"/"+config.ConfigFileName, conf)
	if err != nil {
		return errConfigWriteFailed
	}

	err = register(c.Name, c.RepoType, project)
	if err != nil { // Erase work if we failed
		os.RemoveAll("./" + c.Name)
	}
	return err
}

func initShared() (proj project.ProjectConf, name string, gitConf config.GitConf, err error) {
	proj, err = ensureRootLocation()
	if err != nil {
		return
	}

	name = prompter.Name()
	if _, ok := proj.Registered[name]; ok {
		err = errAlreadyExists
		return
	}

	if RawURL != "" {
		uri, _ := url.Parse(RawURL) // Error is already guaranteed by validation
		gitURI := *uri
		git.CloneRepo(gitURI, name)
		register := exec.Command("iterum", "register", name)
		rename := exec.Command("iterum", "set", "Name", name)
		util.RunCommand(register, "./")
		util.RunCommand(rename, "./"+name)
		os.Exit(0)
	}

	os.Mkdir("./"+name, 0755)

	// Guaranteed to be correct, so no checking needed
	var gitPlatform, _ = config.NewGitPlatform(prompter.GitPlatform())
	var gitProtocol, _ = config.NewGitProtocol(prompter.GitProtocol())

	gitConf = config.GitConf{
		Platform: gitPlatform,
		Protocol: gitProtocol,
	}
	return
}

func finalizeShared(conf config.Configurable, project project.ProjectConf) {
	err := writeConfigAndUpdate(conf, project)
	if err != nil {
		log.Fatal(err.Error())
	}
	initVersionTracking(conf)
}

func createUnitRun(cmd *cobra.Command, args []string) {
	proj, name, gitConf, err := initShared()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Guaranteed to be okay through validation
	var unitType, _ = unit.NewUnitType(prompter.UnitType())
	var unitConfig = unit.NewUnitConf(name)
	unitConfig.UnitType = unitType
	unitConfig.Git = gitConf

	finalizeShared(&unitConfig, proj)
}

func createFlowRun(cmd *cobra.Command, args []string) {
	proj, name, gitConf, err := initShared()
	if err != nil {
		log.Fatal(err.Error())
	}

	var flowConfig = flow.NewFlowConf(name)
	flowConfig.Git = gitConf

	finalizeShared(&flowConfig, proj)
}
