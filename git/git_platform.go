package git

import (
	"bytes"
	"errors"
	"fmt"
)

// Platform is used for defining which git platform is used
type Platform string

// Enum-like values allowed for git.Platform type
const (
	Github    Platform = "github"
	Gitlab    Platform = "gitlab"
	Bitbucket Platform = "bitbucket"
	None      Platform = "NOPLATFORM"
)

// NewPlatform creates a new Platform instance and validates it
func NewPlatform(rawPlatform string) (Platform, error) {
	var gp Platform = Platform(rawPlatform)
	return gp, gp.IsValid()
}

// IsValid checks the validity of the Platform
func (gp Platform) IsValid() error {
	switch gp {
	case Github, Gitlab, Bitbucket, None:
		return nil
	}
	return errors.New("Error: Invalid git.Platform type")
}

// String converts Platform to string representation
func (gp Platform) String() string {
	return string(gp)
}

// AllowedVariables returns a formatted string on how to set this type with the set command
func (gp Platform) AllowedVariables() string {
	var buf bytes.Buffer
	// Excluding None, because should not be set to none later
	fmt.Fprintf(&buf, "Platform            { %v, %v, %v }\n", Github, Gitlab, Bitbucket)
	return buf.String()
}
