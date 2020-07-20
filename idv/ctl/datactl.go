package ctl

import (
	"errors"
	"regexp"

	"github.com/iterum-provenance/cli/idv/ctl/storage"
	"github.com/iterum-provenance/cli/idv/ctl/storage/credentials"
	"github.com/iterum-provenance/cli/util"
)

const (
	// DefaultDaemonURL is the default url tried when none is specified
	DefaultDaemonURL = "http://localhost:3000/"
)

// DataCTL is the structure that a data config files can be parsed into (.icf) files
type DataCTL struct {
	Name        string              `yaml:"name" json:"name"`
	Description string              `yaml:"description" json:"description"`
	DaemonURL   string              `yaml:"daemon" json:"daemon"`
	Backend     storage.Backend     `yaml:"backend" json:"backend"`
	Credentials credentials.Storage `yaml:"credentials" json:"credentials"`
}

// rawDataCTL is supposed to only be used in case of parsing a DataCTL, pass this to ParseFromRaw
type rawDataCTL struct {
	Name        string                 `yaml:"name" json:"name"`
	Description string                 `yaml:"description" json:"description"`
	DaemonURL   string                 `yaml:"daemon" json:"daemon"`
	Backend     storage.Backend        `yaml:"backend" json:"backend"`
	Credentials map[string]interface{} `yaml:"credentials" json:"credentials"`
}

// NewDataCTL creates a new instance of DataCTL, used for creating new data repos
func NewDataCTL(name, description string, backend storage.Backend) DataCTL {
	var creds credentials.Storage
	switch backend {
	case storage.Local:
		creds = credentials.NewLocal()
	case storage.AmazonS3:
		creds = credentials.NewAmazonS3()
	case storage.CloudStore:
		creds = credentials.NewCloudStore()
	}
	return DataCTL{
		Name:        name,
		Description: description,
		Backend:     backend,
		DaemonURL:   DefaultDaemonURL,
		Credentials: creds,
	}
}

// IsValid checks the validity of the struct
func (d DataCTL) IsValid() error {
	rexp, _ := regexp.Compile("[ \t\n\r]")
	if rexp.ReplaceAllString(d.Name, "") != d.Name {
		return errors.New("Error: Name contains whitespace, use '-' instead")
	}
	if d.DaemonURL == "" {
		return errors.New("Error: No valid DaemonURL found")
	}
	err := util.Verify(d.Backend, nil)
	err = util.Verify(d.Credentials, err)
	return err
}

// parseFromRaw takes an interface and tries to parse it as a DataCTL
func (d *DataCTL) parseFromRaw(raw rawDataCTL) error {
	errBackend := raw.Backend.IsValid()
	creds, errCred := credentials.Parse(raw.Credentials, raw.Backend)
	d.Name = raw.Name
	d.Backend = raw.Backend
	d.DaemonURL = raw.DaemonURL
	d.Credentials = creds

	return util.ReturnFirstErr(errBackend, errCred, d.IsValid())
}

// ParseFromMap takes an interface->interface map and tries to parse it as DataCTL
func (d *DataCTL) ParseFromMap(m map[string]interface{}) error {
	var raw rawDataCTL
	raw.Name = m["name"].(string)
	raw.Description = m["description"].(string)
	if _, ok := m["daemon"]; ok {
		raw.DaemonURL = m["daemon"].(string)
	} else {
		raw.DaemonURL = DefaultDaemonURL
	}
	raw.Backend = storage.Backend(m["backend"].(string))
	raw.Credentials = m["credentials"].(map[string]interface{})
	return d.parseFromRaw(raw)
}

// ParseFromFile tries to parse a .ifc file written in yaml
func (d *DataCTL) ParseFromFile(filepath string) error {
	var raw rawDataCTL
	if err := util.ReadYAMLFile(filepath, &raw); err != nil {
		return errors.New("Error: Could not parse yaml file")
	}
	if raw.DaemonURL == "" {
		raw.DaemonURL = DefaultDaemonURL
	}
	return d.parseFromRaw(raw)
}

// WriteToFile writes the datactl to a YAML file
func (d DataCTL) WriteToFile(filepath string) error {
	return util.WriteYAMLFile(filepath, d)
}

// GetStorageLocation returns a string path or URL to where data is located, based on the backend
func (d DataCTL) GetStorageLocation() string {
	return d.Credentials.GetLocation()
}

// ToReport returns a string version that is presentable to a user descbring this ctl
func (d DataCTL) ToReport() string {
	report := ""
	report += "{\n"
	report += "\tName: " + d.Name + "\n"
	report += "\tDescription: " + d.Description + "\n"
	report += "\tDaemon: " + d.DaemonURL + "\n"
	report += "\tBackend: " + d.Backend.String() + "\n"
	report += "\tLocation: " + d.GetStorageLocation() + "\n"
	report += "}\n"
	return report
}
