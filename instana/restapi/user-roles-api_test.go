package restapi_test

import (
	"strings"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

func TestValidMinimalUserRole(t *testing.T) {
	userRole := UserRole{
		ID:   "test-user-role-id",
		Name: "Test User Role",
	}

	if "test-user-role-id" != userRole.GetID() {
		t.Fatalf("Expected to get correct ID but got %s", userRole.GetID())
	}

	if err := userRole.Validate(); err != nil {
		t.Fatalf("Expected valid user role but got validation error %s", err)
	}
}

func TestValidFullUserRole(t *testing.T) {
	userRole := UserRole{
		ID:                                "test-user-role-id",
		Name:                              "Test User Role",
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

	if "test-user-role-id" != userRole.GetID() {
		t.Fatalf("Expected to get correct ID but got %s", userRole.GetID())
	}

	if err := userRole.Validate(); err != nil {
		t.Fatalf("Expected valid user role but got validation error %s", err)
	}
}

func TestInvalidUserRoleBecauseOfMissingId(t *testing.T) {
	userRole := UserRole{
		Name: "Test userRole",
	}

	if err := userRole.Validate(); err == nil || !strings.Contains(err.Error(), "ID") {
		t.Fatalf("Expected invalid userRole because of missing ID")
	}
}

func TestInvalidUserRoleBecauseOfMissingName(t *testing.T) {
	userRole := UserRole{
		ID: "test-user-role-id",
	}

	if err := userRole.Validate(); err == nil || !strings.Contains(err.Error(), "Name") {
		t.Fatalf("Expected invalid userRole because of missing Name")
	}
}
