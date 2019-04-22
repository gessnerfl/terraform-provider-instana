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
	ruleBindingID          = "test-rule-binding-id"
	ruleBindingText        = "test-text"
	ruleBindingDescription = "test-rule-binding-description"
	ruleBindingQuery       = "entity.type:jvm"
	ruldBindingRuleID      = "rule-id-1"
)

func TestSuccessfulGetOneRuleBinding(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewRuleBindingResource(client)
	ruleBinding := makeTestRuleBinding()
	serializedJSON, _ := json.Marshal(ruleBinding)

	client.EXPECT().GetOne(gomock.Eq(ruleBinding.ID), gomock.Eq(restapi.RuleBindingsResourcePath)).Return(serializedJSON, nil)

	data, err := sut.GetOne(ruleBinding.ID)

	if err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}

	if !cmp.Equal(ruleBinding, data) {
		t.Fatalf(testutils.ExpectedUnmarshalledJSONWithStruct, ruleBinding, data, cmp.Diff(ruleBinding, data))
	}
}

func TestFailedGetOneRuleBindingBecauseOfErrorFromRestClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewRuleBindingResource(client)

	client.EXPECT().GetOne(gomock.Eq(ruleBindingID), gomock.Eq(restapi.RuleBindingsResourcePath)).Return(nil, errors.New("error during test"))

	_, err := sut.GetOne(ruleBindingID)

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
}

func TestFailedGetOneRuleBindingBecauseOfInvalidJsonArray(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewRuleBindingResource(client)

	client.EXPECT().GetOne(gomock.Eq(ruleBindingID), gomock.Eq(restapi.RuleBindingsResourcePath)).Return([]byte("[{ \"invalid\" : \"data\" }]"), nil)

	_, err := sut.GetOne(ruleBindingID)

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
}

func TestFailedGetOneRuleBindingBecauseOfInvalidJsonObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewRuleBindingResource(client)

	client.EXPECT().GetOne(gomock.Eq(ruleBindingID), gomock.Eq(restapi.RuleBindingsResourcePath)).Return([]byte("{ \"invalid\" : \"data\" }"), nil)

	_, err := sut.GetOne(ruleBindingID)

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
}

func TestFailedGetOneRuleBindingBecauseOfNoJsonAsResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewRuleBindingResource(client)

	client.EXPECT().GetOne(gomock.Eq(ruleBindingID), gomock.Eq(restapi.RuleBindingsResourcePath)).Return([]byte("Invalid Data"), nil)

	_, err := sut.GetOne(ruleBindingID)

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
}

func TestSuccessfulUpsertOfRuleBinding(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewRuleBindingResource(client)
	ruleBinding := makeTestRuleBinding()
	serializedJSON, _ := json.Marshal(ruleBinding)

	client.EXPECT().Put(gomock.Eq(ruleBinding), gomock.Eq(restapi.RuleBindingsResourcePath)).Return(serializedJSON, nil)

	result, err := sut.Upsert(ruleBinding)

	if err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}

	if !cmp.Equal(ruleBinding, result) {
		t.Fatalf(testutils.ExpectedUnmarshalledJSONWithStruct, ruleBinding, result, cmp.Diff(ruleBinding, result))
	}
}

func TestFailedUpsertOfRuleBindingBecauseOfInvalidRuleBinding(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewRuleBindingResource(client)
	ruleBinding := restapi.RuleBinding{
		Enabled:        false,
		Triggering:     false,
		Severity:       restapi.SeverityWarning.GetAPIRepresentation(),
		Text:           ruleBindingText,
		Description:    ruleBindingDescription,
		ExpirationTime: 60000,
		Query:          ruleBindingQuery,
		RuleIds:        []string{ruldBindingRuleID},
	}

	client.EXPECT().Put(gomock.Eq(ruleBinding), gomock.Eq(restapi.RuleBindingsResourcePath)).Times(0)

	_, err := sut.Upsert(ruleBinding)

	if err == nil {
		t.Fatal(testutils.ExpectedError)
	}
}

func TestFailedUpsertOfRuleBindingBecauseOfInvalidResponseMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewRuleBindingResource(client)
	ruleBinding := makeTestRuleBinding()

	client.EXPECT().Put(gomock.Eq(ruleBinding), gomock.Eq(restapi.RuleBindingsResourcePath)).Return([]byte("invalid response"), nil)

	_, err := sut.Upsert(ruleBinding)

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
}

func TestFailedUpsertOfRuleBindingBecauseOfInvalidRuleBindingInResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewRuleBindingResource(client)
	ruleBinding := makeTestRuleBinding()

	client.EXPECT().Put(gomock.Eq(ruleBinding), gomock.Eq(restapi.RuleBindingsResourcePath)).Return([]byte("{ \"invalid\" : \"rule binding\" }"), nil)

	_, err := sut.Upsert(ruleBinding)

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
}

func TestFailedUpsertOfRuleBindingBecauseOfClientError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewRuleBindingResource(client)
	ruleBinding := makeTestRuleBinding()

	client.EXPECT().Put(gomock.Eq(ruleBinding), gomock.Eq(restapi.RuleBindingsResourcePath)).Return(nil, errors.New("Error during test"))

	_, err := sut.Upsert(ruleBinding)

	if err == nil {
		t.Fatal(testutils.ExpectedError)
	}
}

func TestSuccessfulDeleteOfRuleBindingByRuleBinding(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewRuleBindingResource(client)
	ruleBinding := makeTestRuleBinding()

	client.EXPECT().Delete(gomock.Eq(ruleBindingID), gomock.Eq(restapi.RuleBindingsResourcePath)).Return(nil)

	err := sut.Delete(ruleBinding)

	if err != nil {
		t.Fatalf("Expected no error got %s", err)
	}
}

func TestFailedDeleteOfRuleBindingByRuleBinding(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewRuleBindingResource(client)
	ruleBinding := makeTestRuleBinding()

	client.EXPECT().Delete(gomock.Eq(ruleBindingID), gomock.Eq(restapi.RuleBindingsResourcePath)).Return(errors.New("Error during test"))

	err := sut.Delete(ruleBinding)

	if err == nil {
		t.Fatal(testutils.ExpectedError)
	}
}

func makeTestRuleBinding() restapi.RuleBinding {
	return restapi.RuleBinding{
		ID:             ruleBindingID,
		Enabled:        false,
		Triggering:     false,
		Severity:       restapi.SeverityWarning.GetAPIRepresentation(),
		Text:           ruleBindingText,
		Description:    ruleBindingDescription,
		ExpirationTime: 60000,
		Query:          ruleBindingQuery,
		RuleIds:        []string{ruldBindingRuleID},
	}
}
