package manager

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
	fmt.Println(json.MarshalIndent(jsonResponse, "", "  "))

	return nil
}
