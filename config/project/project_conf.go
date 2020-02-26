package project

import (
	"errors"
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
