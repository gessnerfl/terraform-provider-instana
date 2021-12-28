package restapi_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/require"
)

func TestShouldReturnTrueForAllSupportedAggregations(t *testing.T) {
	for _, v := range SupportedAggregations {
		require.True(t, SupportedAggregations.IsSupported(v))
	}
}

func TestShouldReturnFalseForAllNonSupportedAggregations(t *testing.T) {
	for _, v := range []string{"FOO", "BAR", "INVALID"} {
		require.False(t, SupportedAggregations.IsSupported(Aggregation(v)))
	}
}

func TestShouldReturnSupportedAggregationsAsStringSlice(t *testing.T) {
	expected := []string{"SUM", "MEAN", "MAX", "MIN", "P25", "P50", "P75", "P90", "P95", "P98", "P99", "P99_9", "P99_99", "DISTRIBUTION", "DISTINCT_COUNT", "SUM_POSITIVE"}
	require.Equal(t, expected, SupportedAggregations.ToStringSlice())
}
