package parser

import (
	"errors"

	"github.com/Mantsje/iterum-cli/config"
	"github.com/Mantsje/iterum-cli/config/data"
	"github.com/Mantsje/iterum-cli/config/flow"
	"github.com/Mantsje/iterum-cli/config/project"
	"github.com/Mantsje/iterum-cli/config/unit"

	"github.com/Mantsje/iterum-cli/util"
)

// ParseUnitConfig tries to parse a config file into a UnitConfig
func ParseUnitConfig(filepath string) (unit.UnitConf, error) {
	var unit unit.UnitConf
	if err := util.ReadJSONFile(filepath, &unit); err != nil {
		return unit, errors.New("Error: Could not parse UnitConf")
	}
	if err := unit.IsValid(); err != nil {
		return unit, err
	}
	return unit, nil
}

// ParseFlowConfig tries to parse a config file into a FlowConfig
func ParseFlowConfig(filepath string) (flow.FlowConf, error) {
	var flow flow.FlowConf
	if err := util.ReadJSONFile(filepath, &flow); err != nil {
		return flow, errors.New("Error: Could not parse FlowConf")
	}
	if err := flow.IsValid(); err != nil {
		return flow, err
	}
	return flow, nil
}

// ParseProjectConfig tries to parse a config file into a ProjectConfig
func ParseProjectConfig(filepath string) (project.ProjectConf, error) {
	var project project.ProjectConf
	if err := util.ReadJSONFile(filepath, &project); err != nil {
		return project, errors.New("Error: Could not parse ProjectConf")
	}
	if err := project.IsValid(); err != nil {
		return project, err
	}
	return project, nil
}

// ParseDataConfig tries to parse a config file into a DataConfig
func ParseDataConfig(filepath string) (data.DataConf, error) {
	var data data.DataConf
	if err := util.ReadJSONFile(filepath, &data); err != nil {
		return data, errors.New("Error: Could not parse DataConf")
	}
	if err := data.IsValid(); err != nil {
		return data, err
	}
	return data, nil
}

// ParseConfigFile parses a configfile found at stringpath and
// returns the parsed object, the type and in case of problems an error
func ParseConfigFile(filepath string) (interface{}, config.RepoType, error) {
	if util.FileExists(filepath) {
		var repo = struct {
			RepoType config.RepoType
		}{}
		if err := util.ReadJSONFile(filepath, &repo); err != nil {
			return nil, "", errors.New("Error: Could not find repo type in config file")
		}
		switch repo.RepoType {
		case config.Unit:
			unit, err := ParseUnitConfig(filepath)
			return unit, repo.RepoType, err
		case config.Flow:
			flow, err := ParseFlowConfig(filepath)
			return flow, repo.RepoType, err
		case config.Project:
			project, err := ParseProjectConfig(filepath)
			return project, repo.RepoType, err
		case config.Data:
			data, err := ParseDataConfig(filepath)
			return data, repo.RepoType, err
		}
	}
	return nil, "", errors.New("Error: Could not find target file")
}
