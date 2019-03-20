package resources_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi/resources"
	mocks "github.com/gessnerfl/terraform-provider-instana/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
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
		t.Fatalf("Expected no error but got %s", err)
	}

	if !cmp.Equal(ruleBinding, data) {
		t.Fatalf("Expected json to be unmarshalled to %v but got %v; diff %s", ruleBinding, data, cmp.Diff(ruleBinding, data))
	}
}

func TestFailedGetOneRuleBindingBecauseOfErrorFromRestClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewRuleBindingResource(client)
	ruleBindingID := "test-rule-binding-id"

	client.EXPECT().GetOne(gomock.Eq(ruleBindingID), gomock.Eq(restapi.RuleBindingsResourcePath)).Return(nil, errors.New("error during test"))

	_, err := sut.GetOne(ruleBindingID)

	if err == nil {
		t.Fatalf("Expected to get error")
	}
}

func TestFailedGetOneRuleBindingBecauseOfInvalidJsonArray(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewRuleBindingResource(client)
	ruleBindingID := "test-rule-binding-id"

	client.EXPECT().GetOne(gomock.Eq(ruleBindingID), gomock.Eq(restapi.RuleBindingsResourcePath)).Return([]byte("[{ \"invalid\" : \"data\" }]"), nil)

	_, err := sut.GetOne(ruleBindingID)

	if err == nil {
		t.Fatalf("Expected to get error")
	}
}

func TestFailedGetOneRuleBindingBecauseOfInvalidJsonObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewRuleBindingResource(client)
	ruleBindingID := "test-rule-binding-id"

	client.EXPECT().GetOne(gomock.Eq(ruleBindingID), gomock.Eq(restapi.RuleBindingsResourcePath)).Return([]byte("{ \"invalid\" : \"data\" }"), nil)

	_, err := sut.GetOne(ruleBindingID)

	if err == nil {
		t.Fatalf("Expected to get error")
	}
}

func TestFailedGetOneRuleBindingBecauseOfNoJsonAsResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewRuleBindingResource(client)
	ruleBindingID := "test-rule-binding-id"

	client.EXPECT().GetOne(gomock.Eq(ruleBindingID), gomock.Eq(restapi.RuleBindingsResourcePath)).Return([]byte("Invalid Data"), nil)

	_, err := sut.GetOne(ruleBindingID)

	if err == nil {
		t.Fatalf("Expected to get error")
	}
}

func TestSuccessfulGetAllRuleBindings(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewRuleBindingResource(client)
	ruleBinding1 := makeTestRuleBindingWithCounter(1)
	ruleBinding2 := makeTestRuleBindingWithCounter(2)
	ruleBindings := []restapi.RuleBinding{ruleBinding1, ruleBinding2}
	serializedJSON, _ := json.Marshal(ruleBindings)

	client.EXPECT().GetAll(gomock.Eq(restapi.RuleBindingsResourcePath)).Return(serializedJSON, nil)

	data, err := sut.GetAll()

	if err != nil {
		t.Fatalf("Expected no error but got %s", err)
	}

	if !cmp.Equal(ruleBindings, data) {
		t.Fatalf("Expected json to be unmarshalled to %v but got %v; diff %s", ruleBindings, data, cmp.Diff(ruleBindings, data))
	}
}

func TestFailedGetAllRuleBindingsWithErrorFromRestClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewRuleBindingResource(client)

	client.EXPECT().GetAll(gomock.Eq(restapi.RuleBindingsResourcePath)).Return(nil, errors.New("error during test"))

	_, err := sut.GetAll()

	if err == nil {
		t.Fatalf("Expected to get error")
	}
}

func TestFailedGetAllRuleBindingsWithInvalidJsonArray(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewRuleBindingResource(client)

	client.EXPECT().GetAll(gomock.Eq(restapi.RuleBindingsResourcePath)).Return([]byte("[{ \"invalid\" : \"data\" }]"), nil)

	_, err := sut.GetAll()

	if err == nil {
		t.Fatalf("Expected to get error")
	}
}

func TestFailedGetAllRuleBindingWithInvalidJsonObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewRuleBindingResource(client)

	client.EXPECT().GetAll(gomock.Eq(restapi.RuleBindingsResourcePath)).Return([]byte("{ \"invalid\" : \"data\" }"), nil)

	_, err := sut.GetAll()

	if err == nil {
		t.Fatalf("Expected to get error")
	}
}

func TestFailedGetAllRuleBindingsWithNoJsonAsResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewRuleBindingResource(client)

	client.EXPECT().GetAll(gomock.Eq(restapi.RuleBindingsResourcePath)).Return([]byte("Invalid Data"), nil)

	_, err := sut.GetAll()

	if err == nil {
		t.Fatalf("Expected to get error")
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
		t.Fatalf("Expected no error but got %s", err)
	}

	if !cmp.Equal(ruleBinding, result) {
		t.Fatalf("Expected json to be unmarshalled to %v but got %v; diff %s", ruleBinding, result, cmp.Diff(ruleBinding, result))
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
		Severity:       1,
		Text:           "test-text",
		Description:    "test-description",
		ExpirationTime: 60000,
		Query:          "entity.type:jvm",
		RuleIds:        []string{"test-rule-id"},
	}

	client.EXPECT().Put(gomock.Eq(ruleBinding), gomock.Eq(restapi.RuleBindingsResourcePath)).Times(0)

	_, err := sut.Upsert(ruleBinding)

	if err == nil {
		t.Fatal("Expected to get error")
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
		t.Fatalf("Expected to get error")
	}
}

func TestFailedUpsertOfRuleBindingBecauseOfInvalidRuleInResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewRuleBindingResource(client)
	ruleBinding := makeTestRuleBinding()

	client.EXPECT().Put(gomock.Eq(ruleBinding), gomock.Eq(restapi.RuleBindingsResourcePath)).Return([]byte("{ \"invalid\" : \"rule\" }"), nil)

	_, err := sut.Upsert(ruleBinding)

	if err == nil {
		t.Fatalf("Expected to get error")
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
		t.Fatal("Expected to get error")
	}
}

func TestSuccessfulDeleteOfRuleBindingByRuleBinding(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewRuleBindingResource(client)
	ruleBinding := makeTestRuleBinding()

	client.EXPECT().Delete(gomock.Eq("test-rule-binding-id-1"), gomock.Eq(restapi.RuleBindingsResourcePath)).Return(nil)

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

	client.EXPECT().Delete(gomock.Eq("test-rule-binding-id-1"), gomock.Eq(restapi.RuleBindingsResourcePath)).Return(errors.New("Error during test"))

	err := sut.Delete(ruleBinding)

	if err == nil {
		t.Fatal("Expected to get error")
	}
}

func makeTestRuleBinding() restapi.RuleBinding {
	return makeTestRuleBindingWithCounter(1)
}

func makeTestRuleBindingWithCounter(counter int) restapi.RuleBinding {
	id := fmt.Sprintf("test-rule-binding-id-%d", counter)
	text := fmt.Sprintf("Test Rule Binding Text %d", counter)
	description := fmt.Sprintf("Test Rule Binding Description %d", counter)
	return restapi.RuleBinding{
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
