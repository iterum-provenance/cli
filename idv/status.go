package idv

import (
	"regexp"
	"strings"

	"github.com/Mantsje/iterum-cli/idv/ctl"
)

// Contains all functionality that is status related

// Status returns information about the currently staged files
func Status(fullPath, localPath bool) (report string, err error) {
	defer _returnErrOnPanic(&err)()
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
	defer _returnErrOnPanic(&err)()
	EnsureByPanic(EnsureLOCAL, "")
	var local Commit
	parseLOCAL(&local)
	report = local.FormatFiles(selector, "{\n\t", "\n}", "< Empty Data Set >", "\n\t", fullPath)
	return
}

// LsBranches returns a string report specifying the list of currently known branches
func LsBranches() (report string, err error) {
	defer _returnErrOnPanic(&err)()
	EnsureByPanic(EnsureTREE, "")
	EnsureByPanic(EnsureBRANCH, "")
	var history VTree
	var branch Branch
	parseTREE(&history)
	parseBRANCH(&branch)
	report = ""
	for key, val := range history.Branches {
		if key == branch.Hash {
			report += "* "
		}
		report += val + strings.Repeat(" ", 20-len(val)) + key.String() + "\n"
	}
	return
}

// LsCommits returns a string report specifying the list of commits on the current branch
func LsCommits() (report string, err error) {
	defer _returnErrOnPanic(&err)()
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
		report += node.Name + strings.Repeat(" ", 36-len(node.Name)) + hash.String() + "\n"
		if hash == node.Parent { // Root has parent as itself, stop infinite loop
			break
		}
		hash = node.Parent
		node, _ = history.Tree[hash]
	}
	return
}

// Inspect returns the config file found in the same folder used for adding/committing and pusing data etc
func Inspect() (report string, err error) {
	defer _returnErrOnPanic(&err)()
	EnsureByPanic(EnsureConfig, "")
	var ctl ctl.DataCTL
	ctl.ParseFromFile(configPath) // No error is ensured, so no need to catch it
	report += "Data configuration:\n"
	report += "{\n"
	report += "\tName: " + ctl.Name + "\n"
	report += "\tBackend: " + ctl.Backend.String() + "\n"
	report += "\tLocation: " + ctl.GetStorageLocation() + "\n"
	report += "}\n"
	return
}
