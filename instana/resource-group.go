package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
)

//ResourceInstanaGroup the name of the terraform-provider-instana resource to manage groups for role based access control
const ResourceInstanaGroup = "instana_rbac_group"

const (
	//GroupFieldName constant value for the schema field name
	GroupFieldName = "name"
	//GroupFieldFullName constant value for the schema field full_name
	GroupFieldFullName = "full_name"
	//GroupFieldMembers constant value for the schema field members
	GroupFieldMembers = "members"
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
		Required:    false,
		Optional:    true,
		Description: "The scope bindings to restrict access to applications",
		MaxItems:    groupMaxNumberOfSetElements,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	GroupFieldPermissionSetInfraDFQFilter: {
		Type:        schema.TypeString,
		Required:    false,
		Optional:    true,
		Description: "The scope binding for the dynamic filter query to restrict access to infrastructure assets",
	},
	GroupFieldPermissionSetKubernetesClusterUUIDs: {
		Type:        schema.TypeSet,
		Required:    false,
		Optional:    true,
		Description: "The scope bindings to restrict access to Kubernetes Clusters",
		MaxItems:    groupMaxNumberOfSetElements,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	GroupFieldPermissionSetKubernetesNamespaceUIDs: {
		Type:        schema.TypeSet,
		Required:    false,
		Optional:    true,
		Description: "The scope bindings to restrict access to Kubernetes namespaces",
		MaxItems:    groupMaxNumberOfSetElements,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	GroupFieldPermissionSetMobileAppIDs: {
		Type:        schema.TypeSet,
		Required:    false,
		Optional:    true,
		Description: "The scope bindings to restrict access to mobile apps",
		MaxItems:    groupMaxNumberOfSetElements,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	GroupFieldPermissionSetWebsiteIDs: {
		Type:        schema.TypeSet,
		Required:    false,
		Optional:    true,
		Description: "The scope bindings to restrict access to websites",
		MaxItems:    groupMaxNumberOfSetElements,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	GroupFieldPermissionSetPermissions: {
		Type:        schema.TypeSet,
		Required:    false,
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
		Required:    false,
		Optional:    true,
		Description: "The email address of the group member",
	},
}

var groupSchema = map[string]*schema.Schema{
	GroupFieldName: {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the Group",
	},
	GroupFieldFullName: {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The full name of the group. The field is computed and contains the name which is sent to instana. The computation depends on the configured default_name_prefix and default_name_suffix at provider level",
	},
	GroupFieldMembers: {
		Type:        schema.TypeSet,
		Required:    false,
		Optional:    true,
		Description: "The members of the group",
		MaxItems:    1024,
		Elem: &schema.Resource{
			Schema: groupMemberSchema,
		},
	},
	GroupFieldPermissionSet: {
		Type:        schema.TypeList,
		Required:    false,
		Optional:    true,
		MaxItems:    1,
		Description: "The permission set of the group",
		Elem: &schema.Resource{
			Schema: groupPermissionSet,
		},
	},
}

//NewGroupResourceHandle creates the resource handle for RBAC Groups
func NewGroupResourceHandle() ResourceHandle {
	return &groupResource{
		metaData: ResourceMetaData{
			ResourceName:  ResourceInstanaGroup,
			Schema:        groupSchema,
			SchemaVersion: 0,
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
	return []schema.StateUpgrader{}
}

func (r *groupResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource {
	return api.Groups()
}

func (r *groupResource) SetComputedFields(d *schema.ResourceData) {
	//No computed fields defined
}

func (r *groupResource) UpdateState(d *schema.ResourceData, obj restapi.InstanaDataObject, formatter utils.ResourceNameFormatter) error {
	group := obj.(*restapi.Group)

	members := r.convertGroupMembersToState(group)
	permissions := r.convertPermissionSetToState(group)

	d.Set(GroupFieldName, formatter.UndoFormat(group.Name))
	d.Set(GroupFieldFullName, group.Name)
	d.Set(GroupFieldMembers, members)
	d.Set(GroupFieldPermissionSet, []interface{}{permissions})

	d.SetId(group.ID)
	return nil
}

func (r *groupResource) convertGroupMembersToState(obj *restapi.Group) []map[string]interface{} {
	result := make([]map[string]interface{}, len(obj.Members))
	for i, v := range obj.Members {
		groupMap := map[string]interface{}{
			GroupFieldMemberUserID: v.UserID,
			GroupFieldMemberEmail:  v.Email,
		}
		result[i] = groupMap
	}
	return result
}

func (r *groupResource) convertPermissionSetToState(obj *restapi.Group) map[string]interface{} {
	permissionSet := obj.PermissionSet
	var infraDfQ string
	if obj.PermissionSet.InfraDFQFilter != nil {
		infraDfQ = permissionSet.InfraDFQFilter.ScopeID
	}
	return map[string]interface{}{
		GroupFieldPermissionSetApplicationIDs:          r.convertScopeBindingSliceToState(permissionSet.ApplicationIDs),
		GroupFieldPermissionSetInfraDFQFilter:          infraDfQ,
		GroupFieldPermissionSetKubernetesClusterUUIDs:  r.convertScopeBindingSliceToState(permissionSet.KubernetesClusterUUIDs),
		GroupFieldPermissionSetKubernetesNamespaceUIDs: r.convertScopeBindingSliceToState(permissionSet.KubernetesNamespaceUIDs),
		GroupFieldPermissionSetMobileAppIDs:            r.convertScopeBindingSliceToState(permissionSet.MobileAppIDs),
		GroupFieldPermissionSetWebsiteIDs:              r.convertScopeBindingSliceToState(permissionSet.WebsiteIDs),
		GroupFieldPermissionSetPermissions:             permissionSet.Permissions,
	}
}

func (r *groupResource) convertScopeBindingSliceToState(value []restapi.ScopeBinding) []interface{} {
	result := make([]interface{}, len(value))
	for i, v := range value {
		result[i] = v.ScopeID
	}
	return result
}

func (r *groupResource) MapStateToDataObject(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (restapi.InstanaDataObject, error) {
	name := r.computeFullNameString(d, formatter)
	members := r.convertStateToGroupMembers(d)
	permissionSet := r.convertStateToPermissionSet(d)
	return &restapi.Group{
		ID:            d.Id(),
		Name:          name,
		Members:       members,
		PermissionSet: *permissionSet,
	}, nil
}

func (r *groupResource) computeFullNameString(d *schema.ResourceData, formatter utils.ResourceNameFormatter) string {
	if d.HasChange(GroupFieldName) {
		return formatter.Format(d.Get(GroupFieldName).(string))
	}
	return d.Get(GroupFieldFullName).(string)
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
	return &restapi.APIPermissionSetWithRoles{}
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
	return []restapi.ScopeBinding{}
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
	return []restapi.InstanaPermission{}
}
