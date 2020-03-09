package credentials

import (
	"errors"
)

// AmazonS3 holds credentials for accessing AmazonS3 service
type AmazonS3 struct {
}

// IsValid checks the validity of this structure
func (c AmazonS3) IsValid() error {
	return errors.New("AmazonS3.IsValid N.I")
}

// ParseAmazonS3 tries to parse an interface as this credential storage
func ParseAmazonS3(raw map[interface{}]interface{}) (AmazonS3, error) {
	return AmazonS3{}, errors.New("ParseAmazonS3 N.I")
}
