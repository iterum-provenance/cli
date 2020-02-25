package config

import (
	"errors"
	"regexp"

	common "github.com/Mantsje/iterum-cli/config/common"
)

// ProjectConf contains the config for the root folder of an iterum project
type ProjectConf struct {
	Name              string
	RepoType          common.RepoType
	ProjectType       ProjectType
	Git               common.GitConf
	Registered        map[string]common.RepoType // map from name to type for each separete version controlled part of the project (think submodules+root)
	ValidDependencies map[string]bool            // a map of different external dependencies that were verified
}

// NewProjectConf creates a new ProjectConf instance and sets up defaults
func NewProjectConf(name string) ProjectConf {
	var pc = ProjectConf{
		Name:              name,
		RepoType:          common.Project,
		Git:               common.NewGitConf(),
		Registered:        make(map[string]common.RepoType),
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
		err = common.Verify(val, err)
	}
	err = common.Verify(pc.Git, err)
	err = common.Verify(pc.ProjectType, err)
	err = common.Verify(pc.RepoType, err)
	return err
}
