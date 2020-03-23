package idv

import (
	"errors"
	"os"

	"github.com/Mantsje/iterum-cli/idv/ctl"
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
	if err != nil {
		return err
	}
	history.WriteToFolder(remoteFolder)
	linkTREE(history, false)

	mbranch := NewBranch(masterBranchName)
	commit := NewRootCommit(mbranch.Hash)
	mbranch.HEAD = commit.Hash
	commit.WriteToFolder(remoteFolder)
	mbranch.WriteToFolder(remoteFolder)

	newHistory, err := pushBranchedCommit(ctl.Name, mbranch, commit, Stagemap{})
	if err != nil {
		return err
	}
	writeTREE(newHistory)

	trackCommit(commit, mbranch, true)

	return
}
