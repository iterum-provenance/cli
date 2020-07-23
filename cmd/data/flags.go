package data

// IsCommit states whether `iterum data checkout` args refers to a commit rather than a branch
var IsCommit bool

// IsHash states whether `iterum data checkout` arg refers to a a branch/commit hash rather than name
var IsHash bool

// AsSelector determines whether the passed argument for `iterum data rm` should be used as a selector over the committed files
var AsSelector bool

// Unstage determines whether the current `iterum data rm` should also unstage files staged for adds/updates
var Unstage bool

// Branches determines whether `iterum data ls` should show a list of branches rather than files
var Branches bool

// Commits determines whether `iterum data ls` should show a list of commits in order, rather than files
var Commits bool

// ShowFullPath determines whether `iterum data ls/status` should show the real filepath rather than just the name
var ShowFullPath bool

// ShowLocalPath determines whether `iterum data status` should show the local filepath to the staged files rather than just the name
var ShowLocalPath bool

// CommitHashOrName stores a hash or name passed to `iterum data branch/checkout`
var CommitHashOrName string

// ShowFiles determines whether files are displayed before downloading them
var ShowFiles bool

// NoPrompt determines whether files are downloaded without prompting for assurance
var NoPrompt bool
