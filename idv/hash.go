package idv

import (
	"math/rand"
	"time"
)

// hash is the internal type used to represent hashes
// its a string alias which can be randomly generated
type hash string

// allowedChars are the possible characters used in hashes
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

// toPath converts a hash to a file stored in the `.idv` folder of a data repo
func (h hash) toPath(local bool, extension string) string {
	if local {
		return localFolder + h.String() + extension
	}
	return remoteFolder + h.String() + extension
}

// toBranchPath returns target/folder/hash.branch
func (h hash) toBranchPath(local bool) string {
	return h.toPath(local, branchFileExt)
}

// toCommitPath returns target/folder/hash.commit
func (h hash) toCommitPath(local bool) string {
	return h.toPath(local, commitFileExt)
}
