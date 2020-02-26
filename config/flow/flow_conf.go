package flow

import (
	"errors"
	"regexp"

	"github.com/Mantsje/iterum-cli/config"
	"github.com/Mantsje/iterum-cli/config/git"
)

// FlowConf contains the config for a flow folder in an iterum project
type FlowConf struct {
	Name     string
	RepoType config.RepoType
	Git      git.GitConf
}

// NewFlowConf instantiates a new FlowConf and sets up defaults
func NewFlowConf(name string) FlowConf {
	return FlowConf{
		Name:     name,
		RepoType: config.Flow,
	}
}

// IsValid validates all elements of the FlowConf
func (fc FlowConf) IsValid() error {
	rexp, err := regexp.Compile("[ \t\n\r]")
	if err == nil && rexp.ReplaceAllString(fc.Name, "") != fc.Name {
		err = errors.New("Error: Name of flow contains whitespace which is illegal")
	}
	err = config.Verify(fc.RepoType, err)
	err = config.Verify(fc.Git, err)
	return err
}
