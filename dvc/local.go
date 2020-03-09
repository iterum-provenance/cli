package dvc

// Config the config stuff that iterum uses to keep state of a data versioned repo
// Things like current commit, branch, etc
type Config struct {
	CurrentBranch hash
	CurrentCommit hash
}
