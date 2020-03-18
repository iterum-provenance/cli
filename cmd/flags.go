package cmd

// List of flags accessible throughout commands
// Not all flags are shared by all commands, but they need to be global vars
// This means no duplication and so if some of them share we need a shared file

// FromURL states whether `iterum create` should clone an existing repo rather than creating a new one
var FromURL bool

// NoRemote is a boolean flag used to skip pushing to remote git
var NoRemote bool
