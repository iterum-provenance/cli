package cmd

import "errors"

var (
	errNoProject error = errors.New("Error: Either this folder is not (part of) an iterum project or the .conf file is corrupted")
	errNotRoot   error = errors.New("Error: Current folder is not root of an Iterum project")

	errIllegalUpdate error = errors.New("Error: Setting variable resulted in invalid conf, likely invalid value")

	errAlreadyExists error = errors.New("Error: folder with this name already exists in this project, if it is an iterum component use `iterum register`")
	errMalformedURL  error = errors.New("Error: url flag could not be parsed")

	errRegistrationClash  error = errors.New("Error: This name is already registered in this project")
	errRegistrationFailed error = errors.New("Error: (De)Registering project component failed")
	errNotDeregisterable  error = errors.New("Error: Could not deregister component, likely because it was not registered to begin with")
	errProjectNesting     error = errors.New("Error: Cannot register another project under a project")
	errConfigNotFound     error = errors.New("Error: Could not find iterum.conf (for component)")

	errIndiscriptiveName  error = errors.New("Error: Name should be descriptive, alphanumeric, and (ideally) -(dash) separated")
	errContainsWhitespace error = errors.New("Error: Name contains whitespace which is illegal")

	errConfigWriteFailed error = errors.New("Error: Writing config to file failed, setting variable failed")
)
