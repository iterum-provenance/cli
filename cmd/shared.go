package cmd

import (
	"log"
	"os"

	"github.com/iterum-provenance/cli/config"
	"github.com/iterum-provenance/cli/config/parser"
	"github.com/iterum-provenance/cli/consts"
	"github.com/iterum-provenance/cli/git"
	"github.com/iterum-provenance/cli/util"
)

// Package for shared functions specifically related to the base CLI functionality

// Make sure we are in an iterum component folder
func ensureIterumComponent(relPath string) (interface{}, config.RepoType, error) {
	conf, repo, err := parser.ParseConfigFile("./" + relPath + "/" + consts.ConfigFilePath)
	if err != nil {
		return conf, repo, errNoComponent
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
		err := util.WriteJSONFile(base.Name+"/"+consts.ConfigFilePath, conf)
		if err != nil {
			log.Fatal(errConfigWriteFailed)
		}
	} else {
		git.CreateRepo(commitMsg, git.None, path)
	}
}

// Creates the necessary folders for any iterum component: ./name and ./name/.iterum
func createComponentFolder(name string) {
	path := name + "/" + consts.ConfigFolder
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
	}
}
