package config

import "errors"

// ProjectType is used for defining which platform is used
type ProjectType string

// Enum-like values allowed for ProjectType type
const (
	Local       ProjectType = "local"
	Distributed ProjectType = "distributed"
)

// NewProjectType creates a new ProjectType instance and validates it
func NewProjectType(rawProjectType string) (ProjectType, error) {
	var pt ProjectType = ProjectType(rawProjectType)
	return pt, pt.IsValid()
}

// IsValid checks the validity of the ProjectType
func (pt ProjectType) IsValid() error {
	switch pt {
	case Local, Distributed:
		return nil
	}
	return errors.New("Error: Invalid ProjectType")
}
