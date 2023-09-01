package restapi_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/require"
)

func TestShouldReturnSupportedWebsiteImpactMeasurementMethodsAsStringSlice(t *testing.T) {
	expected := []string{"AGGREGATED", "PER_WINDOW"}
	require.Equal(t, expected, SupportedWebsiteImpactMeasurementMethods.ToStringSlice())
}
