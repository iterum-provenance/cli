package credentials

import (
	"errors"
)

// AmazonS3 holds credentials for accessing AmazonS3 service
type AmazonS3 struct {
}

// IsValid checks the validity of this structure
func (a AmazonS3) IsValid() error {
	return errors.New("AmazonS3.IsValid N.I")
}

// GetLocation returns a string path or url to where the data is located based on the backend
func (a AmazonS3) GetLocation() string {
	return ""
}

// ParseAmazonS3 tries to parse an interface as this credential storage
func ParseAmazonS3(raw map[string]interface{}) (AmazonS3, error) {
	return AmazonS3{}, errors.New("ParseAmazonS3 N.I")
}
