package restapi_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/require"
)

func TestShouldReturnSupportedRelationTypesAsStringSlice(t *testing.T) {
	expected := []string{"USER", "API_TOKEN", "ROLE", "TEAM", "GLOBAL"}
	require.Equal(t, expected, SupportedRelationTypes.ToStringSlice())
}
