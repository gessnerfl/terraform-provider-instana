package restapi_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/require"
)

func TestShouldReturnSupportedLogLevelsAsStringSlice(t *testing.T) {
	expected := []string{"WARN", "ERROR", "ANY"}
	require.Equal(t, expected, SupportedLogLevels.ToStringSlice())
}
