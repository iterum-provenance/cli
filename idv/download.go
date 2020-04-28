package idv

import (
	"regexp"

	"github.com/iterum-provenance/iterum-go/util"
)

// Download data from this repository onto this local machine
func Download(selector *regexp.Regexp) (err error) {
	defer util.ReturnErrOnPanic(&err)()
	EnsureByPanic(EnsureSetup, "")

	return
}
