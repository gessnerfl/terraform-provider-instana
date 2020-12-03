package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

// Name of the resource to manage application alert
const ResourceApplicationAlertConfigs = "instana_application_alert_config"

const (
	ApplicationAlertConfigsFieldAlertName     = "alert_name"
	ApplicationAlertConfigsFieldFullAlertName = "full_alert_name"
	ApplicationAlertConfigsFieldApplicationId = "application_id"
	ApplicationAlertConfigsFieldDescription   = "description"

	ApplicationAlertConfigsFieldRuleAlertType   = "rule_alert_type"
	ApplicationAlertConfigsFieldRuleMetricName  = "rule_metric_name"
	ApplicationAlertConfigsFieldRuleAggregation = "rule_aggregation"
	// Addition for alert_type = logs
	ApplicationAlertConfigsFieldRuleMessage  = "rule_message"
	ApplicationAlertConfigsFieldRuleOperator = "rule_operator"
	ApplicationAlertConfigsFieldRuleLevel    = "rule_level"

	ApplicationAlertConfigsFieldAlertChannelIds = "integration_ids"

	ApplicationAlertConfigsFieldThresholdValue    = "threshold_value"
	ApplicationAlertConfigsFieldThresholdType     = "threshold_type"
	ApplicationAlertConfigsFieldThresholdOperator = "threshold_operator"

	ApplicationAlertConfigsFieldTagFilter = "tag_filter"
)

var ApplicationAlertConfigSchemaAlertName = &schema.Schema{
	Type:         schema.TypeString,
	Required:     true,
	Description:  "Configures the alert name of the application alerting configuration",
	ValidateFunc: validation.StringLenBetween(1, 200),
}

var ApplicationAlertConfigSchemaFullAlertName = &schema.Schema{
	Type:        schema.TypeString,
	Computed:    true,
	Description: "The the full alert name field of the application Alert Config. The field is computed and contains the name which is sent to instana. The computation depends on the configured default_name_prefix and default_name_suffix at provider level",
}

var ApplicationAlertConfigApplicationId = &schema.Schema{
	Type:         schema.TypeString,
	Required:     true,
	Description:  "Configures the application id of the application to attache the alert to",
	ValidateFunc: validation.StringLenBetween(1, 200),
}
var ApplicationAlertConfigSchemaDescription = &schema.Schema{
	Type:         schema.TypeString,
	Required:     true,
	Description:  "Configures the description name of the application alerting configuration",
	ValidateFunc: validation.StringLenBetween(1, 200),
}
var ApplicationAlertConfigsRuleAlertType = &schema.Schema{
	Type:        schema.TypeString,
	Required:    true,
	Description: "Type of the Alert Rule",
	// TODO: Just support slowness for now
	ValidateFunc: validation.StringInSlice([]string{"slowness", "logs"}, true),
}

var ApplicationAlertConfigsRuleMetricName = &schema.Schema{
	Type:         schema.TypeString,
	Required:     true,
	Description:  "Name of the Application Alert Config",
	ValidateFunc: validation.StringLenBetween(0, 256),
}

var ApplicationAlertConfigsRuleAggregation = &schema.Schema{
	Type:        schema.TypeString,
	Required:    true,
	Description: "Aggregation for given Metrics",
	ValidateFunc: validation.StringInSlice([]string{"sum", "mean", "max", "min", "p25",
		"p50", "p75", "p90", "p95", "p98", "p99", "distinct_count"}, true),
}

var ApplicationAlertConfigsRuleMessage = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	Description:  "Message of the alert",
	ValidateFunc: validation.StringLenBetween(0, 256),
}

var ApplicationAlertConfigsRuleLevel = &schema.Schema{
	Type:         schema.TypeString,
	Optional:     true,
	Description:  "Rule Alerting Level",
	ValidateFunc: validation.StringInSlice([]string{"WARN", "ERROR"}, true),
}

var ApplicationAlertConfigsIntergrationIds = &schema.Schema{
	Type:     schema.TypeSet,
	MinItems: 0,
	MaxItems: 1024,
	Elem: &schema.Schema{
		Type: schema.TypeString,
	},
	Required:    true,
	Description: "Configures the list of Alertchannel IDs",
}
var ApplicationAlertConfigsThresholdValue = &schema.Schema{
	Type:        schema.TypeFloat,
	Required:    true,
	Description: "The expected condition value to fulfill the rule",
}

var ApplicationAlertConfigsThresholdType = &schema.Schema{
	Type:        schema.TypeString,
	Required:    true,
	Description: "The threshold Type (Currently just 'staticThreshold' is allowed) ",
}

var ApplicationAlertConfigsThresholdOperator = &schema.Schema{
	Type:        schema.TypeString,
	Optional:    true,
	Description: "The Operator to compare Threshold",
}

const (
	TagFilterName         = "name"
	TagFilterEntity       = "entity"
	TagFilterStringValue  = "string_value"
	TagFilterNumberValue  = "number_value"
	TagFilterBooleanValue = "boolean_value"
	TagFilterOperator     = "operator"
)

var ApplicationFilterOperator = &schema.Schema{
	Type:     schema.TypeString,
	Optional: true,
	ValidateFunc: validation.StringInSlice(
		[]string{"EQUALS",
			"CONTAINS",
			"LESS_THAN",
			"LESS_OR_EQUAL_THAN",
			"GREATER_THAN",
			"GREATER_OR_EQUAL_THAN",
			"NOT_EMPTY",
			"NOT_EQUAL",
			"IS_EMPTY",
			"NOT_BLANK",
			"IS_BLANK",
			"STARTS_WITH",
			"ENDS_WITH",
			"NOT_STARTS_WITH",
			"NOT_ENDS_WITH",
		}, false),
}

var ApplicationAlertConfigsTagFilter = &schema.Schema{
	Type:     schema.TypeSet,
	Optional: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			TagFilterName: {
				Type:     schema.TypeString,
				Required: true,
			},
			TagFilterEntity: {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"NOT_APPLICABLE", "DESTINATION", "SOURCE"}, false),
			},
			TagFilterStringValue: {
				Type:     schema.TypeString,
				Optional: true,
			},
			TagFilterNumberValue: {
				Type:     schema.TypeFloat,
				Optional: true,
				Required: false,
			},
			TagFilterBooleanValue: {
				Type:     schema.TypeBool,
				Required: false,
				Optional: true,
			},
			TagFilterOperator: ApplicationFilterOperator,
		},
	},
}

//NewAlertingConfigResourceHandle creates the resource handle for Application Alerting Configuration
func NewApplicationAlertConfigsResourceHandle() *ResourceHandle {
	return &ResourceHandle{
		ResourceName: ResourceApplicationAlertConfigs,
		Schema: map[string]*schema.Schema{
			ApplicationAlertConfigsFieldAlertName:         ApplicationAlertConfigSchemaAlertName,
			ApplicationAlertConfigsFieldFullAlertName:     ApplicationAlertConfigSchemaFullAlertName,
			ApplicationAlertConfigsFieldDescription:       ApplicationAlertConfigSchemaDescription,
			ApplicationAlertConfigsFieldApplicationId:     ApplicationAlertConfigApplicationId,
			ApplicationAlertConfigsFieldRuleAlertType:     ApplicationAlertConfigsRuleAlertType,
			ApplicationAlertConfigsFieldRuleMetricName:    ApplicationAlertConfigsRuleMetricName,
			ApplicationAlertConfigsFieldRuleAggregation:   ApplicationAlertConfigsRuleAggregation,
			ApplicationAlertConfigsFieldRuleMessage:       ApplicationAlertConfigsRuleMessage,
			ApplicationAlertConfigsFieldRuleOperator:      ApplicationFilterOperator,
			ApplicationAlertConfigsFieldRuleLevel:         ApplicationAlertConfigsRuleLevel,
			ApplicationAlertConfigsFieldAlertChannelIds:   ApplicationAlertConfigsIntergrationIds,
			ApplicationAlertConfigsFieldThresholdValue:    ApplicationAlertConfigsThresholdValue,
			ApplicationAlertConfigsFieldThresholdType:     ApplicationAlertConfigsThresholdType,
			ApplicationAlertConfigsFieldThresholdOperator: ApplicationAlertConfigsThresholdOperator,
			ApplicationAlertConfigsFieldTagFilter:         ApplicationAlertConfigsTagFilter,
		},
		SchemaVersion: 1,
		//StateUpgraders: // None at the momen

		RestResourceFactory: func(api restapi.InstanaAPI) restapi.RestResource {
			return api.ApplicationAlertConfigs()
		},
		UpdateState: updateStateForApplicationAlertConfigs,
		// TODO
		MapStateToDataObject: mapStateToDataObjectForApplicationAlertConfigs,
	}
}

func updateStateForApplicationAlertConfigs(d *schema.ResourceData, obj restapi.InstanaDataObject) error {
	config := obj.(restapi.ApplicationAlertConfigs)

	d.Set(ApplicationAlertConfigsFieldFullAlertName, config.AlertName)
	d.Set(ApplicationAlertConfigsFieldApplicationId, config.ApplicationId)
	d.Set(ApplicationAlertConfigsFieldRuleAlertType, config.Rule.AlertType)
	d.Set(ApplicationAlertConfigsFieldRuleAggregation, config.Rule.Aggregation)
	d.Set(ApplicationAlertConfigsFieldRuleMessage, config.Rule.Message)
	d.Set(ApplicationAlertConfigsFieldRuleOperator, config.Rule.Operator)
	d.Set(ApplicationAlertConfigsFieldRuleLevel, config.Rule.Level)
	d.Set(ApplicationAlertConfigsFieldRuleMetricName, config.Rule.MetricName)
	d.Set(ApplicationAlertConfigsFieldAlertChannelIds, config.AlertChannelIds)
	d.Set(ApplicationAlertConfigsFieldThresholdValue, config.Threshold.Value)
	d.Set(ApplicationAlertConfigsFieldThresholdType, config.Threshold.Type)
	d.Set(ApplicationAlertConfigsFieldThresholdOperator, config.Threshold.Operator)
	d.Set(ApplicationAlertConfigsFieldTagFilter, config.TagFilters)
	d.SetId(config.ID)
	return nil
}

func mapStateToDataObjectForApplicationAlertConfigs(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (restapi.InstanaDataObject, error) {
	name := computeFullApplicationAlertConfigsAlertNameString(d, formatter)
	//panic: interface conversion: interface {} is *schema.Set, not []restapi.ApplicationAlertConfigsTagFilter

	list := d.Get(ApplicationAlertConfigsFieldTagFilter).(*schema.Set).List()
	tagFilters := make([]restapi.ApplicationAlertConfigsTagFilter, len(list))
	for i := range list {
		tagFilterMap := list[i].(map[string]interface{})

		tagFilters[i] = restapi.ApplicationAlertConfigsTagFilter{
			Type:         "TAG_FILTER",
			Name:         tagFilterMap[TagFilterName].(string),
			StringValue:  tagFilterMap[TagFilterStringValue].(string),
			NumberValue:  tagFilterMap[TagFilterNumberValue].(float64),
			BooleanValue: tagFilterMap[TagFilterBooleanValue].(bool),
			Operator:     tagFilterMap[TagFilterOperator].(string),
			Entity:       tagFilterMap[TagFilterEntity].(string),
		}

		//panic: interface conversion: interface {} is map[string]interface {}, not restapi.ApplicationAlertConfigsTagFilter
	}

	return restapi.ApplicationAlertConfigs{
		ID:            d.Id(),
		AlertName:     name,
		ApplicationId: d.Get(ApplicationAlertConfigsFieldApplicationId).(string),
		Rule: restapi.ApplicationAlertConfigsRule{
			AlertType:   d.Get(ApplicationAlertConfigsFieldRuleAlertType).(string),
			Aggregation: d.Get(ApplicationAlertConfigsFieldRuleAggregation).(string),
			MetricName:  d.Get(ApplicationAlertConfigsFieldRuleMetricName).(string),
			Message:     d.Get(ApplicationAlertConfigsFieldRuleMessage).(string),
			Operator:    d.Get(ApplicationAlertConfigsFieldRuleOperator).(string),
			Level:       d.Get(ApplicationAlertConfigsFieldRuleLevel).(string),
		},
		Threshold: restapi.Threshold{
			Type:        d.Get(ApplicationAlertConfigsFieldThresholdType).(string),
			Operator:    d.Get(ApplicationAlertConfigsFieldThresholdOperator).(string),
			LastUpdated: 0,
			Value:       d.Get(ApplicationAlertConfigsFieldThresholdValue).(float64),
		},
		Description: d.Get(ApplicationAlertConfigsFieldDescription).(string),
		TagFilters:  tagFilters,

		AlertChannelIds: ReadStringSetParameterFromResource(d, ApplicationAlertConfigsFieldAlertChannelIds),
		Severity:        5,
	}, nil

}

// Prefix and Suffix Alert Name
func computeFullApplicationAlertConfigsAlertNameString(d *schema.ResourceData, formatter utils.ResourceNameFormatter) string {
	if d.HasChange(ApplicationAlertConfigsFieldAlertName) {
		return formatter.Format(d.Get(ApplicationAlertConfigsFieldAlertName).(string))
	}
	return d.Get(ApplicationAlertConfigsFieldFullAlertName).(string)
}

/*
func convertSupportedEventTypesToStringSlice() []string {
	result := make([]string, len(restapi.SupportedAlertEventTypes))
	for i, t := range restapi.SupportedAlertEventTypes {
		result[i] = string(t)
	}
	return result
}

func alertingChannelConfigSchemaV0() *schema.Resource {
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
*/
