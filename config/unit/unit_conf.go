package config

import (
	"errors"
	"regexp"

	common "github.com/Mantsje/iterum-cli/config/common"
)

// UnitConf contains the config for a unit folder in an iterum project
type UnitConf struct {
	Name     string
	RepoType common.RepoType
	UnitType UnitType
	Git      common.GitConf
}

// NewUnitConf instantiates a new UnitConf and sets up defaults
func NewUnitConf(name string) UnitConf {
	return UnitConf{
		Name:     name,
		RepoType: common.Unit,
	}
}

// IsValid validates all elements of the UnitConf
func (uc UnitConf) IsValid() error {
	rexp, err := regexp.Compile("[ \t\n\r]")
	if err == nil && rexp.ReplaceAllString(uc.Name, "") != uc.Name {
		err = errors.New("Error: Name of unit contains whitespace which is illegal")
	}
	err = common.Verify(uc.RepoType, err)
	err = common.Verify(uc.UnitType, err)
	err = common.Verify(uc.Git, err)
	return err
}
