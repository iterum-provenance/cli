package cmd

import (
	"log"

	"github.com/Mantsje/iterum-cli/config"
	"github.com/Mantsje/iterum-cli/config/parser"
	"github.com/Mantsje/iterum-cli/config/project"
	"github.com/Mantsje/iterum-cli/git"
	"github.com/Mantsje/iterum-cli/util"
)

// Package for shared functions specifically related to the CLI functionality

// Make sure we are in an iterum project root
func ensureRootLocation() (project.ProjectConf, error) {
	conf, repo, err := ensureIterumComponent()
	if err != nil {
		return project.ProjectConf{}, err
	}
	if repo != config.Project {
		return project.ProjectConf{}, errNotRoot
	}
	return conf.(project.ProjectConf), nil
}

// Make sure we are in an iterum component folder
func ensureIterumComponent() (interface{}, config.RepoType, error) {
	conf, repo, err := parser.ParseConfigFile(config.ConfigFileName)
	if err != nil {
		return conf, repo, errNoProject
	}
	return conf, repo, nil
}

// Possibly move this to `git` package
func initVersionTracking(conf config.Configurable) {
	base := conf.GetBaseConf()
	path := "./" + base.Name
	commitMsg := "Creation of Iterum " + base.RepoType.String() + " `" + base.Name + "`"

	if !NoRemote {
		uri := git.CreateRepo(commitMsg, base.Git.Platform, path)
		base.Git.URI = uri
		err := util.JSONWriteFile(base.Name+"/"+config.ConfigFileName, conf)
		if err != nil {
			log.Fatal(errConfigWriteFailed)
		}
	} else {
		git.CreateRepo(commitMsg, git.None, path)
	}
}
