package data

// Recursive states whether `iterum data add/rm` should recurse into folders found in the specified arguments
var Recursive bool

// ShowExcluded states whether the list of excluded files should be shown
var ShowExcluded bool

// Exclusions holds files and folders to be excluded from adding/removing, etc
var Exclusions []string

// IsCommit states whether `iterum data checkout` args refers to a commit rather than a branch
var IsCommit bool

// IsHash states whether `iterum data checkout` arg refers to a a branch/commit hash rather than name
var IsHash bool
