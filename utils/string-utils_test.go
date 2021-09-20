package utils_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/stretchr/testify/require"
)

func TestShouldReturnTrueWhenStringIsEmpty(t *testing.T) {
	require.True(t, IsBlank(""))
}

func TestShouldReturnTrueWhenStringContainsOnlySpaces(t *testing.T) {
	require.True(t, IsBlank("     "))
}

func TestShouldReturnFalseWhenStringContainsNonWhitespaceCharacters(t *testing.T) {
	require.False(t, IsBlank("  ba  "))
}

func TestShouldCreateRandomString(t *testing.T) {
	require.Equal(t, 64, len(RandomString(64)))
}

func TestShouldCreateStringPointerFromString(t *testing.T) {
	value := "string"

	require.Equal(t, &value, StringPtr(value))
}
