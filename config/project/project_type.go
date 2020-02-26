package project

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

// InferProjectType tries to infer a ProjectType from a string
func InferProjectType(raw string) (ProjectType, error) {
	if raw == "l" { // Shorthand for local
		return NewProjectType(string(Local))
	} else if raw == "d" { // Shorthand for distribbuted
		return NewProjectType(string(Distributed))
	}
	return NewProjectType(raw)
}

// IsValid checks the validity of the ProjectType
func (pt ProjectType) IsValid() error {
	switch pt {
	case Local, Distributed:
		return nil
	}
	return errors.New("Error: Invalid ProjectType")
}

// String converts ProjectType to string representation
func (pt ProjectType) String() string {
	return string(pt)
}
