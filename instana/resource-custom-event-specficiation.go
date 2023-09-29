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

	CustomEventSpecificationFieldRules                       = "rules"
	CustomEventSpecificationFieldEntityCountRule             = "entity_count"
	CustomEventSpecificationFieldEntityCountVerificationRule = "entity_count_verification"
	CustomEventSpecificationFieldEntityVerificationRule      = "entity_verification"
	CustomEventSpecificationFieldHostAvailabilityRule        = "host_availability"
	CustomEventSpecificationFieldSystemRule                  = "system"
	CustomEventSpecificationFieldThresholdRule               = "threshold"
	customEventSpecificationThresholdRulePath                = "rules.0.threshold.0."

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
		"rules.0.entity_verification",
		"rules.0.system",
		"rules.0.threshold",
	}

	customEventSpecificationThresholdRuleMetricNameOrPattern = []string{
		customEventSpecificationThresholdRulePath + CustomEventSpecificationThresholdRuleFieldMetricName,
		customEventSpecificationThresholdRulePath + CustomEventSpecificationThresholdRuleFieldMetricPattern,
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
								MaxItems:    1,
								Optional:    true,
								Description: "Threshold rule configuration",
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
														ValidateFunc: validation.StringInSlice([]string{"is", "contains", "startsWith", "endsWith"}, false),
														Description:  "The metric pattern operator (e.g is, contains, startsWith, endsWith)",
													},
												},
											},
											ExactlyOneOf: customEventSpecificationThresholdRuleMetricNameOrPattern,
										},

										CustomEventSpecificationThresholdRuleFieldRollup: {
											Type:         schema.TypeInt,
											Optional:     true,
											ValidateFunc: validation.IntAtLeast(1),
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

func (r *customEventSpecificationResource) MetaData() *ResourceMetaData {
	return &r.metaData
}

func (r *customEventSpecificationResource) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{}
}

func (r *customEventSpecificationResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.CustomEventSpecification] {
	return api.CustomEventSpecifications()
}

func (r *customEventSpecificationResource) SetComputedFields(_ *schema.ResourceData) error {
	return nil
}

func (r *customEventSpecificationResource) UpdateState(d *schema.ResourceData, customEventSpecification *restapi.CustomEventSpecification) error {
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
		CustomEventSpecificationFieldRules:          []interface{}{ruleData},
	}

	d.SetId(customEventSpecification.ID)
	return tfutils.UpdateState(d, eventData)
}

func (r *customEventSpecificationResource) mapRuleToState(rule restapi.RuleSpecification) (map[string]interface{}, error) {
	severity, err := ConvertSeverityFromInstanaAPIToTerraformRepresentation(rule.Severity)
	if err != nil {
		return nil, err
	}

	if rule.DType == restapi.EntityCountRuleType {
		return r.mapEntityCountRuleToState(rule, severity)
	}
	if rule.DType == restapi.EntityCountVerificationRuleType {
		return r.mapEntityCountVerificationRuleToState(rule, severity)
	}
	if rule.DType == restapi.EntityVerificationRuleType {
		return r.mapEntityVerificationRuleToState(rule, severity)
	}
	if rule.DType == restapi.HostAvailabilityRuleType {
		return r.mapHostAvailabilityRuleToState(rule, severity)
	}
	if rule.DType == restapi.SystemRuleType {
		return r.mapSystemRuleToState(rule, severity)
	}
	if rule.DType == restapi.ThresholdRuleType {
		return r.mapThresholdRuleToState(rule, severity)
	}
	return nil, fmt.Errorf("unsupported rule type %s", rule.DType)
}

func (r *customEventSpecificationResource) mapEntityCountRuleToState(rule restapi.RuleSpecification, severity string) (map[string]interface{}, error) {
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

func (r *customEventSpecificationResource) mapEntityCountVerificationRuleToState(rule restapi.RuleSpecification, severity string) (map[string]interface{}, error) {
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

func (r *customEventSpecificationResource) mapEntityVerificationRuleToState(rule restapi.RuleSpecification, severity string) (map[string]interface{}, error) {
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

func (r *customEventSpecificationResource) mapHostAvailabilityRuleToState(rule restapi.RuleSpecification, severity string) (map[string]interface{}, error) {
	ruleData := map[string]interface{}{
		CustomEventSpecificationRuleFieldSeverity:                         severity,
		CustomEventSpecificationRuleFieldOfflineDuration:                  rule.OfflineDuration,
		CustomEventSpecificationHostAvailabilityRuleFieldMetricCloseAfter: rule.CloseAfter,
	}
	if rule.TagFilter != nil {
		normalizedTagFilterString, err := tagfilter.MapTagFilterToNormalizedString(rule.TagFilter.(restapi.TagFilterExpressionElement))
		if err != nil {
			return nil, err
		}
		ruleData[ApplicationConfigFieldTagFilter] = normalizedTagFilterString
	}

	return map[string]interface{}{CustomEventSpecificationFieldHostAvailabilityRule: []interface{}{ruleData}}, nil
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
		CustomEventSpecificationFieldThresholdRule: []interface{}{
			map[string]interface{}{
				CustomEventSpecificationRuleFieldSeverity:               severity,
				CustomEventSpecificationThresholdRuleFieldMetricName:    rule.MetricName,
				CustomEventSpecificationThresholdRuleFieldMetricPattern: metricPattern,
				CustomEventSpecificationThresholdRuleFieldRollup:        rule.Rollup,
				CustomEventSpecificationThresholdRuleFieldWindow:        rule.Window,
				CustomEventSpecificationThresholdRuleFieldAggregation:   rule.Aggregation,
				CustomEventSpecificationRuleFieldConditionOperator:      rule.ConditionOperator,
				CustomEventSpecificationRuleFieldConditionValue:         rule.ConditionValue,
			},
		},
	}, nil
}

func (r *customEventSpecificationResource) MapStateToDataObject(d *schema.ResourceData) (*restapi.CustomEventSpecification, error) {
	rule, err := r.mapRuleFromState(d.Get(CustomEventSpecificationFieldRules).([]interface{})[0].(map[string]interface{}))
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

func (r *customEventSpecificationResource) mapRuleFromState(ruleData map[string]interface{}) (restapi.RuleSpecification, error) {
	if rule, ok := ruleData[CustomEventSpecificationFieldEntityCountRule]; ok && len(rule.([]interface{})) > 0 {
		return r.mapEntityCountRuleFromState(rule.([]interface{})[0].(map[string]interface{}))
	}
	if rule, ok := ruleData[CustomEventSpecificationFieldEntityCountVerificationRule]; ok && len(rule.([]interface{})) > 0 {
		return r.mapEntityCountVerificationRuleFromState(rule.([]interface{})[0].(map[string]interface{}))
	}
	if rule, ok := ruleData[CustomEventSpecificationFieldEntityVerificationRule]; ok && len(rule.([]interface{})) > 0 {
		return r.mapEntityVerificationRuleFromState(rule.([]interface{})[0].(map[string]interface{}))
	}
	if rule, ok := ruleData[CustomEventSpecificationFieldHostAvailabilityRule]; ok && len(rule.([]interface{})) > 0 {
		return r.mapHostAvailabilityRuleFromState(rule.([]interface{})[0].(map[string]interface{}))
	}
	if rule, ok := ruleData[CustomEventSpecificationFieldSystemRule]; ok && len(rule.([]interface{})) > 0 {
		return r.mapSystemRuleFromState(rule.([]interface{})[0].(map[string]interface{}))
	}
	if rule, ok := ruleData[CustomEventSpecificationFieldThresholdRule]; ok && len(rule.([]interface{})) > 0 {
		return r.mapThresholdRuleFromState(rule.([]interface{})[0].(map[string]interface{}))
	}

	return restapi.RuleSpecification{}, errors.New("no supported rule defined")
}

func (r *customEventSpecificationResource) mapEntityCountRuleFromState(rule map[string]interface{}) (restapi.RuleSpecification, error) {
	severity, err := ConvertSeverityFromTerraformToInstanaAPIRepresentation(rule[CustomEventSpecificationRuleFieldSeverity].(string))
	if err != nil {
		return restapi.RuleSpecification{}, err
	}
	return restapi.RuleSpecification{
		DType:             restapi.EntityCountVerificationRuleType,
		Severity:          severity,
		ConditionOperator: GetPointerFromMap[string](rule, CustomEventSpecificationRuleFieldConditionValue),
		ConditionValue:    GetPointerFromMap[float64](rule, CustomEventSpecificationRuleFieldConditionOperator),
	}, nil
}

func (r *customEventSpecificationResource) mapEntityCountVerificationRuleFromState(rule map[string]interface{}) (restapi.RuleSpecification, error) {
	severity, err := ConvertSeverityFromTerraformToInstanaAPIRepresentation(rule[CustomEventSpecificationRuleFieldSeverity].(string))
	if err != nil {
		return restapi.RuleSpecification{}, err
	}
	return restapi.RuleSpecification{
		DType:               restapi.EntityCountVerificationRuleType,
		Severity:            severity,
		ConditionOperator:   GetPointerFromMap[string](rule, CustomEventSpecificationRuleFieldConditionValue),
		ConditionValue:      GetPointerFromMap[float64](rule, CustomEventSpecificationRuleFieldConditionOperator),
		MatchingEntityLabel: GetPointerFromMap[string](rule, CustomEventSpecificationRuleFieldMatchingEntityLabel),
		MatchingEntityType:  GetPointerFromMap[string](rule, CustomEventSpecificationRuleFieldMatchingEntityType),
		MatchingOperator:    GetPointerFromMap[string](rule, CustomEventSpecificationRuleFieldMatchingOperator),
	}, nil
}

func (r *customEventSpecificationResource) mapEntityVerificationRuleFromState(rule map[string]interface{}) (restapi.RuleSpecification, error) {
	severity, err := ConvertSeverityFromTerraformToInstanaAPIRepresentation(rule[CustomEventSpecificationRuleFieldSeverity].(string))
	if err != nil {
		return restapi.RuleSpecification{}, err
	}
	return restapi.RuleSpecification{
		DType:               restapi.EntityVerificationRuleType,
		Severity:            severity,
		MatchingEntityLabel: GetPointerFromMap[string](rule, CustomEventSpecificationRuleFieldMatchingEntityLabel),
		MatchingEntityType:  GetPointerFromMap[string](rule, CustomEventSpecificationRuleFieldMatchingEntityType),
		MatchingOperator:    GetPointerFromMap[string](rule, CustomEventSpecificationRuleFieldMatchingOperator),
		OfflineDuration:     GetPointerFromMap[int](rule, CustomEventSpecificationRuleFieldOfflineDuration),
	}, nil
}

func (r *customEventSpecificationResource) mapHostAvailabilityRuleFromState(rule map[string]interface{}) (restapi.RuleSpecification, error) {
	severity, err := ConvertSeverityFromTerraformToInstanaAPIRepresentation(rule[CustomEventSpecificationRuleFieldSeverity].(string))
	if err != nil {
		return restapi.RuleSpecification{}, err
	}
	var tagFilter restapi.TagFilterExpressionElement
	if tagFilterString, ok := rule[CustomEventSpecificationHostAvailabilityRuleFieldTagFilter]; ok {
		tagFilter, err = r.mapTagFilterStringToAPIModel(tagFilterString.(string))
		if err != nil {
			return restapi.RuleSpecification{}, err
		}
	}

	return restapi.RuleSpecification{
		DType:           restapi.HostAvailabilityRuleType,
		Severity:        severity,
		OfflineDuration: GetPointerFromMap[int](rule, CustomEventSpecificationRuleFieldOfflineDuration),
		CloseAfter:      GetPointerFromMap[int](rule, CustomEventSpecificationHostAvailabilityRuleFieldMetricCloseAfter),
		TagFilter:       tagFilter,
	}, nil
}

func (r *customEventSpecificationResource) mapTagFilterStringToAPIModel(input string) (restapi.TagFilterExpressionElement, error) {
	parser := tagfilter.NewParser()
	expr, err := parser.Parse(input)
	if err != nil {
		return nil, err
	}

	mapper := tagfilter.NewMapper()
	return mapper.ToAPIModel(expr), nil
}

func (r *customEventSpecificationResource) mapSystemRuleFromState(rule map[string]interface{}) (restapi.RuleSpecification, error) {
	severity, err := ConvertSeverityFromTerraformToInstanaAPIRepresentation(rule[CustomEventSpecificationRuleFieldSeverity].(string))
	if err != nil {
		return restapi.RuleSpecification{}, err
	}
	return restapi.RuleSpecification{
		DType:        restapi.SystemRuleType,
		Severity:     severity,
		SystemRuleID: GetPointerFromMap[string](rule, CustomEventSpecificationSystemRuleFieldSystemRuleId),
	}, nil
}

func (r *customEventSpecificationResource) mapThresholdRuleFromState(rule map[string]interface{}) (restapi.RuleSpecification, error) {
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

	return restapi.RuleSpecification{
		DType:             restapi.ThresholdRuleType,
		Severity:          severity,
		MetricName:        GetPointerFromMap[string](rule, CustomEventSpecificationThresholdRuleFieldMetricName),
		MetricPattern:     metricPattern,
		Rollup:            GetPointerFromMap[int](rule, CustomEventSpecificationThresholdRuleFieldRollup),
		Window:            GetPointerFromMap[int](rule, CustomEventSpecificationThresholdRuleFieldWindow),
		Aggregation:       aggregation,
		ConditionOperator: GetPointerFromMap[string](rule, CustomEventSpecificationRuleFieldConditionOperator),
		ConditionValue:    GetPointerFromMap[float64](rule, CustomEventSpecificationRuleFieldConditionValue),
	}, nil
}
