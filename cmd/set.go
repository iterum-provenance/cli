package cmd

import (
	"fmt"
	"strings"

	"github.com/Mantsje/iterum-cli/config"
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

func writeConf(conf config.Validatable) {
	if err := conf.IsValid(); err != nil {
		fmt.Println("Error: Setting variable resulted in invalid conf, likely invalid value")
		return
	}
	if err := util.JSONWriteFile(config.ConfigFileName, conf); err != nil {
		fmt.Println("Error: Writing config to file failed, setting variable failed")
	}
}

func setRun(cmd *cobra.Command, args []string) {
	_conf, repo, err := parser.ParseConfigFile(config.ConfigFileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	variable := strings.Split(args[0], ".")
	value := args[1]
	var conf config.Settable
	var roConf config.Validatable

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
	}
	err = conf.Set(variable, value)
	if err != nil {
		fmt.Println(err.Error())
	}
	writeConf(roConf)
}

func lsRun(cmd *cobra.Command, args []string) {
	_conf, repo, err := parser.ParseConfigFile(config.ConfigFileName)
	if err != nil {
		fmt.Println(err)
		return
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
