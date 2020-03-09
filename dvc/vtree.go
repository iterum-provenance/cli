package dvc

import (
	"fmt"

	"github.com/Mantsje/iterum-cli/util"
)

// vtreeNode is a node in the VTree, each node represents 1 commit
type vtreeNode struct {
	Name     string
	Branch   hash
	Children []hash
	Parent   hash
}

// VTree holds a version tree for data versioning
type VTree map[hash]vtreeNode

// NewVTree instantiates a new version tree and sets the root node
func NewVTree() VTree {
	v := make(VTree)
	node := vtreeNode{
		Name:     "root",
		Branch:   "master",
		Children: []hash{},
		Parent:   hash("root"),
	}
	v[hash("root")] = node
	return v
}

// Add appends a new (leaf) commit to the tree, returns an error on failure
func (v VTree) Add(c Commit) error {
	node := vtreeNode{
		Name:     c.Name,
		Branch:   c.Branch,
		Children: []hash{},
		Parent:   c.Parent,
	}
	if _, ok := v[c.Hash]; ok {
		return fmt.Errorf("Error: Commit hash clash in VTree: %v, could not add commit", c.Hash)
	}
	v[c.Hash] = node
	return nil
}

// WriteToFolder writes the tree to the specified folder.
// Name of file is and should be determined by the tree structure
func (v VTree) WriteToFolder(folderPath string) error {
	fullPath := folderPath + "/history.vtree"
	return util.WriteJSONFile(fullPath, v)
}

// ParseFromFile tries to parse a history.vtree file
func (v *VTree) ParseFromFile(filepath string) error {
	if err := util.ReadJSONFile(filepath, v); err != nil {
		return fmt.Errorf("Error: Could not parse VTree due to `%v`", err)
	}
	return nil
}
