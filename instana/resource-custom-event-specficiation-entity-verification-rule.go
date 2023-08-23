package instana

import (
	"context"
	"github.com/gessnerfl/terraform-provider-instana/tfutils"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceInstanaCustomEventSpecificationEntityVerificationRule the name of the terraform-provider-instana resource to manage custom event specifications with entity verification rule
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

// EntityVerificationRuleEntityType the fix entity_type of entity verification rules
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
		ValidateFunc: validation.StringInSlice(restapi.SupportedMatchingOperators.TerrafromSupportedValues(), false),
		StateFunc: func(val interface{}) string {
			operator, _ := restapi.SupportedMatchingOperators.FromTerraformValue(val.(string))
			return operator.InstanaAPIValue()
		},
		Description: "The operator which should be applied for matching the label for the given entity (e.g. IS, CONTAINS, STARTS_WITH, ENDS_WITH, NONE)",
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

// NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle creates a new ResourceHandle for the terraform resource of custom event specifications with entity verification rules
func NewCustomEventSpecificationWithEntityVerificationRuleResourceHandle() ResourceHandle[*restapi.CustomEventSpecification] {
	commons := &customEventSpecificationCommons{}
	return &customEventSpecificationWithEntityVerificationRuleResource{
		metaData: ResourceMetaData{
			ResourceName:  ResourceInstanaCustomEventSpecificationEntityVerificationRule,
			Schema:        MergeSchemaMap(defaultCustomEventSchemaFields, entityVerificationRuleSchemaFields),
			SchemaVersion: 4,
		},
		commons: commons,
	}
}

type customEventSpecificationWithEntityVerificationRuleResource struct {
	metaData ResourceMetaData
	commons  *customEventSpecificationCommons
}

func (r *customEventSpecificationWithEntityVerificationRuleResource) MetaData() *ResourceMetaData {
	return &r.metaData
}

func (r *customEventSpecificationWithEntityVerificationRuleResource) StateUpgraders() []schema.StateUpgrader {
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
		{
			Type:    r.schemaV2().CoreConfigSchema().ImpliedType(),
			Upgrade: r.migrateCustomEventConfigWithEntityVerificationRuleToVersion3ByChangingMatchingOperatorToInstanaRepresentation,
			Version: 2,
		},
		{
			Type:    r.schemaV3().CoreConfigSchema().ImpliedType(),
			Upgrade: r.commons.migrateCustomEventConfigFullStateFromV2toV3AndRemoveFullname,
			Version: 3,
		},
	}
}

func (r *customEventSpecificationWithEntityVerificationRuleResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource[*restapi.CustomEventSpecification] {
	return api.CustomEventSpecifications()
}

func (r *customEventSpecificationWithEntityVerificationRuleResource) SetComputedFields(d *schema.ResourceData) error {
	return d.Set(CustomEventSpecificationFieldEntityType, EntityVerificationRuleEntityType)
}

func (r *customEventSpecificationWithEntityVerificationRuleResource) UpdateState(d *schema.ResourceData, customEventSpecification *restapi.CustomEventSpecification) error {
	data := r.commons.getDataForBasicCustomEventSpecification(customEventSpecification)

	ruleSpec := customEventSpecification.Rules[0]
	severity, err := ConvertSeverityFromInstanaAPIToTerraformRepresentation(ruleSpec.Severity)
	if err != nil {
		return err
	}
	matchingOperator, err := ruleSpec.MatchingOperatorType()
	if err != nil {
		return err
	}
	data[CustomEventSpecificationRuleSeverity] = severity
	data[EntityVerificationRuleFieldMatchingEntityLabel] = ruleSpec.MatchingEntityLabel
	data[EntityVerificationRuleFieldMatchingEntityType] = ruleSpec.MatchingEntityType
	data[EntityVerificationRuleFieldMatchingOperator] = matchingOperator.InstanaAPIValue()
	data[EntityVerificationRuleFieldOfflineDuration] = ruleSpec.OfflineDuration

	d.SetId(customEventSpecification.ID)
	return tfutils.UpdateState(d, data)
}

func (r *customEventSpecificationWithEntityVerificationRuleResource) MapStateToDataObject(d *schema.ResourceData) (*restapi.CustomEventSpecification, error) {
	severity, err := ConvertSeverityFromTerraformToInstanaAPIRepresentation(d.Get(CustomEventSpecificationRuleSeverity).(string))
	if err != nil {
		return &restapi.CustomEventSpecification{}, err
	}
	entityLabel := d.Get(EntityVerificationRuleFieldMatchingEntityLabel).(string)
	entityType := d.Get(EntityVerificationRuleFieldMatchingEntityType).(string)

	matchingOperatorString := d.Get(EntityVerificationRuleFieldMatchingOperator).(string)
	matchingOperator, err := restapi.SupportedMatchingOperators.FromTerraformValue(matchingOperatorString)
	if err != nil {
		return &restapi.CustomEventSpecification{}, err
	}
	offlineDuration := d.Get(EntityVerificationRuleFieldOfflineDuration).(int)

	rule := restapi.NewEntityVerificationRuleSpecification(entityLabel, entityType, matchingOperator.InstanaAPIValue(), offlineDuration, severity)

	customEventSpecification := r.commons.createCustomEventSpecificationFromResourceData(d)
	customEventSpecification.Rules = []restapi.RuleSpecification{rule}
	return customEventSpecification, nil
}

func (r *customEventSpecificationWithEntityVerificationRuleResource) schemaV0() *schema.Resource {
	return &schema.Resource{
		Schema: MergeSchemaMap(defaultCustomEventSchemaFieldsV0, entityVerificationRuleSchemaFields),
	}
}

func (r *customEventSpecificationWithEntityVerificationRuleResource) schemaV1() *schema.Resource {
	return &schema.Resource{
		Schema: MergeSchemaMap(defaultCustomEventSchemaFieldsV1, entityVerificationRuleSchemaFields),
	}
}

func (r *customEventSpecificationWithEntityVerificationRuleResource) schemaV2() *schema.Resource {
	return &schema.Resource{
		Schema: MergeSchemaMap(defaultCustomEventSchemaFieldsV1, entityVerificationRuleSchemaFields),
	}
}

func (r *customEventSpecificationWithEntityVerificationRuleResource) migrateCustomEventConfigWithEntityVerificationRuleToVersion3ByChangingMatchingOperatorToInstanaRepresentation(_ context.Context, rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	v, ok := rawState[EntityVerificationRuleFieldMatchingOperator]
	if ok {
		operator, err := restapi.SupportedMatchingOperators.FromTerraformValue(v.(string))
		if err != nil {
			return rawState, err
		}
		rawState[EntityVerificationRuleFieldMatchingOperator] = operator.InstanaAPIValue()
	}
	return rawState, nil
}

func (r *customEventSpecificationWithEntityVerificationRuleResource) schemaV3() *schema.Resource {
	return &schema.Resource{
		Schema: MergeSchemaMap(defaultCustomEventSchemaFieldsV2, entityVerificationRuleSchemaFields),
	}
}
