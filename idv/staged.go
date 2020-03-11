package idv

import (
	"errors"
	"fmt"

	"github.com/Mantsje/iterum-cli/util"
)

// Staging is a structure mapping idv filepaths to absolute filepaths on this machine
type Staging map[string]string

// WriteToFile writes the config to the specified file.
func (s Staging) WriteToFile() error {
	return util.WriteJSONFile(stagedFilePath, s)
}

// ParseFromFile tries to parse a idv config file
func (s *Staging) ParseFromFile(filepath string) error {
	if err := util.ReadJSONFile(filepath, s); err != nil {
		return fmt.Errorf("Error: Could not parse %v due to `%v`", stagedFileName, err)
	}
	return nil
}

// Verify ensures all file pointers in the staged file actually exist
func (s Staging) Verify() (err error) {
	for _, absPath := range s {
		if !util.FileExists(absPath) {
			err = errors.New("Error: Staged filelist contains non-existent files")
			break
		}
	}
	return
}
