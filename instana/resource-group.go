package instana

import (
	"context"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
)

// ResourceInstanaGroup the name of the terraform-provider-instana resource to manage groups for role based access control
const ResourceInstanaGroup = "instana_rbac_group"

const (
	//GroupFieldName constant value for the schema field name
	GroupFieldName = "name"
	//GroupFieldFullName constant value for the schema field full_name
	//Deprecated
	GroupFieldFullName = "full_name"
	//GroupFieldMembers constant value for the schema field members
	GroupFieldMembers = "member"
	//GroupFieldMemberEmail constant value for the schema field email
	GroupFieldMemberEmail = "email"
	//GroupFieldMemberUserID constant value for the schema field user_id
	GroupFieldMemberUserID = "user_id"
	//GroupFieldPermissionSet constant value for the schema field permission_set
	GroupFieldPermissionSet = "permission_set"
	//GroupFieldPermissionSetApplicationIDs constant value for the schema field application_ids
	GroupFieldPermissionSetApplicationIDs = "application_ids"
	//GroupFieldPermissionSetInfraDFQFilter constant value for the schema field infra_dfq_filter
	GroupFieldPermissionSetInfraDFQFilter = "infra_dfq_filter"
	//GroupFieldPermissionSetKubernetesClusterUUIDs constant value for the schema field kubernetes_cluster_uuids
	GroupFieldPermissionSetKubernetesClusterUUIDs = "kubernetes_cluster_uuids"
	//GroupFieldPermissionSetKubernetesNamespaceUIDs constant value for the schema field kubernetes_namespaces_uuids
	GroupFieldPermissionSetKubernetesNamespaceUIDs = "kubernetes_namespaces_uuids"
	//GroupFieldPermissionSetMobileAppIDs constant value for the schema field mobile_app_ids
	GroupFieldPermissionSetMobileAppIDs = "mobile_app_ids"
	//GroupFieldPermissionSetWebsiteIDs constant value for the schema field website_ids
	GroupFieldPermissionSetWebsiteIDs = "website_ids"
	//GroupFieldPermissionSetPermissions constant value for the schema field permissions
	GroupFieldPermissionSetPermissions = "permissions"

	groupMaxNumberOfSetElements = 1024
)

var groupPermissionSet = map[string]*schema.Schema{
	GroupFieldPermissionSetApplicationIDs: {
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "The scope bindings to restrict access to applications",
		MaxItems:    groupMaxNumberOfSetElements,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	GroupFieldPermissionSetInfraDFQFilter: {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The scope binding for the dynamic filter query to restrict access to infrastructure assets",
	},
	GroupFieldPermissionSetKubernetesClusterUUIDs: {
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "The scope bindings to restrict access to Kubernetes Clusters",
		MaxItems:    groupMaxNumberOfSetElements,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	GroupFieldPermissionSetKubernetesNamespaceUIDs: {
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "The scope bindings to restrict access to Kubernetes namespaces",
		MaxItems:    groupMaxNumberOfSetElements,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	GroupFieldPermissionSetMobileAppIDs: {
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "The scope bindings to restrict access to mobile apps",
		MaxItems:    groupMaxNumberOfSetElements,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	GroupFieldPermissionSetWebsiteIDs: {
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "The scope bindings to restrict access to websites",
		MaxItems:    groupMaxNumberOfSetElements,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	GroupFieldPermissionSetPermissions: {
		Type:        schema.TypeSet,
		Optional:    true,
		Description: "The permissions assigned which should be assigned to the users of the group",
		MaxItems:    groupMaxNumberOfSetElements,
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validation.StringInSlice(restapi.SupportedInstanaPermissions.ToStringSlice(), false),
		},
	},
}

var groupMemberSchema = map[string]*schema.Schema{
	GroupFieldMemberUserID: {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The user id of the group member",
	},
	GroupFieldMemberEmail: {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The email address of the group member",
	},
}

var groupSchemaName = &schema.Schema{
	Type:        schema.TypeString,
	Required:    true,
	Description: "The name of the Group",
}

// Deprecated
var groupSchemaFullName = &schema.Schema{
	Type:        schema.TypeString,
	Computed:    true,
	Description: "The full name of the group. The field is computed and contains the name which is sent to instana. The computation depends on the configured default_name_prefix and default_name_suffix at provider level",
}
var groupSchemaMembers = &schema.Schema{
	Type:        schema.TypeSet,
	Optional:    true,
	Description: "The members of the group",
	MaxItems:    1024,
	Elem: &schema.Resource{
		Schema: groupMemberSchema,
	},
}
var groupSchemaPermissionSet = &schema.Schema{
	Type:        schema.TypeList,
	Optional:    true,
	MaxItems:    1,
	Description: "The permission set of the group",
	Elem: &schema.Resource{
		Schema: groupPermissionSet,
	},
}

var groupSchema = map[string]*schema.Schema{
	GroupFieldName:          groupSchemaName,
	GroupFieldMembers:       groupSchemaMembers,
	GroupFieldPermissionSet: groupSchemaPermissionSet,
}

// NewGroupResourceHandle creates the resource handle for RBAC Groups
func NewGroupResourceHandle() ResourceHandle[*restapi.Group] {
	return &groupResource{
		metaData: ResourceMetaData{
			ResourceName:     ResourceInstanaGroup,
			Schema:           groupSchema,
			SchemaVersion:    1,
			SkipIDGeneration: true,
		},
	}
}

type groupResource struct {
	metaData ResourceMetaData
}

func (r *groupResource) MetaData() *ResourceMetaData {
	return &r.metaData
}

func (r *groupResource) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{
		{
			Type:    r.groupSchemaV0().CoreConfigSchema().ImpliedType(),
			Upgrade: r.groupStateUpgradeV0,
			Version: 0,
		},
	}
}

func (r *groupResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.Group] {
	return api.Groups()
}

func (r *groupResource) SetComputedFields(_ *schema.ResourceData) error {
	return nil
}

func (r *groupResource) UpdateState(d *schema.ResourceData, group *restapi.Group) error {
	data := map[string]interface{}{
		GroupFieldName: group.Name,
	}

	members := r.convertGroupMembersToState(group)
	if members != nil {
		data[GroupFieldMembers] = members
	}
	if !group.PermissionSet.IsEmpty() {
		permissions := r.convertPermissionSetToState(group)
		data[GroupFieldPermissionSet] = permissions
	}

	d.SetId(group.ID)
	return tfutils.UpdateState(d, data)
}

func (r *groupResource) convertGroupMembersToState(obj *restapi.Group) *schema.Set {
	result := make([]interface{}, len(obj.Members))
	for i, v := range obj.Members {
		var email string
		if v.Email != nil {
			email = *v.Email
		}
		groupMap := map[string]interface{}{
			GroupFieldMemberUserID: v.UserID,
			GroupFieldMemberEmail:  email,
		}
		result[i] = groupMap
	}
	if len(result) > 0 {
		return schema.NewSet(schema.HashResource(groupSchema[GroupFieldMembers].Elem.(*schema.Resource)), result)
	}
	return nil
}

func (r *groupResource) convertPermissionSetToState(obj *restapi.Group) []interface{} {
	permissionSet := obj.PermissionSet

	m := make(map[string]interface{})
	if obj.PermissionSet.InfraDFQFilter != nil && len(obj.PermissionSet.InfraDFQFilter.ScopeID) > 0 {
		m[GroupFieldPermissionSetInfraDFQFilter] = permissionSet.InfraDFQFilter.ScopeID
	}
	m[GroupFieldPermissionSetApplicationIDs] = r.convertScopeBindingSliceToState(permissionSet.ApplicationIDs)
	m[GroupFieldPermissionSetKubernetesClusterUUIDs] = r.convertScopeBindingSliceToState(permissionSet.KubernetesClusterUUIDs)
	m[GroupFieldPermissionSetKubernetesNamespaceUIDs] = r.convertScopeBindingSliceToState(permissionSet.KubernetesNamespaceUIDs)
	m[GroupFieldPermissionSetMobileAppIDs] = r.convertScopeBindingSliceToState(permissionSet.MobileAppIDs)
	m[GroupFieldPermissionSetWebsiteIDs] = r.convertScopeBindingSliceToState(permissionSet.WebsiteIDs)
	m[GroupFieldPermissionSetPermissions] = permissionSet.Permissions
	return []interface{}{m}
}

func (r *groupResource) convertScopeBindingSliceToState(value []restapi.ScopeBinding) *schema.Set {
	result := make([]interface{}, len(value))
	for i, v := range value {
		result[i] = v.ScopeID
	}
	return schema.NewSet(schema.HashString, result)
}

func (r *groupResource) MapStateToDataObject(d *schema.ResourceData) (*restapi.Group, error) {
	members := r.convertStateToGroupMembers(d)
	permissionSet := r.convertStateToPermissionSet(d)
	return &restapi.Group{
		ID:            d.Id(),
		Name:          d.Get(GroupFieldName).(string),
		Members:       members,
		PermissionSet: *permissionSet,
	}, nil
}

func (r *groupResource) convertStateToGroupMembers(d *schema.ResourceData) []restapi.APIMember {
	if val, ok := d.GetOk(GroupFieldMembers); ok {
		if set, ok := val.(*schema.Set); ok {
			slice := set.List()
			result := make([]restapi.APIMember, len(slice))
			for i, setValue := range slice {
				if valueMap, ok := setValue.(map[string]interface{}); ok {
					member := restapi.APIMember{UserID: valueMap[GroupFieldMemberUserID].(string)}
					if email, ok := valueMap[GroupFieldMemberEmail].(string); ok {
						member.Email = &email
					}
					result[i] = member
				} else {
					log.Printf("WARN: group member cannot be read; %v\n", setValue)
				}
			}
			return result
		}
		log.Println("WARN: group member state cannot be read")
	}
	return []restapi.APIMember{}
}

func (r *groupResource) convertStateToPermissionSet(d *schema.ResourceData) *restapi.APIPermissionSetWithRoles {
	if val, ok := d.GetOk(GroupFieldPermissionSet); ok {
		if permissionSetSlice, ok := val.([]interface{}); ok && len(permissionSetSlice) == 1 {
			if permissionSet, ok := permissionSetSlice[0].(map[string]interface{}); ok {
				return &restapi.APIPermissionSetWithRoles{
					ApplicationIDs:          r.convertStateToSliceOfScopeBinding(GroupFieldPermissionSetApplicationIDs, permissionSet[GroupFieldPermissionSetApplicationIDs]),
					InfraDFQFilter:          r.convertStateToScopeBindingPointer(GroupFieldPermissionSetInfraDFQFilter, permissionSet[GroupFieldPermissionSetInfraDFQFilter]),
					KubernetesClusterUUIDs:  r.convertStateToSliceOfScopeBinding(GroupFieldPermissionSetKubernetesClusterUUIDs, permissionSet[GroupFieldPermissionSetKubernetesClusterUUIDs]),
					KubernetesNamespaceUIDs: r.convertStateToSliceOfScopeBinding(GroupFieldPermissionSetKubernetesNamespaceUIDs, permissionSet[GroupFieldPermissionSetKubernetesNamespaceUIDs]),
					MobileAppIDs:            r.convertStateToSliceOfScopeBinding(GroupFieldPermissionSetMobileAppIDs, permissionSet[GroupFieldPermissionSetMobileAppIDs]),
					WebsiteIDs:              r.convertStateToSliceOfScopeBinding(GroupFieldPermissionSetWebsiteIDs, permissionSet[GroupFieldPermissionSetWebsiteIDs]),
					Permissions:             r.convertStateToPermissions(GroupFieldPermissionSetPermissions, permissionSet[GroupFieldPermissionSetPermissions]),
				}
			}
		}
		log.Println("WARN: permission_set state cannot be read")
	}
	emptyScopeBinding := make([]restapi.ScopeBinding, 0)
	return &restapi.APIPermissionSetWithRoles{
		ApplicationIDs:          emptyScopeBinding,
		KubernetesNamespaceUIDs: emptyScopeBinding,
		KubernetesClusterUUIDs:  emptyScopeBinding,
		WebsiteIDs:              emptyScopeBinding,
		MobileAppIDs:            emptyScopeBinding,
		Permissions:             make([]restapi.InstanaPermission, 0),
	}
}

func (r *groupResource) convertStateToSliceOfScopeBinding(attribute string, val interface{}) []restapi.ScopeBinding {
	if set, ok := val.(*schema.Set); ok {
		slice := set.List()
		result := make([]restapi.ScopeBinding, len(slice))
		for i, v := range slice {
			result[i] = restapi.ScopeBinding{ScopeID: v.(string)}
		}
		return result
	}
	log.Printf("WARN: %s state cannot be read\n", attribute)
	return make([]restapi.ScopeBinding, 0)
}

func (r *groupResource) convertStateToScopeBindingPointer(attribute string, val interface{}) *restapi.ScopeBinding {
	if v, ok := val.(string); ok {
		return &restapi.ScopeBinding{ScopeID: v}
	}
	log.Printf("WARN: %s state cannot be read\n", attribute)
	return nil
}

func (r *groupResource) convertStateToPermissions(attribute string, val interface{}) []restapi.InstanaPermission {
	if set, ok := val.(*schema.Set); ok {
		slice := set.List()
		result := make([]restapi.InstanaPermission, len(slice))
		for i, v := range slice {
			result[i] = restapi.InstanaPermission(v.(string))
		}
		return result
	}
	log.Printf("WARN: %s state cannot be read\n", attribute)
	return make([]restapi.InstanaPermission, 0)
}

func (r *groupResource) groupStateUpgradeV0(_ context.Context, state map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	if _, ok := state[GroupFieldFullName]; ok {
		state[GroupFieldName] = state[GroupFieldFullName]
		delete(state, GroupFieldFullName)
	}
	return state, nil
}

func (r *groupResource) groupSchemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			GroupFieldName:          groupSchemaName,
			GroupFieldFullName:      groupSchemaFullName,
			GroupFieldMembers:       groupSchemaMembers,
			GroupFieldPermissionSet: groupSchemaPermissionSet,
		},
	}
}
