package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceInstanaCustomEventSpecificationSystemRule the name of the terraform-provider-instana resource to manage custom event specifications with system rule
const ResourceInstanaCustomEventSpecificationSystemRule = "instana_custom_event_spec_system_rule"

const (
	//SystemRuleSpecificationSystemRuleID constant value for the schema field system_rule_id of a rule specification
	SystemRuleSpecificationSystemRuleID = ruleFieldPrefix + "system_rule_id"
)

// SystemRuleEntityType the fix entity_type of entity verification rules
const SystemRuleEntityType = "any"

var systemRuleSchemaFields = map[string]*schema.Schema{
	CustomEventSpecificationFieldEntityType: {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The computed entity type of a entity verification rule 'any'",
	},
	SystemRuleSpecificationSystemRuleID: {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Configures the system rule id for the system rule of the custom event specification",
	},
}

// NewCustomEventSpecificationWithSystemRuleResourceHandle creates a new ResourceHandle for the terraform resource of custom event specifications with system rules
func NewCustomEventSpecificationWithSystemRuleResourceHandle() ResourceHandle[*restapi.CustomEventSpecification] {
	commons := &customEventSpecificationCommons{}
	return &customEventSpecificationWithSystemRuleResource{
		metaData: ResourceMetaData{
			ResourceName:  ResourceInstanaCustomEventSpecificationSystemRule,
			Schema:        MergeSchemaMap(defaultCustomEventSchemaFields, systemRuleSchemaFields),
			SchemaVersion: 2,
		},
		commons: commons,
	}
}

type customEventSpecificationWithSystemRuleResource struct {
	metaData ResourceMetaData
	commons  *customEventSpecificationCommons
}

func (r *customEventSpecificationWithSystemRuleResource) MetaData() *ResourceMetaData {
	return &r.metaData
}

func (r *customEventSpecificationWithSystemRuleResource) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{
		{
			Type:    r.schemaV0().CoreConfigSchema().ImpliedType(),
			Upgrade: r.commons.migrateCustomEventConfigFullNameInStateFromV0toV1,
			Version: 0,
		},
		{
			Type:    r.schemaV1().CoreConfigSchema().ImpliedType(),
			Upgrade: r.commons.migrateCustomEventConfigFullStateFromV1toV2AndRemoveDownstreamConfiguration,
			Version: 1,
		},
	}
}

func (r *customEventSpecificationWithSystemRuleResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.CustomEventSpecification] {
	return api.CustomEventSpecifications()
}

func (r *customEventSpecificationWithSystemRuleResource) SetComputedFields(d *schema.ResourceData) error {
	return d.Set(CustomEventSpecificationFieldEntityType, SystemRuleEntityType)
}

func (r *customEventSpecificationWithSystemRuleResource) UpdateState(d *schema.ResourceData, customEventSpecification *restapi.CustomEventSpecification, formatter utils.ResourceNameFormatter) error {
	data := r.commons.getDataForBasicCustomEventSpecification(customEventSpecification, formatter)

	ruleSpec := customEventSpecification.Rules[0]
	severity, err := ConvertSeverityFromInstanaAPIToTerraformRepresentation(ruleSpec.Severity)
	if err != nil {
		return err
	}
	data[CustomEventSpecificationRuleSeverity] = severity
	data[SystemRuleSpecificationSystemRuleID] = ruleSpec.SystemRuleID

	d.SetId(customEventSpecification.ID)
	return tfutils.UpdateState(d, data)
}

func (r *customEventSpecificationWithSystemRuleResource) MapStateToDataObject(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (*restapi.CustomEventSpecification, error) {
	severity, err := ConvertSeverityFromTerraformToInstanaAPIRepresentation(d.Get(CustomEventSpecificationRuleSeverity).(string))
	if err != nil {
		return &restapi.CustomEventSpecification{}, err
	}
	systemRuleID := d.Get(SystemRuleSpecificationSystemRuleID).(string)
	rule := restapi.NewSystemRuleSpecification(systemRuleID, severity)

	customEventSpecification := r.commons.createCustomEventSpecificationFromResourceData(d, formatter)
	customEventSpecification.Rules = []restapi.RuleSpecification{rule}
	return customEventSpecification, nil
}

func (r *customEventSpecificationWithSystemRuleResource) schemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: MergeSchemaMap(defaultCustomEventSchemaFieldsV0, systemRuleSchemaFields),
	}
}

func (r *customEventSpecificationWithSystemRuleResource) schemaV1() *schema.Resource {
	return &schema.Resource{
		Schema: MergeSchemaMap(defaultCustomEventSchemaFieldsV1, systemRuleSchemaFields),
	}
}
