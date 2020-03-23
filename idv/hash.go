package idv

import (
	"math/rand"
	"time"
)

type hash string

const allowedChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// newHash generate a new hash of given length
func newHash(length int) hash {
	b := make([]byte, length)
	for i := range b {
		b[i] = allowedChars[rand.Int63()%int64(len(allowedChars))]
	}
	return hash(b)
}

func (h hash) String() string {
	return string(h)
}

func (h hash) toPath(local bool, extension string) string {
	if local {
		return localFolder + h.String() + extension
	}
	return remoteFolder + h.String() + extension
}

func (h hash) toBranchPath(local bool) string {
	return h.toPath(local, branchFileExt)
}

func (h hash) toCommitPath(local bool) string {
	return h.toPath(local, commitFileExt)
}
