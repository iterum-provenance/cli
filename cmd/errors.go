package cmd

import "errors"

var (
	errNoComponent error = errors.New("Error: Either this folder is not (root of) an iterum component or the .conf file is corrupted")

	errIllegalUpdate error = errors.New("Error: Setting variable resulted in invalid conf, likely invalid value")

	errMalformedURL error = errors.New("Error: url flag could not be parsed")

	errComponentNesting error = errors.New("Error: Cannot create another component under an existing iterum component")
	errConfigNotFound   error = errors.New("Error: Could not find iterum.conf (for component)")

	errIndiscriptiveName  error = errors.New("Error: Name should be descriptive, alphanumeric, and (ideally) -(dash) separated")
	errContainsWhitespace error = errors.New("Error: Name contains whitespace which is illegal")

	errConfigWriteFailed error = errors.New("Error: Writing config to file failed, setting variable failed")
)
