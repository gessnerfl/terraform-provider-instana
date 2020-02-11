package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform/helper/schema"
)

//ResourceInstanaCustomEventSpecificationSystemRule the name of the terraform-provider-instana resource to manage custom event specifications with system rule
const ResourceInstanaCustomEventSpecificationSystemRule = "instana_custom_event_spec_system_rule"

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

//NewCustomEventSpecificationWithSystemRuleResourceHandle creates a new ResourceHandle for the terraform resource of custom event specifications with system rules
func NewCustomEventSpecificationWithSystemRuleResourceHandle() *ResourceHandle {
	return &ResourceHandle{
		ResourceName:  ResourceInstanaCustomEventSpecificationSystemRule,
		Schema:        mergeSchemaMap(defaultCustomEventSchemaFields, systemRuleSchemaFields),
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    customEventSpecificationWithSystemRuleSchemaV0().CoreConfigSchema().ImpliedType(),
				Upgrade: migrateCustomEventConfigFullNameInStateFromV0toV1,
				Version: 0,
			},
		},
		RestResourceFactory:  func(api restapi.InstanaAPI) restapi.RestResource { return api.CustomEventSpecifications() },
		UpdateState:          updateStateForCustomEventSpecificationWithSystemRule,
		MapStateToDataObject: mapStateToDataObjectForCustomEventSpecificationWithSystemRule,
		SetComputedFields: func(d *schema.ResourceData) {
			d.Set(CustomEventSpecificationFieldEntityType, SystemRuleEntityType)
		},
	}
}

func updateStateForCustomEventSpecificationWithSystemRule(d *schema.ResourceData, obj restapi.InstanaDataObject) error {
	customEventSpecification := obj.(restapi.CustomEventSpecification)
	return updateStateForBasicCustomEventSpecification(d, customEventSpecification, mapSystemRuleToTerraformState)
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

func mapStateToDataObjectForCustomEventSpecificationWithSystemRule(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (restapi.InstanaDataObject, error) {
	return createCustomEventSpecificationFromResourceData(d, formatter, mapSystemRuleToInstanaAPIModel)
}

func mapSystemRuleToInstanaAPIModel(d *schema.ResourceData) (restapi.RuleSpecification, error) {
	severity, err := ConvertSeverityFromTerraformToInstanaAPIRepresentation(d.Get(CustomEventSpecificationRuleSeverity).(string))
	if err != nil {
		return restapi.RuleSpecification{}, err
	}
	systemRuleID := d.Get(SystemRuleSpecificationSystemRuleID).(string)
	return restapi.NewSystemRuleSpecification(systemRuleID, severity), nil
}

func customEventSpecificationWithSystemRuleSchemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: mergeSchemaMap(defaultCustomEventSchemaFieldsV0, systemRuleSchemaFields),
	}
}
