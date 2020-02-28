package git

import (
	"bytes"
	"errors"
	"fmt"
)

// GitPlatform is used for defining which platform is used
type GitPlatform string

// Enum-like values allowed for GitPlatform type
const (
	Github    GitPlatform = "github"
	Gitlab    GitPlatform = "gitlab"
	Bitbucket GitPlatform = "bitbucket"
	None      GitPlatform = "NOPLATFORM"
)

// NewGitPlatform creates a new GitPlatform instance and validates it
func NewGitPlatform(rawPlatform string) (GitPlatform, error) {
	var gp GitPlatform = GitPlatform(rawPlatform)
	return gp, gp.IsValid()
}

// IsValid checks the validity of the GitPlatform
func (gp GitPlatform) IsValid() error {
	switch gp {
	case Github, Gitlab, Bitbucket, None:
		return nil
	}
	return errors.New("Error: Invalid GitPlatform type")
}

// String converts GitPlatform to string representation
func (gp GitPlatform) String() string {
	return string(gp)
}

// AllowedVariables returns a formatted string on how to set this type with the set command
func (gp GitPlatform) AllowedVariables() string {
	var buf bytes.Buffer
	// Excluding None, because should not be set to none later
	fmt.Fprintf(&buf, "Platform            { %v, %v, %v }\n", Github, Gitlab, Bitbucket)
	return buf.String()
}
