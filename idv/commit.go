package idv

import (
	"errors"
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

// NewRootCommit creates a bare, empty root commit for initializing a new repo
func NewRootCommit(branch hash) Commit {
	return Commit{
		Parent:      "",
		Branch:      branch,
		Name:        rootCommitName,
		Description: "Root commit of this repository",
		Hash:        newHash(32),
		Files:       []string{},
		Diff:        diff{[]string{}, []string{}, []string{}},
		Deprecated:  deprecated{false, ""},
	}
}

// WriteToFolder writes the commit to the specified folder.
// Name of file is and should be determined by the commit structure
func (c Commit) WriteToFolder(folderPath string) error {
	fullPath := folderPath + "/" + c.Hash.String() + commitFileExt
	return util.WriteJSONFile(fullPath, c)
}

// ParseFromFile tries to parse a .commit file
func (c *Commit) ParseFromFile(filepath string) error {
	if err := util.ReadJSONFile(filepath, c); err != nil {
		return fmt.Errorf("Error: Could not parse Commit due to `%v`", err)
	}
	return nil
}

// ToFilePath returns a path to this commit being: .idv/{local, remote}/commithash.extension
// local indicates which of the 2 folders to use
func (c Commit) ToFilePath(local bool) string {
	if local {
		return localFolder + c.Hash.String() + commitFileExt
	}
	return remoteFolder + c.Hash.String() + commitFileExt
}

// FormatFiles returns the list of files as 1 long string interleaved with \n
// Usefull for display and filtering, head is prepended, tail is appended, isEmpty
// is used if no files are listed as follows: `headisEmptytail`
func (c Commit) FormatFiles(head string, tail string, ifEmpty string, delim string) string {
	out := head
	for _, file := range c.Files {
		out += file + delim
	}
	if len(c.Files) == 0 {
		out += ifEmpty
	}
	out += tail + "\n"
	return out
}

// FormatDiff returns the list of changed files in the Diff as 1 long string interleaved with \n
// Usefull for display and filtering, head is prepended, tail is appended, isEmpty
// is used if no files are listed as follows: `headisEmptytail\n`
// else the strings produced are `OP file\n`, where OP is one of { U(pdated), R(emoved), A(dded)}
func (c Commit) FormatDiff(head string, tail string, ifEmpty string, delim string) string {
	out := head
	for _, file := range c.Diff.Added {
		out += "A    " + file + delim
	}
	for _, file := range c.Diff.Updated {
		out += "U    " + file + delim
	}
	for _, file := range c.Diff.Removed {
		out += "R    " + file + delim
	}
	if out == head {
		out += ifEmpty
	}
	out += tail + "\n"
	return out
}

// Add adds a new file to this commit
func (c *Commit) Add(file string) error {
	return errors.New("Error: Adding not yet implemented")
	// TODO ensure it is actually new, prevent clashes
	// c.Diff.Added = append(c.Diff.Added, file)
	// return nil
}

// Remove removes a file from this commit
func (c *Commit) Remove(file string) error {
	return errors.New("Error: Removing not yet implemented")
	// TODO ensure it was existent
	// c.Diff.Removed = append(c.Diff.Removed, file)
	// return nil
}

// Update updates a in this commit
func (c *Commit) Update(file string) error {
	return errors.New("Error: Updating not yet implemented")
	// TODO ensure it is actually new, prevent clashes
	// c.Diff.Updated = append(c.Diff.Updated, file)
	// return nil
}

// ToBoolMap converts a slice into a map[T]true, used for searching
func toBoolMap(slice []interface{}) map[interface{}]bool {
	m := make(map[interface{}]bool)
	for _, elem := range slice {
		m[elem] = true
	}
	return m
}

// ToSlice converts a boolmap into a slice if map[elem] == true
func toSlice(m map[interface{}]bool) []interface{} {
	var slice []interface{}
	for key, val := range m {
		if val {
			slice = append(slice, key)
		}
	}
	return slice
}
