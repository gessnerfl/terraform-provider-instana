package instana_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/stretchr/testify/assert"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/gessnerfl/terraform-provider-instana/utils"
)

const resourceUserRoleDefinitionTemplate = `
provider "instana" {
  api_token = "test-token"
  endpoint = "localhost:{{PORT}}"
}

resource "instana_user_role" "example" {
  name = "name"
  can_configure_service_mapping = true
  can_configure_eum_applications = true
  can_configure_mobile_app_monitoring = true
  can_configure_users = true
  can_install_new_agents = true
  can_see_usage_information = true
  can_configure_integrations = true
  can_see_on_premise_license_information = true
  can_configure_roles = true
  can_configure_custom_alerts = true
  can_configure_api_tokens = true
  can_configure_agent_run_mode = true
  can_view_audit_log = true
  can_configure_objectives = true
  can_configure_agents = true
  can_configure_authentication_methods = true
  can_configure_applications = true
  can_configure_teams = true
  restricted_access = true
  can_configure_releases = true
  can_configure_log_management = true
  can_create_public_custom_dashboards = true
  can_view_logs = true
  can_view_trace_details = true
}
`

const userRoleApiPath = restapi.UserRolesResourcePath + "/{id}"
const testUserRoleDefinition = "instana_user_role.example"
const valueTrue = "true"
const userRoleID = "user-role-id"
const viewFilterFieldValue = "view filter"
const userRoleNameFieldValue = "name"

var userRolePermissionFields = []string{
	UserRoleFieldCanConfigureServiceMapping,
	UserRoleFieldCanConfigureEumApplications,
	UserRoleFieldCanConfigureMobileAppMonitoring,
	UserRoleFieldCanConfigureUsers,
	UserRoleFieldCanInstallNewAgents,
	UserRoleFieldCanSeeUsageInformation,
	UserRoleFieldCanConfigureIntegrations,
	UserRoleFieldCanSeeOnPremiseLicenseInformation,
	UserRoleFieldCanConfigureRoles,
	UserRoleFieldCanConfigureCustomAlerts,
	UserRoleFieldCanConfigureAPITokens,
	UserRoleFieldCanConfigureAgentRunMode,
	UserRoleFieldCanViewAuditLog,
	UserRoleFieldCanConfigureObjectives,
	UserRoleFieldCanConfigureAgents,
	UserRoleFieldCanConfigureAuthenticationMethods,
	UserRoleFieldCanConfigureApplications,
	UserRoleFieldCanConfigureTeams,
	UserRoleFieldRestrictedAccess,
	UserRoleFieldCanConfigureReleases,
	UserRoleFieldCanConfigureLogManagement,
	UserRoleFieldCanCreatePublicCustomDashboards,
	UserRoleFieldCanViewLogs,
	UserRoleFieldCanViewTraceDetails,
}

func TestCRUDOfUserRoleResourceWithMockServer(t *testing.T) {
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPut, userRoleApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodDelete, userRoleApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodGet, userRoleApiPath, func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		json := strings.ReplaceAll(`
		{
			"id" : "{{id}}",
			"name" : "name",
			"canConfigureServiceMapping" : true,
			"canConfigureEumApplications" : true,
			"canConfigureMobileAppMonitoring" : true,
			"canConfigureUsers" : true,
			"canInstallNewAgents" : true,
			"canSeeUsageInformation" : true,
			"canConfigureIntegrations" : true,
			"canSeeOnPremLicenseInformation" : true,
			"canConfigureRoles" : true,
			"canConfigureCustomAlerts" : true,
			"canConfigureApiTokens" : true,
			"canConfigureAgentRunMode" : true,
			"canViewAuditLog" : true,
			"canConfigureObjectives" : true,
			"canConfigureAgents" : true,
			"canConfigureAuthenticationMethods" : true,
			"canConfigureApplications" : true,
			"canConfigureTeams" : true,
			"restrictedAccess" : true,
			"canConfigureReleases" : true,
			"canConfigureLogManagement" : true,
			"canCreatePublicCustomDashboards" : true,
			"canViewLogs" : true,
			"canViewTraceDetails" : true
		}
		`, "{{id}}", vars["id"])
		w.Header().Set(contentType, r.Header.Get(contentType))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()

	resourceUserRoleDefinition := strings.ReplaceAll(resourceUserRoleDefinitionTemplate, "{{PORT}}", strconv.Itoa(httpServer.GetPort()))

	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: resourceUserRoleDefinition,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testUserRoleDefinition, "id"),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldName, userRoleNameFieldValue),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanConfigureServiceMapping, valueTrue),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanConfigureEumApplications, valueTrue),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanConfigureMobileAppMonitoring, valueTrue),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanConfigureUsers, valueTrue),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanInstallNewAgents, valueTrue),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanSeeUsageInformation, valueTrue),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanConfigureIntegrations, valueTrue),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanSeeOnPremiseLicenseInformation, valueTrue),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanConfigureRoles, valueTrue),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanConfigureCustomAlerts, valueTrue),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanConfigureAPITokens, valueTrue),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanConfigureAgentRunMode, valueTrue),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanViewAuditLog, valueTrue),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanConfigureObjectives, valueTrue),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanConfigureAgents, valueTrue),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanConfigureAuthenticationMethods, valueTrue),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanConfigureApplications, valueTrue),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanConfigureTeams, valueTrue),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldRestrictedAccess, valueTrue),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanConfigureReleases, valueTrue),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanConfigureLogManagement, valueTrue),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanCreatePublicCustomDashboards, valueTrue),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanViewLogs, valueTrue),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanViewTraceDetails, valueTrue),
				),
			},
		},
	})
}

func TestUserRoleSchemaDefinitionIsValid(t *testing.T) {
	schema := NewUserRoleResourceHandle().MetaData().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schema, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(UserRoleFieldName)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(UserRoleFieldCanConfigureServiceMapping, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(UserRoleFieldCanConfigureEumApplications, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(UserRoleFieldCanConfigureMobileAppMonitoring, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(UserRoleFieldCanConfigureUsers, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(UserRoleFieldCanInstallNewAgents, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(UserRoleFieldCanSeeUsageInformation, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(UserRoleFieldCanConfigureIntegrations, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(UserRoleFieldCanSeeOnPremiseLicenseInformation, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(UserRoleFieldCanConfigureRoles, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(UserRoleFieldCanConfigureCustomAlerts, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(UserRoleFieldCanConfigureAPITokens, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(UserRoleFieldCanConfigureAgentRunMode, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(UserRoleFieldCanViewAuditLog, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(UserRoleFieldCanConfigureObjectives, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(UserRoleFieldCanConfigureAgents, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(UserRoleFieldCanConfigureAuthenticationMethods, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(UserRoleFieldCanConfigureApplications, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(UserRoleFieldCanConfigureTeams, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(UserRoleFieldRestrictedAccess, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(UserRoleFieldCanConfigureReleases, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(UserRoleFieldCanConfigureLogManagement, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(UserRoleFieldCanCreatePublicCustomDashboards, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(UserRoleFieldCanViewLogs, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(UserRoleFieldCanViewTraceDetails, false)
}

func TestUserRoleResourceShouldHaveSchemaVersionOne(t *testing.T) {
	assert.Equal(t, 1, NewUserRoleResourceHandle().MetaData().SchemaVersion)
}

func TestUserRoleResourceShouldHaveOneStateUpgraderForVersion0(t *testing.T) {
	assert.Equal(t, 1, len(NewUserRoleResourceHandle().StateUpgraders()))
	assert.Equal(t, 0, NewUserRoleResourceHandle().StateUpgraders()[0].Version)
}

func TestShouldDeleteValueOfImplicitViewFilterFieldWhenMigratingToVersion1AndValueIsSet(t *testing.T) {
	field := "implicit_view_filter"
	rawData := make(map[string]interface{})
	rawData[field] = "value"
	meta := "dummy"

	result, err := NewUserRoleResourceHandle().StateUpgraders()[0].Upgrade(rawData, meta)

	assert.Nil(t, err)
	_, found := result[field]
	assert.False(t, found)
}

func TestShouldDoNothingWhenMigratingToVersion1AndImplicitViewFilterValueIsSet(t *testing.T) {
	rawData := make(map[string]interface{})
	meta := "dummy"

	result, err := NewUserRoleResourceHandle().StateUpgraders()[0].Upgrade(rawData, meta)

	assert.Nil(t, err)
	assert.Equal(t, rawData, result)
}

func TestShouldReturnCorrectResourceNameForUserroleResource(t *testing.T) {
	name := NewUserRoleResourceHandle().MetaData().ResourceName

	assert.Equal(t, name, "instana_user_role")
}

func TestShouldUpdateBasicFieldsOfTerraformResourceStateFromModelForUserRole(t *testing.T) {
	testHelper := NewTestHelper(t)
	sut := NewUserRoleResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)
	userRole := restapi.UserRole{
		ID:   userRoleID,
		Name: userRoleNameFieldValue,
	}

	err := sut.UpdateState(resourceData, &userRole)

	assert.Nil(t, err)
	assert.Equal(t, userRoleID, resourceData.Id())
	assert.Equal(t, userRoleNameFieldValue, resourceData.Get(UserRoleFieldName))
	assert.False(t, resourceData.Get(UserRoleFieldCanConfigureServiceMapping).(bool))
	assert.False(t, resourceData.Get(UserRoleFieldCanConfigureEumApplications).(bool))
	assert.False(t, resourceData.Get(UserRoleFieldCanConfigureMobileAppMonitoring).(bool))
	assert.False(t, resourceData.Get(UserRoleFieldCanConfigureUsers).(bool))
	assert.False(t, resourceData.Get(UserRoleFieldCanInstallNewAgents).(bool))
	assert.False(t, resourceData.Get(UserRoleFieldCanSeeUsageInformation).(bool))
	assert.False(t, resourceData.Get(UserRoleFieldCanConfigureIntegrations).(bool))
	assert.False(t, resourceData.Get(UserRoleFieldCanSeeOnPremiseLicenseInformation).(bool))
	assert.False(t, resourceData.Get(UserRoleFieldCanConfigureRoles).(bool))
	assert.False(t, resourceData.Get(UserRoleFieldCanConfigureCustomAlerts).(bool))
	assert.False(t, resourceData.Get(UserRoleFieldCanConfigureAPITokens).(bool))
	assert.False(t, resourceData.Get(UserRoleFieldCanConfigureAgentRunMode).(bool))
	assert.False(t, resourceData.Get(UserRoleFieldCanViewAuditLog).(bool))
	assert.False(t, resourceData.Get(UserRoleFieldCanConfigureObjectives).(bool))
	assert.False(t, resourceData.Get(UserRoleFieldCanConfigureAgents).(bool))
	assert.False(t, resourceData.Get(UserRoleFieldCanConfigureAuthenticationMethods).(bool))
	assert.False(t, resourceData.Get(UserRoleFieldCanConfigureApplications).(bool))
	assert.False(t, resourceData.Get(UserRoleFieldCanConfigureTeams).(bool))
	assert.False(t, resourceData.Get(UserRoleFieldRestrictedAccess).(bool))
	assert.False(t, resourceData.Get(UserRoleFieldCanConfigureReleases).(bool))
	assert.False(t, resourceData.Get(UserRoleFieldCanConfigureLogManagement).(bool))
	assert.False(t, resourceData.Get(UserRoleFieldCanCreatePublicCustomDashboards).(bool))
	assert.False(t, resourceData.Get(UserRoleFieldCanViewLogs).(bool))
	assert.False(t, resourceData.Get(UserRoleFieldCanViewTraceDetails).(bool))
}

func TestShouldUpdateCanConfigureServiceMappingPermissionOfTerraformResourceStateFromModelForUserRole(t *testing.T) {
	userRole := restapi.UserRole{
		ID:                         userRoleID,
		Name:                       userRoleNameFieldValue,
		CanConfigureServiceMapping: true,
	}

	testSingleUserRolePermissionSet(t, userRole, UserRoleFieldCanConfigureServiceMapping)
}

func TestShouldUpdateCanConfigureEumApplicationsPermissionOfTerraformResourceStateFromModelForUserRole(t *testing.T) {
	userRole := restapi.UserRole{
		ID:                          userRoleID,
		Name:                        userRoleNameFieldValue,
		CanConfigureEumApplications: true,
	}

	testSingleUserRolePermissionSet(t, userRole, UserRoleFieldCanConfigureEumApplications)
}

func TestShouldUpdateCanConfigureMobileAppMonitoringPermissionOfTerraformResourceStateFromModelForUserRole(t *testing.T) {
	userRole := restapi.UserRole{
		ID:                              userRoleID,
		Name:                            userRoleNameFieldValue,
		CanConfigureMobileAppMonitoring: true,
	}

	testSingleUserRolePermissionSet(t, userRole, UserRoleFieldCanConfigureMobileAppMonitoring)
}

func TestShouldUpdateCanConfigureUsersPermissionOfTerraformResourceStateFromModelForUserRole(t *testing.T) {
	userRole := restapi.UserRole{
		ID:                userRoleID,
		Name:              userRoleNameFieldValue,
		CanConfigureUsers: true,
	}

	testSingleUserRolePermissionSet(t, userRole, UserRoleFieldCanConfigureUsers)
}

func TestShouldUpdateCanInstallNewAgentsPermissionOfTerraformResourceStateFromModelForUserRole(t *testing.T) {
	userRole := restapi.UserRole{
		ID:                  userRoleID,
		Name:                userRoleNameFieldValue,
		CanInstallNewAgents: true,
	}

	testSingleUserRolePermissionSet(t, userRole, UserRoleFieldCanInstallNewAgents)
}

func TestShouldUpdateCanSeeUsageInformationPermissionOfTerraformResourceStateFromModelForUserRole(t *testing.T) {
	userRole := restapi.UserRole{
		ID:                     userRoleID,
		Name:                   userRoleNameFieldValue,
		CanSeeUsageInformation: true,
	}

	testSingleUserRolePermissionSet(t, userRole, UserRoleFieldCanSeeUsageInformation)
}

func TestShouldUpdateCanConfigureIntegrationsPermissionOfTerraformResourceStateFromModelForUserRole(t *testing.T) {
	userRole := restapi.UserRole{
		ID:                       userRoleID,
		Name:                     userRoleNameFieldValue,
		CanConfigureIntegrations: true,
	}

	testSingleUserRolePermissionSet(t, userRole, UserRoleFieldCanConfigureIntegrations)
}

func TestShouldUpdateCanSeeOnPremiseLicenseInformationPermissionOfTerraformResourceStateFromModelForUserRole(t *testing.T) {
	userRole := restapi.UserRole{
		ID:                                userRoleID,
		Name:                              userRoleNameFieldValue,
		CanSeeOnPremiseLicenseInformation: true,
	}

	testSingleUserRolePermissionSet(t, userRole, UserRoleFieldCanSeeOnPremiseLicenseInformation)
}

func TestShouldUpdateCanConfigureRolesPermissionOfTerraformResourceStateFromModelForUserRole(t *testing.T) {
	userRole := restapi.UserRole{
		ID:                userRoleID,
		Name:              userRoleNameFieldValue,
		CanConfigureRoles: true,
	}

	testSingleUserRolePermissionSet(t, userRole, UserRoleFieldCanConfigureRoles)
}

func TestShouldUpdateCanConfigureCustomAlertsPermissionOfTerraformResourceStateFromModelForUserRole(t *testing.T) {
	userRole := restapi.UserRole{
		ID:                       userRoleID,
		Name:                     userRoleNameFieldValue,
		CanConfigureCustomAlerts: true,
	}

	testSingleUserRolePermissionSet(t, userRole, UserRoleFieldCanConfigureCustomAlerts)
}

func TestShouldUpdateCanConfigureAPITokensPermissionOfTerraformResourceStateFromModelForUserRole(t *testing.T) {
	userRole := restapi.UserRole{
		ID:                    userRoleID,
		Name:                  userRoleNameFieldValue,
		CanConfigureAPITokens: true,
	}

	testSingleUserRolePermissionSet(t, userRole, UserRoleFieldCanConfigureAPITokens)
}

func TestShouldUpdateCanConfigureAgentRunModePermissionOfTerraformResourceStateFromModelForUserRole(t *testing.T) {
	userRole := restapi.UserRole{
		ID:                       userRoleID,
		Name:                     userRoleNameFieldValue,
		CanConfigureAgentRunMode: true,
	}

	testSingleUserRolePermissionSet(t, userRole, UserRoleFieldCanConfigureAgentRunMode)
}

func TestShouldUpdateCanViewAuditLogPermissionOfTerraformResourceStateFromModelForUserRole(t *testing.T) {
	userRole := restapi.UserRole{
		ID:              userRoleID,
		Name:            userRoleNameFieldValue,
		CanViewAuditLog: true,
	}

	testSingleUserRolePermissionSet(t, userRole, UserRoleFieldCanViewAuditLog)
}

func TestShouldUpdateCanConfigureObjectivesPermissionOfTerraformResourceStateFromModelForUserRole(t *testing.T) {
	userRole := restapi.UserRole{
		ID:                     userRoleID,
		Name:                   userRoleNameFieldValue,
		CanConfigureObjectives: true,
	}

	testSingleUserRolePermissionSet(t, userRole, UserRoleFieldCanConfigureObjectives)
}

func TestShouldUpdateCanConfigureAgentsPermissionOfTerraformResourceStateFromModelForUserRole(t *testing.T) {
	userRole := restapi.UserRole{
		ID:                 userRoleID,
		Name:               userRoleNameFieldValue,
		CanConfigureAgents: true,
	}

	testSingleUserRolePermissionSet(t, userRole, UserRoleFieldCanConfigureAgents)
}

func TestShouldUpdateCanConfigureAuthenticationMethodsPermissionOfTerraformResourceStateFromModelForUserRole(t *testing.T) {
	userRole := restapi.UserRole{
		ID:                                userRoleID,
		Name:                              userRoleNameFieldValue,
		CanConfigureAuthenticationMethods: true,
	}

	testSingleUserRolePermissionSet(t, userRole, UserRoleFieldCanConfigureAuthenticationMethods)
}

func TestShouldUpdateCanConfigureApplicationsPermissionOfTerraformResourceStateFromModelForUserRole(t *testing.T) {
	userRole := restapi.UserRole{
		ID:                       userRoleID,
		Name:                     userRoleNameFieldValue,
		CanConfigureApplications: true,
	}

	testSingleUserRolePermissionSet(t, userRole, UserRoleFieldCanConfigureApplications)
}

func TestShouldUpdateCanConfigureTeamsPermissionOfTerraformResourceStateFromModelForUserRole(t *testing.T) {
	userRole := restapi.UserRole{
		ID:                userRoleID,
		Name:              userRoleNameFieldValue,
		CanConfigureTeams: true,
	}

	testSingleUserRolePermissionSet(t, userRole, UserRoleFieldCanConfigureTeams)
}

func TestShouldUpdateRestrictedAccessPermissionOfTerraformResourceStateFromModelForUserRole(t *testing.T) {
	userRole := restapi.UserRole{
		ID:               userRoleID,
		Name:             userRoleNameFieldValue,
		RestrictedAccess: true,
	}

	testSingleUserRolePermissionSet(t, userRole, UserRoleFieldRestrictedAccess)
}

func TestShouldUpdateCanConfigureReleasesPermissionOfTerraformResourceStateFromModelForUserRole(t *testing.T) {
	userRole := restapi.UserRole{
		ID:                   userRoleID,
		Name:                 userRoleNameFieldValue,
		CanConfigureReleases: true,
	}

	testSingleUserRolePermissionSet(t, userRole, UserRoleFieldCanConfigureReleases)
}

func TestShouldUpdateCanConfigureLogManagementPermissionOfTerraformResourceStateFromModelForUserRole(t *testing.T) {
	userRole := restapi.UserRole{
		ID:                        userRoleID,
		Name:                      userRoleNameFieldValue,
		CanConfigureLogManagement: true,
	}

	testSingleUserRolePermissionSet(t, userRole, UserRoleFieldCanConfigureLogManagement)
}

func TestShouldUpdateCanCreatePublicCustomDashboardsPermissionOfTerraformResourceStateFromModelForUserRole(t *testing.T) {
	userRole := restapi.UserRole{
		ID:                              userRoleID,
		Name:                            userRoleNameFieldValue,
		CanCreatePublicCustomDashboards: true,
	}

	testSingleUserRolePermissionSet(t, userRole, UserRoleFieldCanCreatePublicCustomDashboards)
}

func TestShouldUpdateCanViewLogsPermissionOfTerraformResourceStateFromModelForUserRole(t *testing.T) {
	userRole := restapi.UserRole{
		ID:          userRoleID,
		Name:        userRoleNameFieldValue,
		CanViewLogs: true,
	}

	testSingleUserRolePermissionSet(t, userRole, UserRoleFieldCanViewLogs)
}

func TestShouldUpdateCanViewTraceDetailsPermissionOfTerraformResourceStateFromModelForUserRole(t *testing.T) {
	userRole := restapi.UserRole{
		ID:                  userRoleID,
		Name:                userRoleNameFieldValue,
		CanViewTraceDetails: true,
	}

	testSingleUserRolePermissionSet(t, userRole, UserRoleFieldCanViewTraceDetails)
}

func testSingleUserRolePermissionSet(t *testing.T, userRole restapi.UserRole, expectedPermissionField string) {
	testHelper := NewTestHelper(t)
	sut := NewUserRoleResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(sut)

	err := sut.UpdateState(resourceData, &userRole)

	assert.Nil(t, err)
	assert.True(t, resourceData.Get(expectedPermissionField).(bool))
	for _, permissionField := range userRolePermissionFields {
		if permissionField != expectedPermissionField {
			assert.False(t, resourceData.Get(permissionField).(bool))
		}
	}
}

func TestShouldConvertStateOfUserRoleTerraformResourceToDataModel(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewUserRoleResourceHandle()

	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	resourceData.SetId(userRoleID)
	resourceData.Set(UserRoleFieldName, userRoleNameFieldValue)
	resourceData.Set(UserRoleFieldCanConfigureServiceMapping, true)
	resourceData.Set(UserRoleFieldCanConfigureEumApplications, true)
	resourceData.Set(UserRoleFieldCanConfigureMobileAppMonitoring, true)
	resourceData.Set(UserRoleFieldCanConfigureUsers, true)
	resourceData.Set(UserRoleFieldCanInstallNewAgents, true)
	resourceData.Set(UserRoleFieldCanSeeUsageInformation, true)
	resourceData.Set(UserRoleFieldCanConfigureIntegrations, true)
	resourceData.Set(UserRoleFieldCanSeeOnPremiseLicenseInformation, true)
	resourceData.Set(UserRoleFieldCanConfigureRoles, true)
	resourceData.Set(UserRoleFieldCanConfigureCustomAlerts, true)
	resourceData.Set(UserRoleFieldCanConfigureAPITokens, true)
	resourceData.Set(UserRoleFieldCanConfigureAgentRunMode, true)
	resourceData.Set(UserRoleFieldCanViewAuditLog, true)
	resourceData.Set(UserRoleFieldCanConfigureObjectives, true)
	resourceData.Set(UserRoleFieldCanConfigureAgents, true)
	resourceData.Set(UserRoleFieldCanConfigureAuthenticationMethods, true)
	resourceData.Set(UserRoleFieldCanConfigureApplications, true)
	resourceData.Set(UserRoleFieldCanConfigureTeams, true)
	resourceData.Set(UserRoleFieldRestrictedAccess, true)
	resourceData.Set(UserRoleFieldCanConfigureReleases, true)
	resourceData.Set(UserRoleFieldCanConfigureLogManagement, true)
	resourceData.Set(UserRoleFieldCanCreatePublicCustomDashboards, true)
	resourceData.Set(UserRoleFieldCanViewLogs, true)
	resourceData.Set(UserRoleFieldCanViewTraceDetails, true)

	model, err := resourceHandle.MapStateToDataObject(resourceData, utils.NewResourceNameFormatter("prefix ", " suffix"))

	assert.Nil(t, err)
	assert.IsType(t, &restapi.UserRole{}, model, "Model should be an alerting channel")
	assert.Equal(t, userRoleID, model.GetID())
	assert.Equal(t, userRoleNameFieldValue, model.(*restapi.UserRole).Name)
	assert.True(t, model.(*restapi.UserRole).CanConfigureServiceMapping)
	assert.True(t, model.(*restapi.UserRole).CanConfigureEumApplications)
	assert.True(t, model.(*restapi.UserRole).CanConfigureMobileAppMonitoring)
	assert.True(t, model.(*restapi.UserRole).CanConfigureUsers)
	assert.True(t, model.(*restapi.UserRole).CanInstallNewAgents)
	assert.True(t, model.(*restapi.UserRole).CanSeeUsageInformation)
	assert.True(t, model.(*restapi.UserRole).CanConfigureIntegrations)
	assert.True(t, model.(*restapi.UserRole).CanSeeOnPremiseLicenseInformation)
	assert.True(t, model.(*restapi.UserRole).CanConfigureRoles)
	assert.True(t, model.(*restapi.UserRole).CanConfigureCustomAlerts)
	assert.True(t, model.(*restapi.UserRole).CanConfigureAPITokens)
	assert.True(t, model.(*restapi.UserRole).CanConfigureAgentRunMode)
	assert.True(t, model.(*restapi.UserRole).CanViewAuditLog)
	assert.True(t, model.(*restapi.UserRole).CanConfigureObjectives)
	assert.True(t, model.(*restapi.UserRole).CanConfigureAgents)
	assert.True(t, model.(*restapi.UserRole).CanConfigureAuthenticationMethods)
	assert.True(t, model.(*restapi.UserRole).CanConfigureApplications)
	assert.True(t, model.(*restapi.UserRole).CanConfigureTeams)
	assert.True(t, model.(*restapi.UserRole).RestrictedAccess)
	assert.True(t, model.(*restapi.UserRole).CanConfigureReleases)
	assert.True(t, model.(*restapi.UserRole).CanConfigureLogManagement)
	assert.True(t, model.(*restapi.UserRole).CanCreatePublicCustomDashboards)
	assert.True(t, model.(*restapi.UserRole).CanViewLogs)
	assert.True(t, model.(*restapi.UserRole).CanViewTraceDetails)
}
