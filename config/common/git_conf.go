package config

// GitConf contains all git-related configuration settings for units, flows and projects
type GitConf struct {
	platform GitPlatform
	protocol GitProtocol
}
