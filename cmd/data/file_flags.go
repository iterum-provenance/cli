package data

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/Mantsje/iterum-cli/util"
)

// Contains functionality to handle the file specifier related flags of the `iterum data` subcommands

// The list of flags:

// Recursive states whether `iterum data add/rm` should recurse into folders found in the specified arguments
var Recursive bool

// ShowExcluded states whether the list of excluded files should be shown
var ShowExcluded bool

// Exclusions holds files and folders to be excluded from adding/removing, etc
var Exclusions []string

// ------------------------ Gathering Files ------------------------ //

// Cleans up a possibly incomplete path to a complete and absolute one,
// crashes the program if file does not exist
func cleanPath(path string) string {
	if !util.IsFileOrDir(path) {
		log.Fatal(fmt.Errorf("Error: `%v` is not a valid path", path))
	}
	path, _ = filepath.Abs(path)
	return filepath.Clean(path)
}

// filesInDir returns the list of files found in dirPath
// expects it to be a dir, recursive denotes to recurse into the folder
func filesInDir(dirPath string, recurse bool) (files []string) {
	contents, _ := ioutil.ReadDir(dirPath)
	for _, info := range contents {
		if info.IsDir() {
			if recurse {
				files = append(files, filesInDir(dirPath+"/"+info.Name(), recurse)...)
			}
		} else { // Only add files
			files = append(files, cleanPath(dirPath+"/"+info.Name()))
		}
	}
	return files
}

func getAllFiles(filesOrDirs []string) []string {
	files := []string{}
	for _, path := range filesOrDirs { // iterate over passed args
		path = cleanPath(path)
		info, _ := os.Stat(path)
		if info.IsDir() {
			files = append(files, filesInDir(path, Recursive)...)
		} else { // Just add the file
			files = append(files, path)
		}
	}
	return files
}

// ------------------------ Filtering Files ------------------------ //

// filterFilesWith uses a selector to filter a set of files from a list of files and returns both the filtered
// and excluded slice of filepaths.
func filterFilesWith(selector *regexp.Regexp, files []string) (filtered, excluded []string) {
	filtered = []string{}
	excluded = []string{}

	for _, p := range files {
		if !selector.MatchString(p) {
			filtered = append(filtered, p)
		} else if ShowExcluded {
			excluded = append(excluded, p)
		}
	}
	return
}

// printExclusions is the function used to resolve the ShowExluded flag, it will print non-included files
func printExclusions(selector *regexp.Regexp, exclusions []string) {
	fmt.Println("Files excluded using regexp:")
	fmt.Printf("\t%v\n", selector.String())
	fmt.Println("Excluded files")
	fmt.Println("{")
	for _, file := range exclusions {
		fmt.Printf("\t'%v'\n", file)
	}
	fmt.Println("}")
}

// exclude resolves the Exclusions and ShowExcluded flags that `data rm/add` share
func exclude(files []string) []string {
	selector := ""
	first := true
	for _, excl := range Exclusions { // Build a valid regular expression from the passed args
		if isValidSelector(excl) {
			if !first {
				selector += "|" + "(" + excl + ")"
			} else {
				selector += "(" + excl + ")"
			}
		}
		first = false
	}
	if selector == "" { // Don't filter if no selector, because all files match empty expression
		return files
	}
	exclutor, _ := regexp.Compile(selector)
	files, excluded := filterFilesWith(exclutor, files)

	if ShowExcluded {
		printExclusions(exclutor, excluded)
	}

	return files
}
