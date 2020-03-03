package prompter

import (
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
)

// Name asks the user for input to be used as a name
func Name() string {
	return prompt(promptui.Prompt{
		Label: "Enter a name",
		Validate: func(input string) error {
			input = strings.TrimSpace(input)
			if len(input) < 4 {
				return errIndiscriptiveName
			}
			rexp, _ := regexp.Compile("[ \t\n\r]")
			if rexp.ReplaceAllString(input, "") != input {
				return errContainsWhitespace
			}
			return nil
		},
	})
}
