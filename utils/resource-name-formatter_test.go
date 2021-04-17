package utils_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/stretchr/testify/assert"
)

const input = "Test String"
const prefix = "prefix"
const suffix = "suffix"

func TestShouldAppendPrefixAndSuffixToName(t *testing.T) {
	expectedResult := prefix + " " + input + " " + suffix
	inst := NewResourceNameFormatter(prefix, suffix)

	result := inst.Format(input)

	assert.Equal(t, expectedResult, result)
}

func TestShouldNotAppendPrefixWhenNoPrefixIsDefined(t *testing.T) {
	expectedResult := input + " " + suffix
	inst := NewResourceNameFormatter("", suffix)

	result := inst.Format(input)

	assert.Equal(t, expectedResult, result)
}

func TestShouldNotAppendSuffixWhenNoSuffixIsDefined(t *testing.T) {
	expectedResult := prefix + " " + input
	inst := NewResourceNameFormatter(prefix, "")

	result := inst.Format(input)

	assert.Equal(t, expectedResult, result)
}

func TestShouldRemoveSuffixFromNameWhenTerraformManagedResourceNameFormatterIsUsed(t *testing.T) {
	expectedResult := "Test String"
	input := prefix + " " + input + " " + suffix
	inst := NewResourceNameFormatter(prefix, suffix)

	result := inst.UndoFormat(input)

	assert.Equal(t, expectedResult, result)
}

func TestShouldNotRemovePrefixStringWhenPrefixStringAppearsInTheMiddle(t *testing.T) {
	expectedResult := "Test prefix String"
	input := expectedResult
	inst := NewResourceNameFormatter(prefix, suffix)

	result := inst.UndoFormat(input)

	assert.Equal(t, expectedResult, result)
}

func TestShouldNotRemoveSuffixStringWhenSuffixStringAppearsInTheMiddle(t *testing.T) {
	expectedResult := "Test suffix String"
	input := expectedResult
	inst := NewResourceNameFormatter(prefix, suffix)

	result := inst.UndoFormat(input)

	assert.Equal(t, expectedResult, result)
}
