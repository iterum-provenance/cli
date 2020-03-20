package idv

import (
	"fmt"

	"github.com/Mantsje/iterum-cli/idv/ctl"
	"github.com/Mantsje/iterum-cli/util"
	"github.com/prometheus/common/log"
)

// This file contains code related to pushing and pulling to and from the remote data storage

// ApplyCommit finalizes the currently staged changes and submits it to the daemon
func ApplyCommit(name, description string) (err error) {
	defer _returnErrOnPanic(&err)()
	EnsureByPanic(EnsureLOCALIsBranchHead, "")
	EnsureByPanic(EnsureChanges, "")
	EnsureByPanic(EnsureConfig, "")
	log.Warn("TODO: Should ensure latest vtree file")

	var ctl ctl.DataCTL
	ctl.ParseFromFile(configPath) // No error is ensured, so no need to catch it

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

	local.WriteToFolder(remoteFolder)
	branch.WriteToFolder(remoteFolder)  // should go in case of not pushing a branch
	history.WriteToFolder(remoteFolder) // should go afterwards

	linkTREE(history, false)
	trackBranchHead(branch, true)

	log.Warn("TODO: Create multipart form of all data that needs to be send")
	log.Warn("TODO: pass all (necessary) data to the Daemon")
	log.Warn("TODO: accept response of updated .vtree and .branch file")
	return
}

// PullVTree pulls the latest .vtree for the dataset of this repo
func PullVTree() (err error) {
	defer _returnErrOnPanic(&err)()
	EnsureByPanic(EnsureIDVRepo, "")
	EnsureByPanic(EnsureNoBranchOffs, "")
	EnsureByPanic(EnsureConfig, "")
	var ctl ctl.DataCTL
	ctl.ParseFromFile(configPath)
	history, err := getVTree(ctl.Name)
	writeTREE(history)
	fmt.Println(history)
	return
}
