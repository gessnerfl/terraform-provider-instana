package utils

import (
	"math/rand"
	"strings"
	"time"
)

//IsBlank returns true when the provided string is empty or only contains whitespace characters
func IsBlank(input string) bool {
	return len(strings.TrimSpace(input)) == 0
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

//RandomString creates a random string of the given length
func RandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

//StringPtr converts a string to a string pointer
func StringPtr(input string) *string {
	return &input
}
