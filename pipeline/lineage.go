package pipeline

import (
	"net/url"
	"path"

	"github.com/cheggaaa/pb/v3"

	"github.com/iterum-provenance/cli/util"
)

// Lineage retrieves all lineage information associated with the pipeline hash
func Lineage(phash, folder string, daemonURL *url.URL) (err error) {
	defer util.ReturnErrOnPanic(&err)()

	// Set global endpoint for all requests
	daemonURL.Path = path.Join(daemonURL.Path, "pipelines", phash, "lineage")

	var fragmentIds []string
	err = getJSON(daemonURL, &fragmentIds)
	util.PanicIfErr(err, "")

	bar := pb.StartNew(len(fragmentIds))
	daemonURL.Path = path.Join(daemonURL.Path, "toremoveuseingdir")
	for _, id := range fragmentIds {
		bar.Increment()
		var fragLineage interface{}
		daemonURL.Path = path.Join(path.Dir(daemonURL.Path), id)
		err = getJSON(daemonURL, &fragLineage)
		util.PanicIfErr(err, "")
		util.WriteJSONFile(path.Join(folder, id), fragLineage)
	}
	bar.Finish()

	return nil
}
