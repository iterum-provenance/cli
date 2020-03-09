package project

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"

	"github.com/Mantsje/iterum-cli/config"
	"github.com/Mantsje/iterum-cli/util"
)

// ProjectConf contains the config for the root folder of an iterum project
type ProjectConf struct {
	config.Conf
	ProjectType ProjectType
	Registered  map[string]config.RepoType // map from name to type for each separete version controlled part of the project (think submodules+root)
}

// NewProjectConf creates a new ProjectConf instance and sets up defaults
func NewProjectConf(name string) ProjectConf {
	var pc = ProjectConf{
		Conf:       config.NewConf(name, config.Project),
		Registered: make(map[string]config.RepoType),
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
		err = util.Verify(val, err)
	}
	err = util.Verify(pc.Git, err)
	err = util.Verify(pc.ProjectType, err)
	err = util.Verify(pc.RepoType, err)
	return err
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

// ParseFromFile tries to parse a config file into this ProjectConfig
func (pc *ProjectConf) ParseFromFile(filepath string) error {
	if err := util.ReadJSONFile(filepath, pc); err != nil {
		return errors.New("Error: Could not parse ProjectConf")
	}
	if err := pc.IsValid(); err != nil {
		return err
	}
	return nil
}
