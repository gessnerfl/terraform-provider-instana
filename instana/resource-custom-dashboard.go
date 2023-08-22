package instana

import (
	"context"
	"encoding/json"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceInstanaCustomDashboard the name of the terraform-provider-instana resource to manage custom dashboards
const ResourceInstanaCustomDashboard = "instana_custom_dashboard"

const (
	//CustomDashboardFieldTitle constant value for the schema field title
	CustomDashboardFieldTitle = "title"
	//CustomDashboardFieldFullTitle constant value for the computed schema field full_title
	//Deprecated
	CustomDashboardFieldFullTitle = "full_title"
	//CustomDashboardFieldAccessRule constant value for the schema field access_rule
	CustomDashboardFieldAccessRule = "access_rule"
	//CustomDashboardFieldAccessRuleAccessType constant value for the schema field access_rule.access_type
	CustomDashboardFieldAccessRuleAccessType = "access_type"
	//CustomDashboardFieldAccessRuleRelatedID constant value for the schema field access_rule.related_id
	CustomDashboardFieldAccessRuleRelatedID = "related_id"
	//CustomDashboardFieldAccessRuleRelationType constant value for the schema field access_rule.relation_type
	CustomDashboardFieldAccessRuleRelationType = "relation_type"
	//CustomDashboardFieldWidgets constant value for the schema field widgets
	CustomDashboardFieldWidgets = "widgets"
)

var (
	customDashboardSchemaTitle = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The title of the custom dashboard",
	}
	//Deprecated
	customDashboardSchemaFullTitle = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The full title of the custom dashboard. The field is computed and contains the name which is sent to instana. The computation depends on the configured default_name_prefix and default_name_suffix at provider level",
	}
	customDashboardSchemaAccessRule = &schema.Schema{
		Type:        schema.TypeList,
		Required:    true,
		Description: "The access rules applied to the custom dashboard",
		MinItems:    1,
		MaxItems:    64,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				CustomDashboardFieldAccessRuleAccessType: {
					Type:         schema.TypeString,
					Required:     true,
					Description:  "The access type of the given access rule",
					ValidateFunc: validation.StringInSlice(restapi.SupportedAccessTypes.ToStringSlice(), false),
				},
				CustomDashboardFieldAccessRuleRelatedID: {
					Type:         schema.TypeString,
					Optional:     true,
					Description:  "The id of the related entity (user, api_token, etc.) of the given access rule",
					ValidateFunc: validation.StringLenBetween(0, 64),
				},
				CustomDashboardFieldAccessRuleRelationType: {
					Type:         schema.TypeString,
					Required:     true,
					Description:  "The relation type of the given access rule",
					ValidateFunc: validation.StringInSlice(restapi.SupportedRelationTypes.ToStringSlice(), false),
				},
			},
		},
	}
	customDashboardSchemaWidgets = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The json array containing the widgets configured for the custom dashboard",
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			return NormalizeJSONString(old) == NormalizeJSONString(new)
		},
		StateFunc: func(val interface{}) string {
			return NormalizeJSONString(val.(string))
		},
	}
)

// NewCustomDashboardResourceHandle creates the resource handle for RBAC Groups
func NewCustomDashboardResourceHandle() ResourceHandle[*restapi.CustomDashboard] {
	return &customDashboardResource{
		metaData: ResourceMetaData{
			ResourceName: ResourceInstanaCustomDashboard,
			Schema: map[string]*schema.Schema{
				CustomDashboardFieldTitle:      customDashboardSchemaTitle,
				CustomDashboardFieldAccessRule: customDashboardSchemaAccessRule,
				CustomDashboardFieldWidgets:    customDashboardSchemaWidgets,
			},
			SchemaVersion: 1,
		},
	}
}

type customDashboardResource struct {
	metaData ResourceMetaData
}

func (r *customDashboardResource) MetaData() *ResourceMetaData {
	return &r.metaData
}

func (r *customDashboardResource) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{
		{
			Type:    r.schemaV0().CoreConfigSchema().ImpliedType(),
			Upgrade: r.stateUpgradeV0,
			Version: 0,
		},
	}
}

func (r *customDashboardResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.CustomDashboard] {
	return api.CustomDashboards()
}

func (r *customDashboardResource) SetComputedFields(_ *schema.ResourceData) error {
	return nil
}

func (r *customDashboardResource) UpdateState(d *schema.ResourceData, dashboard *restapi.CustomDashboard, _ utils.ResourceNameFormatter) error {
	widgetsBytes, _ := dashboard.Widgets.MarshalJSON()
	widgets := NormalizeJSONString(string(widgetsBytes))

	d.SetId(dashboard.ID)
	return tfutils.UpdateState(d, map[string]interface{}{
		CustomDashboardFieldTitle:      dashboard.Title,
		CustomDashboardFieldWidgets:    widgets,
		CustomDashboardFieldAccessRule: r.mapAccessRuleToState(dashboard),
	})
}

func (r *customDashboardResource) mapAccessRuleToState(dashboard *restapi.CustomDashboard) []map[string]interface{} {
	result := make([]map[string]interface{}, len(dashboard.AccessRules))
	for i, r := range dashboard.AccessRules {
		result[i] = map[string]interface{}{
			CustomDashboardFieldAccessRuleAccessType:   string(r.AccessType),
			CustomDashboardFieldAccessRuleRelatedID:    r.RelatedID,
			CustomDashboardFieldAccessRuleRelationType: string(r.RelationType),
		}
	}
	return result
}

func (r *customDashboardResource) MapStateToDataObject(d *schema.ResourceData, _ utils.ResourceNameFormatter) (*restapi.CustomDashboard, error) {
	accessRules := r.mapAccessRulesFromState(d)

	widgets := d.Get(CustomDashboardFieldWidgets).(string)
	return &restapi.CustomDashboard{
		ID:          d.Id(),
		Title:       d.Get(CustomDashboardFieldTitle).(string),
		AccessRules: accessRules,
		Widgets:     json.RawMessage(widgets),
	}, nil
}

func (r *customDashboardResource) mapAccessRulesFromState(d *schema.ResourceData) []restapi.AccessRule {
	if val, ok := d.GetOk(CustomDashboardFieldAccessRule); ok {
		rules := val.([]interface{})
		result := make([]restapi.AccessRule, len(rules))

		for i, r := range rules {
			ruleMap := r.(map[string]interface{})
			var relatedId *string
			if val, ok := ruleMap[CustomDashboardFieldAccessRuleRelatedID]; ok && !utils.IsBlank(val.(string)) {
				relatedIdStr := val.(string)
				relatedId = &relatedIdStr
			}
			rule := restapi.AccessRule{
				AccessType:   restapi.AccessType(ruleMap[CustomDashboardFieldAccessRuleAccessType].(string)),
				RelatedID:    relatedId,
				RelationType: restapi.RelationType(ruleMap[CustomDashboardFieldAccessRuleRelationType].(string)),
			}
			result[i] = rule
		}
		return result
	}
	return []restapi.AccessRule{}
}

func (r *customDashboardResource) stateUpgradeV0(_ context.Context, state map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	if _, ok := state[CustomDashboardFieldFullTitle]; ok {
		state[CustomDashboardFieldTitle] = state[CustomDashboardFieldFullTitle]
		delete(state, CustomDashboardFieldFullTitle)
	}
	return state, nil
}

func (r *customDashboardResource) schemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			CustomDashboardFieldTitle:      customDashboardSchemaTitle,
			CustomDashboardFieldFullTitle:  customDashboardSchemaFullTitle,
			CustomDashboardFieldAccessRule: customDashboardSchemaAccessRule,
			CustomDashboardFieldWidgets:    customDashboardSchemaWidgets,
		},
	}
}
