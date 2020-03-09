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
	LOCAL string = IDVFolder + "LOCAL"
	// HEAD points to the commit that local repo thinks is the head of the current branch remotely,
	// essentially the parent of LOCAL
	HEAD string = IDVFolder + "HEAD"
	// BRANCH points to the branch that local operates on
	BRANCH string = IDVFolder + "BRANCH"
)

// vtreeFileName is the name of the stored version tree file
const vtreeFileName string = remoteFolder + "history" + vtreeFileExt

// Filenames of the files stored in .idv/local/
const (
	// curCommitFile is the name of the stored version tree file
	curCommitFileName string = localFolder + "current" + commitFileExt

	// curBranchFile is the name of the stored version tree file
	curBranchFileName string = localFolder + "current" + branchFileExt

	// curVTreeFile is the name of the stored version tree file
	curVTreeFileName string = localFolder + "current" + vtreeFileExt
)

// File extensions for idv files
const (
	commitFileExt string = ".commit"
	branchFileExt string = ".branch"
	vtreeFileExt  string = ".vtree"
)
