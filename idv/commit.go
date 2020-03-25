package idv

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/Mantsje/iterum-cli/util"
)

// structure that defines whether a commit is deprecated and why
type deprecated struct {
	Value  bool   `json:"value"`
	Reason string `json:"reason"`
}

// diff defines the updates included in this commit
type diff struct {
	Added   []string `json:"added"`
	Removed []string `json:"removed"`
	Updated []string `json:"updated"`
}

// Commit internally defines a data version commit file
type Commit struct {
	Parent      hash       `json:"parent"`
	Branch      hash       `json:"branch"`
	Hash        hash       `json:"hash"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Files       []string   `json:"files"`
	Diff        diff       `json:"diff"`
	Deprecated  deprecated `json:"deprecated"`
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

// ---------- ---------- ---------- Formatting/Representation ---------- ---------- ---------- \\

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

// ---------- ---------- ---------- Add/Update ---------- ---------- ---------- \\

// addOrUpdate adds or updates files to this commit
// These files needs to be guaranteed to exist before passing
func (c *Commit) addOrUpdate(paths []string) (adds, updates map[string]string) {
	adds = make(map[string]string)
	updates = make(map[string]string)
	fileMap := c.filesToNameMap(c.Files)
	addMap := c.filesToNameMap(c.Diff.Added)
	updateMap := c.filesToNameMap(c.Diff.Updated)
	removeMap := c.filesToNameMap(c.Diff.Removed)
	for _, path := range paths {
		filename := filepath.Base(path)
		if _, ok := addMap[filename]; ok { // If already staged for add
			continue
		} else if _, ok := updateMap[filename]; ok { // If already staged for update
			continue
		} else if _, ok := removeMap[filename]; ok { // If file already staged for removal (skip)
			continue
		} else if _, ok := fileMap[filename]; ok { // If already in existing files, but unstaged -> update
			idvpath := c.pathToIDVPath(path)
			c.Diff.Updated = append(c.Diff.Updated, idvpath)
			updateMap[filename] = idvpath
			updates[idvpath] = path
		} else { // Completely new and unseen file -> add
			idvpath := c.pathToIDVPath(path)
			c.Diff.Added = append(c.Diff.Added, idvpath)
			addMap[filename] = idvpath
			adds[idvpath] = path
		}
	}
	return
}

// ---------- ---------- ---------- Removing ---------- ---------- ---------- \\

// _remove actually removes files from the c.Files and (possibly) c.Diff, behave slightly differently based on the passed parameters.
// unstage denotes whether to also remove staged updates/removes.
// asFiles denotes whether the passed strings are already filenames or still paths. if so they are converted first
func (c *Commit) _remove(staging []string, asFiles, unstage bool) (removals, unstages int) {
	fileMap := c.filesToNameMap(c.Files)
	addMap := c.filesToNameMap(c.Diff.Added)
	updateMap := c.filesToNameMap(c.Diff.Updated)
	removeMap := c.filesToNameMap(c.Diff.Removed)
	for _, pathOrFile := range staging {
		var filename string = pathOrFile
		if asFiles {
			filename = filepath.Base(pathOrFile)
		}
		if _, ok := addMap[filename]; unstage && ok { // If unstage && already staged for add
			addMap[filename] = ""
			unstages++
		} else if _, ok := updateMap[filename]; ok { // If already staged for update
			if unstage {
				updateMap[filename] = ""
				unstages++
			} else {
				// Skip removes of files staged to be updated unless explicitly said to do so
				continue
			}
		}
		if _, ok := removeMap[filename]; ok {
			continue // No need to remove already removed files
		}
		if _, ok := fileMap[filename]; ok { // If already in existing files
			c.Diff.Removed = append(c.Diff.Removed, fileMap[filename])
			removeMap[filename] = fileMap[filename] // add it to removeMap in case of multiple deletion of the same file in 1 statement
			removals++
		}
		// unseen/untracked file  --> do nothing
	}
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
	var matchingFiles, rmAdds, rmUpdates []string
	matchingFiles = util.SelectMatching(selector, c.Files)
	removeMap := c.filesToNameMap(c.Diff.Removed)
	updateMap := c.filesToNameMap(c.Diff.Updated)
	if unstage {
		c.Diff.Added, rmAdds = util.FilterSplit(selector, c.Diff.Added)
		c.Diff.Updated, rmUpdates = util.FilterSplit(selector, c.Diff.Updated)
	}
	for _, file := range matchingFiles {
		if _, ok := removeMap[c.idvPathToName(file)]; ok {
			continue
		}
		if _, ok := updateMap[c.idvPathToName(file)]; !unstage && ok {
			continue // Don't delete files staged for updates unless explicitly specified
		}

		c.Diff.Removed = append(c.Diff.Removed, file)
		removals++
	}
	unstages = len(rmAdds) + len(rmUpdates)
	return
}

// ---------- ---------- ---------- Unstaging ---------- ---------- ---------- \\

// unstage removes files that were stages for removal, updates and/or adds
func (c *Commit) unstage(selector *regexp.Regexp) (unstaged int) {
	var unstagedRemoves, unstagedAdds, unstagedUpdates int
	c.Diff.Added, unstagedAdds = util.Filter(selector, c.Diff.Added)
	c.Diff.Updated, unstagedUpdates = util.Filter(selector, c.Diff.Updated)
	c.Diff.Removed, unstagedRemoves = util.Filter(selector, c.Diff.Removed)
	unstaged = unstagedAdds + unstagedUpdates + unstagedRemoves
	return
}

// ---------- ---------- ---------- Finalizing commit ---------- ---------- ---------- \\

// applyStaged applies all the staged changes to the filelist of the commit, such that it is ready for pushing
func (c *Commit) applyStaged() (err error) {
	fileMap := c.filesToNameMap(c.Files)
	for _, addition := range c.Diff.Added {
		fileMap[c.idvPathToName(addition)] = addition
	}
	for _, update := range c.Diff.Updated {
		if _, ok := fileMap[c.idvPathToName(update)]; ok {
			fileMap[c.idvPathToName(update)] = update
		} else {
			return fmt.Errorf("Error: staged update is not an update, because original is not there: %v", update)
		}
	}
	for _, removal := range c.Diff.Removed {
		if _, ok := fileMap[c.idvPathToName(removal)]; ok {
			fileMap[c.idvPathToName(removal)] = ""
		} else {
			return fmt.Errorf("Error: staged removal is invalid, original is not there: %v", removal)
		}
	}
	c.Files = c.fileMaptoSlice(fileMap)
	return
}

// ---------- ---------- ---------- Utility functions ---------- ---------- ---------- \\

// idvPathToName converts a path (/home/user/path/to/file.extension) into an idvpath (data/file.extension/c.Hash)
func (c Commit) pathToIDVPath(path string) string {
	return filepath.Base(path) + "/" + c.Hash.String()
}

// idvPathToName converts an idv path (data/filename/commithash) into just the filename
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
func (c Commit) fileMaptoSlice(m map[string]string) (slice []string) {
	slice = []string{}
	for _, val := range m {
		if val != "" {
			slice = append(slice, val)
		}
	}
	return
}

// containsChanges returns a bool stating whether this commit has staged anything
func (c Commit) containsChanges() bool {
	return !(len(c.Diff.Added) == 0 && len(c.Diff.Updated) == 0 && len(c.Diff.Removed) == 0)
}

// ---------- ---------- ---------- Merging functions ---------- ---------- ---------- \\

// autoMerge automatically merges all the staged changes of a given commit into c
// from 					Denotes the commit to take all staged changes from
// Files staged for REMOVE are removed regardless of updated versions. it is omitted if already removed
// Files staged for ADD are either added, or update the current version in case of an existing one
// Files staged for UPDATE are updated or in case the original is gone added
func (c *Commit) autoMerge(from Commit) {
	fileMap := c.filesToNameMap(c.Files)
	for _, removal := range from.Diff.Removed {
		filename := from.idvPathToName(removal)
		if matchedFile, ok := fileMap[filename]; ok {
			if matchedFile == removal {
				c.Diff.Removed = append(c.Diff.Removed, removal)
			} else {
				fmt.Printf("An updated version of %v was found but still removed\n", filename)
				c.Diff.Removed = append(c.Diff.Removed, removal)
			}
		} else {
			fmt.Printf("%v already removed, skipping\n", filename)
		}
	}
	for _, updated := range from.Diff.Updated {
		filename := from.idvPathToName(updated)
		if matchedFile, ok := fileMap[filename]; ok {
			if matchedFile != updated {
				fmt.Printf("%v was staged for UPDATE, but updated in subsequent commit(s), staging UPDATE over this file instead\n", filename)
			}
			c.Diff.Updated = append(c.Diff.Updated, updated)
		} else {
			fmt.Printf("%v was staged for UPDATE, but removed in subsequent commit(s), staging for ADD instead\n", filename)
			c.Diff.Added = append(c.Diff.Added, updated)
		}
	}
	for _, addition := range from.Diff.Added {
		filename := from.idvPathToName(addition)
		if _, ok := fileMap[filename]; ok {
			fmt.Printf("%v was staged for ADD, but already present, staging for UPDATE instead\n", filename)
			c.Diff.Updated = append(c.Diff.Updated, addition)
		} else {
			c.Diff.Added = append(c.Diff.Added, addition)
		}
	}
}
