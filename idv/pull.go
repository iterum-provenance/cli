package idv

import (
	"fmt"

	"github.com/iterum-provenance/cli/idv/ctl"
	"github.com/iterum-provenance/cli/util"
)

func attemptMergeToHead(ctl ctl.DataCTL, local Commit, remoteBranch Branch) (remoteHead, newLocal Commit, err error) {
	defer util.ReturnErrOnPanic(&err)()

	remoteHead, err = getCommit(ctl, remoteBranch.HEAD)
	util.PanicIfErr(err, "")

	newLocal = NewCommit(remoteHead, remoteBranch.Hash, "", "")
	newLocal.autoMerge(local)

	return remoteHead, newLocal, err
}

// handlePullWhilstCheckedOut takes care of updating the locally stored version tree whilst
// the user is checked out to some commit. This means just overwriting the current tree, as it can only
// be more extensive than it was before, therefore not breaking anything.
func handlePullWhilstCheckedOut(ctl ctl.DataCTL, remoteHistory VTree, localCommit Commit) (err error) {
	defer util.ReturnErrOnPanic(&err)()
	EnsureByPanic(EnsureNoChanges, "") // just in case

	// It was safe to pull
	writeTREE(remoteHistory)
	return
}

func handlePullWhilstUnbranched(ctl ctl.DataCTL, remoteHistory VTree, localCommit Commit) (err error) {
	defer util.ReturnErrOnPanic(&err)()

	remoteBranch, err := getBranch(ctl, localCommit.Branch)
	util.PanicIfErr(err, "")
	if remoteBranch.HEAD != localCommit.Parent { // If we are behind the remote
		fmt.Println("Behind remote, attempting merge...")
		// It might not be safe to pull
		// Create a new commit that is ahead of the remote, persisting old changes
		remoteHead, newLocal, err := attemptMergeToHead(ctl, localCommit, remoteBranch)
		util.PanicIfErr(err, "")

		var stagemap Stagemap
		parseStagemap(&stagemap)
		err = stagemap.verifyAndSyncWithCommit(newLocal)
		util.PanicIfErr(err, "")

		errStagemap := stagemap.WriteToFile()
		errBranch := remoteBranch.WriteToFolder(remoteFolder)
		errHead := remoteHead.WriteToFolder(remoteFolder)
		errLocal := newLocal.WriteToFolder(localFolder)
		util.PanicIfErr(util.ReturnFirstErr(errStagemap, errBranch, errHead, errLocal), "")
		linkBRANCH(remoteBranch, false)
		linkHEAD(remoteHead)
		linkLOCAL(newLocal)

		// oldLocal is overriden by newLocal as it takes over the commit hash, so this code should be no longer needed
		// oldLocal := localCommit
		// err = os.Remove(oldLocal.ToFilePath(true)) // because we have a successful new local
		// util.PanicIfErr(err, "")

	} else {
		fmt.Println("Up to date")
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
	defer util.ReturnErrOnPanic(&err)()
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
	defer util.ReturnErrOnPanic(&err)()
	EnsureByPanic(EnsureSetup, "")
	EnsureByPanic(EnsureLOCAL, "")
	EnsureByPanic(EnsureConfig, "")

	// Gather required parts
	var ctl ctl.DataCTL
	ctl.ParseFromFile(configPath)
	remoteHistory, err := getVTree(ctl)
	util.PanicIfErr(err, "")
	var localCommit Commit
	parseLOCAL(&localCommit)

	if err = EnsureLatestCommit(); err == errNotHEAD {
		// If we're checked out to some other commit
		fmt.Println("Pulling whilst checked out, merging vtrees...")
		err = handlePullWhilstCheckedOut(ctl, remoteHistory, localCommit)
	} else if err != nil { // something else went wrong
		util.PanicIfErr(err, "")
	} else {
		if isBranched() {
			// We are branched off onto a new branch locally, safe to pull
			fmt.Println("Pulling whilst branched off, merging vtrees and updating locals...")
			err = handlePullWhilstBranched(ctl, remoteHistory, localCommit)
		} else {
			fmt.Println("Pulling whilst working on remotely known branch...")
			err = handlePullWhilstUnbranched(ctl, remoteHistory, localCommit)
		}
	}

	util.PanicIfErr(err, "")

	return
}
