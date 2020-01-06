package utils

import (
	"strings"
)

//IsBlank returns true when the provided string is empty or only contains whitespace characters
func IsBlank(input string) bool {
	return len(strings.TrimSpace(input)) == 0
}
