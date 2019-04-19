package instana

import (
	"errors"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

const (
	//RuleBindingFieldEnabled constant value for the schema field enabled
	RuleBindingFieldEnabled = "enabled"
	//RuleBindingFieldTriggering constant value for the schema field triggering
	RuleBindingFieldTriggering = "triggering"
	//RuleBindingFieldSeverity constant value for the schema field severity
	RuleBindingFieldSeverity = "severity"
	//RuleBindingFieldText constant value for the schema field text
	RuleBindingFieldText = "text"
	//RuleBindingFieldDescription constant value for the schema field description
	RuleBindingFieldDescription = "description"
	//RuleBindingFieldExpirationTime constant value for the schema field expiration_time
	RuleBindingFieldExpirationTime = "expiration_time"
	//RuleBindingFieldQuery constant value for the schema field query
	RuleBindingFieldQuery = "query"
	//RuleBindingFieldRuleIds constant value for the schema field rule_ids
	RuleBindingFieldRuleIds = "rule_ids"
)

//CreateResourceRuleBinding creates the resource definition for the instana api endpoint for Rule Bindings
func CreateResourceRuleBinding() *schema.Resource {
	return &schema.Resource{
		Read:   ReadRuleBinding,
		Create: CreateRuleBinding,
		Update: UpdateRuleBinding,
		Delete: DeleteRuleBinding,

		Schema: map[string]*schema.Schema{
			RuleBindingFieldEnabled: &schema.Schema{
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
				Description: "Configures if the rule binding is enabled or not",
			},
			RuleBindingFieldTriggering: &schema.Schema{
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
				Description: "Configures the issue should trigger an incident",
			},
			RuleBindingFieldSeverity: &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{SeverityWarning.terraformRepresentation, SeverityCritical.terraformRepresentation}, false),
				Description:  "Configures the severity of the issue",
			},
			RuleBindingFieldText: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Configures the title of the rule binding",
			},
			RuleBindingFieldDescription: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Configures the description text of the rule binding",
			},
			RuleBindingFieldExpirationTime: &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Configures the expiration time (grace period) to wait before the issue is closed",
			},
			RuleBindingFieldQuery: &schema.Schema{
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "Configures the dynamic focus query for the rule binding",
			},
			RuleBindingFieldRuleIds: &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Configures the list of rule ids which should be considered by the rule binding",
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
	return updateRuleBindingState(d, ruleBinding)
}

//CreateRuleBinding creates the configured rule binding through the Instana API and updates the resource state.
func CreateRuleBinding(d *schema.ResourceData, meta interface{}) error {
	d.SetId(RandomID())
	return UpdateRuleBinding(d, meta)
}

//UpdateRuleBinding updates the configured rule binding through the Instana API and updates the resource state.
func UpdateRuleBinding(d *schema.ResourceData, meta interface{}) error {
	instanaAPI := meta.(restapi.InstanaAPI)
	ruleBinding, err := createRuleBindingFromResourceData(d)
	if err != nil {
		return err
	}
	updatedRuleBinding, err := instanaAPI.RuleBindings().Upsert(ruleBinding)
	if err != nil {
		return err
	}
	return updateRuleBindingState(d, updatedRuleBinding)
}

//DeleteRuleBinding deletes the configured rule binding through the Instana API and deletes the resource state.
func DeleteRuleBinding(d *schema.ResourceData, meta interface{}) error {
	instanaAPI := meta.(restapi.InstanaAPI)
	ruleBinding, err := createRuleBindingFromResourceData(d)
	if err != nil {
		return err
	}
	err = instanaAPI.RuleBindings().DeleteByID(ruleBinding.ID)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func createRuleBindingFromResourceData(d *schema.ResourceData) (restapi.RuleBinding, error) {
	severity, err := ConvertSeverityFromTerraformToInstanaAPIRepresentation(d.Get(RuleBindingFieldSeverity).(string))
	if err != nil {
		return restapi.RuleBinding{}, err
	}

	return restapi.RuleBinding{
		ID:             d.Id(),
		Enabled:        d.Get(RuleBindingFieldEnabled).(bool),
		Triggering:     d.Get(RuleBindingFieldTriggering).(bool),
		Severity:       severity,
		Text:           d.Get(RuleBindingFieldText).(string),
		Description:    d.Get(RuleBindingFieldDescription).(string),
		ExpirationTime: d.Get(RuleBindingFieldExpirationTime).(int),
		Query:          d.Get(RuleBindingFieldQuery).(string),
		RuleIds:        ReadStringArrayParameterFromResource(d, RuleBindingFieldRuleIds),
	}, nil
}

func updateRuleBindingState(d *schema.ResourceData, ruleBinding restapi.RuleBinding) error {
	severity, err := ConvertSeverityFromInstanaAPIToTerraformRepresentation(ruleBinding.Severity)
	if err != nil {
		return err
	}

	d.Set(RuleBindingFieldEnabled, ruleBinding.Enabled)
	d.Set(RuleBindingFieldTriggering, ruleBinding.Triggering)
	d.Set(RuleBindingFieldSeverity, severity)
	d.Set(RuleBindingFieldText, ruleBinding.Text)
	d.Set(RuleBindingFieldDescription, ruleBinding.Description)
	d.Set(RuleBindingFieldExpirationTime, ruleBinding.ExpirationTime)
	d.Set(RuleBindingFieldQuery, ruleBinding.Query)
	d.Set(RuleBindingFieldRuleIds, ruleBinding.RuleIds)

	d.SetId(ruleBinding.ID)
	return nil
}
