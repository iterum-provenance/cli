package credentials

import (
	"fmt"

	"github.com/iterum-provenance/cli/util"
)

// Local holds required credentials for storing data locally
type Local struct {
	Path string `yaml:"path" json:"path"` // Absolute path to folder
}

// IsValid checks the validity of this structure
func (l Local) IsValid() error {
	if util.DirExists(l.Path) {
		return nil
	}
	return fmt.Errorf("Error: %v is not an existing directory", l.Path)
}

// ParseLocal tries to parse an interface as this credential storage
func ParseLocal(raw map[string]interface{}) (Local, error) {
	l := Local{
		Path: raw["path"].(string),
	}
	return l, l.IsValid()
}

// GetLocation returns a string path or url to where the data is located based on the backend
func (l Local) GetLocation() string {
	return l.Path
}
