package prompter

import (
	"github.com/iterum-provenance/cli/idv/ctl/storage"
	"github.com/manifoldco/promptui"
)

// StorageType lets the user pick a storage backend
func StorageType() string {
	return pick(promptui.Select{
		Label: "What kind of storage backend should this data set use",
		Items: []storage.Backend{
			storage.Local,
			storage.AmazonS3,
			storage.CloudStore,
		},
	})
}
