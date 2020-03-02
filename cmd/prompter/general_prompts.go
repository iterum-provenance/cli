package prompter

import (
	"errors"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
)

var (
	errIndiscriptiveName  error = errors.New("Error: Name should be descriptive, alphanumeric, and (ideally) -(dash) separated")
	errContainsWhitespace error = errors.New("Error: Name contains whitespace which is illegal")
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
