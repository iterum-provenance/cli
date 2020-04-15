package prompter

import (
	"errors"
	"strings"

	"github.com/iterum-provenance/cli/git"
	"github.com/manifoldco/promptui"
)

var (
	errSlashInPath error = errors.New("Error: Path should not start nor end with a '/'")
)

// Protocol asks the user to pick a Protocol
func Protocol() string {
	return pick(promptui.Select{
		Label: "Which git protocol should be used",
		Items: []git.Protocol{
			git.SSH,
			git.HTTPS,
		},
	})
}

// Platform asks the user to pick a Platform
func Platform() string {
	return pick(promptui.Select{
		Label: "Which git platform should be used",
		Items: []git.Platform{
			git.Github,
			git.Gitlab,
			git.Bitbucket,
		},
	})
}

// GitRepoPath asks the user to provide a valid path to a repo
func GitRepoPath() string {
	return prompt(promptui.Prompt{
		Label: "Please provide the path to the repo AFTER the platform so in https://github.com/User/repo-name, only provide User/repo-name",
		Validate: func(input string) error {
			input = strings.TrimSpace(input)
			if len(input) > 0 && (input[0] == '/' || input[len(input)-1] == '/') {
				return errSlashInPath
			}
			return nil
		},
	})
}
