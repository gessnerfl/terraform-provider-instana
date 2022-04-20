package restapi_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/require"
)

func TestShouldReturnTrueForAllSupportedWebsiteImpactMeasurementMethods(t *testing.T) {
	for _, v := range SupportedWebsiteImpactMeasurementMethods {
		require.True(t, SupportedWebsiteImpactMeasurementMethods.IsSupported(v))
	}
}

func TestShouldReturnFalseForAllNonSupportedWebsiteImpactMeasurementMethods(t *testing.T) {
	for _, v := range []string{"FOO", "BAR", "INVALID"} {
		require.False(t, SupportedWebsiteImpactMeasurementMethods.IsSupported(WebsiteImpactMeasurementMethod(v)))
	}
}

func TestShouldReturnSupportedWebsiteImpactMeasurementMethodsAsStringSlice(t *testing.T) {
	expected := []string{"AGGREGATED", "PER_WINDOW"}
	require.Equal(t, expected, SupportedWebsiteImpactMeasurementMethods.ToStringSlice())
}
