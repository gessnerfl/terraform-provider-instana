package restapi

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShouldReturnSupportedCustomPayloadTypesAsStringSlice(t *testing.T) {
	expected := []string{"staticString", "dynamic"}
	require.Equal(t, expected, SupportedCustomPayloadTypes.ToStringSlice())
}
