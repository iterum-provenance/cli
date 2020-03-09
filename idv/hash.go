package idv

import (
	"math/rand"
	"time"
)

type hash string

const allowedChars = "abcdefghijklmnopqrstuvwxyz0123456789"

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
