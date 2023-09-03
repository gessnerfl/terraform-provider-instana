package instana

import (
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceInstanaCustomEventSpecification the name of the terraform-provider-instana resource to manage custom event specifications
const ResourceInstanaCustomEventSpecification = "instana_custom_event_specification"

const (
	CustomEventSpecificationFieldRule                   = "rule"
	CustomEventSpecificationFieldEntityVerificationRule = "entity_verification"
	CustomEventSpecificationFieldSystemRule             = "system"
	CustomEventSpecificationFieldThresholdRule          = "threshold"
)

var (
	customEventSpecificationRuleTypeKeys = []string{
		"rule.0.entity_verification",
		"rule.0.system",
		"rule.0.threshold",
	}
)

// NewCustomEventSpecificationResourceHandle creates a new ResourceHandle for the terraform resource of custom event specifications
func NewCustomEventSpecificationResourceHandle() ResourceHandle {
	commons := &customEventSpecificationCommons{}
	return &customEventSpecificationResource{
		metaData: ResourceMetaData{
			ResourceName: ResourceInstanaCustomEventSpecification,
			Schema: map[string]*schema.Schema{
				CustomEventSpecificationFieldName:           customEventSpecificationSchemaName,
				CustomEventSpecificationFieldQuery:          customEventSpecificationSchemaQuery,
				CustomEventSpecificationFieldTriggering:     customEventSpecificationSchemaTriggering,
				CustomEventSpecificationFieldDescription:    customEventSpecificationSchemaDescription,
				CustomEventSpecificationFieldExpirationTime: customEventSpecificationSchemaExpirationTime,
				CustomEventSpecificationFieldEnabled:        customEventSpecificationSchemaEnabled,
				CustomEventSpecificationFieldRule: {
					Type:        schema.TypeList,
					MinItems:    1,
					MaxItems:    1,
					Required:    true,
					Description: "Indicates the type of rule this custom event specification is about.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							CustomEventSpecificationFieldEntityVerificationRule: {
								Type:        schema.TypeList,
								MinItems:    0,
								MaxItems:    1,
								Optional:    true,
								Description: "Entity verification rule configuration",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										CustomEventSpecificationRuleSeverity: customEventSpecificationSchemaRuleSeverity,
										EntityVerificationRuleFieldMatchingEntityType: {
											Type:        schema.TypeString,
											Required:    true,
											Description: "The type of the matching entity",
										},
										EntityVerificationRuleFieldMatchingOperator: {
											Type:     schema.TypeString,
											Required: true,
											ValidateFunc: validation.StringInSlice([]string{
												"is",
												"contains",
												"startsWith",
												"endsWith"}, false),
											Description: "The operator which should be applied for matching the label for the given entity (e.g. is, contains, startsWith, endsWith)",
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
									},
								},
								ExactlyOneOf: customEventSpecificationRuleTypeKeys,
							},
							CustomEventSpecificationFieldSystemRule: {
								Type:        schema.TypeList,
								MinItems:    0,
								MaxItems:    1,
								Optional:    true,
								Description: "System rule configuration",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										CustomEventSpecificationRuleSeverity: customEventSpecificationSchemaRuleSeverity,
										SystemRuleSpecificationSystemRuleID: {
											Type:        schema.TypeString,
											Required:    true,
											Description: "Configures the system rule id for the system rule of the custom event specification",
										},
									},
								},
								ExactlyOneOf: customEventSpecificationRuleTypeKeys,
							},
							CustomEventSpecificationFieldThresholdRule: {
								Type:        schema.TypeList,
								MinItems:    0,
								MaxItems:    1,
								Optional:    true,
								Description: "System rule configuration",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										CustomEventSpecificationRuleSeverity: customEventSpecificationSchemaRuleSeverity,

										ThresholdRuleFieldMetricName: {
											Type:        schema.TypeString,
											Required:    false,
											Optional:    true,
											Description: "The metric name of the rule",
										},

										ThresholdRuleFieldRollup: {
											Type:        schema.TypeInt,
											Required:    false,
											Optional:    true,
											Description: "The rollup of the metric",
										},
										ThresholdRuleFieldWindow: {
											Type:        schema.TypeInt,
											Required:    false,
											Optional:    true,
											Description: "The time window where the condition has to be fulfilled",
										},
										ThresholdRuleFieldAggregation: {
											Type:         schema.TypeString,
											Required:     false,
											Optional:     true,
											ValidateFunc: validation.StringInSlice(restapi.SupportedAggregationTypes.ToStringSlice(), false),
											Description:  "The aggregation type (e.g. sum, avg)",
										},
										ThresholdRuleFieldConditionOperator: {
											Type:         schema.TypeString,
											Required:     true,
											ValidateFunc: validation.StringInSlice(restapi.SupportedConditionOperators.TerrafromSupportedValues(), false),
											StateFunc: func(val interface{}) string {
												operator, _ := restapi.SupportedConditionOperators.FromTerraformValue(val.(string))
												return operator.InstanaAPIValue()
											},
											Description: "The condition operator (e.g >, <)",
										},
										ThresholdRuleFieldConditionValue: {
											Type:        schema.TypeFloat,
											Required:    true,
											Description: "The expected condition value to fulfill the rule",
										},
										ThresholdRuleFieldMetricPatternPrefix: {
											Type:        schema.TypeString,
											Required:    false,
											Optional:    true,
											Description: "The metric pattern prefix of a dynamic built-in metrics",
										},
										ThresholdRuleFieldMetricPatternPostfix: {
											Type:        schema.TypeString,
											Required:    false,
											Optional:    true,
											Description: "The metric pattern postfix of a dynamic built-in metrics",
										},
										ThresholdRuleFieldMetricPatternPlaceholder: {
											Type:        schema.TypeString,
											Required:    false,
											Optional:    true,
											Description: "The metric pattern placeholer/condition value of a dynamic built-in metrics",
										},
										ThresholdRuleFieldMetricPatternOperator: {
											Type:         schema.TypeString,
											Required:     false,
											Optional:     true,
											ValidateFunc: validation.StringInSlice(restapi.SupportedMetricPatternOperatorTypes.ToStringSlice(), false),
											Description:  "The condition operator (e.g >, <)",
										},
									},
								},
								ExactlyOneOf: customEventSpecificationRuleTypeKeys,
							},
						},
					},
				},
			},
			SchemaVersion: 3,
		},
		commons: commons,
	}
}

type customEventSpecificationResource struct {
	metaData ResourceMetaData
	commons  *customEventSpecificationCommons
}

func (r *customEventSpecificationResource) MetaData() *ResourceMetaData {
	return &r.metaData
}

func (r *customEventSpecificationResource) StateUpgraders() []schema.StateUpgrader {
	return []schema.StateUpgrader{}
}

func (r *customEventSpecificationResource) GetRestResource(api restapi.InstanaAPI) restapi.RestResource {
	return api.CustomEventSpecifications()
}

func (r *customEventSpecificationResource) SetComputedFields(d *schema.ResourceData) {
	d.Set(CustomEventSpecificationFieldEntityType, EntityVerificationRuleEntityType)
}

func (r *customEventSpecificationResource) UpdateState(d *schema.ResourceData, obj restapi.InstanaDataObject, formatter utils.ResourceNameFormatter) error {
	customEventSpecification := obj.(*restapi.CustomEventSpecification)
	r.commons.updateStateForBasicCustomEventSpecification(d, customEventSpecification, formatter)

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
	d.Set(EntityVerificationRuleFieldMatchingOperator, matchingOperator.InstanaAPIValue())
	d.Set(EntityVerificationRuleFieldOfflineDuration, ruleSpec.OfflineDuration)
	return nil
}

func (r *customEventSpecificationResource) MapStateToDataObject(d *schema.ResourceData, formatter utils.ResourceNameFormatter) (restapi.InstanaDataObject, error) {
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

	customEventSpecification := r.commons.createCustomEventSpecificationFromResourceData(d, formatter)
	customEventSpecification.Rules = []restapi.RuleSpecification{rule}
	return customEventSpecification, nil
}
