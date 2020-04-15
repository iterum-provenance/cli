package idv

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/iterum-provenance/cli/idv/ctl"
	"github.com/iterum-provenance/cli/util"
)

var (
	errNotHEAD = errors.New("Error: LOCAL is not direct child of BRANCH.HEAD")
)

// Catches panics, expects them to be of type error, then stores it in the pointer as recovery
func _returnErrOnPanic(perr *error) func() {
	return func() {
		if r := recover(); r != nil {
			*perr = r.(error)
		}
	}
}

// EnsureByPanic takes one of the other Ensure functions
// and panics in case of an error rather than returning it
func EnsureByPanic(ensurerFunc func() error, customMsg string) {
	err := ensurerFunc()
	util.PanicIfErr(err, customMsg)
}

// EnsureIDVRepo makes sure that we are in an idv repository
func EnsureIDVRepo() error {
	if !util.IsFileOrDir(IDVFolder) {
		return errors.New("Error: Not an idv repository")
	}
	return nil
}

func _ensurePath(path string) error {
	errs := []error{}
	errs = append(errs, EnsureIDVRepo())
	if !util.IsFileOrDir(path) {
		errs = append(errs, fmt.Errorf("Error: Either %v does not exist, or does not point to a file", path))
	}
	return util.ReturnFirstErr(errs...)
}

// EnsureLOCAL makes sure that LOCAL points to a file
func EnsureLOCAL() error {
	return _ensurePath(LOCAL)
}

// EnsureHEAD makes sure that HEAD points to a file
func EnsureHEAD() error {
	return _ensurePath(HEAD)
}

// EnsureTREE makes sure that TREE points to a file
func EnsureTREE() error {
	return _ensurePath(TREE)
}

// EnsureBRANCH makes sure that BRANCH points to a file
func EnsureBRANCH() error {
	return _ensurePath(BRANCH)
}

// EnsureNoBranchOffs makes sure that there are no uncommitted branch and vtree files in .idv/local
func EnsureNoBranchOffs() error {
	files, err := ioutil.ReadDir(localFolder)
	if err != nil {
		return err
	}

	for _, file := range files {
		ext := filepath.Ext(file.Name())
		// If there are branch or vtree files in local, that means new branches are created and uncommitted
		if ext == branchFileExt || ext == vtreeFileExt {
			return errors.New("Error: there are uncommitted changes/branches in this repository")
		}
	}
	return nil
}

// EnsureLatestCommit makes sure that the current LOCAL is ahead of the BRANCH.HEAD
func EnsureLatestCommit() error {
	errLocal := EnsureLOCAL()
	errBranch := EnsureBRANCH()
	var local Commit
	errLocalParse := local.ParseFromFile(LOCAL)
	var branch Branch
	errBranchParse := branch.ParseFromFile(BRANCH)
	e := util.ReturnFirstErr(errLocal, errBranch, errLocalParse, errBranchParse)
	if e != nil {
		return e
	}

	if local.Parent != branch.HEAD {
		return errNotHEAD
	}
	return nil
}

// EnsureNoStaged makes sure that there are no uncommitted staged changes in LOCAL
func EnsureNoStaged() error {
	isLatest := EnsureLatestCommit()
	if isLatest != nil { // If it is not the latest commit, this function should not error
		return nil
	}
	err := EnsureLOCAL()
	if err != nil {
		return err
	}

	var local Commit
	err = local.ParseFromFile(LOCAL)
	if err == nil && local.containsChanges() {
		return errors.New("Error: Current LOCAL commit contains staged changes that are not yet committed")
	}
	return err
}

// EnsureNoChanges errors if there are uncommitted changes staged
func EnsureNoChanges() error {
	errBranchoffs := EnsureNoBranchOffs()
	errStaged := EnsureNoStaged()
	return util.ReturnFirstErr(errBranchoffs, errStaged)
}

// EnsureChanges errors if the current LOCAL commit has no staged changes
func EnsureChanges() error {
	err := EnsureLOCAL()
	if err != nil {
		return err
	}

	var local Commit
	err = local.ParseFromFile(LOCAL)
	if err == nil && !local.containsChanges() {
		return errors.New("Error: Current LOCAL commit contains no staged changes")
	}

	return err
}

// EnsureLOCALIsBranchHead checks whether LOCAL is ahead of the current branch's branch.HEAD
func EnsureLOCALIsBranchHead() error {
	errLocal := EnsureLOCAL()
	errBranch := EnsureBRANCH()
	var local Commit
	errLocalParse := local.ParseFromFile(LOCAL)
	var branch Branch
	errBranchParse := branch.ParseFromFile(BRANCH)
	e := util.ReturnFirstErr(errLocal, errBranch, errLocalParse, errBranchParse)
	if e != nil {
		return e
	}

	// if the parent of LOCAL is current branch.HEAD (meaning LOCAL is latest possible commit)
	if local.Parent != branch.HEAD {
		return errNotHEAD
	}
	return nil
}

// EnsureConfig checks whether there is a file called idv-config.yaml that can be parsed as config
func EnsureConfig() error {
	err := EnsureIDVRepo()
	if err != nil {
		return err
	}
	if !util.FileExists(configPath) {
		return errors.New("Error: idv-config.yaml not found")
	}
	var ctl ctl.DataCTL
	return ctl.ParseFromFile(configPath)
}

// EnsureSetup checks whether the repo is (correctly) setup by checking all the supposed symlinks
func EnsureSetup() (err error) {
	err = EnsureIDVRepo()
	if err != nil {
		return
	}

	_, errBRANCH := os.Lstat(BRANCH)
	_, errLOCAL := os.Lstat(LOCAL)
	_, errTREE := os.Lstat(TREE)
	_, errHEAD := os.Lstat(HEAD)
	err = util.ReturnFirstErr(errBRANCH, errLOCAL, errTREE, errHEAD)
	if err != nil {
		return
	}

	_, errBRANCH = os.Stat(BRANCH)
	_, errLOCAL = os.Stat(LOCAL)
	_, errTREE = os.Stat(TREE)
	_, errHEAD = os.Stat(HEAD)
	err = util.ReturnFirstErr(errBRANCH, errLOCAL, errTREE, errHEAD)

	return
}
