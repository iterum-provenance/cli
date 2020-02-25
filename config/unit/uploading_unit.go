package config

// UploadingUnitConf contains the config for a Uploading unit folder in an iterum project
type UploadingUnitConf struct {
	UnitConf
}

// NewUploadingUnitConf instantiates a new UploadingUnitConf and sets up defaults
func NewUploadingUnitConf(name string) UploadingUnitConf {
	// Cannot use full UploadingUnitConf{ field: val } syntax due to type composition in Go
	var uuc = UploadingUnitConf{NewUnitConf(name)}
	uuc.UnitType = UploadingUnit
	return uuc
}
