// Package util contains may different smaller bits and pieces that help the other packages.
// Notably are read and write from/to files for YAML and JSON structures.
// Furthermore some more file checking measures and the most used one: errors.go for some shorthand error handling
//
// Used throughout many parts of Iterum is the ReturnErrOnPanic(&returnValueError) and PanicIfErr(err, "") combination
// The former ensures that upon a panic an error is expected. Instead of panicking the error is returned from the function
// This avoids many if err != nil checks by collapsing it into a single statement
package util
