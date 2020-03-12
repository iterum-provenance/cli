package idv

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/Mantsje/iterum-cli/util"
	"github.com/prometheus/common/log"
)

// Catches panics, expects them to be of type error, then stores it in the pointer as recovery
func _returnErrOnPanic(perr *error) func() {
	return func() {
		if r := recover(); r != nil {
			*perr = r.(error)
		}
	}
}

// Initialize instantiates a new data repo and makes appropriate .idv folder structure
func Initialize() (err error) {
	defer _returnErrOnPanic(&err)()
	err = EnsureIDVRepo()
	if err == nil {
		return errors.New("Error: Cannot initialize idv repo. Reason: Already a repo")
	}

	// Setup folderstructure
	os.MkdirAll(localFolder, 0755)
	os.MkdirAll(remoteFolder, 0755)

	// Pulling necessary info
	log.Warnln("Still need to pull initial stuff from remote, after that delete dummy.go")
	dummyPull()

	var history VTree
	parseVTree(remoteFolder+vtreeFileName, &history)

	// Search for master branch in VTree to find hash. Then parse it into mbranch
	var mbranch Branch
	for branchHash, branchName := range history.Branches {
		if branchName == masterBranchName {
			parseBranch(branchHash.toBranchPath(false), &mbranch)
			break
		}
	}

	// Set up symlinks to locations
	var head Commit
	linkBRANCH(mbranch, false)
	parseCommit(mbranch.HEAD.toCommitPath(false), &head)
	linkHEAD(head)

	// Create current local commit being a child of HEAD and link it
	local := NewCommit(head, mbranch.Hash, "", "")
	err = local.WriteToFolder(localFolder)
	util.PanicIfErr(err, "")
	linkLOCAL(local)

	err = Stagemap{}.WriteToFile()
	util.PanicIfErr(err, "")

	return
}

// Status returns information about the currently staged files
func Status(fullPath, localPath bool) (report string, err error) {
	defer _returnErrOnPanic(&err)()
	err = EnsureLOCAL()
	util.PanicIfErr(err, "")
	var local Commit
	parseLOCAL(&local)
	if localPath {
		stagemap := parseStagemap()
		report = local.FormatDiff("{\n\t", "\n}", "< No Staged Files >", "\n\t", fullPath, stagemap)
	} else {
		report = local.FormatDiff("{\n\t", "\n}", "< No Staged Files >", "\n\t", fullPath, nil)
	}
	return
}

// Ls lists all data in the current commit
func Ls(selector *regexp.Regexp, fullPath bool) (report string, err error) {
	defer _returnErrOnPanic(&err)()
	err = EnsureLOCAL()
	util.PanicIfErr(err, "")
	var local Commit
	parseLOCAL(&local)
	report = local.FormatFiles(selector, "{\n\t", "\n}", "< Empty Data Set >", "\n\t", fullPath)
	return
}

// AddFiles stages new files to be added and existing files as Updates, expects a list of absolute file paths
func AddFiles(files []string) (adds, updates int, err error) {
	defer _returnErrOnPanic(&err)()
	err = EnsureLOCAL()
	util.PanicIfErr(err, "")
	var local Commit
	parseLOCAL(&local)
	for _, file := range files {
		if !util.FileExists(file) {
			panic(fmt.Errorf("Error: %v is either non-existent, or a directory", file))
		}
	}
	addMap, updateMap := local.AddOrUpdate(files)
	adds = len(addMap)
	updates = len(updateMap)
	stagemap := parseStagemap()
	fmt.Println(stagemap)
	err = stagemap.Update(addMap)
	util.PanicIfErr(err, "")
	err = stagemap.Update(updateMap)
	util.PanicIfErr(err, "")
	verifyAndUpdateStagemap(local, stagemap)
	writeLOCAL(local)
	return
}

// RemoveFiles stages files for removal from the dataset
// files is expected to be a list of paths on this machine
// names can be random strings that are matched against names in the commit
func RemoveFiles(files []string, names []string, unstage bool) (removals, unstages int, err error) {
	defer _returnErrOnPanic(&err)()
	err = EnsureLOCAL()
	util.PanicIfErr(err, "")
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
	defer _returnErrOnPanic(&err)()
	err = EnsureLOCAL()
	util.PanicIfErr(err, "")
	var local Commit
	parseLOCAL(&local)
	removals, unstages = local.removeWithSelector(selector, unstage)
	verifyAndUpdateStagemap(local, nil)
	writeLOCAL(local)
	return
}

// Unstage unstages adds/updates/removes of files that match the selector
func Unstage(selector *regexp.Regexp) (unstaged int, err error) {
	defer _returnErrOnPanic(&err)()
	err = EnsureLOCAL()
	util.PanicIfErr(err, "")
	var local Commit
	parseLOCAL(&local)
	unstaged = local.unstage(selector)
	verifyAndUpdateStagemap(local, nil)
	writeLOCAL(local)
	return
}

// ApplyCommit finalizes the currently staged changes and submits it to the daemon
func ApplyCommit() {

}

// BranchCommit branches off of the current commit onto a new branch
func BranchCommit() {

}

// Checkout from the current branch onto another branch/commit
func Checkout() {

}

// Download data from this repository onto this local machine
func Download(selector regexp.Regexp) {

}

// Pull the latest .vtree and if no staged changes checkout onto HEAD of branch
func Pull() {

}
