package prompter

import (
	"github.com/Mantsje/iterum-cli/config"
	"github.com/manifoldco/promptui"
)

// GitProtocol asks the user to pick a GitProtocol
func GitProtocol() string {
	return pick(promptui.Select{
		Label: "Which git protocol should be used",
		Items: []config.GitProtocol{
			config.SSH,
			config.HTTPS,
		},
	})
}

// GitPlatform asks the user to pick a GitPlatform
func GitPlatform() string {
	return pick(promptui.Select{
		Label: "Which git platform should be used",
		Items: []config.GitPlatform{
			config.Github,
			config.Gitlab,
			config.Bitbucket,
		},
	})
}
