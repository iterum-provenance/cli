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
func (c Local) IsValid() error {
	if util.DirExists(c.Path) {
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
