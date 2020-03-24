package idv

import (
	"os"

	"github.com/Mantsje/iterum-cli/idv/ctl"
	"github.com/Mantsje/iterum-cli/util"
)

func attemptMergeToHead(ctl ctl.DataCTL, local Commit, remoteBranch Branch) (remoteHead, newLocal Commit, err error) {
	defer _returnErrOnPanic(&err)()

	remoteHead, err = getCommit(remoteBranch.HEAD, ctl.Name)
	util.PanicIfErr(err, "")

	newLocal = NewCommit(remoteHead, remoteBranch.Hash, "", "")
	removeUpdatedRemovals := true // REMOVE files that were staged for REMOVE but UPDATED in unseen commits
	addRemovedUpdates := true     // ADD files that were staged for UPDATE but REMOVED during unseen commits
	updatePresentAdds := true     // UPDATE files that were staged for ADD but already ADDED during unseen commits
	newLocal.autoMerge(local, removeUpdatedRemovals, addRemovedUpdates, updatePresentAdds)
	return
}

// handlePullWhilstCheckedOut takes care of updating the locally stored version tree whilst
// the user is checked out to some commit. This means just overwriting the current tree, as it can only
// be more extensive than it was before, therefore not breaking anything.
func handlePullWhilstCheckedOut(ctl ctl.DataCTL, remoteHistory VTree, localCommit Commit) (err error) {
	defer _returnErrOnPanic(&err)()
	EnsureByPanic(EnsureNoChanges, "") // just in case

	// It was safe to pull
	writeTREE(remoteHistory)
	return
}

func handlePullWhilstUnbranched(ctl ctl.DataCTL, remoteHistory VTree, localCommit Commit) (err error) {
	defer _returnErrOnPanic(&err)()

	remoteBranch, err := getBranch(localCommit.Branch, ctl.Name)
	util.PanicIfErr(err, "")
	if remoteBranch.HEAD != localCommit.Parent { // If we are behind the remote
		// It might not be safe to pull
		// Create a new commit that is ahead of the remote, persisting old changes
		remoteHead, newLocal, err := attemptMergeToHead(ctl, localCommit, remoteBranch)
		util.PanicIfErr(err, "")

		remoteBranch.WriteToFolder(remoteFolder)
		remoteHead.WriteToFolder(remoteFolder)
		newLocal.WriteToFolder(localFolder)
		linkBRANCH(remoteBranch, false)
		linkHEAD(remoteHead)
		linkLOCAL(newLocal)

		err = os.Remove(localCommit.ToFilePath(true)) // because we have a successful new local
		util.PanicIfErr(err, "")

	}
	// Safe to pull and write because we are ahead of the remote branch now
	writeTREE(remoteHistory)

	return
}

// handlePullWhilstBranched takes care of pulling the version tree and writing it to disk
// whilst the user is working on a new local branch not yet pushed to the remote
// Essentially it becomes save to pull the file, however because currently the CLI keeps a
// local copy of vtree that is updated when we branch, we need to update this new remoteVTree as well
func handlePullWhilstBranched(ctl ctl.DataCTL, remoteHistory VTree, localCommit Commit) (err error) {
	defer _returnErrOnPanic(&err)()
	EnsureByPanic(EnsureBRANCH, "")

	// update local history with any remote changes, then add the branch off information too
	var branch Branch
	parseBRANCH(&branch)
	newLocalHistory := remoteHistory.duplicate()
	errBranch := newLocalHistory.addBranch(branch)
	errCommit := newLocalHistory.addCommit(localCommit)
	util.PanicIfErr(util.ReturnFirstErr(errBranch, errCommit), "")

	writeTREE(newLocalHistory) // in this case TREE is in localFolder, other cases it in remoteFolder
	remoteHistory.WriteToFolder(remoteFolder)

	return
}

// Pull gets the latest vtree from the daemon and attempts to resolve any conflicts arising locally from this
func Pull() (err error) {
	defer _returnErrOnPanic(&err)()
	EnsureByPanic(EnsureSetup, "")
	EnsureByPanic(EnsureLOCAL, "")
	EnsureByPanic(EnsureConfig, "")

	// Gather required parts
	var ctl ctl.DataCTL
	ctl.ParseFromFile(configPath)
	remoteHistory, err := getVTree(ctl.Name)
	util.PanicIfErr(err, "")
	var localCommit Commit
	parseLOCAL(&localCommit)

	if err = EnsureLatestCommit(); err == errNotHEAD {
		// If we're checked out to some other commit
		err = handlePullWhilstCheckedOut(ctl, remoteHistory, localCommit)
	} else if err != nil {
		util.PanicIfErr(err, "")
	}

	if isBranched() {
		// We are branched off onto a new branch locally, safe to pull
		err = handlePullWhilstBranched(ctl, remoteHistory, localCommit)
	} else {
		err = handlePullWhilstUnbranched(ctl, remoteHistory, localCommit)
	}
	util.PanicIfErr(err, "")

	return
}
