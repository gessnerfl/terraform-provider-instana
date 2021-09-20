package utils_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/stretchr/testify/require"
)

func TestShouldCreateBoolPointerFromBool(t *testing.T) {
	value := true

	require.Equal(t, &value, BoolPtr(value))
}
