package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

//ResourceInstanaCustomEventSpecificationEntityVerificationRule the name of the terraform-provider-instana resource to manage custom event specifications with entity verification rule
const ResourceInstanaCustomEventSpecificationEntityVerificationRule = "instana_custom_event_spec_entity_verification_rule"

const (
	//EntityVerificationRuleFieldMatchingEntityType constant value for the schema field matching_entity_type
	EntityVerificationRuleFieldMatchingEntityType = ruleFieldPrefix + "matching_entity_type"
	//EntityVerificationRuleFieldMatchingOperator constant value for the schema field matching_operator
	EntityVerificationRuleFieldMatchingOperator = ruleFieldPrefix + "matching_operator"
	//EntityVerificationRuleFieldMatchingEntityLabel constant value for the schema field matching_entity_label
	EntityVerificationRuleFieldMatchingEntityLabel = ruleFieldPrefix + "matching_entity_label"
	//EntityVerificationRuleFieldOfflineDuration constant value for the schema field offline_duration
	EntityVerificationRuleFieldOfflineDuration = ruleFieldPrefix + "offline_duration"
)

//EntityVerificationRuleEntityType the fix entity_type of entity verification rules
const EntityVerificationRuleEntityType = "host"

var entityVerificationRuleSchemaFields = map[string]*schema.Schema{
	CustomEventSpecificationFieldEntityType: {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The computed entity type of a entity verification rule 'host'",
	},
	EntityVerificationRuleFieldMatchingEntityType: {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The type of the matching entity",
	},
	EntityVerificationRuleFieldMatchingOperator: {
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice(restapi.SupportedMatchingOperatorTypes.TerrafromSupportedValues(), false),
		Description:  "The operator which should be applied for matching the label for the given entity (e.g. IS, CONTAINS, STARTS_WITH, ENDS_WITH, NONE)",
	},
	EntityVerificationRuleFieldMatchingEntityLabel: {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The label of the matching entity",
	},
	EntityVerificationRuleFieldOfflineDuration: {
		Type:        schema.TypeInt,
		Required:    true,
		Description: "The duration after which the matching entity is considered to be offline",
	},
}

//NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle creates a new ResourceHandle for the terraform resource of custom event specifications with entity verification rules
func NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle() *ResourceHandle {
	return &ResourceHandle{
		ResourceName:  ResourceInstanaCustomEventSpecificationEntityVerificationRule,
		Schema:        mergeSchemaMap(defaultCustomEventSchemaFields, entityVerificationRuleSchemaFields),
		SchemaVersion: 2,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    customEventSpecificationWithEntityVerificationRuleSchemaV0().CoreConfigSchema().ImpliedType(),
				Upgrade: migrateCustomEventConfigFullNameInStateFromV0toV1,
				Version: 0,
			},
			{
				Type:    customEventSpecificationWithEntityVerificationRuleSchemaV1().CoreConfigSchema().ImpliedType(),
				Upgrade: migrateCustomEventConfigFullStateFromV1toV2AndRemoveDownstreamConfiguration,
				Version: 1,
			},
		},
		RestResourceFactory:  func(api restapi.InstanaAPI) restapi.RestResource { return api.CustomEventSpecifications() },
		UpdateState:          updateStateForCustomEventSpecificationWithEntityVerificationRule,
		MapStateToDataObject: mapStateToDataObjectForCustomEventSpecificationWithEntityVerificationRule,
		SetComputedFields: func(d *schema.ResourceData) {
			d.Set(CustomEventSpecificationFieldEntityType, EntityVerificationRuleEntityType)
		},
	}
}

func updateStateForCustomEventSpecificationWithEntityVerificationRule(d *schema.ResourceData, obj restapi.InstanaDataObject) error {
	customEventSpecification := obj.(restapi.CustomEventSpecification)
	updateStateForBasicCustomEventSpecification(d, customEventSpecification)

	ruleSpec := customEventSpecification.Rules[0]
	severity, err := ConvertSeverityFromInstanaAPIToTerraformRepresentation(ruleSpec.Severity)
	if err != nil {
		return err
	}
	matchingOperator, err := ruleSpec.MatchingOperatorType()
	if err != nil {
		return err
	}
	d.Set(CustomEventSpecificationRuleSeverity, severity)
	d.Set(EntityVerificationRuleFieldMatchingEntityLabel, ruleSpec.MatchingEntityLabel)
	d.Set(EntityVerificationRuleFieldMatchingEntityType, ruleSpec.MatchingEntityType)
	d.Set(EntityVerificationRuleFieldMatchingOperator, matchingOperator.TerraformRepresentation)
	d.Set(EntityVerificationRuleFieldOfflineDuration, ruleSpec.OfflineDuration)
	return nil
}

func mapStateToDataObjectForCustomEventSpecificationWithEntityVerificationRule(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (restapi.InstanaDataObject, error) {
	severity, err := ConvertSeverityFromTerraformToInstanaAPIRepresentation(d.Get(CustomEventSpecificationRuleSeverity).(string))
	if err != nil {
		return restapi.CustomEventSpecification{}, err
	}
	entityLabel := d.Get(EntityVerificationRuleFieldMatchingEntityLabel).(string)
	entityType := d.Get(EntityVerificationRuleFieldMatchingEntityType).(string)

	matchingOperatorString := d.Get(EntityVerificationRuleFieldMatchingOperator).(string)
	matchingOperator, err := restapi.SupportedMatchingOperatorTypes.ForTerraformRepresentation(matchingOperatorString)
	if err != nil {
		return restapi.CustomEventSpecification{}, err
	}
	offlineDuration := d.Get(EntityVerificationRuleFieldOfflineDuration).(int)

	rule := restapi.NewEntityVerificationRuleSpecification(entityLabel, entityType, matchingOperator.InstanaRepresentation, offlineDuration, severity)

	customEventSpecification := createCustomEventSpecificationFromResourceData(d, formatter)
	customEventSpecification.Rules = []restapi.RuleSpecification{rule}
	return customEventSpecification, nil
}

func customEventSpecificationWithEntityVerificationRuleSchemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: mergeSchemaMap(defaultCustomEventSchemaFieldsV0, entityVerificationRuleSchemaFields),
	}
}

func customEventSpecificationWithEntityVerificationRuleSchemaV1() *schema.Resource {
	return &schema.Resource{
		Schema: mergeSchemaMap(defaultCustomEventSchemaFieldsV1, entityVerificationRuleSchemaFields),
	}
}
