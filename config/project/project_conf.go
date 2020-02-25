package config

import (
	"github.com/Mantsje/iterum-cli/config"
)

// ProjectConf contains the config for the root folder of an iterum project
type ProjectConf struct {
	git config.GitConf
}
