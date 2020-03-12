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

func (c Commit) _convertToDisplayName(file string, fullPath bool, stagemap Stagemap) string {
	if fullPath {
		return file
	}
	if stagemap != nil {
		return stagemap[file]
	}
	return c.idvPathToName(file)
}

// FormatFiles returns the list of files as 1 long string interleaved with \n
// Usefull for display and filtering, head is prepended, tail is appended, isEmpty
// is used if no files are listed as follows: `headisEmptytail`
func (c Commit) FormatFiles(selector *regexp.Regexp, head string, tail string, ifEmpty string, delim string, fullPath bool) string {
	out := head
	for _, file := range c.Files {
		if selector.MatchString(c.idvPathToName(file)) {
			out += c._convertToDisplayName(file, fullPath, nil) + delim
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
func (c Commit) FormatDiff(head string, tail string, ifEmpty string, delim string, fullPath bool, stagemap Stagemap) string {
	out := head
	for _, file := range c.Diff.Added {
		out += "A    " + c._convertToDisplayName(file, fullPath, stagemap) + delim
	}
	for _, file := range c.Diff.Updated {
		out += "U    " + c._convertToDisplayName(file, fullPath, stagemap) + delim
	}
	for _, file := range c.Diff.Removed {
		out += "R    " + c._convertToDisplayName(file, fullPath, stagemap) + delim
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
func (c *Commit) AddOrUpdate(paths []string) (adds, updates map[string]string) {
	adds = make(map[string]string)
	updates = make(map[string]string)
	fileMap := c.filesToNameMap(c.Files)
	addMap := c.filesToNameMap(c.Diff.Added)
	updateMap := c.filesToNameMap(c.Diff.Updated)
	for _, path := range paths {
		filename := filepath.Base(path)
		if _, ok := addMap[filename]; ok { // If already stage for add
			continue
		} else if _, ok := updateMap[filename]; ok { // If already staged for update
			continue
		} else if _, ok := fileMap[filename]; ok { // If already in existing files
			idvpath := c.pathToIDVPath(path)
			c.update(idvpath)
			updates[idvpath] = path
		} else { // Completely new and unseen file
			idvpath := c.pathToIDVPath(path)
			c.add(idvpath)
			adds[idvpath] = path
		}
	}
	return
}

// add adds a new file to this commit
func (c *Commit) add(idvpath string) {
	c.Diff.Added = append(c.Diff.Added, idvpath)
}

// update updates a in this commit
func (c *Commit) update(idvpath string) {
	c.Diff.Updated = append(c.Diff.Updated, idvpath)
}

// _remove actually removes files from the c.Files and (possibly) c.Diff, behave slightly differently based on the passed parameters.
// unstage denotes whether to also remove staged updates/removes.
// asFiles denotes whether the passed strings are already filenames or still paths. if so they are converted first
func (c *Commit) _remove(staging []string, asFiles, unstage bool) (removals, unstages int) {
	fileMap := c.filesToNameMap(c.Files)
	addMap := c.filesToNameMap(c.Diff.Added)
	updateMap := c.filesToNameMap(c.Diff.Updated)
	for _, pathOrFile := range staging {
		var filename string
		if asFiles {
			filename = filepath.Base(pathOrFile)
		}
		if unstage {
			if _, ok := addMap[filename]; ok { // If already staged for add
				addMap[filename] = ""
				unstages++
			} else if _, ok := updateMap[filename]; ok { // If already staged for update
				updateMap[filename] = ""
				unstages++
			}
		}
		if _, ok := fileMap[filename]; ok { // If already in existing files
			c.Diff.Removed = append(c.Diff.Removed, fileMap[filename])
			fileMap[filename] = ""
			removals++
		}
		// unseen/untracked file  --> do nothing
	}
	c.Files = c.fileMaptoSlice(fileMap)
	c.Diff.Added = c.fileMaptoSlice(addMap)
	c.Diff.Updated = c.fileMaptoSlice(updateMap)
	return
}

// removeNames removes each file in the passed list of names from this commit
// expects the paths to be complete, absolute and valid
func (c *Commit) removeFiles(paths []string, unstage bool) (removals, unstages int) {
	return c._remove(paths, true, unstage)
}

// removeNames removes each name in the passed list of names from this commit
func (c *Commit) removeNames(names []string, unstage bool) (removals, unstages int) {
	return c._remove(names, false, unstage)
}

// removeWithSelector stages removal of each file that matches the selector
func (c *Commit) removeWithSelector(selector *regexp.Regexp, unstage bool) (removals, unstages int) {
	var rmFiles, rmAdds, rmUpdates int
	c.Files, rmFiles = util.Filter(selector, c.Files)
	if unstage {
		c.Diff.Added, rmAdds = util.Filter(selector, c.Diff.Added)
		c.Diff.Updated, rmUpdates = util.Filter(selector, c.Diff.Updated)
	}
	removals = rmFiles
	unstages = rmAdds + rmUpdates
	return
}

// unstage removes files that were stages for removal, updates and/or adds
func (c *Commit) unstage(selector *regexp.Regexp) (unstaged int) {
	var unstagedRemoves, unstagedAdds, unstagedUpdates int
	c.Diff.Added, unstagedAdds = util.Filter(selector, c.Diff.Added)
	c.Diff.Updated, unstagedUpdates = util.Filter(selector, c.Diff.Updated)
	c.Diff.Removed, unstagedRemoves = util.Filter(selector, c.Diff.Removed)
	unstaged = unstagedAdds + unstagedUpdates + unstagedRemoves
	return
}

func (c Commit) pathToIDVPath(path string) string {
	return dataFolder + filepath.Base(path) + "/" + c.Hash.String()
}

func (c Commit) idvPathToName(path string) string {
	return filepath.Base(filepath.Dir(path))
}

// filesToNameMap converts a slice of filepaths into a map[filename]filepath, usable for searching
func (c Commit) filesToNameMap(slice []string) map[string]string {
	m := make(map[string]string)
	for _, elem := range slice {
		m[c.idvPathToName(elem)] = elem
	}
	return m
}

// fileMaptoSlice converts a fileMap into a slice if map[elem] is a valid path
func (c Commit) fileMaptoSlice(m map[string]string) []string {
	var slice []string
	for _, val := range m {
		if val != "" {
			slice = append(slice, val)
		}
	}
	return slice
}
