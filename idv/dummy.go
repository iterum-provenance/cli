package idv

func dummyPull() {
	branch := NewBranch(masterBranchName)
	commit := NewRootCommit(branch.Hash)
	branch.HEAD = commit.Hash
	history := NewVTree(commit, branch)
	history.WriteToFolder(remoteFolder)
	commit.WriteToFolder(remoteFolder)
	branch.WriteToFolder(remoteFolder)
}
