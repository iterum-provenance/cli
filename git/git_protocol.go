package git

import (
	"bytes"
	"errors"
	"fmt"
)

// Protocol is used for defining which protocol to use
type Protocol string

// Enum-like values allowed for Protocol type
const (
	SSH   Protocol = "ssh"
	HTTPS Protocol = "https"
)

// NewProtocol creates a new Protocol instance and validates it
func NewProtocol(rawProtocol string) (Protocol, error) {
	var gp Protocol = Protocol(rawProtocol)
	return gp, gp.IsValid()
}

// IsValid checks the validity of the Protocol
func (gp Protocol) IsValid() error {
	switch gp {
	case SSH, HTTPS:
		return nil
	}
	return errors.New("Error: Invalid Protocol type")
}

// String converts Protocol to string representation
func (gp Protocol) String() string {
	return string(gp)
}

// AllowedVariables returns a formatted string on how to set this type with the set command
func (gp Protocol) AllowedVariables() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Protocol            { %v, %v }\n", SSH, HTTPS)
	return buf.String()
}
