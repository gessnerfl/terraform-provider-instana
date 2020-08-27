package restapi_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

func TestShouldSuccessfullyUnmarshalUserRole(t *testing.T) {
	userRole := UserRole{
		ID:                                "role-id",
		Name:                              "role-name",
		CanConfigureServiceMapping:        true,
		CanConfigureEumApplications:       true,
		CanConfigureMobileAppMonitoring:   true,
		CanConfigureUsers:                 true,
		CanInstallNewAgents:               true,
		CanSeeUsageInformation:            true,
		CanConfigureIntegrations:          true,
		CanSeeOnPremiseLicenseInformation: true,
		CanConfigureRoles:                 true,
		CanConfigureCustomAlerts:          true,
		CanConfigureAPITokens:             true,
		CanConfigureAgentRunMode:          true,
		CanViewAuditLog:                   true,
		CanConfigureObjectives:            true,
		CanConfigureAgents:                true,
		CanConfigureAuthenticationMethods: true,
		CanConfigureApplications:          true,
		CanConfigureTeams:                 true,
		RestrictedAccess:                  true,
		CanConfigureReleases:              true,
		CanConfigureLogManagement:         true,
		CanCreatePublicCustomDashboards:   true,
		CanViewLogs:                       true,
		CanViewTraceDetails:               true,
	}

	serializedJSON, _ := json.Marshal(userRole)

	result, err := NewUserRoleUnmarshaller().Unmarshal(serializedJSON)

	assert.Nil(t, err)

	assert.Equal(t, userRole, result)
}

func TestShouldFailToUnmarshalUserRoleWhenResponseIsAJsonArray(t *testing.T) {
	response := `["foo","bar"]`

	_, err := NewUserRoleUnmarshaller().Unmarshal([]byte(response))

	assert.NotNil(t, err)
}

func TestShouldFailToUnmarshalUserRoleWhenResponseIsNotAJsonMessage(t *testing.T) {
	response := `foo bar`

	_, err := NewUserRoleUnmarshaller().Unmarshal([]byte(response))

	assert.NotNil(t, err)
}

func TestShouldReturnEmptyUserRoleWhenJsonObjectIsReturnWhereNoFiledMatches(t *testing.T) {
	response := `{"foo" : "bar" }`

	result, err := NewUserRoleUnmarshaller().Unmarshal([]byte(response))

	assert.Nil(t, err)
	assert.Equal(t, UserRole{}, result)
}
