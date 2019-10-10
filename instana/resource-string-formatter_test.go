package instana_test

import (
	"reflect"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
)

func TestShouldCreateInstanceOfNoopResourceStringFormatterWhenFormattingIsNotActivated(t *testing.T) {
	inst := NewResourceStringFormatter(false)

	if instType := reflect.TypeOf(inst).String(); instType != "*instana.noopResourceStringFormatter" {
		t.Fatalf("Instance should be of type noopResourceStringFormatter but was %s", instType)
	}
}

func TestShouldCreateInstanceOfTerraformManagedResourceStringFormatterWhenFormattingIsActivated(t *testing.T) {
	inst := NewResourceStringFormatter(true)

	if instType := reflect.TypeOf(inst).String(); instType != "*instana.terraformManagedResourceStringFormatter" {
		t.Fatalf("Instance should be of type terraformManagedResourceStringFormatter but was %s", instType)
	}
}

func TestShouldNotAppendSuffixToNameWhenNoopResourceStringFormatterIsUsed(t *testing.T) {
	input := "Test String"
	inst := NewResourceStringFormatter(false)

	result := inst.FormatName(input)

	if input != result {
		t.Fatalf("Input should not be modified by noopResourceStringFormatter: %s vs %s", input, result)
	}
}

func TestShouldNotRemoveSuffixFromNameWhenNoopResourceStringFormatterIsUsed(t *testing.T) {
	input := "Test String" + TerraformManagedResourceNameSuffix
	inst := NewResourceStringFormatter(false)

	result := inst.UndoFormatName(input)

	if input != result {
		t.Fatalf("Input should not be modified by noopResourceStringFormatter: %s vs %s", input, result)
	}
}

func TestShouldNotAppendSuffixToDescriptionWhenNoopResourceStringFormatterIsUsed(t *testing.T) {
	input := "Test String"
	inst := NewResourceStringFormatter(false)

	result := inst.FormatDescription(input)

	if input != result {
		t.Fatalf("Input should not be modified by noopResourceStringFormatter: %s vs %s", input, result)
	}
}

func TestShouldNotRemoveSuffixFromDescriptionWhenNoopResourceStringFormatterIsUsed(t *testing.T) {
	input := "Test String" + TerraformManagedResourceDescriptionSuffix
	inst := NewResourceStringFormatter(false)

	result := inst.UndoFormatDescription(input)

	if input != result {
		t.Fatalf("Input should not be modified by noopResourceStringFormatter: %s vs %s", input, result)
	}
}

func TestShouldAppendSuffixToNameWhenTerraformManagedResourceStringFormatterIsUsed(t *testing.T) {
	input := "Test String"
	expectedResult := input + TerraformManagedResourceNameSuffix
	inst := NewResourceStringFormatter(true)

	result := inst.FormatName(input)

	if expectedResult != result {
		t.Fatalf("Input string should be appended by terraformManagedResourceStringFormatter when formatting name: %s vs %s", expectedResult, result)
	}
}

func TestShouldRemoveSuffixFromNameWhenTerraformManagedResourceStringFormatterIsUsed(t *testing.T) {
	expectedResult := "Test String"
	input := expectedResult + TerraformManagedResourceNameSuffix
	inst := NewResourceStringFormatter(true)

	result := inst.UndoFormatName(input)

	if expectedResult != result {
		t.Fatalf("Suffix string should be trimmed by terraformManagedResourceStringFormatter when undoing format of name: %s vs %s", expectedResult, result)
	}
}

func TestShouldAppendSuffixToDescriptionWhenTerraformManagedResourceStringFormatterIsUsed(t *testing.T) {
	input := "Test String"
	expectedResult := input + TerraformManagedResourceDescriptionSuffix
	inst := NewResourceStringFormatter(true)

	result := inst.FormatDescription(input)

	if expectedResult != result {
		t.Fatalf("String should be appended by terraformManagedResourceStringFormatter when formatting description: %s vs %s", expectedResult, result)
	}
}

func TestShouldRemoveSuffixFromDescriptionWhenTerraformManagedResourceStringFormatterIsUsed(t *testing.T) {
	expectedResult := "Test String"
	input := expectedResult + TerraformManagedResourceDescriptionSuffix
	inst := NewResourceStringFormatter(true)

	result := inst.UndoFormatDescription(input)

	if expectedResult != result {
		t.Fatalf("Suffix string should be trimmed by terraformManagedResourceStringFormatter when undoing format of description: %s vs %s", expectedResult, result)
	}
}
