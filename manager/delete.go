package manager

import (
	"net/url"
	"path"

	"github.com/iterum-provenance/cli/util"
)

// DeleteAllPipelines kill all left over jobs and completed pipeline elements
func DeleteAllPipelines(managerURL *url.URL) (err error) {
	defer util.ReturnErrOnPanic(&err)()

	// Set target endpoint
	managerURL.Path = path.Join(managerURL.Path, "pipelines")

	err = delete(managerURL)
	util.PanicIfErr(err, "")

	return nil
}

// DeletePipeline kills and removes all parts of a pipeline, except for MinIO and RabbitMQ
func DeletePipeline(phash string, managerURL *url.URL) (err error) {
	defer util.ReturnErrOnPanic(&err)()

	// Set target endpoint
	managerURL.Path = path.Join(managerURL.Path, "pipelines", phash)

	err = delete(managerURL)
	util.PanicIfErr(err, "")

	return nil
}
