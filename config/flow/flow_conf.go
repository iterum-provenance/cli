package flow

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"

	"github.com/Mantsje/iterum-cli/config"
	"github.com/Mantsje/iterum-cli/util"
)

// FlowConf contains the config for an Iterum flow component
type FlowConf struct {
	config.Conf
}

// NewFlowConf instantiates a new FlowConf and sets up defaults
func NewFlowConf(name string) FlowConf {
	return FlowConf{
		Conf: config.NewConf(name, config.Flow),
	}
}

// IsValid validates all elements of the FlowConf
func (fc FlowConf) IsValid() error {
	rexp, err := regexp.Compile("[ \t\n\r]")
	if err == nil && rexp.ReplaceAllString(fc.Name, "") != fc.Name {
		err = errors.New("Error: Name of flow contains whitespace which is illegal")
	}
	err = util.Verify(fc.RepoType, err)
	err = util.Verify(fc.Git, err)
	return err
}

// AllowedVariables returns a formatted string on how to set this type with the set command
func (fc FlowConf) AllowedVariables() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "\n")
	fmt.Fprintf(&buf, "Name                string\n")
	fmt.Fprintf(&buf, fc.Git.AllowedVariables())
	return buf.String()
}

// ParseFromFile tries to parse a config file into this FlowConfig
func (fc *FlowConf) ParseFromFile(filepath string) error {
	if err := util.ReadJSONFile(filepath, fc); err != nil {
		return errors.New("Error: Could not parse FlowConf")
	}
	if err := fc.IsValid(); err != nil {
		return err
	}
	return nil
}
