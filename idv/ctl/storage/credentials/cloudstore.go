package credentials

import (
	"errors"
)

// CloudStore holds credentials for Google Cloud Store
type CloudStore struct {
	Message string `yaml:"message" json:"message"` // Remove this once implementing
}

// NewCloudStore instantiates a clean empty CloudStore storage backend credentials struct
func NewCloudStore() CloudStore {
	return CloudStore{
		Message: "CloudStore storage NOT IMPLEMENTED YET",
	}
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
