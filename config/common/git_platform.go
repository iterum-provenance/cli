package config

import "errors"

// GitPlatform is used for defining which platform is used
type GitPlatform string

// Enum-like values allowed for GitPlatform type
const (
	Github    GitPlatform = "github"
	Gitlab    GitPlatform = "gitlab"
	Bitbucket GitPlatform = "bitbucket"
)

// NewGitPlatform creates a new GitPlatform instance and validates it
func NewGitPlatform(rawPlatform string) (GitPlatform, error) {
	var gp GitPlatform = GitPlatform(rawPlatform)
	return gp, gp.IsValid()
}

// IsValid checks the validity of the GitPlatform
func (gp GitPlatform) IsValid() error {
	switch gp {
	case Github, Gitlab, Bitbucket:
		return nil
	}
	return errors.New("Error: Invalid GitPlatform type")
}
