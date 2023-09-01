package restapi_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

const (
	customEventMetricName     = "threshold-rule-metric-name"
	customEventRollup         = 40000
	customEventConditionValue = 1.2

	customEventMatchingEntityLabel = "custom-event-matching-entity-label"
	customEventMatchingEntityType  = "custom-event-matching-entity-type"
	customEventOfflineDuration     = 60000

	valueInvalid = "invalid"
)

func TestShouldReturnConditionOperatorTypeOfThresholdRuleWhenValidInstanaWebRestAPIConditionOperatorTypoIsProvided(t *testing.T) {
	metricName := customEventMetricName
	conditionOperator := ConditionOperatorEquals.InstanaAPIValue()
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

	val, err := rule.ConditionOperatorType()

	assert.Nil(t, err)
	assert.NotNil(t, val)
	assert.Equal(t, ConditionOperatorEquals, val)
}

func TestShouldReturnErrorWheConditionOperatorTypeOfThresholdnRuleIsNotAValidInstanaWebRestAPIValue(t *testing.T) {
	metricName := customEventMetricName
	conditionOperator := "invalid"
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

	val, err := rule.ConditionOperatorType()

	assert.NotNil(t, err)
	assert.Nil(t, val)
}

func TestShouldReturnNilWhenRuleSpecificationDoesNotSupportConditionOperatorTypesAndConditionOperatorTypeIsRequested(t *testing.T) {
	entityLabel := customEventMatchingEntityLabel
	entityType := customEventMatchingEntityType
	operator := MatchingOperatorIs.InstanaAPIValue()
	offlineDuration := customEventOfflineDuration
	rule := RuleSpecification{
		DType:               EntityVerificationRuleType,
		Severity:            SeverityWarning.GetAPIRepresentation(),
		MatchingEntityLabel: &entityLabel,
		MatchingEntityType:  &entityType,
		MatchingOperator:    &operator,
		OfflineDuration:     &offlineDuration,
	}

	val, err := rule.ConditionOperatorType()

	assert.Nil(t, err)
	assert.Nil(t, val)
}

func TestShouldReturnMatchingOperatorTypeForEntityVerificationRuleWhenValidInstanaWebRestAPIMatchingOperatorTypoIsProvided(t *testing.T) {
	entityLabel := customEventMatchingEntityLabel
	entityType := customEventMatchingEntityType
	operator := MatchingOperatorIs.InstanaAPIValue()
	offlineDuration := customEventOfflineDuration
	rule := RuleSpecification{
		DType:               EntityVerificationRuleType,
		Severity:            SeverityWarning.GetAPIRepresentation(),
		MatchingEntityLabel: &entityLabel,
		MatchingEntityType:  &entityType,
		MatchingOperator:    &operator,
		OfflineDuration:     &offlineDuration,
	}

	val, err := rule.MatchingOperatorType()

	assert.Nil(t, err)
	assert.NotNil(t, val)
	assert.Equal(t, MatchingOperatorIs, val)
}

func TestShouldReturnErrorWhenMatchingOperatorTypeForEntityVerificationRuleIsNotAValidInstanaWebRestAPIValue(t *testing.T) {
	entityLabel := customEventMatchingEntityLabel
	entityType := customEventMatchingEntityType
	operator := "invalid"
	offlineDuration := customEventOfflineDuration
	rule := RuleSpecification{
		DType:               EntityVerificationRuleType,
		Severity:            SeverityWarning.GetAPIRepresentation(),
		MatchingEntityLabel: &entityLabel,
		MatchingEntityType:  &entityType,
		MatchingOperator:    &operator,
		OfflineDuration:     &offlineDuration,
	}

	val, err := rule.MatchingOperatorType()

	assert.NotNil(t, err)
	assert.Nil(t, val)
}

func TestShouldReturnNilWhenRuleSpecificationDoesNotSupportMatchingOperatorTypesAndMatchingOperatorTypeIsRequested(t *testing.T) {
	metricName := customEventMetricName
	conditionOperator := ConditionOperatorEquals.InstanaAPIValue()
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

	val, err := rule.MatchingOperatorType()

	assert.Nil(t, err)
	assert.Nil(t, val)
}

func TestShouldConvertSupportedAggregationTypesToSliceOfString(t *testing.T) {
	expectedResult := []string{string(AggregationSum), string(AggregationAvg), string(AggregationMin), string(AggregationMax)}
	result := SupportedAggregationTypes.ToStringSlice()

	assert.Equal(t, expectedResult, result)
}

func TestShouldReturnTrueForAllSupportedMetricPatternTypes(t *testing.T) {
	for _, ot := range SupportedMetricPatternOperatorTypes {
		t.Run(fmt.Sprintf("TestShouldReturnTrueForAllSupportedMetricPatternTypes%s", ot), createTestCaseForVerifyingIfAIsSupportedReturnsTrueForAllSupportedMetricPatternTypes(ot))
	}
}

func createTestCaseForVerifyingIfAIsSupportedReturnsTrueForAllSupportedMetricPatternTypes(ot MetricPatternOperatorType) func(*testing.T) {
	return func(t *testing.T) {
		assert.True(t, SupportedMetricPatternOperatorTypes.IsSupported(ot))
	}
}

func TestShouldReturnFalseWhenMetricPatternOperatorTypeIsNotSupported(t *testing.T) {
	assert.False(t, SupportedMetricPatternOperatorTypes.IsSupported("invalid"))
}

func TestShouldConvertMetricPatternOperatorTypesToStringSlice(t *testing.T) {
	assert.Equal(t, []string{"is", "contains", "any", "startsWith", "endsWith"}, SupportedMetricPatternOperatorTypes.ToStringSlice())
}
