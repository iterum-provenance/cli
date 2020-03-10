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

	return
}

// Status returns information about the currently staged files
func Status() (report string, err error) {
	defer _returnErrOnPanic(&err)()
	err = EnsureLOCAL()
	util.PanicIfErr(err, "")
	var local Commit
	parseLOCAL(&local)
	report = local.FormatDiff("{", "}", "< No Staged Files >", "\n\t")
	return
}

// Add stages new files to be added and existing files as Updates
func Add(files []string) (err error) {
	defer _returnErrOnPanic(&err)()
	err = EnsureLOCAL()
	util.PanicIfErr(err, "")
	var local Commit
	parseLOCAL(&local)
	for _, file := range files {
		if !util.IsFolderOrDir(file) {
			panic(fmt.Errorf("Error: %v is not a file or directory", file))
		}
		info, _ := os.Stat(file)
		fmt.Println(info.Name())
	}
	return
}

// Remove stages files for removal from the dataset
func Remove(files []string) {

}

// Ls lists all data in the current commit
func Ls(selector regexp.Regexp) {

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
