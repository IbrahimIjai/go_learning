package helpers

import (
	"crypto/rand"
	"math/big"
)

const (
	alphabet   = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	CodeLength = 7 // 62^7 ≈ 3.5 trillion combinations
)

// NewCode returns a random base62 code.
//
// We use random codes + retry-on-collision rather than a base62-encoded
// counter: counters leak how many URLs exist and are enumerable. crypto/rand
// keeps codes unpredictable.
func NewCode() (string, error) {
	b := make([]byte, CodeLength)
	max := big.NewInt(int64(len(alphabet)))
	for i := range b {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b[i] = alphabet[n.Int64()]
	}
	return string(b), nil
}

// Valid reports whether s looks like a code we could have issued. Cheap
// validation before hitting the store saves a DB roundtrip for garbage input.
func Valid(s string) bool {
	if len(s) != CodeLength {
		return false
	}
	for _, c := range s {
		if !isBase62(c) {
			return false
		}
	}
	return true
}

func isBase62(c rune) bool {
	return (c >= '0' && c <= '9') ||
		(c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z')
}
