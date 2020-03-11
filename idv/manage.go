package idv

import (
	"errors"
	"os"

	"github.com/Mantsje/iterum-cli/util"
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

// writeLOCAL takes a commit and writes it to wherever LOCAL points
func writeLOCAL(commit Commit) {
	err := commit.writeToFile(LOCAL)
	util.PanicIfErr(err, "")
}
