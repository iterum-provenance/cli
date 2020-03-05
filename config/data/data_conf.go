package data

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"

	"github.com/Mantsje/iterum-cli/config"
)

// DataConf contains the config for a data component in an iterum project
type DataConf struct {
	config.Conf
}

// NewDataConf instantiates a new DataConf and sets up defaults
func NewDataConf(name string) DataConf {
	return DataConf{
		Conf: config.NewConf(name, config.Data),
	}
}

// IsValid validates all elements of the DataConf
func (dc DataConf) IsValid() error {
	rexp, err := regexp.Compile("[ \t\n\r]")
	if err == nil && rexp.ReplaceAllString(dc.Name, "") != dc.Name {
		err = errors.New("Error: Name of data component contains whitespace which is illegal")
	}
	err = config.Verify(dc.RepoType, err)
	err = config.Verify(dc.Git, err)
	return err
}

// AllowedVariables returns a formatted string on how to set this type with the set command
func (dc DataConf) AllowedVariables() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "\n")
	fmt.Fprintf(&buf, "Name                string\n")
	fmt.Fprintf(&buf, dc.Git.AllowedVariables())
	return buf.String()
}
