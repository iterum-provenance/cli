package idv

import (
	"fmt"

	"github.com/iterum-provenance/cli/util"
)

// This file contains functionality related to switching branches/commits
// This means branching off of a commit as well as checking out to other branches and commits

// BranchFromCommit branches off of the current commit onto a new branch
func BranchFromCommit(branchName, commitHashOrName string, isHash bool) (err error) {
	defer _returnErrOnPanic(&err)()
	EnsureByPanic(EnsureNoChanges, "")
	EnsureByPanic(EnsureTREE, "")

	var history VTree
	parseTREE(&history)
	var branchRoot Commit
	if commitHashOrName == "" { // If no commit passed, branch HEAD
		err = EnsureHEAD()
		util.PanicIfErr(err, "")
		parseHEAD(&branchRoot)
	} else { // If a commit passed, parse it based on its hash or name
		var rootHash hash
		if isHash {
			rootHash = hash(commitHashOrName)
		} else {
			rootHash, err = history.getCommitHashByName(commitHashOrName)
			util.PanicIfErr(err, "")
		}
		if !history.isExistingCommit(rootHash) {
			return fmt.Errorf("%v is not an existing commit, cannot branch of non-existent commit", rootHash)
		}
		branchRoot = pullParseCommit(rootHash)
	}

	newBranch, headCommit, err := history.branchOff(branchRoot, branchName)
	util.PanicIfErr(err, "")
	writeLOCAL(headCommit)
	err = newBranch.WriteToFolder(localFolder)
	util.PanicIfErr(err, "")
	err = history.WriteToFolder(localFolder)
	util.PanicIfErr(err, "")
	linkTREE(history, true)
	linkBRANCH(newBranch, true)
	return
}

// Checkout from the current branch onto another branch/commit
func Checkout(nameOrHash string, isCommit bool, isHash bool) (err error) {
	defer _returnErrOnPanic(&err)()
	EnsureByPanic(EnsureNoChanges, "")
	EnsureByPanic(EnsureTREE, "")

	targetHash := hash(nameOrHash)
	targetName := nameOrHash
	var history VTree
	parseTREE(&history)
	if isCommit {
		if !isHash {
			targetHash, err = history.getCommitHashByName(targetName)
			util.PanicIfErr(err, "Error: passed commit name is not part of the version tree. Make sure you have the latest version")
		}
		commit := pullParseCommit(targetHash)
		branch := pullParseBranch(commit.Branch)
		trackCommit(commit, branch, false)
	} else {
		if !isHash {
			targetHash, err = history.getBranchHashByName(targetName)
			util.PanicIfErr(err, "Error: passed branch name is not part of the version tree, are you sure have the latest version?")
		}
		branch := pullParseBranch(targetHash)
		trackBranchHead(branch, true)
	}

	return
}
