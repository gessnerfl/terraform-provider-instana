package restapi_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/require"
)

func TestShouldReturnTrueForAllSupportedApplicationAlertEvaluationTypes(t *testing.T) {
	for _, v := range SupportedApplicationAlertEvaluationTypes {
		require.True(t, SupportedApplicationAlertEvaluationTypes.IsSupported(v))
	}
}

func TestShouldReturnFalseForAllNonSupportedApplicationAlertEvaluationTypes(t *testing.T) {
	for _, v := range []string{"FOO", "BAR", "INVALID"} {
		require.False(t, SupportedApplicationAlertEvaluationTypes.IsSupported(ApplicationAlertEvaluationType(v)))
	}
}

func TestShouldReturnSupportedApplicationAlertEvaluationTypesAsStringSlice(t *testing.T) {
	expected := []string{"PER_AP", "PER_AP_SERVICE", "PER_AP_ENDPOINT"}
	require.Equal(t, expected, SupportedApplicationAlertEvaluationTypes.ToStringSlice())
}
