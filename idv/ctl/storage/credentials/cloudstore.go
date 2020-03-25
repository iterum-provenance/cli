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

// GetLocation returns a string path or url to where the data is located based on the backend
func (c CloudStore) GetLocation() string {
	return ""
}

// ParseCloudStore tries to parse an interface as this credential storage
func ParseCloudStore(raw map[string]interface{}) (CloudStore, error) {
	return CloudStore{}, errors.New("ParseCloudStore not Implemented")
}
