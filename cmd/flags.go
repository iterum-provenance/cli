package cmd

// List of flags accessible throughout commands
// Not all flags are shared by all commands, but they need to be global vars
// This means no duplication and so if some of them share we need a shared file

// RawURL contains the raw url found in the optional url flag
var RawURL string

// NoRemote is a boolean flag used to skip pushing to remote git
var NoRemote bool
