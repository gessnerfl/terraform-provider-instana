package instana_test

import (
	"fmt"
	"github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
	"testing"
)

const tagFilterExpressionKey = "entity.type EQUALS 'foo'"
const validTagFilterExpression = "entity.type EQUALS 'foo'"
const invalidTagFilterExpression = "entity.type bla bla bla"

func TestTagFilterSchema(t *testing.T) {
	for k, tagFilterExpressionSchema := range map[string]*schema.Schema{"optional": instana.OptionalTagFilterExpressionSchema, "required": instana.RequiredTagFilterExpressionSchema} {
		t.Run(fmt.Sprintf("DiffSuppressFunc of %s TagFilterExpression Schema should return true when value can be normalized and old and new normalized value are equal", k), createTestOfDiffSuppressFuncOfTagFilterShouldReturnTrueWhenValueCanBeNormalizedAndOldAndNewNormalizedValueAreEqual(tagFilterExpressionSchema))
		t.Run(fmt.Sprintf("DiffSuppressFunc of %s TagFilterExpression Schema should return false when value can be normalized and old and new normalized value are not equal", k), createTestOfDiffSuppressFuncOfTagFilterShouldReturnFalseWhenValueCanBeNormalizedAndOldAndNewNormalizedValueAreNotEqual(tagFilterExpressionSchema))
		t.Run(fmt.Sprintf("DiffSuppressFunc of %s TagFilterExpression Schema should return true when value can be normalized and old and new value are equal", k), createTestOfDiffSuppressFuncOfTagFilterShouldReturnTrueWhenValueCannotBeNormalizedAndOldAndNewValueAreEqual(tagFilterExpressionSchema))
		t.Run(fmt.Sprintf("DiffSuppressFunc of %s TagFilterExpression Schema should return false when value cannot be normalized and old and new value are not equal", k), createTestOfDiffSuppressFuncOfTagFilterShouldReturnFalseWhenValueCannotBeNormalizedAndOldAndNewValueAreNotEqual(tagFilterExpressionSchema))
		t.Run(fmt.Sprintf("StateFunc of %s TagFilterExpression Schema should return normalized value when value can be normalized", k), createTestOfStateFuncOfTagFilterShouldReturnNormalizedValueWhenValueCanBeNormalized(tagFilterExpressionSchema))
		t.Run(fmt.Sprintf("StateFunc of %s TagFilterExpression Schema should return provided value when value cannot be normalized", k), createTestOfStateFuncOfTagFilterShouldReturnProvidedValueWhenValueCannotBeNormalized(tagFilterExpressionSchema))
		t.Run(fmt.Sprintf("ValidateFunc of %s TagFilterExpression Schema should return no errors and warnings when value can be parsed", k), createTestOfValidateFuncOfTagFilterShouldReturnNoErrorsAndWarningsWhenValueCanBeParsed(tagFilterExpressionSchema))
		t.Run(fmt.Sprintf("ValidateFunc of %s TagFilterExpression Schema should return one error and no warnings when value can be parsed", k), createTestOfValidateFuncOfTagFilterShouldReturnOneErrorAndNoWarningsWhenValueCannotBeParsed(tagFilterExpressionSchema))
	}
}

func createTestOfDiffSuppressFuncOfTagFilterShouldReturnTrueWhenValueCanBeNormalizedAndOldAndNewNormalizedValueAreEqual(tagFilterSchema *schema.Schema) func(t *testing.T) {
	return func(t *testing.T) {
		oldValue := expressionEntityTypeDestEqValue
		newValue := "entity.type  EQUALS    'foo'"

		require.True(t, tagFilterSchema.DiffSuppressFunc(tagFilterExpressionKey, oldValue, newValue, nil))
	}
}

func createTestOfDiffSuppressFuncOfTagFilterShouldReturnFalseWhenValueCanBeNormalizedAndOldAndNewNormalizedValueAreNotEqual(tagFilterSchema *schema.Schema) func(t *testing.T) {
	return func(t *testing.T) {
		oldValue := expressionEntityTypeSrcEqValue
		newValue := validTagFilterExpression

		require.False(t, tagFilterSchema.DiffSuppressFunc(tagFilterExpressionKey, oldValue, newValue, nil))
	}
}

func createTestOfDiffSuppressFuncOfTagFilterShouldReturnTrueWhenValueCannotBeNormalizedAndOldAndNewValueAreEqual(tagFilterSchema *schema.Schema) func(t *testing.T) {
	return func(t *testing.T) {
		invalidValue := invalidTagFilterExpression

		require.True(t, tagFilterSchema.DiffSuppressFunc(tagFilterExpressionKey, invalidValue, invalidValue, nil))
	}
}

func createTestOfDiffSuppressFuncOfTagFilterShouldReturnFalseWhenValueCannotBeNormalizedAndOldAndNewValueAreNotEqual(tagFilterSchema *schema.Schema) func(t *testing.T) {
	return func(t *testing.T) {
		oldValue := invalidTagFilterExpression
		newValue := "entity.type foo foo foo"

		require.False(t, tagFilterSchema.DiffSuppressFunc(tagFilterExpressionKey, oldValue, newValue, nil))
	}
}

func createTestOfStateFuncOfTagFilterShouldReturnNormalizedValueWhenValueCanBeNormalized(tagFilterSchema *schema.Schema) func(t *testing.T) {
	return func(t *testing.T) {
		expectedValue := expressionEntityTypeDestEqValue
		newValue := validTagFilterExpression

		require.Equal(t, expectedValue, tagFilterSchema.StateFunc(newValue))
	}
}

func createTestOfStateFuncOfTagFilterShouldReturnProvidedValueWhenValueCannotBeNormalized(tagFilterSchema *schema.Schema) func(t *testing.T) {
	return func(t *testing.T) {
		value := invalidTagFilterExpression

		require.Equal(t, value, tagFilterSchema.StateFunc(value))
	}
}

func createTestOfValidateFuncOfTagFilterShouldReturnNoErrorsAndWarningsWhenValueCanBeParsed(tagFilterSchema *schema.Schema) func(t *testing.T) {
	return func(t *testing.T) {
		value := validTagFilterExpression

		warns, errs := tagFilterSchema.ValidateFunc(value, tagFilterExpressionKey)
		require.Empty(t, warns)
		require.Empty(t, errs)
	}
}

func createTestOfValidateFuncOfTagFilterShouldReturnOneErrorAndNoWarningsWhenValueCannotBeParsed(tagFilterSchema *schema.Schema) func(t *testing.T) {
	return func(t *testing.T) {
		value := invalidTagFilterExpression

		warns, errs := tagFilterSchema.ValidateFunc(value, tagFilterExpressionKey)
		require.Empty(t, warns)
		require.Len(t, errs, 1)
	}
}
