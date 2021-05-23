package instana_test

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana"

	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

const (
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
