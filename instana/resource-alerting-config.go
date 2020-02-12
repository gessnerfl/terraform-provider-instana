package instana

import (
	"strings"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

//ResourceInstanaAlertingConfig the name of the terraform-provider-instana resource to manage alerting configurations
const ResourceInstanaAlertingConfig = "instana_alerting_config"

const (
	//AlertingConfigFieldAlertName constant value for the schema field alert_name
	AlertingConfigFieldAlertName = "alert_name"
	//AlertingConfigFieldFullAlertName constant value for the schema field full_alert_name
	AlertingConfigFieldFullAlertName = "full_alert_name"
	//AlertingConfigFieldIntegrationIds constant value for the schema field integration_ids
	AlertingConfigFieldIntegrationIds = "integration_ids"
	//AlertingConfigFieldCustomPayload constant value for the schema field custom_payload
	AlertingConfigFieldCustomPayload = "custom_payload"
	//AlertingConfigFieldEventFilterQuery constant value for the schema field event_filter_query
	AlertingConfigFieldEventFilterQuery = "event_filter_query"
	//AlertingConfigFieldEventFilterEventTypes constant value for the schema field event_filter_event_types
	AlertingConfigFieldEventFilterEventTypes = "event_filter_event_types"
	//AlertingConfigFieldEventFilterRuleIDs constant value for the schema field event_filter_rule_ids
	AlertingConfigFieldEventFilterRuleIDs = "event_filter_rule_ids"
)

var supportedEventTypes = convertSupportedEventTypesToStringSlice()

//NewAlertingConfigResourceHandle creates the resource handle for Alerting Configuration
func NewAlertingConfigResourceHandle() *ResourceHandle {
	return &ResourceHandle{
		ResourceName: ResourceInstanaAlertingConfig,
		Schema: map[string]*schema.Schema{
			AlertingConfigFieldAlertName: &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Configures the alert name of the alerting configuration",
				ValidateFunc: validation.StringLenBetween(1, 200),
			},
			AlertingConfigFieldFullAlertName: &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The the full alert name field of the alerting configuration. The field is computed and contains the name which is sent to instana. The computation depends on the configured default_name_prefix and default_name_suffix at provider level",
			},
			AlertingConfigFieldIntegrationIds: &schema.Schema{
				Type:     schema.TypeList,
				MinItems: 0,
				MaxItems: 1024,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "Configures the list of Integration IDs (Alerting Channels).",
			},
			AlertingConfigFieldEventFilterQuery: &schema.Schema{
				Type:         schema.TypeString,
				Required:     false,
				Optional:     true,
				Description:  "Configures a filter query to to filter rules or event types for a limited set of entities",
				ValidateFunc: validation.StringLenBetween(0, 2048),
			},
			AlertingConfigFieldEventFilterEventTypes: &schema.Schema{
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
			AlertingConfigFieldEventFilterRuleIDs: &schema.Schema{
				Type:     schema.TypeList,
				MinItems: 0,
				MaxItems: 6,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:      false,
				Optional:      true,
				ConflictsWith: []string{AlertingConfigFieldEventFilterEventTypes},
				Description:   "Configures the list of Rule IDs which should trigger an alert.",
			},
			AlertingConfigFieldCustomPayload: &schema.Schema{
				Type:         schema.TypeString,
				Required:     false,
				Optional:     true,
				Description:  "Configures a custom payload for the alerting configuration",
				ValidateFunc: validation.StringLenBetween(0, 65536),
			},
		},
		RestResourceFactory:  func(api restapi.InstanaAPI) restapi.RestResource { return api.AlertingConfigurations() },
		UpdateState:          updateStateForAlertingConfig,
		MapStateToDataObject: mapStateToDataObjectForAlertingConfig,
	}
}

func updateStateForAlertingConfig(d *schema.ResourceData, obj restapi.InstanaDataObject) error {
	config := obj.(restapi.AlertingConfiguration)
	d.Set(AlertingConfigFieldFullAlertName, config.AlertName)
	d.Set(AlertingConfigFieldIntegrationIds, config.IntegrationIDs)
	d.Set(AlertingConfigFieldCustomPayload, config.CustomPayload)
	d.Set(AlertingConfigFieldEventFilterQuery, config.EventFilteringConfiguration.Query)
	d.Set(AlertingConfigFieldEventFilterEventTypes, convertEventTypesToHarmonizedStringRepresentation(config.EventFilteringConfiguration.EventTypes))
	d.Set(AlertingConfigFieldEventFilterRuleIDs, config.EventFilteringConfiguration.RuleIDs)
	d.SetId(config.ID)
	return nil
}

func convertEventTypesToHarmonizedStringRepresentation(input []restapi.AlertEventType) []string {
	result := make([]string, len(input))
	for i, v := range input {
		value := strings.ToUpper(string(v))
		result[i] = value
	}
	return result
}

func mapStateToDataObjectForAlertingConfig(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (restapi.InstanaDataObject, error) {
	name := computeFullAlertingConfigAlertNameString(d, formatter)
	customPayload := GetStringPointerFromResourceData(d, AlertingConfigFieldCustomPayload)
	query := GetStringPointerFromResourceData(d, AlertingConfigFieldEventFilterQuery)

	return restapi.AlertingConfiguration{
		ID:             d.Id(),
		AlertName:      name,
		IntegrationIDs: ReadStringArrayParameterFromResource(d, AlertingConfigFieldIntegrationIds),
		CustomPayload:  customPayload,
		EventFilteringConfiguration: restapi.EventFilteringConfiguration{
			Query:      query,
			RuleIDs:    ReadStringArrayParameterFromResource(d, AlertingConfigFieldEventFilterRuleIDs),
			EventTypes: readEventTypesFromResourceData(d),
		},
	}, nil
}

func readEventTypesFromResourceData(d *schema.ResourceData) []restapi.AlertEventType {
	rawData := ReadStringArrayParameterFromResource(d, AlertingConfigFieldEventFilterEventTypes)
	result := make([]restapi.AlertEventType, len(rawData))
	for i, v := range rawData {
		value := strings.ToUpper(v)
		result[i] = restapi.AlertEventType(value)
	}
	return result
}

func computeFullAlertingConfigAlertNameString(d *schema.ResourceData, formatter utils.ResourceNameFormatter) string {
	if d.HasChange(AlertingConfigFieldAlertName) {
		return formatter.Format(d.Get(AlertingConfigFieldAlertName).(string))
	}
	return d.Get(AlertingConfigFieldFullAlertName).(string)
}

func convertSupportedEventTypesToStringSlice() []string {
	result := make([]string, len(restapi.SupportedAlertEventTypes))
	for i, t := range restapi.SupportedAlertEventTypes {
		result[i] = string(t)
	}
	return result
}
