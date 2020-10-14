package restapi_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

func TestShouldPovideInstanaSupportedValuesOfConditionOperatorTypes(t *testing.T) {
	expectedResult := []string{ConditionOperatorEquals.InstanaAPIValue(), ConditionOperatorNotEqual.InstanaAPIValue(), ConditionOperatorLessThan.InstanaAPIValue(), ConditionOperatorLessThanOrEqual.InstanaAPIValue(), ConditionOperatorGreaterThan.InstanaAPIValue(), ConditionOperatorGreaterThanOrEqual.InstanaAPIValue()}
	result := SupportedConditionOperators.InstanaAPISupportedValues()

	assert.Equal(t, expectedResult, result)
}

func TestShouldPovideTerraformSupportedValuesOfConditionOperatorTypes(t *testing.T) {
	expectedResult := []string{}
	for _, op := range SupportedConditionOperators {
		expectedResult = append(expectedResult, op.TerraformSupportedValues()...)
	}
	result := SupportedConditionOperators.TerrafromSupportedValues()

	assert.Equal(t, expectedResult, result)
}

func TestShouldReturnTheConditionOperatorTypeForAllSupportedInstanaWebRestAPIValues(t *testing.T) {
	for _, op := range SupportedConditionOperators {
		t.Run(fmt.Sprintf("TestShouldReturnTheConditionOperatorTypeForAllSupportedInstanaWebRestAPIValues%s", op), createTestCaseForTestShouldReturnTheConditionOperatorTypeForAllSupportedInstanaWebRestAPIValues(op))
	}
}

func createTestCaseForTestShouldReturnTheConditionOperatorTypeForAllSupportedInstanaWebRestAPIValues(op ConditionOperator) func(*testing.T) {
	return func(t *testing.T) {
		val, err := SupportedConditionOperators.FromInstanaAPIValue(op.InstanaAPIValue())

		assert.Nil(t, err)
		assert.Equal(t, op, val)
	}
}

func TestShouldReturnErrorWhenTheConditionOperatorTypeIsNotASupportedInstanaWebRestAPIValue(t *testing.T) {
	val, err := SupportedConditionOperators.FromInstanaAPIValue("invalid")

	assert.NotNil(t, err)
	assert.Equal(t, ConditionOperatorEquals, val)
}

func TestShouldReturnTheConditionOperatorTypeForAllSupportedTerraformProviderValues(t *testing.T) {
	for _, op := range SupportedConditionOperators {
		for _, v := range op.TerraformSupportedValues() {
			t.Run(fmt.Sprintf("TestShouldReturnTheConditionOperatorTypeForAllSupportedTerraformProviderValues%s", v), createTestCaseForTestShouldReturnTheConditionOperatorTypeForAllSupportedTerraformProviderValues(op, v))
		}
	}
}

func createTestCaseForTestShouldReturnTheConditionOperatorTypeForAllSupportedTerraformProviderValues(op ConditionOperator, value string) func(*testing.T) {
	return func(t *testing.T) {
		val, err := SupportedConditionOperators.FromTerraformValue(value)

		assert.Nil(t, err)
		assert.Equal(t, op, val)
	}
}

func TestShouldReturnErrorWhenTheConditionOperatorTypeIsNotASupportedInstanaTerraformProviderValue(t *testing.T) {
	val, err := SupportedConditionOperators.FromTerraformValue("invalid")

	assert.NotNil(t, err)
	assert.Equal(t, ConditionOperatorEquals, val)
}
