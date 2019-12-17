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

	customEventMatchingEntityLabel = "custom-event-matching-entity-label"
	customEventMatchingEntityType  = "custom-event-matching-entity-type"
	customEventOfflineDuration     = 60000

	valueInvalid = "invalid"

	messagePartExactlyOneRule    = "exactly one rule"
	messagePartIntegrationId     = "integration id"
	messagePartConditionOperator = "condition operator"
)

func TestShouldValidateMinimalCustemEventSpecificationWithSystemRule(t *testing.T) {
	systemRuleId := customEventSystemRuleID
	spec := CustomEventSpecification{
		ID:         customEventID,
		Name:       customEventName,
		EntityType: customEventEntityType,
		Rules:      []RuleSpecification{RuleSpecification{DType: SystemRuleType, SystemRuleID: &systemRuleId}},
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
	systemRuleId := customEventSystemRuleID

	spec := CustomEventSpecification{
		ID:             customEventID,
		Name:           customEventName,
		EntityType:     customEventEntityType,
		Query:          &query,
		Description:    &description,
		ExpirationTime: &expirationTime,
		Triggering:     true,
		Enabled:        true,
		Rules:          []RuleSpecification{RuleSpecification{DType: SystemRuleType, SystemRuleID: &systemRuleId}},
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
	systemRuleId := customEventSystemRuleID
	spec := CustomEventSpecification{
		Name:       customEventName,
		EntityType: customEventEntityType,
		Rules:      []RuleSpecification{RuleSpecification{DType: SystemRuleType, SystemRuleID: &systemRuleId}},
	}

	if err := spec.Validate(); err == nil || !strings.Contains(err.Error(), "ID") {
		t.Fatal("Expected validate to fail as ID is not provided")
	}
}

func TestFailToValidateCustemEventSpecificationWithSystemRuleWhenNameIsMissing(t *testing.T) {
	systemRuleId := customEventSystemRuleID
	spec := CustomEventSpecification{
		ID:         customEventID,
		EntityType: customEventEntityType,
		Rules:      []RuleSpecification{RuleSpecification{DType: SystemRuleType, SystemRuleID: &systemRuleId}},
	}

	if err := spec.Validate(); err == nil || !strings.Contains(err.Error(), "name") {
		t.Fatal("Expected validate to fail as name is not provided")
	}
}

func TestFailToValidateCustemEventSpecificationWithSystemRuleWhenEntityTypeIsMissing(t *testing.T) {
	systemRuleId := customEventSystemRuleID
	spec := CustomEventSpecification{
		ID:    customEventID,
		Name:  customEventName,
		Rules: []RuleSpecification{RuleSpecification{DType: SystemRuleType, SystemRuleID: &systemRuleId}},
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

	if err := spec.Validate(); err == nil || !strings.Contains(err.Error(), messagePartExactlyOneRule) {
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

	if err := spec.Validate(); err == nil || !strings.Contains(err.Error(), messagePartExactlyOneRule) {
		t.Fatal("Expected validate to fail as no rule is provided")
	}
}

func TestFailToValidateCustemEventSpecificationWhenMultipleRulesAreProvided(t *testing.T) {
	systemRuleId := customEventSystemRuleID
	spec := CustomEventSpecification{
		ID:         customEventID,
		Name:       customEventName,
		EntityType: customEventEntityType,
		Rules:      []RuleSpecification{RuleSpecification{DType: SystemRuleType, SystemRuleID: &systemRuleId}, RuleSpecification{DType: SystemRuleType, SystemRuleID: &systemRuleId}},
	}

	if err := spec.Validate(); err == nil || !strings.Contains(err.Error(), messagePartExactlyOneRule) {
		t.Fatal("Expected validation to fail as multiple rules are provided")
	}
}

func TestFailToValidateCustemEventSpecificationWhenRuleTypeIsNotSupported(t *testing.T) {
	spec := CustomEventSpecification{
		ID:         customEventID,
		Name:       customEventName,
		EntityType: customEventEntityType,
		Rules:      []RuleSpecification{RuleSpecification{DType: "invalid"}},
	}

	if err := spec.Validate(); err == nil || !strings.Contains(err.Error(), "Unsupported rule type") {
		t.Fatal("Expected validation to fail as rule type is not supported")
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
	systemRuleId := customEventSystemRuleID
	spec := CustomEventSpecification{
		ID:         customEventID,
		Name:       customEventName,
		EntityType: customEventEntityType,
		Rules:      []RuleSpecification{RuleSpecification{DType: SystemRuleType, SystemRuleID: &systemRuleId}},
		Downstream: &EventSpecificationDownstream{},
	}

	if err := spec.Validate(); err == nil || !strings.Contains(err.Error(), messagePartIntegrationId) {
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
	systemRuleId := customEventSystemRuleID
	rule := RuleSpecification{SystemRuleID: &systemRuleId}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "type of rule") {
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

	if err := downstream.Validate(); err == nil || !strings.Contains(err.Error(), messagePartIntegrationId) {
		t.Fatal("Expected to fail to validate event specification downstream as integration ids is nil")
	}
}

func TestShouldFailToValidateEventSpecificationDownstreamWhenNoIntegrationIDIsProvided(t *testing.T) {
	downstream := EventSpecificationDownstream{
		IntegrationIds: []string{},
	}

	if err := downstream.Validate(); err == nil || !strings.Contains(err.Error(), messagePartIntegrationId) {
		t.Fatal("Expected to fail to validate event specification downstream as no integration id is provided")
	}
}

func TestShouldValidateFullThresholdRuleSpecificationWithWindowRollupAndAggregation(t *testing.T) {
	metricName := customEventMetricName
	aggregation := customEventAggregation
	conditionOperator := customEventConditionOperator
	conditionValue := customEventConditionValue
	window := customEventWindow
	rollup := customEventRollup
	rule := RuleSpecification{
		DType:             ThresholdRuleType,
		Severity:          SeverityWarning.GetAPIRepresentation(),
		MetricName:        &metricName,
		Window:            &window,
		Rollup:            &rollup,
		Aggregation:       &aggregation,
		ConditionOperator: &conditionOperator,
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
	metricName := customEventMetricName
	conditionOperator := customEventConditionOperator
	return func(t *testing.T) {
		conditionValue := customEventConditionValue
		window := customEventWindow
		rule := RuleSpecification{
			DType:             ThresholdRuleType,
			Severity:          SeverityWarning.GetAPIRepresentation(),
			MetricName:        &metricName,
			Window:            &window,
			Aggregation:       &aggregation,
			ConditionOperator: &conditionOperator,
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
		metricName := customEventMetricName
		conditionOperator := customEventConditionOperator
		aggregation := customEventAggregation
		conditionValue := customEventConditionValue
		window := customEventWindow
		rule := RuleSpecification{
			DType:             ThresholdRuleType,
			Severity:          SeverityWarning.GetAPIRepresentation(),
			MetricName:        &metricName,
			Window:            &window,
			Aggregation:       &aggregation,
			ConditionOperator: &conditionOperator,
			ConditionValue:    &conditionValue,
		}

		if err := rule.Validate(); err != nil {
			t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
		}
	}
}

func TestShouldValidateMinimalThresholdRuleSpecificationWithRollup(t *testing.T) {
	metricName := customEventMetricName
	conditionOperator := customEventConditionOperator
	rollup := customEventRollup
	conditionValue := customEventConditionValue
	rule := RuleSpecification{
		DType:             ThresholdRuleType,
		Severity:          SeverityWarning.GetAPIRepresentation(),
		MetricName:        &metricName,
		Rollup:            &rollup,
		ConditionOperator: &conditionOperator,
		ConditionValue:    &conditionValue,
	}

	if err := rule.Validate(); err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}
}

func TestShouldFailToValidateThresholdRuleSpecificationWithWindowWhenMetricNameIsMissing(t *testing.T) {
	conditionOperator := customEventConditionOperator
	aggregation := customEventAggregation
	conditionValue := customEventConditionValue
	window := customEventWindow
	rule := RuleSpecification{
		DType:             ThresholdRuleType,
		Severity:          SeverityWarning.GetAPIRepresentation(),
		Window:            &window,
		Aggregation:       &aggregation,
		ConditionOperator: &conditionOperator,
		ConditionValue:    &conditionValue,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "metric name") {
		t.Fatal("Expected to fail to validate threshold rule as no metric name is provided")
	}
}

func TestShouldFailToValidateThresholdRuleSpecificationWithWindowWhenMetricNameIsBlank(t *testing.T) {
	metricName := ""
	conditionOperator := customEventConditionOperator
	aggregation := customEventAggregation
	conditionValue := customEventConditionValue
	window := customEventWindow
	rule := RuleSpecification{
		DType:             ThresholdRuleType,
		Severity:          SeverityWarning.GetAPIRepresentation(),
		MetricName:        &metricName,
		Window:            &window,
		Aggregation:       &aggregation,
		ConditionOperator: &conditionOperator,
		ConditionValue:    &conditionValue,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "metric name") {
		t.Fatal("Expected to fail to validate threshold rule as no metric name is provided")
	}
}

func TestShouldFailToValidateThresholdRuleSpecificationWhenNeitherRollupNorWindowIsDefined(t *testing.T) {
	metricName := customEventMetricName
	conditionOperator := customEventConditionOperator
	aggregation := customEventAggregation
	conditionValue := customEventConditionValue
	rule := RuleSpecification{
		DType:             ThresholdRuleType,
		Severity:          SeverityWarning.GetAPIRepresentation(),
		MetricName:        &metricName,
		Aggregation:       &aggregation,
		ConditionOperator: &conditionOperator,
		ConditionValue:    &conditionValue,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "either rollup or window") {
		t.Fatal("Expected to fail to validate threshold rule as no window is provided")
	}
}

func TestShouldFailToValidateThresholdRuleSpecificationWithRollupAndWindowAreZero(t *testing.T) {
	metricName := customEventMetricName
	conditionOperator := customEventConditionOperator
	window := 0
	rollup := 0
	aggregation := customEventAggregation
	conditionValue := customEventConditionValue
	rule := RuleSpecification{
		DType:             ThresholdRuleType,
		Severity:          SeverityWarning.GetAPIRepresentation(),
		MetricName:        &metricName,
		Rollup:            &rollup,
		Window:            &window,
		Aggregation:       &aggregation,
		ConditionOperator: &conditionOperator,
		ConditionValue:    &conditionValue,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "either rollup or window") {
		t.Fatal("Expected to fail to validate threshold rule as no window is provided")
	}
}

func TestShouldFailToValidateThresholdRuleSpecificationWithWindowWhenAggregationIsMissingConditionOperator(t *testing.T) {
	metricName := customEventMetricName
	conditionOperator := customEventConditionOperator
	conditionValue := customEventConditionValue
	window := customEventWindow
	rule := RuleSpecification{
		DType:             ThresholdRuleType,
		Severity:          SeverityWarning.GetAPIRepresentation(),
		MetricName:        &metricName,
		Window:            &window,
		ConditionOperator: &conditionOperator,
		ConditionValue:    &conditionValue,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "aggregation") {
		t.Fatal("Expected to fail to validate threshold rule as no aggregation is provided")
	}
}

func TestShouldFailToValidateThresholdRuleSpecificationWithWindowWhenAggregationIsNotValid(t *testing.T) {
	metricName := customEventMetricName
	conditionOperator := customEventConditionOperator
	aggregation := AggregationType(valueInvalid)
	conditionValue := customEventConditionValue
	window := customEventWindow
	rule := RuleSpecification{
		DType:             ThresholdRuleType,
		Severity:          SeverityWarning.GetAPIRepresentation(),
		MetricName:        &metricName,
		Window:            &window,
		Aggregation:       &aggregation,
		ConditionOperator: &conditionOperator,
		ConditionValue:    &conditionValue,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "aggregation") {
		t.Fatal("Expected to fail to validate threshold rule as no aggregation is provided")
	}
}

func TestShouldFailToValidateThresholdRuleSpecificationWithWindowWhenConditionOperatorIsMissing(t *testing.T) {
	metricName := customEventMetricName
	aggregation := customEventAggregation
	conditionValue := customEventConditionValue
	window := customEventWindow
	rule := RuleSpecification{
		DType:          ThresholdRuleType,
		Severity:       SeverityWarning.GetAPIRepresentation(),
		MetricName:     &metricName,
		Window:         &window,
		Aggregation:    &aggregation,
		ConditionValue: &conditionValue,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), messagePartConditionOperator) {
		t.Fatal("Expected to fail to validate threshold rule as no condition operator is provided")
	}
}

func TestShouldFailToValidateThresholdRuleSpecificationWithWindowWhenConditionOperatorIsNotValid(t *testing.T) {
	metricName := customEventMetricName
	aggregation := customEventAggregation
	conditionOperator := ConditionOperatorType(valueInvalid)
	conditionValue := customEventConditionValue
	window := customEventWindow
	rule := RuleSpecification{
		DType:             ThresholdRuleType,
		Severity:          SeverityWarning.GetAPIRepresentation(),
		MetricName:        &metricName,
		Window:            &window,
		Aggregation:       &aggregation,
		ConditionOperator: &conditionOperator,
		ConditionValue:    &conditionValue,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), messagePartConditionOperator) {
		t.Fatal("Expected to fail to validate threshold rule as no condition operator is not valid")
	}
}

func TestShouldValidateEntityVerificationRuleSpecificationWhenAllRequiredFieldsAreProvided(t *testing.T) {
	entityLabel := customEventMatchingEntityLabel
	entityType := customEventMatchingEntityType
	operator := MatchingOperatorIs
	offlineDuration := customEventOfflineDuration
	rule := RuleSpecification{
		DType:               EntityVerificationRuleType,
		Severity:            SeverityWarning.GetAPIRepresentation(),
		MatchingEntityLabel: &entityLabel,
		MatchingEntityType:  &entityType,
		MatchingOperator:    &operator,
		OfflineDuration:     &offlineDuration,
	}

	if err := rule.Validate(); err != nil {
		t.Fatal("Expected to successfully validate entity verification rule")
	}
}

func TestShouldFaileToValidateEntityVerificationRuleSpecificationWhenEntityLabelIsMissing(t *testing.T) {
	entityType := customEventMatchingEntityType
	operator := MatchingOperatorIs
	offlineDuration := customEventOfflineDuration
	rule := RuleSpecification{
		DType:              EntityVerificationRuleType,
		MatchingEntityType: &entityType,
		MatchingOperator:   &operator,
		OfflineDuration:    &offlineDuration,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "matching entity label") {
		t.Fatal("Expected to fail to validate entity verification rule as matching entity label is missing")
	}
}

func TestShouldFaileToValidateEntityVerificationRuleSpecificationWhenEntityLabelIsBlank(t *testing.T) {
	entityLabel := ""
	entityType := customEventMatchingEntityType
	operator := MatchingOperatorIs
	offlineDuration := customEventOfflineDuration
	rule := RuleSpecification{
		DType:               EntityVerificationRuleType,
		MatchingEntityLabel: &entityLabel,
		MatchingEntityType:  &entityType,
		MatchingOperator:    &operator,
		OfflineDuration:     &offlineDuration,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "matching entity label") {
		t.Fatal("Expected to fail to validate entity verification rule as matching entity label is blank")
	}
}

func TestShouldFaileToValidateEntityVerificationRuleSpecificationWhenEntityTypeIsMissing(t *testing.T) {
	entityLabel := customEventMatchingEntityLabel
	operator := MatchingOperatorIs
	offlineDuration := customEventOfflineDuration
	rule := RuleSpecification{
		DType:               EntityVerificationRuleType,
		MatchingEntityLabel: &entityLabel,
		MatchingOperator:    &operator,
		OfflineDuration:     &offlineDuration,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "matching entity type") {
		t.Fatal("Expected to fail to validate entity verification rule as matching entity type is blank")
	}
}

func TestShouldFaileToValidateEntityVerificationRuleSpecificationWhenEntityTypeIsBlank(t *testing.T) {
	entityLabel := customEventMatchingEntityLabel
	entityType := ""
	operator := MatchingOperatorIs
	offlineDuration := customEventOfflineDuration
	rule := RuleSpecification{
		DType:               EntityVerificationRuleType,
		MatchingEntityLabel: &entityLabel,
		MatchingEntityType:  &entityType,
		MatchingOperator:    &operator,
		OfflineDuration:     &offlineDuration,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "matching entity type") {
		t.Fatal("Expected to fail to validate entity verification rule as matching entity type is blank")
	}
}

func TestShouldFaileToValidateEntityVerificationRuleSpecificationWhenMatchingOperatorIsMissing(t *testing.T) {
	entityLabel := customEventMatchingEntityLabel
	entityType := customEventMatchingEntityType
	offlineDuration := customEventOfflineDuration
	rule := RuleSpecification{
		DType:               EntityVerificationRuleType,
		MatchingEntityType:  &entityType,
		MatchingEntityLabel: &entityLabel,
		OfflineDuration:     &offlineDuration,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "matching operator") {
		t.Fatal("Expected to fail to validate entity verification rule as matching operator is missing")
	}
}

func TestShouldFaileToValidateEntityVerificationRuleSpecificationWhenMatchingOpertatorIsNotSupported(t *testing.T) {
	entityLabel := customEventMatchingEntityLabel
	entityType := customEventMatchingEntityType
	operator := MatchingOperatorType("Invalid")
	offlineDuration := customEventOfflineDuration
	rule := RuleSpecification{
		DType:               EntityVerificationRuleType,
		MatchingEntityLabel: &entityLabel,
		MatchingEntityType:  &entityType,
		MatchingOperator:    &operator,
		OfflineDuration:     &offlineDuration,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "matching operator") {
		t.Fatal("Expected to fail to validate entity verification rule as matching operator is not supported")
	}
}

func TestShouldFaileToValidateEntityVerificationRuleSpecificationWhenOfflineDurationIsNotSupported(t *testing.T) {
	entityLabel := customEventMatchingEntityLabel
	entityType := customEventMatchingEntityType
	operator := MatchingOperatorIs
	rule := RuleSpecification{
		DType:               EntityVerificationRuleType,
		MatchingEntityLabel: &entityLabel,
		MatchingEntityType:  &entityType,
		MatchingOperator:    &operator,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "offline duration") {
		t.Fatal("Expected to fail to validate entity verification rule as offline duration is missing")
	}
}

func TestShouldConvertSupportedAggregationTypesToSliceOfString(t *testing.T) {
	expectedResult := []string{string(AggregationSum), string(AggregationAvg), string(AggregationMin), string(AggregationMax)}
	result := SupportedAggregationTypes.ToStringSlice()

	if !cmp.Equal(result, expectedResult) {
		t.Fatal("Expected to get slice of strings for supported aggregations")
	}
}

func TestShouldConvertSupportedConditionOperatorTypesToSliceOfString(t *testing.T) {
	expectedResult := []string{string(ConditionOperatorEquals), string(ConditionOperatorNotEqual), string(ConditionOperatorLessThan), string(ConditionOperatorLessThanOrEqual), string(ConditionOperatorGreaterThan), string(ConditionOperatorGreaterThanOrEqual)}
	result := SupportedConditionOperatorTypes.ToStringSlice()

	if !cmp.Equal(result, expectedResult) {
		t.Fatal("Expected to get slice of strings for supported condition operators")
	}
}

func TestShouldConvertMatchingOperatorTypesToSliceOfString(t *testing.T) {
	expectedResult := []string{string(MatchingOperatorIs), string(MatchingOperatorContains), string(MatchingOperatorStartsWith), string(MatchingOperatorEndsWith), string(MatchingOperatorNone)}
	result := SupportedMatchingOperatorTypes.ToStringSlice()

	if !cmp.Equal(result, expectedResult) {
		t.Fatal("Expected to get slice of strings for supported matching operators")
	}
}
