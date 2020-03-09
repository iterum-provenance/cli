package dvc

import (
	"fmt"

	"github.com/Mantsje/iterum-cli/util"
)

// Branch internally defines a data version commit file
type Branch struct {
	Name string
	HEAD hash
	Hash hash
}

// NewBranch creates and initializes a new branch
func NewBranch(name string) Branch {
	return Branch{
		Name: name,
		Hash: newHash(32),
		HEAD: "",
	}
}

// WriteToFolder writes the branch to the specified folder.
// Name of file is and should be determined by the branch structure
func (b Branch) WriteToFolder(folderPath string) error {
	fullPath := folderPath + "/" + b.Hash.String() + ".branch"
	return util.WriteJSONFile(fullPath, b)
}

// ParseFromFile tries to parse a .branch file
func (b *Branch) ParseFromFile(filepath string) error {
	if err := util.ReadJSONFile(filepath, b); err != nil {
		return fmt.Errorf("Error: Could not parse Branch due to `%v`", err)
	}
	return nil
}
