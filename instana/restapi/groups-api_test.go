package restapi_test

import (
	"fmt"
	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	groupId                    = "group-id"
	groupName                  = "group-name"
	groupMember1ID             = "group-member-1-id"
	groupMember1Email          = "group-member-1-email"
	groupMember2ID             = "group-member-2-id"
	groupPermissionScopeID     = "group-permission-scope-id"
	groupPermissionScopeRoleID = "group-permission-scope-role-id"
)

func TestShouldSuccessfullyValidateGroup(t *testing.T) {
	scopeBinding := ScopeBinding{ScopeID: groupPermissionScopeID}
	scopeBindings := []ScopeBinding{scopeBinding}
	group := Group{
		ID:      groupId,
		Name:    groupName,
		Members: []APIMember{{UserID: groupMember1ID}, {UserID: groupMember2ID}},
		PermissionSet: APIPermissionSetWithRoles{
			ApplicationIDs:          scopeBindings,
			InfraDFQFilter:          &scopeBinding,
			KubernetesClusterUUIDs:  scopeBindings,
			KubernetesNamespaceUIDs: scopeBindings,
			MobileAppIDs:            scopeBindings,
			WebsiteIDs:              scopeBindings,
			Permissions:             []InstanaPermission{PermissionCanConfigureAgents, PermissionCanConfigureCustomAlerts},
		},
	}

	err := group.Validate()

	require.NoError(t, err)
}

func TestShouldSuccessfullyValidateMinimalGroup(t *testing.T) {
	group := Group{
		ID:   groupId,
		Name: groupName,
	}

	err := group.Validate()

	require.NoError(t, err)
}

func TestShouldFailToValidateGroupWhenNameIsMissing(t *testing.T) {
	group := Group{
		ID: groupId,
	}

	err := group.Validate()

	require.Error(t, err)
	require.Contains(t, err.Error(), "name is missing")
}

func TestShouldFailToValidateGroupWhenNameIsBlank(t *testing.T) {
	group := Group{
		ID:   groupId,
		Name: " ",
	}

	err := group.Validate()

	require.Error(t, err)
	require.Contains(t, err.Error(), "name is missing")
}

func TestShouldSuccessfullyValidateGroupWithFullMemberDefinition(t *testing.T) {
	email := groupMember1Email
	group := Group{
		ID:      groupId,
		Name:    groupName,
		Members: []APIMember{{UserID: groupMember1ID, Email: &email}},
	}

	err := group.Validate()

	require.NoError(t, err)
}

func TestShouldFailToValidateGroupWithGroupMemberWithoutUserID(t *testing.T) {
	email := groupMember1Email
	group := Group{
		ID:      groupId,
		Name:    groupName,
		Members: []APIMember{{Email: &email}},
	}

	err := group.Validate()

	require.Error(t, err)
	require.Contains(t, err.Error(), "userId of group member")
}

func TestShouldSuccessfullyValidateFullScopeBindingOfInfrastructureDFQ(t *testing.T) {
	roleID := groupPermissionScopeRoleID
	scopeBinding := ScopeBinding{ScopeID: groupPermissionScopeID, ScopeRoleID: &roleID}
	group := Group{
		ID:   groupId,
		Name: groupName,
		PermissionSet: APIPermissionSetWithRoles{
			InfraDFQFilter: &scopeBinding,
		},
	}

	err := group.Validate()

	require.NoError(t, err)
}

func TestShouldSuccessfullyValidateScopeBindingOfInfrastructureDFQWhenScopeIDIsMissing(t *testing.T) {
	roleID := groupPermissionScopeRoleID
	scopeBinding := ScopeBinding{ScopeRoleID: &roleID}
	group := Group{
		ID:   groupId,
		Name: groupName,
		PermissionSet: APIPermissionSetWithRoles{
			InfraDFQFilter: &scopeBinding,
		},
	}

	err := group.Validate()

	require.NoError(t, err)
}

func TestShouldSuccessfullyValidateScopeBindingOfInfrastructureDFQWhenScopeIDIsBlank(t *testing.T) {
	roleID := groupPermissionScopeRoleID
	scopeBinding := ScopeBinding{ScopeID: " ", ScopeRoleID: &roleID}
	group := Group{
		ID:   groupId,
		Name: groupName,
		PermissionSet: APIPermissionSetWithRoles{
			InfraDFQFilter: &scopeBinding,
		},
	}

	err := group.Validate()

	require.NoError(t, err)
}

func TestSliceOfScopeBindings(t *testing.T) {
	testSet := map[string]func(ps *APIPermissionSetWithRoles, bindings *[]ScopeBinding){
		"ApplicationIDs":          func(ps *APIPermissionSetWithRoles, bindings *[]ScopeBinding) { ps.ApplicationIDs = *bindings },
		"KubernetesClusterUUIDs":  func(ps *APIPermissionSetWithRoles, bindings *[]ScopeBinding) { ps.KubernetesClusterUUIDs = *bindings },
		"KubernetesNamespaceUIDs": func(ps *APIPermissionSetWithRoles, bindings *[]ScopeBinding) { ps.KubernetesNamespaceUIDs = *bindings },
		"MobileAppIDs":            func(ps *APIPermissionSetWithRoles, bindings *[]ScopeBinding) { ps.MobileAppIDs = *bindings },
		"WebsiteIDs":              func(ps *APIPermissionSetWithRoles, bindings *[]ScopeBinding) { ps.WebsiteIDs = *bindings },
	}

	testShouldSuccessfulValidateScopeBindings(t, testSet)
	testShouldSuccessfullyValidateScopeBindingsWhenScopeIDIsMissingOfBlank(t, testSet)
}

func testShouldSuccessfulValidateScopeBindings(t *testing.T, testSet map[string]func(ps *APIPermissionSetWithRoles, bindings *[]ScopeBinding)) {
	for attribute, mapping := range testSet {
		t.Run(fmt.Sprintf("TestShouldSuccessfullyValidateScopeBindingsOf%sOfPermissionSetOfGroup", attribute), func(t *testing.T) {
			roleID := groupPermissionScopeRoleID
			scopeBindings := []ScopeBinding{{ScopeID: groupPermissionScopeID, ScopeRoleID: &roleID}}
			permissionSet := APIPermissionSetWithRoles{}
			mapping(&permissionSet, &scopeBindings)
			group := Group{
				ID:            groupId,
				Name:          groupName,
				PermissionSet: permissionSet,
			}

			err := group.Validate()

			require.NoError(t, err)
		})
	}
}

func testShouldSuccessfullyValidateScopeBindingsWhenScopeIDIsMissingOfBlank(t *testing.T, testSet map[string]func(ps *APIPermissionSetWithRoles, bindings *[]ScopeBinding)) {
	for attribute, mapping := range testSet {
		for testCase, scopeID := range map[string]string{"Missing": "", "Blank": ""} {
			t.Run(fmt.Sprintf("TestShouldFailToValidateScopeBindingsOf%sOfPermissionSetOfGroupWhenIDIs%s", attribute, testCase), func(t *testing.T) {
				roleID := groupPermissionScopeRoleID
				scopeBindings := []ScopeBinding{{ScopeID: scopeID, ScopeRoleID: &roleID}}
				permissionSet := APIPermissionSetWithRoles{}
				mapping(&permissionSet, &scopeBindings)
				group := Group{
					ID:            groupId,
					Name:          groupName,
					PermissionSet: permissionSet,
				}

				err := group.Validate()

				require.NoError(t, err)
			})
		}
	}
}

func TestShouldSuccessfullyValidatePermissionsOfGroupWhenPermissionsAreSupported(t *testing.T) {
	group := Group{
		ID:   groupId,
		Name: groupName,
		PermissionSet: APIPermissionSetWithRoles{
			Permissions: SupportedInstanaPermissions,
		},
	}

	err := group.Validate()

	require.NoError(t, err)
}

func TestShouldFailToValidatePermissionsOfGroupWhenPermissionIsNotSupported(t *testing.T) {
	group := Group{
		ID:   groupId,
		Name: groupName,
		PermissionSet: APIPermissionSetWithRoles{
			Permissions: []InstanaPermission{InstanaPermission("invalid")},
		},
	}

	err := group.Validate()

	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid is not a supported")
}

func TestShouldReturnTrueWhenTheProvidedInstanaPermissionIsSupported(t *testing.T) {
	for _, v := range SupportedInstanaPermissions {
		t.Run(fmt.Sprintf("TestShouldReturnTrueWhenCheckingIfPermission%sIsSupported", v), func(t *testing.T) {
			require.True(t, SupportedInstanaPermissions.IsSupported(v))
		})
	}
}

func TestShouldReturnFalseWhenTheProvidedInstanaPermissionIsSupported(t *testing.T) {
	require.False(t, SupportedInstanaPermissions.IsSupported("invalid"))
}

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
