package utils_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/stretchr/testify/assert"
)

func TestShouldReturnTrueWhenStringIsEmpty(t *testing.T) {
	assert.True(t, IsBlank(""))
}

func TestShouldReturnTrueWhenStringContainsOnlySpaces(t *testing.T) {
	assert.True(t, IsBlank("     "))
}

func TestShouldReturnFalseWhenStringContainsNonWhitespaceCharacters(t *testing.T) {
	assert.False(t, IsBlank("  ba  "))
}

func TestShouldCreateRandomString(t *testing.T) {
	assert.Equal(t, 64, len(RandomString(64)))
}
