package config

import "errors"

// GitProtocol is used for defining which protocol to use
type GitProtocol string

// Enum-like values allowed for GitProtocol type
const (
	SSH   GitProtocol = "ssh"
	HTTPS GitProtocol = "https"
)

// IsValid checks the validity of the GitProtocol
func (gp GitProtocol) IsValid() error {
	switch gp {
	case SSH, HTTPS:
		return nil
	}
	return errors.New("Invalid GitProtocol type")
}
