package utils

import (
	"math/rand"
)

func GenerateRandomCode(length int) string {
	var letterRune = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, length)

	for i := range b {
		b[i] = letterRune[rand.Intn(len(letterRune))]
	}

	return string(b)
}
