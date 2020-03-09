package credentials

import (
	"errors"

	"github.com/Mantsje/iterum-cli/ctlfiles/storage"
	"github.com/Mantsje/iterum-cli/util"
)

// Storage is de general interface that holds
// the required credentials for connecting to a storage backend
// These are different for each platform
type Storage interface {
	util.Validatable
}

// Parse tries to parse an interface into a credential type based on the backend
func Parse(raw map[interface{}]interface{}, backend storage.Backend) (Storage, error) {
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
