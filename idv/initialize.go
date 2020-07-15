package idv

import (
	"errors"
	"os"

	"github.com/iterum-provenance/cli/idv/ctl"
	"github.com/iterum-provenance/iterum-go/util"
)

// Initialize instantiates a new data repo and makes appropriate .idv folder structure
func Initialize() (err error) {
	defer util.ReturnErrOnPanic(&err)()
	notAlreadyARepoTest := EnsureIDVRepo()
	if notAlreadyARepoTest == nil {
		return errors.New("Error: Cannot initialize idv repo. Reason: Already a repo")
	}

	// Setup folderstructure
	os.MkdirAll(localFolder, 0755)
	os.MkdirAll(remoteFolder, 0755)

	return
}

// Setup sets up the necessary stuff at the Daemon and locally links all necessary symlinks
func Setup() (err error) {
	defer util.ReturnErrOnPanic(&err)()
	EnsureByPanic(EnsureIDVRepo, "")
	EnsureByPanic(EnsureConfig, "")

	errNotSetup := EnsureSetup()
	if errNotSetup == nil { // Meaning this repo has already been setup
		return errors.New("Error: Repo already set up")
	}

	var ctl ctl.DataCTL
	parseConfig(configPath, &ctl)

	errPosting := postDataset(ctl)
	if errPosting != nil && errPosting != errConflictingDataset {
		return errPosting
	}

	history, err := getVTree(ctl)
	util.PanicIfErr(err, "")
	history.WriteToFolder(remoteFolder)
	linkTREE(history, false)

	mbranchHash, err := history.getBranchHashByName(masterBranchName)
	util.PanicIfErr(err, "")

	mbranch, err := getBranch(ctl, mbranchHash)
	util.PanicIfErr(err, "")
	mbranch.WriteToFolder(remoteFolder)
	commit, err := getCommit(ctl, mbranch.HEAD)
	util.PanicIfErr(err, "")
	commit.WriteToFolder(remoteFolder)

	trackCommit(commit, mbranch, true)
	return
}

// Apply sets up the necessary stuff at the Daemon using idv-config.yaml
func Apply() (err error) {
	defer util.ReturnErrOnPanic(&err)()
	EnsureByPanic(EnsureSetup, "")
	EnsureByPanic(EnsureConfig, "")

	var ctl ctl.DataCTL
	parseConfig(configPath, &ctl)

	errPosting := postDataset(ctl)
	if errPosting != nil && errPosting != errConflictingDataset {
		return errPosting
	}
	return
}
