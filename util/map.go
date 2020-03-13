package util

import "errors"

// MapContains returns whether a map contains a certain element
// asValue denotes whether to check for values (if false will search for key)
func MapContains(m map[interface{}]interface{}, target interface{}, asValue bool) bool {
	if asValue {
		for _, val := range m {
			if val == target {
				return true
			}
		}
	} else {
		_, ok := m[target]
		return ok
	}
	return false
}

// GetKeyByValue returns the key of a certain value of a map. It returns the first instance it finds that
// has the given value. It returns an err if the value is not in the map
func GetKeyByValue(m map[interface{}]interface{}, target interface{}) (interface{}, error) {
	for key, val := range m {
		if val == target {
			return key, nil
		}
	}
	return nil, errors.New("Error: Value was not found in passed map")
}
