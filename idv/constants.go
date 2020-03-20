package idv

// useful folders folders for idv
const (
	// IDVFolder is the name of the folder these files are stored in (like .git)
	IDVFolder string = ".idv/"
	// remoteFolder is the folder that holds copies of remote files, these files should not be edited
	remoteFolder string = IDVFolder + "remote/"
	// localFolder is the folder non-remote files and updates are stored in (to-be-committed)
	localFolder string = IDVFolder + "local/"
)

const (
	// LOCAL points to CurrentCommitFile, the local commit containing staged updates not yet pushed
	// Or in case of checking out to a commit which is not the head of a branch, it points to HEAD
	LOCAL string = IDVFolder + "LOCAL"
	// HEAD points to the commit that local repo thinks is the head of the current branch remotely,
	// essentially the parent of LOCAL
	HEAD string = IDVFolder + "HEAD"
	// BRANCH points to the branch that local operates on
	BRANCH string = IDVFolder + "BRANCH"
	// TREE points to the .vtree that local thinks is the most recent version of .vtree
	TREE string = IDVFolder + "TREE"
)

// vtreeFileName is the name of the stored version tree file
const vtreeFileName string = "history" + vtreeFileExt

// Filenames of the files stored in .idv/local/
const (
	// curCommitFilePath is the name of the stored version tree file
	curCommitFilePath string = localFolder + "current" + commitFileExt

	// curBranchFilePath is the name of the stored version tree file
	curBranchFilePath string = localFolder + "current" + branchFileExt

	// vtreeFilePath is the name of the stored version tree file
	vtreeFilePath string = remoteFolder + vtreeFileName

	configPath string = "idv-config.yaml"
)

// File extensions for idv files
const (
	commitFileExt string = ".commit"
	branchFileExt string = ".branch"
	vtreeFileExt  string = ".vtree"
)

// Defaults for roots and masters
const (
	masterBranchName string = "master"
	rootCommitName   string = "<root>"
)

const (
	dataFolder     string = "data/"
	stagedFileName string = "local.staged"
	stagedFilePath string = localFolder + stagedFileName
)
