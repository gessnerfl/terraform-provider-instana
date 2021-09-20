package utils_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/stretchr/testify/require"
)

func TestShouldCreateInt64PointerFromInt64(t *testing.T) {
	value := int64(123)

	require.Equal(t, &value, Int64Ptr(value))
}
