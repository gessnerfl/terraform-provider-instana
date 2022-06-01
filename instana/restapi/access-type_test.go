package restapi_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/require"
)

func TestShouldReturnTrueForAllSupportedAccessTypes(t *testing.T) {
	for _, v := range SupportedAccessTypes {
		require.True(t, SupportedAccessTypes.IsSupported(v))
	}
}

func TestShouldReturnFalseForAllNonSupportedAccessTypes(t *testing.T) {
	for _, v := range []string{"FOO", "BAR", "INVALID"} {
		require.False(t, SupportedAccessTypes.IsSupported(AccessType(v)))
	}
}

func TestShouldReturnSupportedAccessTypesAsStringSlice(t *testing.T) {
	expected := []string{"READ", "READ_WRITE"}
	require.Equal(t, expected, SupportedAccessTypes.ToStringSlice())
}
