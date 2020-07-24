package pipeline

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/iterum-provenance/cli/util"
	"github.com/prometheus/common/log"
)

// Download retrieves a file from the result set of a pipeline run
func Download(phash, filename, folder string, daemonURL *url.URL) (err error) {
	defer util.ReturnErrOnPanic(&err)()

	// Set target endpoint
	daemonURL.Path = path.Join(daemonURL.Path, "pipelines", phash, "results", filename)

	resp, err := http.Get(daemonURL.String())
	util.PanicIfErr(err, "")

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("Download failed and returned statuscode %v", resp.StatusCode)
	}

	filepath := path.Join(folder, filename)
	log.Infof("Flushing file to '%v'", filepath)
	fhandle, err := os.Create(filepath)
	util.PanicIfErr(err, "")
	defer fhandle.Close()
	_, err = io.Copy(fhandle, resp.Body)
	util.PanicIfErr(err, "")

	return nil
}
