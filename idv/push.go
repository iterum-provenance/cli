package idv

import (
	"github.com/iterum-provenance/cli/idv/ctl"
	"github.com/iterum-provenance/cli/util"
)

// This file contains code related to pushing and pulling to and from the remote data storage

// ApplyCommit finalizes the currently staged changes and submits it to the daemon
func ApplyCommit(name, description string) (err error) {
	defer _returnErrOnPanic(&err)()
	EnsureByPanic(EnsureSetup, "")
	EnsureByPanic(EnsureLOCALIsBranchHead, "")
	EnsureByPanic(EnsureChanges, "")
	EnsureByPanic(EnsureConfig, "")

	var ctl ctl.DataCTL
	ctl.ParseFromFile(configPath) // No error is ensured, so no need to catch it

	var local Commit
	parseLOCAL(&local)
	var branch Branch
	parseBRANCH(&branch)
	var stagemap Stagemap
	parseStagemap(&stagemap)

	err = stagemap.verifyAndSyncWithCommit(local)
	util.PanicIfErr(err, "")

	err = local.applyStaged()
	util.PanicIfErr(err, "")
	local.Name = name
	local.Description = description

	local.writeToFile(tempCommitPath)

	var updatedHistory VTree
	var updatedBranch Branch
	if isBranched() {
		updatedBranch, updatedHistory, err = postBranchedCommit(ctl.Name, branch, local, stagemap)
	} else {
		updatedBranch, updatedHistory, err = postCommit(ctl.Name, local, stagemap)
	}
	util.PanicIfErr(err, "")

	local.WriteToFolder(remoteFolder)
	updatedBranch.WriteToFolder(remoteFolder)
	updatedHistory.WriteToFolder(remoteFolder)

	linkTREE(updatedHistory, false)
	trackBranchHead(updatedBranch, true)

	return
}
