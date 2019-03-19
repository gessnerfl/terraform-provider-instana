package instana

import (
	"errors"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform/helper/schema"
)

//RuleBindingFieldEnabled constant value for the schema field enabled
const RuleBindingFieldEnabled = "enabled"

//RuleBindingFieldTriggering constant value for the schema field triggering
const RuleBindingFieldTriggering = "triggering"

//RuleBindingFieldSeverity constant value for the schema field severity
const RuleBindingFieldSeverity = "severity"

//RuleBindingFieldText constant value for the schema field text
const RuleBindingFieldText = "text"

//RuleBindingFieldDescription constant value for the schema field description
const RuleBindingFieldDescription = "description"

//RuleBindingFieldExpirationTime constant value for the schema field expiration_time
const RuleBindingFieldExpirationTime = "expiration_time"

//RuleBindingFieldQuery constant value for the schema field query
const RuleBindingFieldQuery = "query"

//RuleBindingFieldRuleIds constant value for the schema field rule_ids
const RuleBindingFieldRuleIds = "rule_ids"

//CreateResourceRuleBinding creates the resource definition for the instana api endpoint for Rule Bindings
func CreateResourceRuleBinding() *schema.Resource {
	return &schema.Resource{
		Read:   ReadRuleBinding,
		Create: CreateRuleBinding,
		Update: UpdateRuleBinding,
		Delete: DeleteRuleBinding,

		Schema: map[string]*schema.Schema{
			RuleBindingFieldEnabled: &schema.Schema{
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			RuleBindingFieldTriggering: &schema.Schema{
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			RuleBindingFieldSeverity: &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			RuleBindingFieldText: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			RuleBindingFieldDescription: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			RuleBindingFieldExpirationTime: &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			RuleBindingFieldQuery: &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
			RuleBindingFieldRuleIds: &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

//ReadRuleBinding reads the rule binding with the given id from the Instana API and updates the resource state.
func ReadRuleBinding(d *schema.ResourceData, meta interface{}) error {
	instanaAPI := meta.(restapi.InstanaAPI)
	ruleBindingID := d.Id()
	if len(ruleBindingID) == 0 {
		return errors.New("ID of rule binding is missing")
	}
	ruleBinding, err := instanaAPI.RuleBindings().GetOne(ruleBindingID)
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

//CreateRuleBinding creates the configured rule binding through the Instana API and updates the resource state.
func CreateRuleBinding(d *schema.ResourceData, meta interface{}) error {
	d.SetId(RandomID())
	return UpdateRuleBinding(d, meta)
}

//UpdateRuleBinding updates the configured rule binding through the Instana API and updates the resource state.
func UpdateRuleBinding(d *schema.ResourceData, meta interface{}) error {
	instanaAPI := meta.(restapi.InstanaAPI)
	ruleBinding := createRuleBindingFromResourceData(d)
	updatedRuleBinding, err := instanaAPI.RuleBindings().Upsert(ruleBinding)
	if err != nil {
		return err
	}
	updateRuleBindingState(d, updatedRuleBinding)
	return nil
}

//DeleteRuleBinding deletes the configured rule binding through the Instana API and deletes the resource state.
func DeleteRuleBinding(d *schema.ResourceData, meta interface{}) error {
	instanaAPI := meta.(restapi.InstanaAPI)
	ruleBinding := createRuleBindingFromResourceData(d)
	err := instanaAPI.RuleBindings().DeleteByID(ruleBinding.ID)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func createRuleBindingFromResourceData(d *schema.ResourceData) restapi.RuleBinding {
	return restapi.RuleBinding{
		ID:             d.Id(),
		Enabled:        d.Get(RuleBindingFieldEnabled).(bool),
		Triggering:     d.Get(RuleBindingFieldTriggering).(bool),
		Severity:       d.Get(RuleBindingFieldSeverity).(int),
		Text:           d.Get(RuleBindingFieldText).(string),
		Description:    d.Get(RuleBindingFieldDescription).(string),
		ExpirationTime: d.Get(RuleBindingFieldExpirationTime).(int),
		Query:          d.Get(RuleBindingFieldQuery).(string),
		RuleIds:        ReadStringArrayParameterFromResource(d, RuleBindingFieldRuleIds),
	}
}

func updateRuleBindingState(d *schema.ResourceData, ruleBinding restapi.RuleBinding) {
	d.Set(RuleBindingFieldEnabled, ruleBinding.Enabled)
	d.Set(RuleBindingFieldTriggering, ruleBinding.Triggering)
	d.Set(RuleBindingFieldSeverity, ruleBinding.Severity)
	d.Set(RuleBindingFieldText, ruleBinding.Text)
	d.Set(RuleBindingFieldDescription, ruleBinding.Description)
	d.Set(RuleBindingFieldExpirationTime, ruleBinding.ExpirationTime)
	d.Set(RuleBindingFieldQuery, ruleBinding.Query)
	d.Set(RuleBindingFieldRuleIds, ruleBinding.RuleIds)

	d.SetId(ruleBinding.ID)
}
