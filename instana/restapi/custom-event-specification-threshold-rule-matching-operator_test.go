package restapi_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

func TestShouldPovideInstanaSupportedValuesOfMatchingOperatorTypes(t *testing.T) {
	expectedResult := []string{MatchingOperatorIs.InstanaAPIValue(), MatchingOperatorContains.InstanaAPIValue(), MatchingOperatorStartsWith.InstanaAPIValue(), MatchingOperatorEndsWith.InstanaAPIValue()}
	result := SupportedMatchingOperators.InstanaAPISupportedValues()

	assert.Equal(t, expectedResult, result)
}

func TestShouldPovideTerraformSupportedValuesOfMatchingOperatorTypes(t *testing.T) {
	expectedResult := []string{}
	for _, mo := range SupportedMatchingOperators {
		expectedResult = append(expectedResult, mo.TerraformSupportedValues()...)
	}
	result := SupportedMatchingOperators.TerrafromSupportedValues()

	assert.Equal(t, expectedResult, result)
}

func TestShouldReturnTheMatchingOperatorTypeForAllSupportedInstanaWebRestAPIValues(t *testing.T) {
	for _, mo := range SupportedMatchingOperators {
		t.Run(fmt.Sprintf("TestShouldReturnTheMatchingOperatorTypeForAllSupportedInstanaWebRestAPIValues%s", mo), createTestCaseForTestShouldReturnTheMatchingOperatorTypeForAllSupportedInstanaWebRestAPIValues(mo))
	}
}

func createTestCaseForTestShouldReturnTheMatchingOperatorTypeForAllSupportedInstanaWebRestAPIValues(mo MatchingOperator) func(*testing.T) {
	return func(t *testing.T) {
		val, err := SupportedMatchingOperators.FromInstanaAPIValue(mo.InstanaAPIValue())

		assert.Nil(t, err)
		assert.Equal(t, mo, val)
	}
}

func TestShouldReturnErrorWhenTheMatchingOperatorTypeIsNotASupportedInstanaWebRestAPIValue(t *testing.T) {
	val, err := SupportedMatchingOperators.FromInstanaAPIValue("invalid")

	assert.NotNil(t, err)
	assert.Equal(t, MatchingOperatorIs, val)
}

func TestShouldReturnTheMatchingOperatorTypeForAllSupportedTerraformProviderValues(t *testing.T) {
	for _, mo := range SupportedMatchingOperators {
		for _, v := range mo.TerraformSupportedValues() {
			t.Run(fmt.Sprintf("TestShouldReturnTheMatchingOperatorTypeForAllSupportedTerraformProviderValues%s", v), createTestCaseForTestShouldReturnTheMatchingOperatorTypeForAllSupportedTerraformProviderValues(mo, v))
		}
	}
}

func createTestCaseForTestShouldReturnTheMatchingOperatorTypeForAllSupportedTerraformProviderValues(mo MatchingOperator, value string) func(*testing.T) {
	return func(t *testing.T) {
		val, err := SupportedMatchingOperators.FromTerraformValue(value)

		assert.Nil(t, err)
		assert.Equal(t, mo, val)
	}
}

func TestShouldReturnErrorWhenTheMatchingOperatorTypeIsNotASupportedInstanaTerraformProviderValue(t *testing.T) {
	val, err := SupportedMatchingOperators.FromTerraformValue("invalid")

	assert.NotNil(t, err)
	assert.Equal(t, MatchingOperatorIs, val)
}
