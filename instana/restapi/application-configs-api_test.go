package restapi_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

func TestShouldReturnTrueForAllSupportedApplicationConfigScopes(t *testing.T) {
	for _, scope := range SupportedApplicationConfigScopes {
		t.Run(fmt.Sprintf("TestShouldReturnTrueForSupportedApplicationConfigScope%s", string(scope)), createTestCaseToVerifySupportedApplicationConfigScope(scope))
	}
}

func createTestCaseToVerifySupportedApplicationConfigScope(scope ApplicationConfigScope) func(t *testing.T) {
	return func(t *testing.T) {
		require.True(t, SupportedApplicationConfigScopes.IsSupported(scope))
	}
}

func TestShouldReturnfalseWhenApplicationConfigScopeIsNotSupported(t *testing.T) {
	require.False(t, SupportedApplicationConfigScopes.IsSupported(ApplicationConfigScope(valueInvalid)))
}

func TestShouldReturnStringRepresentationOfSupporedApplicationConfigScopes(t *testing.T) {
	require.Equal(t, []string{"INCLUDE_NO_DOWNSTREAM", "INCLUDE_IMMEDIATE_DOWNSTREAM_DATABASE_AND_MESSAGING", "INCLUDE_ALL_DOWNSTREAM"}, SupportedApplicationConfigScopes.ToStringSlice())
}

func TestShouldReturnTrueForAllSupportedMatcherExpressionEntities(t *testing.T) {
	for _, entity := range SupportedMatcherExpressionEntities {
		t.Run(fmt.Sprintf("TestShouldReturnTrueForSupportedMatcherExpressionEntity%s", string(entity)), createTestCaseToVerifySupportedMatcherExpressionEntity(entity))
	}
}

func createTestCaseToVerifySupportedMatcherExpressionEntity(entity MatcherExpressionEntity) func(t *testing.T) {
	return func(t *testing.T) {
		require.True(t, SupportedMatcherExpressionEntities.IsSupported(entity))
	}
}

func TestShouldReturnfalseWhenMatcherExpressionEntityIsNotSupported(t *testing.T) {
	require.False(t, SupportedMatcherExpressionEntities.IsSupported(MatcherExpressionEntity(valueInvalid)))
}

func TestShouldReturnStringRepresentationOfSupporedMatcherExpressionEntities(t *testing.T) {
	require.Equal(t, []string{"SOURCE", "DESTINATION", "NOT_APPLICABLE"}, SupportedMatcherExpressionEntities.ToStringSlice())
}
