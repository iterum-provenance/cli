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
// All of these functions panic in case of an error panici(err)
// The caller is responsible for deferring a catch

// _symlink symlinks a file and panics on fail
func _symlink(src string, target string) {
	if !util.FileExists(src) {
		panic(errors.New("Error: cannot symlink non persisted-to-disk structure"))
	}
	err := os.Symlink("../"+src, target)
	util.PanicIfErr(err, "")
}

// linkHEAD symlinks the HEAD of the repo and panics on fail
func linkHEAD(remoteCommit Commit) {
	path := remoteCommit.ToFilePath(false)
	_symlink(path, HEAD)
}

// linkLOCAL symlinks the LOCAL of the repo and panics on fail
func linkLOCAL(localCommit Commit) {
	path := localCommit.ToFilePath(true)
	_symlink(path, LOCAL)
}

// linkBRANCH symlinks the BRANCH of the repo and panics on fail
func linkBRANCH(branch Branch, isLocal bool) {
	path := branch.ToFilePath(isLocal)
	_symlink(path, BRANCH)
}

// _parse calls ParseFromFile on passed interface.
// On success, p contains data, panics on fail
func _parse(path string, p util.Parseable) {
	err := p.ParseFromFile(path)
	util.PanicIfErr(err, "")
}

// parseBranch parses a branch into pointer, else panics
func parseBranch(path string, branch *Branch) {
	_parse(path, branch)
}

// parseCommit parses a commit into pointer, else panics
func parseCommit(path string, commit *Commit) {
	_parse(path, commit)
}

// parseVTree parses a vtree into pointer, else panics
func parseVTree(path string, history *VTree) {
	_parse(path, history)
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

func parseStagemap() (stagemap Stagemap) {
	_parse(stagedFilePath, &stagemap)
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
		stagemap = parseStagemap()
	}
	err := stagemap.verifyAndSyncWithCommit(commit)
	util.PanicIfErr(err, "")
	err = stagemap.WriteToFile()
	util.PanicIfErr(err, "")
}

// initLOCAL initializes LOCAL by inheriting from the currently tracked HEAD
func initLOCAL() {
	// Ensure dependencies for this
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

	initLOCAL()
}

// trackBranchHead takes a branch and sets HEAD and BRANCH to this branch ands its head
// it also initializes LOCAL with a commit inheriting from HEAD
func trackBranchHead(branch Branch) {
	err := EnsureIDVRepo()
	util.PanicIfErr(err, "")

	// Parse the head of this branch
	var head Commit
	parseCommit(branch.HEAD.toCommitPath(false), &head)

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

func pullParseCommit(h hash) (commit Commit) {
	commit = Commit{}
	err := _pullAndParse(remoteFolder+h.String()+commitFileExt, &commit)
	util.PanicIfErr(err, "")
	return
}

func pullParseBranch(h hash) (branch Branch) {
	branch = Branch{}
	err := _pullAndParse(remoteFolder+h.String()+branchFileExt, &branch)
	util.PanicIfErr(err, "")
	return
}
