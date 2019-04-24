package restapi

import "errors"

//CustomEventSpecificationResource represents the REST resource of custom event specification at Instana
type CustomEventSpecificationResource interface {
	GetOne(id string) (CustomEventSpecification, error)
	Upsert(spec CustomEventSpecification) (CustomEventSpecification, error)
	Delete(spec CustomEventSpecification) error
	DeleteByID(specID string) error
}

//RuleType custom type representing the type of the custom event specification rule
type RuleType string

const (
	//SystemRuleType const for RuleType of System
	SystemRuleType = "system"
	//ThresholdRuleType const for RuleType of Threshold
	ThresholdRuleType = "threshold"
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

//ConditionOperatorType custom type representing a condition operator of a custom event specification rule
type ConditionOperatorType string

//ConditionOperatorTypes custom type representing a slice of ConditionOperatorType
type ConditionOperatorTypes []ConditionOperatorType

//ToStringSlice Returns the string representations fo the condition operators
func (types ConditionOperatorTypes) ToStringSlice() []string {
	result := make([]string, len(types))
	for i, v := range types {
		result[i] = string(v)
	}
	return result
}

const (
	//ConditionOperatorEquals const for a equals (==) condition operator
	ConditionOperatorEquals = ConditionOperatorType("==")
	//ConditionOperatorNotEqual const for a not equal (!=) condition operator
	ConditionOperatorNotEqual = ConditionOperatorType("!=")
	//ConditionOperatorLessThan const for a less than (<) condition operator
	ConditionOperatorLessThan = ConditionOperatorType("<")
	//ConditionOperatorLessThanOrEqual const for a less than or equal (<=) condition operator
	ConditionOperatorLessThanOrEqual = ConditionOperatorType("<=")
	//ConditionOperatorGreaterThan const for a greater than (>) condition operator
	ConditionOperatorGreaterThan = ConditionOperatorType(">")
	//ConditionOperatorGreaterThanOrEqual const for a greater than or equal (<=) condition operator
	ConditionOperatorGreaterThanOrEqual = ConditionOperatorType(">=")
)

//SupportedConditionOperatorTypes slice of supported aggregation types
var SupportedConditionOperatorTypes = ConditionOperatorTypes{ConditionOperatorEquals, ConditionOperatorNotEqual, ConditionOperatorLessThan, ConditionOperatorLessThanOrEqual, ConditionOperatorGreaterThan, ConditionOperatorGreaterThanOrEqual}

//IsSupportedConditionOperatorType check if the provided condition operator type is supported
func IsSupportedConditionOperatorType(operator ConditionOperatorType) bool {
	for _, v := range SupportedConditionOperatorTypes {
		if v == operator {
			return true
		}
	}
	return false
}

//NewSystemRuleSpecification creates a new instance of System Rule
func NewSystemRuleSpecification(systemRuleID string, severity int) RuleSpecification {
	return RuleSpecification{
		DType:        SystemRuleType,
		SystemRuleID: systemRuleID,
		Severity:     severity,
	}
}

//RuleSpecification representation of a rule specification for a CustomEventSpecification
type RuleSpecification struct {
	//Common Fields
	DType    RuleType `json:"ruleType"`
	Severity int      `json:"severity"`

	//System Rule fields
	SystemRuleID string `json:"systemRuleId"`

	//Threshold Rule fields
	MetricName                         string                `json:"metricName"`
	Rollup                             *int                  `json:"rollup"`
	Window                             *int                  `json:"window"`
	Aggregation                        *AggregationType      `json:"aggregation"`
	ConditionOperator                  ConditionOperatorType `json:"conditionOperator"`
	ConditionValue                     *float64              `json:"conditionValue"`
	AggregationForNonPercentileMetric  bool                  `json:"aggregationForNonPercentileMetric"`
	EitherRollupOrWindowAndAggregation bool                  `json:"eitherRollupOrWindowAndAggregation"`
}

//Validate Rule interface implementation for SystemRule
func (r *RuleSpecification) Validate() error {
	if len(r.DType) == 0 {
		return errors.New("type of system rule is missing")
	}
	if r.DType == SystemRuleType {
		return r.validateSystemRule()
	} else if r.DType == ThresholdRuleType {
		return r.validateThresholdRule()
	}
	return nil
}

func (r *RuleSpecification) validateSystemRule() error {
	if len(r.SystemRuleID) == 0 {
		return errors.New("id of system rule is missing")
	}
	return nil
}

func (r *RuleSpecification) validateThresholdRule() error {
	if len(r.MetricName) == 0 {
		return errors.New("metric name of threshold rule is missing")
	}
	if r.Window == nil && r.Rollup == nil || r.Window != nil && r.Rollup != nil {
		return errors.New("either rollup or window and condition must be defined")
	}

	if r.Window != nil && (r.Aggregation == nil || !IsSupportedAggregationType(*r.Aggregation)) {
		return errors.New("aggregation type of threshold rule is mission or not valid")
	}

	if !IsSupportedConditionOperatorType(r.ConditionOperator) {
		return errors.New("condition operator of threshold rule is missing or not valid")
	}

	return nil
}

//EventSpecificationDownstream definition of downstream reporting for the event specification
type EventSpecificationDownstream struct {
	IntegrationIds                []string `json:"integrationIds"`
	BroadcastToAllAlertingConfigs bool     `json:"broadcastToAllAlertingConfigs"`
}

//Validate validates the consitency of an EventSpecificationDownstream
func (d EventSpecificationDownstream) Validate() error {
	if len(d.IntegrationIds) == 0 {
		return errors.New("At least one integration id must be defined for a downstream specification")
	}
	return nil
}

//CustomEventSpecification is the representation of a custom event specification in Instana
type CustomEventSpecification struct {
	ID             string                        `json:"id"`
	Name           string                        `json:"name"`
	EntityType     string                        `json:"entityType"`
	Query          *string                       `json:"query"`
	Triggering     bool                          `json:"triggering"`
	Description    *string                       `json:"description"`
	ExpirationTime *int                          `json:"expirationTime"`
	Enabled        bool                          `json:"enabled"`
	Rules          []RuleSpecification           `json:"rules"`
	Downstream     *EventSpecificationDownstream `json:"downstream"`
}

//GetID implemention of the interface InstanaDataObject
func (spec CustomEventSpecification) GetID() string {
	return spec.ID
}

//Validate implementation of the interface InstanaDataObject to verify if data object is correct
func (spec CustomEventSpecification) Validate() error {
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

	if spec.Downstream != nil {
		if err := spec.Downstream.Validate(); err != nil {
			return err
		}
	}
	return nil
}
