// Package pipeline contains all the actual functionalities called from the `cmd/pipeline` combra-commands
// this package only contains the API calls that are called upon both the Iterum Manager and Daemon.
// Through here pipelines can be submitted, cancelled, results can be purged or retrieved and lineage information
// can be retrieved in order to generate lineage trees for each and every input and output element.
//
// requests.go contains the general functions for the various http calls.
package pipeline
