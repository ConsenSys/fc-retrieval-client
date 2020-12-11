package prng

import (
	"crypto/rand"
)

// TODO initially, just use the standard library to generate numbers
// We need to investigate to see how secure / insecure this is


// GenerateRandomBytes generates zero or more random numbers
func GenerateRandomBytes(b []byte) {
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

}