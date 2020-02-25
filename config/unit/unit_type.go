package config

import "errors"

// UnitType is used for defining what type of unit this module contains
type UnitType string

// Enum-like values allowed for UnitType type
const (
	DownloadingUnit UnitType = "downloading"
	ProcessingUnit  UnitType = "processing"
	UploadingUnit   UnitType = "uploading"
)

// NewUnitType creates a new UnitType instance and validates it
func NewUnitType(unitType string) (UnitType, error) {
	var ut UnitType = UnitType(unitType)
	return ut, ut.IsValid()
}

// IsValid checks the validity of the UnitType
func (ut UnitType) IsValid() error {
	switch ut {
	case DownloadingUnit, ProcessingUnit, UploadingUnit:
		return nil
	}
	return errors.New("Error: Invalid UnitType")
}
