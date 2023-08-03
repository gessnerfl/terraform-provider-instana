package instana

import (
	"context"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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

	downstreamFieldPrefix                                           = "downstream_"
	customEventSpecificationDownstreamIntegrationIds                = downstreamFieldPrefix + "integration_ids"
	customEventSpecificationDownstreamBroadcastToAllAlertingConfigs = downstreamFieldPrefix + "broadcast_to_all_alerting_configs"
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
var customEventSpecificationSchemaDownstreamIntegrationIds = &schema.Schema{
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
var customEventSpecificationSchemaDownstreamBroadcastToAllAlertingConfigs = &schema.Schema{
	Type:        schema.TypeBool,
	Default:     true,
	Optional:    true,
	Description: "Configures the downstream reporting should be sent to all integrations",
}
var customEventSpecificationSchemaRuleSeverity = &schema.Schema{
	Type:         schema.TypeString,
	Required:     true,
	ValidateFunc: validation.StringInSlice(restapi.SupportedSeverities.TerraformRepresentations(), false),
	Description:  "Configures the severity of the rule of the custom event specification",
}

var defaultCustomEventSchemaFieldsV0 = map[string]*schema.Schema{
	CustomEventSpecificationFieldName:                               customEventSpecificationSchemaName,
	CustomEventSpecificationFieldQuery:                              customEventSpecificationSchemaQuery,
	CustomEventSpecificationFieldTriggering:                         customEventSpecificationSchemaTriggering,
	CustomEventSpecificationFieldDescription:                        customEventSpecificationSchemaDescription,
	CustomEventSpecificationFieldExpirationTime:                     customEventSpecificationSchemaExpirationTime,
	CustomEventSpecificationFieldEnabled:                            customEventSpecificationSchemaEnabled,
	customEventSpecificationDownstreamIntegrationIds:                customEventSpecificationSchemaDownstreamIntegrationIds,
	customEventSpecificationDownstreamBroadcastToAllAlertingConfigs: customEventSpecificationSchemaDownstreamBroadcastToAllAlertingConfigs,
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
	customEventSpecificationDownstreamIntegrationIds:                customEventSpecificationSchemaDownstreamIntegrationIds,
	customEventSpecificationDownstreamBroadcastToAllAlertingConfigs: customEventSpecificationSchemaDownstreamBroadcastToAllAlertingConfigs,
	CustomEventSpecificationRuleSeverity:                            customEventSpecificationSchemaRuleSeverity,
}

var defaultCustomEventSchemaFields = map[string]*schema.Schema{
	CustomEventSpecificationFieldName:           customEventSpecificationSchemaName,
	CustomEventSpecificationFieldFullName:       customEventSpecificationSchemaFullName,
	CustomEventSpecificationFieldQuery:          customEventSpecificationSchemaQuery,
	CustomEventSpecificationFieldTriggering:     customEventSpecificationSchemaTriggering,
	CustomEventSpecificationFieldDescription:    customEventSpecificationSchemaDescription,
	CustomEventSpecificationFieldExpirationTime: customEventSpecificationSchemaExpirationTime,
	CustomEventSpecificationFieldEnabled:        customEventSpecificationSchemaEnabled,
	CustomEventSpecificationRuleSeverity:        customEventSpecificationSchemaRuleSeverity,
}

type customEventSpecificationCommons struct{}

func (c *customEventSpecificationCommons) createCustomEventSpecificationFromResourceData(d *schema.ResourceData, formatter utils.ResourceNameFormatter) *restapi.CustomEventSpecification {
	name := c.computeFullCustomEventNameString(d, formatter)
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
	return &apiModel
}

func (c *customEventSpecificationCommons) computeFullCustomEventNameString(d *schema.ResourceData, formatter utils.ResourceNameFormatter) string {
	if d.HasChange(CustomEventSpecificationFieldName) {
		return formatter.Format(d.Get(CustomEventSpecificationFieldName).(string))
	}
	return d.Get(CustomEventSpecificationFieldFullName).(string)
}

func (c *customEventSpecificationCommons) getDataForBasicCustomEventSpecification(spec *restapi.CustomEventSpecification, formatter utils.ResourceNameFormatter) map[string]interface{} {
	return map[string]interface{}{
		CustomEventSpecificationFieldName:           formatter.UndoFormat(spec.Name),
		CustomEventSpecificationFieldFullName:       spec.Name,
		CustomEventSpecificationFieldQuery:          spec.Query,
		CustomEventSpecificationFieldEntityType:     spec.EntityType,
		CustomEventSpecificationFieldTriggering:     spec.Triggering,
		CustomEventSpecificationFieldDescription:    spec.Description,
		CustomEventSpecificationFieldExpirationTime: spec.ExpirationTime,
		CustomEventSpecificationFieldEnabled:        spec.Enabled,
	}
}

func (c *customEventSpecificationCommons) migrateCustomEventConfigFullNameInStateFromV0toV1(_ context.Context, rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	rawState[CustomEventSpecificationFieldFullName] = rawState[CustomEventSpecificationFieldName]
	return rawState, nil
}

func (c *customEventSpecificationCommons) migrateCustomEventConfigFullStateFromV1toV2AndRemoveDownstreamConfiguration(_ context.Context, rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	delete(rawState, customEventSpecificationDownstreamIntegrationIds)
	delete(rawState, customEventSpecificationDownstreamBroadcastToAllAlertingConfigs)
	return rawState, nil
}
