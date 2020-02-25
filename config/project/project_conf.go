package config

import (
	"errors"
	"regexp"

	common "github.com/Mantsje/iterum-cli/config/common"
)

// ProjectConf contains the config for the root folder of an iterum project
type ProjectConf struct {
	name              string
	repoType          common.RepoType
	projectType       ProjectType
	git               common.GitConf
	registered        map[string]common.RepoType // map from name to type for each separete version controlled part of the project (think submodules+root)
	validDependencies map[string]bool            // a map of different external dependencies that were verified
}

// NewProjectConf creates a new ProjectConf instance and sets up defaults
func NewProjectConf(name string) ProjectConf {
	var pc = ProjectConf{
		name:              name,
		repoType:          common.Project,
		git:               common.NewGitConf(),
		registered:        make(map[string]common.RepoType),
		validDependencies: make(map[string]bool),
	}
	return pc
}

// IsValid checks the validity of the ProjectConf
func (pc ProjectConf) IsValid() error {
	rexp, err := regexp.Compile("[ \t\n\r]")
	if err == nil && rexp.ReplaceAllString(pc.name, "") != pc.name {
		err = errors.New("Error: Name of project contains whitespace which is illegal")
	}
	for _, val := range pc.registered {
		err = common.Verify(val, err)
	}
	err = common.Verify(pc.git, err)
	err = common.Verify(pc.projectType, err)
	err = common.Verify(pc.repoType, err)
	return err
}
