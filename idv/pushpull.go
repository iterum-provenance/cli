package idv

import (
	"errors"
	"os"

	"github.com/Mantsje/iterum-cli/util"
	"github.com/prometheus/common/log"
)

// This file contains code related to pushing and pulling to and from the remote data storage

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
	linkTREE(history, false)

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

// ApplyCommit finalizes the currently staged changes and submits it to the daemon
func ApplyCommit(name, description string) (err error) {
	defer _returnErrOnPanic(&err)()
	EnsureByPanic(EnsureLOCALIsBranchHead, "")
	log.Warn("TODO: Should ensure latest vtree file")

	var local Commit
	parseLOCAL(&local)
	var branch Branch
	parseBRANCH(&branch)
	var history VTree
	parseTREE(&history)
	var stagemap Stagemap
	parseStagemap(&stagemap)

	err = local.applyStaged()
	util.PanicIfErr(err, "")
	local.Name = name
	local.Description = description

	// Both of the following 2 statements should be performed at daemon (except when branched maybe)
	branch.HEAD = local.Hash
	history.add(local)
	log.Warn("TODO: Create multipart form of all data that needs to be send")
	log.Warn("TODO: pass all (necessary) data to the Daemon")
	log.Warn("TODO: accept response of updated .vtree and .branch file (dummy this first)")

	local.WriteToFolder(remoteFolder)
	branch.WriteToFolder(remoteFolder)
	history.WriteToFolder(remoteFolder)

	linkTREE(history, false)
	trackBranchHead(branch)
	return
}

// Pull the latest .vtree and if no staged changes checkout onto HEAD of branch
func Pull() {

}
