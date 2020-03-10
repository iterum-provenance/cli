package idv

func dummyPull() {
	branch := NewBranch(masterBranchName)
	commit := NewRootCommit(branch.Hash)
	branch.HEAD = commit.Hash
	history := NewVTree(commit, branch)
	history.WriteToFolder(".idv/remote")
	commit.WriteToFolder(".idv/remote")
	branch.WriteToFolder(".idv/remote")
}
