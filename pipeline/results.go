package pipeline

import (
	"fmt"
	"net/url"
	"path"

	"github.com/iterum-provenance/cli/util"
)

// Results retrieves the names of all files in the results of a pipeline execution
func Results(phash string, daemonURL *url.URL) (err error) {
	defer util.ReturnErrOnPanic(&err)()

	// Set target endpoint
	daemonURL.Path = path.Join(daemonURL.Path, "pipelines", phash, "results")

	var files []string
	err = getJSON(daemonURL, &files)
	util.PanicIfErr(err, "")

	for _, file := range files {
		fmt.Println(file)
	}
	return nil
}
