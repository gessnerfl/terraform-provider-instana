package utils_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/stretchr/testify/assert"
)

func TestShouldReturnTrueWhenStringIsEmpty(t *testing.T) {
	if !IsBlank("") {
		t.Fatal("Expected to return true for empty string")
	}
}

func TestShouldReturnTrueWhenStringContainsOnlySpaces(t *testing.T) {
	if !IsBlank("    ") {
		t.Fatal("Expected to return true for string containing only spaces")
	}
}

func TestShouldReturnFalseWhenStringContainsNonWhitespaceCharacters(t *testing.T) {
	if IsBlank("  ba  ") {
		t.Fatal("Expected to return false for string containing non whitespaces")
	}
}

func TestShouldCreateRandomString(t *testing.T) {
	assert.Equal(t, 64, len(RandomString(64)))
}
