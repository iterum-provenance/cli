package prompter

import (
	"log"

	"github.com/manifoldco/promptui"
)

// prompt asks the user for input
func prompt(prompt promptui.Prompt) string {
	result, err := prompt.Run()
	if err != nil {
		log.Fatal("Prompt failed due to: ", err)
	}
	return result
}

// pick asks the user to pick an option (select is keyword)
func pick(prompt promptui.Select) string {
	_, value, err := prompt.Run()
	if err != nil {
		log.Fatal("Select failed due to: ", err)
	}
	return value
}
