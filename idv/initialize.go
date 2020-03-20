package idv

import (
	"errors"
	"fmt"
	"os"

	"github.com/Mantsje/iterum-cli/idv/ctl"
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
	parseVTree(vtreeFilePath, &history)
	linkTREE(history, false)

	// Search for master branch in VTree to find hash. Then parse it into mbranch
	var mbranch Branch
	for branchHash, branchName := range history.Branches {
		if branchName == masterBranchName {
			parseBranch(branchHash.toBranchPath(false), &mbranch)
			break
		}
	}

	trackBranchHead(mbranch, true)
	return
}

// Setup sets up the necessary stuff at the Daemon and locally links all necessary symlinks
func Setup() (err error) {
	defer _returnErrOnPanic(&err)()
	EnsureByPanic(EnsureIDVRepo, "")
	EnsureByPanic(EnsureConfig, "")

	errNotSetup := EnsureSetup()
	if errNotSetup == nil { // Meaning this repo has already been setup
		return errors.New("Error: Repo already set up")
	}
	fmt.Println(errNotSetup)

	var ctl ctl.DataCTL
	parseConfig(configPath, &ctl)

	errPosting := postDataset(ctl)
	if errPosting != nil && errPosting != errConflictingDataset {
		return errPosting
	}

	return
}
