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
