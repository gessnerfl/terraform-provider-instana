package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

const (
	//SystemRuleSpecificationSystemRuleID constant value for the schema field system_rule_id of a rule specification
	SystemRuleSpecificationSystemRuleID = ruleFieldPrefix + "system_rule_id"
)

//SystemRuleEntityType the fix entity_type of entity verification rules
const SystemRuleEntityType = "any"

var systemRuleSchemaFields = map[string]*schema.Schema{
	CustomEventSpecificationFieldEntityType: &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The computed entity type of a entity verification rule 'any'",
	},
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

		Schema:        mergeSchemaMap(defaultCustomEventSchemaFields, systemRuleSchemaFields),
		SchemaVersion: 1,
		MigrateState:  CreateMigrateCustomEventConfigStateFunction(make(map[int](func(inst *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error)))),
	}
}

func createReadCustomEventSpecificationWithSystemRule() func(*schema.ResourceData, interface{}) error {
	return createCustomEventSpecificationReadFunc(mapSystemRuleToTerraformState)
}

func createCreateCustomEventSpecificationWithSystemRule() func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		updateFunc := createCustomEventSpecificationUpdateFunc(mapSystemRuleToInstanaAPIModel, mapSystemRuleToTerraformState)

		d.SetId(RandomID())
		d.Set(CustomEventSpecificationFieldEntityType, SystemRuleEntityType)
		return updateFunc(d, meta)
	}
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
