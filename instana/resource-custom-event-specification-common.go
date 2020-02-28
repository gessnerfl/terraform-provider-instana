package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

const (
	//CustomEventSpecificationFieldName constant value for the schema field name
	CustomEventSpecificationFieldName = "name"
	//CustomEventSpecificationFieldFullName constant value for the schema field full_name. The field is computed and contains the name which is sent to instana. The computation depends on the configured default_name_prefix and default_name_suffix at provider level
	CustomEventSpecificationFieldFullName = "full_name"
	//CustomEventSpecificationFieldEntityType constant value for the schema field entity type
	CustomEventSpecificationFieldEntityType = "entity_type"
	//CustomEventSpecificationFieldQuery constant value for the schema field query
	CustomEventSpecificationFieldQuery = "query"
	//CustomEventSpecificationFieldTriggering constant value for the schema field triggering
	CustomEventSpecificationFieldTriggering = "triggering"
	//CustomEventSpecificationFieldDescription constant value for the schema field description
	CustomEventSpecificationFieldDescription = "description"
	//CustomEventSpecificationFieldExpirationTime constant value for the schema field expiration_time
	CustomEventSpecificationFieldExpirationTime = "expiration_time"
	//CustomEventSpecificationFieldEnabled constant value for the schema field enabled
	CustomEventSpecificationFieldEnabled = "enabled"

	ruleFieldPrefix = "rule_"

	//CustomEventSpecificationRuleSeverity constant value for the schema field severity of a rule specification
	CustomEventSpecificationRuleSeverity = ruleFieldPrefix + "severity"

	downstreamFieldPrefix = "downstream_"

	//CustomEventSpecificationDownstreamIntegrationIds constant value for the schema field integration_ids of a event specification downstream
	CustomEventSpecificationDownstreamIntegrationIds = downstreamFieldPrefix + "integration_ids"
	//CustomEventSpecificationDownstreamBroadcastToAllAlertingConfigs constant value for the schema field broadcast_to_all_alerting_configs of a event specification downstream
	CustomEventSpecificationDownstreamBroadcastToAllAlertingConfigs = downstreamFieldPrefix + "broadcast_to_all_alerting_configs"
)

var customEventSpecificationSchemaName = &schema.Schema{
	Type:        schema.TypeString,
	Required:    true,
	Description: "Configures the name of the custom event specification",
}
var customEventSpecificationSchemaFullName = &schema.Schema{
	Type:        schema.TypeString,
	Computed:    true,
	Description: "The computed full name of the custom event specification. The field contains the name which is sent to instana. The computation depends on the configured default_name_prefix and default_name_suffix at provider level",
}
var customEventSpecificationSchemaQuery = &schema.Schema{
	Type:        schema.TypeString,
	Required:    false,
	Optional:    true,
	Description: "Configures the dynamic focus query for the custom event specification",
}
var customEventSpecificationSchemaTriggering = &schema.Schema{
	Type:        schema.TypeBool,
	Default:     false,
	Optional:    true,
	Description: "Configures the custom event specification should trigger an incident",
}
var customEventSpecificationSchemaDescription = &schema.Schema{
	Type:        schema.TypeString,
	Required:    false,
	Optional:    true,
	Description: "Configures the description text of the custom event specification",
}
var customEventSpecificationSchemaExpirationTime = &schema.Schema{
	Type:        schema.TypeInt,
	Required:    false,
	Optional:    true,
	Description: "Configures the expiration time (grace period) to wait before the issue is closed",
}
var customEventSpecificationSchemaEnabled = &schema.Schema{
	Type:        schema.TypeBool,
	Default:     true,
	Optional:    true,
	Description: "Configures if the custom event specification is enabled or not",
}
var customEventSpecificationSchemaDownstreamIntegrationIdsLegacy = &schema.Schema{
	Type:     schema.TypeList,
	Required: false,
	Optional: true,
	MinItems: 0,
	MaxItems: 16,
	Elem: &schema.Schema{
		Type: schema.TypeString,
	},
	Description: "Configures the list of integration ids which should be used for downstream reporting",
}

//CustomEventSpecificationSchemaDownstreamIntegrationIds schema definition for integration ids
var CustomEventSpecificationSchemaDownstreamIntegrationIds = &schema.Schema{
	Type:     schema.TypeSet,
	Required: false,
	Optional: true,
	MinItems: 0,
	MaxItems: 16,
	Elem: &schema.Schema{
		Type: schema.TypeString,
	},
	Description: "Configures the list of integration ids which should be used for downstream reporting",
}
var customEventSpecificationSchemaDownstreamBroadcastToAllAlertingConfigs = &schema.Schema{
	Type:        schema.TypeBool,
	Default:     true,
	Optional:    true,
	Description: "Configures the downstream reporting should be sent to all integrations",
}
var customEventSpecificationSchemaRuleSeverity = &schema.Schema{
	Type:         schema.TypeString,
	Required:     true,
	ValidateFunc: validation.StringInSlice([]string{restapi.SeverityWarning.GetTerraformRepresentation(), restapi.SeverityCritical.GetTerraformRepresentation()}, false),
	Description:  "Configures the severity of the rule of the custom event specification",
}

var defaultCustomEventSchemaFieldsV0 = map[string]*schema.Schema{
	CustomEventSpecificationFieldName:                               customEventSpecificationSchemaName,
	CustomEventSpecificationFieldQuery:                              customEventSpecificationSchemaQuery,
	CustomEventSpecificationFieldTriggering:                         customEventSpecificationSchemaTriggering,
	CustomEventSpecificationFieldDescription:                        customEventSpecificationSchemaDescription,
	CustomEventSpecificationFieldExpirationTime:                     customEventSpecificationSchemaExpirationTime,
	CustomEventSpecificationFieldEnabled:                            customEventSpecificationSchemaEnabled,
	CustomEventSpecificationDownstreamIntegrationIds:                customEventSpecificationSchemaDownstreamIntegrationIdsLegacy,
	CustomEventSpecificationDownstreamBroadcastToAllAlertingConfigs: customEventSpecificationSchemaDownstreamBroadcastToAllAlertingConfigs,
	CustomEventSpecificationRuleSeverity:                            customEventSpecificationSchemaRuleSeverity,
}

var defaultCustomEventSchemaFieldsV1 = map[string]*schema.Schema{
	CustomEventSpecificationFieldName:                               customEventSpecificationSchemaName,
	CustomEventSpecificationFieldFullName:                           customEventSpecificationSchemaFullName,
	CustomEventSpecificationFieldQuery:                              customEventSpecificationSchemaQuery,
	CustomEventSpecificationFieldTriggering:                         customEventSpecificationSchemaTriggering,
	CustomEventSpecificationFieldDescription:                        customEventSpecificationSchemaDescription,
	CustomEventSpecificationFieldExpirationTime:                     customEventSpecificationSchemaExpirationTime,
	CustomEventSpecificationFieldEnabled:                            customEventSpecificationSchemaEnabled,
	CustomEventSpecificationDownstreamIntegrationIds:                customEventSpecificationSchemaDownstreamIntegrationIdsLegacy,
	CustomEventSpecificationDownstreamBroadcastToAllAlertingConfigs: customEventSpecificationSchemaDownstreamBroadcastToAllAlertingConfigs,
	CustomEventSpecificationRuleSeverity:                            customEventSpecificationSchemaRuleSeverity,
}

var defaultCustomEventSchemaFields = map[string]*schema.Schema{
	CustomEventSpecificationFieldName:                               customEventSpecificationSchemaName,
	CustomEventSpecificationFieldFullName:                           customEventSpecificationSchemaFullName,
	CustomEventSpecificationFieldQuery:                              customEventSpecificationSchemaQuery,
	CustomEventSpecificationFieldTriggering:                         customEventSpecificationSchemaTriggering,
	CustomEventSpecificationFieldDescription:                        customEventSpecificationSchemaDescription,
	CustomEventSpecificationFieldExpirationTime:                     customEventSpecificationSchemaExpirationTime,
	CustomEventSpecificationFieldEnabled:                            customEventSpecificationSchemaEnabled,
	CustomEventSpecificationDownstreamIntegrationIds:                CustomEventSpecificationSchemaDownstreamIntegrationIds,
	CustomEventSpecificationDownstreamBroadcastToAllAlertingConfigs: customEventSpecificationSchemaDownstreamBroadcastToAllAlertingConfigs,
	CustomEventSpecificationRuleSeverity:                            customEventSpecificationSchemaRuleSeverity,
}

func mergeSchemaMap(mapA map[string]*schema.Schema, mapB map[string]*schema.Schema) map[string]*schema.Schema {
	mergedMap := make(map[string]*schema.Schema)

	for k, v := range mapA {
		mergedMap[k] = v
	}
	for k, v := range mapB {
		mergedMap[k] = v
	}

	return mergedMap
}

func createCustomEventSpecificationFromResourceData(d *schema.ResourceData, formatter utils.ResourceNameFormatter) restapi.CustomEventSpecification {
	name := computeFullCustomEventNameString(d, formatter)
	apiModel := restapi.CustomEventSpecification{
		ID:             d.Id(),
		Name:           name,
		EntityType:     d.Get(CustomEventSpecificationFieldEntityType).(string),
		Query:          GetStringPointerFromResourceData(d, CustomEventSpecificationFieldQuery),
		Triggering:     d.Get(CustomEventSpecificationFieldTriggering).(bool),
		Description:    GetStringPointerFromResourceData(d, CustomEventSpecificationFieldDescription),
		ExpirationTime: GetIntPointerFromResourceData(d, CustomEventSpecificationFieldExpirationTime),
		Enabled:        d.Get(CustomEventSpecificationFieldEnabled).(bool),
	}

	downstreamIntegrationIds := ReadStringSetParameterFromResource(d, CustomEventSpecificationDownstreamIntegrationIds)
	if len(downstreamIntegrationIds) > 0 {
		apiModel.Downstream = &restapi.EventSpecificationDownstream{
			IntegrationIds:                downstreamIntegrationIds,
			BroadcastToAllAlertingConfigs: d.Get(CustomEventSpecificationDownstreamBroadcastToAllAlertingConfigs).(bool),
		}
	}
	return apiModel
}

func computeFullCustomEventNameString(d *schema.ResourceData, formatter utils.ResourceNameFormatter) string {
	if d.HasChange(CustomEventSpecificationFieldName) {
		return formatter.Format(d.Get(CustomEventSpecificationFieldName).(string))
	}
	return d.Get(CustomEventSpecificationFieldFullName).(string)
}

func updateStateForBasicCustomEventSpecification(d *schema.ResourceData, spec restapi.CustomEventSpecification) {
	d.SetId(spec.ID)
	d.Set(CustomEventSpecificationFieldFullName, spec.Name)
	d.Set(CustomEventSpecificationFieldQuery, spec.Query)
	d.Set(CustomEventSpecificationFieldEntityType, spec.EntityType)
	d.Set(CustomEventSpecificationFieldTriggering, spec.Triggering)
	d.Set(CustomEventSpecificationFieldDescription, spec.Description)
	d.Set(CustomEventSpecificationFieldExpirationTime, spec.ExpirationTime)
	d.Set(CustomEventSpecificationFieldEnabled, spec.Enabled)

	if spec.Downstream != nil {
		d.Set(CustomEventSpecificationDownstreamIntegrationIds, schema.NewSet(schema.HashString, ConvertStringToInterfaceSlice(spec.Downstream.IntegrationIds)))
		d.Set(CustomEventSpecificationDownstreamBroadcastToAllAlertingConfigs, spec.Downstream.BroadcastToAllAlertingConfigs)
	}
}

func migrateCustomEventConfigFullNameInStateFromV0toV1(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	rawState[CustomEventSpecificationFieldFullName] = rawState[CustomEventSpecificationFieldName]
	return rawState, nil
}

func migrateCustomEventIntegrationIdsInStateFromV1toV2(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	integrationIds, found := rawState[CustomEventSpecificationDownstreamIntegrationIds]
	if found {
		rawState[CustomEventSpecificationDownstreamIntegrationIds] = schema.NewSet(schema.HashSchema(CustomEventSpecificationSchemaDownstreamIntegrationIds.Elem.(*schema.Schema)), integrationIds.([]interface{}))
	}
	return rawState, nil
}
