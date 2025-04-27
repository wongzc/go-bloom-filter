// Package randomstring provides utilities to generate random strings for testing purposes.
package randomstring

import (
	"math/rand"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// generate random string of length 1-25
func RandString() string {
	n := rand.Intn(25) + 1
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}