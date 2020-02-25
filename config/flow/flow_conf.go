package config

import (
	"errors"
	"regexp"

	common "github.com/Mantsje/iterum-cli/config/common"
)

// FlowConf contains the config for a flow folder in an iterum project
type FlowConf struct {
	name     string
	repoType common.RepoType
	git      common.GitConf
}

// NewFlowConf instantiates a new FlowConf and sets up defaults
func NewFlowConf(name string) FlowConf {
	return FlowConf{
		name:     name,
		repoType: common.Flow,
	}
}

// IsValid validates all elements of the FlowConf
func (fc FlowConf) IsValid() error {
	rexp, err := regexp.Compile("[ \t\n\r]")
	if err == nil && rexp.ReplaceAllString(fc.name, "") != fc.name {
		err = errors.New("Error: Name of flow contains whitespace which is illegal")
	}
	err = common.Verify(fc.repoType, err)
	err = common.Verify(fc.git, err)
	return err
}
