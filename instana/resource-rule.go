package instana

import (
	"errors"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform/helper/schema"
)

//RuleFieldName constant value for the schema field name
const RuleFieldName = "name"

//RuleFieldEntityType constant value for the schema field entity_type
const RuleFieldEntityType = "entity_type"

//RuleFieldMetricName constant value for the schema field metric_name
const RuleFieldMetricName = "metric_name"

//RuleFieldRollup constant value for the schema field rollup
const RuleFieldRollup = "rollup"

//RuleFieldWindow constant value for the schema field window
const RuleFieldWindow = "window"

//RuleFieldAggregation constant value for the schema field aggregation
const RuleFieldAggregation = "aggregation"

//RuleFieldConditionOperator constant value for the schema field condition_operator
const RuleFieldConditionOperator = "condition_operator"

//RuleFieldConditionValue constant value for the schema field condition_value
const RuleFieldConditionValue = "condition_value"

//CreateResourceRule creates the resource definition for the resource instana_rule
func CreateResourceRule() *schema.Resource {
	return &schema.Resource{
		Create: ResourceRuleCreate,
		Read:   ResourceRuleRead,
		Update: ResourceRuleUpdate,
		Delete: ResourceRuleDelete,

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
				Description: "The metric name of the rult",
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
				Type:        schema.TypeString,
				Required:    true,
				Description: "The aggregation type (e.g. sum, avg)",
			},
			RuleFieldConditionOperator: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The aggregation operator (e.g >, <)",
			},
			RuleFieldConditionValue: &schema.Schema{
				Type:        schema.TypeFloat,
				Required:    true,
				Description: "The expected aggregation value to fulfill the rule",
			},
		},
	}
}

//ResourceRuleCreate defines the create operation for the resource instana_rule
func ResourceRuleCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId(RandomID())
	return ResourceRuleUpdate(d, meta)
}

//ResourceRuleRead defines the read operation for the resource instana_rule
func ResourceRuleRead(d *schema.ResourceData, meta interface{}) error {
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

//ResourceRuleUpdate defines the update operation for the resource instana_rule
func ResourceRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	instanaAPI := meta.(restapi.InstanaAPI)
	rule := createRuleFromResourceData(d)
	updatedRule, err := instanaAPI.Rules().Upsert(rule)
	if err != nil {
		return err
	}
	updateRuleState(d, updatedRule)
	return nil
}

//ResourceRuleDelete defines the delete operation for the resource instana_rule
func ResourceRuleDelete(d *schema.ResourceData, meta interface{}) error {
	instanaAPI := meta.(restapi.InstanaAPI)
	rule := createRuleFromResourceData(d)
	err := instanaAPI.Rules().Delete(rule)
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
		Aggregation:       d.Get(RuleFieldAggregation).(string),
		ConditionOperator: d.Get(RuleFieldConditionOperator).(string),
		ConditionValue:    d.Get(RuleFieldConditionValue).(float64),
	}
}

func updateRuleState(d *schema.ResourceData, rule restapi.Rule) {
	d.Set(RuleFieldName, rule.Name)
	d.Set(RuleFieldEntityType, rule.EntityType)
	d.Set(RuleFieldMetricName, rule.MetricName)
	d.Set(RuleFieldRollup, rule.Rollup)
	d.Set(RuleFieldWindow, rule.Window)
	d.Set(RuleFieldAggregation, rule.Aggregation)
	d.Set(RuleFieldConditionOperator, rule.ConditionOperator)
	d.Set(RuleFieldConditionValue, rule.ConditionValue)

	d.SetId(rule.ID)
}
