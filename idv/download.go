package idv

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"

	"github.com/cheggaaa/pb"
	"github.com/iterum-provenance/cli/idv/ctl"
	"github.com/iterum-provenance/cli/util"
	"github.com/manifoldco/promptui"
)

// Download data from this repository onto this local machine
func Download(selector *regexp.Regexp, folder string, show, prompt bool) (err error) {
	defer util.ReturnErrOnPanic(&err)()
	EnsureByPanic(EnsureIDVRepo, "")
	EnsureByPanic(EnsureSetup, "")
	EnsureByPanic(EnsureLOCAL, "")
	var local Commit
	parseLOCAL(&local)
	var ctl ctl.DataCTL
	parseConfig("idv-config.yaml", &ctl)

	// Match all files in the current commit
	matches := []string{}
	for _, file := range local.Files {
		if selector.MatchString(path.Dir(file)) {
			matches = append(matches, file)
		}
	}

	// Print files to be downloaded
	if show {
		for _, file := range matches {
			fmt.Println(path.Dir(file))
		}
		fmt.Printf("    Downloading %v files\n", len(matches))
	}

	// Prompt the user whether he wants to continue with this
	if prompt {
		prompt := promptui.Select{
			Label: fmt.Sprintf("Do you want to download %v files into '%v'?", len(matches), folder),
			Items: []string{"No", "Yes"},
		}
		_, answer, err := prompt.Run()
		util.PanicIfErr(err, "")
		if answer == "No" {
			fmt.Println("Cancelled")
			return nil
		}
	}

	// Actually download all the files
	bar := pb.StartNew(len(matches))
	defer bar.Finish()
	strDaemonURL := ctl.DaemonURL
	daemonURL, err := url.Parse(strDaemonURL)
	util.PanicIfErr(err, "")
	daemonURL.Path = path.Join(daemonURL.Path, ctl.Name, "file", "fakefile", "fakecommit")

	for _, file := range matches {
		bar.Increment()
		// update the url to the new one
		daemonURL.Path = path.Join(path.Dir(path.Dir(daemonURL.Path)), file)

		resp, err := http.Get(daemonURL.String())
		util.PanicIfErr(err, "")

		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			return fmt.Errorf("Download failed and returned statuscode %v", resp.StatusCode)
		}

		filename := path.Dir(file)
		chash := path.Base(file)
		filepath := path.Join(folder, chash+"_"+filename)
		fmt.Printf("Downloading %v to '%v'\n", file, filepath)
		fhandle, err := os.Create(filepath)
		util.PanicIfErr(err, "")
		_, err = io.Copy(fhandle, resp.Body)
		util.PanicIfErr(err, "")
		fhandle.Close()
	}

	return nil
}
