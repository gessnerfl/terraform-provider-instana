package instana

import (
	"context"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"
	"strings"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceInstanaAlertingConfig the name of the terraform-provider-instana resource to manage alerting configurations
const ResourceInstanaAlertingConfig = "instana_alerting_config"

const (
	//AlertingConfigFieldAlertName constant value for the schema field alert_name
	AlertingConfigFieldAlertName = "alert_name"
	//AlertingConfigFieldFullAlertName constant value for the schema field full_alert_name
	//Deprecated
	AlertingConfigFieldFullAlertName = "full_alert_name"
	//AlertingConfigFieldIntegrationIds constant value for the schema field integration_ids
	AlertingConfigFieldIntegrationIds = "integration_ids"
	//AlertingConfigFieldEventFilterQuery constant value for the schema field event_filter_query
	AlertingConfigFieldEventFilterQuery = "event_filter_query"
	//AlertingConfigFieldEventFilterEventTypes constant value for the schema field event_filter_event_types
	AlertingConfigFieldEventFilterEventTypes = "event_filter_event_types"
	//AlertingConfigFieldEventFilterRuleIDs constant value for the schema field event_filter_rule_ids
	AlertingConfigFieldEventFilterRuleIDs = "event_filter_rule_ids"
)

var supportedEventTypes = convertSupportedEventTypesToStringSlice()

func convertSupportedEventTypesToStringSlice() []string {
	result := make([]string, len(restapi.SupportedAlertEventTypes))
	for i, t := range restapi.SupportedAlertEventTypes {
		result[i] = string(t)
	}
	return result
}

// AlertingConfigSchemaAlertName schema field definition of instana_alerting_config field alert_name
var AlertingConfigSchemaAlertName = &schema.Schema{
	Type:         schema.TypeString,
	Required:     true,
	Description:  "Configures the alert name of the alerting configuration",
	ValidateFunc: validation.StringLenBetween(1, 200),
}

// AlertingConfigSchemaFullAlertName schema field definition of instana_alerting_config field full_alert_name
// Deprecated
var AlertingConfigSchemaFullAlertName = &schema.Schema{
	Type:        schema.TypeString,
	Computed:    true,
	Description: "The the full alert name field of the alerting configuration. The field is computed and contains the name which is sent to instana. The computation depends on the configured default_name_prefix and default_name_suffix at provider level",
}

// AlertingConfigSchemaIntegrationIds schema field definition of instana_alerting_config field integration_ids
var AlertingConfigSchemaIntegrationIds = &schema.Schema{
	Type:     schema.TypeSet,
	MinItems: 0,
	MaxItems: 1024,
	Elem: &schema.Schema{
		Type: schema.TypeString,
	},
	Required:    true,
	Description: "Configures the list of Integration IDs (Alerting Channels).",
}

// AlertingConfigSchemaEventFilterQuery schema field definition of instana_alerting_config field event_filter_query
var AlertingConfigSchemaEventFilterQuery = &schema.Schema{
	Type:         schema.TypeString,
	Required:     false,
	Optional:     true,
	Description:  "Configures a filter query to to filter rules or event types for a limited set of entities",
	ValidateFunc: validation.StringLenBetween(0, 2048),
}

// AlertingConfigSchemaEventFilterEventTypes schema field definition of instana_alerting_config field event_filter_event_types
var AlertingConfigSchemaEventFilterEventTypes = &schema.Schema{
	Type:     schema.TypeSet,
	MinItems: 0,
	MaxItems: len(supportedEventTypes),
	Elem: &schema.Schema{
		Type:         schema.TypeString,
		ValidateFunc: validation.StringInSlice(supportedEventTypes, false),
	},
	Required:      false,
	Optional:      true,
	ConflictsWith: []string{AlertingConfigFieldEventFilterRuleIDs},
	Description:   "Configures the list of Event Types IDs which should trigger an alert.",
}

// AlertingConfigSchemaEventFilterRuleIDs schema field definition of instana_alerting_config field event_filter_rule_ids
var AlertingConfigSchemaEventFilterRuleIDs = &schema.Schema{
	Type:     schema.TypeSet,
	MinItems: 0,
	MaxItems: 1024,
	Elem: &schema.Schema{
		Type: schema.TypeString,
	},
	Required:      false,
	Optional:      true,
	ConflictsWith: []string{AlertingConfigFieldEventFilterEventTypes},
	Description:   "Configures the list of Rule IDs which should trigger an alert.",
}

// NewAlertingConfigResourceHandle creates the resource handle for Alerting Configuration
func NewAlertingConfigResourceHandle() ResourceHandle[*restapi.AlertingConfiguration] {
	return &alertingConfigResource{
		metaData: ResourceMetaData{
			ResourceName: ResourceInstanaAlertingConfig,
			Schema: map[string]*schema.Schema{
				AlertingConfigFieldAlertName:             AlertingConfigSchemaAlertName,
				AlertingConfigFieldIntegrationIds:        AlertingConfigSchemaIntegrationIds,
				AlertingConfigFieldEventFilterQuery:      AlertingConfigSchemaEventFilterQuery,
				AlertingConfigFieldEventFilterEventTypes: AlertingConfigSchemaEventFilterEventTypes,
				AlertingConfigFieldEventFilterRuleIDs:    AlertingConfigSchemaEventFilterRuleIDs,
			},
			SchemaVersion: 2,
		},
	}
}

type alertingConfigResource struct {
	metaData ResourceMetaData
}

func (r *alertingConfigResource) MetaData() *ResourceMetaData {
	return &r.metaData
}

func (r *alertingConfigResource) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{
		{
			Type: r.schemaV0().CoreConfigSchema().ImpliedType(),
			Upgrade: func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
				return rawState, nil
			},
			Version: 0,
		},
		{
			Type:    r.schemaV1().CoreConfigSchema().ImpliedType(),
			Upgrade: r.stateUpgradeV1,
			Version: 1,
		},
	}
}

func (r *alertingConfigResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.AlertingConfiguration] {
	return api.AlertingConfigurations()
}

func (r *alertingConfigResource) SetComputedFields(_ *schema.ResourceData) error {
	return nil
}

func (r *alertingConfigResource) UpdateState(d *schema.ResourceData, config *restapi.AlertingConfiguration) error {
	d.SetId(config.ID)
	return tfutils.UpdateState(d, map[string]interface{}{
		AlertingConfigFieldAlertName:             config.AlertName,
		AlertingConfigFieldIntegrationIds:        config.IntegrationIDs,
		AlertingConfigFieldEventFilterQuery:      config.EventFilteringConfiguration.Query,
		AlertingConfigFieldEventFilterEventTypes: r.convertEventTypesToHarmonizedStringRepresentation(config.EventFilteringConfiguration.EventTypes),
		AlertingConfigFieldEventFilterRuleIDs:    config.EventFilteringConfiguration.RuleIDs,
	})
}

func (r *alertingConfigResource) convertEventTypesToHarmonizedStringRepresentation(input []restapi.AlertEventType) []string {
	result := make([]string, len(input))
	for i, v := range input {
		value := strings.ToLower(string(v))
		result[i] = value
	}
	return result
}

func (r *alertingConfigResource) MapStateToDataObject(d *schema.ResourceData) (*restapi.AlertingConfiguration, error) {
	query := GetStringPointerFromResourceData(d, AlertingConfigFieldEventFilterQuery)

	return &restapi.AlertingConfiguration{
		ID:             d.Id(),
		AlertName:      d.Get(AlertingConfigFieldAlertName).(string),
		IntegrationIDs: ReadStringSetParameterFromResource(d, AlertingConfigFieldIntegrationIds),
		EventFilteringConfiguration: restapi.EventFilteringConfiguration{
			Query:      query,
			RuleIDs:    ReadStringSetParameterFromResource(d, AlertingConfigFieldEventFilterRuleIDs),
			EventTypes: r.readEventTypesFromResourceData(d),
		},
	}, nil
}

func (r *alertingConfigResource) readEventTypesFromResourceData(d *schema.ResourceData) []restapi.AlertEventType {
	rawData := ReadStringSetParameterFromResource(d, AlertingConfigFieldEventFilterEventTypes)
	result := make([]restapi.AlertEventType, len(rawData))
	for i, v := range rawData {
		value := strings.ToLower(v)
		result[i] = restapi.AlertEventType(value)
	}
	return result
}

func (r *alertingConfigResource) schemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			AlertingConfigFieldAlertName:     AlertingConfigSchemaAlertName,
			AlertingConfigFieldFullAlertName: AlertingConfigSchemaFullAlertName,
			AlertingConfigFieldIntegrationIds: {
				Type:     schema.TypeList,
				MinItems: 0,
				MaxItems: 1024,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "Configures the list of Integration IDs (Alerting Channels).",
			},
			AlertingConfigFieldEventFilterQuery: AlertingConfigSchemaEventFilterQuery,
			AlertingConfigFieldEventFilterEventTypes: {
				Type:     schema.TypeList,
				MinItems: 0,
				MaxItems: 6,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(supportedEventTypes, false),
				},
				Required:      false,
				Optional:      true,
				ConflictsWith: []string{AlertingConfigFieldEventFilterRuleIDs},
				Description:   "Configures the list of Event Types IDs which should trigger an alert.",
			},
			AlertingConfigFieldEventFilterRuleIDs: {
				Type:     schema.TypeList,
				MinItems: 0,
				MaxItems: 1024,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:      false,
				Optional:      true,
				ConflictsWith: []string{AlertingConfigFieldEventFilterEventTypes},
				Description:   "Configures the list of Rule IDs which should trigger an alert.",
			},
		},
	}
}

func (r *alertingConfigResource) stateUpgradeV1(_ context.Context, state map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	if _, ok := state[AlertingConfigFieldFullAlertName]; ok {
		state[AlertingConfigFieldAlertName] = state[AlertingConfigFieldFullAlertName]
		delete(state, AlertingConfigFieldFullAlertName)
	}
	return state, nil
}

func (r *alertingConfigResource) schemaV1() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			AlertingConfigFieldAlertName:             AlertingConfigSchemaAlertName,
			AlertingConfigFieldFullAlertName:         AlertingConfigSchemaFullAlertName,
			AlertingConfigFieldIntegrationIds:        AlertingConfigSchemaIntegrationIds,
			AlertingConfigFieldEventFilterQuery:      AlertingConfigSchemaEventFilterQuery,
			AlertingConfigFieldEventFilterEventTypes: AlertingConfigSchemaEventFilterEventTypes,
			AlertingConfigFieldEventFilterRuleIDs:    AlertingConfigSchemaEventFilterRuleIDs,
		},
	}
}
