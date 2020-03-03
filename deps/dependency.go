package deps

import (
	"errors"
	"log"
	"os/exec"
)

// Dependencies holds all the Iterum dependency structs in 1 neat slice
var Dependencies = []Dep{}

// Dep contains all information about a dependency
type Dep struct {
	Cmd  string
	Name string
}

// IsValid checks the validity of the Dependency
func (d Dep) IsValid() error {
	for _, r := range Dependencies {
		if r == d {
			return nil
		}
	}
	return errors.New("Error: Invalid Dep")
}

// IsUsable checks whether the go program can use this dep using os.exec
func (d Dep) IsUsable() bool {
	_, err := exec.LookPath(d.Cmd)
	return err == nil
}

// EnsureDep checks whether a given dependency is met
func EnsureDep(dep Dep) {
	if !dep.IsUsable() {
		log.Fatal(dep.Name, " dependency is not accessible to iterum, run `iterum check` to verify")
	}
}

// Register updates the Dependencies slice with the passed Deps
func Register(deps ...Dep) {
	for _, d := range deps {
		Dependencies = append(Dependencies, d)
	}
}
