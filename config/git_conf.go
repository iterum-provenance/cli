package config

import (
	"bytes"
	"fmt"
	"net/url"
)

// GitConf contains all git-related configuration settings for units, flows and projects
type GitConf struct {
	Platform GitPlatform
	Protocol GitProtocol
	URI      url.URL
}

// NewGitConf creates a new and fully empty instance of GitConf
func NewGitConf() GitConf {
	return GitConf{}
}

// IsValid checks the validity of the GitConf
func (gc GitConf) IsValid() error {
	err := Verify(gc.Platform, nil)
	err = Verify(gc.Protocol, err)
	return err
}

// Set sets a field in this conf based on a string, rather than knowing the exact type
func (gc *GitConf) Set(variable []string, value interface{}) error {
	return SetField(gc, variable, value)
}

// AllowedVariables returns a formatted string on how to set this type with the set command
func (gc GitConf) AllowedVariables() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Git\n")
	fmt.Fprintf(&buf, "    .%v\n", gc.Platform.AllowedVariables())
	fmt.Fprintf(&buf, "    .%v\n", gc.Protocol.AllowedVariables())
	return buf.String()
}
