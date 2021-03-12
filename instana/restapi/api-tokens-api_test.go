package restapi_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/assert"
)

const (
	apiTokenAccessGrantingToken = "api-token-access-granting-token"
	apiTokenName                = "api-token-name"
)

func TestValidMinimalAPIToken(t *testing.T) {
	apiToken := APIToken{
		AccessGrantingToken: apiTokenAccessGrantingToken,
		Name:                apiTokenName,
	}

	err := apiToken.Validate()
	assert.Nil(t, err)
}

func TestValidFullAPIToken(t *testing.T) {
	apiToken := APIToken{
		ID:                                   "id",
		AccessGrantingToken:                  apiTokenAccessGrantingToken,
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

	assert.Equal(t, "id", apiToken.GetID())

	err := apiToken.Validate()
	assert.Nil(t, err)
}

func TestInvalidAPITokenBecauseOfMissingId(t *testing.T) {
	apiToken := APIToken{
		Name: apiTokenName,
	}

	err := apiToken.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "ID")
}

func TestInvalidAPITokenBecauseOfMissingName(t *testing.T) {
	apiToken := APIToken{
		ID: apiTokenAccessGrantingToken,
	}

	err := apiToken.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Name")
}
