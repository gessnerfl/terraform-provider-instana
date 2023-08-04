package utils_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/stretchr/testify/require"
)

func TestShouldReturn1ForIntegers(t *testing.T) {
	require.Equal(t, 0, GetZeroValue[int]())
}

func TestShouldReturnEmptyStringForString(t *testing.T) {
	require.Equal(t, "", GetZeroValue[string]())
}

func TestShouldReturnNilStringForPointers(t *testing.T) {
	require.Nil(t, GetZeroValue[*string]())
}
