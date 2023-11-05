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
	//EntityCountRuleType const for rule type of Entity Count
	EntityCountRuleType = "entity_count"
	//EntityCountVerificationRuleType const for rule type of Entity Count Verification
	EntityCountVerificationRuleType = "entity_count_verification"
	//HostAvailabilityRuleType const for rule type of Host Availability
	HostAvailabilityRuleType = "host_availability"
)

// MetricPattern representation of a metric pattern for dynamic built-in metrics for CustomEventSpecification
type MetricPattern struct {
	Prefix      string  `json:"prefix"`
	Postfix     *string `json:"postfix"`
	Placeholder *string `json:"placeholder"`
	Operator    string  `json:"operator"`
}

// RuleSpecification representation of a rule specification for a CustomEventSpecification
type RuleSpecification struct {
	//Common Fields
	DType    string `json:"ruleType"`
	Severity int    `json:"severity"`

	//System Rule fields
	SystemRuleID *string `json:"systemRuleId"`

	//Threshold Rule fields
	MetricName        *string        `json:"metricName"`
	Rollup            *int           `json:"rollup"`
	Window            *int           `json:"window"`
	Aggregation       *string        `json:"aggregation"`
	ConditionOperator *string        `json:"conditionOperator"`
	ConditionValue    *float64       `json:"conditionValue"`
	MetricPattern     *MetricPattern `json:"metricPattern"`

	//Entity Verification Rule
	MatchingEntityType  *string    `json:"matchingEntityType"`
	MatchingOperator    *string    `json:"matchingOperator"`
	MatchingEntityLabel *string    `json:"matchingEntityLabel"`
	OfflineDuration     *int       `json:"offlineDuration"`
	CloseAfter          *int       `json:"closeAfter"`
	TagFilter           *TagFilter `jsom:"tagFilter"`
}

// CustomEventSpecification is the representation of a custom event specification in Instana
type CustomEventSpecification struct {
	ID                  string              `json:"id"`
	Name                string              `json:"name"`
	EntityType          string              `json:"entityType"`
	Query               *string             `json:"query"`
	Triggering          bool                `json:"triggering"`
	Description         *string             `json:"description"`
	ExpirationTime      *int                `json:"expirationTime"`
	Enabled             bool                `json:"enabled"`
	RuleLogicalOperator string              `json:"ruleLogicalOperator"`
	Rules               []RuleSpecification `json:"rules"`
}

// GetIDForResourcePath implementation of the interface InstanaDataObject
func (spec *CustomEventSpecification) GetIDForResourcePath() string {
	return spec.ID
}
