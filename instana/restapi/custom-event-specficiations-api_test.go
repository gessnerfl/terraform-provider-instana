package restapi_test

import (
	"strings"
	"testing"

	"github.com/gessnerfl/terraform-provider-instana/testutils"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

const (
	customEventID           = "custom-event-id"
	customEventName         = "custom-event-name"
	customEventEntityType   = "custom-event-entity-type"
	customEventQuery        = "custom-event-query"
	customEventDescription  = "custom-event-description"
	customEventSystemRuleID = "system-rule-id"
)

func TestShouldValidateMinimalCustemEventSpecificationWithSystemRule(t *testing.T) {
	spec := CustomEventSpecification{
		ID:         customEventID,
		Name:       customEventName,
		EntityType: customEventEntityType,
		Rules:      []RuleSpecification{RuleSpecification{DType: SystemRuleType, SystemRuleID: customEventSystemRuleID}},
	}

	if err := spec.Validate(); err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}

	if customEventID != spec.GetID() {
		t.Fatal("Expected GetID returns the correct id of the custom event specification")
	}
}

func TestShouldValidateFullCustomEventSpecificationWithSystemRule(t *testing.T) {
	query := customEventQuery
	description := customEventDescription
	expirationTime := 1234

	spec := CustomEventSpecification{
		ID:             customEventID,
		Name:           customEventName,
		EntityType:     customEventEntityType,
		Query:          &query,
		Description:    &description,
		ExpirationTime: &expirationTime,
		Triggering:     true,
		Enabled:        true,
		Rules:          []RuleSpecification{RuleSpecification{DType: SystemRuleType, SystemRuleID: customEventSystemRuleID}},
		Downstream: &EventSpecificationDownstream{
			IntegrationIds:                []string{"downstream-integration-id"},
			BroadcastToAllAlertingConfigs: true,
		},
	}

	if err := spec.Validate(); err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}
}

func TestFailToValidateCustemEventSpecificationWithSystemRuleWhenIDIsMissing(t *testing.T) {
	spec := CustomEventSpecification{
		Name:       customEventName,
		EntityType: customEventEntityType,
		Rules:      []RuleSpecification{RuleSpecification{DType: SystemRuleType, SystemRuleID: customEventSystemRuleID}},
	}

	if err := spec.Validate(); err == nil || !strings.Contains(err.Error(), "ID") {
		t.Fatal("Expected validate to fail as ID is not provided")
	}
}

func TestFailToValidateCustemEventSpecificationWithSystemRuleWhenNameIsMissing(t *testing.T) {
	spec := CustomEventSpecification{
		ID:         customEventID,
		EntityType: customEventEntityType,
		Rules:      []RuleSpecification{RuleSpecification{DType: SystemRuleType, SystemRuleID: customEventSystemRuleID}},
	}

	if err := spec.Validate(); err == nil || !strings.Contains(err.Error(), "name") {
		t.Fatal("Expected validate to fail as name is not provided")
	}
}

func TestFailToValidateCustemEventSpecificationWithSystemRuleWhenEntityTypeIsMissing(t *testing.T) {
	spec := CustomEventSpecification{
		ID:    customEventID,
		Name:  customEventName,
		Rules: []RuleSpecification{RuleSpecification{DType: SystemRuleType, SystemRuleID: customEventSystemRuleID}},
	}

	if err := spec.Validate(); err == nil || !strings.Contains(err.Error(), "entity type") {
		t.Fatal("Expected validate to fail as entity type is not provided")
	}
}

func TestFailToValidateCustemEventSpecificationWhenNoRuleIsNil(t *testing.T) {
	spec := CustomEventSpecification{
		ID:         customEventID,
		Name:       customEventName,
		EntityType: customEventEntityType,
		Rules:      nil,
	}

	if err := spec.Validate(); err == nil || !strings.Contains(err.Error(), "exactly one rule") {
		t.Fatal("Expected validate to fail as no rule is provided")
	}
}

func TestFailToValidateCustemEventSpecificationWhenNoRuleIsProvided(t *testing.T) {
	spec := CustomEventSpecification{
		ID:         customEventID,
		Name:       customEventName,
		EntityType: customEventEntityType,
		Rules:      []RuleSpecification{},
	}

	if err := spec.Validate(); err == nil || !strings.Contains(err.Error(), "exactly one rule") {
		t.Fatal("Expected validate to fail as no rule is provided")
	}
}

func TestFailToValidateCustemEventSpecificationWhenMultipleRulesAreProvided(t *testing.T) {
	spec := CustomEventSpecification{
		ID:         customEventID,
		Name:       customEventName,
		EntityType: customEventEntityType,
		Rules:      []RuleSpecification{RuleSpecification{DType: SystemRuleType, SystemRuleID: customEventSystemRuleID}, RuleSpecification{DType: SystemRuleType, SystemRuleID: customEventSystemRuleID}},
	}

	if err := spec.Validate(); err == nil || !strings.Contains(err.Error(), "exactly one rule") {
		t.Fatal("Expected validate to fail as no id of the second system rule is not provided")
	}
}

func TestFailToValidateCustemEventSpecificationWhenTheProvidedRuleIsNotValid(t *testing.T) {
	spec := CustomEventSpecification{
		ID:         customEventID,
		Name:       customEventName,
		EntityType: customEventEntityType,
		Rules:      []RuleSpecification{RuleSpecification{DType: SystemRuleType}},
	}

	if err := spec.Validate(); err == nil || !strings.Contains(err.Error(), "id of system rule") {
		t.Fatal("Expected validate to fail as no id of the second system rule is not provided")
	}
}

func TestFailToValidateCustemEventSpecificationWhenDownstreamSpecificationIsNotValid(t *testing.T) {
	spec := CustomEventSpecification{
		ID:         customEventID,
		Name:       customEventName,
		EntityType: customEventEntityType,
		Rules:      []RuleSpecification{RuleSpecification{DType: SystemRuleType, SystemRuleID: customEventSystemRuleID}},
		Downstream: &EventSpecificationDownstream{},
	}

	if err := spec.Validate(); err == nil || !strings.Contains(err.Error(), "integration id") {
		t.Fatal("Expected validate to fail as no integration id is provided for the downstream specification")
	}
}

func TestShouldSuccessfullyValidateSystemRule(t *testing.T) {
	rule := NewSystemRuleSpecification(customEventSystemRuleID, 1000)

	if err := rule.Validate(); err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}
}

func TestShouldFailToValidateSystemRuleWhenSystemRuleIDIsMissing(t *testing.T) {
	rule := RuleSpecification{DType: SystemRuleType}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "id of system rule") {
		t.Fatal("Expected to fail to validate system rule as no system rule id is provided")
	}
}

func TestShouldFailToValidateSystemRuleWhenRuleTypeIsMissing(t *testing.T) {
	rule := RuleSpecification{SystemRuleID: customEventSystemRuleID}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "type of system rule") {
		t.Fatal("Expected to fail to validate system rule as no system rule id is provided")
	}
}

func TestShouldSuccessfullyValidateEventSpecificationDownstream(t *testing.T) {
	rule := EventSpecificationDownstream{
		IntegrationIds:                []string{"integration-id-1"},
		BroadcastToAllAlertingConfigs: true,
	}

	if err := rule.Validate(); err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}
}

func TestShouldFailToValidateEventSpecificationDownstreamWhenNoIntegrationIDsAreNil(t *testing.T) {
	rule := EventSpecificationDownstream{
		IntegrationIds: nil,
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "integration id") {
		t.Fatal("Expected to fail to validate event specification downstream as integration ids is nil")
	}
}

func TestShouldFailToValidateEventSpecificationDownstreamWhenNoIntegrationIDIsProvided(t *testing.T) {
	rule := EventSpecificationDownstream{
		IntegrationIds: []string{},
	}

	if err := rule.Validate(); err == nil || !strings.Contains(err.Error(), "integration id") {
		t.Fatal("Expected to fail to validate event specification downstream as no integration id is provided")
	}
}
