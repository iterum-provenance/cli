package idv

import (
	"errors"
	"fmt"

	"github.com/Mantsje/iterum-cli/util"
)

// Stagemap is a structure mapping idv filepaths to absolute filepaths on this machine
type Stagemap map[string]string

// WriteToFile writes the config to the specified file.
func (s Stagemap) WriteToFile() error {
	return util.WriteJSONFile(stagedFilePath, s)
}

// ParseFromFile tries to parse a idv config file
func (s *Stagemap) ParseFromFile(filepath string) error {
	if err := util.ReadJSONFile(filepath, s); err != nil {
		return fmt.Errorf("Error: Could not parse %v due to `%v`", stagedFileName, err)
	}
	return s.Verify()
}

// Verify ensures all file pointers in the staged file actually exist
func (s Stagemap) Verify() (err error) {
	for _, absPath := range s {
		if !util.FileExists(absPath) {
			err = errors.New("Error: Staged filelist contains non-existent files")
			break
		}
	}
	return
}

// verifyAndSyncWithCommit uses the passed commit to update the stagemap by removing dangling file maps and ensuring mappings for all staged changes
func (s *Stagemap) verifyAndSyncWithCommit(commit Commit) (err error) {
	// Check if all staged are mapped
	allStages := append(commit.Diff.Added, commit.Diff.Updated...)
	for _, idvfile := range allStages {
		if _, ok := (*s)[idvfile]; !ok { // if file is not in stagemap
			return fmt.Errorf("Error: Staged file not found locally, try unstaging and adding/updating again: %v", idvfile)
		}
	}
	// Check if no mapped or unstaged, if so remove them (this is save, because it it redundant)
	addMap := commit.filesToNameMap(commit.Diff.Added)
	updateMap := commit.filesToNameMap(commit.Diff.Updated)
	for key := range *s {
		keyAsName := commit.idvPathToName(key)
		if _, ok := addMap[keyAsName]; !ok { // if stagemap file is not in staged adds
			if _, ok := updateMap[keyAsName]; !ok { // if stagemap file is also not in staged updates
				// Then we don't need it and can remove it
				delete(*s, key)
			}
		}

	}
	return
}

// update takes a map similar to Stagemap and adds its entries
func (s *Stagemap) update(m map[string]string) (err error) {
	for key, val := range m {
		if _, ok := (*s)[key]; ok {
			return errors.New("Error: Duplicate stagemap key in Update of stagemap")
		}
		(*s)[key] = val
	}
	return
}
