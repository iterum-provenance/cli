package data

import (
	"regexp"

	"github.com/Mantsje/iterum-cli/util"
)

// Checks whether a passed selector is parseable as regex
func isValidSelector(arg string) bool {
	_, err := regexp.Compile(arg)
	if err != nil {
		return false
	}
	return true
}

// Checks whether a passed file/folder arg exists
func isValidLocation(arg string) bool {
	return util.IsFolderOrDir(arg)
}

// ParseCurrentIDVState reads in the currently stored dvc information as a general setup
// func ParseCurrentIDVState() (idv.VTree, idv.Branch, idv.Config, error) {
// var history idv.VTree
// var branch idv.Branch
// var config idv.Config
// var local idv.Commit
// history.ParseFromFile(idv.VTreeFile)
// branch.ParseFromFile(idv.CurrentBranch)
// local.ParseFromFile(idv.CurrentCommit)
// idv.LOCAL
// config.ParseFromFile(constants.IDVConfigFileName)
// }
