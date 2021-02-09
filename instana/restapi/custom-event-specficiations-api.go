package restapi

import (
	"errors"

	"github.com/gessnerfl/terraform-provider-instana/utils"
)

const (
	//EventSpecificationBasePath path to Event Specification settings of Instana RESTful API
	EventSpecificationBasePath = EventSettingsBasePath + "/event-specifications"
	//CustomEventSpecificationResourcePath path to Custom Event Specification settings resource of Instana RESTful API
	CustomEventSpecificationResourcePath = EventSpecificationBasePath + "/custom"
)

//Severity representation of the severity in both worlds Instana API and Terraform Provider
type Severity struct {
	apiRepresentation       int
	terraformRepresentation string
}

//GetAPIRepresentation returns the integer representation of the Instana API
func (s Severity) GetAPIRepresentation() int { return s.apiRepresentation }

//GetTerraformRepresentation returns the string representation of the Terraform Provider
func (s Severity) GetTerraformRepresentation() string { return s.terraformRepresentation }

//SeverityCritical representation of the critical severity
var SeverityCritical = Severity{apiRepresentation: 10, terraformRepresentation: "critical"}

//SeverityWarning representation of the warning severity
var SeverityWarning = Severity{apiRepresentation: 5, terraformRepresentation: "warning"}

//RuleType custom type representing the type of the custom event specification rule
type RuleType string

const (
	//SystemRuleType const for RuleType of System
	SystemRuleType = "system"
	//ThresholdRuleType const for RuleType of Threshold
	ThresholdRuleType = "threshold"
	//EntityVerificationRuleType const for RuleType of Entity Verification
	EntityVerificationRuleType = "entity_verification"
)

//AggregationType custom type representing an aggregation of a custom event specification rule
type AggregationType string

//AggregationTypes custom type representing a slice of AggregationType
type AggregationTypes []AggregationType

//ToStringSlice Returns the string representations fo the aggregations
func (types AggregationTypes) ToStringSlice() []string {
	result := make([]string, len(types))
	for i, v := range types {
		result[i] = string(v)
	}
	return result
}

const (
	//AggregationSum const for a sum aggregation
	AggregationSum = AggregationType("sum")
	//AggregationAvg const for a avg aggregation
	AggregationAvg = AggregationType("avg")
	//AggregationMin const for a min aggregation
	AggregationMin = AggregationType("min")
	//AggregationMax const for a max aggregation
	AggregationMax = AggregationType("max")
)

//SupportedAggregationTypes slice of supported aggregation types
var SupportedAggregationTypes = AggregationTypes{AggregationSum, AggregationAvg, AggregationMin, AggregationMax}

//IsSupportedAggregationType check if the provided aggregation type is supported
func IsSupportedAggregationType(aggregation AggregationType) bool {
	for _, v := range SupportedAggregationTypes {
		if v == aggregation {
			return true
		}
	}
	return false
}

//MetricPatternOperatorType the operator type of the metric pattern of a dynamic built-in metric
type MetricPatternOperatorType string

//MetricPatternOperatorTypes type definition of a slice of MetricPatternOperatorType
type MetricPatternOperatorTypes []MetricPatternOperatorType

//ToStringSlice Returns the string representations of the metric pattern operator types
func (types MetricPatternOperatorTypes) ToStringSlice() []string {
	result := make([]string, len(types))
	for i, v := range types {
		result[i] = string(v)
	}
	return result
}

//IsSupported checks if the given value is a valid representation of a supported MetricPatternOperatorType
func (types MetricPatternOperatorTypes) IsSupported(val MetricPatternOperatorType) bool {
	for _, t := range types {
		if t == val {
			return true
		}
	}
	return false
}

//Constant values of all supported MetricPatternOperatorTypes
const (
	//MetricPatternOperatorTypeIs constant value for the metric pattern operator type 'is'
	MetricPatternOperatorTypeIs = MetricPatternOperatorType("is")
	//MetricPatternOperatorTypeContains constant value for the metric pattern operator type 'contains'
	MetricPatternOperatorTypeContains = MetricPatternOperatorType("contains")
	//MetricPatternOperatorTypeAny constant value for the metric pattern operator type 'any'
	MetricPatternOperatorTypeAny = MetricPatternOperatorType("any")
	//MetricPatternOperatorTypeStartsWith constant value for the metric pattern operator type 'startsWith'
	MetricPatternOperatorTypeStartsWith = MetricPatternOperatorType("startsWith")
	//MetricPatternOperatorTypeEndsWith constant value for the metric pattern operator type 'endsWith'
	MetricPatternOperatorTypeEndsWith = MetricPatternOperatorType("endsWith")
)

//SupportedMetricPatternOperatorTypes slice of all supported MetricPatternOperatorTypes of the Instana Web Rest API
var SupportedMetricPatternOperatorTypes = MetricPatternOperatorTypes{MetricPatternOperatorTypeIs, MetricPatternOperatorTypeContains, MetricPatternOperatorTypeAny, MetricPatternOperatorTypeStartsWith, MetricPatternOperatorTypeEndsWith}

//NewSystemRuleSpecification creates a new instance of System Rule
func NewSystemRuleSpecification(systemRuleID string, severity int) RuleSpecification {
	return RuleSpecification{
		DType:        SystemRuleType,
		SystemRuleID: &systemRuleID,
		Severity:     severity,
	}
}

//NewEntityVerificationRuleSpecification creates a new instance of Entity Verification Rule
func NewEntityVerificationRuleSpecification(matchingEntityLabel string, matchingEntityType string, matchingOperator string, offlineDuration int, severity int) RuleSpecification {
	return RuleSpecification{
		DType:               EntityVerificationRuleType,
		MatchingEntityLabel: &matchingEntityLabel,
		MatchingEntityType:  &matchingEntityType,
		MatchingOperator:    &matchingOperator,
		OfflineDuration:     &offlineDuration,
		Severity:            severity,
	}
}

//MetricPattern representation of a metric pattern for dynamic built-in metrics for CustomEventSpecification
type MetricPattern struct {
	Prefix      string                    `json:"prefix"`
	Postfix     *string                   `json:"postfix"`
	Placeholder *string                   `json:"placeholder"`
	Operator    MetricPatternOperatorType `json:"operator"`
}

//Validate checks if the given MetricPattern is consistent
func (m *MetricPattern) Validate() error {
	if utils.IsBlank(m.Prefix) {
		return errors.New("Metric pattern prefix is missing")
	}
	if !SupportedMetricPatternOperatorTypes.IsSupported(m.Operator) {
		return errors.New("Metric pattern operator is not supported")
	}
	return nil
}

//RuleSpecification representation of a rule specification for a CustomEventSpecification
type RuleSpecification struct {
	//Common Fields
	DType    RuleType `json:"ruleType"`
	Severity int      `json:"severity"`

	//System Rule fields
	SystemRuleID *string `json:"systemRuleId"`

	//Threshold Rule fields
	MetricName        *string          `json:"metricName"`
	Rollup            *int             `json:"rollup"`
	Window            *int             `json:"window"`
	Aggregation       *AggregationType `json:"aggregation"`
	ConditionOperator *string          `json:"conditionOperator"`
	ConditionValue    *float64         `json:"conditionValue"`
	MetricPattern     *MetricPattern   `json:"metricPattern"`

	//Entity Verification Rule
	MatchingEntityType  *string `json:"matchingEntityType"`
	MatchingOperator    *string `json:"matchingOperator"`
	MatchingEntityLabel *string `json:"matchingEntityLabel"`
	OfflineDuration     *int    `json:"offlineDuration"`
}

//ConditionOperatorType returns the ConditionOperator for the given Instana Web REST API representation when available. In case of invalid values an error will be returned
func (r *RuleSpecification) ConditionOperatorType() (ConditionOperator, error) {
	if r.ConditionOperator != nil {
		operator, err := SupportedConditionOperators.FromInstanaAPIValue(*r.ConditionOperator)
		if err != nil {
			return nil, err
		}
		return operator, nil
	}
	return nil, nil
}

//MatchingOperatorType returns the MatchingOperatorType for the given Instana Web REST API representation when available. In case of invalid values an error will be returned
func (r *RuleSpecification) MatchingOperatorType() (MatchingOperator, error) {
	if r.MatchingOperator != nil {
		operator, err := SupportedMatchingOperators.FromInstanaAPIValue(*r.MatchingOperator)
		if err != nil {
			return nil, err
		}
		return operator, nil
	}
	return nil, nil
}

//Validate Rule interface implementation for SystemRule
func (r *RuleSpecification) Validate() error {
	if len(r.DType) == 0 {
		return errors.New("type of rule is missing")
	}
	if r.DType == SystemRuleType {
		return r.validateSystemRule()
	} else if r.DType == ThresholdRuleType {
		return r.validateThresholdRule()
	} else if r.DType == EntityVerificationRuleType {
		return r.validateEntityVerificationRule()
	}
	return errors.New("Unsupported rule type " + string(r.DType))
}

func (r *RuleSpecification) validateSystemRule() error {
	if r.SystemRuleID == nil || len(*r.SystemRuleID) == 0 {
		return errors.New("id of system rule is missing")
	}
	return nil
}

func (r *RuleSpecification) validateThresholdRule() error {
	if ((r.MetricName == nil || utils.IsBlank(*r.MetricName)) && r.MetricPattern == nil) || (r.MetricName != nil && !utils.IsBlank(*r.MetricName) && r.MetricPattern != nil) {
		return errors.New("either metric name or metric pattern of threshold rule needs to be defined")
	}
	if (r.Window == nil && r.Rollup == nil) || (r.Window != nil && r.Rollup != nil && *r.Window == 0 && *r.Rollup == 0) {
		return errors.New("either rollup or window and condition must be defined")
	}

	if r.Window != nil && (r.Aggregation == nil || !IsSupportedAggregationType(*r.Aggregation)) {
		return errors.New("aggregation type of threshold rule is mission or not valid")
	}

	if r.ConditionOperator == nil || !SupportedConditionOperators.IsSupportedInstanaAPIConditionOperator(*r.ConditionOperator) {
		return errors.New("condition operator of threshold rule is missing or not valid")
	}
	if r.MetricPattern != nil {
		return r.MetricPattern.Validate()
	}
	return nil
}

func (r *RuleSpecification) validateEntityVerificationRule() error {
	if r.MatchingEntityLabel == nil || len(*r.MatchingEntityLabel) == 0 {
		return errors.New("matching entity label of entity verification rule is missing")
	}
	if r.MatchingEntityType == nil || len(*r.MatchingEntityType) == 0 {
		return errors.New("matching entity type of entity verification rule is missing")
	}
	if r.MatchingOperator == nil || !SupportedMatchingOperators.IsSupportedInstanaAPIMatchingOperator(*r.MatchingOperator) {
		return errors.New("matching operator of entity verification rule is missing or not valid")
	}
	if r.OfflineDuration == nil {
		return errors.New("offline duration of entity verification rule is missing")
	}
	return nil
}

//CustomEventSpecification is the representation of a custom event specification in Instana
type CustomEventSpecification struct {
	ID             string              `json:"id"`
	Name           string              `json:"name"`
	EntityType     string              `json:"entityType"`
	Query          *string             `json:"query"`
	Triggering     bool                `json:"triggering"`
	Description    *string             `json:"description"`
	ExpirationTime *int                `json:"expirationTime"`
	Enabled        bool                `json:"enabled"`
	Rules          []RuleSpecification `json:"rules"`
}

//GetID implemention of the interface InstanaDataObject
func (spec *CustomEventSpecification) GetID() string {
	return spec.ID
}

//Validate implementation of the interface InstanaDataObject to verify if data object is correct
func (spec *CustomEventSpecification) Validate() error {
	if len(spec.ID) == 0 {
		return errors.New("ID is missing")
	}
	if len(spec.Name) == 0 {
		return errors.New("name is missing")
	}
	if len(spec.EntityType) == 0 {
		return errors.New("entity type is missing")
	}
	if len(spec.Rules) != 1 {
		return errors.New("exactly one rule must be defined")
	}
	for _, r := range spec.Rules {
		if err := r.Validate(); err != nil {
			return err
		}
	}
	return nil
}
