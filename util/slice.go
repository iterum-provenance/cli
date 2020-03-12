package util

import "regexp"

// Filter removes all matching elements from a string slice based on a regexp selector
func Filter(selector *regexp.Regexp, slice []string) (out []string, filtered int) {
	for _, elem := range slice {
		if selector.MatchString(elem) {
			filtered++
		} else {
			out = append(out, elem)
		}
	}
	return
}
