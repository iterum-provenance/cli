package config

import (
	"errors"
	"regexp"

	common "github.com/Mantsje/iterum-cli/config/common"
)

// FlowConf contains the config for a flow folder in an iterum project
type FlowConf struct {
	Name     string
	RepoType common.RepoType
	Git      common.GitConf
}

// NewFlowConf instantiates a new FlowConf and sets up defaults
func NewFlowConf(name string) FlowConf {
	return FlowConf{
		Name:     name,
		RepoType: common.Flow,
	}
}

// IsValid validates all elements of the FlowConf
func (fc FlowConf) IsValid() error {
	rexp, err := regexp.Compile("[ \t\n\r]")
	if err == nil && rexp.ReplaceAllString(fc.Name, "") != fc.Name {
		err = errors.New("Error: Name of flow contains whitespace which is illegal")
	}
	err = common.Verify(fc.RepoType, err)
	err = common.Verify(fc.Git, err)
	return err
}
