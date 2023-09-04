package instana

import (
	"errors"
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceInstanaCustomEventSpecification the name of the terraform-provider-instana resource to manage custom event specifications
const ResourceInstanaCustomEventSpecification = "instana_custom_event_specification"

const (
	CustomEventSpecificationFieldRule                   = "rule"
	CustomEventSpecificationFieldEntityVerificationRule = "entity_verification"
	CustomEventSpecificationFieldSystemRule             = "system"
	CustomEventSpecificationFieldThresholdRule          = "threshold"
	customEventSpecificationThresholdRulePath           = "rule.0.threshold.0."

	CustomEventSpecificationRuleFieldSeverity                              = "severity"
	CustomEventSpecificationEntityVerificationRuleFieldMatchingEntityType  = "matching_entity_type"
	CustomEventSpecificationEntityVerificationRuleFieldMatchingOperator    = "matching_operator"
	CustomEventSpecificationEntityVerificationRuleFieldMatchingEntityLabel = "matching_entity_label"
	CustomEventSpecificationEntityVerificationRuleFieldOfflineDuration     = "offline_duration"
	CustomEventSpecificationSystemRuleFieldSystemRuleId                    = "system_rule_id"
	CustomEventSpecificationThresholdRuleFieldMetricName                   = "metric_name"
	CustomEventSpecificationThresholdRuleFieldRollup                       = "rollup"
	CustomEventSpecificationThresholdRuleFieldWindow                       = "window"
	CustomEventSpecificationThresholdRuleFieldAggregation                  = "aggregation"
	CustomEventSpecificationThresholdRuleFieldConditionOperator            = "condition_operator"
	CustomEventSpecificationThresholdRuleFieldConditionValue               = "condition_value"
	CustomEventSpecificationThresholdRuleFieldMetricPattern                = "metric_pattern"
	CustomEventSpecificationThresholdRuleFieldMetricPatternPrefix          = "prefix"
	CustomEventSpecificationThresholdRuleFieldMetricPatternPostfix         = "postfix"
	CustomEventSpecificationThresholdRuleFieldMetricPatternPlaceholder     = "placeholder"
	CustomEventSpecificationThresholdRuleFieldMetricPatternOperator        = "operator"
)

var (
	customEventSpecificationRuleTypeKeys = []string{
		"rule.0.entity_verification",
		"rule.0.system",
		"rule.0.threshold",
	}

	customEventSpecificationThresholdRuleMetricNameOrPattern = []string{
		customEventSpecificationThresholdRulePath + CustomEventSpecificationThresholdRuleFieldMetricName,
		customEventSpecificationThresholdRulePath + CustomEventSpecificationThresholdRuleFieldMetricPattern,
	}

	customEventSpecificationThresholdRuleWindowOrRollup = []string{
		customEventSpecificationThresholdRulePath + CustomEventSpecificationThresholdRuleFieldRollup,
		customEventSpecificationThresholdRulePath + CustomEventSpecificationThresholdRuleFieldWindow,
	}
)

// NewCustomEventSpecificationResourceHandle creates a new ResourceHandle for the terraform resource of custom event specifications
func NewCustomEventSpecificationResourceHandle() ResourceHandle {
	commons := &customEventSpecificationCommons{}
	return &customEventSpecificationResource{
		metaData: ResourceMetaData{
			ResourceName: ResourceInstanaCustomEventSpecification,
			Schema: map[string]*schema.Schema{
				CustomEventSpecificationFieldName:           customEventSpecificationSchemaName,
				CustomEventSpecificationFieldQuery:          customEventSpecificationSchemaQuery,
				CustomEventSpecificationFieldTriggering:     customEventSpecificationSchemaTriggering,
				CustomEventSpecificationFieldDescription:    customEventSpecificationSchemaDescription,
				CustomEventSpecificationFieldExpirationTime: customEventSpecificationSchemaExpirationTime,
				CustomEventSpecificationFieldEnabled:        customEventSpecificationSchemaEnabled,
				CustomEventSpecificationFieldRule: {
					Type:        schema.TypeList,
					MinItems:    1,
					MaxItems:    1,
					Required:    true,
					Description: "Indicates the type of rule this custom event specification is about.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							CustomEventSpecificationFieldEntityVerificationRule: {
								Type:        schema.TypeList,
								MinItems:    0,
								MaxItems:    1,
								Optional:    true,
								Description: "Entity verification rule configuration",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										CustomEventSpecificationRuleFieldSeverity: customEventSpecificationSchemaRuleSeverity,
										CustomEventSpecificationEntityVerificationRuleFieldMatchingEntityType: {
											Type:        schema.TypeString,
											Required:    true,
											Description: "The type of the matching entity",
										},
										CustomEventSpecificationEntityVerificationRuleFieldMatchingOperator: {
											Type:     schema.TypeString,
											Required: true,
											ValidateFunc: validation.StringInSlice([]string{
												"is",
												"contains",
												"startsWith",
												"endsWith"}, false),
											Description: "The operator which should be applied for matching the label for the given entity (e.g. is, contains, startsWith, endsWith)",
										},
										CustomEventSpecificationEntityVerificationRuleFieldMatchingEntityLabel: {
											Type:        schema.TypeString,
											Required:    true,
											Description: "The label of the matching entity",
										},
										CustomEventSpecificationEntityVerificationRuleFieldOfflineDuration: {
											Type:        schema.TypeInt,
											Required:    true,
											Description: "The duration after which the matching entity is considered to be offline",
										},
									},
								},
								ExactlyOneOf: customEventSpecificationRuleTypeKeys,
							},
							CustomEventSpecificationFieldSystemRule: {
								Type:        schema.TypeList,
								MinItems:    0,
								MaxItems:    1,
								Optional:    true,
								Description: "System rule configuration",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										CustomEventSpecificationRuleFieldSeverity: customEventSpecificationSchemaRuleSeverity,
										CustomEventSpecificationSystemRuleFieldSystemRuleId: {
											Type:        schema.TypeString,
											Required:    true,
											Description: "Configures the system rule id for the system rule of the custom event specification",
										},
									},
								},
								ExactlyOneOf: customEventSpecificationRuleTypeKeys,
							},
							CustomEventSpecificationFieldThresholdRule: {
								Type:        schema.TypeList,
								MinItems:    0,
								MaxItems:    1,
								Optional:    true,
								Description: "System rule configuration",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										CustomEventSpecificationRuleFieldSeverity: customEventSpecificationSchemaRuleSeverity,

										CustomEventSpecificationThresholdRuleFieldMetricName: {
											Type:         schema.TypeString,
											Required:     false,
											Optional:     true,
											Description:  "The metric name of the rule",
											ExactlyOneOf: customEventSpecificationThresholdRuleMetricNameOrPattern,
										},

										CustomEventSpecificationThresholdRuleFieldMetricPattern: {
											Type:        schema.TypeList,
											MinItems:    0,
											MaxItems:    1,
											Optional:    true,
											Description: "The metric pattern of the rule",
											Elem: &schema.Resource{
												Schema: map[string]*schema.Schema{
													CustomEventSpecificationThresholdRuleFieldMetricPatternPrefix: {
														Type:        schema.TypeString,
														Required:    false,
														Optional:    true,
														Description: "The metric pattern prefix of a dynamic built-in metrics",
													},
													CustomEventSpecificationThresholdRuleFieldMetricPatternPostfix: {
														Type:        schema.TypeString,
														Required:    false,
														Optional:    true,
														Description: "The metric pattern postfix of a dynamic built-in metrics",
													},
													CustomEventSpecificationThresholdRuleFieldMetricPatternPlaceholder: {
														Type:        schema.TypeString,
														Required:    false,
														Optional:    true,
														Description: "The metric pattern placeholer/condition value of a dynamic built-in metrics",
													},
													CustomEventSpecificationThresholdRuleFieldMetricPatternOperator: {
														Type:         schema.TypeString,
														Required:     false,
														Optional:     true,
														RequiredWith: []string{customEventSpecificationThresholdRulePath + CustomEventSpecificationThresholdRuleFieldMetricPatternPrefix},
														ValidateFunc: validation.StringInSlice(restapi.SupportedMetricPatternOperatorTypes.ToStringSlice(), false),
														Description:  "The metric pattern operator (e.g >, <)",
													},
												},
											},
											ExactlyOneOf: customEventSpecificationThresholdRuleMetricNameOrPattern,
										},

										CustomEventSpecificationThresholdRuleFieldRollup: {
											Type:         schema.TypeInt,
											Required:     false,
											Optional:     true,
											Description:  "The rollup of the metric",
											ExactlyOneOf: customEventSpecificationThresholdRuleWindowOrRollup,
										},
										CustomEventSpecificationThresholdRuleFieldWindow: {
											Type:         schema.TypeInt,
											Required:     false,
											Optional:     true,
											Description:  "The time window where the condition has to be fulfilled",
											ExactlyOneOf: customEventSpecificationThresholdRuleWindowOrRollup,
										},
										CustomEventSpecificationThresholdRuleFieldAggregation: {
											Type:         schema.TypeString,
											Required:     false,
											Optional:     true,
											ValidateFunc: validation.StringInSlice(restapi.SupportedAggregationTypes.ToStringSlice(), false),
											Description:  "The aggregation type (e.g. sum, avg)",
										},
										CustomEventSpecificationThresholdRuleFieldConditionOperator: {
											Type:     schema.TypeString,
											Required: true,
											ValidateFunc: validation.StringInSlice([]string{
												">",
												">=",
												"<",
												"<=",
												"=",
												"!=",
											}, false),
											Description: "The condition operator (e.g >, <)",
										},
										CustomEventSpecificationThresholdRuleFieldConditionValue: {
											Type:        schema.TypeFloat,
											Required:    true,
											Description: "The expected condition value to fulfill the rule",
										},
									},
								},
								ExactlyOneOf: customEventSpecificationRuleTypeKeys,
							},
						},
					},
				},
			},
			SchemaVersion: 3,
		},
		commons: commons,
	}
}

type customEventSpecificationResource struct {
	metaData ResourceMetaData
	commons  *customEventSpecificationCommons
}

func (r *customEventSpecificationResource) MetaData() *ResourceMetaData {
	return &r.metaData
}

func (r *customEventSpecificationResource) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{}
}

func (r *customEventSpecificationResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource {
	return api.CustomEventSpecifications()
}

func (r *customEventSpecificationResource) SetComputedFields(d *schema.ResourceData) {
}

func (r *customEventSpecificationResource) UpdateState(d *schema.ResourceData, obj restapi.InstanaDataObject, _ utils.ResourceNameFormatter) error {
	customEventSpecification := obj.(*restapi.CustomEventSpecification)

	d.SetId(customEventSpecification.ID)

	ruleSpec := customEventSpecification.Rules[0]
	ruleData, err := r.mapRuleToState(ruleSpec)
	if err != nil {
		return err
	}

	eventData := map[string]interface{}{
		CustomEventSpecificationFieldName:           customEventSpecification.Name,
		CustomEventSpecificationFieldQuery:          customEventSpecification.Query,
		CustomEventSpecificationFieldEntityType:     customEventSpecification.EntityType,
		CustomEventSpecificationFieldTriggering:     customEventSpecification.Triggering,
		CustomEventSpecificationFieldDescription:    customEventSpecification.Description,
		CustomEventSpecificationFieldExpirationTime: customEventSpecification.ExpirationTime,
		CustomEventSpecificationFieldEnabled:        customEventSpecification.Enabled,
		CustomEventSpecificationFieldRule:           []interface{}{ruleData},
	}

	return tfutils.UpdateState(d, eventData)
}

func (r *customEventSpecificationResource) mapRuleToState(rule restapi.RuleSpecification) (map[string]interface{}, error) {
	severity, err := ConvertSeverityFromInstanaAPIToTerraformRepresentation(rule.Severity)
	if err != nil {
		return nil, err
	}

	if rule.DType == restapi.EntityVerificationRuleType {
		return r.mapEntityEventSpecificationRuleToState(rule, severity)
	} else if rule.DType == restapi.SystemRuleType {
		return r.mapSystemRuleToState(rule, severity)
	} else if rule.DType == restapi.ThresholdRuleType {
		return r.mapThresholdRuleToState(rule, severity)
	} else {
		return nil, fmt.Errorf("unsupported rule type %s", rule.DType)
	}
}

func (r *customEventSpecificationResource) mapEntityEventSpecificationRuleToState(rule restapi.RuleSpecification, severity string) (map[string]interface{}, error) {
	return map[string]interface{}{
		CustomEventSpecificationFieldEntityVerificationRule: []interface{}{
			map[string]interface{}{
				CustomEventSpecificationRuleFieldSeverity:                              severity,
				CustomEventSpecificationEntityVerificationRuleFieldMatchingEntityLabel: rule.MatchingEntityLabel,
				CustomEventSpecificationEntityVerificationRuleFieldMatchingEntityType:  rule.MatchingEntityType,
				CustomEventSpecificationEntityVerificationRuleFieldMatchingOperator:    rule.MatchingOperator,
				CustomEventSpecificationEntityVerificationRuleFieldOfflineDuration:     rule.OfflineDuration,
			},
		},
	}, nil
}

func (r *customEventSpecificationResource) mapSystemRuleToState(rule restapi.RuleSpecification, severity string) (map[string]interface{}, error) {
	return map[string]interface{}{
		CustomEventSpecificationFieldSystemRule: []interface{}{
			map[string]interface{}{
				CustomEventSpecificationRuleFieldSeverity:           severity,
				CustomEventSpecificationSystemRuleFieldSystemRuleId: rule.SystemRuleID,
			},
		},
	}, nil
}

func (r *customEventSpecificationResource) mapThresholdRuleToState(rule restapi.RuleSpecification, severity string) (map[string]interface{}, error) {
	var metricPattern []interface{}
	if rule.MetricPattern != nil {
		metricPattern = []interface{}{
			map[string]interface{}{
				CustomEventSpecificationThresholdRuleFieldMetricPatternPrefix:      rule.MetricPattern.Prefix,
				CustomEventSpecificationThresholdRuleFieldMetricPatternPostfix:     rule.MetricPattern.Postfix,
				CustomEventSpecificationThresholdRuleFieldMetricPatternPlaceholder: rule.MetricPattern.Placeholder,
				CustomEventSpecificationThresholdRuleFieldMetricPatternOperator:    rule.MetricPattern.Operator,
			},
		}
	}

	return map[string]interface{}{
		CustomEventSpecificationFieldEntityVerificationRule: []interface{}{
			map[string]interface{}{
				CustomEventSpecificationRuleFieldSeverity:                   severity,
				CustomEventSpecificationThresholdRuleFieldMetricName:        rule.MetricName,
				CustomEventSpecificationThresholdRuleFieldMetricPattern:     metricPattern,
				CustomEventSpecificationThresholdRuleFieldRollup:            rule.Rollup,
				CustomEventSpecificationThresholdRuleFieldWindow:            rule.Window,
				CustomEventSpecificationThresholdRuleFieldAggregation:       rule.Aggregation,
				CustomEventSpecificationThresholdRuleFieldConditionOperator: rule.ConditionOperator,
				CustomEventSpecificationThresholdRuleFieldConditionValue:    rule.ConditionValue,
			},
		},
	}, nil
}

func (r *customEventSpecificationResource) MapStateToDataObject(d *schema.ResourceData, _ utils.ResourceNameFormatter) (restapi.InstanaDataObject, error) {
	severity, err := ConvertSeverityFromTerraformToInstanaAPIRepresentation(d.Get(CustomEventSpecificationRuleSeverity).(string))
	if err != nil {
		return &restapi.CustomEventSpecification{}, err
	}

	rule, err := r.mapRuleFromState(severity, d.Get(CustomEventSpecificationFieldRule).([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return nil, err
	}

	apiModel := restapi.CustomEventSpecification{
		ID:             d.Id(),
		Name:           d.Get(CustomEventSpecificationFieldName).(string),
		EntityType:     d.Get(CustomEventSpecificationFieldEntityType).(string),
		Query:          GetStringPointerFromResourceData(d, CustomEventSpecificationFieldQuery),
		Triggering:     d.Get(CustomEventSpecificationFieldTriggering).(bool),
		Description:    GetStringPointerFromResourceData(d, CustomEventSpecificationFieldDescription),
		ExpirationTime: GetIntPointerFromResourceData(d, CustomEventSpecificationFieldExpirationTime),
		Enabled:        d.Get(CustomEventSpecificationFieldEnabled).(bool),
		Rules:          []restapi.RuleSpecification{rule},
	}
	return &apiModel, nil
}

func (r *customEventSpecificationResource) mapRuleFromState(severity int, ruleData map[string]interface{}) (restapi.RuleSpecification, error) {
	if rule, ok := ruleData[CustomEventSpecificationFieldEntityVerificationRule]; ok {
		return r.mapEntityVerificationRuleFromState(rule.([]interface{})[0].(map[string]interface{}), severity), nil
	}
	if rule, ok := ruleData[CustomEventSpecificationFieldSystemRule]; ok {
		return r.mapSystemRuleFromState(rule.([]interface{})[0].(map[string]interface{}), severity), nil
	}
	if rule, ok := ruleData[CustomEventSpecificationFieldThresholdRule]; ok {
		return r.mapThresholdRuleFromState(rule.([]interface{})[0].(map[string]interface{}), severity), nil
	}

	return restapi.RuleSpecification{}, errors.New("no supported rule defined")
}

func (r *customEventSpecificationResource) mapEntityVerificationRuleFromState(rule map[string]interface{}, severity int) restapi.RuleSpecification {
	return restapi.RuleSpecification{
		Severity:            severity,
		MatchingEntityLabel: GetPointerFromMap[string](rule, CustomEventSpecificationEntityVerificationRuleFieldMatchingEntityLabel),
		MatchingEntityType:  GetPointerFromMap[string](rule, CustomEventSpecificationEntityVerificationRuleFieldMatchingEntityType),
		MatchingOperator:    GetPointerFromMap[string](rule, CustomEventSpecificationEntityVerificationRuleFieldMatchingOperator),
		OfflineDuration:     GetPointerFromMap[int](rule, CustomEventSpecificationEntityVerificationRuleFieldOfflineDuration),
	}
}

func (r *customEventSpecificationResource) mapSystemRuleFromState(rule map[string]interface{}, severity int) restapi.RuleSpecification {
	return restapi.RuleSpecification{
		Severity:     severity,
		SystemRuleID: GetPointerFromMap[string](rule, CustomEventSpecificationSystemRuleFieldSystemRuleId),
	}
}

func (r *customEventSpecificationResource) mapThresholdRuleFromState(rule map[string]interface{}, severity int) restapi.RuleSpecification {
	var metricPattern *restapi.MetricPattern
	if val, ok := rule[CustomEventSpecificationThresholdRuleFieldMetricPattern]; ok {
		metricPatternData := val.([]interface{})[0].(map[string]interface{})
		metricPatternObj := restapi.MetricPattern{
			Prefix:      metricPatternData[CustomEventSpecificationThresholdRuleFieldMetricPatternPrefix].(string),
			Postfix:     GetPointerFromMap[string](rule, CustomEventSpecificationThresholdRuleFieldMetricPatternPostfix),
			Placeholder: GetPointerFromMap[string](rule, CustomEventSpecificationThresholdRuleFieldMetricPatternPlaceholder),
			Operator:    restapi.MetricPatternOperatorType(metricPatternData[CustomEventSpecificationThresholdRuleFieldMetricPatternOperator].(string)),
		}
		metricPattern = &metricPatternObj
	}

	var aggregation *restapi.AggregationType
	if val, ok := rule[CustomEventSpecificationThresholdRuleFieldAggregation]; ok {
		agg := restapi.AggregationType(val.(string))
		aggregation = &agg
	}

	return restapi.RuleSpecification{
		Severity:          severity,
		MetricName:        GetPointerFromMap[string](rule, CustomEventSpecificationThresholdRuleFieldMetricName),
		MetricPattern:     metricPattern,
		Rollup:            GetPointerFromMap[int](rule, CustomEventSpecificationThresholdRuleFieldRollup),
		Window:            GetPointerFromMap[int](rule, CustomEventSpecificationThresholdRuleFieldWindow),
		Aggregation:       aggregation,
		ConditionOperator: GetPointerFromMap[string](rule, CustomEventSpecificationThresholdRuleFieldConditionOperator),
		ConditionValue:    GetPointerFromMap[float64](rule, CustomEventSpecificationThresholdRuleFieldConditionValue),
	}
}
