package parser

import (
	"errors"

	"github.com/iterum-provenance/cli/config"
	"github.com/iterum-provenance/cli/config/data"
	"github.com/iterum-provenance/cli/config/flow"
	"github.com/iterum-provenance/cli/config/unit"

	"github.com/iterum-provenance/iterum-go/util"
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
		case config.Data:
			data := data.DataConf{}
			err = data.ParseFromFile(filepath)
			conf = data
		}
		return
	}
	return nil, "", errors.New("Error: Could not find target file")
}
