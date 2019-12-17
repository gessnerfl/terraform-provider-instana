package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/hashicorp/terraform/terraform"
)

const (
	//EntityVerificationRuleFieldMatchingEntityLabel constant value for the schema field matching_entity_label
	EntityVerificationRuleFieldMatchingEntityLabel = ruleFieldPrefix + "matching_entity_label"
	//EntityVerificationRuleFieldMatchingEntityType constant value for the schema field matching_entity_type
	EntityVerificationRuleFieldMatchingEntityType = ruleFieldPrefix + "matching_entity_type"
	//EntityVerificationRuleFieldMatchingOperator constant value for the schema field matching_operator
	EntityVerificationRuleFieldMatchingOperator = ruleFieldPrefix + "matching_operator"
	//EntityVerificationRuleFieldOfflineDuration constant value for the schema field offline_duration
	EntityVerificationRuleFieldOfflineDuration = ruleFieldPrefix + "offline_duration"
)

//EntityVerificationRuleEntityType the fix entity_type of entity verification rules
const EntityVerificationRuleEntityType = "host"

var entityVerificationRuleSchemaFields = map[string]*schema.Schema{
	CustomEventSpecificationFieldEntityType: &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The computed entity type of a entity verification rule 'host'",
	},
	EntityVerificationRuleFieldMatchingEntityLabel: &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The label of the matching entity",
	},
	EntityVerificationRuleFieldMatchingEntityType: &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The type of the matching entity",
	},
	EntityVerificationRuleFieldMatchingOperator: &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice(restapi.SupportedMatchingOperatorTypes.ToStringSlice(), false),
		Description:  "The operator which should be applied for matching the label for the given entity (e.g. IS, CONTAINS, STARTS_WITH, ENDS_WITH, NONE)",
	},
	EntityVerificationRuleFieldOfflineDuration: &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "The duration after which the matching entity is considered to be offline",
	},
}

//CreateResourceCustomEventSpecificationWithEntityVerificationRule creates the resource definition for the instana api endpoint for Custom Event Specifications for Threshold rules
func CreateResourceCustomEventSpecificationWithEntityVerificationRule() *schema.Resource {
	return &schema.Resource{
		Read:   createReadCustomEventSpecificationWithEntityVerificationRule(),
		Create: createCreateCustomEventSpecificationWithEntityVerificationRule(),
		Update: createUpdateCustomEventSpecificationWithEntityVerificationRule(),
		Delete: createDeleteCustomEventSpecificationWithEntityVerificationRule(),

		Schema:        mergeSchemaMap(defaultCustomEventSchemaFields, entityVerificationRuleSchemaFields),
		SchemaVersion: 1,
		MigrateState:  CreateMigrateCustomEventConfigStateFunction(make(map[int](func(inst *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error)))),
	}
}

func createReadCustomEventSpecificationWithEntityVerificationRule() func(*schema.ResourceData, interface{}) error {
	return createCustomEventSpecificationReadFunc(mapEntityVerificationRuleToTerraformState)
}

func createCreateCustomEventSpecificationWithEntityVerificationRule() func(*schema.ResourceData, interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		updateFunc := createCustomEventSpecificationUpdateFunc(mapEntityVerificationRuleToInstanaAPIModel, mapEntityVerificationRuleToTerraformState)

		d.SetId(RandomID())
		d.Set(CustomEventSpecificationFieldEntityType, EntityVerificationRuleEntityType)
		return updateFunc(d, meta)
	}
}

func createUpdateCustomEventSpecificationWithEntityVerificationRule() func(*schema.ResourceData, interface{}) error {
	return createCustomEventSpecificationUpdateFunc(mapEntityVerificationRuleToInstanaAPIModel, mapEntityVerificationRuleToTerraformState)
}

func createDeleteCustomEventSpecificationWithEntityVerificationRule() func(*schema.ResourceData, interface{}) error {
	return createCustomEventSpecificationDeleteFunc(mapEntityVerificationRuleToInstanaAPIModel)
}

func mapEntityVerificationRuleToInstanaAPIModel(d *schema.ResourceData) (restapi.RuleSpecification, error) {
	severity, err := ConvertSeverityFromTerraformToInstanaAPIRepresentation(d.Get(CustomEventSpecificationRuleSeverity).(string))
	if err != nil {
		return restapi.RuleSpecification{}, err
	}
	entityLabel := d.Get(EntityVerificationRuleFieldMatchingEntityLabel).(string)
	entityType := d.Get(EntityVerificationRuleFieldMatchingEntityType).(string)
	operator := restapi.MatchingOperatorType(d.Get(EntityVerificationRuleFieldMatchingOperator).(string))
	offlineDuration := d.Get(EntityVerificationRuleFieldOfflineDuration).(int)

	return restapi.NewEntityVerificationRuleSpecification(entityLabel, entityType, operator, offlineDuration, severity), nil
}

func mapEntityVerificationRuleToTerraformState(d *schema.ResourceData, spec restapi.CustomEventSpecification) error {
	ruleSpec := spec.Rules[0]
	severity, err := ConvertSeverityFromInstanaAPIToTerraformRepresentation(ruleSpec.Severity)
	if err != nil {
		return err
	}

	d.Set(CustomEventSpecificationRuleSeverity, severity)
	d.Set(EntityVerificationRuleFieldMatchingEntityLabel, ruleSpec.MatchingEntityLabel)
	d.Set(EntityVerificationRuleFieldMatchingEntityType, ruleSpec.MatchingEntityType)
	d.Set(EntityVerificationRuleFieldMatchingOperator, ruleSpec.MatchingOperator)
	d.Set(EntityVerificationRuleFieldOfflineDuration, ruleSpec.OfflineDuration)
	return nil
}
