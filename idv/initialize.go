package idv

import (
	"errors"
	"os"

	"github.com/prometheus/common/log"
)

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
