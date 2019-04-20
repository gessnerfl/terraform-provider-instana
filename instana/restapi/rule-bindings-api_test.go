package restapi_test

import (
	"strings"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

const (
	ruleBindingId          = "rule-binding-id"
	ruleBindingText        = "rule-binding-text"
	ruleBindingDescription = "rule-binding-description"
	ruleBindingQuery       = "entity.type:jvm"
	ruleBindingRuleId      = "rule-id"
)

func TestValidRuleBinding(t *testing.T) {
	ruleBinding := RuleBinding{
		ID:             ruleBindingId,
		Enabled:        false,
		Triggering:     false,
		Severity:       1,
		Text:           ruleBindingText,
		Description:    ruleBindingDescription,
		ExpirationTime: 60000,
		Query:          ruleBindingQuery,
		RuleIds:        []string{ruleBindingRuleId},
	}

	if ruleBindingId != ruleBinding.GetID() {
		t.Fatalf("Expected to get correct ID but got %s", ruleBinding.GetID())
	}

	if err := ruleBinding.Validate(); err != nil {
		t.Fatalf("Expected valid rule binding got validation error %s", err)
	}
}

func TestInvalidRuleBindingBecauseOfMissingId(t *testing.T) {
	ruleBinding := RuleBinding{
		Enabled:        false,
		Triggering:     false,
		Severity:       1,
		Text:           ruleBindingText,
		Description:    ruleBindingDescription,
		ExpirationTime: 60000,
		Query:          ruleBindingQuery,
		RuleIds:        []string{ruleBindingRuleId},
	}

	if err := ruleBinding.Validate(); err == nil || !strings.Contains(err.Error(), "ID") {
		t.Fatalf("Expected invalid rule binding because of missing ID")
	}
}

func TestInvalidRuleBindingBecauseOfMissingText(t *testing.T) {
	ruleBinding := RuleBinding{
		ID:             ruleBindingId,
		Enabled:        false,
		Triggering:     false,
		Severity:       1,
		Description:    ruleBindingDescription,
		ExpirationTime: 60000,
		Query:          ruleBindingQuery,
		RuleIds:        []string{ruleBindingRuleId},
	}

	if err := ruleBinding.Validate(); err == nil || !strings.Contains(err.Error(), "Text") {
		t.Fatalf("Expected invalid rule binding because of missing Text")
	}
}

func TestInvalidRuleBindingBecauseOfMissingDescription(t *testing.T) {
	ruleBinding := RuleBinding{
		ID:             ruleBindingId,
		Enabled:        false,
		Triggering:     false,
		Severity:       1,
		Text:           ruleBindingText,
		ExpirationTime: 60000,
		Query:          ruleBindingQuery,
		RuleIds:        []string{ruleBindingRuleId},
	}

	if err := ruleBinding.Validate(); err == nil || !strings.Contains(err.Error(), "Description") {
		t.Fatalf("Expected invalid rule binding because of missing Description")
	}
}

func TestInvalidRuleBindingBecauseOfMissingRuleIds(t *testing.T) {
	ruleBinding := RuleBinding{
		ID:             ruleBindingId,
		Enabled:        false,
		Triggering:     false,
		Severity:       1,
		Text:           ruleBindingText,
		Description:    ruleBindingDescription,
		ExpirationTime: 60000,
		Query:          ruleBindingQuery,
	}

	if err := ruleBinding.Validate(); err == nil || !strings.Contains(err.Error(), "RuleIds") {
		t.Fatalf("Expected invalid rule binding because of missing RuleIds")
	}
}
