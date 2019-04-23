package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	//SystemRuleSpecificationSystemRuleID constant value for the schema field system_rule_id of a rule specification
	SystemRuleSpecificationSystemRuleID = ruleFieldPrefix + "system_rule_id"
)

var systemRuleSchemaFields = map[string]*schema.Schema{
	SystemRuleSpecificationSystemRuleID: &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Configures the system rule id for the system rule of the custom event specification",
	},
}

//CreateResourceCustomSystemEventSpecification creates the resource definition for the instana api endpoint for Custom Event Specifications for System rules
func CreateResourceCustomSystemEventSpecification() *schema.Resource {
	return &schema.Resource{
		Read:   createReadCustomSystemEventSpecification(),
		Create: createCreateCustomSystemEventSpecification(),
		Update: createUpdateCustomSystemEventSpecification(),
		Delete: createDeleteCustomSystemEventSpecification(),

		Schema: createCustomEventSpecificationSchema(systemRuleSchemaFields),
	}
}

func createReadCustomSystemEventSpecification() func(*schema.ResourceData, interface{}) error {
	return createCustomEventSpecificationReadFunc(mapCustomSystemEventToTerraformState)
}

func createCreateCustomSystemEventSpecification() func(*schema.ResourceData, interface{}) error {
	return createCustomEventSpecificationCreateFunc(mapCustomSystemEventToInstanaAPIModel, mapCustomSystemEventToTerraformState)
}

func createUpdateCustomSystemEventSpecification() func(*schema.ResourceData, interface{}) error {
	return createCustomEventSpecificationUpdateFunc(mapCustomSystemEventToInstanaAPIModel, mapCustomSystemEventToTerraformState)
}

func createDeleteCustomSystemEventSpecification() func(d *schema.ResourceData, meta interface{}) error {
	return createCustomEventSpecificationDeleteFunc(mapCustomSystemEventToInstanaAPIModel)
}

func mapCustomSystemEventToInstanaAPIModel(d *schema.ResourceData) (restapi.RuleSpecification, error) {
	severity, err := ConvertSeverityFromTerraformToInstanaAPIRepresentation(d.Get(CustomEventSpecificationRuleSeverity).(string))
	if err != nil {
		return restapi.RuleSpecification{}, err
	}
	systemRuleID := d.Get(SystemRuleSpecificationSystemRuleID).(string)
	return restapi.NewSystemRuleSpecification(systemRuleID, severity), nil
}

func mapCustomSystemEventToTerraformState(d *schema.ResourceData, spec restapi.CustomEventSpecification) error {
	ruleSpec := spec.Rules[0]
	severity, err := ConvertSeverityFromInstanaAPIToTerraformRepresentation(ruleSpec.Severity)
	if err != nil {
		return err
	}

	d.Set(CustomEventSpecificationRuleSeverity, severity)
	d.Set(SystemRuleSpecificationSystemRuleID, ruleSpec.SystemRuleID)
	return nil
}
