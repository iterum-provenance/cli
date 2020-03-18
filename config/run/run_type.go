package run

import (
	"bytes"
	"errors"
	"fmt"
)

// RunType is used for defining what kind of deployment is used
type RunType string

// Enum-like values allowed for RunType type
const (
	Local       RunType = "local"
	Distributed RunType = "distributed"
)

// NewRunType creates a new RunType instance and validates it
func NewRunType(rawRunType string) (RunType, error) {
	var rt RunType = RunType(rawRunType)
	return rt, rt.IsValid()
}

// InferRunType tries to infer a RunType from a string
func InferRunType(raw string) (RunType, error) {
	if raw == "l" { // Shorthand for local
		return NewRunType(string(Local))
	} else if raw == "d" { // Shorthand for distribbuted
		return NewRunType(string(Distributed))
	}
	return NewRunType(raw)
}

// IsValid checks the validity of the RunType
func (rt RunType) IsValid() error {
	switch rt {
	case Local, Distributed:
		return nil
	}
	return errors.New("Error: Invalid RunType")
}

// String converts RunType to string representation
func (rt RunType) String() string {
	return string(rt)
}

// AllowedVariables returns a formatted string on how to set this type with the set command
func (rt RunType) AllowedVariables() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "RunType         { %v, %v }\n", Distributed, Local)
	return buf.String()
}
