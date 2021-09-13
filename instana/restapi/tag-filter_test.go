package restapi_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/require"
)

func TestShouldReturnTrueForAllSupportedTagFilterEntityTypes(t *testing.T) {
	for _, v := range SupportedTagFilterEntities {
		require.True(t, SupportedTagFilterEntities.IsSupported(v))
	}
}

func TestShouldReturnFalseForAllNonSupportedTagFilterEntityTypes(t *testing.T) {
	for _, v := range []string{"FOO", "BAR", "INVALID"} {
		require.False(t, SupportedTagFilterEntities.IsSupported(TagFilterEntity(v)))
	}
}

func TestShouldConvertTagFilterEntitiesToStringSlice(t *testing.T) {
	expectedResult := []string{"SOURCE", "DESTINATION", "NOT_APPLICABLE"}
	require.Equal(t, expectedResult, SupportedTagFilterEntities.ToStringSlice())
}

func TestShouldReturnTrueForAllSupportedComparisonOperators(t *testing.T) {
	for _, v := range SupportedComparisonOperators {
		require.True(t, SupportedComparisonOperators.IsSupported(v))
	}
}

func TestShouldReturnFalseForAllNonSupportedComparisonOperators(t *testing.T) {
	for _, v := range append(SupportedUnaryExpressionOperators, "INVALID_OPERATOR") {
		require.False(t, SupportedComparisonOperators.IsSupported(v))
	}
}

func TestShouldReturnTrueForAllSupportedUnaryExpressionOperators(t *testing.T) {
	for _, v := range SupportedUnaryExpressionOperators {
		require.True(t, SupportedUnaryExpressionOperators.IsSupported(v))
	}
}

func TestShouldReturnFalseForAllNonSupportedUnaryExpressionOperators(t *testing.T) {
	for _, v := range append(SupportedComparisonOperators, "INVALID_OPERATOR") {
		require.False(t, SupportedUnaryExpressionOperators.IsSupported(v))
	}
}
