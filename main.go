package main

import (
	"github.com/Mantsje/iterum-cli/cmd"
	"github.com/Mantsje/iterum-cli/container"
	"github.com/Mantsje/iterum-cli/git"
)

// Registers the dependencies of the carious packages under the global deps.Dependencies slice
func registerDependencies() {
	git.RegisterDeps()
	container.RegisterDeps()
}

// see https://github.com/spf13/cobra for help
func main() {
	registerDependencies()
	cmd.Execute()
}
