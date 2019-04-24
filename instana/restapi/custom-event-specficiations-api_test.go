package restapi_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gessnerfl/terraform-provider-instana/testutils"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/google/go-cmp/cmp"
)

const (
	customEventID                = "custom-event-id"
	customEventName              = "custom-event-name"
	customEventEntityType        = "custom-event-entity-type"
	customEventQuery             = "custom-event-query"
	customEventDescription       = "custom-event-description"
	customEventSystemRuleID      = "system-rule-id"
	customEventMetricName        = "threshold-rule-metric-name"
	customEventWindow            = 60000
	customEventRollup            = 40000
	customEventAggregation       = AggregationSum
	customEventConditionOperator = ConditionOperatorEquals
	customEventConditionValue    = 1.2
)

func TestShouldValidateMinimalCustemEventSpecificationWithSystemRule(t *testing.T) {
	spec := CustomEventSpecification{
		ID:         customEventID,
		Name:       customEventName,
		EntityType: customEventEntityType,
		Rules:      []RuleSpecification{RuleSpecification{DType: SystemRuleType, SystemRuleID: customEventSystemRuleID}},
	}

	if err := spec.Validate(); err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}

	if customEventID != spec.GetID() {
		t.Fatal("Expected GetID returns the correct id of the custom event specification")
	}
}

func TestShouldValidateFullCustomEventSpecificationWithSystemRule(t *testing.T) {
	query := customEventQuery
	description := customEventDescription
	expirationTime := 1234

	spec := CustomEventSpecification{
		ID:             customEventID,
		Name:           customEventName,
		EntityType:     customEventEntityType,
		Query:          &query,
		Description:    &description,
		ExpirationTime: &expirationTime,
		Triggering:     true,
		Enabled:        true,
		Rules:          []RuleSpecification{RuleSpecification{DType: SystemRuleType, SystemRuleID: customEventSystemRuleID}},
		Downstream: &EventSpecificationDownstream{
			IntegrationIds:                []string{"downstream-integration-id"},
			BroadcastToAllAlertingConfigs: true,
		},
	}

	if err := spec.Validate(); err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}
}

func TestFailToValidateCustemEventSpecificationWithSystemRuleWhenIDIsMissing(t *testing.T) {
	spec := CustomEventSpecification{
		Name:       customEventName,
		EntityType: customEventEntityType,
		Rules:      []RuleSpecification{RuleSpecification{DType: SystemRuleType, SystemRuleID: customEventSystemRuleID}},
	}

	if err := spec.Validate(); err == nil || !strings.Contains(err.Error(), "ID") {
		t.Fatal("Expected validate to fail as ID is not provided")
	}
}

func TestFailToValidateCustemEventSpecificationWithSystemRuleWhenNameIsMissing(t *testing.T) {
	spec := CustomEventSpecification{
		ID:         customEventID,
		EntityType: customEventEntityType,
		Rules:      []RuleSpecification{RuleSpecification{DType: SystemRuleType, SystemRuleID: customEventSystemRuleID}},
	}

	if err := spec.Validate(); err == nil || !strings.Contains(err.Error(), "name") {
		t.Fatal("Expected validate to fail as name is not provided")
	}
}

func TestFailToValidateCustemEventSpecificationWithSystemRuleWhenEntityTypeIsMissing(t *testing.T) {
	spec := CustomEventSpecification{
		ID:    customEventID,
		Name:  customEventName,
		Rules: []RuleSpecification{RuleSpecification{DType: SystemRuleType, SystemRuleID: customEventSystemRuleID}},
	}

	if err := spec.Validate(); err == nil || !strings.Contains(err.Error(), "entity type") {
		t.Fatal("Expected validate to fail as entity type is not provided")
	}
}

func TestFailToValidateCustemEventSpecificationWhenNoRuleIsNil(t *testing.T) {
	spec := CustomEventSpecification{
		ID:         customEventID,
		Name:       customEventName,
		EntityType: customEventEntityType,
		Rules:      nil,
	}

	if err := spec.Validate(); err == nil || !strings.Contains(err.Error(), "exactly one rule") {
		t.Fatal("Expected validate to fail as no rule is provided")
	}
}

func TestFailToValidateCustemEventSpecificationWhenNoRuleIsProvided(t *testing.T) {
	spec := CustomEventSpecification{
		ID:         customEventID,
		Name:       customEventName,
		EntityType: customEventEntityType,
		Rules:      []RuleSpecification{},
	}

	if err := spec.Validate(); err == nil || !strings.Contains(err.Error(), "exactly one rule") {
		t.Fatal("Expected validate to fail as no rule is provided")
	}
}

func TestFailToValidateCustemEventSpecificationWhenMultipleRulesAreProvided(t *testing.T) {
	spec := CustomEventSpecification{
		ID:         customEventID,
		Name:       customEventName,
		EntityType: customEventEntityType,
		Rules:      []RuleSpecification{RuleSpecification{DType: SystemRuleType, SystemRuleID: customEventSystemRuleID}, RuleSpecification{DType: SystemRuleType, SystemRuleID: customEventSystemRuleID}},
	}

	if err := spec.Validate(); err == nil || !strings.Contains(err.Error(), "exactly one rule") {
		t.Fatal("Expected validate to fail as no id of the second system rule is not provided")
	}
}

func TestFailToValidateCustemEventSpecificationWhenTheProvidedRuleIsNotValid(t *testing.T) {
	spec := CustomEventSpecification{
		ID:         customEventID,
		Name:       customEventName,
		EntityType: customEventEntityType,
		Rules:      []RuleSpecification{RuleSpecification{DType: SystemRuleType}},
	}

	if err := spec.Validate(); err == nil || !strings.Contains(err.Error(), "id of system rule") {
		t.Fatal("Expected validate to fail as no id of the second system rule is not provided")
	}
}

func TestFailToValidateCustemEventSpecificationWhenDownstreamSpecificationIsNotValid(t *testing.T) {
	spec := CustomEventSpecification{
		ID:         customEventID,
		Name:       customEventName,
		EntityType: customEventEntityType,
		Rules:      []RuleSpecification{RuleSpecification{DType: SystemRuleType, SystemRuleID: customEventSystemRuleID}},
		Downstream: &EventSpecificationDownstream{},
	}

	if err := spec.Validate(); err == nil || !strings.Contains(err.Error(), "integration id") {
		t.Fatal("Expected validate to fail as no integration id is provided for the downstream specification")
	}
}

func TestShouldSuccessfullyValidateSystemRule(t *testing.T) {
	rule := NewSystemRuleSpecification(customEventSystemRuleID, 1000)

	if err := rule.Validate(); err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}
}

func TestShouldFailToValidateSystemRuleWhenSystemRuleIDIsMissing(t *testing.T) {
	rule := RuleSpecification{DType: SystemRuleType}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "id of system rule") {
		t.Fatal("Expected to fail to validate system rule as no system rule id is provided")
	}
}

func TestShouldFailToValidateSystemRuleWhenRuleTypeIsMissing(t *testing.T) {
	rule := RuleSpecification{SystemRuleID: customEventSystemRuleID}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "type of system rule") {
		t.Fatal("Expected to fail to validate system rule as no system rule id is provided")
	}
}

func TestShouldSuccessfullyValidateEventSpecificationDownstream(t *testing.T) {
	downstream := EventSpecificationDownstream{
		IntegrationIds:                []string{"integration-id-1"},
		BroadcastToAllAlertingConfigs: true,
	}

	if err := downstream.Validate(); err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}
}

func TestShouldFailToValidateEventSpecificationDownstreamWhenNoIntegrationIDsAreNil(t *testing.T) {
	downstream := EventSpecificationDownstream{
		IntegrationIds: nil,
	}

	if err := downstream.Validate(); err == nil || !strings.Contains(err.Error(), "integration id") {
		t.Fatal("Expected to fail to validate event specification downstream as integration ids is nil")
	}
}

func TestShouldFailToValidateEventSpecificationDownstreamWhenNoIntegrationIDIsProvided(t *testing.T) {
	downstream := EventSpecificationDownstream{
		IntegrationIds: []string{},
	}

	if err := downstream.Validate(); err == nil || !strings.Contains(err.Error(), "integration id") {
		t.Fatal("Expected to fail to validate event specification downstream as no integration id is provided")
	}
}

func TestShouldValidateFullThresholdRuleSpecificationWithWindowAndAggregation(t *testing.T) {
	aggregation := customEventAggregation
	conditionValue := customEventConditionValue
	window := customEventWindow
	rule := RuleSpecification{
		DType:             ThresholdRuleType,
		Severity:          SeverityWarning.GetAPIRepresentation(),
		MetricName:        customEventMetricName,
		Window:            &window,
		Aggregation:       &aggregation,
		ConditionOperator: customEventConditionOperator,
		ConditionValue:    &conditionValue,
	}

	if err := rule.Validate(); err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}
}

func TestShouldSuccessfullyValidateMinimalThresholdRuleSpecificationForAllSupportedAggregations(t *testing.T) {
	for _, a := range SupportedAggregationTypes {
		t.Run(fmt.Sprintf("TestShouldSuccessfullyValidateMinimalThresholdRuleForAggregation%s", a), createTestCaseForSuccessfullValidateMinimalThresholdRuleForAggregation(a))
	}
}

func createTestCaseForSuccessfullValidateMinimalThresholdRuleForAggregation(aggregation AggregationType) func(*testing.T) {
	return func(t *testing.T) {
		conditionValue := customEventConditionValue
		window := customEventWindow
		rule := RuleSpecification{
			DType:             ThresholdRuleType,
			Severity:          SeverityWarning.GetAPIRepresentation(),
			MetricName:        customEventMetricName,
			Window:            &window,
			Aggregation:       &aggregation,
			ConditionOperator: customEventConditionOperator,
			ConditionValue:    &conditionValue,
		}

		if err := rule.Validate(); err != nil {
			t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
		}
	}
}

func TestShouldSuccessfullyValidateMinimalThresholdRuleSpecificationForAllSupportedConditionOperators(t *testing.T) {
	for _, o := range SupportedConditionOperatorTypes {
		t.Run(fmt.Sprintf("TestShouldSuccessfullyValidateMinimalThresholdRuleForConditionOperator%s", o), createTestCaseForSuccessfullValidateMinimalThresholdRuleForConditionOperators(o))
	}
}

func createTestCaseForSuccessfullValidateMinimalThresholdRuleForConditionOperators(operator ConditionOperatorType) func(*testing.T) {
	return func(t *testing.T) {
		aggregation := customEventAggregation
		conditionValue := customEventConditionValue
		window := customEventWindow
		rule := RuleSpecification{
			DType:             ThresholdRuleType,
			Severity:          SeverityWarning.GetAPIRepresentation(),
			MetricName:        customEventMetricName,
			Window:            &window,
			Aggregation:       &aggregation,
			ConditionOperator: operator,
			ConditionValue:    &conditionValue,
		}

		if err := rule.Validate(); err != nil {
			t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
		}
	}
}

func TestShouldValidateMinimalThresholdRuleSpecificationWithRollup(t *testing.T) {
	rollup := customEventRollup
	conditionValue := customEventConditionValue
	rule := RuleSpecification{
		DType:             ThresholdRuleType,
		Severity:          SeverityWarning.GetAPIRepresentation(),
		MetricName:        customEventMetricName,
		Rollup:            &rollup,
		ConditionOperator: customEventConditionOperator,
		ConditionValue:    &conditionValue,
	}

	if err := rule.Validate(); err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}
}

func TestShouldValidateFullThresholdRuleSpecificationWithRollup(t *testing.T) {
	rollup := customEventRollup
	conditionValue := customEventConditionValue
	rule := RuleSpecification{
		DType:                              ThresholdRuleType,
		Severity:                           SeverityWarning.GetAPIRepresentation(),
		MetricName:                         customEventMetricName,
		Rollup:                             &rollup,
		ConditionOperator:                  customEventConditionOperator,
		ConditionValue:                     &conditionValue,
		AggregationForNonPercentileMetric:  true,
		EitherRollupOrWindowAndAggregation: true,
	}

	if err := rule.Validate(); err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}
}

func TestShouldFailToValidateThresholdRuleSpecificationWithWindowWhenMetricNameIsMissing(t *testing.T) {
	aggregation := customEventAggregation
	conditionValue := customEventConditionValue
	window := customEventWindow
	rule := RuleSpecification{
		DType:             ThresholdRuleType,
		Severity:          SeverityWarning.GetAPIRepresentation(),
		Window:            &window,
		Aggregation:       &aggregation,
		ConditionOperator: customEventConditionOperator,
		ConditionValue:    &conditionValue,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "metric name") {
		t.Fatal("Expected to fail to validate threshold rule as no metric name is provided")
	}
}

func TestShouldFailToValidateThresholdRuleSpecificationWithWindowWhenNeitherNorRollupWindowIsDefined(t *testing.T) {
	aggregation := customEventAggregation
	conditionValue := customEventConditionValue
	rule := RuleSpecification{
		DType:             ThresholdRuleType,
		Severity:          SeverityWarning.GetAPIRepresentation(),
		MetricName:        customEventMetricName,
		Aggregation:       &aggregation,
		ConditionOperator: customEventConditionOperator,
		ConditionValue:    &conditionValue,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "either rollup or window") {
		t.Fatal("Expected to fail to validate threshold rule as no window is provided")
	}
}

func TestShouldFailToValidateThresholdRuleSpecificationWithWindowWhenRollupAndWindowIsDefined(t *testing.T) {
	window := customEventWindow
	rollup := customEventRollup
	aggregation := customEventAggregation
	conditionValue := customEventConditionValue
	rule := RuleSpecification{
		DType:             ThresholdRuleType,
		Severity:          SeverityWarning.GetAPIRepresentation(),
		MetricName:        customEventMetricName,
		Rollup:            &rollup,
		Window:            &window,
		Aggregation:       &aggregation,
		ConditionOperator: customEventConditionOperator,
		ConditionValue:    &conditionValue,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "either rollup or window") {
		t.Fatal("Expected to fail to validate threshold rule as no window is provided")
	}
}

func TestShouldFailToValidateThresholdRuleSpecificationWithWindowWhenAggregationIsMissingConditionOperator(t *testing.T) {
	conditionValue := customEventConditionValue
	window := customEventWindow
	rule := RuleSpecification{
		DType:             ThresholdRuleType,
		Severity:          SeverityWarning.GetAPIRepresentation(),
		MetricName:        customEventMetricName,
		Window:            &window,
		ConditionOperator: customEventConditionOperator,
		ConditionValue:    &conditionValue,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "aggregation") {
		t.Fatal("Expected to fail to validate threshold rule as no aggregation is provided")
	}
}

func TestShouldFailToValidateThresholdRuleSpecificationWithWindowWhenAggregationIsNotValid(t *testing.T) {
	aggregation := AggregationType("invalid")
	conditionValue := customEventConditionValue
	window := customEventWindow
	rule := RuleSpecification{
		DType:             ThresholdRuleType,
		Severity:          SeverityWarning.GetAPIRepresentation(),
		MetricName:        customEventMetricName,
		Window:            &window,
		Aggregation:       &aggregation,
		ConditionOperator: customEventConditionOperator,
		ConditionValue:    &conditionValue,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "aggregation") {
		t.Fatal("Expected to fail to validate threshold rule as no aggregation is provided")
	}
}

func TestShouldFailToValidateThresholdRuleSpecificationWithWindowWhenConditionOperatorIsMissing(t *testing.T) {
	aggregation := customEventAggregation
	conditionValue := customEventConditionValue
	window := customEventWindow
	rule := RuleSpecification{
		DType:          ThresholdRuleType,
		Severity:       SeverityWarning.GetAPIRepresentation(),
		MetricName:     customEventMetricName,
		Window:         &window,
		Aggregation:    &aggregation,
		ConditionValue: &conditionValue,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "condition operator") {
		t.Fatal("Expected to fail to validate threshold rule as no condition operator is provided")
	}
}

func TestShouldFailToValidateThresholdRuleSpecificationWithWindowWhenConditionOperatorIsNotValid(t *testing.T) {
	aggregation := customEventAggregation
	conditionOperator := ConditionOperatorType("invalid")
	conditionValue := customEventConditionValue
	window := customEventWindow
	rule := RuleSpecification{
		DType:             ThresholdRuleType,
		Severity:          SeverityWarning.GetAPIRepresentation(),
		MetricName:        customEventMetricName,
		Window:            &window,
		Aggregation:       &aggregation,
		ConditionOperator: conditionOperator,
		ConditionValue:    &conditionValue,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "condition operator") {
		t.Fatal("Expected to fail to validate threshold rule as no condition operator is not valid")
	}
}

func TestShouldFailToValidateThresholdRuleSpecificationWithRollupWhenConditionOperatorIsMissing(t *testing.T) {
	conditionValue := customEventConditionValue
	rollup := customEventRollup
	rule := RuleSpecification{
		DType:          ThresholdRuleType,
		Severity:       SeverityWarning.GetAPIRepresentation(),
		MetricName:     customEventMetricName,
		Rollup:         &rollup,
		ConditionValue: &conditionValue,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "condition operator") {
		t.Fatal("Expected to fail to validate threshold rule as no condition operator is provided")
	}
}

func TestShouldFailToValidateThresholdRuleSpecificationWithRollupWhenConditionOperatorIsNotValid(t *testing.T) {
	conditionOperator := ConditionOperatorType("invalid")
	conditionValue := customEventConditionValue
	rollup := customEventRollup
	rule := RuleSpecification{
		DType:             ThresholdRuleType,
		Severity:          SeverityWarning.GetAPIRepresentation(),
		MetricName:        customEventMetricName,
		Rollup:            &rollup,
		ConditionOperator: conditionOperator,
		ConditionValue:    &conditionValue,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "condition operator") {
		t.Fatal("Expected to fail to validate threshold rule as no condition operator is not valid")
	}
}

func TestShouldConvertSupportedAggregationTypesToSliceOfString(t *testing.T) {
	expectedResult := []string{string(AggregationSum), string(AggregationAvg), string(AggregationMin), string(AggregationMax)}
	result := SupportedAggregationTypes.ToStringSlice()

	if !cmp.Equal(result, expectedResult) {
		t.Fatal("Expected to get slice of strings for supported aggregations")
	}
}

func TestShouldConvertSupportedConditionOperatorTypessToSliceOfString(t *testing.T) {
	expectedResult := []string{string(ConditionOperatorEquals), string(ConditionOperatorNotEqual), string(ConditionOperatorLessThan), string(ConditionOperatorLessThanOrEqual), string(ConditionOperatorGreaterThan), string(ConditionOperatorGreaterThanOrEqual)}
	result := SupportedConditionOperatorTypes.ToStringSlice()

	if !cmp.Equal(result, expectedResult) {
		t.Fatal("Expected to get slice of strings for supported condition operators")
	}
}
