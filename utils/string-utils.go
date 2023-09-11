package utils

import (
	"regexp"
	"strings"
)

// IsBlank returns true when the provided string is empty or only contains whitespace characters
func IsBlank(input string) bool {
	return len(strings.TrimSpace(input)) == 0
}

// StringPtr converts a string to a string pointer
func StringPtr(input string) *string {
	return &input
}

var multiSpacesRegexp = regexp.MustCompile("( ){2,}")

// RemoveNewLinesAndTabs removes all new lines and tabs from a string
func RemoveNewLinesAndTabs(input string) string {
	spacesOnly := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(input, "\n\r", " "), "\n", " "), "\t", " ")
	return multiSpacesRegexp.ReplaceAllString(spacesOnly, " ")
}
