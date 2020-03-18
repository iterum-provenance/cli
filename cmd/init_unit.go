package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/Mantsje/iterum-cli/cmd/prompter"
	"github.com/Mantsje/iterum-cli/config/unit"
)

func init() {
	initCmd.AddCommand(initUnitCmd)
}

var initUnitCmd = &cobra.Command{
	Use:   "unit",
	Short: "Initialize a new Iterum unit component",
	Long:  `Initialize or pull a new Iterum unit component`,
	Run:   initUnitRun,
}

func initUnitRun(cmd *cobra.Command, args []string) {
	name, gitConf, err := initCreate()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Guaranteed to be okay through validation
	var unitType, _ = unit.NewUnitType(prompter.UnitType())
	var unitConfig = unit.NewUnitConf(name)
	unitConfig.UnitType = unitType
	unitConfig.Git = gitConf

	finalizeCreate(&unitConfig)
}
