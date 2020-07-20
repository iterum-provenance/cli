package manager

import (
	"net/url"
	"path"

	"github.com/iterum-provenance/cli/util"
)

// TerminateAllPipelines kill all left over jobs and completed pipeline elements
func TerminateAllPipelines(managerURL *url.URL) (err error) {
	defer util.ReturnErrOnPanic(&err)()

	// Set target endpoint
	managerURL.Path = path.Join(managerURL.Path, "pipelines")

	err = delete(managerURL)
	util.PanicIfErr(err, "")

	return nil
}

// TerminatePipeline kills and removes all parts of a pipeline, except for MinIO and RabbitMQ
func TerminatePipeline(phash string, managerURL *url.URL) (err error) {
	defer util.ReturnErrOnPanic(&err)()

	// Set target endpoint
	managerURL.Path = path.Join(managerURL.Path, "pipelines", phash)

	err = delete(managerURL)
	util.PanicIfErr(err, "")

	return nil
}
