package restapi

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShouldReturnTrueForAllSupportedCustomPayloadTypes(t *testing.T) {
	for _, v := range SupportedCustomPayloadTypes {
		require.True(t, SupportedCustomPayloadTypes.IsSupported(v))
	}
}

func TestShouldReturnFalseForAllNonSupportedCustomPayloadTypes(t *testing.T) {
	for _, v := range []string{"FOO", "BAR", "INVALID"} {
		require.False(t, SupportedCustomPayloadTypes.IsSupported(CustomPayloadType(v)))
	}
}

func TestShouldReturnSupportedCustomPayloadTypesAsStringSlice(t *testing.T) {
	expected := []string{"staticString", "dynamic"}
	require.Equal(t, expected, SupportedCustomPayloadTypes.ToStringSlice())
}
