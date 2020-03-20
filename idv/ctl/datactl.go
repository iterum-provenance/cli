package ctl

import (
	"errors"
	"regexp"

	"github.com/Mantsje/iterum-cli/idv/ctl/storage"
	"github.com/Mantsje/iterum-cli/idv/ctl/storage/credentials"
	"github.com/Mantsje/iterum-cli/util"
)

// DataCTL is the structure that a data config files can be parsed into (.icf) files
type DataCTL struct {
	Name        string              `yaml:"name" json:"name"`
	Description string              `yaml:"description" json:"description"`
	Backend     storage.Backend     `yaml:"backend" json:"backend"`
	Credentials credentials.Storage `yaml:"credentials" json:"credentials"`
}

// IsValid checks the validity of the struct
func (d DataCTL) IsValid() error {
	rexp, _ := regexp.Compile("[ \t\n\r]")
	if rexp.ReplaceAllString(d.Name, "") != d.Name {
		return errors.New("Error: Name contains whitespace, use '-' instead")
	}
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

// WriteToFile writes the datactl to a YAML file
func (d DataCTL) WriteToFile(filepath string) error {
	return util.WriteYAMLFile(filepath, d)
}

// GetStorageLocation returns a string path or URL to where data is located, based on the backend
func (d DataCTL) GetStorageLocation() string {
	return d.Credentials.GetLocation()
}
