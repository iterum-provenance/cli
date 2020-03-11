package idv

import (
	"errors"

	"github.com/Mantsje/iterum-cli/util"
)

// EnsureIDVRepo makes sure that we are in an idv repository
func EnsureIDVRepo() error {
	if !util.IsFileOrDir(IDVFolder) {
		return errors.New("Error: Not an idv repository")
	}
	return nil
}

// EnsureLOCAL makes sure that LOCAL points to a file
func EnsureLOCAL() error {
	errs := []error{}
	errs = append(errs, EnsureIDVRepo())
	if !util.IsFileOrDir(LOCAL) {
		errs = append(errs, errors.New("Error: Either .idv/LOCAL does not exist, or does not point to a file"))
	}
	return util.ReturnFirstErr(errs...)
}
