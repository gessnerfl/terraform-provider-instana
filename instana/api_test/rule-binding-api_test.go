package api_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/api"
	"github.com/google/go-cmp/cmp"
	"github.com/petergtz/pegomock"
)

func TestValidRuleBinding(t *testing.T) {
	ruleBinding := makeTestRuleBinding()

	if "test-rule-binding-id-1" != ruleBinding.GetID() {
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

func TestSuccessfulDeleteOfRuleBindingByRuleBinding(t *testing.T) {
	client := NewMockRestClient(pegomock.WithT(t))
	sut := NewRuleBindingAPI(client)
	ruleBinding := makeTestRuleBinding()

	err := sut.Delete(ruleBinding)

	if err != nil {
		t.Errorf("Expected no error got %s", err)
	}

	client.VerifyWasCalledOnce().Delete("test-rule-binding-id-1", "/ruleBindings")
}
func TestFailedDeleteOfRuleBindingByRuleBinding(t *testing.T) {
	client := NewMockRestClient(pegomock.WithT(t))
	sut := NewRuleBindingAPI(client)
	ruleBinding := makeTestRuleBinding()

	pegomock.When(client.Delete("test-rule-binding-id-1", "/ruleBindings")).ThenReturn(errors.New("Error during test"))

	err := sut.Delete(ruleBinding)

	if err == nil {
		t.Error("Expected to get error")
	}

	client.VerifyWasCalledOnce().Delete("test-rule-binding-id-1", "/ruleBindings")
}

func TestSuccessfulUpsertOfRuleBinding(t *testing.T) {
	client := NewMockRestClient(pegomock.WithT(t))
	sut := NewRuleBindingAPI(client)
	ruleBinding := makeTestRuleBinding()

	err := sut.Upsert(ruleBinding)

	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}

	client.VerifyWasCalledOnce().Put(ruleBinding, "/ruleBindings")
}

func TestFailedUpsertOfRuleBindingBecauseOfInvalidRuleBinding(t *testing.T) {
	client := NewMockRestClient(pegomock.WithT(t))
	sut := NewRuleBindingAPI(client)
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

	err := sut.Upsert(ruleBinding)

	if err == nil {
		t.Error("Expected to get error")
	}

	client.VerifyWasCalled(pegomock.Never()).Put(ruleBinding, "/ruleBindings")
}

func TestFailedUpsertOfRuleBindingBecauseOfClientError(t *testing.T) {
	client := NewMockRestClient(pegomock.WithT(t))
	sut := NewRuleBindingAPI(client)
	ruleBinding := makeTestRuleBinding()

	pegomock.When(client.Put(ruleBinding, "/ruleBindings")).ThenReturn(errors.New("Error during test"))

	err := sut.Upsert(ruleBinding)

	if err == nil {
		t.Error("Expected to get error")
	}

	client.VerifyWasCalledOnce().Put(ruleBinding, "/ruleBindings")
}

func TestSuccessfulGetAllRuleBindings(t *testing.T) {
	client := NewMockRestClient(pegomock.WithT(t))
	sut := NewRuleBindingAPI(client)
	ruleBinding1 := makeTestRuleBindingWithCounter(1)
	ruleBinding2 := makeTestRuleBindingWithCounter(2)
	ruleBindings := []RuleBinding{ruleBinding1, ruleBinding2}
	serializedJson, _ := json.Marshal(ruleBindings)

	pegomock.When(client.GetAll("/ruleBindings")).ThenReturn(serializedJson, nil)

	data, err := sut.GetAll()

	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}

	if !cmp.Equal(ruleBindings, data) {
		t.Errorf("Expected json to be unmarshalled to %v but got %v; diff %s", ruleBindings, data, cmp.Diff(ruleBindings, data))
	}

	client.VerifyWasCalledOnce().GetAll("/ruleBindings")
}

func TestFailedGetAllRuleBindingsWithErrorFromRestClient(t *testing.T) {
	client := NewMockRestClient(pegomock.WithT(t))
	sut := NewRuleBindingAPI(client)

	pegomock.When(client.GetAll("/ruleBindings")).ThenReturn(nil, errors.New("error during test"))

	_, err := sut.GetAll()

	if err == nil {
		t.Errorf("Expected to get error")
	}

	client.VerifyWasCalledOnce().GetAll("/ruleBindings")
}

func TestFailedGetAllRuleBindingsWithInvalidJsonArray(t *testing.T) {
	client := NewMockRestClient(pegomock.WithT(t))
	sut := NewRuleBindingAPI(client)

	pegomock.When(client.GetAll("/ruleBindings")).ThenReturn([]byte("[{ \"invalid\" : \"data\" }]"), nil)

	_, err := sut.GetAll()

	if err == nil {
		t.Errorf("Expected to get error")
	}

	client.VerifyWasCalledOnce().GetAll("/ruleBindings")
}

func TestFailedGetAllRuleBindingWithInvalidJsonObject(t *testing.T) {
	client := NewMockRestClient(pegomock.WithT(t))
	sut := NewRuleBindingAPI(client)

	pegomock.When(client.GetAll("/ruleBindings")).ThenReturn([]byte("{ \"invalid\" : \"data\" }"), nil)

	_, err := sut.GetAll()

	if err == nil {
		t.Errorf("Expected to get error")
	}

	client.VerifyWasCalledOnce().GetAll("/ruleBindings")
}

func TestFailedGetAllRuleBindingsWithNoJsonAsResponse(t *testing.T) {
	client := NewMockRestClient(pegomock.WithT(t))
	sut := NewRuleBindingAPI(client)

	pegomock.When(client.GetAll("/ruleBindings")).ThenReturn([]byte("Invalid Data"), nil)

	_, err := sut.GetAll()

	if err == nil {
		t.Errorf("Expected to get error")
	}

	client.VerifyWasCalledOnce().GetAll("/ruleBindings")
}

func makeTestRuleBinding() RuleBinding {
	return makeTestRuleBindingWithCounter(1)
}

func makeTestRuleBindingWithCounter(counter int) RuleBinding {
	id := fmt.Sprintf("test-rule-binding-id-%d", counter)
	text := fmt.Sprintf("Test Rule Binding Text %d", counter)
	description := fmt.Sprintf("Test Rule Binding Description %d", counter)
	return RuleBinding{
		ID:             id,
		Enabled:        false,
		Triggering:     false,
		Severity:       1,
		Text:           text,
		Description:    description,
		ExpirationTime: 60000,
		Query:          "entity.type:jvm",
		RuleIds:        []string{"test-rule-id"},
	}
}
