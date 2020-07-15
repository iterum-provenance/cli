package idv

import (
	"log"
	"regexp"
	"strings"

	"github.com/iterum-provenance/cli/idv/ctl"
	"github.com/iterum-provenance/iterum-go/util"
)

// Contains all functionality that is status related

// Status returns information about the currently staged files
func Status(fullPath, localPath bool) (report string, err error) {
	defer util.ReturnErrOnPanic(&err)()
	EnsureByPanic(EnsureLOCAL, "")

	var local Commit
	parseLOCAL(&local)
	if localPath {
		var stagemap Stagemap
		parseStagemap(&stagemap)
		report = local.FormatDiff("{\n\t", "\n}", "< No Staged Files >", "\n\t", fullPath, stagemap)
	} else {
		report = local.FormatDiff("{\n\t", "\n}", "< No Staged Files >", "\n\t", fullPath, nil)
	}
	return
}

// Ls lists all data in the current commit
func Ls(selector *regexp.Regexp, fullPath bool) (report string, err error) {
	defer util.ReturnErrOnPanic(&err)()
	EnsureByPanic(EnsureLOCAL, "")
	var local Commit
	parseLOCAL(&local)
	report = local.FormatFiles(selector, "{\n\t", "\n}", "< Empty Data Set >", "\n\t", fullPath)
	return
}

// LsBranches returns a string report specifying the list of currently known branches
func LsBranches() (report string, err error) {
	defer util.ReturnErrOnPanic(&err)()
	EnsureByPanic(EnsureTREE, "")
	EnsureByPanic(EnsureBRANCH, "")
	var history VTree
	var branch Branch
	parseTREE(&history)
	parseBRANCH(&branch)
	report = ""
	for key, val := range history.Branches {
		hashOffset := 40
		if key == branch.Hash {
			report += "* "
			hashOffset -= 2
		}
		report += val + strings.Repeat(" ", hashOffset-len(val)) + key.String() + "\n"
	}
	return
}

// LsCommits returns a string report specifying the list of commits on the current branch
func LsCommits() (report string, err error) {
	defer util.ReturnErrOnPanic(&err)()
	EnsureByPanic(EnsureHEAD, "")
	EnsureByPanic(EnsureTREE, "")
	var head Commit
	var history VTree
	parseHEAD(&head)
	parseTREE(&history)
	report = ""
	branch := head.Branch
	hash := head.Hash
	node, _ := history.Tree[hash]
	for node.Branch == branch {
		report += node.Name + strings.Repeat(" ", 40-len(node.Name)) + hash.String() + "\n"
		if node.Parent == "" || hash == node.Parent {
			break
		}
		hash = node.Parent
		node, _ = history.Tree[hash]
	}
	return
}

// Inspect returns the config file found in the same folder used for adding/committing and pusing data etc
func Inspect() (report string, err error) {
	defer util.ReturnErrOnPanic(&err)()
	EnsureByPanic(EnsureConfig, "")
	var ctl ctl.DataCTL
	ctl.ParseFromFile(configPath) // No error is ensured, so no need to catch it
	remoteCtl, err := getConfig(ctl)
	log.Println(remoteCtl)
	log.Println(err)
	util.PanicIfErr(err, "")
	datasets, err := getDatasets(ctl)
	util.PanicIfErr(err, "")

	report += "Local Data configuration:\n"
	report += ctl.ToReport()
	report += "\nRemote configuration for this dataset:\n"
	report += remoteCtl.ToReport()
	report += "\nDatasets known to Daemon:\n"
	for _, set := range datasets {
		report += "   - " + set + "\n"
	}
	return
}
