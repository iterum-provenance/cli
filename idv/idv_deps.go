package idv

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/Mantsje/iterum-cli/util"
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
		errs = append(errs, fmt.Errorf("Error: Either .idv/%v does not exist, or does not point to a file", path))
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

// EnsureNoStaged makes sure that there are no uncommitted staged changes in LOCAL
func EnsureNoStaged() error {
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

// EnsureLOCALIsBranchHead checks whether LOCAL is the current branch's branch HEAD
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
	// if current LOCAL is not branch.HEAD (in case we already branched)
	// nor the parent of LOCAL is current branch.HEAD (meaning LOCAL is latest possible commit))
	if local.Hash != branch.HEAD && local.Parent != branch.HEAD {
		return errors.New("Error: Cannot create child commit of commits from the past unless you branch off of them")
	}
	return nil
}
