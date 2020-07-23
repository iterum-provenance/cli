// Package storage contains the different supported storage backend and their files
package storage

import "fmt"

// Backend is an enum defining which backend is used for data storage
type Backend string

// Enum-like values for type
const (
	Local      Backend = "Local"       //"local-storage"
	AmazonS3   Backend = "AmazonS3"    //"amazon-s3"
	CloudStore Backend = "GoogleCloud" //"cloud-store"
)

// NewBackend creates a new instance and validates it
func NewBackend(backend string) (Backend, error) {
	b := Backend(backend)
	return b, b.IsValid()
}

// IsValid checks the validity
func (b Backend) IsValid() error {
	switch b {
	case Local, AmazonS3, CloudStore:
		return nil
	}
	return fmt.Errorf("Error: %v is not a valid Backend, pick one of { %v, %v, %v }", b, Local, AmazonS3, CloudStore)
}

func (b Backend) String() string {
	return string(b)
}
