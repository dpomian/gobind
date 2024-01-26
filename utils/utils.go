package utils

import (
	"math/rand"
	"time"
)

// Generates a random string of a given length using a specific character set
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	var b []byte
	for i := 0; i < length; i++ {
		b = append(b, charset[seededRand.Intn(len(charset))])
	}
	return string(b)
}
