package instana

import (
	"errors"

	"github.com/gessnerfl/terraform-provider-instana/instana/api"
	"github.com/hashicorp/terraform/helper/schema"
)

const fieldID = "identifier"
const fieldName = "name"
const fieldEntityType = "entity_type"
const fieldMetricName = "metric_name"
const fieldRollup = "rollup"
const fieldWindow = "window"
const fieldAggregation = "aggregation"
const fieldConditionOperator = "condition_operator"
const fieldConditionValue = "condition_value"

func createResourceRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceRuleCreate,
		Read:   resourceRuleRead,
		Update: resourceRuleUpdate,
		Delete: resourceRuleDelete,

		Schema: map[string]*schema.Schema{
			fieldID: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			fieldName: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			fieldEntityType: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			fieldMetricName: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			fieldRollup: &schema.Schema{
				Type:     schema.TypeInt,
				Required: false,
			},
			fieldWindow: &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			fieldAggregation: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			fieldConditionOperator: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			fieldConditionValue: &schema.Schema{
				Type:     schema.TypeFloat,
				Required: true,
			},
		},
	}
}

func resourceRuleCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceRuleUpdate(d, meta)
}

func resourceRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.RestClient)
	ruleID := d.Id()
	if len(ruleID) == 0 {
		return errors.New("ID of rule is missing")
	}
	rule, err := api.NewRuleAPI(*client).GetOne(ruleID)
	if err != nil {
		return err
	}
	updateRuleState(d, rule)
	return nil
}

func resourceRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.RestClient)
	rule := createRuleFromResourceData(d)
	updatedRule, err := api.NewRuleAPI(*client).Upsert(rule)
	if err != nil {
		return err
	}
	updateRuleState(d, updatedRule)
	return nil
}

func resourceRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*api.RestClient)
	rule := createRuleFromResourceData(d)
	err := api.NewRuleAPI(*client).Delete(rule)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func createRuleFromResourceData(d *schema.ResourceData) api.Rule {
	return api.Rule{
		ID:                d.Get(fieldID).(string),
		Name:              d.Get(fieldName).(string),
		EntityType:        d.Get(fieldEntityType).(string),
		MetricName:        d.Get(fieldMetricName).(string),
		Rollup:            d.Get(fieldRollup).(int),
		Window:            d.Get(fieldWindow).(int),
		Aggregation:       d.Get(fieldAggregation).(string),
		ConditionOperator: d.Get(fieldConditionOperator).(string),
		ConditionValue:    d.Get(fieldConditionValue).(float32),
	}
}

func updateRuleState(d *schema.ResourceData, rule api.Rule) {
	d.Set(fieldID, rule.ID)
	d.Set(fieldName, rule.Name)
	d.Set(fieldEntityType, rule.EntityType)
	d.Set(fieldMetricName, rule.MetricName)
	d.Set(fieldRollup, rule.Rollup)
	d.Set(fieldWindow, rule.Window)
	d.Set(fieldAggregation, rule.Aggregation)
	d.Set(fieldConditionOperator, rule.ConditionOperator)
	d.Set(fieldConditionValue, rule.ConditionValue)

	d.SetId(rule.ID)
}
