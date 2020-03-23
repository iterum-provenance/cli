package idv

import (
	"errors"
	"os"

	"github.com/Mantsje/iterum-cli/idv/ctl"
	"github.com/Mantsje/iterum-cli/util"
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

	var ctl ctl.DataCTL
	parseConfig(configPath, &ctl)

	errPosting := postDataset(ctl)
	if errPosting != nil && errPosting != errConflictingDataset {
		return errPosting
	}

	history, err := getVTree(ctl.Name)
	util.PanicIfErr(err, "")
	history.WriteToFolder(remoteFolder)
	linkTREE(history, false)

	mbranchHash, err := history.getBranchHashByName(masterBranchName)
	util.PanicIfErr(err, "")

	mbranch, err := getBranch(mbranchHash, ctl.Name)
	util.PanicIfErr(err, "")
	mbranch.WriteToFolder(remoteFolder)
	commit, err := getCommit(mbranch.HEAD, ctl.Name)
	util.PanicIfErr(err, "")
	commit.WriteToFolder(remoteFolder)

	trackCommit(commit, mbranch, true)
	return
}
