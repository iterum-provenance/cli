package credentials

// Local holds required credentials for storing data locally
type Local struct {
	Path string `yaml:"path" json:"path"` // Absolute path to folder
}

// IsValid checks the validity of this structure
func (l Local) IsValid() error {
	return nil
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
