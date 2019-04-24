package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

const (
	//ThresholdRuleFieldMetricName constant value for the schema field metric_name
	ThresholdRuleFieldMetricName = ruleFieldPrefix + "metric_name"
	//ThresholdRuleFieldRollup constant value for the schema field rollup
	ThresholdRuleFieldRollup = ruleFieldPrefix + "rollup"
	//ThresholdRuleFieldWindow constant value for the schema field window
	ThresholdRuleFieldWindow = ruleFieldPrefix + "window"
	//ThresholdRuleFieldAggregation constant value for the schema field aggregation
	ThresholdRuleFieldAggregation = ruleFieldPrefix + "aggregation"
	//ThresholdRuleFieldConditionOperator constant value for the schema field condition_operator
	ThresholdRuleFieldConditionOperator = ruleFieldPrefix + "condition_operator"
	//ThresholdRuleFieldConditionValue constant value for the schema field condition_value
	ThresholdRuleFieldConditionValue = ruleFieldPrefix + "condition_value"
)

var thresholdRuleSchemaFields = map[string]*schema.Schema{
	ThresholdRuleFieldMetricName: &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The metric name of the rule",
	},
	ThresholdRuleFieldRollup: &schema.Schema{
		Type:          schema.TypeInt,
		Required:      false,
		Optional:      true,
		ConflictsWith: []string{ThresholdRuleFieldWindow},
		Description:   "The rollup of the metric",
	},
	ThresholdRuleFieldWindow: &schema.Schema{
		Type:          schema.TypeInt,
		Required:      false,
		Optional:      true,
		ConflictsWith: []string{ThresholdRuleFieldRollup},
		Description:   "The time window where the condition has to be fulfilled",
	},
	ThresholdRuleFieldAggregation: &schema.Schema{
		Type:         schema.TypeString,
		Required:     false,
		Optional:     true,
		ValidateFunc: validation.StringInSlice(restapi.SupportedAggregationTypes.ToStringSlice(), false),
		Description:  "The aggregation type (e.g. sum, avg)",
	},
	ThresholdRuleFieldConditionOperator: &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice(restapi.SupportedConditionOperatorTypes.ToStringSlice(), false),
		Description:  "The condition operator (e.g >, <)",
	},
	ThresholdRuleFieldConditionValue: &schema.Schema{
		Type:        schema.TypeFloat,
		Required:    true,
		Description: "The expected condition value to fulfill the rule",
	},
}

//CreateResourceCustomEventSpecificationWithThresholdRule creates the resource definition for the instana api endpoint for Custom Event Specifications for Threshold rules
func CreateResourceCustomEventSpecificationWithThresholdRule() *schema.Resource {
	return &schema.Resource{
		Read:   createReadCustomEventSpecificationWithThresholdRule(),
		Create: createCreateCustomEventSpecificationWithThresholdRule(),
		Update: createUpdateCustomEventSpecificationWithThresholdRule(),
		Delete: createDeleteCustomEventSpecificationWithThresholdRule(),

		Schema: createCustomEventSpecificationSchema(thresholdRuleSchemaFields),
	}
}

func createReadCustomEventSpecificationWithThresholdRule() func(*schema.ResourceData, interface{}) error {
	return createCustomEventSpecificationReadFunc(mapThresholdRuleToTerraformState)
}

func createCreateCustomEventSpecificationWithThresholdRule() func(*schema.ResourceData, interface{}) error {
	return createCustomEventSpecificationCreateFunc(mapThresholdRuleToInstanaAPIModel, mapThresholdRuleToTerraformState)
}

func createUpdateCustomEventSpecificationWithThresholdRule() func(*schema.ResourceData, interface{}) error {
	return createCustomEventSpecificationUpdateFunc(mapThresholdRuleToInstanaAPIModel, mapThresholdRuleToTerraformState)
}

func createDeleteCustomEventSpecificationWithThresholdRule() func(*schema.ResourceData, interface{}) error {
	return createCustomEventSpecificationDeleteFunc(mapThresholdRuleToInstanaAPIModel)
}

func mapThresholdRuleToInstanaAPIModel(d *schema.ResourceData) (restapi.RuleSpecification, error) {
	severity, err := ConvertSeverityFromTerraformToInstanaAPIRepresentation(d.Get(CustomEventSpecificationRuleSeverity).(string))
	if err != nil {
		return restapi.RuleSpecification{}, err
	}
	return restapi.RuleSpecification{
		DType:                              restapi.ThresholdRuleType,
		Severity:                           severity,
		MetricName:                         d.Get(ThresholdRuleFieldMetricName).(string),
		Rollup:                             GetIntPointerFromResourceData(d, ThresholdRuleFieldRollup),
		Window:                             GetIntPointerFromResourceData(d, ThresholdRuleFieldWindow),
		Aggregation:                        getAggregationTypePointerFromResourceData(d, ThresholdRuleFieldAggregation),
		ConditionOperator:                  restapi.ConditionOperatorType(d.Get(ThresholdRuleFieldConditionOperator).(string)),
		ConditionValue:                     GetFloat64PointerFromResourceData(d, ThresholdRuleFieldConditionValue),
		AggregationForNonPercentileMetric:  true,
		EitherRollupOrWindowAndAggregation: true,
	}, nil
}

func getAggregationTypePointerFromResourceData(d *schema.ResourceData, key string) *restapi.AggregationType {
	val, ok := d.GetOk(key)
	if ok {
		value := restapi.AggregationType(val.(string))
		return &value
	}
	return nil
}

func mapThresholdRuleToTerraformState(d *schema.ResourceData, spec restapi.CustomEventSpecification) error {
	ruleSpec := spec.Rules[0]
	severity, err := ConvertSeverityFromInstanaAPIToTerraformRepresentation(ruleSpec.Severity)
	if err != nil {
		return err
	}

	d.Set(CustomEventSpecificationRuleSeverity, severity)
	d.Set(ThresholdRuleFieldMetricName, ruleSpec.MetricName)
	d.Set(ThresholdRuleFieldRollup, ruleSpec.Rollup)
	d.Set(ThresholdRuleFieldWindow, ruleSpec.Window)
	d.Set(ThresholdRuleFieldAggregation, ruleSpec.Aggregation)
	d.Set(ThresholdRuleFieldConditionOperator, ruleSpec.ConditionOperator)
	d.Set(ThresholdRuleFieldConditionValue, ruleSpec.ConditionValue)
	return nil
}
