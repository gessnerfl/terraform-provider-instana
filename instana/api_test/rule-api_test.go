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

func TestValidRule(t *testing.T) {
	rule := makeTestRule()

	if "test-rule-id-1" != rule.GetID() {
		t.Errorf("Expected to get correct ID but got %s", rule.GetID())
	}

	if err := rule.Validate(); err != nil {
		t.Errorf("Expected valid rule got validation error %s", err)
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
		t.Errorf("Expected invalid rule because of missing ID")
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
		t.Errorf("Expected invalid rule because of missing Name")
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
		t.Errorf("Expected invalid rule because of missing EntityType")
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
		t.Errorf("Expected invalid rule because of missing MetricName")
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
		t.Errorf("Expected invalid rule because of missing Aggregation")
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
		t.Errorf("Expected invalid rule because of missing ConditionOperator")
	}
}

func TestSuccessfulDeleteOfRuleByRule(t *testing.T) {
	client := NewMockRestClient(pegomock.WithT(t))
	sut := NewRuleAPI(client)
	rule := makeTestRule()

	err := sut.Delete(rule)

	if err != nil {
		t.Errorf("Expected no error got %s", err)
	}

	client.VerifyWasCalledOnce().Delete("test-rule-id-1", "/rules")
}
func TestFailedDeleteOfRuleByRule(t *testing.T) {
	client := NewMockRestClient(pegomock.WithT(t))
	sut := NewRuleAPI(client)
	rule := makeTestRule()

	pegomock.When(client.Delete("test-rule-id-1", "/rules")).ThenReturn(errors.New("Error during test"))

	err := sut.Delete(rule)

	if err == nil {
		t.Error("Expected to get error")
	}

	client.VerifyWasCalledOnce().Delete("test-rule-id-1", "/rules")
}

func TestSuccessfulUpsertOfRule(t *testing.T) {
	client := NewMockRestClient(pegomock.WithT(t))
	sut := NewRuleAPI(client)
	rule := makeTestRule()
	serializedJson, _ := json.Marshal(rule)

	pegomock.When(client.Put(rule, "/rules")).ThenReturn(serializedJson, nil)

	result, err := sut.Upsert(rule)

	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}

	if !cmp.Equal(rule, result) {
		t.Errorf("Expected json to be unmarshalled to %v but got %v; diff %s", rule, result, cmp.Diff(result, result))
	}

	client.VerifyWasCalledOnce().Put(rule, "/rules")
}

func TestFailedUpsertOfRuleBecauseOfClientError(t *testing.T) {
	client := NewMockRestClient(pegomock.WithT(t))
	sut := NewRuleAPI(client)
	rule := makeTestRule()

	pegomock.When(client.Put(rule, "/rules")).ThenReturn(nil, errors.New("Error during test"))

	_, err := sut.Upsert(rule)

	if err == nil {
		t.Error("Expected to get error")
	}

	client.VerifyWasCalledOnce().Put(rule, "/rules")
}

func TestFailedUpsertOfRuleBecauseOfInvalidResponseMessage(t *testing.T) {
	client := NewMockRestClient(pegomock.WithT(t))
	sut := NewRuleAPI(client)
	rule := makeTestRule()

	pegomock.When(client.Put(rule, "/rules")).ThenReturn([]byte("invalid response"), nil)

	_, err := sut.Upsert(rule)

	if err == nil {
		t.Errorf("Expected to get error")
	}

	client.VerifyWasCalledOnce().Put(rule, "/rules")
}

func TestFailedUpsertOfRuleBecauseOfInvalidRuleInResponse(t *testing.T) {
	client := NewMockRestClient(pegomock.WithT(t))
	sut := NewRuleAPI(client)
	rule := makeTestRule()

	pegomock.When(client.Put(rule, "/rules")).ThenReturn([]byte("{ \"invalid\" : \"rule\" }"), nil)

	_, err := sut.Upsert(rule)

	if err == nil {
		t.Errorf("Expected to get error")
	}

	client.VerifyWasCalledOnce().Put(rule, "/rules")
}

func TestFailedUpsertOfRuleBecauseOfInvalidRuleProvided(t *testing.T) {
	client := NewMockRestClient(pegomock.WithT(t))
	sut := NewRuleAPI(client)
	rule := Rule{
		Name:              "Test Rule",
		EntityType:        "test",
		MetricName:        "test.metric",
		Rollup:            0,
		Window:            300000,
		Aggregation:       "sum",
		ConditionValue:    0,
		ConditionOperator: ">",
	}

	_, err := sut.Upsert(rule)

	if err == nil {
		t.Error("Expected to get error")
	}

	client.VerifyWasCalled(pegomock.Never()).Put(rule, "/rules")
}

func TestSuccessfulGetOneRule(t *testing.T) {
	client := NewMockRestClient(pegomock.WithT(t))
	sut := NewRuleAPI(client)
	rule := makeTestRule()
	serializedJson, _ := json.Marshal(rule)

	pegomock.When(client.GetOne(rule.ID, "/rules")).ThenReturn(serializedJson, nil)

	data, err := sut.GetOne(rule.ID)

	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}

	if !cmp.Equal(rule, data) {
		t.Errorf("Expected json to be unmarshalled to %v but got %v; diff %s", rule, data, cmp.Diff(rule, data))
	}

	client.VerifyWasCalledOnce().GetOne(rule.ID, "/rules")
}

func TestFailedGetOneRuleBecauseOfErrorFromRestClient(t *testing.T) {
	client := NewMockRestClient(pegomock.WithT(t))
	sut := NewRuleAPI(client)
	ruleId := "test-rule-id"

	pegomock.When(client.GetOne(ruleId, "/rules")).ThenReturn(nil, errors.New("error during test"))

	_, err := sut.GetOne(ruleId)

	if err == nil {
		t.Errorf("Expected to get error")
	}

	client.VerifyWasCalledOnce().GetOne(ruleId, "/rules")
}

func TestFailedGetOneRuleBecauseOfInvalidJsonArray(t *testing.T) {
	client := NewMockRestClient(pegomock.WithT(t))
	sut := NewRuleAPI(client)
	ruleId := "test-rule-id"

	pegomock.When(client.GetOne(ruleId, "/rules")).ThenReturn([]byte("[{ \"invalid\" : \"data\" }]"), nil)

	_, err := sut.GetOne(ruleId)

	if err == nil {
		t.Errorf("Expected to get error")
	}

	client.VerifyWasCalledOnce().GetOne(ruleId, "/rules")
}

func TestFailedGetOneRuleBecauseOfInvalidJsonObject(t *testing.T) {
	client := NewMockRestClient(pegomock.WithT(t))
	sut := NewRuleAPI(client)
	ruleId := "test-rule-id"

	pegomock.When(client.GetOne(ruleId, "/rules")).ThenReturn([]byte("{ \"invalid\" : \"data\" }"), nil)

	_, err := sut.GetOne(ruleId)

	if err == nil {
		t.Errorf("Expected to get error")
	}

	client.VerifyWasCalledOnce().GetOne(ruleId, "/rules")
}

func TestFailedGetOneRuleBecauseOfNoJsonAsResponse(t *testing.T) {
	client := NewMockRestClient(pegomock.WithT(t))
	sut := NewRuleAPI(client)
	ruleId := "test-rule-id"

	pegomock.When(client.GetOne(ruleId, "/rules")).ThenReturn([]byte("Invalid Data"), nil)

	_, err := sut.GetOne(ruleId)

	if err == nil {
		t.Errorf("Expected to get error")
	}

	client.VerifyWasCalledOnce().GetOne(ruleId, "/rules")
}

func TestSuccessfulGetAllRules(t *testing.T) {
	client := NewMockRestClient(pegomock.WithT(t))
	sut := NewRuleAPI(client)
	rule1 := makeTestRuleWithCounter(1)
	rule2 := makeTestRuleWithCounter(2)
	rules := []Rule{rule1, rule2}
	serializedJson, _ := json.Marshal(rules)

	pegomock.When(client.GetAll("/rules")).ThenReturn(serializedJson, nil)

	data, err := sut.GetAll()

	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}

	if !cmp.Equal(rules, data) {
		t.Errorf("Expected json to be unmarshalled to %v but got %v; diff %s", rules, data, cmp.Diff(rules, data))
	}

	client.VerifyWasCalledOnce().GetAll("/rules")
}

func TestFailedGetAllRulesWithErrorFromRestClient(t *testing.T) {
	client := NewMockRestClient(pegomock.WithT(t))
	sut := NewRuleAPI(client)

	pegomock.When(client.GetAll("/rules")).ThenReturn(nil, errors.New("error during test"))

	_, err := sut.GetAll()

	if err == nil {
		t.Errorf("Expected to get error")
	}

	client.VerifyWasCalledOnce().GetAll("/rules")
}

func TestFailedGetAllRulesWithInvalidJsonArray(t *testing.T) {
	client := NewMockRestClient(pegomock.WithT(t))
	sut := NewRuleAPI(client)

	pegomock.When(client.GetAll("/rules")).ThenReturn([]byte("[{ \"invalid\" : \"data\" }]"), nil)

	_, err := sut.GetAll()

	if err == nil {
		t.Errorf("Expected to get error")
	}

	client.VerifyWasCalledOnce().GetAll("/rules")
}

func TestFailedGetAllRulesWithInvalidJsonObject(t *testing.T) {
	client := NewMockRestClient(pegomock.WithT(t))
	sut := NewRuleAPI(client)

	pegomock.When(client.GetAll("/rules")).ThenReturn([]byte("{ \"invalid\" : \"data\" }"), nil)

	_, err := sut.GetAll()

	if err == nil {
		t.Errorf("Expected to get error")
	}

	client.VerifyWasCalledOnce().GetAll("/rules")
}

func TestFailedGetAllRulesWithNoJsonAsResponse(t *testing.T) {
	client := NewMockRestClient(pegomock.WithT(t))
	sut := NewRuleAPI(client)

	pegomock.When(client.GetAll("/rules")).ThenReturn([]byte("Invalid Data"), nil)

	_, err := sut.GetAll()

	if err == nil {
		t.Errorf("Expected to get error")
	}

	client.VerifyWasCalledOnce().GetAll("/rules")
}

func makeTestRule() Rule {
	return makeTestRuleWithCounter(1)
}

func makeTestRuleWithCounter(counter int) Rule {
	id := fmt.Sprintf("test-rule-id-%d", counter)
	name := fmt.Sprintf("Test Rule %d", counter)
	return Rule{
		ID:                id,
		Name:              name,
		EntityType:        "test",
		MetricName:        "test.metric",
		Rollup:            0,
		Window:            300000,
		Aggregation:       "sum",
		ConditionOperator: ">",
		ConditionValue:    0,
	}
}
