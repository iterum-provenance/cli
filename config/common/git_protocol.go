package config

import "errors"

// GitProtocol is used for defining which protocol to use
type GitProtocol string

// Enum-like values allowed for GitProtocol type
const (
	SSH   GitProtocol = "ssh"
	HTTPS GitProtocol = "https"
)

// NewGitProtocol creates a new GitProtocol instance and validates it
func NewGitProtocol(rawProtocol string) (GitProtocol, error) {
	var gp GitProtocol = GitProtocol(rawProtocol)
	return gp, gp.IsValid()
}

// IsValid checks the validity of the GitProtocol
func (gp GitProtocol) IsValid() error {
	switch gp {
	case SSH, HTTPS:
		return nil
	}
	return errors.New("Error: Invalid GitProtocol type")
}
