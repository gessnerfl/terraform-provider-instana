package instana_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"net/http"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana"

	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

const fullRBACGroupDefinitionTemplate = `
resource "instana_rbac_group" "example" {
  name = "name %d"
  member { 
	user_id = "user_1_id"
	email = "user_1_email"
  }
  member { 
	user_id = "user_2_id"
	email = "user_2_email"
  }

  permission_set {
	application_ids = [ "app_id1", "app_id2" ]
	kubernetes_cluster_uuids = [ "k8s_cluster_id1", "k8s_cluster_id2" ]
	kubernetes_namespaces_uuids = [ "k8s_namespace_id1", "k8s_namespace_id2" ]
	mobile_app_ids = [ "mobile_app_id1", "mobile_app_id2" ]
	website_ids = [ "website_id1", "website_id2" ]
    infra_dfq_filter = "infra_dfq"
	permissions = [ "CAN_CONFIGURE_APPLICATIONS", "CAN_CONFIGURE_AGENTS" ]
  }
}
`

const minimalRBACGroupDefinitionTemplate = `
resource "instana_rbac_group" "example" {
  name = "name %d"
}
`

const rbacGroupDefinitionWithPermissionsAndApplicationIDsAssignedTemplate = `
resource "instana_rbac_group" "example" {
  name = "name %d"

  permission_set {
	application_ids = [ "app_id1", "app_id2" ]
	permissions = [ "CAN_CONFIGURE_APPLICATIONS", "CAN_CONFIGURE_AGENTS" ]
  }
}
`

const (
	groupApiPath              = restapi.GroupsResourcePath + "/{internal-id}"
	testGroupDefinition       = "instana_rbac_group.example"
	defaultGroupID            = "group_id"
	defaultGroupName          = "group_name"
	defaultGroupFullName      = "prefix group_name suffix"
	defaultGroupMember1UserID = "user_1_id"
	defaultGroupMember1Email  = "user_1_email"
	defaultGroupMember2UserID = "user_2_id"
	defaultGroupMember2Email  = "user_2_email"
	defaultScope1ID           = "scope_1_id"
	defaultScope2ID           = "scope_2_id"
)

func TestCRUDOfMinimalRBACGroupResourceWithMockServer(t *testing.T) {
	serverResponseTemplate := `
		{
			"id" : "%s",
			"name" : "prefix name %d suffix",
			"members": [],
		    "permissionSet": {
				"permissions": [],
				"applicationIds": [],
				"kubernetesClusterUUIDs": [],
				"kubernetesNamespaceUIDs": [],
				"websiteIds": [],
				"mobileAppIds": [],
				"infraDfqFilter": {
				  "scopeId": "",
				  "scopeRoleId": "-600"
				}
			}
		}
		`

	testStepsFactory := func(httpServerPort int, resourceID string) []resource.TestStep {
		return []resource.TestStep{
			createMinimalRbacGroupResourceTestStep(httpServerPort, 0, resourceID),
			testStepImportWithCustomID(testGroupDefinition, resourceID),
			createMinimalRbacGroupResourceTestStep(httpServerPort, 1, resourceID),
			testStepImportWithCustomID(testGroupDefinition, resourceID),
		}
	}

	executeRBACGroupIntegrationTest(t, serverResponseTemplate, testStepsFactory)
}

func createMinimalRbacGroupResourceTestStep(httpPort int, iteration int, id string) resource.TestStep {
	return resource.TestStep{
		Config: appendProviderConfig(fmt.Sprintf(minimalRBACGroupDefinitionTemplate, iteration), httpPort),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(testGroupDefinition, "id", id),
			resource.TestCheckResourceAttr(testGroupDefinition, GroupFieldName, formatResourceName(iteration)),
			resource.TestCheckResourceAttr(testGroupDefinition, GroupFieldFullName, formatResourceFullName(iteration)),
		),
	}
}

func TestCRUDOfRBACGroupResourceWithPermissionsAndApplicationIdsAssignedUsingMockServer(t *testing.T) {
	serverResponseTemplate := `
		{
			"id" : "%s",
			"name" : "prefix name %d suffix",
			"members": [],
		    "permissionSet": {
				"permissions": [ "CAN_CONFIGURE_APPLICATIONS", "CAN_CONFIGURE_AGENTS" ],
				"applicationIds": [ { "scopeId" : "app_id1" },  { "scopeId" : "app_id2" } ]
			}
		}
		`

	testStepsFactory := func(httpServerPort int, resourceID string) []resource.TestStep {
		return []resource.TestStep{
			createRbacGroupResourceWithPermissionsAndApplicationIdsAssignedTestStep(httpServerPort, 0, resourceID),
			testStepImportWithCustomID(testGroupDefinition, resourceID),
			createRbacGroupResourceWithPermissionsAndApplicationIdsAssignedTestStep(httpServerPort, 1, resourceID),
			testStepImportWithCustomID(testGroupDefinition, resourceID),
		}
	}

	executeRBACGroupIntegrationTest(t, serverResponseTemplate, testStepsFactory)
}

func createRbacGroupResourceWithPermissionsAndApplicationIdsAssignedTestStep(httpPort int, iteration int, id string) resource.TestStep {
	return resource.TestStep{
		Config: appendProviderConfig(fmt.Sprintf(rbacGroupDefinitionWithPermissionsAndApplicationIDsAssignedTemplate, iteration), httpPort),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(testGroupDefinition, "id", id),
			resource.TestCheckResourceAttr(testGroupDefinition, GroupFieldName, formatResourceName(iteration)),
			resource.TestCheckResourceAttr(testGroupDefinition, GroupFieldFullName, formatResourceFullName(iteration)),
			testCheckPermissionSetStringSetField(GroupFieldPermissionSetApplicationIDs, 0, "app_id1"),
			testCheckPermissionSetStringSetField(GroupFieldPermissionSetApplicationIDs, 1, "app_id2"),
			testCheckPermissionSetStringSetField(GroupFieldPermissionSetPermissions, 0, string(restapi.PermissionCanConfigureAgents)),
			testCheckPermissionSetStringSetField(GroupFieldPermissionSetPermissions, 1, string(restapi.PermissionCanConfigureApplications)),
		),
	}
}

func TestCRUDOfFullRBACGroupResourceWithMockServer(t *testing.T) {
	serverResponseTemplate := `
		{
			"id" : "%s",
			"name" : "prefix name %d suffix",
			"members" : [ 
				{ "userId" : "user_1_id", "email" : "user_1_email" }, 
				{ "userId" : "user_2_id", "email" : "user_2_email" } 
			],
			"permissionSet" : {
				"applicationIds" : [ { "scopeId" : "app_id1" },  { "scopeId" : "app_id2" } ],
				"kubernetesClusterUUIDs" : [ { "scopeId" : "k8s_cluster_id1" },  { "scopeId" : "k8s_cluster_id2" } ],
				"kubernetesNamespaceUIDs" : [ { "scopeId" : "k8s_namespace_id1" },  { "scopeId" : "k8s_namespace_id2" } ],
				"mobileAppIds" : [ { "scopeId" : "mobile_app_id1" },  { "scopeId" : "mobile_app_id2" } ],
				"websiteIds" : [ { "scopeId" : "website_id1" },  { "scopeId" : "website_id2" } ],
   			 	"infraDfqFilter" : { "scopeId" : "infra_dfq" },
				"permissions" : [ "CAN_CONFIGURE_APPLICATIONS", "CAN_CONFIGURE_AGENTS" ]
			}
		}
		`

	testStepsFactory := func(httpServerPort int, resourceID string) []resource.TestStep {
		return []resource.TestStep{
			createFullRbacGroupResourceTestStep(httpServerPort, 0, resourceID),
			testStepImportWithCustomID(testGroupDefinition, resourceID),
			createFullRbacGroupResourceTestStep(httpServerPort, 1, resourceID),
			testStepImportWithCustomID(testGroupDefinition, resourceID),
		}
	}

	executeRBACGroupIntegrationTest(t, serverResponseTemplate, testStepsFactory)
}

func createFullRbacGroupResourceTestStep(httpPort int, iteration int, id string) resource.TestStep {
	return resource.TestStep{
		Config: appendProviderConfig(fmt.Sprintf(fullRBACGroupDefinitionTemplate, iteration), httpPort),
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(testGroupDefinition, "id", id),
			resource.TestCheckResourceAttr(testGroupDefinition, GroupFieldName, formatResourceName(iteration)),
			resource.TestCheckResourceAttr(testGroupDefinition, GroupFieldFullName, formatResourceFullName(iteration)),
			resource.TestCheckResourceAttr(testGroupDefinition, fmt.Sprintf("%s.%d.%s", GroupFieldMembers, 0, GroupFieldMemberUserID), defaultGroupMember1UserID),
			resource.TestCheckResourceAttr(testGroupDefinition, fmt.Sprintf("%s.%d.%s", GroupFieldMembers, 0, GroupFieldMemberEmail), defaultGroupMember1Email),
			resource.TestCheckResourceAttr(testGroupDefinition, fmt.Sprintf("%s.%d.%s", GroupFieldMembers, 1, GroupFieldMemberUserID), defaultGroupMember2UserID),
			resource.TestCheckResourceAttr(testGroupDefinition, fmt.Sprintf("%s.%d.%s", GroupFieldMembers, 1, GroupFieldMemberEmail), defaultGroupMember2Email),
			testCheckPermissionSetStringSetField(GroupFieldPermissionSetApplicationIDs, 0, "app_id1"),
			testCheckPermissionSetStringSetField(GroupFieldPermissionSetApplicationIDs, 1, "app_id2"),
			testCheckPermissionSetStringSetField(GroupFieldPermissionSetKubernetesClusterUUIDs, 0, "k8s_cluster_id1"),
			testCheckPermissionSetStringSetField(GroupFieldPermissionSetKubernetesClusterUUIDs, 1, "k8s_cluster_id2"),
			testCheckPermissionSetStringSetField(GroupFieldPermissionSetKubernetesNamespaceUIDs, 0, "k8s_namespace_id1"),
			testCheckPermissionSetStringSetField(GroupFieldPermissionSetKubernetesNamespaceUIDs, 1, "k8s_namespace_id2"),
			testCheckPermissionSetStringSetField(GroupFieldPermissionSetMobileAppIDs, 0, "mobile_app_id1"),
			testCheckPermissionSetStringSetField(GroupFieldPermissionSetMobileAppIDs, 1, "mobile_app_id2"),
			testCheckPermissionSetStringSetField(GroupFieldPermissionSetWebsiteIDs, 0, "website_id1"),
			testCheckPermissionSetStringSetField(GroupFieldPermissionSetWebsiteIDs, 1, "website_id2"),
			testCheckPermissionSetStringSetField(GroupFieldPermissionSetPermissions, 0, string(restapi.PermissionCanConfigureAgents)),
			testCheckPermissionSetStringSetField(GroupFieldPermissionSetPermissions, 1, string(restapi.PermissionCanConfigureApplications)),
			resource.TestCheckResourceAttr(testGroupDefinition, fmt.Sprintf("%s.0.%s", GroupFieldPermissionSet, GroupFieldPermissionSetInfraDFQFilter), "infra_dfq"),
		),
	}
}

func testCheckPermissionSetStringSetField(attribute string, idx int, value string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttr(testGroupDefinition, fmt.Sprintf("%s.0.%s.%d", GroupFieldPermissionSet, attribute, idx), value)
}

type rbacGroupTestStepsFactory func(httpServerPort int, resourceID string) []resource.TestStep

func executeRBACGroupIntegrationTest(t *testing.T, serverResponseTemplate string, testStepsFactory rbacGroupTestStepsFactory) {
	id := RandomID()
	testutils.DeactivateTLSServerCertificateVerification()
	httpServer := testutils.NewTestHTTPServer()
	httpServer.AddRoute(http.MethodPost, restapi.GroupsResourcePath, func(w http.ResponseWriter, r *http.Request) {
		group := &restapi.Group{}
		err := json.NewDecoder(r.Body).Decode(group)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			r.Write(bytes.NewBufferString("Failed to get request"))
		} else {
			group.ID = id
			w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(group)
		}
	})
	httpServer.AddRoute(http.MethodPut, groupApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodDelete, groupApiPath, testutils.EchoHandlerFunc)
	httpServer.AddRoute(http.MethodGet, groupApiPath, func(w http.ResponseWriter, r *http.Request) {
		modCount := httpServer.GetCallCount(http.MethodPut, restapi.GroupsResourcePath+"/"+id)
		json := fmt.Sprintf(serverResponseTemplate, id, modCount)
		w.Header().Set(contentType, r.Header.Get(contentType))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(json))
	})
	httpServer.Start()
	defer httpServer.Close()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testProviderFactory,
		Steps:             testStepsFactory(httpServer.GetPort(), id),
	})
}

func TestResourceGroupDefinition(t *testing.T) {
	resource := NewGroupResourceHandle()

	schemaMap := resource.MetaData().Schema

	schemaAssert := testutils.NewTerraformSchemaAssert(schemaMap, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(GroupFieldName)
	schemaAssert.AssertSchemaIsComputedAndOfTypeString(GroupFieldFullName)
	verifyGroupMemberSchema(t, schemaMap[GroupFieldMembers])
	verifyGroupPermissionSetSchema(t, schemaMap[GroupFieldPermissionSet])
}

func verifyGroupMemberSchema(t *testing.T, groupMemberSchema *schema.Schema) {
	require.False(t, groupMemberSchema.Required)
	require.True(t, groupMemberSchema.Optional)
	require.Equal(t, 1024, groupMemberSchema.MaxItems)
	require.Equal(t, 0, groupMemberSchema.MinItems)
	require.Equal(t, schema.TypeSet, groupMemberSchema.Type)
	require.IsType(t, &schema.Resource{}, groupMemberSchema.Elem)

	elemSchema := groupMemberSchema.Elem.(*schema.Resource).Schema
	schemaAssert := testutils.NewTerraformSchemaAssert(elemSchema, t)
	schemaAssert.AssertSchemaIsRequiredAndOfTypeString(GroupFieldMemberUserID)
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(GroupFieldMemberEmail)
}

func verifyGroupPermissionSetSchema(t *testing.T, permissionSetSchema *schema.Schema) {
	require.False(t, permissionSetSchema.Required)
	require.True(t, permissionSetSchema.Optional)
	require.Equal(t, 1, permissionSetSchema.MaxItems)
	require.Equal(t, 0, permissionSetSchema.MinItems)
	require.Equal(t, schema.TypeList, permissionSetSchema.Type)
	require.IsType(t, &schema.Resource{}, permissionSetSchema.Elem)

	elemSchema := permissionSetSchema.Elem.(*schema.Resource).Schema
	schemaAssert := testutils.NewTerraformSchemaAssert(elemSchema, t)
	verifyOptionalSetOfStringWithUpTo1024Elements(t, elemSchema[GroupFieldPermissionSetApplicationIDs])
	schemaAssert.AssertSchemaIsOptionalAndOfTypeString(GroupFieldPermissionSetInfraDFQFilter)
	verifyOptionalSetOfStringWithUpTo1024Elements(t, elemSchema[GroupFieldPermissionSetKubernetesClusterUUIDs])
	verifyOptionalSetOfStringWithUpTo1024Elements(t, elemSchema[GroupFieldPermissionSetKubernetesClusterUUIDs])
	verifyOptionalSetOfStringWithUpTo1024Elements(t, elemSchema[GroupFieldPermissionSetMobileAppIDs])
	verifyOptionalSetOfStringWithUpTo1024Elements(t, elemSchema[GroupFieldPermissionSetWebsiteIDs])
	verifyOptionalSetOfStringWithUpTo1024Elements(t, elemSchema[GroupFieldPermissionSetPermissions])
}

func verifyOptionalSetOfStringWithUpTo1024Elements(t *testing.T, s *schema.Schema) {
	require.False(t, s.Required)
	require.True(t, s.Optional)
	require.Equal(t, 1024, s.MaxItems)
	require.Equal(t, 0, s.MinItems)
	require.Equal(t, schema.TypeSet, s.Type)
	require.IsType(t, &schema.Schema{}, s.Elem)
	require.IsType(t, schema.TypeString, s.Elem.(*schema.Schema).Type)
}

func TestShouldReturnCorrectResourceNameForGroups(t *testing.T) {
	name := NewGroupResourceHandle().MetaData().ResourceName

	require.Equal(t, "instana_rbac_group", name)
}

func TestGroupResourceShouldHaveSchemaVersionZero(t *testing.T) {
	require.Equal(t, 0, NewGroupResourceHandle().MetaData().SchemaVersion)
}

func TestUpdateStateOfGroupResource(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewGroupResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	member1Email := defaultGroupMember1Email
	member2Email := defaultGroupMember2Email
	scopeBinding1 := restapi.ScopeBinding{ScopeID: defaultScope1ID}
	scopeBinding2 := restapi.ScopeBinding{ScopeID: defaultScope2ID}
	group := restapi.Group{
		ID:   defaultGroupID,
		Name: defaultGroupName,
		Members: []restapi.APIMember{
			{UserID: defaultGroupMember1UserID, Email: &member1Email},
			{UserID: defaultGroupMember2UserID, Email: &member2Email},
		},
		PermissionSet: restapi.APIPermissionSetWithRoles{
			ApplicationIDs:          []restapi.ScopeBinding{scopeBinding1, scopeBinding2},
			InfraDFQFilter:          &scopeBinding1,
			KubernetesClusterUUIDs:  []restapi.ScopeBinding{scopeBinding1, scopeBinding2},
			KubernetesNamespaceUIDs: []restapi.ScopeBinding{scopeBinding1, scopeBinding2},
			WebsiteIDs:              []restapi.ScopeBinding{scopeBinding1, scopeBinding2},
			MobileAppIDs:            []restapi.ScopeBinding{scopeBinding1, scopeBinding2},
			Permissions:             []restapi.InstanaPermission{restapi.PermissionCanConfigureAgentRunMode, restapi.PermissionCanConfigureAPITokens},
		},
	}

	err := resourceHandle.UpdateState(resourceData, &group, testHelper.ResourceFormatter())

	require.NoError(t, err)
	require.Equal(t, defaultGroupID, resourceData.Id())
	require.Equal(t, defaultGroupName, resourceData.Get(GroupFieldName))
	expectedGroupMembers := []interface{}{
		map[string]interface{}{
			GroupFieldMemberUserID: defaultGroupMember1UserID,
			GroupFieldMemberEmail:  defaultGroupMember1Email,
		},
		map[string]interface{}{
			GroupFieldMemberUserID: defaultGroupMember2UserID,
			GroupFieldMemberEmail:  defaultGroupMember2Email,
		},
	}
	require.Equal(t, expectedGroupMembers, resourceData.Get(GroupFieldMembers).(*schema.Set).List())
	require.IsType(t, []interface{}{}, resourceData.Get(GroupFieldPermissionSet))
	permissionSetSlice := resourceData.Get(GroupFieldPermissionSet).([]interface{})
	require.Len(t, permissionSetSlice, 1)
	require.IsType(t, map[string]interface{}{}, permissionSetSlice[0])
	permissionSet := permissionSetSlice[0].(map[string]interface{})
	expectedScopeSet := []interface{}{defaultScope2ID, defaultScope1ID}
	require.Equal(t, expectedScopeSet, permissionSet[GroupFieldPermissionSetApplicationIDs].(*schema.Set).List())
	require.Equal(t, expectedScopeSet, permissionSet[GroupFieldPermissionSetKubernetesNamespaceUIDs].(*schema.Set).List())
	require.Equal(t, expectedScopeSet, permissionSet[GroupFieldPermissionSetKubernetesClusterUUIDs].(*schema.Set).List())
	require.Equal(t, defaultScope1ID, permissionSet[GroupFieldPermissionSetInfraDFQFilter])
	require.Equal(t, expectedScopeSet, permissionSet[GroupFieldPermissionSetWebsiteIDs].(*schema.Set).List())
	require.Equal(t, expectedScopeSet, permissionSet[GroupFieldPermissionSetMobileAppIDs].(*schema.Set).List())
	require.Equal(t, []interface{}{string(restapi.PermissionCanConfigureAPITokens), string(restapi.PermissionCanConfigureAgentRunMode)}, permissionSet[GroupFieldPermissionSetPermissions].(*schema.Set).List())
}

func TestShouldUpdateStateWhenNoGroupMembersAndAnEmptyPermissionSetIsProvided(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewGroupResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)

	group := restapi.Group{
		ID:   defaultGroupID,
		Name: defaultGroupName,
	}

	err := resourceHandle.UpdateState(resourceData, &group, testHelper.ResourceFormatter())

	require.NoError(t, err)
	require.Equal(t, defaultGroupID, resourceData.Id())
	require.Equal(t, defaultGroupName, resourceData.Get(GroupFieldName))
	emptySlice := make([]interface{}, 0)
	require.Equal(t, emptySlice, resourceData.Get(GroupFieldMembers).(*schema.Set).List())
	permissionSetSlice := resourceData.Get(GroupFieldPermissionSet).([]interface{})
	require.Len(t, permissionSetSlice, 0)
}

func TestGroupResourceShouldReadModelFromState(t *testing.T) {
	testHelper := NewTestHelper(t)
	resourceHandle := NewGroupResourceHandle()
	resourceData := testHelper.CreateEmptyResourceDataForResourceHandle(resourceHandle)
	members := []interface{}{
		map[string]interface{}{
			GroupFieldMemberUserID: defaultGroupMember1UserID,
			GroupFieldMemberEmail:  defaultGroupMember1Email,
		},
		map[string]interface{}{
			GroupFieldMemberUserID: defaultGroupMember2UserID,
			GroupFieldMemberEmail:  defaultGroupMember2Email,
		},
	}
	scopeSet := []interface{}{defaultScope2ID, defaultScope1ID}
	permissionSet := []interface{}{
		map[string]interface{}{
			GroupFieldPermissionSetApplicationIDs:          schema.NewSet(schema.HashString, scopeSet),
			GroupFieldPermissionSetKubernetesClusterUUIDs:  schema.NewSet(schema.HashString, scopeSet),
			GroupFieldPermissionSetKubernetesNamespaceUIDs: schema.NewSet(schema.HashString, scopeSet),
			GroupFieldPermissionSetInfraDFQFilter:          defaultScope1ID,
			GroupFieldPermissionSetMobileAppIDs:            schema.NewSet(schema.HashString, scopeSet),
			GroupFieldPermissionSetWebsiteIDs:              schema.NewSet(schema.HashString, scopeSet),
			GroupFieldPermissionSetPermissions:             schema.NewSet(schema.HashString, []interface{}{string(restapi.PermissionCanConfigureAPITokens), string(restapi.PermissionCanConfigureAgentRunMode)}),
		},
	}
	resourceData.SetId(defaultGroupID)
	resourceData.Set(GroupFieldName, defaultGroupName)
	resourceData.Set(GroupFieldFullName, defaultGroupFullName)
	resourceData.Set(GroupFieldMembers, members)
	resourceData.Set(GroupFieldPermissionSet, permissionSet)

	result, err := resourceHandle.MapStateToDataObject(resourceData, testHelper.ResourceFormatter())

	require.NoError(t, err)
	require.IsType(t, &restapi.Group{}, result)
	group := result.(*restapi.Group)
	require.Equal(t, defaultGroupID, group.GetIDForResourcePath())
	require.Equal(t, defaultGroupID, group.ID)
	require.Equal(t, defaultGroupFullName, group.Name)
	member1Email := defaultGroupMember1Email
	member2Email := defaultGroupMember2Email
	expectedMembers := []restapi.APIMember{
		{UserID: defaultGroupMember1UserID, Email: &member1Email},
		{UserID: defaultGroupMember2UserID, Email: &member2Email},
	}
	require.Equal(t, expectedMembers, group.Members)
}
