package config

// GitConf contains all git-related configuration settings for units, flows and projects
type GitConf struct {
	platform GitPlatform
	protocol GitProtocol
}

// NewGitConf creates a new and fully empty instance of GitConf
func NewGitConf() GitConf {
	return GitConf{}
}

// IsValid checks the validity of the GitConf
func (gc GitConf) IsValid() error {
	err := Verify(gc.platform, nil)
	err = Verify(gc.protocol, err)
	return err
}
