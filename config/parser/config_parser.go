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

// ParseConfigFile parses a configfile found at stringpath and
// returns the parsed object, the type and in case of problems an error
func ParseConfigFile(filepath string) (conf interface{}, repo config.RepoType, err error) {
	if util.FileExists(filepath) {
		var repoWr = struct {
			RepoType config.RepoType
		}{}
		if err := util.ReadJSONFile(filepath, &repoWr); err != nil {
			return nil, "", errors.New("Error: Could not find repo type in config file")
		}
		repo = repoWr.RepoType
		switch repo {
		case config.Unit:
			unit := unit.UnitConf{}
			err = unit.ParseFromFile(filepath)
			conf = unit
		case config.Flow:
			flow := flow.FlowConf{}
			err = flow.ParseFromFile(filepath)
			conf = flow
		case config.Project:
			project := project.ProjectConf{}
			err = project.ParseFromFile(filepath)
			conf = project
		case config.Data:
			data := data.DataConf{}
			err = data.ParseFromFile(filepath)
			conf = data
		}
		return
	}
	return nil, "", errors.New("Error: Could not find target file")
}
