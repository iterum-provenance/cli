package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/Mantsje/iterum-cli/config"
	"github.com/Mantsje/iterum-cli/config/data"
	"github.com/Mantsje/iterum-cli/config/flow"
	"github.com/Mantsje/iterum-cli/config/parser"
	"github.com/Mantsje/iterum-cli/config/project"
	"github.com/Mantsje/iterum-cli/config/unit"
	"github.com/Mantsje/iterum-cli/util"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(setCmd)
	setCmd.AddCommand(lsCmd)
	setCmd.SetUsageFunc(setUsage)
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List possible variables for current folder",
	Long:  `Lists the possible variables of project/unit/flow variables that can be set based on the current path`,
	Run:   lsRun,
}

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Update values of variables",
	Long:  `Update the values of project/unit/flow variables based on the current value of $PWD`,
	Args:  cobra.ExactArgs(2),
	Run:   setRun,
}

func writeConf(conf util.Validatable) {
	if err := conf.IsValid(); err != nil {
		log.Fatal(errIllegalUpdate)
	}
	if err := util.WriteJSONFile(config.ConfigFilePath, conf); err != nil {
		log.Fatal(errConfigWriteFailed)
	}
}

func setRun(cmd *cobra.Command, args []string) {
	_conf, repo, err := parser.ParseConfigFile(config.ConfigFilePath)
	if err != nil {
		log.Fatal(err)
	}
	variable := strings.Split(args[0], ".")
	value := args[1]

	var conf config.Settable
	var roConf util.Validatable

	switch repo {
	case config.Unit:
		u := _conf.(unit.UnitConf)
		conf = &u
		roConf = &u
	case config.Flow:
		f := _conf.(flow.FlowConf)
		conf = &f
		roConf = &f
	case config.Project:
		p := _conf.(project.ProjectConf)
		conf = &p
		roConf = &p
	case config.Data:
		d := _conf.(data.DataConf)
		conf = &d
		roConf = &d
	}
	err = conf.Set(variable, value)
	if err != nil {
		log.Fatal(err)
	}
	writeConf(roConf)
}

func lsRun(cmd *cobra.Command, args []string) {
	_conf, repo, err := parser.ParseConfigFile(config.ConfigFilePath)
	if err != nil {
		log.Fatal(err)
	}
	switch repo {
	case config.Unit:
		u := _conf.(unit.UnitConf)
		fmt.Println("\nUnit config found, the following variables can be set:")
		fmt.Println(u.AllowedVariables())
	case config.Flow:
		f := _conf.(flow.FlowConf)
		fmt.Println("\nFlow config found, the following variables can be set:")
		fmt.Println(f.AllowedVariables())
	case config.Project:
		p := _conf.(project.ProjectConf)
		fmt.Println("\nProject config found, the following variables can be set:")
		fmt.Println(p.AllowedVariables())
	case config.Data:
		d := _conf.(data.DataConf)
		fmt.Println("\nData config found, the following variables can be set:")
		fmt.Println(d.AllowedVariables())
	}
}

func setUsage(cmd *cobra.Command) error {
	fmt.Println(`Usage:
	iterum set [variable] [value]
	or
	iterum set ls

Flags: 
	-h, --help	help for set

Examples:
	iterum set git.protocol https`)
	return nil
}
