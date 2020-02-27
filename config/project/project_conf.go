package project

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"

	"github.com/Mantsje/iterum-cli/config"
	"github.com/Mantsje/iterum-cli/config/git"
)

// ProjectConf contains the config for the root folder of an iterum project
type ProjectConf struct {
	Name              string
	RepoType          config.RepoType
	ProjectType       ProjectType
	Git               git.GitConf
	Registered        map[string]config.RepoType // map from name to type for each separete version controlled part of the project (think submodules+root)
	ValidDependencies map[string]bool            // a map of different external dependencies that were verified
}

// NewProjectConf creates a new ProjectConf instance and sets up defaults
func NewProjectConf(name string) ProjectConf {
	var pc = ProjectConf{
		Name:              name,
		RepoType:          config.Project,
		Git:               git.NewGitConf(),
		Registered:        make(map[string]config.RepoType),
		ValidDependencies: make(map[string]bool),
	}
	return pc
}

// IsValid checks the validity of the ProjectConf
func (pc ProjectConf) IsValid() error {
	rexp, err := regexp.Compile("[ \t\n\r]")
	if err == nil && rexp.ReplaceAllString(pc.Name, "") != pc.Name {
		err = errors.New("Error: Name of project contains whitespace which is illegal")
	}
	for _, val := range pc.Registered {
		err = config.Verify(val, err)
	}
	err = config.Verify(pc.Git, err)
	err = config.Verify(pc.ProjectType, err)
	err = config.Verify(pc.RepoType, err)
	return err
}

// Set sets a field in this conf based on a string, rather than knowing the exact type
func (pc *ProjectConf) Set(variable []string, value interface{}) error {
	return config.SetField(pc, variable, value)
}

// AllowedVariables returns a formatted string on how to set this type with the set command
func (pc ProjectConf) AllowedVariables() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "\n")
	fmt.Fprintf(&buf, "Name                string\n")
	fmt.Fprintf(&buf, pc.ProjectType.AllowedVariables())
	fmt.Fprintf(&buf, pc.Git.AllowedVariables())
	return buf.String()
}
