package cmd

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"github.com/Mantsje/iterum-cli/config"
	common_conf "github.com/Mantsje/iterum-cli/config/common"
	project_conf "github.com/Mantsje/iterum-cli/config/project"
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

var namePrompt = promptui.Prompt{
	Label: "Enter a name for your project",
	Validate: func(input string) error {
		input = strings.TrimSpace(input)
		if len(input) < 4 {
			return errors.New("Name should be descriptive, alphanumeric, and (ideally) -(dash) separated")
		}
		rexp, _ := regexp.Compile("[ \t\n\r]")
		if rexp.ReplaceAllString(input, "") != input {
			return errors.New("Name contains whitespace which is illegal")
		}
		return nil
	},
}

var projectTypePrompt = promptui.Select{
	Label: "In what setting will this project run",
	Items: []project_conf.ProjectType{
		project_conf.Local,
		project_conf.Distributed,
	},
}

var gitProtocolPrompt = promptui.Select{
	Label: "Which git protocol should be used",
	Items: []common_conf.GitProtocol{
		common_conf.SSH,
		common_conf.HTTPS,
	},
}

var gitPlatformPrompt = promptui.Select{
	Label: "Which git platform should be used",
	Items: []common_conf.GitPlatform{
		common_conf.Github,
		common_conf.Gitlab,
		common_conf.Bitbucket,
	},
}

func runPrompt(prompt promptui.Prompt) string {
	result, err := prompt.Run()
	if err != nil {
		fmt.Print("Prompt failed due to: ")
		fmt.Println(err)
		return ""
	}
	return result
}

func runSelect(prompt promptui.Select) string {
	_, value, err := prompt.Run()
	if err != nil {
		fmt.Print("Select failed due to: ")
		fmt.Println(err)
		return ""
	}
	return value
}

func initRun(cmd *cobra.Command, args []string) {
	var name string = runPrompt(namePrompt)
	if util.FileExists(configName) || util.FileExists(name+"/"+configName) {
		fmt.Println(errors.New("Error: current or ./*project-name* is already (part of) an iterum project, you cannot start another here"))
	} else {
		// Guaranteed to be correct, so no checking needed
		var projectType, _ = project_conf.ParseProjectType(runSelect(projectTypePrompt))
		var gitPlatform, _ = common_conf.NewGitPlatform(runSelect(gitPlatformPrompt))
		var gitProtocol, _ = common_conf.NewGitProtocol(runSelect(gitProtocolPrompt))

		var projectConfig = project_conf.NewProjectConf(name)
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
