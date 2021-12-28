package restapi_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/require"
)

func TestShouldReturnTrueForAllSupportedThresholdOperators(t *testing.T) {
	for _, v := range SupportedThresholdOperators {
		require.True(t, SupportedThresholdOperators.IsSupported(v))
	}
}

func TestShouldReturnFalseForAllNonSupportedThresholdOperators(t *testing.T) {
	for _, v := range []string{"FOO", "BAR", "INVALID"} {
		require.False(t, SupportedThresholdOperators.IsSupported(ThresholdOperator(v)))
	}
}

func TestShouldReturnSupportedThresholdOperatorsAsStringSlice(t *testing.T) {
	expected := []string{">", ">=", "<", "<="}
	require.Equal(t, expected, SupportedThresholdOperators.ToStringSlice())
}

func TestShouldReturnTrueForAllSupportedThresholdSeasonalities(t *testing.T) {
	for _, v := range SupportedThresholdSeasonalities {
		require.True(t, SupportedThresholdSeasonalities.IsSupported(v))
	}
}

func TestShouldReturnFalseForAllNonSupportedThresholdSeasonalities(t *testing.T) {
	for _, v := range []string{"FOO", "BAR", "INVALID"} {
		require.False(t, SupportedThresholdSeasonalities.IsSupported(ThresholdSeasonality(v)))
	}
}

func TestShouldReturnSupportedThresholdSeasonalitiesAsStringSlice(t *testing.T) {
	expected := []string{"WEEKLY", "DAILY"}
	require.Equal(t, expected, SupportedThresholdSeasonalities.ToStringSlice())
}
