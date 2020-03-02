package unit

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"

	"github.com/Mantsje/iterum-cli/config"
)

// UnitConf contains the config for a unit folder in an iterum project
type UnitConf struct {
	config.Conf
	UnitType UnitType
}

// NewUnitConf instantiates a new UnitConf and sets up defaults
func NewUnitConf(name string) UnitConf {
	return UnitConf{
		Conf: config.NewConf(name, config.Unit),
	}
}

// IsValid validates all elements of the UnitConf
func (uc UnitConf) IsValid() error {
	rexp, err := regexp.Compile("[ \t\n\r]")
	if err == nil && rexp.ReplaceAllString(uc.Name, "") != uc.Name {
		err = errors.New("Error: Name of unit contains whitespace which is illegal")
	}
	err = config.Verify(uc.RepoType, err)
	err = config.Verify(uc.UnitType, err)
	err = config.Verify(uc.Git, err)
	return err
}

// AllowedVariables returns a formatted string on how to set this type with the set command
func (uc UnitConf) AllowedVariables() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "\n")
	fmt.Fprintf(&buf, "Name                string\n")
	fmt.Fprintf(&buf, uc.UnitType.AllowedVariables())
	fmt.Fprintf(&buf, uc.Git.AllowedVariables())
	return buf.String()
}
