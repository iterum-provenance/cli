package cmd

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/Mantsje/iterum-cli/config/git"
	"github.com/Mantsje/iterum-cli/config/project"
	"github.com/Mantsje/iterum-cli/config/unit"
	"github.com/manifoldco/promptui"
)

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

var namePrompt = promptui.Prompt{
	Label: "Enter a name",
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

var unitTypePrompt = promptui.Select{
	Label: "What kind of unit will this be",
	Items: []unit.UnitType{
		unit.ProcessingUnit,
		unit.UploadingUnit,
		unit.DownloadingUnit,
	},
}

var projectTypePrompt = promptui.Select{
	Label: "In what setting will this project run",
	Items: []project.ProjectType{
		project.Local,
		project.Distributed,
	},
}

var gitProtocolPrompt = promptui.Select{
	Label: "Which git protocol should be used",
	Items: []git.GitProtocol{
		git.SSH,
		git.HTTPS,
	},
}

var gitPlatformPrompt = promptui.Select{
	Label: "Which git platform should be used",
	Items: []git.GitPlatform{
		git.Github,
		git.Gitlab,
		git.Bitbucket,
	},
}
