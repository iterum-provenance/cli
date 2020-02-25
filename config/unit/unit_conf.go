package config

import (
	"errors"
	"regexp"

	common "github.com/Mantsje/iterum-cli/config/common"
)

// UnitConf contains the config for a unit folder in an iterum project
type UnitConf struct {
	name     string
	repoType common.RepoType
	unitType UnitType
	git      common.GitConf
}

// NewUnitConf instantiates a new UnitConf and sets up defaults
func NewUnitConf(name string) UnitConf {
	return UnitConf{
		name:     name,
		repoType: common.Unit,
	}
}

// IsValid validates all elements of the UnitConf
func (uc UnitConf) IsValid() error {
	rexp, err := regexp.Compile("[ \t\n\r]")
	if err == nil && rexp.ReplaceAllString(uc.name, "") != uc.name {
		err = errors.New("Error: Name of unit contains whitespace which is illegal")
	}
	err = common.Verify(uc.repoType, err)
	err = common.Verify(uc.unitType, err)
	err = common.Verify(uc.git, err)
	return err
}
