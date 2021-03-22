package restapi_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/require"
)

const (
	apiTokenInternalID          = "internal-id"
	apiTokenAccessGrantingToken = "access-granting-token"
	apiTokenName                = "api-token-name"
)

func TestValidMinimalAPIToken(t *testing.T) {
	apiToken := APIToken{
		InternalID:          apiTokenInternalID,
		AccessGrantingToken: apiTokenInternalID,
		Name:                apiTokenName,
	}

	require.Equal(t, apiTokenInternalID, apiToken.GetIDForResourcePath())
	err := apiToken.Validate()
	require.Nil(t, err)
}

func TestValidFullAPIToken(t *testing.T) {
	apiToken := APIToken{
		ID:                                   "id",
		AccessGrantingToken:                  apiTokenAccessGrantingToken,
		InternalID:                           apiTokenInternalID,
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

	require.Equal(t, apiTokenInternalID, apiToken.GetIDForResourcePath())

	err := apiToken.Validate()
	require.Nil(t, err)
}

func TestInvalidAPITokenWhenInternalIdIsMissing(t *testing.T) {
	apiToken := APIToken{
		AccessGrantingToken: apiTokenAccessGrantingToken,
		Name:                apiTokenName,
	}

	err := apiToken.Validate()

	require.Error(t, err)
	require.Contains(t, err.Error(), "Internal ID")
}

func TestInvalidAPITokenWhenInternalIdIsBlank(t *testing.T) {
	apiToken := APIToken{
		AccessGrantingToken: apiTokenAccessGrantingToken,
		InternalID:          " ",
		Name:                apiTokenName,
	}

	err := apiToken.Validate()

	require.Error(t, err)
	require.Contains(t, err.Error(), "Internal ID")
}

func TestInvalidAPITokenWhenAccessGrantingTokenIsMissing(t *testing.T) {
	apiToken := APIToken{
		InternalID: apiTokenInternalID,
		Name:       apiTokenName,
	}

	err := apiToken.Validate()

	require.Error(t, err)
	require.Contains(t, err.Error(), "Access Granting Token")
}

func TestInvalidAPITokenWhenAccessGrantingTokenIsBlank(t *testing.T) {
	apiToken := APIToken{
		InternalID:          apiTokenInternalID,
		AccessGrantingToken: " ",
		Name:                apiTokenName,
	}

	err := apiToken.Validate()

	require.Error(t, err)
	require.Contains(t, err.Error(), "Access Granting Token")
}

func TestInvalidAPITokenWhenNameIsMissing(t *testing.T) {
	apiToken := APIToken{
		InternalID:          apiTokenInternalID,
		AccessGrantingToken: apiTokenAccessGrantingToken,
	}

	err := apiToken.Validate()

	require.Error(t, err)
	require.Contains(t, err.Error(), "Name")
}

func TestInvalidAPITokenWhenNameIsBlank(t *testing.T) {
	apiToken := APIToken{
		InternalID:          apiTokenInternalID,
		AccessGrantingToken: apiTokenAccessGrantingToken,
		Name:                " ",
	}

	err := apiToken.Validate()

	require.Error(t, err)
	require.Contains(t, err.Error(), "Name")
}
