package idv

import "regexp"

// Download data from this repository onto this local machine
func Download(selector *regexp.Regexp) (err error) {
	defer _returnErrOnPanic(&err)()
	EnsureByPanic(EnsureSetup, "")

	return
}
