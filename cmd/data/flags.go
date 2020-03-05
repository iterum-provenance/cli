package data

// Recursive states whether `iterum data add/rm` should recurse into folders found in the specified arguments
var Recursive bool

// IsCommit states whether `iterum data checkout` args refers to a commit rather than a branch
var IsCommit bool

// IsHash states whether `iterum data checkout` arg refers to a a branch/commit hash rather than name
var IsHash bool
