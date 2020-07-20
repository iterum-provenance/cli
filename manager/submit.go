package manager

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/gaarkeuken/go-lib/util"
)

// SubmitPipeline is the function that submits a json specification of a pipeline to the manager at managerURL
func SubmitPipeline(filepath string, managerURL *url.URL) (err error) {
	defer util.ReturnErrOnPanic(&err)()

	// Set target endpoint
	managerURL.Path = path.Join(managerURL.Path, "submit_pipeline_actor")

	// Open configuration file
	file, err := os.Open(filepath)
	util.PanicIfErr(err, "")

	resp, err := http.Post(managerURL.String(), "application/json", file)
	util.PanicIfErr(err, "")
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		return fmt.Errorf("Request returned status code %v", resp.StatusCode)
	}

	var jsonResponse interface{}
	data, err := ioutil.ReadAll(resp.Body)
	util.PanicIfErr(err, "")
	json.Unmarshal(data, &jsonResponse)
	fmt.Println(jsonResponse)

	return nil
}
