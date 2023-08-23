package instana

import (
	"context"
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceInstanaWebsiteAlertConfig the name of the terraform-provider-instana resource to manage website alert configs
const ResourceInstanaWebsiteAlertConfig = "instana_website_alert_config"

const (
	//WebsiteAlertConfigFieldAlertChannelIDs constant value for field alerting_channel_ids of resource instana_website_alert_config
	WebsiteAlertConfigFieldAlertChannelIDs = "alert_channel_ids"
	//WebsiteAlertConfigFieldWebsiteID constant value for field websites.website_id of resource instana_website_alert_config
	WebsiteAlertConfigFieldWebsiteID = "website_id"
	//WebsiteAlertConfigFieldCustomPayloadFields constant value for field custom_payload_fields of resource instana_website_alert_config
	WebsiteAlertConfigFieldCustomPayloadFields = "custom_payload_field"
	//WebsiteAlertConfigFieldCustomPayloadFieldsKey constant value for field custom_payload_fields.key of resource instana_website_alert_config
	WebsiteAlertConfigFieldCustomPayloadFieldsKey = "key"
	//WebsiteAlertConfigFieldCustomPayloadFieldsValue constant value for field custom_payload_fields.value of resource instana_website_alert_config
	WebsiteAlertConfigFieldCustomPayloadFieldsValue = "value"
	//WebsiteAlertConfigFieldDescription constant value for field description of resource instana_website_alert_config
	WebsiteAlertConfigFieldDescription = "description"
	//WebsiteAlertConfigFieldGranularity constant value for field granularity of resource instana_website_alert_config
	WebsiteAlertConfigFieldGranularity = "granularity"
	//WebsiteAlertConfigFieldName constant value for field name of resource instana_website_alert_config
	WebsiteAlertConfigFieldName = "name"
	//WebsiteAlertConfigFieldFullName constant value for field full_name of resource instana_website_alert_config
	WebsiteAlertConfigFieldFullName = "full_name"

	//WebsiteAlertConfigFieldRule constant value for field rule of resource instana_website_alert_config
	WebsiteAlertConfigFieldRule = "rule"
	//WebsiteAlertConfigFieldRuleMetricName constant value for field rule.*.metric_name of resource instana_website_alert_config
	WebsiteAlertConfigFieldRuleMetricName = "metric_name"
	//WebsiteAlertConfigFieldRuleAggregation constant value for field rule.*.aggregation of resource instana_website_alert_config
	WebsiteAlertConfigFieldRuleAggregation = "aggregation"
	//WebsiteAlertConfigFieldRuleOperator constant value for field rule.*.operator of resource instana_website_alert_config
	WebsiteAlertConfigFieldRuleOperator = "operator"
	//WebsiteAlertConfigFieldRuleValue constant value for field rule.*.value of resource instana_website_alert_config
	WebsiteAlertConfigFieldRuleValue = "value"
	//WebsiteAlertConfigFieldRuleSlowness constant value for field rule.slowness of resource instana_website_alert_config
	WebsiteAlertConfigFieldRuleSlowness = "slowness"
	//WebsiteAlertConfigFieldRuleStatusCode constant value for field rule.status_code of resource instana_website_alert_config
	WebsiteAlertConfigFieldRuleStatusCode = "status_code"
	//WebsiteAlertConfigFieldRuleThroughput constant value for field rule.throughput of resource instana_website_alert_config
	WebsiteAlertConfigFieldRuleThroughput = "throughput"
	//WebsiteAlertConfigFieldRuleSpecificJsError constant value for field rule.specific_js_error of resource instana_website_alert_config
	WebsiteAlertConfigFieldRuleSpecificJsError = "specific_js_error"

	//WebsiteAlertConfigFieldSeverity constant value for field severity of resource instana_website_alert_config
	WebsiteAlertConfigFieldSeverity = "severity"
	//WebsiteAlertConfigFieldTagFilter constant value for field tag_filter of resource instana_website_alert_config
	WebsiteAlertConfigFieldTagFilter = "tag_filter"

	//WebsiteAlertConfigFieldTimeThreshold constant value for field time_threshold of resource instana_website_alert_config
	WebsiteAlertConfigFieldTimeThreshold = "time_threshold"
	//WebsiteAlertConfigFieldTimeThresholdTimeWindow constant value for field time_threshold.time_window of resource instana_website_alert_config
	WebsiteAlertConfigFieldTimeThresholdTimeWindow = "time_window"
	//WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequence constant value for field time_threshold.user_impact_of_violations_in_sequence of resource instana_website_alert_config
	WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequence = "user_impact_of_violations_in_sequence"
	//WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequenceImpactMeasurementMethod constant value for field time_threshold.user_impact_of_violations_in_sequence.impact_measurement_method of resource instana_website_alert_config
	WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequenceImpactMeasurementMethod = "impact_measurement_method"
	//WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequenceUserPercentage constant value for field time_threshold.user_impact_of_violations_in_sequence.user_percentage of resource instana_website_alert_config
	WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequenceUserPercentage = "user_percentage"
	//WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequenceUsers constant value for field time_threshold.user_impact_of_violations_in_sequence.users of resource instana_website_alert_config
	WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequenceUsers = "users"
	//WebsiteAlertConfigFieldTimeThresholdViolationsInPeriod constant value for field time_threshold.violations_in_period of resource instana_website_alert_config
	WebsiteAlertConfigFieldTimeThresholdViolationsInPeriod = "violations_in_period"
	//WebsiteAlertConfigFieldTimeThresholdViolationsInPeriodViolations constant value for field time_threshold.violations_in_period.violations of resource instana_website_alert_config
	WebsiteAlertConfigFieldTimeThresholdViolationsInPeriodViolations = "violations"
	//WebsiteAlertConfigFieldTimeThresholdViolationsInSequence constant value for field time_threshold.violations_in_sequence of resource instana_website_alert_config
	WebsiteAlertConfigFieldTimeThresholdViolationsInSequence = "violations_in_sequence"
	//WebsiteAlertConfigFieldTriggering constant value for field triggering of resource instana_website_alert_config
	WebsiteAlertConfigFieldTriggering = "triggering"
)

var (
	websiteAlertConfigSchemaAlertChannelIDs = &schema.Schema{
		Type:     schema.TypeSet,
		MinItems: 0,
		MaxItems: 1024,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Optional:    true,
		Description: "List of IDs of alert channels defined in Instana.",
	}
	websiteAlertConfigSchemaCustomPayloadFields = &schema.Schema{
		Type: schema.TypeSet,
		Set: func(i interface{}) int {
			return schema.HashString(i.(map[string]interface{})[WebsiteAlertConfigFieldCustomPayloadFieldsKey])
		},
		Optional: true,
		MinItems: 0,
		MaxItems: 20,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				WebsiteAlertConfigFieldCustomPayloadFieldsKey: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The key of the custom payload field",
				},
				WebsiteAlertConfigFieldCustomPayloadFieldsValue: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "The value of the custom payload field",
				},
			},
		},
		Description: "An optional list of custom payload fields (static key/value pairs added to the event)",
	}
	websiteAlertConfigSchemaDescription = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		Description:  "The description text of the website alert config",
		ValidateFunc: validation.StringLenBetween(0, 65536),
	}
	websiteAlertConfigSchemaGranularity = &schema.Schema{
		Type:         schema.TypeInt,
		Optional:     true,
		Default:      restapi.Granularity600000,
		ValidateFunc: validation.IntInSlice(restapi.SupportedGranularities.ToIntSlice()),
		Description:  "The evaluation granularity used for detection of violations of the defined threshold. In other words, it defines the size of the tumbling window used",
	}
	websiteAlertConfigSchemaName = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		Description:  "Name for the website alert configuration",
		ValidateFunc: validation.StringLenBetween(0, 256),
	}
	websiteAlertConfigSchemaFullName = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The full name field of the website alert config. The field is computed and contains the name which is sent to Instana. The computation depends on the configured default_name_prefix and default_name_suffix at provider level",
	}
	websiteAlertConfigSchemaRule = &schema.Schema{
		Type:        schema.TypeList,
		MinItems:    1,
		MaxItems:    1,
		Required:    true,
		Description: "Indicates the type of rule this alert configuration is about.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				WebsiteAlertConfigFieldRuleSlowness: {
					Type:        schema.TypeList,
					MinItems:    0,
					MaxItems:    1,
					Optional:    true,
					Description: "Rule based on the slowness of the configured alert configuration target",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							WebsiteAlertConfigFieldRuleMetricName:  websiteAlertConfigSchemaRuleMetricName,
							WebsiteAlertConfigFieldRuleAggregation: websiteAlertConfigSchemaRequiredRuleAggregation,
						},
					},
					ExactlyOneOf: websiteAlertConfigRuleTypeKeys,
				},
				WebsiteAlertConfigFieldRuleSpecificJsError: {
					Type:        schema.TypeList,
					MinItems:    0,
					MaxItems:    1,
					Optional:    true,
					Description: "Rule based on a specific javascript error of the configured alert configuration target",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							WebsiteAlertConfigFieldRuleMetricName:  websiteAlertConfigSchemaRuleMetricName,
							WebsiteAlertConfigFieldRuleAggregation: websiteAlertConfigSchemaOptionalRuleAggregation,
							WebsiteAlertConfigFieldRuleOperator:    websiteAlertConfigSchemaRuleOperator,
							WebsiteAlertConfigFieldRuleValue: {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "The value identify the specific javascript error",
							},
						},
					},
					ExactlyOneOf: websiteAlertConfigRuleTypeKeys,
				},
				WebsiteAlertConfigFieldRuleStatusCode: {
					Type:        schema.TypeList,
					MinItems:    0,
					MaxItems:    1,
					Optional:    true,
					Description: "Rule based on the HTTP status code of the configured alert configuration target",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							WebsiteAlertConfigFieldRuleMetricName:  websiteAlertConfigSchemaRuleMetricName,
							WebsiteAlertConfigFieldRuleAggregation: websiteAlertConfigSchemaOptionalRuleAggregation,
							WebsiteAlertConfigFieldRuleOperator:    websiteAlertConfigSchemaRuleOperator,
							WebsiteAlertConfigFieldRuleValue: {
								Type:        schema.TypeString,
								Required:    true,
								Description: "The value identify the specific http status code",
							},
						},
					},
					ExactlyOneOf: websiteAlertConfigRuleTypeKeys,
				},
				WebsiteAlertConfigFieldRuleThroughput: {
					Type:        schema.TypeList,
					MinItems:    0,
					MaxItems:    1,
					Optional:    true,
					Description: "Rule based on the throughput of the configured alert configuration target",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							WebsiteAlertConfigFieldRuleMetricName:  websiteAlertConfigSchemaRuleMetricName,
							WebsiteAlertConfigFieldRuleAggregation: websiteAlertConfigSchemaOptionalRuleAggregation,
						},
					},
					ExactlyOneOf: websiteAlertConfigRuleTypeKeys,
				},
			},
		},
	}
	websiteAlertConfigSchemaSeverity = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice(restapi.SupportedSeverities.TerraformRepresentations(), false),
		Description:  "The severity of the alert when triggered",
	}
	websiteAlertConfigSchemaTagFilter = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The tag filter of the website alert config",
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			normalized, err := tagfilter.Normalize(new)
			if err == nil {
				return normalized == old
			}
			return old == new
		},
		StateFunc: func(val interface{}) string {
			normalized, err := tagfilter.Normalize(val.(string))
			if err == nil {
				return normalized
			}
			return val.(string)
		},
		ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
			v := val.(string)
			if _, err := tagfilter.NewParser().Parse(v); err != nil {
				errs = append(errs, fmt.Errorf("%q is not a valid tag filter; %s", key, err))
			}

			return
		},
	}
	websiteAlertConfigSchemaTimeThreshold = &schema.Schema{
		Type:        schema.TypeList,
		MinItems:    1,
		MaxItems:    1,
		Required:    true,
		Description: "Indicates the type of violation of the defined threshold.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequence: {
					Type:        schema.TypeList,
					MinItems:    0,
					MaxItems:    1,
					Optional:    true,
					Description: "Time threshold base on user impact of violations in sequence",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							WebsiteAlertConfigFieldTimeThresholdTimeWindow: websiteAlertConfigSchemaOptionalTimeThresholdTimeWindow,
							WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequenceImpactMeasurementMethod: {
								Type:         schema.TypeString,
								Required:     true,
								ValidateFunc: validation.StringInSlice(restapi.SupportedWebsiteImpactMeasurementMethods.ToStringSlice(), false),
								Description:  "The impact method of the time threshold based on user impact of violations in sequence",
							},
							WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequenceUserPercentage: {
								Type:         schema.TypeFloat,
								Optional:     true,
								ValidateFunc: validation.FloatBetween(0.0, 1.0),
								Description:  "The percentage of impacted users of the time threshold based on user impact of violations in sequence",
							},
							WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequenceUsers: {
								Type:         schema.TypeInt,
								Optional:     true,
								ValidateFunc: validation.IntAtLeast(1),
								Description:  "The number of impacted users of the time threshold based on user impact of violations in sequence",
							},
						},
					},
					ExactlyOneOf: websiteAlertConfigTimeThresholdTypeKeys,
				},
				WebsiteAlertConfigFieldTimeThresholdViolationsInPeriod: {
					Type:        schema.TypeList,
					MinItems:    0,
					MaxItems:    1,
					Optional:    true,
					Description: "Time threshold base on violations in period",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							WebsiteAlertConfigFieldTimeThresholdTimeWindow: websiteAlertConfigSchemaOptionalTimeThresholdTimeWindow,
							WebsiteAlertConfigFieldTimeThresholdViolationsInPeriodViolations: {
								Type:         schema.TypeInt,
								Optional:     true,
								ValidateFunc: validation.IntBetween(1, 12),
								Description:  "The violations appeared in the period",
							},
						},
					},
					ExactlyOneOf: websiteAlertConfigTimeThresholdTypeKeys,
				},
				WebsiteAlertConfigFieldTimeThresholdViolationsInSequence: {
					Type:        schema.TypeList,
					MinItems:    0,
					MaxItems:    1,
					Optional:    true,
					Description: "Time threshold base on violations in sequence",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							WebsiteAlertConfigFieldTimeThresholdTimeWindow: websiteAlertConfigSchemaOptionalTimeThresholdTimeWindow,
						},
					},
					ExactlyOneOf: websiteAlertConfigTimeThresholdTypeKeys,
				},
			},
		},
	}
	websiteAlertConfigSchemaTriggering = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Optional flag to indicate whether also an Incident is triggered or not. The default is false",
	}
	websiteAlertConfigSchemaWebsiteID = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		Description:  "Unique ID of the website",
		ValidateFunc: validation.StringLenBetween(0, 64),
	}
	websiteAlertConfigRuleTypeKeys = []string{
		"rule.0.specific_js_error",
		"rule.0.slowness",
		"rule.0.status_code",
		"rule.0.throughput",
	}
	websiteAlertConfigSchemaRuleMetricName = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The metric name of the website alert rule",
	}
	websiteAlertConfigSchemaRequiredRuleAggregation = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice(restapi.SupportedAggregations.ToStringSlice(), true),
		Description:  "The aggregation function of the website alert rule",
	}
	websiteAlertConfigSchemaOptionalRuleAggregation = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		ValidateFunc: validation.StringInSlice(restapi.SupportedAggregations.ToStringSlice(), true),
		Description:  "The aggregation function of the website alert rule",
	}
	websiteAlertConfigSchemaRuleOperator = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		Description:  "The operator which will be applied to evaluate this rule",
		ValidateFunc: validation.StringInSlice(restapi.SupportedExpressionOperators.ToStringSlice(), true),
	}
	websiteAlertConfigTimeThresholdTypeKeys = []string{
		"time_threshold.0.user_impact_of_violations_in_sequence",
		"time_threshold.0.violations_in_period",
		"time_threshold.0.violations_in_sequence",
	}
	websiteAlertConfigSchemaOptionalTimeThresholdTimeWindow = &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Description: "The time window if the time threshold",
	}
)

var websiteAlertConfigResourceSchema = map[string]*schema.Schema{
	WebsiteAlertConfigFieldAlertChannelIDs:     websiteAlertConfigSchemaAlertChannelIDs,
	WebsiteAlertConfigFieldCustomPayloadFields: websiteAlertConfigSchemaCustomPayloadFields,
	WebsiteAlertConfigFieldDescription:         websiteAlertConfigSchemaDescription,
	WebsiteAlertConfigFieldGranularity:         websiteAlertConfigSchemaGranularity,
	WebsiteAlertConfigFieldName:                websiteAlertConfigSchemaName,
	WebsiteAlertConfigFieldRule:                websiteAlertConfigSchemaRule,
	WebsiteAlertConfigFieldSeverity:            websiteAlertConfigSchemaSeverity,
	WebsiteAlertConfigFieldTagFilter:           websiteAlertConfigSchemaTagFilter,
	ResourceFieldThreshold:                     thresholdSchema,
	WebsiteAlertConfigFieldTimeThreshold:       websiteAlertConfigSchemaTimeThreshold,
	WebsiteAlertConfigFieldTriggering:          websiteAlertConfigSchemaTriggering,
	WebsiteAlertConfigFieldWebsiteID:           websiteAlertConfigSchemaWebsiteID,
}

// NewWebsiteAlertConfigResourceHandle creates the resource handle for Website Alert Configs
func NewWebsiteAlertConfigResourceHandle() ResourceHandle[*restapi.WebsiteAlertConfig] {
	return &websiteAlertConfigResource{
		metaData: ResourceMetaData{
			ResourceName:     ResourceInstanaWebsiteAlertConfig,
			Schema:           websiteAlertConfigResourceSchema,
			SkipIDGeneration: true,
			SchemaVersion:    1,
		},
	}
}

type websiteAlertConfigResource struct {
	metaData ResourceMetaData
}

func (r *websiteAlertConfigResource) MetaData() *ResourceMetaData {
	return &r.metaData
}

func (r *websiteAlertConfigResource) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{
		{
			Type:    r.websiteAlertConfigSchemaV0().CoreConfigSchema().ImpliedType(),
			Upgrade: r.websiteAlertConfigStateUpgradeV0,
			Version: 0,
		},
	}
}

func (r *websiteAlertConfigResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.WebsiteAlertConfig] {
	return api.WebsiteAlertConfig()
}

func (r *websiteAlertConfigResource) SetComputedFields(_ *schema.ResourceData) error {
	return nil
}

func (r *websiteAlertConfigResource) UpdateState(d *schema.ResourceData, config *restapi.WebsiteAlertConfig) error {
	severity, err := ConvertSeverityFromInstanaAPIToTerraformRepresentation(config.Severity)
	if err != nil {
		return err
	}
	var normalizedTagFilterString *string
	if config.TagFilterExpression != nil {
		normalizedTagFilterString, err = tagfilter.MapTagFilterToNormalizedString(config.TagFilterExpression.(restapi.TagFilterExpressionElement))
		if err != nil {
			return err
		}
	}

	d.SetId(config.ID)
	return tfutils.UpdateState(d, map[string]interface{}{
		WebsiteAlertConfigFieldAlertChannelIDs:     config.AlertChannelIDs,
		WebsiteAlertConfigFieldCustomPayloadFields: r.mapCustomPayloadFieldsToSchema(config),
		WebsiteAlertConfigFieldDescription:         config.Description,
		WebsiteAlertConfigFieldGranularity:         config.Granularity,
		WebsiteAlertConfigFieldName:                config.Name,
		WebsiteAlertConfigFieldRule:                r.mapRuleToSchema(config),
		WebsiteAlertConfigFieldSeverity:            severity,
		WebsiteAlertConfigFieldTagFilter:           normalizedTagFilterString,
		ResourceFieldThreshold:                     newThresholdMapper().toState(&config.Threshold),
		WebsiteAlertConfigFieldTimeThreshold:       r.mapTimeThresholdToSchema(config),
		WebsiteAlertConfigFieldTriggering:          config.Triggering,
		WebsiteAlertConfigFieldWebsiteID:           config.WebsiteID,
	})
}

func (r *websiteAlertConfigResource) mapCustomPayloadFieldsToSchema(config *restapi.WebsiteAlertConfig) []map[string]string {
	result := make([]map[string]string, len(config.CustomerPayloadFields))
	for i, v := range config.CustomerPayloadFields {
		field := make(map[string]string)
		field[WebsiteAlertConfigFieldCustomPayloadFieldsKey] = v.Key
		field[WebsiteAlertConfigFieldCustomPayloadFieldsValue] = string(v.Value)
		result[i] = field
	}
	return result
}

func (r *websiteAlertConfigResource) mapRuleToSchema(config *restapi.WebsiteAlertConfig) []map[string]interface{} {
	ruleAttribute := make(map[string]interface{})
	ruleAttribute[WebsiteAlertConfigFieldRuleMetricName] = config.Rule.MetricName

	if config.Rule.Aggregation != nil {
		ruleAttribute[WebsiteAlertConfigFieldRuleAggregation] = string(*config.Rule.Aggregation)
	}
	if config.Rule.Operator != nil {
		ruleAttribute[WebsiteAlertConfigFieldRuleOperator] = string(*config.Rule.Operator)
	}
	if config.Rule.Value != nil {
		ruleAttribute[WebsiteAlertConfigFieldRuleValue] = *config.Rule.Value
	}

	alertType := r.mapAlertTypeToSchema(config.Rule.AlertType)
	rule := make(map[string]interface{})
	rule[alertType] = []interface{}{ruleAttribute}
	result := make([]map[string]interface{}, 1)
	result[0] = rule
	return result
}

func (r *websiteAlertConfigResource) mapAlertTypeToSchema(alertType string) string {
	if alertType == "specificJsError" {
		return WebsiteAlertConfigFieldRuleSpecificJsError
	} else if alertType == "statusCode" {
		return WebsiteAlertConfigFieldRuleStatusCode
	}
	return alertType
}

func (r *websiteAlertConfigResource) mapTimeThresholdToSchema(config *restapi.WebsiteAlertConfig) []map[string]interface{} {
	timeThresholdConfig := make(map[string]interface{})

	if config.TimeThreshold.TimeWindow != nil {
		timeThresholdConfig[WebsiteAlertConfigFieldTimeThresholdTimeWindow] = config.TimeThreshold.TimeWindow
	}
	if config.TimeThreshold.Violations != nil {
		timeThresholdConfig[WebsiteAlertConfigFieldTimeThresholdViolationsInPeriodViolations] = int(*config.TimeThreshold.Violations)
	}
	if config.TimeThreshold.ImpactMeasurementMethod != nil {
		timeThresholdConfig[WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequenceImpactMeasurementMethod] = string(*config.TimeThreshold.ImpactMeasurementMethod)
	}
	if config.TimeThreshold.Users != nil {
		timeThresholdConfig[WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequenceUsers] = int(*config.TimeThreshold.Users)
	}
	if config.TimeThreshold.UserPercentage != nil {
		timeThresholdConfig[WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequenceUserPercentage] = *config.TimeThreshold.UserPercentage
	}

	timeThresholdType := r.mapTimeThresholdTypeToSchema(config.TimeThreshold.Type)
	timeThreshold := make(map[string]interface{})
	timeThreshold[timeThresholdType] = []interface{}{timeThresholdConfig}
	result := make([]map[string]interface{}, 1)
	result[0] = timeThreshold
	return result
}

func (r *websiteAlertConfigResource) mapTimeThresholdTypeToSchema(input string) string {
	if input == "userImpactOfViolationsInSequence" {
		return WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequence
	} else if input == "violationsInPeriod" {
		return WebsiteAlertConfigFieldTimeThresholdViolationsInPeriod
	} else if input == "violationsInSequence" {
		return WebsiteAlertConfigFieldTimeThresholdViolationsInSequence
	}
	return input
}

func (r *websiteAlertConfigResource) MapStateToDataObject(d *schema.ResourceData) (*restapi.WebsiteAlertConfig, error) {
	severity, err := ConvertSeverityFromTerraformToInstanaAPIRepresentation(d.Get(WebsiteAlertConfigFieldSeverity).(string))
	if err != nil {
		return nil, err
	}

	var tagFilter restapi.TagFilterExpressionElement
	tagFilterStr, ok := d.GetOk(WebsiteAlertConfigFieldTagFilter)
	if ok {
		tagFilter, err = r.mapTagFilterExpressionFromSchema(tagFilterStr.(string))
		if err != nil {
			return &restapi.WebsiteAlertConfig{}, err
		}
	}

	threshold := newThresholdMapper().fromState(d)

	return &restapi.WebsiteAlertConfig{
		ID:                    d.Id(),
		AlertChannelIDs:       ReadStringSetParameterFromResource(d, WebsiteAlertConfigFieldAlertChannelIDs),
		CustomerPayloadFields: r.mapCustomPayloadFieldsFromSchema(d),
		Description:           d.Get(WebsiteAlertConfigFieldDescription).(string),
		Granularity:           restapi.Granularity(d.Get(WebsiteAlertConfigFieldGranularity).(int)),
		Name:                  d.Get(WebsiteMonitoringConfigFieldName).(string),
		Rule:                  *r.mapRuleFromSchema(d),
		Severity:              severity,
		TagFilterExpression:   tagFilter,
		Threshold:             *threshold,
		TimeThreshold:         *r.mapTimeThresholdFromSchema(d),
		Triggering:            d.Get(WebsiteAlertConfigFieldTriggering).(bool),
		WebsiteID:             d.Get(WebsiteAlertConfigFieldWebsiteID).(string),
	}, nil
}

func (r *websiteAlertConfigResource) mapTagFilterExpressionFromSchema(input string) (restapi.TagFilterExpressionElement, error) {
	parser := tagfilter.NewParser()
	expr, err := parser.Parse(input)
	if err != nil {
		return nil, err
	}

	mapper := tagfilter.NewMapper()
	return mapper.ToAPIModel(expr), nil
}

func (r *websiteAlertConfigResource) mapCustomPayloadFieldsFromSchema(d *schema.ResourceData) []restapi.CustomPayloadField[restapi.StaticStringCustomPayloadFieldValue] {
	val := d.Get(WebsiteAlertConfigFieldCustomPayloadFields)
	if val != nil {
		fields := val.(*schema.Set).List()
		result := make([]restapi.CustomPayloadField[restapi.StaticStringCustomPayloadFieldValue], len(fields))
		for i, v := range fields {
			field := v.(map[string]interface{})
			customPayloadFieldType := restapi.StaticCustomPayloadType
			key := field[WebsiteAlertConfigFieldCustomPayloadFieldsKey].(string)
			value := field[WebsiteAlertConfigFieldCustomPayloadFieldsValue].(string)
			result[i] = restapi.CustomPayloadField[restapi.StaticStringCustomPayloadFieldValue]{
				Type:  customPayloadFieldType,
				Key:   key,
				Value: restapi.StaticStringCustomPayloadFieldValue(value),
			}
		}
		return result
	}
	return []restapi.CustomPayloadField[restapi.StaticStringCustomPayloadFieldValue]{}
}

func (r *websiteAlertConfigResource) mapRuleFromSchema(d *schema.ResourceData) *restapi.WebsiteAlertRule {
	ruleSlice := d.Get(WebsiteAlertConfigFieldRule).([]interface{})
	rule := ruleSlice[0].(map[string]interface{})
	for alertType, v := range rule {
		configSlice := v.([]interface{})
		if len(configSlice) == 1 {
			config := configSlice[0].(map[string]interface{})
			return r.mapRuleConfigFromSchema(config, alertType)
		}
	}
	return &restapi.WebsiteAlertRule{}
}

func (r *websiteAlertConfigResource) mapRuleConfigFromSchema(config map[string]interface{}, alertType string) *restapi.WebsiteAlertRule {
	var aggregationPtr *restapi.Aggregation
	if aggregationStr, ok := config[WebsiteAlertConfigFieldRuleAggregation]; ok {
		aggregation := restapi.Aggregation(aggregationStr.(string))
		aggregationPtr = &aggregation
	}
	var valuePtr *string
	if v, ok := config[WebsiteAlertConfigFieldRuleValue]; ok {
		value := v.(string)
		valuePtr = &value
	}
	var operatorPtr *restapi.ExpressionOperator
	if v, ok := config[WebsiteAlertConfigFieldRuleOperator]; ok {
		operator := restapi.ExpressionOperator(v.(string))
		operatorPtr = &operator
	}
	return &restapi.WebsiteAlertRule{
		AlertType:   r.mapAlertTypeFromSchema(alertType),
		MetricName:  config[WebsiteAlertConfigFieldRuleMetricName].(string),
		Aggregation: aggregationPtr,
		Operator:    operatorPtr,
		Value:       valuePtr,
	}
}

func (r *websiteAlertConfigResource) mapAlertTypeFromSchema(alertType string) string {
	if alertType == WebsiteAlertConfigFieldRuleSpecificJsError {
		return "specificJsError"
	} else if alertType == WebsiteAlertConfigFieldRuleStatusCode {
		return "statusCode"
	}
	return alertType
}

func (r *websiteAlertConfigResource) mapTimeThresholdFromSchema(d *schema.ResourceData) *restapi.WebsiteTimeThreshold {
	timeThresholdSlice := d.Get(WebsiteAlertConfigFieldTimeThreshold).([]interface{})
	timeThreshold := timeThresholdSlice[0].(map[string]interface{})
	for timeThresholdType, v := range timeThreshold {
		configSlice := v.([]interface{})
		if len(configSlice) == 1 {
			config := configSlice[0].(map[string]interface{})
			var timeWindowPtr *int64
			if v, ok := config[WebsiteAlertConfigFieldTimeThresholdTimeWindow]; ok {
				timeWindow := int64(v.(int))
				timeWindowPtr = &timeWindow
			}
			var violationsPtr *int32
			if v, ok := config[WebsiteAlertConfigFieldTimeThresholdViolationsInPeriodViolations]; ok {
				violations := int32(v.(int))
				violationsPtr = &violations
			}
			var impactMeasurementMethodPtr *restapi.WebsiteImpactMeasurementMethod
			if v, ok := config[WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequenceImpactMeasurementMethod]; ok {
				impactMeasurementMethod := restapi.WebsiteImpactMeasurementMethod(v.(string))
				impactMeasurementMethodPtr = &impactMeasurementMethod
			}
			var userPercentagePtr *float64
			if v, ok := config[WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequenceUserPercentage]; ok {
				userPercentage := v.(float64)
				userPercentagePtr = &userPercentage
			}
			var usersPtr *int32
			if v, ok := config[WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequenceUsers]; ok {
				users := int32(v.(int))
				usersPtr = &users
			}
			return &restapi.WebsiteTimeThreshold{
				Type:                    r.mapTimeThresholdTypeFromSchema(timeThresholdType),
				TimeWindow:              timeWindowPtr,
				Violations:              violationsPtr,
				ImpactMeasurementMethod: impactMeasurementMethodPtr,
				UserPercentage:          userPercentagePtr,
				Users:                   usersPtr,
			}
		}
	}
	return &restapi.WebsiteTimeThreshold{}
}

func (r *websiteAlertConfigResource) mapTimeThresholdTypeFromSchema(input string) string {
	if input == WebsiteAlertConfigFieldTimeThresholdUserImpactOfViolationsInSequence {
		return "userImpactOfViolationsInSequence"
	} else if input == WebsiteAlertConfigFieldTimeThresholdViolationsInPeriod {
		return "violationsInPeriod"
	} else if input == WebsiteAlertConfigFieldTimeThresholdViolationsInSequence {
		return "violationsInSequence"
	}
	return input
}

func (r *websiteAlertConfigResource) websiteAlertConfigStateUpgradeV0(_ context.Context, state map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	if _, ok := state[WebsiteAlertConfigFieldFullName]; ok {
		state[WebsiteAlertConfigFieldName] = state[WebsiteAlertConfigFieldFullName]
		delete(state, WebsiteAlertConfigFieldFullName)
	}
	return state, nil
}

func (r *websiteAlertConfigResource) websiteAlertConfigSchemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			WebsiteAlertConfigFieldAlertChannelIDs:     websiteAlertConfigSchemaAlertChannelIDs,
			WebsiteAlertConfigFieldCustomPayloadFields: websiteAlertConfigSchemaCustomPayloadFields,
			WebsiteAlertConfigFieldDescription:         websiteAlertConfigSchemaDescription,
			WebsiteAlertConfigFieldGranularity:         websiteAlertConfigSchemaGranularity,
			WebsiteAlertConfigFieldName:                websiteAlertConfigSchemaName,
			WebsiteAlertConfigFieldFullName:            websiteAlertConfigSchemaFullName,
			WebsiteAlertConfigFieldRule:                websiteAlertConfigSchemaRule,
			WebsiteAlertConfigFieldSeverity:            websiteAlertConfigSchemaSeverity,
			WebsiteAlertConfigFieldTagFilter:           websiteAlertConfigSchemaTagFilter,
			ResourceFieldThreshold:                     thresholdSchema,
			WebsiteAlertConfigFieldTimeThreshold:       websiteAlertConfigSchemaTimeThreshold,
			WebsiteAlertConfigFieldTriggering:          websiteAlertConfigSchemaTriggering,
			WebsiteAlertConfigFieldWebsiteID:           websiteAlertConfigSchemaWebsiteID,
		},
	}
}
