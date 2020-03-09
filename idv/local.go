package idv

import (
	"fmt"

	"github.com/Mantsje/iterum-cli/util"
)

// Config the config stuff that iterum uses to keep state of a data versioned repo
// Things like current commit, branch, etc
type Config struct {
	LocalChanged   bool
	CurrentBranch  hash
	OriginalCommit hash // original commit that local is forking on
}

// WriteToFile writes the config to the specified file.
func (c Config) WriteToFile(filePath string) error {
	return util.WriteJSONFile(filePath, c)
}

// ParseFromFile tries to parse a idv config file
func (c *Config) ParseFromFile(filepath string) error {
	if err := util.ReadJSONFile(filepath, c); err != nil {
		return fmt.Errorf("Error: Could not parse config.idv due to `%v`", err)
	}
	return nil
}
