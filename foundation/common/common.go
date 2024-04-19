package common

import "crypto/rand"

func MustGenerateRandomKey(keyLength int) []byte {
	key := make([]byte, keyLength)
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}
	return key
}
