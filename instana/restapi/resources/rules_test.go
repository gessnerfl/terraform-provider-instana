package resources_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi/resources"
	mocks "github.com/gessnerfl/terraform-provider-instana/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewRuleResource(client)
	rule := makeTestRule()

	client.EXPECT().Delete(gomock.Eq("test-rule-id-1"), gomock.Eq("/rules")).Return(nil)

	err := sut.Delete(rule)

	if err != nil {
		t.Errorf("Expected no error got %s", err)
	}
}

func TestFailedDeleteOfRuleByRule(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewRuleResource(client)
	rule := makeTestRule()

	client.EXPECT().Delete(gomock.Eq("test-rule-id-1"), gomock.Eq("/rules")).Return(errors.New("Error during test"))

	err := sut.Delete(rule)

	if err == nil {
		t.Error("Expected to get error")
	}
}

func TestSuccessfulUpsertOfRule(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewRuleResource(client)
	rule := makeTestRule()
	serializedJSON, _ := json.Marshal(rule)

	client.EXPECT().Put(gomock.Eq(rule), gomock.Eq("/rules")).Return(serializedJSON, nil)

	result, err := sut.Upsert(rule)

	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}

	if !cmp.Equal(rule, result) {
		t.Errorf("Expected json to be unmarshalled to %v but got %v; diff %s", rule, result, cmp.Diff(result, result))
	}
}

func TestFailedUpsertOfRuleBecauseOfClientError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewRuleResource(client)
	rule := makeTestRule()

	client.EXPECT().Put(gomock.Eq(rule), gomock.Eq("/rules")).Return(nil, errors.New("Error during test"))

	_, err := sut.Upsert(rule)

	if err == nil {
		t.Error("Expected to get error")
	}
}

func TestFailedUpsertOfRuleBecauseOfInvalidResponseMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewRuleResource(client)
	rule := makeTestRule()

	client.EXPECT().Put(gomock.Eq(rule), gomock.Eq("/rules")).Return([]byte("invalid response"), nil)

	_, err := sut.Upsert(rule)

	if err == nil {
		t.Errorf("Expected to get error")
	}
}

func TestFailedUpsertOfRuleBecauseOfInvalidRuleInResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewRuleResource(client)
	rule := makeTestRule()

	client.EXPECT().Put(gomock.Eq(rule), gomock.Eq("/rules")).Return([]byte("{ \"invalid\" : \"rule\" }"), nil)

	_, err := sut.Upsert(rule)

	if err == nil {
		t.Errorf("Expected to get error")
	}
}

func TestFailedUpsertOfRuleBecauseOfInvalidRuleProvided(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewRuleResource(client)
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

	client.EXPECT().Put(gomock.Eq(rule), gomock.Eq("/rules")).Times(0)

	_, err := sut.Upsert(rule)

	if err == nil {
		t.Error("Expected to get error")
	}
}

func TestSuccessfulGetOneRule(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewRuleResource(client)
	rule := makeTestRule()
	serializedJSON, _ := json.Marshal(rule)

	client.EXPECT().GetOne(gomock.Eq(rule.ID), gomock.Eq("/rules")).Return(serializedJSON, nil)

	data, err := sut.GetOne(rule.ID)

	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}

	if !cmp.Equal(rule, data) {
		t.Errorf("Expected json to be unmarshalled to %v but got %v; diff %s", rule, data, cmp.Diff(rule, data))
	}
}

func TestFailedGetOneRuleBecauseOfErrorFromRestClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewRuleResource(client)
	ruleID := "test-rule-id"

	client.EXPECT().GetOne(gomock.Eq(ruleID), gomock.Eq("/rules")).Return(nil, errors.New("error during test"))

	_, err := sut.GetOne(ruleID)

	if err == nil {
		t.Errorf("Expected to get error")
	}
}

func TestFailedGetOneRuleBecauseOfInvalidJsonArray(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewRuleResource(client)
	ruleID := "test-rule-id"

	client.EXPECT().GetOne(gomock.Eq(ruleID), gomock.Eq("/rules")).Return([]byte("[{ \"invalid\" : \"data\" }]"), nil)

	_, err := sut.GetOne(ruleID)

	if err == nil {
		t.Errorf("Expected to get error")
	}
}

func TestFailedGetOneRuleBecauseOfInvalidJsonObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewRuleResource(client)
	ruleID := "test-rule-id"

	client.EXPECT().GetOne(gomock.Eq(ruleID), gomock.Eq("/rules")).Return([]byte("{ \"invalid\" : \"data\" }"), nil)

	_, err := sut.GetOne(ruleID)

	if err == nil {
		t.Errorf("Expected to get error")
	}
}

func TestFailedGetOneRuleBecauseOfNoJsonAsResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewRuleResource(client)
	ruleID := "test-rule-id"

	client.EXPECT().GetOne(gomock.Eq(ruleID), gomock.Eq("/rules")).Return([]byte("Invalid Data"), nil)

	_, err := sut.GetOne(ruleID)

	if err == nil {
		t.Errorf("Expected to get error")
	}
}

func TestSuccessfulGetAllRules(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewRuleResource(client)
	rule1 := makeTestRuleWithCounter(1)
	rule2 := makeTestRuleWithCounter(2)
	rules := []Rule{rule1, rule2}
	serializedJSON, _ := json.Marshal(rules)

	client.EXPECT().GetAll(gomock.Eq("/rules")).Return(serializedJSON, nil)

	data, err := sut.GetAll()

	if err != nil {
		t.Errorf("Expected no error but got %s", err)
	}

	if !cmp.Equal(rules, data) {
		t.Errorf("Expected json to be unmarshalled to %v but got %v; diff %s", rules, data, cmp.Diff(rules, data))
	}
}

func TestFailedGetAllRulesWithErrorFromRestClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewRuleResource(client)

	client.EXPECT().GetAll(gomock.Eq("/rules")).Return(nil, errors.New("error during test"))

	_, err := sut.GetAll()

	if err == nil {
		t.Errorf("Expected to get error")
	}
}

func TestFailedGetAllRulesWithInvalidJsonArray(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewRuleResource(client)

	client.EXPECT().GetAll(gomock.Eq("/rules")).Return([]byte("[{ \"invalid\" : \"data\" }]"), nil)

	_, err := sut.GetAll()

	if err == nil {
		t.Errorf("Expected to get error")
	}
}

func TestFailedGetAllRulesWithInvalidJsonObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewRuleResource(client)

	client.EXPECT().GetAll(gomock.Eq("/rules")).Return([]byte("{ \"invalid\" : \"data\" }"), nil)

	_, err := sut.GetAll()

	if err == nil {
		t.Errorf("Expected to get error")
	}
}

func TestFailedGetAllRulesWithNoJsonAsResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewRuleResource(client)

	client.EXPECT().GetAll(gomock.Eq("/rules")).Return([]byte("Invalid Data"), nil)

	_, err := sut.GetAll()

	if err == nil {
		t.Errorf("Expected to get error")
	}
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
