package config

import "net/url"

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
