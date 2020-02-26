package config

import (
	"errors"

	common "github.com/Mantsje/iterum-cli/config/common"
	flow "github.com/Mantsje/iterum-cli/config/flow"
	project "github.com/Mantsje/iterum-cli/config/project"
	unit "github.com/Mantsje/iterum-cli/config/unit"

	"github.com/Mantsje/iterum-cli/util"
)

// ParseUnitConfig tries to parse a config file into a UnitConfig
func ParseUnitConfig(filepath string) (unit.UnitConf, error) {
	var unit unit.UnitConf
	if err := util.JSONReadFile(filepath, &unit); err != nil {
		return unit, errors.New("Error: Could not parse UnitConf")
	}
	if err := unit.IsValid(); err != nil {
		return unit, err
	}
	return unit, nil
}

// ParseFlowConfig tries to parse a config file into a UnitConfig
func ParseFlowConfig(filepath string) (flow.FlowConf, error) {
	var flow flow.FlowConf
	if err := util.JSONReadFile(filepath, &flow); err != nil {
		return flow, errors.New("Error: Could not parse FlowConf")
	}
	if err := flow.IsValid(); err != nil {
		return flow, err
	}
	return flow, nil
}

// ParseProjectConfig tries to parse a config file into a UnitConfig
func ParseProjectConfig(filepath string) (project.ProjectConf, error) {
	var project project.ProjectConf
	if err := util.JSONReadFile(filepath, &project); err != nil {
		return project, errors.New("Error: Could not parse ProjectConf")
	}
	if err := project.IsValid(); err != nil {
		return project, err
	}
	return project, nil
}

// ParseConfigFile parses a configfile found at stringpath and
// returns the parsed object, the type and in case of problems an error
func ParseConfigFile(filepath string) (interface{}, common.RepoType, error) {
	if util.FileExists(filepath) {
		var repo = struct {
			RepoType common.RepoType
		}{}
		if err := util.JSONReadFile(filepath, &repo); err != nil {
			return nil, "", errors.New("Error: Could not find repo type in config file")
		}
		switch repo.RepoType {
		case common.Unit:
			unit, err := ParseUnitConfig(filepath)
			return unit, repo.RepoType, err
		case common.Flow:
			flow, err := ParseFlowConfig(filepath)
			return flow, repo.RepoType, err
		case common.Project:
			project, err := ParseProjectConfig(filepath)
			return project, repo.RepoType, err
		}
	}
	return nil, "", errors.New("Error: Could not find target file")
}
