package pipeline

import (
	"errors"
	"fmt"
)

// List of all errors shared by possibly multiple files
var (
	errNotEnoughArgs   = errors.New("Not enough arguments passed")
	errTooManyArgs     = errors.New("Too many arguments passed")
	errInvalidLocation = errors.New("passed location is not a valid location/path")
)

func errWrongAmountOfArgs(nargs int) error {
	return fmt.Errorf("Wrong amount of arguments, wanted %d", nargs)
}

func errInvalidArgs(reason string, args ...string) error {
	tail := args[0]
	for _, val := range args[1:] {
		tail = tail + ", " + val
	}
	return errors.New("Invalid argument(s) specified: " + tail + "\n\t\tReason: " + reason + "\n")
}
