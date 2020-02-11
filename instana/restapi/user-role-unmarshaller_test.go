package restapi_test

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

func TestShouldSuccessfullyUnmarshalUserRole(t *testing.T) {
	userRole := UserRole{
		ID:                                "role-id",
		Name:                              "role-name",
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

	serializedJSON, _ := json.Marshal(userRole)

	result, err := NewUserRoleUnmarshaller().Unmarshal(serializedJSON)

	if err != nil {
		t.Fatal("Expected user role to be successfully unmarshalled")
	}

	if !cmp.Equal(result, userRole) {
		t.Fatalf("Expected user role to be properly unmarshalled, %s", cmp.Diff(result, userRole))
	}
}

func TestShouldFailToUnmarshalUserRoleWhenResponseIsAJsonArray(t *testing.T) {
	response := `["foo","bar"]`

	_, err := NewUserRoleUnmarshaller().Unmarshal([]byte(response))

	if err == nil {
		t.Fatal("Expected unmarshalling to fail")
	}
}

func TestShouldFailToUnmarshalUserRoleWhenResponseIsNotAJsonMessage(t *testing.T) {
	response := `foo bar`

	_, err := NewUserRoleUnmarshaller().Unmarshal([]byte(response))

	if err == nil {
		t.Fatal("Expected unmarshalling to fail")
	}
}

func TestShouldReturnEmptyUserRoleWhenJsonObjectIsReturnWhereNoFiledMatches(t *testing.T) {
	response := `{"foo" : "bar" }`

	result, err := NewUserRoleUnmarshaller().Unmarshal([]byte(response))

	if err != nil {
		t.Fatalf("Expected to successfully unmarshal user role response, %s", err)
	}

	if !cmp.Equal(result, UserRole{}) {
		t.Fatal("Expected empty user role")
	}
}
