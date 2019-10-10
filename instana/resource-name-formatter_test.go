package instana_test

import (
	"reflect"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
)

func TestShouldCreateInstanceOfNoopResourceNameFormatterWhenFormattingIsNotActivated(t *testing.T) {
	inst := NewResourceNameFormatter(false)

	if instType := reflect.TypeOf(inst).String(); instType != "*instana.noopResourceNameFormatter" {
		t.Fatalf("Instance should be of type noopResourceNameFormatter but was %s", instType)
	}
}

func TestShouldCreateInstanceOfTerraformManagedResourceNameFormatterWhenFormattingIsActivated(t *testing.T) {
	inst := NewResourceNameFormatter(true)

	if instType := reflect.TypeOf(inst).String(); instType != "*instana.terraformManagedResourceNameFormatter" {
		t.Fatalf("Instance should be of type terraformManagedResourceNameFormatter but was %s", instType)
	}
}

func TestShouldNotAppendSuffixToNameWhenNoopResourceNameFormatterIsUsed(t *testing.T) {
	input := "Test String"
	inst := NewResourceNameFormatter(false)

	result := inst.Format(input)

	if input != result {
		t.Fatalf("Input should not be modified by noopResourceNameFormatter: %s vs %s", input, result)
	}
}

func TestShouldNotRemoveSuffixFromNameWhenNoopResourceNameFormatterIsUsed(t *testing.T) {
	input := "Test String" + TerraformManagedResourceNameSuffix
	inst := NewResourceNameFormatter(false)

	result := inst.UndoFormat(input)

	if input != result {
		t.Fatalf("Input should not be modified by noopResourceNameFormatter: %s vs %s", input, result)
	}
}

func TestShouldAppendSuffixToNameWhenTerraformManagedResourceNameFormatterIsUsed(t *testing.T) {
	input := "Test String"
	expectedResult := input + TerraformManagedResourceNameSuffix
	inst := NewResourceNameFormatter(true)

	result := inst.Format(input)

	if expectedResult != result {
		t.Fatalf("Input string should be appended by terraformManagedResourceNameFormatter when formatting name: %s vs %s", expectedResult, result)
	}
}

func TestShouldRemoveSuffixFromNameWhenTerraformManagedResourceNameFormatterIsUsed(t *testing.T) {
	expectedResult := "Test String"
	input := expectedResult + TerraformManagedResourceNameSuffix
	inst := NewResourceNameFormatter(true)

	result := inst.UndoFormat(input)

	if expectedResult != result {
		t.Fatalf("Suffix string should be trimmed by terraformManagedResourceNameFormatter when undoing format of name: %s vs %s", expectedResult, result)
	}
}
