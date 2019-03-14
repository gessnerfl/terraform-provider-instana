package instana

import (
	"errors"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi/resources"
	"github.com/hashicorp/terraform/helper/schema"
)

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
				Optional: true,
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
	d.SetId(RandomID())
	return resourceRuleUpdate(d, meta)
}

func resourceRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*restapi.RestClient)
	ruleID := d.Id()
	if len(ruleID) == 0 {
		return errors.New("ID of rule is missing")
	}
	rule, err := resources.NewRuleResource(*client).GetOne(ruleID)
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

func resourceRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*restapi.RestClient)
	rule := createRuleFromResourceData(d)
	updatedRule, err := resources.NewRuleResource(*client).Upsert(rule)
	if err != nil {
		return err
	}
	updateRuleState(d, updatedRule)
	return nil
}

func resourceRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*restapi.RestClient)
	rule := createRuleFromResourceData(d)
	err := resources.NewRuleResource(*client).Delete(rule)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func createRuleFromResourceData(d *schema.ResourceData) restapi.Rule {
	return restapi.Rule{
		ID:                d.Id(),
		Name:              d.Get(fieldName).(string),
		EntityType:        d.Get(fieldEntityType).(string),
		MetricName:        d.Get(fieldMetricName).(string),
		Rollup:            d.Get(fieldRollup).(int),
		Window:            d.Get(fieldWindow).(int),
		Aggregation:       d.Get(fieldAggregation).(string),
		ConditionOperator: d.Get(fieldConditionOperator).(string),
		ConditionValue:    d.Get(fieldConditionValue).(float64),
	}
}

func updateRuleState(d *schema.ResourceData, rule restapi.Rule) {
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
