package credentials

import (
	"errors"

	"github.com/Mantsje/iterum-cli/util"
)

// Local holds required credentials for storing data locally
type Local struct {
	Path string // Absolute path to folder
}

// IsValid checks the validity of this structure
func (l Local) IsValid() error {
	if util.DirExists(l.Path) {
		return nil
	}
	return errors.New("Error: Path is not an existing directory")
}

// ParseLocal tries to parse an interface as this credential storage
func ParseLocal(raw map[interface{}]interface{}) (Local, error) {
	l := Local{
		Path: raw["path"].(string),
	}
	return l, l.IsValid()
}

// GetLocation returns a string path or url to where the data is located based on the backend
func (l Local) GetLocation() string {
	return l.Path
}
