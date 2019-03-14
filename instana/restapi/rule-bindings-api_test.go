package restapi_test

import (
	"strings"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

func TestValidRuleBinding(t *testing.T) {
	ruleBinding := RuleBinding{
		ID:             "test-id",
		Enabled:        false,
		Triggering:     false,
		Severity:       1,
		Text:           "test-text",
		Description:    "test-description",
		ExpirationTime: 60000,
		Query:          "entity.type:jvm",
		RuleIds:        []string{"test-rule-id"},
	}

	if "test-id" != ruleBinding.GetID() {
		t.Errorf("Expected to get correct ID but got %s", ruleBinding.GetID())
	}

	if err := ruleBinding.Validate(); err != nil {
		t.Errorf("Expected valid rule binding got validation error %s", err)
	}
}

func TestInvalidRuleBindingBecauseOfMissingId(t *testing.T) {
	ruleBinding := RuleBinding{
		Enabled:        false,
		Triggering:     false,
		Severity:       1,
		Text:           "test-text",
		Description:    "test-description",
		ExpirationTime: 60000,
		Query:          "entity.type:jvm",
		RuleIds:        []string{"test-rule-id"},
	}

	if err := ruleBinding.Validate(); err == nil || !strings.Contains(err.Error(), "ID") {
		t.Errorf("Expected invalid rule binding because of missing ID")
	}
}

func TestInvalidRuleBindingBecauseOfMissingText(t *testing.T) {
	ruleBinding := RuleBinding{
		ID:             "test-id",
		Enabled:        false,
		Triggering:     false,
		Severity:       1,
		Description:    "test-description",
		ExpirationTime: 60000,
		Query:          "entity.type:jvm",
		RuleIds:        []string{"test-rule-id"},
	}

	if err := ruleBinding.Validate(); err == nil || !strings.Contains(err.Error(), "Text") {
		t.Errorf("Expected invalid rule binding because of missing Text")
	}
}

func TestInvalidRuleBindingBecauseOfMissingDescription(t *testing.T) {
	ruleBinding := RuleBinding{
		ID:             "test-id",
		Enabled:        false,
		Triggering:     false,
		Severity:       1,
		Text:           "test-text",
		ExpirationTime: 60000,
		Query:          "entity.type:jvm",
		RuleIds:        []string{"test-rule-id"},
	}

	if err := ruleBinding.Validate(); err == nil || !strings.Contains(err.Error(), "Description") {
		t.Errorf("Expected invalid rule binding because of missing Description")
	}
}

func TestInvalidRuleBindingBecauseOfMissingRuleIds(t *testing.T) {
	ruleBinding := RuleBinding{
		ID:             "test-id",
		Enabled:        false,
		Triggering:     false,
		Severity:       1,
		Text:           "test-text",
		Description:    "test-description",
		ExpirationTime: 60000,
		Query:          "entity.type:jvm",
	}

	if err := ruleBinding.Validate(); err == nil || !strings.Contains(err.Error(), "RuleIds") {
		t.Errorf("Expected invalid rule binding because of missing RuleIds")
	}
}
