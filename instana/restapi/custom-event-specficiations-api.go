package restapi

const (
	//EventSpecificationBasePath path to Event Specification settings of Instana RESTful API
	EventSpecificationBasePath = EventSettingsBasePath + "/event-specifications"
	//CustomEventSpecificationResourcePath path to Custom Event Specification settings resource of Instana RESTful API
	CustomEventSpecificationResourcePath = EventSpecificationBasePath + "/custom"
)

const (
	//SystemRuleType const for rule type of System
	SystemRuleType = "system"
	//ThresholdRuleType const for rule type of Threshold
	ThresholdRuleType = "threshold"
	//EntityVerificationRuleType const for rule type of Entity Verification
	EntityVerificationRuleType = "entity_verification"
)

// AggregationType custom type representing an aggregation of a custom event specification rule
type AggregationType string

// AggregationTypes custom type representing a slice of AggregationType
type AggregationTypes []AggregationType

// ToStringSlice Returns the string representations fo the aggregations
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

// SupportedAggregationTypes slice of supported aggregation types
var SupportedAggregationTypes = AggregationTypes{AggregationSum, AggregationAvg, AggregationMin, AggregationMax}

// MetricPatternOperatorType the operator type of the metric pattern of a dynamic built-in metric
type MetricPatternOperatorType string

// MetricPatternOperatorTypes type definition of a slice of MetricPatternOperatorType
type MetricPatternOperatorTypes []MetricPatternOperatorType

// ToStringSlice Returns the string representations of the metric pattern operator types
func (types MetricPatternOperatorTypes) ToStringSlice() []string {
	result := make([]string, len(types))
	for i, v := range types {
		result[i] = string(v)
	}
	return result
}

// Constant values of all supported MetricPatternOperatorTypes
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

// SupportedMetricPatternOperatorTypes slice of all supported MetricPatternOperatorTypes of the Instana Web Rest API
var SupportedMetricPatternOperatorTypes = MetricPatternOperatorTypes{MetricPatternOperatorTypeIs, MetricPatternOperatorTypeContains, MetricPatternOperatorTypeAny, MetricPatternOperatorTypeStartsWith, MetricPatternOperatorTypeEndsWith}

// NewSystemRuleSpecification creates a new instance of System Rule
func NewSystemRuleSpecification(systemRuleID string, severity int) RuleSpecification {
	return RuleSpecification{
		DType:        SystemRuleType,
		SystemRuleID: &systemRuleID,
		Severity:     severity,
	}
}

// NewEntityVerificationRuleSpecification creates a new instance of Entity Verification Rule
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

// MetricPattern representation of a metric pattern for dynamic built-in metrics for CustomEventSpecification
type MetricPattern struct {
	Prefix      string                    `json:"prefix"`
	Postfix     *string                   `json:"postfix"`
	Placeholder *string                   `json:"placeholder"`
	Operator    MetricPatternOperatorType `json:"operator"`
}

// RuleSpecification representation of a rule specification for a CustomEventSpecification
type RuleSpecification struct {
	//Common Fields
	DType    string `json:"ruleType"`
	Severity int    `json:"severity"`

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

// ConditionOperatorType returns the ConditionOperator for the given Instana Web REST API representation when available. In case of invalid values an error will be returned
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

// MatchingOperatorType returns the MatchingOperatorType for the given Instana Web REST API representation when available. In case of invalid values an error will be returned
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

// CustomEventSpecification is the representation of a custom event specification in Instana
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

// GetIDForResourcePath implemention of the interface InstanaDataObject
func (spec *CustomEventSpecification) GetIDForResourcePath() string {
	return spec.ID
}
