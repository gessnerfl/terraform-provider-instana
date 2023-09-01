package restapi_test

import (
	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	groupId                = "group-id"
	groupName              = "group-name"
	groupPermissionScopeID = "group-permission-scope-id"
)

func TestShouldReturnSupportedInstanaPermissionsAsString(t *testing.T) {
	expectedResult := []string{
		"CAN_CONFIGURE_APPLICATIONS",
		"CAN_SEE_ON_PREM_LICENE_INFORMATION",
		"CAN_CONFIGURE_EUM_APPLICATIONS",
		"CAN_CONFIGURE_AGENTS",
		"CAN_VIEW_TRACE_DETAILS",
		"CAN_VIEW_LOGS",
		"CAN_CONFIGURE_SESSION_SETTINGS",
		"CAN_CONFIGURE_INTEGRATIONS",
		"CAN_CONFIGURE_GLOBAL_ALERT_CONFIGS",
		"CAN_CONFIGURE_GLOBAL_ALERT_PAYLOAD",
		"CAN_CONFIGURE_MOBILE_APP_MONITORING",
		"CAN_CONFIGURE_API_TOKENS",
		"CAN_CONFIGURE_SERVICE_LEVEL_INDICATORS",
		"CAN_CONFIGURE_AUTHENTICATION_METHODS",
		"CAN_CONFIGURE_RELEASES",
		"CAN_VIEW_AUDIT_LOG",
		"CAN_CONFIGURE_CUSTOM_ALERTS",
		"CAN_CONFIGURE_AGENT_RUN_MODE",
		"CAN_CONFIGURE_SERVICE_MAPPING",
		"CAN_SEE_USAGE_INFORMATION",
		"CAN_EDIT_ALL_ACCESSIBLE_CUSTOM_DASHBOARDS",
		"CAN_CONFIGURE_USERS",
		"CAN_INSTALL_NEW_AGENTS",
		"CAN_CONFIGURE_TEAMS",
		"CAN_CREATE_PUBLIC_CUSTOM_DASHBOARDS",
		"CAN_CONFIGURE_LOG_MANAGEMENT",
		"CAN_VIEW_ACCOUNT_AND_BILLING_INFORMATION",
	}
	assert.Equal(t, expectedResult, SupportedInstanaPermissions.ToStringSlice())
}

func TestShouldReturnIDOfGroupAsIDForAPIPaths(t *testing.T) {
	group := Group{
		ID:   groupId,
		Name: groupName,
	}

	assert.Equal(t, groupId, group.GetIDForResourcePath())
}

func TestShouldReturnTrueWhenPermissionSetIsEmpty(t *testing.T) {
	p := APIPermissionSetWithRoles{}

	require.True(t, p.IsEmpty())

	emptyScopeBinding := ScopeBinding{}
	p = APIPermissionSetWithRoles{
		InfraDFQFilter: &emptyScopeBinding,
	}

	require.True(t, p.IsEmpty())

	emptyScopeBindingSlice := make([]ScopeBinding, 0)
	p = APIPermissionSetWithRoles{
		ApplicationIDs:          emptyScopeBindingSlice,
		KubernetesNamespaceUIDs: emptyScopeBindingSlice,
		KubernetesClusterUUIDs:  emptyScopeBindingSlice,
		InfraDFQFilter:          &emptyScopeBinding,
		MobileAppIDs:            emptyScopeBindingSlice,
		WebsiteIDs:              emptyScopeBindingSlice,
		Permissions:             make([]InstanaPermission, 0),
	}
}

func TestShouldReturnFalseWhenPermissionSetIsNotEmptyWhenApplicationIDsAreSet(t *testing.T) {
	p := APIPermissionSetWithRoles{
		ApplicationIDs: []ScopeBinding{{ScopeID: groupPermissionScopeID}},
	}
	require.False(t, p.IsEmpty())
}

func TestShouldReturnFalseWhenPermissionSetIsNotEmptyWhenKubernetesClusterUUIDsAreSet(t *testing.T) {
	p := APIPermissionSetWithRoles{
		KubernetesClusterUUIDs: []ScopeBinding{{ScopeID: groupPermissionScopeID}},
	}
	require.False(t, p.IsEmpty())
}

func TestShouldReturnFalseWhenPermissionSetIsNotEmptyWhenKubernetesNamespaceUIDsAreSet(t *testing.T) {
	p := APIPermissionSetWithRoles{
		KubernetesNamespaceUIDs: []ScopeBinding{{ScopeID: groupPermissionScopeID}},
	}
	require.False(t, p.IsEmpty())
}

func TestShouldReturnFalseWhenPermissionSetIsNotEmptyWhenMobileAppIDsAreSet(t *testing.T) {
	p := APIPermissionSetWithRoles{
		MobileAppIDs: []ScopeBinding{{ScopeID: groupPermissionScopeID}},
	}
	require.False(t, p.IsEmpty())
}

func TestShouldReturnFalseWhenPermissionSetIsNotEmptyWhenWebsiteIDsAreSet(t *testing.T) {
	p := APIPermissionSetWithRoles{
		WebsiteIDs: []ScopeBinding{{ScopeID: groupPermissionScopeID}},
	}
	require.False(t, p.IsEmpty())
}

func TestShouldReturnFalseWhenPermissionSetIsNotEmptyWhenPermissionsAreSet(t *testing.T) {
	p := APIPermissionSetWithRoles{
		Permissions: []InstanaPermission{PermissionCanConfigureApplications},
	}
	require.False(t, p.IsEmpty())
}

func TestShouldReturnFalseWhenPermissionSetIsNotEmptyWhenInfrastructureDFQIsSet(t *testing.T) {
	scopeBinding := ScopeBinding{ScopeID: groupPermissionScopeID}
	p := APIPermissionSetWithRoles{
		InfraDFQFilter: &scopeBinding,
	}
	require.False(t, p.IsEmpty())
}
