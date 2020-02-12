package utils_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/utils"
	"github.com/stretchr/testify/assert"
)

func TestShouldReturnTrueWhenCheckingForUniqueElementsStringSliceAndSliceIsEmpty(t *testing.T) {
	assert.True(t, StringSliceElementsAreUnique([]string{}))
}

func TestShouldReturnTrueWhenCheckingForUniqueElementsInStringSliceAndSliceContainsUniqueElementsOnly(t *testing.T) {
	assert.True(t, StringSliceElementsAreUnique([]string{"a", "b", "c", "d"}))
}

func TestShouldReturnFalseWhenCheckingForUniqueElementsInStringSliceAndSliceContainsDuplicateElements(t *testing.T) {
	assert.False(t, StringSliceElementsAreUnique([]string{"a", "b", "c", "d", "a"}))
}
