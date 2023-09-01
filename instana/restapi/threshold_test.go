package restapi_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/require"
)

func TestShouldReturnSupportedThresholdOperatorsAsStringSlice(t *testing.T) {
	expected := []string{">", ">=", "<", "<="}
	require.Equal(t, expected, SupportedThresholdOperators.ToStringSlice())
}

func TestShouldReturnSupportedThresholdSeasonalitiesAsStringSlice(t *testing.T) {
	expected := []string{"WEEKLY", "DAILY"}
	require.Equal(t, expected, SupportedThresholdSeasonalities.ToStringSlice())
}
