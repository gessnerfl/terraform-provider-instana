package instana

import (
	"errors"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

const (
	//RuleFieldName constant value for the schema field name
	RuleFieldName = "name"
	//RuleFieldEntityType constant value for the schema field entity_type
	RuleFieldEntityType = "entity_type"
	//RuleFieldMetricName constant value for the schema field metric_name
	RuleFieldMetricName = "metric_name"
	//RuleFieldRollup constant value for the schema field rollup
	RuleFieldRollup = "rollup"
	//RuleFieldWindow constant value for the schema field window
	RuleFieldWindow = "window"
	//RuleFieldAggregation constant value for the schema field aggregation
	RuleFieldAggregation = "aggregation"
	//RuleFieldConditionOperator constant value for the schema field condition_operator
	RuleFieldConditionOperator = "condition_operator"
	//RuleFieldConditionValue constant value for the schema field condition_value
	RuleFieldConditionValue = "condition_value"
)

//CreateResourceRule creates the resource definition for the resource instana_rule
func CreateResourceRule() *schema.Resource {
	return &schema.Resource{
		Create: CreateRule,
		Read:   ReadRule,
		Update: UpdateRule,
		Delete: DeleteRule,

		Schema: map[string]*schema.Schema{
			RuleFieldName: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the rule",
			},
			RuleFieldEntityType: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The entity type of the rule",
			},
			RuleFieldMetricName: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The metric name of the rule",
			},
			RuleFieldRollup: &schema.Schema{
				Type:        schema.TypeInt,
				Required:    false,
				Optional:    true,
				Description: "The rollup of the metric",
			},
			RuleFieldWindow: &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The time window where the condition has to be fulfilled",
			},
			RuleFieldAggregation: &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(restapi.SupportedAggregationTypes.ToStringSlice(), false),
				Description:  "The aggregation type (e.g. sum, avg)",
			},
			RuleFieldConditionOperator: &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(restapi.SupportedConditionOperatorTypes.ToStringSlice(), false),
				Description:  "The aggregation operator (e.g >, <)",
			},
			RuleFieldConditionValue: &schema.Schema{
				Type:        schema.TypeFloat,
				Required:    true,
				Description: "The expected aggregation value to fulfill the rule",
			},
		},
	}
}

//CreateRule defines the create operation for the resource instana_rule
func CreateRule(d *schema.ResourceData, meta interface{}) error {
	d.SetId(RandomID())
	return UpdateRule(d, meta)
}

//ReadRule defines the read operation for the resource instana_rule
func ReadRule(d *schema.ResourceData, meta interface{}) error {
	instanaAPI := meta.(restapi.InstanaAPI)
	ruleID := d.Id()
	if len(ruleID) == 0 {
		return errors.New("ID of rule is missing")
	}
	rule, err := instanaAPI.Rules().GetOne(ruleID)
	if err != nil {
		if err == restapi.ErrEntityNotFound {
			d.SetId("")
			return nil
		}
		return err
	}
	updateRuleState(d, rule)
	return nil
}

//UpdateRule defines the update operation for the resource instana_rule
func UpdateRule(d *schema.ResourceData, meta interface{}) error {
	instanaAPI := meta.(restapi.InstanaAPI)
	rule := createRuleFromResourceData(d)
	updatedRule, err := instanaAPI.Rules().Upsert(rule)
	if err != nil {
		return err
	}
	updateRuleState(d, updatedRule)
	return nil
}

//DeleteRule defines the delete operation for the resource instana_rule
func DeleteRule(d *schema.ResourceData, meta interface{}) error {
	instanaAPI := meta.(restapi.InstanaAPI)
	rule := createRuleFromResourceData(d)
	err := instanaAPI.Rules().DeleteByID(rule.ID)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func createRuleFromResourceData(d *schema.ResourceData) restapi.Rule {
	return restapi.Rule{
		ID:                d.Id(),
		Name:              d.Get(RuleFieldName).(string),
		EntityType:        d.Get(RuleFieldEntityType).(string),
		MetricName:        d.Get(RuleFieldMetricName).(string),
		Rollup:            d.Get(RuleFieldRollup).(int),
		Window:            d.Get(RuleFieldWindow).(int),
		Aggregation:       restapi.AggregationType(d.Get(RuleFieldAggregation).(string)),
		ConditionOperator: restapi.ConditionOperatorType(d.Get(RuleFieldConditionOperator).(string)),
		ConditionValue:    d.Get(RuleFieldConditionValue).(float64),
	}
}

func updateRuleState(d *schema.ResourceData, rule restapi.Rule) {
	d.Set(RuleFieldName, rule.Name)
	d.Set(RuleFieldEntityType, rule.EntityType)
	d.Set(RuleFieldMetricName, rule.MetricName)
	d.Set(RuleFieldRollup, rule.Rollup)
	d.Set(RuleFieldWindow, rule.Window)
	d.Set(RuleFieldAggregation, string(rule.Aggregation))
	d.Set(RuleFieldConditionOperator, string(rule.ConditionOperator))
	d.Set(RuleFieldConditionValue, rule.ConditionValue)

	d.SetId(rule.ID)
}
