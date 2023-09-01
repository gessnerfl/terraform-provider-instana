package restapi_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/require"
)

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

func TestShouldReturnSupportedOperatorsAsStringSlice(t *testing.T) {
	expected := []string{"IS_EMPTY", "NOT_EMPTY", "IS_BLANK", "NOT_BLANK"}
	require.Equal(t, expected, SupportedUnaryExpressionOperators.ToStringSlice())
}
