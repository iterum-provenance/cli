package config

// Configurable structs are those that have an instance of Conf embedded
type Configurable interface {
	GetBaseConf() *Conf
}

// Conf is the base struct for configuration files
type Conf struct {
	Name     string
	RepoType RepoType
	Git      GitConf
}

// NewConf initializes a new Conf instance with the passed args (used by child inits)
func NewConf(name string, repo RepoType) Conf {
	return Conf{
		Name:     name,
		RepoType: repo,
		Git:      NewGitConf(),
	}
}

// Set sets a field in this conf based on a string, rather than knowing the exact type
func (c *Conf) Set(variable []string, value interface{}) error {
	return SetField(c, variable, value)
}

// GetBaseConf returns the embedded Conf
func (c *Conf) GetBaseConf() *Conf {
	return c
}
