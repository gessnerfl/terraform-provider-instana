package restapi_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/require"
)

func TestShouldReturnTrueForAllSupportedGranularities(t *testing.T) {
	for _, v := range SupportedGranularities {
		require.True(t, SupportedGranularities.IsSupported(v))
	}
}

func TestShouldReturnFalseForAllNonSupportedGranularities(t *testing.T) {
	for _, v := range []int32{1, 2, 3} {
		require.False(t, SupportedGranularities.IsSupported(Granularity(v)))
	}
}

func TestShouldReturnSupportedGranularitiesAsIntSlice(t *testing.T) {
	expected := []int{300000, 600000, 900000, 1200000, 1800000}
	require.Equal(t, expected, SupportedGranularities.ToIntSlice())
}
