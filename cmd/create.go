package cmd

import (
	"log"
	"net/url"
	"os"

	"github.com/spf13/cobra"

	"github.com/Mantsje/iterum-cli/config"
	"github.com/Mantsje/iterum-cli/config/flow"
	"github.com/Mantsje/iterum-cli/config/git"
	"github.com/Mantsje/iterum-cli/config/project"
	"github.com/Mantsje/iterum-cli/config/unit"
	"github.com/Mantsje/iterum-cli/util"
)

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.AddCommand(createUnitCmd)
	createCmd.AddCommand(createFlowCmd)
	createCmd.PersistentFlags().StringVarP(&RawURL, "url", "u", "", "url to existing git page")
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
func writeConfigAndUpdate(name string, conf interface{}, repo config.RepoType, project project.ProjectConf) error {
	// Write config
	err := util.JSONWriteFile(name+"/"+config.ConfigFileName, conf)
	if err != nil {
		return errConfigWriteFailed
	}

	err = register(name, repo, project)
	if err != nil { // Erase work if we failed
		os.RemoveAll("./" + name)
	}
	return err
}

func gatherSharedRequirements() (proj project.ProjectConf, name string, gitConf git.GitConf, err error) {
	proj, err = ensureRootLocation()
	if err != nil {
		return
	}
	name = runPrompt(namePrompt)
	if _, ok := proj.Registered[name]; ok {
		err = errAlreadyExists
		return
	}
	os.Mkdir("./"+name, 0755)

	// Guaranteed to be correct, so no checking needed
	var gitPlatform, _ = git.NewGitPlatform(runSelect(gitPlatformPrompt))
	var gitProtocol, _ = git.NewGitProtocol(runSelect(gitProtocolPrompt))
	gitConf = git.GitConf{
		Platform: gitPlatform,
		Protocol: gitProtocol,
	}
	return
}

func createUnitRun(cmd *cobra.Command, args []string) {
	proj, name, gitConf, err := gatherSharedRequirements()
	if err != nil {
		log.Fatal(err.Error())
	}
	// Guaranteed to be okay through validation
	var unitType, _ = unit.NewUnitType(runSelect(unitTypePrompt))

	var unitConfig = unit.NewUnitConf(name)
	unitConfig.UnitType = unitType
	unitConfig.Git = gitConf

	if RawURL != "" { // TODO use this part
		gitURI, _ := url.Parse(RawURL) // Error is already guaranteed by validation
		unitConfig.Git.URI = *gitURI
	}

	err = writeConfigAndUpdate(name, unitConfig, config.Unit, proj)
	if err != nil {
		log.Fatal(err.Error())
	}
	if NoRemote {
		createRepo("Creation of Iterum project", git.None, "./"+name)
	} else {
		uri := createRepo("Creation of Iterum unit `"+name+"`", gitConf.Platform, "./"+name)
		unitConfig.Git.URI = uri
		err := util.JSONWriteFile(name+"/"+config.ConfigFileName, unitConfig)
		if err != nil {
			log.Fatal(errConfigWriteFailed)
		}
	}
}

func createFlowRun(cmd *cobra.Command, args []string) {
	proj, name, gitConf, err := gatherSharedRequirements()
	if err != nil {
		log.Fatal(err.Error())
	}

	var flowConfig = flow.NewFlowConf(name)
	flowConfig.Git = gitConf
	if RawURL != "" { // TODO use this part
		gitURI, _ := url.Parse(RawURL) // Error is already guaranteed by validation
		flowConfig.Git.URI = *gitURI
	}

	err = writeConfigAndUpdate(name, flowConfig, config.Flow, proj)
	if err != nil {
		log.Fatal(err.Error())
	}
	path := "./" + name
	if NoRemote {
		createRepo("Creation of Iterum project", git.None, path)
	} else {
		uri := createRepo("Creation of Iterum flow`"+name+"`", gitConf.Platform, path)
		flowConfig.Git.URI = uri
		err := util.JSONWriteFile(name+"/"+config.ConfigFileName, flowConfig)
		if err != nil {
			log.Fatal(errConfigWriteFailed)
		}
	}
}
