package instana_test

import (
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/google/go-cmp/cmp"
)

func TestRandomID(t *testing.T) {
	id := RandomID()

	if len(id) == 0 {
		t.Fatal("Expected to get a new id generated")
	}
}

func TestReadStringArrayParameterFromResourceWhenParameterIsProvided(t *testing.T) {
	ruleIds := []string{"test1", "test2"}
	data := make(map[string]interface{})
	data[RuleBindingFieldRuleIds] = ruleIds
	resourceData := createRuleBindingResourceData(t, data)
	result := ReadStringArrayParameterFromResource(resourceData, RuleBindingFieldRuleIds)

	if result == nil {
		t.Fatal("Expected result to available")
	}
	if !cmp.Equal(result, ruleIds) {
		t.Fatal("Expected to get rule ids in string array")
	}
}

func TestReadStringArrayParameterFromResourceWhenParameterIsMissing(t *testing.T) {
	resourceData := createEmptyRuleBindingResourceData(t)
	result := ReadStringArrayParameterFromResource(resourceData, RuleBindingFieldRuleIds)

	if result != nil {
		t.Fatal("Expected result to be nil as no data is provided")
	}
}
