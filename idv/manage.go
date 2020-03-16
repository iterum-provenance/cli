package idv

import (
	"errors"
	"fmt"
	"os"

	"github.com/Mantsje/iterum-cli/util"
	"github.com/prometheus/common/log"
)

// Contains functionality used for managing an idv repository
// Non of these are exported and therefore for internal use only
// (Nearly) all of these functions panic in case of an error panic(err)
// The caller is responsible for deferring a function that catches the panic
// These are all helper functions for the actual execution of the versioning commands

// _symlink symlinks a file and panics on fail
func _symlink(src string, target string) (err error) {
	if !util.FileExists(src) {
		return errors.New("Error: Cannot symlink non persisted-to-disk structure")
	}
	if _, err := os.Lstat(target); err == nil {
		if err := os.Remove(target); err != nil {
			return fmt.Errorf("Error: Failed to unlink: %+v", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("Error: Symlink checking failed: %+v", err)
	}
	err = os.Symlink("../"+src, target)
	return err
}

// linkHEAD symlinks the HEAD of the repo and panics on fail
func linkHEAD(remoteCommit Commit) {
	path := remoteCommit.ToFilePath(false)
	err := _symlink(path, HEAD)
	util.PanicIfErr(err, "")
}

// linkLOCAL symlinks the LOCAL of the repo and panics on fail
func linkLOCAL(localCommit Commit) {
	path := localCommit.ToFilePath(true)
	err := _symlink(path, LOCAL)
	util.PanicIfErr(err, "")
}

// linkBRANCH symlinks the BRANCH of the repo and panics on fail
func linkBRANCH(branch Branch, isLocal bool) {
	path := branch.ToFilePath(isLocal)
	err := _symlink(path, BRANCH)
	util.PanicIfErr(err, "")
}

// linkTREE symlinks the TREE of the repo and panics on fail
func linkTREE(vtree VTree, isLocal bool) {
	path := vtree.ToFilePath(isLocal)
	err := _symlink(path, TREE)
	util.PanicIfErr(err, "")
}

// _parse calls ParseFromFile on passed interface.
// On success, p contains data, panics on fail
func _parse(path string, p util.Parseable) (err error) {
	return p.ParseFromFile(path)
}

// parseBranch parses a branch into pointer, else panics
func parseBranch(path string, branch *Branch) {
	err := _parse(path, branch)
	util.PanicIfErr(err, "")
}

// parseCommit parses a commit into pointer, else panics
func parseCommit(path string, commit *Commit) {
	err := _parse(path, commit)
	util.PanicIfErr(err, "")
}

// parseVTree parses a vtree into pointer, else panics
func parseVTree(path string, history *VTree) {
	err := _parse(path, history)
	util.PanicIfErr(err, "")
}

// parseHEAD takes a commit pointer and reads in the HEAD into this
// panics on fail, else commit will be filled with HEAD
func parseHEAD(commit *Commit) {
	parseCommit(HEAD, commit)
}

// parseLOCAL takes a commit pointer and reads in LOCAL into this
// panics on fail, else commit will be filled with LOCAL
func parseLOCAL(commit *Commit) {
	parseCommit(LOCAL, commit)
}

// parseBRANCH takes a branch pointer and reads in BRANCH into this
// panics on fail, else branch will be filled with BRANCH
func parseBRANCH(branch *Branch) {
	parseBranch(BRANCH, branch)
}

func parseTREE(vtree *VTree) {
	parseVTree(TREE, vtree)
}

func parseStagemap(stagemap *Stagemap) {
	err := _parse(stagedFilePath, stagemap)
	util.PanicIfErr(err, "")
	return
}

// writeLOCAL takes a commit and writes it to wherever LOCAL points
func writeLOCAL(commit Commit) {
	err := commit.writeToFile(LOCAL)
	util.PanicIfErr(err, "")
}

// verifyAndUpdateStagemap takes a commit and (possibly) a stagemap and validates
// the one with the other. If no map is passed it reads in the default one.
func verifyAndUpdateStagemap(commit Commit, stagemap Stagemap) {
	if stagemap == nil {
		parseStagemap(&stagemap)
	}
	err := stagemap.verifyAndSyncWithCommit(commit)
	util.PanicIfErr(err, "")
	err = stagemap.WriteToFile()
	util.PanicIfErr(err, "")
}

// initLocalFolder initializes localFolder by inheriting from the currently tracked HEAD
func initLocalFolder() {
	// Ensure dependencies for this
	clearLocalFolder()
	err := EnsureHEAD()
	util.PanicIfErr(err, "")
	err = EnsureBRANCH()
	util.PanicIfErr(err, "")

	var parent Commit
	var branch Branch
	parseHEAD(&parent)
	parseBRANCH(&branch)

	// Create current local commit being a child of HEAD and link it
	local := NewCommit(parent, branch.Hash, "", "")
	err = local.WriteToFolder(localFolder)
	util.PanicIfErr(err, "")
	linkLOCAL(local)

	// Setup empty Stagemap for this local commit
	err = Stagemap{}.WriteToFile()
	util.PanicIfErr(err, "")

}

// trackCommit does the same as trackBranchHead, but it reasons from a commit and tracks that
// Is used to track from a random commit in a tree rather than only branch.HEAD
func trackCommit(commit Commit, branch Branch) {
	err := EnsureIDVRepo()
	util.PanicIfErr(err, "")

	if commit.Branch != branch.Hash {
		panic(fmt.Errorf("Error: Commit %v with branch %v, is not part of passed branch %v", commit.Hash, commit.Branch, branch.Hash))
	}

	// Set up symlinks to locations
	linkBRANCH(branch, false)
	linkHEAD(commit)

	initLocalFolder()
}

// trackBranchHead takes a branch and sets HEAD and BRANCH to this branch ands its head
// it also initializes LOCAL with a commit inheriting from HEAD
func trackBranchHead(branch Branch) {
	err := EnsureIDVRepo()
	util.PanicIfErr(err, "")

	// Parse the head of this branch
	var head Commit = pullParseCommit(branch.HEAD)

	trackCommit(head, branch)
}

// _pullAndParse checks for existance of an idv file. If it exists it parses it, else it first pulls it
func _pullAndParse(path string, p util.Parseable) (err error) {
	if !util.FileExists(path) { // If the file does not exist yet locally, pull them
		log.Warn(fmt.Sprintf("Should pull %v file", path))
		return errors.New("Error: cannot pull files yet")
	}
	_parse(path, p)
	return
}

// pullParseBranch first pulls and then parses the given commit associated with the hash
func pullParseCommit(h hash) (commit Commit) {
	commit = Commit{}
	err := _pullAndParse(remoteFolder+h.String()+commitFileExt, &commit)
	util.PanicIfErr(err, "")
	return
}

// pullParseBranch first pulls and then parses the given branch associated with the hash
func pullParseBranch(h hash) (branch Branch) {
	err := _pullAndParse(remoteFolder+h.String()+branchFileExt, &branch)
	util.PanicIfErr(err, "")
	return
}

// clearLocalFolder removes all files in the idv/local folder such that we can track a new commit
func clearLocalFolder() {
	os.RemoveAll(localFolder)
	os.MkdirAll(localFolder, 0755)
}
