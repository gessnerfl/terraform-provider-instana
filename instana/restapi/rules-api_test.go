package restapi_test

import (
	"strings"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

func TestValidRule(t *testing.T) {
	rule := Rule{
		ID:                "test-rule-id",
		Name:              "Test Rule",
		EntityType:        "test",
		MetricName:        "test.metric",
		Rollup:            0,
		Window:            300000,
		Aggregation:       "sum",
		ConditionOperator: ">",
		ConditionValue:    0,
	}

	if "test-rule-id" != rule.GetID() {
		t.Fatalf("Expected to get correct ID but got %s", rule.GetID())
	}

	if err := rule.Validate(); err != nil {
		t.Fatalf("Expected valid rule got validation error %s", err)
	}
}

func TestInvalidRuleBecauseOfMissingId(t *testing.T) {
	rule := Rule{
		Name:              "Test Rule",
		EntityType:        "test",
		MetricName:        "test.metric",
		Rollup:            0,
		Window:            300000,
		Aggregation:       "sum",
		ConditionOperator: ">",
		ConditionValue:    0,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "ID") {
		t.Fatalf("Expected invalid rule because of missing ID")
	}
}

func TestInvalidRuleBecauseOfMissingName(t *testing.T) {
	rule := Rule{
		ID:                "test-rule-id",
		EntityType:        "test",
		MetricName:        "test.metric",
		Rollup:            0,
		Window:            300000,
		Aggregation:       "sum",
		ConditionOperator: ">",
		ConditionValue:    0,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "Name") {
		t.Fatalf("Expected invalid rule because of missing Name")
	}
}

func TestInvalidRuleBecauseOfMissingEntityType(t *testing.T) {
	rule := Rule{
		ID:                "test-rule-id",
		Name:              "Test Rule",
		MetricName:        "test.metric",
		Rollup:            0,
		Window:            300000,
		Aggregation:       "sum",
		ConditionOperator: ">",
		ConditionValue:    0,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "EntityType") {
		t.Fatalf("Expected invalid rule because of missing EntityType")
	}
}

func TestInvalidRuleBecauseOfMissingMetricName(t *testing.T) {
	rule := Rule{
		ID:                "test-rule-id",
		Name:              "Test Rule",
		EntityType:        "test",
		Rollup:            0,
		Window:            300000,
		Aggregation:       "sum",
		ConditionOperator: ">",
		ConditionValue:    0,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "MetricName") {
		t.Fatalf("Expected invalid rule because of missing MetricName")
	}
}

func TestInvalidRuleBecauseOfMissingAggregation(t *testing.T) {
	rule := Rule{
		ID:                "test-rule-id",
		Name:              "Test Rule",
		EntityType:        "test",
		MetricName:        "test.metric",
		Rollup:            0,
		Window:            300000,
		ConditionOperator: ">",
		ConditionValue:    0,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "Aggregation") {
		t.Fatalf("Expected invalid rule because of missing Aggregation")
	}
}

func TestInvalidRuleBecauseOfMissingConditionOperator(t *testing.T) {
	rule := Rule{
		ID:             "test-rule-id",
		Name:           "Test Rule",
		EntityType:     "test",
		MetricName:     "test.metric",
		Rollup:         0,
		Window:         300000,
		Aggregation:    "sum",
		ConditionValue: 0,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "ConditionOperator") {
		t.Fatalf("Expected invalid rule because of missing ConditionOperator")
	}
}
