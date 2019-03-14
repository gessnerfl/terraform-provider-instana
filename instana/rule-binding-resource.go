package instana

import (
	"errors"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi/resources"
	"github.com/hashicorp/terraform/helper/schema"
)

const fieldEnabled = "enabled"
const fieldTriggering = "triggering"
const fieldSeverity = "severity"
const fieldText = "text"
const fieldDescription = "description"
const fieldExpirationTime = "expiration_time"
const fieldQuery = "query"
const fieldRuleIds = "rule_ids"

func createResourceRuleBinding() *schema.Resource {
	return &schema.Resource{
		Create: resourceRuleBindingCreate,
		Read:   resourceRuleBindingRead,
		Update: resourceRuleBindingUpdate,
		Delete: resourceRuleBindingDelete,

		Schema: map[string]*schema.Schema{
			fieldEnabled: &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
				Default:  true,
			},
			fieldTriggering: &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
				Default:  false,
			},
			fieldSeverity: &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			fieldText: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			fieldDescription: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			fieldExpirationTime: &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			fieldQuery: &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			fieldRuleIds: &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceRuleBindingCreate(d *schema.ResourceData, meta interface{}) error {
	d.SetId(RandomID())
	return resourceRuleBindingUpdate(d, meta)
}

func resourceRuleBindingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*restapi.RestClient)
	ruleBindingID := d.Id()
	if len(ruleBindingID) == 0 {
		return errors.New("ID of rule binding is missing")
	}
	ruleBinding, err := resources.NewRuleBindingResource(*client).GetOne(ruleBindingID)
	if err != nil {
		if err == restapi.ErrEntityNotFound {
			d.SetId("")
			return nil
		}
		return err
	}
	updateRuleBindingState(d, ruleBinding)
	return nil
}

func resourceRuleBindingUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*restapi.RestClient)
	ruleBinding := createRuleBindingFromResourceData(d)
	updatedRuleBinding, err := resources.NewRuleBindingResource(*client).Upsert(ruleBinding)
	if err != nil {
		return err
	}
	updateRuleBindingState(d, updatedRuleBinding)
	return nil
}

func resourceRuleBindingDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*restapi.RestClient)
	ruleBinding := createRuleBindingFromResourceData(d)
	err := resources.NewRuleBindingResource(*client).Delete(ruleBinding)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func createRuleBindingFromResourceData(d *schema.ResourceData) restapi.RuleBinding {
	return restapi.RuleBinding{
		ID:             d.Id(),
		Enabled:        d.Get(fieldEnabled).(bool),
		Triggering:     d.Get(fieldTriggering).(bool),
		Severity:       d.Get(fieldSeverity).(int),
		Text:           d.Get(fieldText).(string),
		Description:    d.Get(fieldDescription).(string),
		ExpirationTime: d.Get(fieldExpirationTime).(int),
		Query:          d.Get(fieldQuery).(string),
		RuleIds:        d.Get(fieldRuleIds).([]string),
	}
}

func updateRuleBindingState(d *schema.ResourceData, ruleBinding restapi.RuleBinding) {
	d.Set(fieldEnabled, ruleBinding.Enabled)
	d.Set(fieldTriggering, ruleBinding.Triggering)
	d.Set(fieldSeverity, ruleBinding.Severity)
	d.Set(fieldText, ruleBinding.Text)
	d.Set(fieldDescription, ruleBinding.Description)
	d.Set(fieldExpirationTime, ruleBinding.ExpirationTime)
	d.Set(fieldQuery, ruleBinding.Query)
	d.Set(fieldRuleIds, ruleBinding.RuleIds)

	d.SetId(ruleBinding.ID)
}
