package instana_test

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/mocks"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
)

var testUserRoleProviders = map[string]terraform.ResourceProvider{
	"instana": Provider(),
}

const resourceUserRoleDefinitionTemplate = `
provider "instana" {
  api_token = "test-token"
  endpoint = "localhost:{{PORT}}"
}

resource "instana_user_role" "example" {
  name = "name"
  implicit_view_filter = "view filter"
  can_configure_service_mapping = true
  can_configure_eum_applications = true
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
}
`

const userRoleApiPath = restapi.UserRolesResourcePath + "/{id}"
const testUserRoleDefinition = "instana_user_role.example"

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
			"implicitViewFilter" : "view filter",
			"canConfigureServiceMapping" : true,
			"canConfigureEumApplications" : true,
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
			"canConfigureApplications" : true
		}
		`, "{{id}}", vars["id"])
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()

	resourceUserRoleDefinition := strings.ReplaceAll(resourceUserRoleDefinitionTemplate, "{{PORT}}", strconv.Itoa(httpServer.GetPort()))

	resource.UnitTest(t, resource.TestCase{
		Providers: testUserRoleProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: resourceUserRoleDefinition,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(testUserRoleDefinition, "id"),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldName, "name"),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldImplicitViewFilter, "view filter"),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanConfigureServiceMapping, "true"),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanConfigureEumApplications, "true"),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanConfigureUsers, "true"),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanInstallNewAgents, "true"),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanSeeUsageInformation, "true"),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanConfigureIntegrations, "true"),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanSeeOnPremiseLicenseInformation, "true"),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanConfigureRoles, "true"),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanConfigureCustomAlerts, "true"),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanConfigureAPITokens, "true"),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanConfigureAgentRunMode, "true"),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanViewAuditLog, "true"),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanConfigureObjectives, "true"),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanConfigureAgents, "true"),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanConfigureAuthenticationMethods, "true"),
					resource.TestCheckResourceAttr(testUserRoleDefinition, UserRoleFieldCanConfigureApplications, "true"),
				),
			},
		},
	})
}

func TestResourceUserRoleDefinition(t *testing.T) {
	resource := CreateResourceUserRole()

	validateUserRoleResourceSchema(resource.Schema, t)

	if resource.Create == nil {
		t.Fatal("Create function expected")
	}
	if resource.Update == nil {
		t.Fatal("Update function expected")
	}
	if resource.Read == nil {
		t.Fatal("Read function expected")
	}
	if resource.Delete == nil {
		t.Fatal("Delete function expected")
	}
}

func validateUserRoleResourceSchema(schemaMap map[string]*schema.Schema, t *testing.T) {
	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(UserRoleFieldName)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(UserRoleFieldImplicitViewFilter)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(UserRoleFieldCanConfigureServiceMapping, false)
	schemaAssert.AssertSchemaIsOfTypeBooleanWithDefault(UserRoleFieldCanConfigureEumApplications, false)
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

}

func TestShouldSuccessfullyReadUserRoleFromInstanaAPIWhenMinimalDataIsReturned(t *testing.T) {
	expectedModel := createMinimalTestUserRoleModel()
	testShouldSuccessfullyReadUserRoleFromInstanaAPI(expectedModel, t)
}

func TestShouldSuccessfullyReadUserRoleFromInstanaAPIWhenFullDataIsReturned(t *testing.T) {
	expectedModel := createFullTestUserRoleModel()
	testShouldSuccessfullyReadUserRoleFromInstanaAPI(expectedModel, t)
}

func testShouldSuccessfullyReadUserRoleFromInstanaAPI(expectedModel restapi.UserRole, t *testing.T) {
	resourceData := NewTestHelper(t).CreateEmptyUserRoleResourceData()
	userRoleID := "user-role-id"
	resourceData.SetId(userRoleID)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRoleApi := mocks.NewMockUserRoleResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().UserRoles().Return(mockUserRoleApi).Times(1)
	mockUserRoleApi.EXPECT().GetOne(gomock.Eq(userRoleID)).Return(expectedModel, nil).Times(1)

	err := ReadUserRole(resourceData, mockInstanaAPI)

	if err != nil {
		t.Fatalf("Expected no error to be returned, %s", err)
	}
	verifyUserRoleModelAppliedToResource(expectedModel, resourceData, t)
}

func TestShouldFailToReadUserRoleFromInstanaAPIWhenIDIsMissing(t *testing.T) {
	resourceData := NewTestHelper(t).CreateEmptyUserRoleResourceData()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	err := ReadUserRole(resourceData, mockInstanaAPI)

	if err == nil || !strings.HasPrefix(err.Error(), "ID of user role") {
		t.Fatal("Expected error to occur because of missing id")
	}
}

func TestShouldFailToReadUserRoleFromInstanaAPIAndDeleteResourceWhenUserRoleDoesNotExist(t *testing.T) {
	resourceData := NewTestHelper(t).CreateEmptyUserRoleResourceData()
	userRoleID := "user-role-id"
	resourceData.SetId(userRoleID)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRoleApi := mocks.NewMockUserRoleResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().UserRoles().Return(mockUserRoleApi).Times(1)
	mockUserRoleApi.EXPECT().GetOne(gomock.Eq(userRoleID)).Return(restapi.UserRole{}, restapi.ErrEntityNotFound).Times(1)

	err := ReadUserRole(resourceData, mockInstanaAPI)

	if err != nil {
		t.Fatalf("Expected no error to be returned, %s", err)
	}
	if len(resourceData.Id()) > 0 {
		t.Fatal("Expected ID to be cleaned to destroy resource")
	}
}

func TestShouldFailToReadUserRoleFromInstanaAPIAndReturnErrorWhenAPICallFails(t *testing.T) {
	resourceData := NewTestHelper(t).CreateEmptyUserRoleResourceData()
	userRoleID := "user-role-id"
	resourceData.SetId(userRoleID)
	expectedError := errors.New("test")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRoleApi := mocks.NewMockUserRoleResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().UserRoles().Return(mockUserRoleApi).Times(1)
	mockUserRoleApi.EXPECT().GetOne(gomock.Eq(userRoleID)).Return(restapi.UserRole{}, expectedError).Times(1)

	err := ReadUserRole(resourceData, mockInstanaAPI)

	if err == nil || err != expectedError {
		t.Fatal("Expected error should be returned")
	}
	if len(resourceData.Id()) == 0 {
		t.Fatal("Expected ID should still be set")
	}
}

func TestShouldCreateUserRoleThroughInstanaAPI(t *testing.T) {
	data := createFullTestUserRoleData()
	resourceData := NewTestHelper(t).CreateUserRoleResourceData(data)
	expectedModel := createFullTestUserRoleModel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRoleApi := mocks.NewMockUserRoleResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().UserRoles().Return(mockUserRoleApi).Times(1)
	mockUserRoleApi.EXPECT().Upsert(gomock.AssignableToTypeOf(restapi.UserRole{})).Return(expectedModel, nil).Times(1)

	err := CreateUserRole(resourceData, mockInstanaAPI)

	if err != nil {
		t.Fatalf("Expected no error to be returned, %s", err)
	}
	verifyUserRoleModelAppliedToResource(expectedModel, resourceData, t)
}

func TestShouldReturnErrorWhenCreateUserRoleFailsThroughInstanaAPI(t *testing.T) {
	data := createFullTestUserRoleData()
	resourceData := NewTestHelper(t).CreateUserRoleResourceData(data)
	expectedError := errors.New("test")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRoleApi := mocks.NewMockUserRoleResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().UserRoles().Return(mockUserRoleApi).Times(1)
	mockUserRoleApi.EXPECT().Upsert(gomock.AssignableToTypeOf(restapi.UserRole{})).Return(restapi.UserRole{}, expectedError).Times(1)

	err := CreateUserRole(resourceData, mockInstanaAPI)

	if err == nil || expectedError != err {
		t.Fatal("Expected definned error to be returned")
	}
}

func TestShouldDeleteUserRoleThroughInstanaAPI(t *testing.T) {
	id := "test-id"
	data := createFullTestUserRoleData()
	resourceData := NewTestHelper(t).CreateUserRoleResourceData(data)
	resourceData.SetId(id)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRoleApi := mocks.NewMockUserRoleResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().UserRoles().Return(mockUserRoleApi).Times(1)
	mockUserRoleApi.EXPECT().DeleteByID(gomock.Eq(id)).Return(nil).Times(1)

	err := DeleteUserRole(resourceData, mockInstanaAPI)

	if err != nil {
		t.Fatalf("Expected no error to be returned, %s", err)
	}
	if len(resourceData.Id()) > 0 {
		t.Fatal("Expected ID to be cleaned to destroy resource")
	}
}

func TestShouldReturnErrorWhenDeleteUserRoleFailsThroughInstanaAPI(t *testing.T) {
	id := "test-id"
	data := createFullTestUserRoleData()
	resourceData := NewTestHelper(t).CreateUserRoleResourceData(data)
	resourceData.SetId(id)
	expectedError := errors.New("test")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRoleApi := mocks.NewMockUserRoleResource(ctrl)
	mockInstanaAPI := mocks.NewMockInstanaAPI(ctrl)

	mockInstanaAPI.EXPECT().UserRoles().Return(mockUserRoleApi).Times(1)
	mockUserRoleApi.EXPECT().DeleteByID(gomock.Eq(id)).Return(expectedError).Times(1)

	err := DeleteUserRole(resourceData, mockInstanaAPI)

	if err == nil || err != expectedError {
		t.Fatal("Expected error to be returned")
	}
	if len(resourceData.Id()) == 0 {
		t.Fatal("Expected ID not to be cleaned to avoid resource is destroy")
	}
}

func verifyUserRoleModelAppliedToResource(model restapi.UserRole, resourceData *schema.ResourceData, t *testing.T) {
	if model.ID != resourceData.Id() {
		t.Fatal("Expected ID to be identical")
	}
	if model.Name != resourceData.Get(UserRoleFieldName).(string) {
		t.Fatal("Expected Name to be identical")
	}
	if model.ImplicitViewFilter != resourceData.Get(UserRoleFieldImplicitViewFilter).(string) {
		t.Fatal("Expected ImplicitViewFilter to be identical")
	}
	if model.CanConfigureServiceMapping != resourceData.Get(UserRoleFieldCanConfigureServiceMapping).(bool) {
		t.Fatal("Expected CanConfigureServiceMapping to be identical")
	}
	if model.CanConfigureEumApplications != resourceData.Get(UserRoleFieldCanConfigureEumApplications).(bool) {
		t.Fatal("Expected CanConfigureEumApplications to be identical")
	}
	if model.CanConfigureUsers != resourceData.Get(UserRoleFieldCanConfigureUsers).(bool) {
		t.Fatal("Expected CanConfigureUsers to be identical")
	}
	if model.CanInstallNewAgents != resourceData.Get(UserRoleFieldCanInstallNewAgents).(bool) {
		t.Fatal("Expected CanInstallNewAgents to be identical")
	}
	if model.CanSeeUsageInformation != resourceData.Get(UserRoleFieldCanSeeUsageInformation).(bool) {
		t.Fatal("Expected CanSeeUsageInformation to be identical")
	}
	if model.CanConfigureIntegrations != resourceData.Get(UserRoleFieldCanConfigureIntegrations).(bool) {
		t.Fatal("Expected CanConfigureIntegrations to be identical")
	}
	if model.CanSeeOnPremiseLicenseInformation != resourceData.Get(UserRoleFieldCanSeeOnPremiseLicenseInformation).(bool) {
		t.Fatal("Expected CanSeeOnPremiseLicenseInformation to be identical")
	}
	if model.CanConfigureCustomAlerts != resourceData.Get(UserRoleFieldCanConfigureCustomAlerts).(bool) {
		t.Fatal("Expected CanConfigureCustomAlerts to be identical")
	}
	if model.CanConfigureAPITokens != resourceData.Get(UserRoleFieldCanConfigureAPITokens).(bool) {
		t.Fatal("Expected CanConfigureAPITokens to be identical")
	}
	if model.CanConfigureAgentRunMode != resourceData.Get(UserRoleFieldCanConfigureAgentRunMode).(bool) {
		t.Fatal("Expected CanConfigureAgentRunMode to be identical")
	}
	if model.CanViewAuditLog != resourceData.Get(UserRoleFieldCanViewAuditLog).(bool) {
		t.Fatal("Expected CanViewAuditLog to be identical")
	}
	if model.CanConfigureObjectives != resourceData.Get(UserRoleFieldCanConfigureObjectives).(bool) {
		t.Fatal("Expected CanConfigureObjectives to be identical")
	}
	if model.CanConfigureAgents != resourceData.Get(UserRoleFieldCanConfigureAgents).(bool) {
		t.Fatal("Expected CanConfigureAgents to be identical")
	}
	if model.CanConfigureAuthenticationMethods != resourceData.Get(UserRoleFieldCanConfigureAuthenticationMethods).(bool) {
		t.Fatal("Expected CanConfigureAuthenticationMethods to be identical")
	}
	if model.CanConfigureApplications != resourceData.Get(UserRoleFieldCanConfigureApplications).(bool) {
		t.Fatal("Expected CanConfigureApplications to be identical")
	}
}

func createFullTestUserRoleModel() restapi.UserRole {
	data := createMinimalTestUserRoleModel()
	data.ImplicitViewFilter = "view filter"
	data.CanConfigureServiceMapping = true
	data.CanConfigureEumApplications = true
	data.CanConfigureUsers = true
	data.CanInstallNewAgents = true
	data.CanSeeUsageInformation = true
	data.CanConfigureIntegrations = true
	data.CanSeeOnPremiseLicenseInformation = true
	data.CanConfigureRoles = true
	data.CanConfigureCustomAlerts = true
	data.CanConfigureAPITokens = true
	data.CanConfigureAgentRunMode = true
	data.CanViewAuditLog = true
	data.CanConfigureObjectives = true
	data.CanConfigureAgents = true
	data.CanConfigureAuthenticationMethods = true
	data.CanConfigureApplications = true
	return data
}

func createMinimalTestUserRoleModel() restapi.UserRole {
	return restapi.UserRole{
		ID:   "id",
		Name: "name",
	}
}

func createFullTestUserRoleData() map[string]interface{} {
	data := make(map[string]interface{})
	data[UserRoleFieldName] = "name"
	data[UserRoleFieldImplicitViewFilter] = "view filter"
	data[UserRoleFieldCanConfigureServiceMapping] = true
	data[UserRoleFieldCanConfigureEumApplications] = true
	data[UserRoleFieldCanConfigureUsers] = true
	data[UserRoleFieldCanInstallNewAgents] = true
	data[UserRoleFieldCanSeeUsageInformation] = true
	data[UserRoleFieldCanConfigureIntegrations] = true
	data[UserRoleFieldCanSeeOnPremiseLicenseInformation] = true
	data[UserRoleFieldCanConfigureRoles] = true
	data[UserRoleFieldCanConfigureCustomAlerts] = true
	data[UserRoleFieldCanConfigureAPITokens] = true
	data[UserRoleFieldCanConfigureAgentRunMode] = true
	data[UserRoleFieldCanViewAuditLog] = true
	data[UserRoleFieldCanConfigureObjectives] = true
	data[UserRoleFieldCanConfigureAgents] = true
	data[UserRoleFieldCanConfigureAuthenticationMethods] = true
	data[UserRoleFieldCanConfigureApplications] = true
	return data
}
