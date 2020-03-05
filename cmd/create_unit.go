package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/Mantsje/iterum-cli/cmd/prompter"
	"github.com/Mantsje/iterum-cli/config/unit"
)

func init() {
	createCmd.AddCommand(createUnitCmd)
}

var createUnitCmd = &cobra.Command{
	Use:   "unit",
	Short: "Create a new unit for this Iterum project",
	Long:  `Create or pull a new unit and add it to this iterum project`,
	Run:   createUnitRun,
}

func createUnitRun(cmd *cobra.Command, args []string) {
	proj, name, gitConf, err := initCreate()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Guaranteed to be okay through validation
	var unitType, _ = unit.NewUnitType(prompter.UnitType())
	var unitConfig = unit.NewUnitConf(name)
	unitConfig.UnitType = unitType
	unitConfig.Git = gitConf

	finalizeCreate(&unitConfig, proj)
}
