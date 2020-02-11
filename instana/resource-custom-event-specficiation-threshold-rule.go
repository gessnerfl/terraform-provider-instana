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
	CustomEventSpecificationFieldEntityType: &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Configures the entity type of the custom event specification",
	},
	ThresholdRuleFieldMetricName: &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The metric name of the rule",
	},
	ThresholdRuleFieldRollup: &schema.Schema{
		Type:        schema.TypeInt,
		Required:    false,
		Optional:    true,
		Description: "The rollup of the metric",
	},
	ThresholdRuleFieldWindow: &schema.Schema{
		Type:        schema.TypeInt,
		Required:    false,
		Optional:    true,
		Description: "The time window where the condition has to be fulfilled",
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

//NewCustomEventSpecificationWithThresholdRuleResourceHandle creates a new ResourceHandle for the terraform resource of custom event specifications with system rules
func NewCustomEventSpecificationWithThresholdRuleResourceHandle() *ResourceHandle {
	return &ResourceHandle{
		ResourceName:  ResourceInstanaCustomEventSpecificationThresholdRule,
		Schema:        mergeSchemaMap(defaultCustomEventSchemaFields, thresholdRuleSchemaFields),
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    customEventSpecificationWithThresholdRuleSchemaV0().CoreConfigSchema().ImpliedType(),
				Upgrade: migrateCustomEventConfigFullNameInStateFromV0toV1,
				Version: 0,
			},
		},
		RestResourceFactory:  func(api restapi.InstanaAPI) restapi.RestResource { return api.CustomEventSpecifications() },
		UpdateState:          updateStateForCustomEventSpecificationWithThresholdRule,
		MapStateToDataObject: mapStateToDataObjectForCustomEventSpecificationWithThresholdRule,
	}
}

func updateStateForCustomEventSpecificationWithThresholdRule(d *schema.ResourceData, obj restapi.InstanaDataObject) error {
	customEventSpecification := obj.(restapi.CustomEventSpecification)
	return updateStateForBasicCustomEventSpecification(d, customEventSpecification, mapThresholdRuleToTerraformState)
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

func mapStateToDataObjectForCustomEventSpecificationWithThresholdRule(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (restapi.InstanaDataObject, error) {
	return createCustomEventSpecificationFromResourceData(d, formatter, mapThresholdRuleToInstanaAPIModel)
}

func mapThresholdRuleToInstanaAPIModel(d *schema.ResourceData) (restapi.RuleSpecification, error) {
	severity, err := ConvertSeverityFromTerraformToInstanaAPIRepresentation(d.Get(CustomEventSpecificationRuleSeverity).(string))
	if err != nil {
		return restapi.RuleSpecification{}, err
	}
	metricName := d.Get(ThresholdRuleFieldMetricName).(string)
	conditionOperator := restapi.ConditionOperatorType(d.Get(ThresholdRuleFieldConditionOperator).(string))
	return restapi.RuleSpecification{
		DType:             restapi.ThresholdRuleType,
		Severity:          severity,
		MetricName:        &metricName,
		Rollup:            GetIntPointerFromResourceData(d, ThresholdRuleFieldRollup),
		Window:            GetIntPointerFromResourceData(d, ThresholdRuleFieldWindow),
		Aggregation:       getAggregationTypePointerFromResourceData(d, ThresholdRuleFieldAggregation),
		ConditionOperator: &conditionOperator,
		ConditionValue:    GetFloat64PointerFromResourceData(d, ThresholdRuleFieldConditionValue),
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

func customEventSpecificationWithThresholdRuleSchemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: mergeSchemaMap(defaultCustomEventSchemaFieldsV0, thresholdRuleSchemaFields),
	}
}
