package utils

import "math/rand"

const characterPool = "ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

func GenerateRandomString(length uint8) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = characterPool[rand.Intn(len(characterPool))]
	}
	return string(b)
}
