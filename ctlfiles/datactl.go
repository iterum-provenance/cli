package ctlfiles

import (
	"errors"

	"github.com/Mantsje/iterum-cli/ctlfiles/storage"
	"github.com/Mantsje/iterum-cli/ctlfiles/storage/credentials"
	"github.com/Mantsje/iterum-cli/util"
)

// DataCTL is the structure that a data config files can be parsed into (.icf) files
type DataCTL struct {
	Name        string
	Backend     storage.Backend
	Credentials credentials.Storage
}

// IsValid checks the validity of the struct
func (d DataCTL) IsValid() error {
	err := util.Verify(d.Backend, nil)
	err = util.Verify(d.Credentials, err)
	return err
}

// ParseFromFile tries to parse a .ifc file written in yaml
func (d *DataCTL) ParseFromFile(filepath string) error {
	// Needed because no cannot parse into interface of DataCTL directly
	var raw struct {
		Name        string
		Backend     storage.Backend
		Credentials interface{}
	}
	if err := util.ReadYAMLFile(filepath, &raw); err != nil {
		return errors.New("Error: Could not parse yaml file")
	}
	errBackend := raw.Backend.IsValid()
	c := raw.Credentials.(map[interface{}]interface{})
	creds, errCred := credentials.Parse(c, raw.Backend)
	d.Name = raw.Name
	d.Backend = raw.Backend
	d.Credentials = creds

	return util.ReturnFirstErr(errBackend, errCred, d.IsValid())
}
