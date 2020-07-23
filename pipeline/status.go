package pipeline

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"

	"github.com/iterum-provenance/cli/util"
)

// StatusAll retrieves the global status of each pipeline known to the manager
func StatusAll(managerURL *url.URL) (err error) {
	defer util.ReturnErrOnPanic(&err)()

	// Set target endpoint
	managerURL.Path = path.Join(managerURL.Path, "pipelines")

	var jsonResponse interface{}
	err = getJSON(managerURL, &jsonResponse)
	util.PanicIfErr(err, "")

	data, err := json.MarshalIndent(jsonResponse, "", "  ")
	util.PanicIfErr(err, fmt.Sprintf("Response generated an error: %v", err))
	fmt.Println(string(data))
	return nil
}

// Status retrieves the status of a specific pipeline known to the manager
func Status(phash string, managerURL *url.URL) (err error) {
	defer util.ReturnErrOnPanic(&err)()

	// Set target endpoint
	managerURL.Path = path.Join(managerURL.Path, "pipelines", phash, "status")

	var jsonResponse interface{}
	err = getJSON(managerURL, &jsonResponse)
	util.PanicIfErr(err, "")

	data, err := json.MarshalIndent(jsonResponse, "", "  ")
	util.PanicIfErr(err, fmt.Sprintf("Response generated an error: %v", err))
	fmt.Println(string(data))
	return nil
}

// Describe prompts the manager for the deployment specification of a certain pipeline
func Describe(phash string, managerURL *url.URL) (err error) {
	defer util.ReturnErrOnPanic(&err)()

	// Set target endpoint
	managerURL.Path = path.Join(managerURL.Path, "pipelines", phash)

	var jsonResponse interface{}
	err = getJSON(managerURL, &jsonResponse)
	util.PanicIfErr(err, "")

	data, err := json.MarshalIndent(jsonResponse, "", "  ")
	util.PanicIfErr(err, fmt.Sprintf("Response generated an error: %v", err))
	fmt.Println(string(data))
	return nil
}
