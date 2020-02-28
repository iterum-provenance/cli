package cmd

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/Mantsje/iterum-cli/config"
	"github.com/Mantsje/iterum-cli/config/parser"
	"github.com/Mantsje/iterum-cli/config/project"
)

// Package for shared command functions

// Make sure we are in an iterum project root
func ensureRootLocation() (project.ProjectConf, error) {
	conf, repo, err := parser.ParseConfigFile(config.ConfigFileName)
	if err != nil {
		return project.ProjectConf{}, errNoProject
	}
	if repo != config.Project {
		return project.ProjectConf{}, errNotRoot
	}
	return conf.(project.ProjectConf), nil
}

// Run an arbitrary command as if you were in a terminal at the given path
func runCommand(cmd *exec.Cmd, path string) string {
	cmd.Dir = path
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(string(out))
	return string(out)
}
