package restapi_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
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

func TestShouldReturnTheProperRespresentationsForSeverityWarning(t *testing.T) {
	assert.Equal(t, 5, SeverityWarning.GetAPIRepresentation())
	assert.Equal(t, "warning", SeverityWarning.GetTerraformRepresentation())
}

func TestShouldReturnTheProperRespresentationsForSeverityCritical(t *testing.T) {
	assert.Equal(t, 10, SeverityCritical.GetAPIRepresentation())
	assert.Equal(t, "critical", SeverityCritical.GetTerraformRepresentation())
}

func TestShouldValidateMinimalCustemEventSpecificationWithSystemRule(t *testing.T) {
	systemRuleId := customEventSystemRuleID
	spec := CustomEventSpecification{
		ID:         customEventID,
		Name:       customEventName,
		EntityType: customEventEntityType,
		Rules:      []RuleSpecification{RuleSpecification{DType: SystemRuleType, SystemRuleID: &systemRuleId}},
	}

	err := spec.Validate()

	assert.Nil(t, err)
	assert.Equal(t, customEventID, spec.GetID())
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

	err := spec.Validate()

	assert.Nil(t, err)
}

func TestFailToValidateCustemEventSpecificationWithSystemRuleWhenIDIsMissing(t *testing.T) {
	systemRuleId := customEventSystemRuleID
	spec := CustomEventSpecification{
		Name:       customEventName,
		EntityType: customEventEntityType,
		Rules:      []RuleSpecification{RuleSpecification{DType: SystemRuleType, SystemRuleID: &systemRuleId}},
	}

	err := spec.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "ID")
}

func TestFailToValidateCustemEventSpecificationWithSystemRuleWhenNameIsMissing(t *testing.T) {
	systemRuleId := customEventSystemRuleID
	spec := CustomEventSpecification{
		ID:         customEventID,
		EntityType: customEventEntityType,
		Rules:      []RuleSpecification{RuleSpecification{DType: SystemRuleType, SystemRuleID: &systemRuleId}},
	}

	err := spec.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "name")
}

func TestFailToValidateCustemEventSpecificationWithSystemRuleWhenEntityTypeIsMissing(t *testing.T) {
	systemRuleId := customEventSystemRuleID
	spec := CustomEventSpecification{
		ID:    customEventID,
		Name:  customEventName,
		Rules: []RuleSpecification{RuleSpecification{DType: SystemRuleType, SystemRuleID: &systemRuleId}},
	}

	err := spec.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "entity type")
}

func TestFailToValidateCustemEventSpecificationWhenNoRuleIsNil(t *testing.T) {
	spec := CustomEventSpecification{
		ID:         customEventID,
		Name:       customEventName,
		EntityType: customEventEntityType,
		Rules:      nil,
	}

	err := spec.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), messagePartExactlyOneRule)
}

func TestFailToValidateCustemEventSpecificationWhenNoRuleIsProvided(t *testing.T) {
	spec := CustomEventSpecification{
		ID:         customEventID,
		Name:       customEventName,
		EntityType: customEventEntityType,
		Rules:      []RuleSpecification{},
	}

	err := spec.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), messagePartExactlyOneRule)
}

func TestFailToValidateCustemEventSpecificationWhenMultipleRulesAreProvided(t *testing.T) {
	systemRuleId := customEventSystemRuleID
	spec := CustomEventSpecification{
		ID:         customEventID,
		Name:       customEventName,
		EntityType: customEventEntityType,
		Rules:      []RuleSpecification{RuleSpecification{DType: SystemRuleType, SystemRuleID: &systemRuleId}, RuleSpecification{DType: SystemRuleType, SystemRuleID: &systemRuleId}},
	}

	err := spec.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), messagePartExactlyOneRule)
}

func TestFailToValidateCustemEventSpecificationWhenRuleTypeIsNotSupported(t *testing.T) {
	spec := CustomEventSpecification{
		ID:         customEventID,
		Name:       customEventName,
		EntityType: customEventEntityType,
		Rules:      []RuleSpecification{RuleSpecification{DType: "invalid"}},
	}

	err := spec.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Unsupported rule type")
}

func TestFailToValidateCustemEventSpecificationWhenTheProvidedRuleIsNotValid(t *testing.T) {
	spec := CustomEventSpecification{
		ID:         customEventID,
		Name:       customEventName,
		EntityType: customEventEntityType,
		Rules:      []RuleSpecification{RuleSpecification{DType: SystemRuleType}},
	}

	err := spec.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "id of system rule")
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

	err := spec.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), messagePartIntegrationId)
}

func TestShouldSuccessfullyValidateSystemRule(t *testing.T) {
	rule := NewSystemRuleSpecification(customEventSystemRuleID, 1000)

	err := rule.Validate()

	assert.Nil(t, err)
}

func TestShouldFailToValidateSystemRuleWhenSystemRuleIDIsMissing(t *testing.T) {
	rule := RuleSpecification{DType: SystemRuleType}

	err := rule.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "id of system rule")
}

func TestShouldFailToValidateSystemRuleWhenRuleTypeIsMissing(t *testing.T) {
	systemRuleId := customEventSystemRuleID
	rule := RuleSpecification{SystemRuleID: &systemRuleId}

	err := rule.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "type of rule")
}

func TestShouldSuccessfullyValidateEventSpecificationDownstream(t *testing.T) {
	downstream := EventSpecificationDownstream{
		IntegrationIds:                []string{"integration-id-1"},
		BroadcastToAllAlertingConfigs: true,
	}

	err := downstream.Validate()

	assert.Nil(t, err)
}

func TestShouldFailToValidateEventSpecificationDownstreamWhenNoIntegrationIDsAreNil(t *testing.T) {
	downstream := EventSpecificationDownstream{
		IntegrationIds: nil,
	}

	err := downstream.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), messagePartIntegrationId)
}

func TestShouldFailToValidateEventSpecificationDownstreamWhenNoIntegrationIDIsProvided(t *testing.T) {
	downstream := EventSpecificationDownstream{
		IntegrationIds: []string{},
	}

	err := downstream.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), messagePartIntegrationId)
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

	err := rule.Validate()

	assert.Nil(t, err)
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

		err := rule.Validate()

		assert.Nil(t, err)
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

		err := rule.Validate()

		assert.Nil(t, err)
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

	err := rule.Validate()

	assert.Nil(t, err)
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

	err := rule.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "metric name")
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

	err := rule.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "metric name")
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

	err := rule.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "either rollup or window")
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

	err := rule.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "either rollup or window")
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

	err := rule.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "aggregation")
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

	err := rule.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "aggregation")
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

	err := rule.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), messagePartConditionOperator)
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

	err := rule.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), messagePartConditionOperator)
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

	err := rule.Validate()

	assert.Nil(t, err)
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

	err := rule.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "matching entity label")
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

	err := rule.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "matching entity label")
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

	err := rule.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "matching entity type")
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

	err := rule.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "matching entity type")
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

	err := rule.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "matching operator")
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

	err := rule.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "matching operator")
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

	err := rule.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "offline duration")
}

func TestShouldConvertSupportedAggregationTypesToSliceOfString(t *testing.T) {
	expectedResult := []string{string(AggregationSum), string(AggregationAvg), string(AggregationMin), string(AggregationMax)}
	result := SupportedAggregationTypes.ToStringSlice()

	assert.Equal(t, expectedResult, result)
}

func TestShouldConvertSupportedConditionOperatorTypesToSliceOfString(t *testing.T) {
	expectedResult := []string{string(ConditionOperatorEquals), string(ConditionOperatorNotEqual), string(ConditionOperatorLessThan), string(ConditionOperatorLessThanOrEqual), string(ConditionOperatorGreaterThan), string(ConditionOperatorGreaterThanOrEqual)}
	result := SupportedConditionOperatorTypes.ToStringSlice()

	assert.Equal(t, expectedResult, result)
}

func TestShouldConvertMatchingOperatorTypesToSliceOfString(t *testing.T) {
	expectedResult := []string{string(MatchingOperatorIs), string(MatchingOperatorContains), string(MatchingOperatorStartsWith), string(MatchingOperatorEndsWith), string(MatchingOperatorNone)}
	result := SupportedMatchingOperatorTypes.ToStringSlice()

	assert.Equal(t, expectedResult, result)
}
