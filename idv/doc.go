// Package idv includes all the actual functionalities related to data versioning.
// Both relevant types and functions are included in this package
//
// Some notable helper files:
//	* constants.go		constants used throughout this package
//  * manage.go			contains many basic functions such as parsing often used files, symlinking those files and other general functioanlities
//  * checks.go 		contains all kinds of different checks, existence checking, verification, state checking, etc.
//  * daemon.go 		API calls executed upon the Daemon. used for all Daemon interaction
//
// Relevant type descriptions:
//  * hash.go			Hash type, simple randomly generated strings
//  * commit.go			Complex structure containing all information and functions related to commits. Most complex file
//  * branch.go			Structure containing a branch and its information
//  * vtree.go			Structure containing an entire version tree and all its references to braches and commits and their structure
//  * stagemap.go		Stagemap holds file references for changes staged for commits. Mapping staged files to the actual paths on the user's machine
//
// The other files contain the subcommand handlers that can be called from cmd/data
// The best way to understand this package is by using the CLI and seeing where its calls trace throughout this package.
// Especially the interaction of the commit.go file is core to many of the functionalities
package idv
