package instana

import (
	"errors"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform/helper/schema"
)

//FieldName constant value for the schema field name
const FieldName = "name"

//FieldEntityType constant value for the schema field entity_type
const FieldEntityType = "entity_type"

//FieldMetricName constant value for the schema field metric_name
const FieldMetricName = "metric_name"

//FieldRollup constant value for the schema field rollup
const FieldRollup = "rollup"

//FieldWindow constant value for the schema field window
const FieldWindow = "window"

//FieldAggregation constant value for the schema field aggregation
const FieldAggregation = "aggregation"

//FieldConditionOperator constant value for the schema field condition_operator
const FieldConditionOperator = "condition_operator"

//FieldConditionValue constant value for the schema field condition_value
const FieldConditionValue = "condition_value"

//CreateResourceRule creates the resource definition for the resource instana_rule
func CreateResourceRule() *schema.Resource {
	return &schema.Resource{
		Create: ResourceRuleCreate,
		Read:   ResourceRuleRead,
		Update: ResourceRuleUpdate,
		Delete: ResourceRuleDelete,

		Schema: map[string]*schema.Schema{
			FieldName: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the rule",
			},
			FieldEntityType: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The entity type of the rule",
			},
			FieldMetricName: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The metric name of the rult",
			},
			FieldRollup: &schema.Schema{
				Type:        schema.TypeInt,
				Required:    false,
				Optional:    true,
				Description: "The rollup of the metric",
			},
			FieldWindow: &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The time window where the condition has to be fulfilled",
			},
			FieldAggregation: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The aggregation type (e.g. sum, avg)",
			},
			FieldConditionOperator: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The aggregation operator (e.g >, <)",
			},
			FieldConditionValue: &schema.Schema{
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
		Name:              d.Get(FieldName).(string),
		EntityType:        d.Get(FieldEntityType).(string),
		MetricName:        d.Get(FieldMetricName).(string),
		Rollup:            d.Get(FieldRollup).(int),
		Window:            d.Get(FieldWindow).(int),
		Aggregation:       d.Get(FieldAggregation).(string),
		ConditionOperator: d.Get(FieldConditionOperator).(string),
		ConditionValue:    d.Get(FieldConditionValue).(float64),
	}
}

func updateRuleState(d *schema.ResourceData, rule restapi.Rule) {
	d.Set(FieldName, rule.Name)
	d.Set(FieldEntityType, rule.EntityType)
	d.Set(FieldMetricName, rule.MetricName)
	d.Set(FieldRollup, rule.Rollup)
	d.Set(FieldWindow, rule.Window)
	d.Set(FieldAggregation, rule.Aggregation)
	d.Set(FieldConditionOperator, rule.ConditionOperator)
	d.Set(FieldConditionValue, rule.ConditionValue)

	d.SetId(rule.ID)
}
