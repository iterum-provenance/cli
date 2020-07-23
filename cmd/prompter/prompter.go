package prompter

import (
	"errors"
	"github.com/prometheus/common/log"

	"github.com/manifoldco/promptui"
)

var (
	errIndiscriptiveName  error = errors.New("Error: Name should be descriptive, alphanumeric, and (ideally) -(dash) separated")
	errContainsWhitespace error = errors.New("Error: Name contains whitespace which is illegal")
	errTooLong            error = errors.New("Error: Text is too long")
)

// prompt asks the user for input
func prompt(prompt promptui.Prompt) string {
	result, err := prompt.Run()
	if err != nil {
		log.Fatalln("Prompt failed due to: ", err)
	}
	return result
}

// pick asks the user to pick an option (select is keyword)
func pick(prompt promptui.Select) string {
	_, value, err := prompt.Run()
	if err != nil {
		log.Fatalln("Select failed due to: ", err)
	}
	return value
}
