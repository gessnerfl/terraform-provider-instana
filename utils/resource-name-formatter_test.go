package utils_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/utils"
)

const input = "Test String"
const prefix = "prefix"
const suffix = "suffix"

func TestShouldAppendPrefixAndSuffixToName(t *testing.T) {
	expectedResult := prefix + " " + input + " " + suffix
	inst := NewResourceNameFormatter(prefix, suffix)

	result := inst.Format(input)

	if expectedResult != result {
		t.Fatalf("Prefix and suffix should be appended when formatting name: %s vs %s", expectedResult, result)
	}
}

func TestShouldRemoveSuffixFromNameWhenTerraformManagedResourceNameFormatterIsUsed(t *testing.T) {
	expectedResult := "Test String"
	input := prefix + " " + input + " " + suffix
	inst := NewResourceNameFormatter(prefix, suffix)

	result := inst.UndoFormat(input)

	if expectedResult != result {
		t.Fatalf("Prefix and suffix should be trimmed when undoing format of name: %s vs %s", expectedResult, result)
	}
}

func TestShouldNotRemovePrefixStringWhenPrefixStringAppearsInTheMiddle(t *testing.T) {
	expectedResult := "Test prefix String"
	input := expectedResult
	inst := NewResourceNameFormatter(prefix, suffix)

	result := inst.UndoFormat(input)

	if expectedResult != result {
		t.Fatalf("Prefix should only be trimmed at the beginning: %s vs %s", expectedResult, result)
	}
}

func TestShouldNotRemoveSuffixStringWhenSuffixStringAppearsInTheMiddle(t *testing.T) {
	expectedResult := "Test suffix String"
	input := expectedResult
	inst := NewResourceNameFormatter(prefix, suffix)

	result := inst.UndoFormat(input)

	if expectedResult != result {
		t.Fatalf("Suffix should only be trimmed at the end: %s vs %s", expectedResult, result)
	}
}
