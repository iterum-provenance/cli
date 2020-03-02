package prompter

import (
	"github.com/Mantsje/iterum-cli/config/project"
	"github.com/Mantsje/iterum-cli/config/unit"
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

// ProjectType lets the user pick a ProjectType
func ProjectType() string {
	return pick(promptui.Select{
		Label: "In what setting will this project run",
		Items: []project.ProjectType{
			project.Local,
			project.Distributed,
		},
	})
}
