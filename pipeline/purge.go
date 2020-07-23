package pipeline

import (
	"net/url"
	"path"

	"github.com/iterum-provenance/cli/util"
)

// PurgePipeline kills and removes all parts of a pipeline, except for MinIO and RabbitMQ
func PurgePipeline(phash string, managerURL *url.URL) (err error) {
	defer util.ReturnErrOnPanic(&err)()

	// Set target endpoint
	managerURL.Path = path.Join(managerURL.Path, "pipelines", phash, "purge")

	err = delete(managerURL)
	util.PanicIfErr(err, "")

	return nil
}
