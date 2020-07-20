package idv

import (
	"fmt"
	"regexp"

	"github.com/iterum-provenance/cli/util"
)

// AddFiles stages new files to be added and existing files as Updates, expects a list of absolute file paths
func AddFiles(files []string) (adds, updates int, err error) {
	defer util.ReturnErrOnPanic(&err)()
	EnsureByPanic(EnsureSetup, "")
	EnsureByPanic(EnsureLOCALIsBranchHead, "")

	var local Commit
	parseLOCAL(&local)
	for _, file := range files {
		if !util.FileExists(file) {
			panic(fmt.Errorf("Error: %v is either non-existent, or a directory", file))
		}
	}
	addMap, updateMap := local.addOrUpdate(files)
	adds = len(addMap)
	updates = len(updateMap)
	var stagemap Stagemap
	parseStagemap(&stagemap)
	err = stagemap.update(addMap)
	util.PanicIfErr(err, "")
	err = stagemap.update(updateMap)
	util.PanicIfErr(err, "")
	verifyAndUpdateStagemap(local, stagemap)
	writeLOCAL(local)
	return
}

// RemoveFiles stages files for removal from the dataset
// files is expected to be a list of paths on this machine
// names can be random strings that are matched against names in the commit
func RemoveFiles(files []string, names []string, unstage bool) (removals, unstages int, err error) {
	defer util.ReturnErrOnPanic(&err)()
	EnsureByPanic(EnsureSetup, "")
	EnsureByPanic(EnsureLOCALIsBranchHead, "")

	var local Commit
	parseLOCAL(&local)
	for _, file := range files {
		if !util.FileExists(file) {
			panic(fmt.Errorf("Error: %v is either non-existent, or a directory", file))
		}
	}
	removals, unstages = local.removeFiles(files, unstage)
	removedNames, unstagedNames := local.removeNames(names, unstage)
	removals += removedNames
	unstages += unstagedNames
	verifyAndUpdateStagemap(local, nil)
	writeLOCAL(local)
	return
}

// RemoveWithSelector removes files from data set and staging based on a regex.
// All files matching the selector are staged for removal, or unstaged in case a staged files
func RemoveWithSelector(selector *regexp.Regexp, unstage bool) (removals, unstages int, err error) {
	defer util.ReturnErrOnPanic(&err)()
	EnsureByPanic(EnsureSetup, "")
	EnsureByPanic(EnsureLOCALIsBranchHead, "")
	var local Commit
	parseLOCAL(&local)
	removals, unstages = local.removeWithSelector(selector, unstage)
	verifyAndUpdateStagemap(local, nil)
	writeLOCAL(local)
	return
}

// Unstage unstages adds/updates/removes of files that match the selector
func Unstage(selector *regexp.Regexp) (unstaged int, err error) {
	defer util.ReturnErrOnPanic(&err)()
	EnsureByPanic(EnsureSetup, "")
	EnsureByPanic(EnsureLOCALIsBranchHead, "")

	var local Commit
	parseLOCAL(&local)
	unstaged = local.unstage(selector)
	verifyAndUpdateStagemap(local, nil)
	writeLOCAL(local)
	return
}
