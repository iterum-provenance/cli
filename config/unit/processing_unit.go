package unit

// ProcessingUnitConf contains the config for a processing unit folder in an iterum project
type ProcessingUnitConf struct {
	UnitConf
}

// NewProcessingUnitConf instantiates a new ProcessingUnitConf and sets up defaults
func NewProcessingUnitConf(name string) ProcessingUnitConf {
	// Cannot use full ProcessingUnitConf{ field: val } syntax due to type composition in Go
	var puc = ProcessingUnitConf{NewUnitConf(name)}
	puc.UnitType = ProcessingUnit
	return puc
}
