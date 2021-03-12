package restapi_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/require"
)

const (
	apiTokenID   = "ID"
	apiTokenName = "api-token-name"
)

func TestValidMinimalAPIToken(t *testing.T) {
	apiToken := APIToken{
		ID:                  apiTokenID,
		AccessGrantingToken: apiTokenID,
		Name:                apiTokenName,
	}

	err := apiToken.Validate()
	require.Nil(t, err)
}

func TestValidFullAPIToken(t *testing.T) {
	apiToken := APIToken{
		ID:                                   apiTokenID,
		AccessGrantingToken:                  apiTokenID,
		InternalID:                           "internal-id",
		Name:                                 apiTokenName,
		CanConfigureServiceMapping:           true,
		CanConfigureEumApplications:          true,
		CanConfigureMobileAppMonitoring:      true,
		CanConfigureUsers:                    true,
		CanInstallNewAgents:                  true,
		CanSeeUsageInformation:               true,
		CanConfigureIntegrations:             true,
		CanSeeOnPremiseLicenseInformation:    true,
		CanConfigureCustomAlerts:             true,
		CanConfigureAPITokens:                true,
		CanConfigureAgentRunMode:             true,
		CanViewAuditLog:                      true,
		CanConfigureAgents:                   true,
		CanConfigureAuthenticationMethods:    true,
		CanConfigureApplications:             true,
		CanConfigureTeams:                    true,
		CanConfigureReleases:                 true,
		CanConfigureLogManagement:            true,
		CanCreatePublicCustomDashboards:      true,
		CanViewLogs:                          true,
		CanViewTraceDetails:                  true,
		CanConfigureSessionSettings:          true,
		CanConfigureServiceLevelIndicators:   true,
		CanConfigureGlobalAlertPayload:       true,
		CanConfigureGlobalAlertConfigs:       true,
		CanViewAccountAndBillingInformation:  true,
		CanEditAllAccessibleCustomDashboards: true,
	}

	require.Equal(t, apiTokenID, apiToken.GetID())

	err := apiToken.Validate()
	require.Nil(t, err)
}

func TestInvalidAPITokenWhenIdIsMissing(t *testing.T) {
	apiToken := APIToken{
		AccessGrantingToken: apiTokenID,
		Name:                apiTokenName,
	}

	err := apiToken.Validate()

	require.Error(t, err)
	require.Contains(t, err.Error(), "ID")
}

func TestInvalidAPITokenWhenAccessGrantingTokenIsMissing(t *testing.T) {
	apiToken := APIToken{
		ID:   apiTokenID,
		Name: apiTokenName,
	}

	err := apiToken.Validate()

	require.Error(t, err)
	require.Contains(t, err.Error(), "Access Granting Token and ID")
}

func TestInvalidAPITokenWhenAccessGrantingTokenAndIDAreNotEqual(t *testing.T) {
	apiToken := APIToken{
		ID:                  apiTokenID,
		AccessGrantingToken: "foo",
		Name:                apiTokenName,
	}

	err := apiToken.Validate()

	require.Error(t, err)
	require.Contains(t, err.Error(), "Access Granting Token and ID")
}

func TestInvalidAPITokenWhenNameIsMissing(t *testing.T) {
	apiToken := APIToken{
		ID:                  apiTokenID,
		AccessGrantingToken: apiTokenID,
	}

	err := apiToken.Validate()

	require.Error(t, err)
	require.Contains(t, err.Error(), "Name")
}
