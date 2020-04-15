package prompter

import (
	"github.com/iterum-provenance/cli/config"
	"github.com/iterum-provenance/cli/config/run"
	"github.com/iterum-provenance/cli/config/unit"
	"github.com/manifoldco/promptui"
)

// UnitType lets the user pick a UnitType
func UnitType() string {
	return pick(promptui.Select{
		Label: "What kind of unit will this be",
		Items: []unit.UnitType{
			unit.ProcessingUnit,
			unit.UploadingUnit,
			unit.DownloadingUnit,
		},
	})
}

// RepoType lets the user pick a RepoType
func RepoType() string {
	return pick(promptui.Select{
		Label: "What kind of Iterum component will this be",
		Items: []config.RepoType{
			config.Data,
			config.Unit,
			config.Flow,
		},
	})
}

// RunType lets the user pick a RunType
func RunType() string {
	return pick(promptui.Select{
		Label: "In what setting will this be run",
		Items: []run.RunType{
			run.Local,
			run.Distributed,
		},
	})
}
