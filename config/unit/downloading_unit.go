package config

// DownloadingUnitConf contains the config for a downloading unit folder in an iterum project
type DownloadingUnitConf struct {
	UnitConf
}

// NewDownloadingUnitConf instantiates a new DownloadingUnitConf and sets up defaults
func NewDownloadingUnitConf(name string) DownloadingUnitConf {
	// Cannot use full DownloadingUnitConf{ field: val } syntax due to type composition in Go
	var duc = DownloadingUnitConf{NewUnitConf(name)}
	duc.unitType = DownloadingUnit
	return duc
}
