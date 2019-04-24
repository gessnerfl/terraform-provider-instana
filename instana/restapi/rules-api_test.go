package restapi_test

import (
	"strings"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

const (
	ruleID                = "rule-id"
	ruleName              = "rule-name"
	ruleEntityType        = "rule-entity-type"
	ruleMetricName        = "rule-metric-name"
	ruleAggregation       = "sum"
	ruleConditionOperator = ">"
)

func TestValidRule(t *testing.T) {
	rule := Rule{
		ID:                ruleID,
		Name:              ruleName,
		EntityType:        ruleEntityType,
		MetricName:        ruleMetricName,
		Rollup:            0,
		Window:            300000,
		Aggregation:       ruleAggregation,
		ConditionOperator: ruleConditionOperator,
		ConditionValue:    0,
	}

	if ruleID != rule.GetID() {
		t.Fatalf("Expected to get correct ID but got %s", rule.GetID())
	}

	if err := rule.Validate(); err != nil {
		t.Fatalf("Expected valid rule got validation error %s", err)
	}
}

func TestShouldFailToValidateRuleWhenIdIsMissing(t *testing.T) {
	rule := Rule{
		Name:              ruleName,
		EntityType:        ruleEntityType,
		MetricName:        ruleMetricName,
		Rollup:            0,
		Window:            300000,
		Aggregation:       ruleAggregation,
		ConditionOperator: ruleConditionOperator,
		ConditionValue:    0,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "ID") {
		t.Fatalf("Expected invalid rule because of missing ID")
	}
}

func TestShouldFailToValidateRuleWhenNameIsMissing(t *testing.T) {
	rule := Rule{
		ID:                ruleID,
		EntityType:        ruleEntityType,
		MetricName:        ruleMetricName,
		Rollup:            0,
		Window:            300000,
		Aggregation:       ruleAggregation,
		ConditionOperator: ruleConditionOperator,
		ConditionValue:    0,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "Name") {
		t.Fatalf("Expected invalid rule because of missing Name")
	}
}

func TestShouldFailToValidateRuleWhenEntityTypeIsMissing(t *testing.T) {
	rule := Rule{
		ID:                ruleID,
		Name:              ruleName,
		MetricName:        ruleMetricName,
		Rollup:            0,
		Window:            300000,
		Aggregation:       ruleAggregation,
		ConditionOperator: ruleConditionOperator,
		ConditionValue:    0,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "EntityType") {
		t.Fatalf("Expected invalid rule because of missing EntityType")
	}
}

func TestShouldFailToValidateRuleWhenMetricNameIsMissing(t *testing.T) {
	rule := Rule{
		ID:                ruleID,
		Name:              ruleName,
		EntityType:        ruleEntityType,
		Rollup:            0,
		Window:            300000,
		Aggregation:       ruleAggregation,
		ConditionOperator: ruleConditionOperator,
		ConditionValue:    0,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "MetricName") {
		t.Fatalf("Expected invalid rule because of missing MetricName")
	}
}

func TestShouldFailToValidateRuleWhenAggregationIsMissing(t *testing.T) {
	rule := Rule{
		ID:                ruleID,
		Name:              ruleName,
		EntityType:        ruleEntityType,
		MetricName:        ruleMetricName,
		Rollup:            0,
		Window:            300000,
		ConditionOperator: ruleConditionOperator,
		ConditionValue:    0,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "Aggregation") {
		t.Fatalf("Expected invalid rule because of missing Aggregation")
	}
}

func TestShouldFailToValidateRuleWhenAggregationIsNotValid(t *testing.T) {
	rule := Rule{
		ID:                ruleID,
		Name:              ruleName,
		EntityType:        ruleEntityType,
		MetricName:        ruleMetricName,
		Rollup:            0,
		Window:            300000,
		Aggregation:       "invalid",
		ConditionOperator: ruleConditionOperator,
		ConditionValue:    0,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "Aggregation") {
		t.Fatalf("Expected invalid rule because of invalid Aggregation")
	}
}

func TestShouldFailToValidateRuleWhenConditionOperatorIsMissing(t *testing.T) {
	rule := Rule{
		ID:             ruleID,
		Name:           ruleName,
		EntityType:     ruleEntityType,
		MetricName:     ruleMetricName,
		Rollup:         0,
		Window:         300000,
		Aggregation:    ruleAggregation,
		ConditionValue: 0,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "ConditionOperator") {
		t.Fatalf("Expected invalid rule because of missing ConditionOperator")
	}
}

func TestShouldFailToValidateRuleWhenConditionOperatorIsNotValid(t *testing.T) {
	rule := Rule{
		ID:                ruleID,
		Name:              ruleName,
		EntityType:        ruleEntityType,
		MetricName:        ruleMetricName,
		Rollup:            0,
		Window:            300000,
		Aggregation:       ruleAggregation,
		ConditionOperator: "invalid",
		ConditionValue:    0,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "ConditionOperator") {
		t.Fatalf("Expected invalid rule because of invalid ConditionOperator")
	}
}
