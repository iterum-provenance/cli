package data

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/iterum-provenance/cli/util"
	"github.com/prometheus/common/log"
)

// ------------------------ Gathering Files ------------------------ //
// Some additional logic for handling selection of files on the user's machines.
// These functions allow the other packages to pose assumptions about file existence and such

// cleanPath cleans up a possibly incomplete path to a complete and absolute one,
// crashes the program if file does not exist
func cleanPath(path string) string {
	if !util.IsFileOrDir(path) {
		log.Fatalln(fmt.Errorf("Error: `%v` is not a valid path", path))
	}
	path, _ = filepath.Abs(path)
	return filepath.Clean(path)
}

// filesInDir returns the list of files found in dirPath
// expects it to be a validated dir, recursive denotes whether to recurse into the folder
func filesInDir(dirPath string, recurse bool) (files []string) {
	contents, _ := ioutil.ReadDir(dirPath)
	for _, info := range contents {
		if info.IsDir() {
			if recurse {
				files = append(files, filesInDir(path.Join(dirPath, info.Name()), recurse)...)
			}
		} else { // Only add files
			files = append(files, cleanPath(path.Join(dirPath, info.Name())))
		}
	}
	return files
}

// getAllFiles returns all files contained in the passed array of folders or directories
func getAllFiles(filesOrDirs []string, recurse bool) []string {
	files := []string{}
	for _, path := range filesOrDirs { // iterate over passed args
		path = cleanPath(path)
		info, _ := os.Stat(path)
		if info.IsDir() {
			files = append(files, filesInDir(path, recurse)...)
		} else { // Just add the file
			files = append(files, path)
		}
	}
	return files
}

// getPaths returns all existing file paths as paths and non-existing files as names
// useful for removing files from commits that are not on the user's machine
func getPaths(args []string) (paths, names []string) {
	for _, arg := range args {
		if isValidLocation(arg) {
			paths = append(paths, arg)
		} else {
			names = append(names, arg)
		}
	}
	return
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
		} else {
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

// buildSelector construct a regexp from a set of selectors,
// if regexp cannot parse it, it is simply omitted
func buildSelector(selectors []string) (r *regexp.Regexp) {
	selector := ""
	first := true
	for _, sel := range selectors { // Build a valid regular expression from the passed args
		if isValidSelector(sel) {
			if !first {
				selector += "|" + "(" + sel + ")"
			} else {
				selector += "(" + sel + ")"
			}
		}
		first = false
	}
	r, _ = regexp.Compile(selector)
	return
}

// exclude resolves the Exclusions and ShowExcluded flags that `data rm/add` share
func exclude(files, exclutors []string, printExcluded bool) []string {
	// Don't filter if no selector, because all files match empty expression
	if len(exclutors) == 0 {
		return files
	}
	exclutor := buildSelector(exclutors)
	files, excluded := filterFilesWith(exclutor, files)

	if printExcluded {
		printExclusions(exclutor, excluded)
	}

	return files
}
