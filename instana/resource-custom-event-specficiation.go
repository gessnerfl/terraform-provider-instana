package instana

import (
	"errors"
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/instana/tagfilter"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceInstanaCustomEventSpecification the name of the terraform-provider-instana resource to manage custom event specifications
const ResourceInstanaCustomEventSpecification = "instana_custom_event_specification"

const (
	CustomEventSpecificationFieldName           = "name"
	CustomEventSpecificationFieldEntityType     = "entity_type"
	CustomEventSpecificationFieldQuery          = "query"
	CustomEventSpecificationFieldTriggering     = "triggering"
	CustomEventSpecificationFieldDescription    = "description"
	CustomEventSpecificationFieldExpirationTime = "expiration_time"
	CustomEventSpecificationFieldEnabled        = "enabled"

	CustomEventSpecificationFieldRuleLogicalOperator         = "rule_logical_operator"
	CustomEventSpecificationFieldRules                       = "rules"
	CustomEventSpecificationFieldEntityCountRule             = "entity_count"
	CustomEventSpecificationFieldEntityCountVerificationRule = "entity_count_verification"
	CustomEventSpecificationFieldEntityVerificationRule      = "entity_verification"
	CustomEventSpecificationFieldHostAvailabilityRule        = "host_availability"
	CustomEventSpecificationFieldSystemRule                  = "system"
	CustomEventSpecificationFieldThresholdRule               = "threshold"

	CustomEventSpecificationRuleFieldSeverity                          = "severity"
	CustomEventSpecificationRuleFieldMatchingEntityType                = "matching_entity_type"
	CustomEventSpecificationRuleFieldMatchingOperator                  = "matching_operator"
	CustomEventSpecificationRuleFieldMatchingEntityLabel               = "matching_entity_label"
	CustomEventSpecificationRuleFieldOfflineDuration                   = "offline_duration"
	CustomEventSpecificationSystemRuleFieldSystemRuleId                = "system_rule_id"
	CustomEventSpecificationThresholdRuleFieldMetricName               = "metric_name"
	CustomEventSpecificationThresholdRuleFieldRollup                   = "rollup"
	CustomEventSpecificationThresholdRuleFieldWindow                   = "window"
	CustomEventSpecificationThresholdRuleFieldAggregation              = "aggregation"
	CustomEventSpecificationRuleFieldConditionOperator                 = "condition_operator"
	CustomEventSpecificationRuleFieldConditionValue                    = "condition_value"
	CustomEventSpecificationThresholdRuleFieldMetricPattern            = "metric_pattern"
	CustomEventSpecificationThresholdRuleFieldMetricPatternPrefix      = "prefix"
	CustomEventSpecificationThresholdRuleFieldMetricPatternPostfix     = "postfix"
	CustomEventSpecificationThresholdRuleFieldMetricPatternPlaceholder = "placeholder"
	CustomEventSpecificationThresholdRuleFieldMetricPatternOperator    = "operator"
	CustomEventSpecificationHostAvailabilityRuleFieldMetricCloseAfter  = "close_after"
	CustomEventSpecificationHostAvailabilityRuleFieldTagFilter         = "tag_filter"
)

var (
	customEventSpecificationRuleTypeKeys = []string{
		"rules.0.entity_count",
		"rules.0.entity_count_verification",
		"rules.0.entity_verification",
		"rules.0.host_availability",
		"rules.0.system",
		"rules.0.threshold",
	}
)

var (
	customEventSpecificationSchemaRuleSeverity = &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice(restapi.SupportedSeverities.TerraformRepresentations(), false),
		Description:  "Configures the severity of the rule of the custom event specification",
	}
	customEventSpecificationSchemaRuleConditionOperator = &schema.Schema{
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
	}
	customEventSpecificationSchemaRuleConditionValue = &schema.Schema{
		Type:        schema.TypeFloat,
		Required:    true,
		Description: "The expected condition value to fulfill the rule",
	}
	customEventSpecificationSchemaRuleMatchingEntityType = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The type of the matching entity",
	}
	customEventSpecificationSchemaRuleMatchingOperator = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		ValidateFunc: validation.StringInSlice([]string{
			"is",
			"contains",
			"startsWith",
			"endsWith"}, false),
		Description: "The operator which should be applied for matching the label for the given entity (e.g. is, contains, startsWith, endsWith)",
	}
	customEventSpecificationSchemaRuleMatchingEntityLabel = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The label of the matching entity",
	}
	customEventSpecificationSchemaRuleOfflineDuration = &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "The duration after which the matching entity is considered to be offline",
	}
)

// NewCustomEventSpecificationResourceHandle creates a new ResourceHandle for the terraform resource of custom event specifications
func NewCustomEventSpecificationResourceHandle() ResourceHandle[*restapi.CustomEventSpecification] {
	return &customEventSpecificationResource{
		metaData: ResourceMetaData{
			ResourceName: ResourceInstanaCustomEventSpecification,
			Schema: map[string]*schema.Schema{
				CustomEventSpecificationFieldName: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Configures the name of the custom event specification",
				},
				CustomEventSpecificationFieldEntityType: {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Configures the entity type of the custom event specification. This value must be set to 'host' for entity verification rules and 'any' in case of system rules.",
				},
				CustomEventSpecificationFieldQuery: {
					Type:        schema.TypeString,
					Required:    false,
					Optional:    true,
					Description: "Configures the dynamic focus query for the custom event specification",
				},
				CustomEventSpecificationFieldTriggering: {
					Type:        schema.TypeBool,
					Default:     false,
					Optional:    true,
					Description: "Configures the custom event specification should trigger an incident",
				},
				CustomEventSpecificationFieldDescription: {
					Type:        schema.TypeString,
					Required:    false,
					Optional:    true,
					Description: "Configures the description text of the custom event specification",
				},
				CustomEventSpecificationFieldExpirationTime: {
					Type:        schema.TypeInt,
					Required:    false,
					Optional:    true,
					Description: "Configures the expiration time (grace period) to wait before the issue is closed",
				},
				CustomEventSpecificationFieldEnabled: {
					Type:        schema.TypeBool,
					Default:     true,
					Optional:    true,
					Description: "Configures if the custom event specification is enabled or not",
				},
				CustomEventSpecificationFieldRuleLogicalOperator: {
					Type:         schema.TypeString,
					Default:      "AND",
					Optional:     true,
					Description:  "The logical operator to be applied when multiple threshold rules are defined",
					ValidateFunc: validation.StringInSlice([]string{"AND", "OR"}, false),
				},
				CustomEventSpecificationFieldRules: {
					Type:        schema.TypeList,
					MinItems:    1,
					MaxItems:    1,
					Required:    true,
					Description: "Indicates the type of rule this custom event specification is about.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							CustomEventSpecificationFieldEntityCountRule: {
								Type:        schema.TypeList,
								MinItems:    0,
								MaxItems:    1,
								Optional:    true,
								Description: "Entity count rule configuration",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										CustomEventSpecificationRuleFieldSeverity:          customEventSpecificationSchemaRuleSeverity,
										CustomEventSpecificationRuleFieldConditionOperator: customEventSpecificationSchemaRuleConditionOperator,
										CustomEventSpecificationRuleFieldConditionValue:    customEventSpecificationSchemaRuleConditionValue,
									},
								},
								ExactlyOneOf: customEventSpecificationRuleTypeKeys,
							},
							CustomEventSpecificationFieldEntityCountVerificationRule: {
								Type:        schema.TypeList,
								MinItems:    0,
								MaxItems:    1,
								Optional:    true,
								Description: "Entity count rule configuration",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										CustomEventSpecificationRuleFieldSeverity:            customEventSpecificationSchemaRuleSeverity,
										CustomEventSpecificationRuleFieldConditionOperator:   customEventSpecificationSchemaRuleConditionOperator,
										CustomEventSpecificationRuleFieldConditionValue:      customEventSpecificationSchemaRuleConditionValue,
										CustomEventSpecificationRuleFieldMatchingEntityType:  customEventSpecificationSchemaRuleMatchingEntityType,
										CustomEventSpecificationRuleFieldMatchingOperator:    customEventSpecificationSchemaRuleMatchingOperator,
										CustomEventSpecificationRuleFieldMatchingEntityLabel: customEventSpecificationSchemaRuleMatchingEntityLabel,
									},
								},
								ExactlyOneOf: customEventSpecificationRuleTypeKeys,
							},
							CustomEventSpecificationFieldEntityVerificationRule: {
								Type:        schema.TypeList,
								MinItems:    0,
								MaxItems:    1,
								Optional:    true,
								Description: "Entity verification rule configuration",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										CustomEventSpecificationRuleFieldSeverity:            customEventSpecificationSchemaRuleSeverity,
										CustomEventSpecificationRuleFieldMatchingEntityType:  customEventSpecificationSchemaRuleMatchingEntityType,
										CustomEventSpecificationRuleFieldMatchingOperator:    customEventSpecificationSchemaRuleMatchingOperator,
										CustomEventSpecificationRuleFieldMatchingEntityLabel: customEventSpecificationSchemaRuleMatchingEntityLabel,
										CustomEventSpecificationRuleFieldOfflineDuration:     customEventSpecificationSchemaRuleOfflineDuration,
									},
								},
								ExactlyOneOf: customEventSpecificationRuleTypeKeys,
							},
							CustomEventSpecificationFieldHostAvailabilityRule: {
								Type:        schema.TypeList,
								MinItems:    0,
								MaxItems:    1,
								Optional:    true,
								Description: "Host availability rule configuration",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										CustomEventSpecificationRuleFieldSeverity:        customEventSpecificationSchemaRuleSeverity,
										CustomEventSpecificationRuleFieldOfflineDuration: customEventSpecificationSchemaRuleOfflineDuration,
										CustomEventSpecificationHostAvailabilityRuleFieldMetricCloseAfter: {
											Type:        schema.TypeInt,
											Required:    true,
											Description: "if a host is offline for longer than the defined period, Instana does not expect the host to reappear anymore, and the event will be closed after the grace period",
										},
										CustomEventSpecificationHostAvailabilityRuleFieldTagFilter: RequiredTagFilterExpressionSchema,
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
								MaxItems:    5,
								Optional:    true,
								Description: "Threshold rule configuration",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										CustomEventSpecificationRuleFieldSeverity: customEventSpecificationSchemaRuleSeverity,

										CustomEventSpecificationThresholdRuleFieldMetricName: {
											Type:        schema.TypeString,
											Required:    false,
											Optional:    true,
											Description: "The metric name of the rule",
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
														Required:    true,
														Description: "The metric pattern prefix of a dynamic built-in metrics",
													},
													CustomEventSpecificationThresholdRuleFieldMetricPatternPostfix: {
														Type:        schema.TypeString,
														Optional:    true,
														Description: "The metric pattern postfix of a dynamic built-in metrics",
													},
													CustomEventSpecificationThresholdRuleFieldMetricPatternPlaceholder: {
														Type:        schema.TypeString,
														Required:    true,
														Description: "The metric pattern placeholder/condition value of a dynamic built-in metrics",
													},
													CustomEventSpecificationThresholdRuleFieldMetricPatternOperator: {
														Type:         schema.TypeString,
														Required:     true,
														ValidateFunc: validation.StringInSlice([]string{"is", "contains", "startsWith", "endsWith", "any"}, false),
														Description:  "The metric pattern operator (e.g is, contains, startsWith, endsWith, any)",
													},
												},
											},
										},

										CustomEventSpecificationThresholdRuleFieldRollup: {
											Type:         schema.TypeInt,
											Optional:     true,
											ValidateFunc: validation.IntAtLeast(0),
											Description:  "The rollup of the metric",
										},
										CustomEventSpecificationThresholdRuleFieldWindow: {
											Type:         schema.TypeInt,
											Required:     true,
											ValidateFunc: validation.IntAtLeast(1),
											Description:  "The time window where the condition has to be fulfilled",
										},
										CustomEventSpecificationThresholdRuleFieldAggregation: {
											Type:         schema.TypeString,
											Required:     true,
											ValidateFunc: validation.StringInSlice([]string{"sum", "avg", "min", "max"}, false),
											Description:  "The aggregation type (e.g. sum, avg)",
										},
										CustomEventSpecificationRuleFieldConditionOperator: customEventSpecificationSchemaRuleConditionOperator,
										CustomEventSpecificationRuleFieldConditionValue:    customEventSpecificationSchemaRuleConditionValue,
									},
								},
								ExactlyOneOf: customEventSpecificationRuleTypeKeys,
							},
						},
					},
				},
			},
			SchemaVersion: 0,
		},
	}
}

type customEventSpecificationResource struct {
	metaData ResourceMetaData
}

func (c *customEventSpecificationResource) MetaData() *ResourceMetaData {
	return &c.metaData
}

func (c *customEventSpecificationResource) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{}
}

func (c *customEventSpecificationResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.CustomEventSpecification] {
	return api.CustomEventSpecifications()
}

func (c *customEventSpecificationResource) SetComputedFields(_ *schema.ResourceData) error {
	return nil
}

func (c *customEventSpecificationResource) UpdateState(d *schema.ResourceData, customEventSpecification *restapi.CustomEventSpecification) error {
	ruleData, err := c.mapRulesToState(customEventSpecification)
	if err != nil {
		return err
	}

	eventData := map[string]interface{}{
		CustomEventSpecificationFieldName:                customEventSpecification.Name,
		CustomEventSpecificationFieldQuery:               customEventSpecification.Query,
		CustomEventSpecificationFieldEntityType:          customEventSpecification.EntityType,
		CustomEventSpecificationFieldTriggering:          customEventSpecification.Triggering,
		CustomEventSpecificationFieldDescription:         customEventSpecification.Description,
		CustomEventSpecificationFieldExpirationTime:      customEventSpecification.ExpirationTime,
		CustomEventSpecificationFieldEnabled:             customEventSpecification.Enabled,
		CustomEventSpecificationFieldRuleLogicalOperator: customEventSpecification.RuleLogicalOperator,
		CustomEventSpecificationFieldRules:               []interface{}{ruleData},
	}

	d.SetId(customEventSpecification.ID)
	return tfutils.UpdateState(d, eventData)
}

func (c *customEventSpecificationResource) mapRulesToState(customEventSpecification *restapi.CustomEventSpecification) (map[string]interface{}, error) {
	if len(customEventSpecification.Rules) > 0 {
		severity, err := ConvertSeverityFromInstanaAPIToTerraformRepresentation(customEventSpecification.Rules[0].Severity)
		if err != nil {
			return nil, err
		}

		if customEventSpecification.Rules[0].DType == restapi.ThresholdRuleType {
			return c.mapThresholdRulesToState(customEventSpecification, severity)
		} else {
			//only a single rule is allowed for non threshold rules types
			rule := customEventSpecification.Rules[0]
			return c.mapNonThresholdRulesToState(rule, severity)
		}
	}
	return map[string]interface{}{}, nil
}

func (c *customEventSpecificationResource) mapNonThresholdRulesToState(rule restapi.RuleSpecification, severity string) (map[string]interface{}, error) {
	if rule.DType == restapi.EntityCountRuleType {
		return c.mapEntityCountRuleToState(rule, severity)
	}
	if rule.DType == restapi.EntityCountVerificationRuleType {
		return c.mapEntityCountVerificationRuleToState(rule, severity)
	}
	if rule.DType == restapi.EntityVerificationRuleType {
		return c.mapEntityVerificationRuleToState(rule, severity)
	}
	if rule.DType == restapi.HostAvailabilityRuleType {
		return c.mapHostAvailabilityRuleToState(rule, severity)
	}
	if rule.DType == restapi.SystemRuleType {
		return c.mapSystemRuleToState(rule, severity)
	}
	return nil, fmt.Errorf("unsupported rule type %s", rule.DType)
}

func (c *customEventSpecificationResource) mapEntityCountRuleToState(rule restapi.RuleSpecification, severity string) (map[string]interface{}, error) {
	return map[string]interface{}{
		CustomEventSpecificationFieldEntityCountRule: []interface{}{
			map[string]interface{}{
				CustomEventSpecificationRuleFieldSeverity:          severity,
				CustomEventSpecificationRuleFieldConditionOperator: rule.ConditionOperator,
				CustomEventSpecificationRuleFieldConditionValue:    rule.ConditionValue,
			},
		},
	}, nil
}

func (c *customEventSpecificationResource) mapEntityCountVerificationRuleToState(rule restapi.RuleSpecification, severity string) (map[string]interface{}, error) {
	return map[string]interface{}{
		CustomEventSpecificationFieldEntityCountVerificationRule: []interface{}{
			map[string]interface{}{
				CustomEventSpecificationRuleFieldSeverity:            severity,
				CustomEventSpecificationRuleFieldConditionOperator:   rule.ConditionOperator,
				CustomEventSpecificationRuleFieldConditionValue:      rule.ConditionValue,
				CustomEventSpecificationRuleFieldMatchingEntityLabel: rule.MatchingEntityLabel,
				CustomEventSpecificationRuleFieldMatchingEntityType:  rule.MatchingEntityType,
				CustomEventSpecificationRuleFieldMatchingOperator:    rule.MatchingOperator,
			},
		},
	}, nil
}

func (c *customEventSpecificationResource) mapEntityVerificationRuleToState(rule restapi.RuleSpecification, severity string) (map[string]interface{}, error) {
	return map[string]interface{}{
		CustomEventSpecificationFieldEntityVerificationRule: []interface{}{
			map[string]interface{}{
				CustomEventSpecificationRuleFieldSeverity:            severity,
				CustomEventSpecificationRuleFieldMatchingEntityLabel: rule.MatchingEntityLabel,
				CustomEventSpecificationRuleFieldMatchingEntityType:  rule.MatchingEntityType,
				CustomEventSpecificationRuleFieldMatchingOperator:    rule.MatchingOperator,
				CustomEventSpecificationRuleFieldOfflineDuration:     rule.OfflineDuration,
			},
		},
	}, nil
}

func (c *customEventSpecificationResource) mapHostAvailabilityRuleToState(rule restapi.RuleSpecification, severity string) (map[string]interface{}, error) {
	ruleData := map[string]interface{}{
		CustomEventSpecificationRuleFieldSeverity:                         severity,
		CustomEventSpecificationRuleFieldOfflineDuration:                  rule.OfflineDuration,
		CustomEventSpecificationHostAvailabilityRuleFieldMetricCloseAfter: rule.CloseAfter,
	}
	if rule.TagFilter != nil {
		normalizedTagFilterString, err := tagfilter.MapTagFilterToNormalizedString(rule.TagFilter)
		if err != nil {
			return nil, err
		}
		ruleData[CustomEventSpecificationHostAvailabilityRuleFieldTagFilter] = normalizedTagFilterString
	}

	return map[string]interface{}{CustomEventSpecificationFieldHostAvailabilityRule: []interface{}{ruleData}}, nil
}

func (c *customEventSpecificationResource) mapSystemRuleToState(rule restapi.RuleSpecification, severity string) (map[string]interface{}, error) {
	return map[string]interface{}{
		CustomEventSpecificationFieldSystemRule: []interface{}{
			map[string]interface{}{
				CustomEventSpecificationRuleFieldSeverity:           severity,
				CustomEventSpecificationSystemRuleFieldSystemRuleId: rule.SystemRuleID,
			},
		},
	}, nil
}

func (c *customEventSpecificationResource) mapThresholdRulesToState(customEventSpecification *restapi.CustomEventSpecification, severity string) (map[string]interface{}, error) {
	rules := make([]interface{}, len(customEventSpecification.Rules))
	for i, r := range customEventSpecification.Rules {
		rule, err := c.mapThresholdRuleToState(r, severity)
		if err != nil {
			return map[string]interface{}{}, err
		}
		rules[i] = rule
	}
	return map[string]interface{}{CustomEventSpecificationFieldThresholdRule: rules}, nil
}

func (c *customEventSpecificationResource) mapThresholdRuleToState(rule restapi.RuleSpecification, severity string) (map[string]interface{}, error) {
	var metricPattern []interface{}
	if rule.DType != restapi.ThresholdRuleType {
		return map[string]interface{}{}, errors.New("invalid rule specification; rules must be of type threshold rule when multiple rules are defined")
	}
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
		CustomEventSpecificationRuleFieldSeverity:               severity,
		CustomEventSpecificationThresholdRuleFieldMetricName:    rule.MetricName,
		CustomEventSpecificationThresholdRuleFieldMetricPattern: metricPattern,
		CustomEventSpecificationThresholdRuleFieldRollup:        rule.Rollup,
		CustomEventSpecificationThresholdRuleFieldWindow:        rule.Window,
		CustomEventSpecificationThresholdRuleFieldAggregation:   rule.Aggregation,
		CustomEventSpecificationRuleFieldConditionOperator:      rule.ConditionOperator,
		CustomEventSpecificationRuleFieldConditionValue:         rule.ConditionValue,
	}, nil
}

func (c *customEventSpecificationResource) MapStateToDataObject(d *schema.ResourceData) (*restapi.CustomEventSpecification, error) {
	rules, err := c.mapRulesFromState(d.Get(CustomEventSpecificationFieldRules).([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return nil, err
	}

	apiModel := restapi.CustomEventSpecification{
		ID:                  d.Id(),
		Name:                d.Get(CustomEventSpecificationFieldName).(string),
		EntityType:          d.Get(CustomEventSpecificationFieldEntityType).(string),
		Query:               GetStringPointerFromResourceData(d, CustomEventSpecificationFieldQuery),
		Triggering:          d.Get(CustomEventSpecificationFieldTriggering).(bool),
		Description:         GetStringPointerFromResourceData(d, CustomEventSpecificationFieldDescription),
		ExpirationTime:      GetIntPointerFromResourceData(d, CustomEventSpecificationFieldExpirationTime),
		Enabled:             d.Get(CustomEventSpecificationFieldEnabled).(bool),
		RuleLogicalOperator: d.Get(CustomEventSpecificationFieldRuleLogicalOperator).(string),
		Rules:               rules,
	}
	return &apiModel, nil
}

func (c *customEventSpecificationResource) mapRulesFromState(ruleData map[string]interface{}) ([]restapi.RuleSpecification, error) {
	if rule, ok := ruleData[CustomEventSpecificationFieldEntityCountRule]; ok && len(rule.([]interface{})) > 0 {
		return c.mapEntityCountRuleFromState(rule.([]interface{})[0].(map[string]interface{}))
	}
	if rule, ok := ruleData[CustomEventSpecificationFieldEntityCountVerificationRule]; ok && len(rule.([]interface{})) > 0 {
		return c.mapEntityCountVerificationRuleFromState(rule.([]interface{})[0].(map[string]interface{}))
	}
	if rule, ok := ruleData[CustomEventSpecificationFieldEntityVerificationRule]; ok && len(rule.([]interface{})) > 0 {
		return c.mapEntityVerificationRuleFromState(rule.([]interface{})[0].(map[string]interface{}))
	}
	if rule, ok := ruleData[CustomEventSpecificationFieldHostAvailabilityRule]; ok && len(rule.([]interface{})) > 0 {
		return c.mapHostAvailabilityRuleFromState(rule.([]interface{})[0].(map[string]interface{}))
	}
	if rule, ok := ruleData[CustomEventSpecificationFieldSystemRule]; ok && len(rule.([]interface{})) > 0 {
		return c.mapSystemRuleFromState(rule.([]interface{})[0].(map[string]interface{}))
	}
	if rule, ok := ruleData[CustomEventSpecificationFieldThresholdRule]; ok && len(rule.([]interface{})) > 0 {
		return c.mapThresholdRulesFromState(rule.([]interface{}))
	}

	return []restapi.RuleSpecification{}, errors.New("no supported rule defined")
}

func (c *customEventSpecificationResource) mapEntityCountRuleFromState(rule map[string]interface{}) ([]restapi.RuleSpecification, error) {
	severity, err := ConvertSeverityFromTerraformToInstanaAPIRepresentation(rule[CustomEventSpecificationRuleFieldSeverity].(string))
	if err != nil {
		return []restapi.RuleSpecification{}, err
	}
	return []restapi.RuleSpecification{{
		DType:             restapi.EntityCountRuleType,
		Severity:          severity,
		ConditionOperator: GetPointerFromMap[string](rule, CustomEventSpecificationRuleFieldConditionOperator),
		ConditionValue:    GetPointerFromMap[float64](rule, CustomEventSpecificationRuleFieldConditionValue),
	}}, nil
}

func (c *customEventSpecificationResource) mapEntityCountVerificationRuleFromState(rule map[string]interface{}) ([]restapi.RuleSpecification, error) {
	severity, err := ConvertSeverityFromTerraformToInstanaAPIRepresentation(rule[CustomEventSpecificationRuleFieldSeverity].(string))
	if err != nil {
		return []restapi.RuleSpecification{}, err
	}
	return []restapi.RuleSpecification{{
		DType:               restapi.EntityCountVerificationRuleType,
		Severity:            severity,
		ConditionOperator:   GetPointerFromMap[string](rule, CustomEventSpecificationRuleFieldConditionOperator),
		ConditionValue:      GetPointerFromMap[float64](rule, CustomEventSpecificationRuleFieldConditionValue),
		MatchingEntityLabel: GetPointerFromMap[string](rule, CustomEventSpecificationRuleFieldMatchingEntityLabel),
		MatchingEntityType:  GetPointerFromMap[string](rule, CustomEventSpecificationRuleFieldMatchingEntityType),
		MatchingOperator:    GetPointerFromMap[string](rule, CustomEventSpecificationRuleFieldMatchingOperator),
	}}, nil
}

func (c *customEventSpecificationResource) mapEntityVerificationRuleFromState(rule map[string]interface{}) ([]restapi.RuleSpecification, error) {
	severity, err := ConvertSeverityFromTerraformToInstanaAPIRepresentation(rule[CustomEventSpecificationRuleFieldSeverity].(string))
	if err != nil {
		return []restapi.RuleSpecification{}, err
	}
	return []restapi.RuleSpecification{{
		DType:               restapi.EntityVerificationRuleType,
		Severity:            severity,
		MatchingEntityLabel: GetPointerFromMap[string](rule, CustomEventSpecificationRuleFieldMatchingEntityLabel),
		MatchingEntityType:  GetPointerFromMap[string](rule, CustomEventSpecificationRuleFieldMatchingEntityType),
		MatchingOperator:    GetPointerFromMap[string](rule, CustomEventSpecificationRuleFieldMatchingOperator),
		OfflineDuration:     GetPointerFromMap[int](rule, CustomEventSpecificationRuleFieldOfflineDuration),
	}}, nil
}

func (c *customEventSpecificationResource) mapHostAvailabilityRuleFromState(rule map[string]interface{}) ([]restapi.RuleSpecification, error) {
	severity, err := ConvertSeverityFromTerraformToInstanaAPIRepresentation(rule[CustomEventSpecificationRuleFieldSeverity].(string))
	if err != nil {
		return []restapi.RuleSpecification{}, err
	}
	var tagFilter *restapi.TagFilter
	if tagFilterString, ok := rule[CustomEventSpecificationHostAvailabilityRuleFieldTagFilter]; ok {
		tagFilter, err = c.mapTagFilterStringToAPIModel(tagFilterString.(string))
		if err != nil {
			return []restapi.RuleSpecification{}, err
		}
	}

	return []restapi.RuleSpecification{{
		DType:           restapi.HostAvailabilityRuleType,
		Severity:        severity,
		OfflineDuration: GetPointerFromMap[int](rule, CustomEventSpecificationRuleFieldOfflineDuration),
		CloseAfter:      GetPointerFromMap[int](rule, CustomEventSpecificationHostAvailabilityRuleFieldMetricCloseAfter),
		TagFilter:       tagFilter,
	}}, nil
}

func (c *customEventSpecificationResource) mapTagFilterStringToAPIModel(input string) (*restapi.TagFilter, error) {
	parser := tagfilter.NewParser()
	expr, err := parser.Parse(input)
	if err != nil {
		return nil, err
	}

	mapper := tagfilter.NewMapper()
	return mapper.ToAPIModel(expr), nil
}

func (c *customEventSpecificationResource) mapSystemRuleFromState(rule map[string]interface{}) ([]restapi.RuleSpecification, error) {
	severity, err := ConvertSeverityFromTerraformToInstanaAPIRepresentation(rule[CustomEventSpecificationRuleFieldSeverity].(string))
	if err != nil {
		return []restapi.RuleSpecification{}, err
	}
	return []restapi.RuleSpecification{{
		DType:        restapi.SystemRuleType,
		Severity:     severity,
		SystemRuleID: GetPointerFromMap[string](rule, CustomEventSpecificationSystemRuleFieldSystemRuleId),
	}}, nil
}

func (c *customEventSpecificationResource) mapThresholdRulesFromState(rules []interface{}) ([]restapi.RuleSpecification, error) {
	result := make([]restapi.RuleSpecification, len(rules))
	for i, r := range rules {
		rule, err := c.mapThresholdRuleFromState(r.(map[string]interface{}))
		if err != nil {
			return result, err
		}
		result[i] = rule
	}
	return result, nil
}

func (c *customEventSpecificationResource) mapThresholdRuleFromState(rule map[string]interface{}) (restapi.RuleSpecification, error) {
	severity, err := ConvertSeverityFromTerraformToInstanaAPIRepresentation(rule[CustomEventSpecificationRuleFieldSeverity].(string))
	if err != nil {
		return restapi.RuleSpecification{}, err
	}

	var metricPattern *restapi.MetricPattern
	if val, ok := rule[CustomEventSpecificationThresholdRuleFieldMetricPattern]; ok && len(val.([]interface{})) > 0 {
		metricPatternData := val.([]interface{})[0].(map[string]interface{})
		metricPatternObj := restapi.MetricPattern{
			Prefix:      metricPatternData[CustomEventSpecificationThresholdRuleFieldMetricPatternPrefix].(string),
			Postfix:     GetPointerFromMap[string](metricPatternData, CustomEventSpecificationThresholdRuleFieldMetricPatternPostfix),
			Placeholder: GetPointerFromMap[string](metricPatternData, CustomEventSpecificationThresholdRuleFieldMetricPatternPlaceholder),
			Operator:    metricPatternData[CustomEventSpecificationThresholdRuleFieldMetricPatternOperator].(string),
		}
		metricPattern = &metricPatternObj
	}

	var aggregation *string
	if val, ok := rule[CustomEventSpecificationThresholdRuleFieldAggregation]; ok {
		agg := val.(string)
		aggregation = &agg
	}

	metricName := GetPointerFromMap[string](rule, CustomEventSpecificationThresholdRuleFieldMetricName)
	if (metricName != nil && metricPattern != nil) || (metricName == nil && metricPattern == nil) {
		return restapi.RuleSpecification{}, errors.New("either metric_metric name or metric_pattern must be defined")
	}

	return restapi.RuleSpecification{
		DType:             restapi.ThresholdRuleType,
		Severity:          severity,
		MetricName:        metricName,
		MetricPattern:     metricPattern,
		Rollup:            GetPointerFromMap[int](rule, CustomEventSpecificationThresholdRuleFieldRollup),
		Window:            GetPointerFromMap[int](rule, CustomEventSpecificationThresholdRuleFieldWindow),
		Aggregation:       aggregation,
		ConditionOperator: GetPointerFromMap[string](rule, CustomEventSpecificationRuleFieldConditionOperator),
		ConditionValue:    GetPointerFromMap[float64](rule, CustomEventSpecificationRuleFieldConditionValue),
	}, nil
}
