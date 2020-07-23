package pipeline

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"

	"github.com/iterum-provenance/cli/util"
)

// SubmitPipeline is the function that submits a json specification of a pipeline to the manager at managerURL
func SubmitPipeline(filepath string, managerURL *url.URL) (err error) {
	defer util.ReturnErrOnPanic(&err)()

	// Set target endpoint
	managerURL.Path = path.Join(managerURL.Path, "pipelines")

	// Open configuration file
	file, err := os.Open(filepath)
	util.PanicIfErr(err, "")

	var jsonResponse interface{}
	err = postJSON(managerURL, file, &jsonResponse)
	util.PanicIfErr(err, "")

	data, err := json.MarshalIndent(jsonResponse, "", "  ")
	util.PanicIfErr(err, fmt.Sprintf("Response generated an error: %v", err))
	fmt.Println(string(data))
	return nil
}
