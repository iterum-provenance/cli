package idv

import (
	"fmt"

	"github.com/iterum-provenance/cli/util"
)

// Some helper functions, only for internal use within this file

func (c commitTree) _toInterfaceMap() (out map[interface{}]interface{}) {
	out = make(map[interface{}]interface{})
	for key, val := range c {
		out[key] = val
	}
	return
}

func (c commitTree) _toHashNameMap() (out map[interface{}]interface{}) {
	out = make(map[interface{}]interface{})
	for hash, node := range c {
		out[hash] = node.Name
	}
	return
}

func (b branchMap) _toInterfaceMap() (out map[interface{}]interface{}) {
	out = make(map[interface{}]interface{})
	for key, val := range b {
		out[key] = val
	}
	return
}

type branchMap map[hash]string

type commitTree map[hash]vtreeNode

// vtreeNode is a node in the VTree, each node represents 1 commit
type vtreeNode struct {
	Name     string `json:"name"`
	Branch   hash   `json:"branch"`
	Children []hash `json:"children"`
	Parent   hash   `json:"parent"`
}

// VTree holds a version tree for data versioning
type VTree struct {
	Tree     commitTree `json:"tree"`
	Branches branchMap  `json:"branches"`
}

// NewVTree instantiates a new version tree and sets the root node
func NewVTree(root Commit, master Branch) VTree {
	v := VTree{make(commitTree), make(branchMap)}
	node := vtreeNode{
		Name:     root.Name,
		Branch:   master.Hash,
		Children: []hash{},
		Parent:   hash(""),
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

// ToFilePath returns a path to this VTree being: .idv/{local, remote}/vtreeFileName
// local indicates which of the 2 folders to use
func (v VTree) ToFilePath(local bool) string {
	if local {
		return localFolder + vtreeFileName
	}
	return remoteFolder + vtreeFileName
}

// Add appends a new (leaf) commit to the tree, returns an error on failure
func (v *VTree) addCommit(c Commit) error {
	node := vtreeNode{
		Name:     c.Name,
		Branch:   c.Branch,
		Children: []hash{},
		Parent:   c.Parent,
	}
	if v.isExistingCommit(c.Hash) {
		return fmt.Errorf("Error: Commit hash already exists in VTree: %v, could not add commit", c.Hash)
	}
	if !v.isExistingBranch(c.Branch) {
		return fmt.Errorf("Error: Commit refers to an unknown Branch in this VTree, could not add commit")
	}
	if !v.isExistingCommit(c.Parent) {
		return fmt.Errorf("Error: Commit's parent is non-existent. Parent: %v", c.Parent)
	}
	v.Tree[c.Hash] = node
	parentNode := v.Tree[c.Parent]
	parentNode.Children = append(parentNode.Children, c.Hash)
	v.Tree[c.Parent] = parentNode
	return nil
}

// addBranch attempts to add a branch to the version tree
func (v *VTree) addBranch(branch Branch) (err error) {
	if v.isExistingBranch(branch.Hash) {
		err = fmt.Errorf("Error: Branch has clash in VTree: %v, could not branch", branch.Hash)
		return
	}
	if !v.isExistingCommit(branch.HEAD) {
		err = fmt.Errorf("Error: Branch HEAD is a non-existing commit in VTree: %v. could not add branch", branch.HEAD)
		return
	}
	v.Branches[branch.Hash] = branch.Name
	return
}

// BranchOff branches from a commit and updates the tree by doing:
// - Create a new Branch structure based on `branchName`
// - Create a copy of `commit` with the new branch as branch
// - Point the parent of the new commit to the original commit
// - Set the HEAD of the new branch to the original commit
// - Returns both created Branch and Commit (since these 2 are local to this function for now)
// If there is an error in any of the processes the VTree is not updated
func (v *VTree) branchOff(commit Commit, branchName string) (branch Branch, branchRoot Commit, err error) {
	branch = NewBranch(branchName)
	branchRoot = NewCommit(commit, branch.Hash, branchName+":"+commit.Name, commit.Name+" as root of "+branchName)
	branch.HEAD = commit.Hash
	if err = v.addBranch(branch); err != nil {
		return
	}
	if err = v.addCommit(branchRoot); err != nil {
		delete(v.Branches, branch.Hash) // remove successful part  of update
		return
	}
	return
}

// _existing checks wheter a hash or name is part of this vtree.
// isHash indicates whether hashOrName should be interpreted as type hash
// isBranch indicates whether hashOrName references a branch (default is commit)
func (v VTree) _existing(hashOrName string, isHash, isBranch bool) bool {
	var target interface{} = hashOrName
	if isHash {
		target = hash(hashOrName)
	}
	if isBranch { // check for branches
		return util.MapContains(v.Branches._toInterfaceMap(), target, !isHash)
	}
	// check for commits
	return util.MapContains(v.Tree._toInterfaceMap(), target, !isHash)
}

// isExistingCommitName checks whether a given name is the name of a commit in the tree
func (v VTree) isExistingCommitName(cname string) bool {
	return v._existing(cname, false, false)
}

// isExistingCommit checks whether a given commit hash exists in the tree
func (v VTree) isExistingCommit(chash hash) bool {
	return v._existing(chash.String(), true, false)
}

// isExistingBranchName checks whether a given name is the name of an existing branch in the tree
func (v VTree) isExistingBranchName(bname string) bool {
	return v._existing(bname, false, true)
}

// isExistingBranch checks whether a given branch hash exists in the tree
func (v VTree) isExistingBranch(bhash hash) bool {
	return v._existing(bhash.String(), true, true)
}

// _getHashByName returns the the hash corresponding to a name value. Errors if name is not found
// it returns the first matching value found
func (v VTree) _getHashByName(m map[interface{}]interface{}, name string) (h hash, err error) {
	out, err := util.GetKeyByValue(m, name)
	if err != nil {
		return
	}
	h = out.(hash)
	return
}

// getCommitHashByName returns the hash corresponding to a commit name, errors if non-existent
func (v VTree) getCommitHashByName(name string) (h hash, err error) {
	return v._getHashByName(v.Tree._toHashNameMap(), name)
}

// getBranchHashByName returns the hash corresponding to a branch name, errors if non-existent
func (v VTree) getBranchHashByName(name string) (h hash, err error) {
	return v._getHashByName(v.Branches._toInterfaceMap(), name)
}

// duplicate makes a copy of the passed tree
func (v VTree) duplicate() VTree {
	copy := VTree{make(commitTree), make(branchMap)}
	for key, val := range v.Tree {
		copy.Tree[key] = val
	}
	for key, val := range v.Branches {
		copy.Branches[key] = val
	}
	return copy
}
