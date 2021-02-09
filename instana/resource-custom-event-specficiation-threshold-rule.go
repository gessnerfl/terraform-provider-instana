package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

//ResourceInstanaCustomEventSpecificationThresholdRule the name of the terraform-provider-instana resource to manage custom event specifications with threshold rule
const ResourceInstanaCustomEventSpecificationThresholdRule = "instana_custom_event_spec_threshold_rule"

const (
	//ThresholdRuleFieldMetricName constant value for the schema field rule_metric_name
	ThresholdRuleFieldMetricName = ruleFieldPrefix + "metric_name"
	//ThresholdRuleFieldRollup constant value for the schema field rule_rollup
	ThresholdRuleFieldRollup = ruleFieldPrefix + "rollup"
	//ThresholdRuleFieldWindow constant value for the schema field rule_window
	ThresholdRuleFieldWindow = ruleFieldPrefix + "window"
	//ThresholdRuleFieldAggregation constant value for the schema field rule_aggregation
	ThresholdRuleFieldAggregation = ruleFieldPrefix + "aggregation"
	//ThresholdRuleFieldConditionOperator constant value for the schema field rule_condition_operator
	ThresholdRuleFieldConditionOperator = ruleFieldPrefix + "condition_operator"
	//ThresholdRuleFieldConditionValue constant value for the schema field rule_condition_value
	ThresholdRuleFieldConditionValue = ruleFieldPrefix + "condition_value"

	thresholdRuleFieldMetricPattern = ruleFieldPrefix + "metric_pattern_"
	//ThresholdRuleFieldMetricPatternPrefix constant value for the schema field rule_metric_pattern_prefix
	ThresholdRuleFieldMetricPatternPrefix = thresholdRuleFieldMetricPattern + "prefix"
	//ThresholdRuleFieldMetricPatternPostfix constant value for the schema field rule_metric_pattern_postfix
	ThresholdRuleFieldMetricPatternPostfix = thresholdRuleFieldMetricPattern + "postfix"
	//ThresholdRuleFieldMetricPatternPlaceholder constant value for the schema field rule_metric_pattern_placeholder
	ThresholdRuleFieldMetricPatternPlaceholder = thresholdRuleFieldMetricPattern + "placeholder"
	//ThresholdRuleFieldMetricPatternOperator constant value for the schema field rule_metric_pattern_operator
	ThresholdRuleFieldMetricPatternOperator = thresholdRuleFieldMetricPattern + "operator"
)

var thresholdRuleSchemaFields = map[string]*schema.Schema{
	CustomEventSpecificationFieldEntityType: {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Configures the entity type of the custom event specification",
	},
	ThresholdRuleFieldMetricName: {
		Type:        schema.TypeString,
		Required:    false,
		Optional:    true,
		Description: "The metric name of the rule",
	},
	ThresholdRuleFieldRollup: {
		Type:        schema.TypeInt,
		Required:    false,
		Optional:    true,
		Description: "The rollup of the metric",
	},
	ThresholdRuleFieldWindow: {
		Type:        schema.TypeInt,
		Required:    false,
		Optional:    true,
		Description: "The time window where the condition has to be fulfilled",
	},
	ThresholdRuleFieldAggregation: {
		Type:         schema.TypeString,
		Required:     false,
		Optional:     true,
		ValidateFunc: validation.StringInSlice(restapi.SupportedAggregationTypes.ToStringSlice(), false),
		Description:  "The aggregation type (e.g. sum, avg)",
	},
	ThresholdRuleFieldConditionOperator: {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice(restapi.SupportedConditionOperators.TerrafromSupportedValues(), false),
		StateFunc: func(val interface{}) string {
			operator, _ := restapi.SupportedConditionOperators.FromTerraformValue(val.(string))
			return operator.InstanaAPIValue()
		},
		Description: "The condition operator (e.g >, <)",
	},
	ThresholdRuleFieldConditionValue: {
		Type:        schema.TypeFloat,
		Required:    true,
		Description: "The expected condition value to fulfill the rule",
	},
	ThresholdRuleFieldMetricPatternPrefix: {
		Type:        schema.TypeString,
		Required:    false,
		Optional:    true,
		Description: "The metric pattern prefix of a dynamic built-in metrics",
	},
	ThresholdRuleFieldMetricPatternPostfix: {
		Type:        schema.TypeString,
		Required:    false,
		Optional:    true,
		Description: "The metric pattern postfix of a dynamic built-in metrics",
	},
	ThresholdRuleFieldMetricPatternPlaceholder: {
		Type:        schema.TypeString,
		Required:    false,
		Optional:    true,
		Description: "The metric pattern placeholer/condition value of a dynamic built-in metrics",
	},
	ThresholdRuleFieldMetricPatternOperator: {
		Type:         schema.TypeString,
		Required:     false,
		Optional:     true,
		ValidateFunc: validation.StringInSlice(restapi.SupportedMetricPatternOperatorTypes.ToStringSlice(), false),
		Description:  "The condition operator (e.g >, <)",
	},
}

//NewCustomEventSpecificationWithThresholdRuleResourceHandle creates a new ResourceHandle for the terraform resource of custom event specifications with system rules
func NewCustomEventSpecificationWithThresholdRuleResourceHandle() *ResourceHandle {
	return &ResourceHandle{
		ResourceName:  ResourceInstanaCustomEventSpecificationThresholdRule,
		Schema:        mergeSchemaMap(defaultCustomEventSchemaFields, thresholdRuleSchemaFields),
		SchemaVersion: 3,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    customEventSpecificationWithThresholdRuleSchemaV0().CoreConfigSchema().ImpliedType(),
				Upgrade: migrateCustomEventConfigFullNameInStateFromV0toV1,
				Version: 0,
			},
			{
				Type:    customEventSpecificationWithThresholdRuleSchemaV1().CoreConfigSchema().ImpliedType(),
				Upgrade: migrateCustomEventConfigFullStateFromV1toV2AndRemoveDownstreamConfiguration,
				Version: 1,
			},
			{
				Type:    customEventSpecificationWithThresholdRuleSchemaV2().CoreConfigSchema().ImpliedType(),
				Upgrade: migrateCustomEventConfigWithThreasholdRuleToVersion3ByChangingConditionOperatorToInstanaRepresentation,
				Version: 2,
			},
		},
		RestResourceFactory:  func(api restapi.InstanaAPI) restapi.RestResource { return api.CustomEventSpecifications() },
		UpdateState:          updateStateForCustomEventSpecificationWithThresholdRule,
		MapStateToDataObject: mapStateToDataObjectForCustomEventSpecificationWithThresholdRule,
	}
}

func updateStateForCustomEventSpecificationWithThresholdRule(d *schema.ResourceData, obj restapi.InstanaDataObject) error {
	customEventSpecification := obj.(*restapi.CustomEventSpecification)
	ruleSpec := customEventSpecification.Rules[0]

	severity, err := ConvertSeverityFromInstanaAPIToTerraformRepresentation(ruleSpec.Severity)
	if err != nil {
		return err
	}
	conditionOperator, err := ruleSpec.ConditionOperatorType()
	if err != nil {
		return err
	}

	updateStateForBasicCustomEventSpecification(d, customEventSpecification)
	d.Set(CustomEventSpecificationRuleSeverity, severity)
	d.Set(ThresholdRuleFieldMetricName, ruleSpec.MetricName)
	d.Set(ThresholdRuleFieldRollup, ruleSpec.Rollup)
	d.Set(ThresholdRuleFieldWindow, ruleSpec.Window)
	d.Set(ThresholdRuleFieldAggregation, ruleSpec.Aggregation)
	d.Set(ThresholdRuleFieldConditionOperator, conditionOperator.InstanaAPIValue())
	d.Set(ThresholdRuleFieldConditionValue, ruleSpec.ConditionValue)

	if ruleSpec.MetricPattern != nil {
		d.Set(ThresholdRuleFieldMetricPatternPrefix, ruleSpec.MetricPattern.Prefix)
		d.Set(ThresholdRuleFieldMetricPatternPostfix, ruleSpec.MetricPattern.Postfix)
		d.Set(ThresholdRuleFieldMetricPatternPlaceholder, ruleSpec.MetricPattern.Placeholder)
		d.Set(ThresholdRuleFieldMetricPatternOperator, ruleSpec.MetricPattern.Operator)
	}
	return nil
}

func mapStateToDataObjectForCustomEventSpecificationWithThresholdRule(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (restapi.InstanaDataObject, error) {
	severity, err := ConvertSeverityFromTerraformToInstanaAPIRepresentation(d.Get(CustomEventSpecificationRuleSeverity).(string))
	if err != nil {
		return &restapi.CustomEventSpecification{}, err
	}
	metricName := d.Get(ThresholdRuleFieldMetricName).(string)
	conditionOperatorString := d.Get(ThresholdRuleFieldConditionOperator).(string)
	conditionOperator, err := restapi.SupportedConditionOperators.FromTerraformValue(conditionOperatorString)
	if err != nil {
		return &restapi.CustomEventSpecification{}, err
	}
	conditionOperatorInstanaValue := conditionOperator.InstanaAPIValue()

	rule := restapi.RuleSpecification{
		DType:             restapi.ThresholdRuleType,
		Severity:          severity,
		MetricName:        &metricName,
		Rollup:            GetIntPointerFromResourceData(d, ThresholdRuleFieldRollup),
		Window:            GetIntPointerFromResourceData(d, ThresholdRuleFieldWindow),
		Aggregation:       getAggregationTypePointerFromResourceData(d, ThresholdRuleFieldAggregation),
		ConditionOperator: &conditionOperatorInstanaValue,
		ConditionValue:    GetFloat64PointerFromResourceData(d, ThresholdRuleFieldConditionValue),
	}

	metricPatternPrefix, ok := d.GetOk(ThresholdRuleFieldMetricPatternPrefix)
	if ok {
		metricPattern := restapi.MetricPattern{
			Prefix:      metricPatternPrefix.(string),
			Postfix:     GetStringPointerFromResourceData(d, ThresholdRuleFieldMetricPatternPostfix),
			Placeholder: GetStringPointerFromResourceData(d, ThresholdRuleFieldMetricPatternPlaceholder),
			Operator:    restapi.MetricPatternOperatorType(d.Get(ThresholdRuleFieldMetricPatternOperator).(string)),
		}
		rule.MetricPattern = &metricPattern
	}

	customEventSpecification := createCustomEventSpecificationFromResourceData(d, formatter)
	customEventSpecification.Rules = []restapi.RuleSpecification{rule}
	return customEventSpecification, nil
}

func getAggregationTypePointerFromResourceData(d *schema.ResourceData, key string) *restapi.AggregationType {
	val, ok := d.GetOk(key)
	if ok {
		value := restapi.AggregationType(val.(string))
		return &value
	}
	return nil
}

func customEventSpecificationWithThresholdRuleSchemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: mergeSchemaMap(defaultCustomEventSchemaFieldsV0, thresholdRuleSchemaFields),
	}
}

func customEventSpecificationWithThresholdRuleSchemaV1() *schema.Resource {
	return &schema.Resource{
		Schema: mergeSchemaMap(defaultCustomEventSchemaFieldsV1, thresholdRuleSchemaFields),
	}
}

func customEventSpecificationWithThresholdRuleSchemaV2() *schema.Resource {
	return &schema.Resource{
		Schema: mergeSchemaMap(defaultCustomEventSchemaFieldsV1, thresholdRuleSchemaFields),
	}
}

func migrateCustomEventConfigWithThreasholdRuleToVersion3ByChangingConditionOperatorToInstanaRepresentation(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	v, ok := rawState[ThresholdRuleFieldConditionOperator]
	if ok {
		operator, err := restapi.SupportedConditionOperators.FromTerraformValue(v.(string))
		if err != nil {
			return rawState, err
		}
		rawState[ThresholdRuleFieldConditionOperator] = operator.InstanaAPIValue()
	}
	return rawState, nil
}
