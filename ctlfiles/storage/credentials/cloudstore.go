package credentials

import (
	"errors"
)

// CloudStore holds credentials for Google Cloud Store
type CloudStore struct {
}

// IsValid checks the validity of this structure
func (c CloudStore) IsValid() error {
	return errors.New("CloudStore.IsValid N.I")
}

// ParseCloudStore tries to parse an interface as this credential storage
func ParseCloudStore(raw map[interface{}]interface{}) (CloudStore, error) {
	return CloudStore{}, errors.New("ParseCloudStore not Implemented")
}
