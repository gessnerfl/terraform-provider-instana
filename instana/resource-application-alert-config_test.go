package instana_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
)

func TestShouldReturnTrueWhenCheckingForSchemaDiffSuppressForTagFilterOfApplicationAlertConfigAndValueCanBeNormalizedAndOldAndNewNormalizedValueAreEqual(t *testing.T) {
	resourceHandle := NewApplicationAlertConfigResourceHandle()
	schema := resourceHandle.MetaData().Schema
	oldValue := expressionEntityTypeDestEqValue
	newValue := "entity.type  EQUALS    'foo'"

	require.True(t, schema[ApplicationAlertConfigFieldTagFilter].DiffSuppressFunc(ApplicationAlertConfigFieldTagFilter, oldValue, newValue, nil))
}

func TestShouldReturnFalseWhenCheckingForSchemaDiffSuppressForTagFilterOfApplicationAlertConfigAndValueCanBeNormalizedAndOldAndNewNormalizedValueAreNotEqual(t *testing.T) {
	resourceHandle := NewApplicationAlertConfigResourceHandle()
	schema := resourceHandle.MetaData().Schema
	oldValue := expressionEntityTypeSrcEqValue
	newValue := validTagFilter

	require.False(t, schema[ApplicationAlertConfigFieldTagFilter].DiffSuppressFunc(ApplicationAlertConfigFieldTagFilter, oldValue, newValue, nil))
}

func TestShouldReturnTrueWhenCheckingForSchemaDiffSuppressForTagFilterOfApplicationAlertConfigAndValueCannotBeNormalizedAndOldAndNewValueAreEqual(t *testing.T) {
	resourceHandle := NewApplicationAlertConfigResourceHandle()
	schema := resourceHandle.MetaData().Schema
	invalidValue := invalidTagFilter

	require.True(t, schema[ApplicationAlertConfigFieldTagFilter].DiffSuppressFunc(ApplicationAlertConfigFieldTagFilter, invalidValue, invalidValue, nil))
}

func TestShouldReturnFalseWhenCheckingForSchemaDiffSuppressForTagFilterOfApplicationAlertConfigAndValueCannotBeNormalizedAndOldAndNewValueAreNotEqual(t *testing.T) {
	resourceHandle := NewApplicationAlertConfigResourceHandle()
	schema := resourceHandle.MetaData().Schema
	oldValue := invalidTagFilter
	newValue := "entity.type foo foo foo"

	require.False(t, schema[ApplicationAlertConfigFieldTagFilter].DiffSuppressFunc(ApplicationAlertConfigFieldTagFilter, oldValue, newValue, nil))
}

func TestShouldReturnNormalizedValueForTagFilterOfApplicationAlertConfigWhenStateFuncIsCalledAndValueCanBeNormalized(t *testing.T) {
	resourceHandle := NewApplicationAlertConfigResourceHandle()
	schema := resourceHandle.MetaData().Schema
	expectedValue := expressionEntityTypeDestEqValue
	newValue := validTagFilter

	require.Equal(t, expectedValue, schema[ApplicationAlertConfigFieldTagFilter].StateFunc(newValue))
}

func TestShouldReturnProvidedValueForTagFilterOfApplicationAlertConfigWhenStateFuncIsCalledAndValueCannotBeNormalized(t *testing.T) {
	resourceHandle := NewApplicationAlertConfigResourceHandle()
	schema := resourceHandle.MetaData().Schema
	value := invalidTagFilter

	require.Equal(t, value, schema[ApplicationAlertConfigFieldTagFilter].StateFunc(value))
}

func TestShouldReturnNoErrorsAndWarningsWhenValidationOfTagFilterOfApplicationAlertConfigIsCalledAndValueCanBeParsed(t *testing.T) {
	resourceHandle := NewApplicationAlertConfigResourceHandle()
	schema := resourceHandle.MetaData().Schema
	value := validTagFilter

	warns, errs := schema[ApplicationAlertConfigFieldTagFilter].ValidateFunc(value, ApplicationAlertConfigFieldTagFilter)
	require.Empty(t, warns)
	require.Empty(t, errs)
}

func TestShouldReturnOneErrorAndNoWarningsWhenValidationOfTagFilterOfApplicationAlertConfigIsCalledAndValueCannotBeParsed(t *testing.T) {
	resourceHandle := NewApplicationAlertConfigResourceHandle()
	schema := resourceHandle.MetaData().Schema
	value := invalidTagFilter

	warns, errs := schema[ApplicationAlertConfigFieldTagFilter].ValidateFunc(value, ApplicationAlertConfigFieldTagFilter)
	require.Empty(t, warns)
	require.Len(t, errs, 1)
}

func TestApplicationAlertConfigResourceShouldHaveSchemaVersionZero(t *testing.T) {
	require.Equal(t, 0, NewApplicationAlertConfigResourceHandle().MetaData().SchemaVersion)
}

func TestApplicationConfigResourceShouldHaveNoStateUpgrader(t *testing.T) {
	resourceHandler := NewApplicationAlertConfigResourceHandle()

	require.Empty(t, resourceHandler.StateUpgraders())
}
