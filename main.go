package main

import (
	"fmt"

	"github.com/Mantsje/iterum-cli/cmd"
	unit "github.com/Mantsje/iterum-cli/config/unit"
)

// see https://github.com/spf13/cobra for help
func main() {
	// pc := project.NewProjectConf("name-name")
	puc := unit.NewProcessingUnitConf("processing-unit")
	err := puc.IsValid()
	fmt.Println(err)

	cmd.Execute()
}
