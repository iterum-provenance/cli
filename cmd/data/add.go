package data

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.PersistentFlags().BoolVarP(&Recursive, "recursive", "r", false, "Descend recursively into passed folders")
	addCmd.PersistentFlags().StringSliceVarP(&Exclusions, "exclude", "x", []string{}, "Exclude files and folders -x selector1 -x selector2,selector3")
	addCmd.PersistentFlags().BoolVarP(&ShowExcluded, "show-excluded", "s", false, "Show list of excluded files")
}

var addCmd = &cobra.Command{
	Use:   "add [file/folder]...",
	Short: "Add or Update files to the current commit",
	Long:  `Stages files to be added to the dataset, or in case of name clash to update those files`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errNotEnoughArgs
		}
		var invalids []string
		for _, arg := range args {
			if !isValidLocation(arg) {
				invalids = append(invalids, arg)
			}
		}
		if len(invalids) != 0 {
			return errInvalidArgs(invalids...)
		}
		return nil
	},
	Run: addRun,
}

func getAllFiles(filesOrDirs []string) []string {
	files := []string{}
	for _, path := range filesOrDirs { // iterate over passed args
		info, _ := os.Stat(path)
		if info.IsDir() {
			if Recursive { // Descend into dir
				err := filepath.Walk(path,
					func(path string, info os.FileInfo, err error) error {
						if err != nil {
							return err
						}
						if !info.IsDir() {
							files = append(files, path)
						}
						return nil
					})
				if err != nil {
					log.Fatal(err)
				}
			} else { // Add only files in this folder, don't recurse
				contents, _ := ioutil.ReadDir(path)
				for _, info := range contents {
					if !info.IsDir() { // Exclude directories, only add files
						files = append(files, info.Name())
					}
				}
			}
		} else { // Just add the file
			files = append(files, path)
		}
	}
	return files
}

func exclude(files []string) []string {
	selector := ""
	first := true
	for _, excl := range Exclusions {
		if isValidSelector(excl) {
			if !first {
				selector += "|" + "(" + excl + ")"
			} else {
				selector += "(" + excl + ")"
			}
		}
		first = false
	}
	exclutor, _ := regexp.Compile(selector)
	whitelisted := []string{}

	if ShowExcluded {
		fmt.Println("Files excluded using regexp:")
		fmt.Printf("\t%v\n", exclutor.String())
		fmt.Println("Excluded files")
		fmt.Println("{")
	}
	for _, p := range files {
		if exclutor.MatchString(p) {
			if ShowExcluded {
				fmt.Printf("\t'%v'\n", p)
			} else {
				whitelisted = append(whitelisted, os.Getenv("PWD")+p)
			}
		}
	}
	if ShowExcluded {
		fmt.Println("}")
	}
	return whitelisted
}

func addRun(cmd *cobra.Command, args []string) {
	log.Println("`iterum data add`")
	allFiles := getAllFiles(args)
	whitelisted := exclude(allFiles)
	fmt.Println(whitelisted)
	// idv.Add(whitelisted)
}
