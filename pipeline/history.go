package pipeline

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"

	"github.com/iterum-provenance/cli/util"
)

// History retrieves the global information of each pipeline known to the daemon
func History(daemonURL *url.URL) (err error) {
	defer util.ReturnErrOnPanic(&err)()

	// Set target endpoint
	daemonURL.Path = path.Join(daemonURL.Path, "pipelines")

	var jsonResponse interface{}
	err = getJSON(daemonURL, &jsonResponse)
	util.PanicIfErr(err, "")

	data, err := json.MarshalIndent(jsonResponse, "", "  ")
	util.PanicIfErr(err, fmt.Sprintf("Response generated an error: %v", err))
	fmt.Println(string(data))
	return nil
}
