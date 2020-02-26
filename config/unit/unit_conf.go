package unit

import (
	"errors"
	"regexp"

	"github.com/Mantsje/iterum-cli/config"
	"github.com/Mantsje/iterum-cli/config/git"
)

// UnitConf contains the config for a unit folder in an iterum project
type UnitConf struct {
	Name     string
	RepoType config.RepoType
	UnitType UnitType
	Git      git.GitConf
}

// NewUnitConf instantiates a new UnitConf and sets up defaults
func NewUnitConf(name string) UnitConf {
	return UnitConf{
		Name:     name,
		RepoType: config.Unit,
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
