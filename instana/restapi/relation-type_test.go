package restapi_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/require"
)

func TestShouldReturnTrueForAllSupportedRelationTypes(t *testing.T) {
	for _, v := range SupportedRelationTypes {
		require.True(t, SupportedRelationTypes.IsSupported(v))
	}
}

func TestShouldReturnFalseForAllNonSupportedRelationTypes(t *testing.T) {
	for _, v := range []string{"FOO", "BAR", "INVALID"} {
		require.False(t, SupportedRelationTypes.IsSupported(RelationType(v)))
	}
}

func TestShouldReturnSupportedRelationTypesAsStringSlice(t *testing.T) {
	expected := []string{"USER", "API_TOKEN", "ROLE", "TEAM", "GLOBAL"}
	require.Equal(t, expected, SupportedRelationTypes.ToStringSlice())
}
