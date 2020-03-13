package idv

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
	notAlreadyARepoTest := EnsureIDVRepo()
	if notAlreadyARepoTest == nil {
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

	trackBranchHead(mbranch)
	return
}

// Status returns information about the currently staged files
func Status(fullPath, localPath bool) (report string, err error) {
	defer _returnErrOnPanic(&err)()
	EnsureByPanic(EnsureLOCALIsBranchHead, "")

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
	EnsureByPanic(EnsureLOCAL, "")
	var local Commit
	parseLOCAL(&local)
	report = local.FormatFiles(selector, "{\n\t", "\n}", "< Empty Data Set >", "\n\t", fullPath)
	return
}

// AddFiles stages new files to be added and existing files as Updates, expects a list of absolute file paths
func AddFiles(files []string) (adds, updates int, err error) {
	defer _returnErrOnPanic(&err)()
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
	stagemap := parseStagemap()
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
	defer _returnErrOnPanic(&err)()
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
	defer _returnErrOnPanic(&err)()
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
	defer _returnErrOnPanic(&err)()
	EnsureByPanic(EnsureLOCALIsBranchHead, "")

	var local Commit
	parseLOCAL(&local)
	unstaged = local.unstage(selector)
	verifyAndUpdateStagemap(local, nil)
	writeLOCAL(local)
	return
}

// BranchFromCommit branches off of the current commit onto a new branch
func BranchFromCommit(branchName, commitHashOrName string, isHash bool) (err error) {
	defer _returnErrOnPanic(&err)()
	EnsureByPanic(EnsureNoChanges, "")
	log.Warn("Should ensure latest vtree file")

	var history VTree
	parseVTree(vtreeFilePath, &history)
	var branchRoot Commit
	if commitHashOrName == "" {
		err = EnsureHEAD()
		util.PanicIfErr(err, "")
		parseHEAD(&branchRoot)
	} else {
		var rootHash hash
		if isHash {
			rootHash = hash(commitHashOrName)
		} else {
			rootHash, err = history.getCommitHashByName(commitHashOrName)
			util.PanicIfErr(err, "")
		}
		if !history.isExistingCommit(rootHash) {
			return fmt.Errorf("%v is not an existing commit, cannot branch of non-existent commit", rootHash)
		}
		rootPath := remoteFolder + rootHash.String() + commitFileExt
		if util.FileExists(rootPath) {
			parseCommit(rootPath, &branchRoot)
		} else {
			log.Warn(fmt.Sprintf("Should pull %v%v file", rootHash, commitFileExt))
			return errors.New("Error: cannot pull files yet")
		}
	}

	newBranch, headCommit, err := history.branchOff(branchRoot, branchName)
	util.PanicIfErr(err, "")
	writeLOCAL(headCommit)
	err = newBranch.WriteToFolder(localFolder)
	util.PanicIfErr(err, "")
	err = history.WriteToFolder(localFolder)
	util.PanicIfErr(err, "")
	linkBRANCH(newBranch, true)
	return
}

// Checkout from the current branch onto another branch/commit
func Checkout(nameOrHash string, isCommit bool, isHash bool) (err error) {
	defer _returnErrOnPanic(&err)()
	EnsureByPanic(EnsureNoChanges, "")

	targetHash := hash(nameOrHash)
	targetName := nameOrHash
	var history VTree
	parseVTree(vtreeFilePath, &history)
	if isCommit {
		if !isHash {
			targetHash, err = history.getCommitHashByName(targetName)
			util.PanicIfErr(err, "Error: passed commit name is not part of the version tree. Make sure you have the latest version")
		}
		commit := pullParseCommit(targetHash)
		branch := pullParseBranch(commit.Branch)
		trackCommit(commit, branch)
	} else {
		if !isHash {
			targetHash, err = history.getBranchHashByName(targetName)
			util.PanicIfErr(err, "Error: passed branch name is not part of the version tree, are you sure have the latest version?")
		}
		branch := pullParseBranch(targetHash)
		trackBranchHead(branch)
	}

	return
}

// ApplyCommit finalizes the currently staged changes and submits it to the daemon
func ApplyCommit() (err error) {
	defer _returnErrOnPanic(&err)()
	EnsureByPanic(EnsureLOCALIsBranchHead, "")
	log.Warn("Should ensure latest vtree file")

	var local Commit
	var branch Branch
	var vtree VTree
	parseLOCAL(&local)
	// Uncomment this once needed for Daemon integration
	// stagemap := parseStagemap()

	// parse potential new VTree (in case we branched)
	// this can all go once we integrate with DAEMON!
	gotTree := false
	files, err := ioutil.ReadDir(localFolder)
	util.PanicIfErr(err, "")
	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if ext == vtreeFileExt {
			parseVTree(file.Name(), &vtree)
			gotTree = true
		}
	}
	// Here stops what we ca throw away
	if !gotTree { // this if can stil go
		parseVTree(vtreeFilePath, &vtree)
	}

	parseBRANCH(&branch)
	err = local.applyStaged()
	util.PanicIfErr(err, "")
	branch.HEAD = local.Hash

	log.Warn("TODO: Create multipart form of all data that needs to be send")
	log.Warn("TODO: pass all (necessary) data to the Daemon")
	log.Warn("TODO: accept response of updated .vtree and .branch file (dummy this first)")
	// Clean up/Remove all files from .idv/local
	// trackBranch(returnedBranch)
	return
}

// Download data from this repository onto this local machine
func Download(selector regexp.Regexp) {

}

// Pull the latest .vtree and if no staged changes checkout onto HEAD of branch
func Pull() {

}
