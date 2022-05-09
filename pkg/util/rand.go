package util

import "crypto/rand"

func GetRandomBytes(l int) []byte {
	b := make([]byte, l)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}
