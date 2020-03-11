package idv

import (
	"fmt"
	"path/filepath"
	"regexp"

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

// writeToFile writes the commit to the specified file.
// Don't use this unless you really know what you're doing, try WriteToFolder instead
func (c Commit) writeToFile(filePath string) error {
	return util.WriteJSONFile(filePath, c)
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
func (c Commit) FormatFiles(selector *regexp.Regexp, head string, tail string, ifEmpty string, delim string) string {
	out := head
	for _, file := range c.Files {
		if selector.MatchString(c.idvPathToName(file)) {
			out += c.idvPathToName(file) + delim
		}
	}
	if out == head {
		out += ifEmpty
	} else {
		out = out[:len(out)-len(delim)]
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
		out += "A    " + c.idvPathToName(file) + delim
	}
	for _, file := range c.Diff.Updated {
		out += "U    " + c.idvPathToName(file) + delim
	}
	for _, file := range c.Diff.Removed {
		out += "R    " + c.idvPathToName(file) + delim
	}
	if out == head {
		out += ifEmpty
	} else {
		out = out[:len(out)-len(delim)]
	}
	out += tail + "\n"
	return out
}

// AddOrUpdate adds or updates files to this commit
// These files needs to be guaranteed to exist before passing
func (c *Commit) AddOrUpdate(paths []string) (adds, updates int) {
	fileMap := c._filesToNameMap(c.Files)
	addMap := c._filesToNameMap(c.Diff.Added)
	updateMap := c._filesToNameMap(c.Diff.Updated)
	for _, path := range paths {
		filename := filepath.Base(path)
		if _, ok := addMap[filename]; ok { // If already stage for add
			continue
		} else if _, ok := updateMap[filename]; ok { // If already staged for update
			continue
		} else if _, ok := fileMap[filename]; ok { // If already in existing files
			c.update(path)
			updates++
		} else { // Completely new and unseen file
			c.add(path)
			adds++
		}
	}
	return
}

// add adds a new file to this commit
func (c *Commit) add(file string) {
	c.Diff.Added = append(c.Diff.Added, c.pathToIDVPath(file))
}

// update updates a in this commit
func (c *Commit) update(file string) {
	c.Diff.Updated = append(c.Diff.Updated, c.pathToIDVPath(file))
}

func (c *Commit) _remove(staging []string, asFiles bool) (removals int) {
	fileMap := c._filesToNameMap(c.Files)
	addMap := c._filesToNameMap(c.Diff.Added)
	updateMap := c._filesToNameMap(c.Diff.Updated)
	for _, pathOrFile := range staging {
		var filename string
		if asFiles {
			filename = filepath.Base(pathOrFile)
		}
		if _, ok := addMap[filename]; ok { // If already staged for add
			addMap[filename] = ""
			removals++
		} else if _, ok := updateMap[filename]; ok { // If already staged for update
			updateMap[filename] = ""
			removals++
		} else if _, ok := fileMap[filename]; ok { // If already in existing files
			c.Diff.Removed = append(c.Diff.Removed, fileMap[filename])
			fileMap[filename] = ""
			removals++
		}
		// unseen/tracked file
	}
	c.Files = c._fileMaptoSlice(fileMap)
	c.Diff.Added = c._fileMaptoSlice(addMap)
	c.Diff.Updated = c._fileMaptoSlice(updateMap)
	return
}

// removeNames removes each file in the passed list of names from this commit
// expects the paths to be complete, absolute and valid
func (c *Commit) removeFiles(paths []string) (removals int) {
	return c._remove(paths, true)
}

// removeNames removes each name in the passed list of names from this commit
func (c *Commit) removeNames(names []string) (removals int) {
	return c._remove(names, false)
}

func (c Commit) pathToIDVPath(path string) string {
	return dataFolder + filepath.Base(path) + "/" + c.Hash.String()
}

func (c Commit) idvPathToName(path string) string {
	return filepath.Base(filepath.Dir(path))
}

// _filesToNameMap converts a slice of filepaths into a map[filename]filepath, usable for searching
func (c Commit) _filesToNameMap(slice []string) map[string]string {
	m := make(map[string]string)
	for _, elem := range slice {
		m[c.idvPathToName(elem)] = elem
	}
	return m
}

// _fileMaptoSlice converts a fileMap into a slice if map[elem] is a valid path
func (c Commit) _fileMaptoSlice(m map[string]string) []string {
	var slice []string
	for _, val := range m {
		if val != "" {
			slice = append(slice, val)
		}
	}
	return slice
}
