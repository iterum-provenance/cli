package config

import "errors"

// RepoType is used for defining what this repo holds
type RepoType string

// Enum-like values allowed for RepoType type
const (
	Unit    RepoType = "unit"
	Flow    RepoType = "flow"
	Project RepoType = "project"
)

// NewRepoType creates a new RepoType instance and validates it
func NewRepoType(repoType string) (RepoType, error) {
	var rt RepoType = RepoType(repoType)
	return rt, rt.IsValid()
}

// IsValid checks the validity of the RepoType
func (rt RepoType) IsValid() error {
	switch rt {
	case Unit, Flow, Project:
		return nil
	}
	return errors.New("Error: Invalid RepoType")
}

func (rt RepoType) String() string {
	return string(rt)
}
