package restapi_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/assert"
)

const (
	userRoleID   = "user-role-id"
	userRoleName = "user-role-name"
)

func TestValidMinimalUserRole(t *testing.T) {
	userRole := UserRole{
		ID:   userRoleID,
		Name: userRoleName,
	}

	assert.Equal(t, userRoleID, userRole.GetID())

	err := userRole.Validate()
	assert.Nil(t, err)
}

func TestValidFullUserRole(t *testing.T) {
	userRole := UserRole{
		ID:                                userRoleID,
		Name:                              userRoleName,
		ImplicitViewFilter:                "Test view filter",
		CanConfigureServiceMapping:        true,
		CanConfigureEumApplications:       true,
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
	}

	assert.Equal(t, userRoleID, userRole.GetID())

	err := userRole.Validate()
	assert.Nil(t, err)
}

func TestInvalidUserRoleBecauseOfMissingId(t *testing.T) {
	userRole := UserRole{
		Name: userRoleName,
	}

	err := userRole.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "ID")
}

func TestInvalidUserRoleBecauseOfMissingName(t *testing.T) {
	userRole := UserRole{
		ID: userRoleID,
	}

	err := userRole.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Name")
}
