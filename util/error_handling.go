package util

import (
	"errors"
	"log"
)

// ReturnFirstErr returns the first not-nil error from a list of errors.
// Used to prevent many copies of if err != nil { return err } when they are indepedent
func ReturnFirstErr(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

// LogIfError logs an error if it occurs, rather than repeating this structure everywhere
func LogIfError(err error) {
	if err != nil {
		log.Println(err)
	}
}

// PanicIfErr panics in case of an error. If custom message it creates a new error based on that and returns that instead
func PanicIfErr(err error, customMessage string) {
	if err != nil {
		if customMessage != "" {
			panic(errors.New(customMessage))
		} else {
			panic(err)
		}
	}
}
