package idv

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

type commitTree map[hash]vtreeNode

// VTree holds a version tree for data versioning
type VTree struct {
	Tree     commitTree
	Branches map[hash]string
}

// NewVTree instantiates a new version tree and sets the root node
func NewVTree(root Commit, master Branch) VTree {
	v := VTree{make(commitTree), make(map[hash]string)}
	node := vtreeNode{
		Name:     root.Name,
		Branch:   master.Hash,
		Children: []hash{},
		Parent:   root.Hash,
	}
	v.Tree[root.Hash] = node
	v.Branches[master.Hash] = master.Name
	return v
}

// WriteToFolder writes the tree to the specified folder.
// Name of file is and should be determined by the tree structure
func (v VTree) WriteToFolder(folderPath string) error {
	fullPath := folderPath + "/" + vtreeFileName
	return util.WriteJSONFile(fullPath, v)
}

// ParseFromFile tries to parse a history.vtree file
func (v *VTree) ParseFromFile(filepath string) error {
	if err := util.ReadJSONFile(filepath, v); err != nil {
		return fmt.Errorf("Error: Could not parse VTree due to `%v`", err)
	}
	return nil
}

// Add appends a new (leaf) commit to the tree, returns an error on failure
func (v *VTree) Add(c Commit) error {
	node := vtreeNode{
		Name:     c.Name,
		Branch:   c.Branch,
		Children: []hash{},
		Parent:   c.Parent,
	}
	if _, ok := v.Tree[c.Hash]; ok {
		return fmt.Errorf("Error: Commit has clash in VTree: %v, could not add commit", c.Hash)
	}
	v.Tree[c.Hash] = node
	return nil
}

// BranchOff branches from a commit and updates the tree by doing:
// - Create a new Branch structure based on `name`
// - Create a copy of `commit` with the new branch as branch
// - Point the parent of the new commit to the original commit
// - Set the HEAD of the new branch to the new commit
// - Returns both created Branch and Commit as these are not written to disk
// If there is an error in any of the processes the VTree is not updated
func (v *VTree) BranchOff(commit Commit, branchName string) (branch Branch, branchRoot Commit, err error) {
	branch = NewBranch(branchName)
	branchRoot = NewCommit(commit, branch.Hash, branchName+":"+commit.Name, commit.Name+" as root of "+branchName)
	branch.HEAD = branchRoot.Hash
	if _, ok := v.Branches[branch.Hash]; ok {
		err = fmt.Errorf("Error: Branch has clash in VTree: %v, could not branch", branch.Hash)
		return
	}
	if err = v.Add(branchRoot); err != nil {
		return
	}
	v.Branches[branch.Hash] = branch.Name
	return
}
