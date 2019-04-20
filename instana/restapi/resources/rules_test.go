package resources_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi/resources"
	mocks "github.com/gessnerfl/terraform-provider-instana/mocks"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
)

const (
	ruleID                = "test-rule-id"
	ruleName              = "Test Rule"
	ruleEntityType        = "entity_type"
	ruleMetricName        = "entity.metric"
	ruleAggregation       = "sum"
	ruleConditionOperator = ">"
)

func TestSuccessfulGetOneRule(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewRuleResource(client)
	rule := makeTestRule()
	serializedJSON, _ := json.Marshal(rule)

	client.EXPECT().GetOne(gomock.Eq(rule.ID), gomock.Eq(restapi.RulesResourcePath)).Return(serializedJSON, nil)

	data, err := sut.GetOne(rule.ID)

	if err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}

	if !cmp.Equal(rule, data) {
		t.Fatalf(testutils.ExpectedUnmarshalledJSONWithStruct, rule, data, cmp.Diff(rule, data))
	}
}

func TestFailedGetOneRuleBecauseOfErrorFromRestClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewRuleResource(client)

	client.EXPECT().GetOne(gomock.Eq(ruleID), gomock.Eq(restapi.RulesResourcePath)).Return(nil, errors.New("error during test"))

	_, err := sut.GetOne(ruleID)

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
}

func TestFailedGetOneRuleBecauseOfInvalidJsonArray(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewRuleResource(client)

	client.EXPECT().GetOne(gomock.Eq(ruleID), gomock.Eq(restapi.RulesResourcePath)).Return([]byte("[{ \"invalid\" : \"data\" }]"), nil)

	_, err := sut.GetOne(ruleID)

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
}

func TestFailedGetOneRuleBecauseOfInvalidJsonObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewRuleResource(client)

	client.EXPECT().GetOne(gomock.Eq(ruleID), gomock.Eq(restapi.RulesResourcePath)).Return([]byte("{ \"invalid\" : \"data\" }"), nil)

	_, err := sut.GetOne(ruleID)

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
}

func TestFailedGetOneRuleBecauseOfNoJsonAsResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewRuleResource(client)

	client.EXPECT().GetOne(gomock.Eq(ruleID), gomock.Eq(restapi.RulesResourcePath)).Return([]byte("Invalid Data"), nil)

	_, err := sut.GetOne(ruleID)

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
}

func TestSuccessfulUpsertOfRule(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewRuleResource(client)
	rule := makeTestRule()
	serializedJSON, _ := json.Marshal(rule)

	client.EXPECT().Put(gomock.Eq(rule), gomock.Eq(restapi.RulesResourcePath)).Return(serializedJSON, nil)

	result, err := sut.Upsert(rule)

	if err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}

	if !cmp.Equal(rule, result) {
		t.Fatalf(testutils.ExpectedUnmarshalledJSONWithStruct, rule, result, cmp.Diff(result, result))
	}
}

func TestFailedUpsertOfRuleBecauseOfClientError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewRuleResource(client)
	rule := makeTestRule()

	client.EXPECT().Put(gomock.Eq(rule), gomock.Eq(restapi.RulesResourcePath)).Return(nil, errors.New("Error during test"))

	_, err := sut.Upsert(rule)

	if err == nil {
		t.Fatal(testutils.ExpectedError)
	}
}

func TestFailedUpsertOfRuleBecauseOfInvalidResponseMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewRuleResource(client)
	rule := makeTestRule()

	client.EXPECT().Put(gomock.Eq(rule), gomock.Eq(restapi.RulesResourcePath)).Return([]byte("invalid response"), nil)

	_, err := sut.Upsert(rule)

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
}

func TestFailedUpsertOfRuleBecauseOfInvalidRuleInResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewRuleResource(client)
	rule := makeTestRule()

	client.EXPECT().Put(gomock.Eq(rule), gomock.Eq(restapi.RulesResourcePath)).Return([]byte("{ \"invalid\" : \"rule\" }"), nil)

	_, err := sut.Upsert(rule)

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
}

func TestFailedUpsertOfRuleBecauseOfInvalidRuleProvided(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewRuleResource(client)
	rule := restapi.Rule{
		Name:              ruleName,
		EntityType:        ruleEntityType,
		MetricName:        ruleMetricName,
		Rollup:            0,
		Window:            300000,
		Aggregation:       ruleAggregation,
		ConditionValue:    0,
		ConditionOperator: ruleConditionOperator,
	}

	client.EXPECT().Put(gomock.Eq(rule), gomock.Eq(restapi.RulesResourcePath)).Times(0)

	_, err := sut.Upsert(rule)

	if err == nil {
		t.Fatal(testutils.ExpectedError)
	}
}

func TestSuccessfulDeleteOfRuleByRule(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewRuleResource(client)
	rule := makeTestRule()

	client.EXPECT().Delete(gomock.Eq(ruleID), gomock.Eq(restapi.RulesResourcePath)).Return(nil)

	err := sut.Delete(rule)

	if err != nil {
		t.Fatalf("Expected no error got %s", err)
	}
}

func TestFailedDeleteOfRuleByRule(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewRuleResource(client)
	rule := makeTestRule()

	client.EXPECT().Delete(gomock.Eq(ruleID), gomock.Eq(restapi.RulesResourcePath)).Return(errors.New("Error during test"))

	err := sut.Delete(rule)

	if err == nil {
		t.Fatal(testutils.ExpectedError)
	}
}

func makeTestRule() restapi.Rule {
	return restapi.Rule{
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
}
