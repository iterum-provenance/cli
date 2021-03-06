// Package credentials contains the different target credentials structures that are linked to the supported storage backends of IDV
package credentials

import (
	"errors"

	"github.com/iterum-provenance/cli/idv/ctl/storage"
	"github.com/iterum-provenance/cli/util"
)

// Storage is de general interface that holds
// the required credentials for connecting to a storage backend
// These are different for each platform
type Storage interface {
	util.Validatable
	GetLocation() string
}

// Parse tries to parse an interface into a credential type based on the backend
func Parse(raw map[string]interface{}, backend storage.Backend) (Storage, error) {
	switch backend {
	case storage.Local:
		return ParseLocal(raw)
	case storage.AmazonS3:
		return ParseAmazonS3(raw)
	case storage.CloudStore:
		return ParseCloudStore(raw)
	}
	return nil, errors.New("Error: No valid backend passed to storage.Parse")
}
