package restapi_test

import (
	"fmt"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/require"
)

func TestShouldReturnTrueForAllSupportedApplicationConfigBoundaryScopes(t *testing.T) {
	for _, scope := range SupportedApplicationConfigBoundaryScopes {
		t.Run(fmt.Sprintf("TestShouldReturnTrueForSupportedBoundaryScope%s", string(scope)), createTestCaseToVerifySupportedApplicationConfigBoundaryScope(scope))
	}
}

func createTestCaseToVerifySupportedApplicationConfigBoundaryScope(scope BoundaryScope) func(t *testing.T) {
	return func(t *testing.T) {
		require.True(t, SupportedApplicationConfigBoundaryScopes.IsSupported(scope))
	}
}

func TestShouldReturnFalseWhenApplicationConfigBoundaryScopeIsNotSupported(t *testing.T) {
	require.False(t, SupportedApplicationConfigBoundaryScopes.IsSupported(BoundaryScope(valueInvalid)))
}

func TestShouldReturnStringRepresentationOfSupportedApplicationConfigBoundaryScopes(t *testing.T) {
	require.Equal(t, []string{"ALL", "INBOUND", "DEFAULT"}, SupportedApplicationConfigBoundaryScopes.ToStringSlice())
}

func TestShouldReturnStringRepresentationOfSupportedApplicationAlertConfigBoundaryScopes(t *testing.T) {
	require.Equal(t, []string{"ALL", "INBOUND"}, SupportedApplicationAlertConfigBoundaryScopes.ToStringSlice())
}
