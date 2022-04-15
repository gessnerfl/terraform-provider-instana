package restapi_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/require"
)

func TestShouldReturnTrueForAllSupportedLogLevels(t *testing.T) {
	for _, v := range SupportedLogLevels {
		require.True(t, SupportedLogLevels.IsSupported(v))
	}
}

func TestShouldReturnFalseForAllNonSupportedLogLevels(t *testing.T) {
	for _, v := range []string{"FOO", "BAR", "INVALID"} {
		require.False(t, SupportedLogLevels.IsSupported(LogLevel(v)))
	}
}

func TestShouldReturnSupportedLogLevelsAsStringSlice(t *testing.T) {
	expected := []string{"WARN", "ERROR", "ANY"}
	require.Equal(t, expected, SupportedLogLevels.ToStringSlice())
}
