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

//CreateResourceCustomEventSpecificationWithSystemRule creates the resource definition for the instana api endpoint for Custom Event Specifications for System rules
func CreateResourceCustomEventSpecificationWithSystemRule() *schema.Resource {
	return &schema.Resource{
		Read:   createReadCustomEventSpecificationWithSystemRule(),
		Create: createCreateCustomEventSpecificationWithSystemRule(),
		Update: createUpdateCustomEventSpecificationWithSystemRule(),
		Delete: createDeleteCustomEventSpecificationWithSystemRule(),

		Schema: createCustomEventSpecificationSchema(systemRuleSchemaFields),
	}
}

func createReadCustomEventSpecificationWithSystemRule() func(*schema.ResourceData, interface{}) error {
	return createCustomEventSpecificationReadFunc(mapSystemRuleToTerraformState)
}

func createCreateCustomEventSpecificationWithSystemRule() func(*schema.ResourceData, interface{}) error {
	return createCustomEventSpecificationCreateFunc(mapSystemRuleToInstanaAPIModel, mapSystemRuleToTerraformState)
}

func createUpdateCustomEventSpecificationWithSystemRule() func(*schema.ResourceData, interface{}) error {
	return createCustomEventSpecificationUpdateFunc(mapSystemRuleToInstanaAPIModel, mapSystemRuleToTerraformState)
}

func createDeleteCustomEventSpecificationWithSystemRule() func(*schema.ResourceData, interface{}) error {
	return createCustomEventSpecificationDeleteFunc(mapSystemRuleToInstanaAPIModel)
}

func mapSystemRuleToInstanaAPIModel(d *schema.ResourceData) (restapi.RuleSpecification, error) {
	severity, err := ConvertSeverityFromTerraformToInstanaAPIRepresentation(d.Get(CustomEventSpecificationRuleSeverity).(string))
	if err != nil {
		return restapi.RuleSpecification{}, err
	}
	systemRuleID := d.Get(SystemRuleSpecificationSystemRuleID).(string)
	return restapi.NewSystemRuleSpecification(systemRuleID, severity), nil
}

func mapSystemRuleToTerraformState(d *schema.ResourceData, spec restapi.CustomEventSpecification) error {
	ruleSpec := spec.Rules[0]
	severity, err := ConvertSeverityFromInstanaAPIToTerraformRepresentation(ruleSpec.Severity)
	if err != nil {
		return err
	}

	d.Set(CustomEventSpecificationRuleSeverity, severity)
	d.Set(SystemRuleSpecificationSystemRuleID, ruleSpec.SystemRuleID)
	return nil
}
