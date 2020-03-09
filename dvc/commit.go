package dvc

import (
	"fmt"

	"github.com/Mantsje/iterum-cli/util"
)

// structure that defines whether a commit is deprecated and why
type deprecated struct {
	Value  bool
	Reason string
}

// diff defines the updates included in this commit
type diff struct {
	Added   []string
	Removed []string
	Updated []string
}

// Commit internally defines a data version commit file
type Commit struct {
	Parent      hash
	Branch      hash
	Hash        hash
	Name        string
	Description string
	Files       []string
	Diff        diff
	Deprecated  deprecated
}

// NewCommit creates a bare new commit
func NewCommit(parent Commit, branch hash, name string, desc string) Commit {
	return Commit{
		Parent:      parent.Hash,
		Branch:      branch,
		Name:        name,
		Description: desc,
		Hash:        newHash(32),
		Files:       parent.Files,
		Diff:        diff{[]string{}, []string{}, []string{}},
		Deprecated:  deprecated{false, ""},
	}
}

// WriteToFolder writes the commit to the specified folder.
// Name of file is and should be determined by the commit structure
func (c Commit) WriteToFolder(folderPath string) error {
	fullPath := folderPath + "/" + c.Hash.String() + ".commit"
	return util.WriteJSONFile(fullPath, c)
}

// ParseFromFile tries to parse a .commit file
func (c *Commit) ParseFromFile(filepath string) error {
	if err := util.ReadJSONFile(filepath, c); err != nil {
		return fmt.Errorf("Error: Could not parse Commit due to `%v`", err)
	}
	return nil
}

// FilesAsString returns the list of files as 1 long string interleaved with \n
// Usefull for display and filtering, head is prepended, tail is appended, isEmpty
// is used if no files are listed as follows: ```headisEmptytail```
func (c Commit) FilesAsString(head string, tail string, ifEmpty string) string {
	out := head
	for _, file := range c.Files {
		out += file + "\n"
	}
	if len(c.Files) == 0 {
		out += ifEmpty
	}
	out += tail
	return out
}
