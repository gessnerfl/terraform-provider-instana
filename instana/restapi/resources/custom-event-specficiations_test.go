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
	customEventID           = "custom-event-id"
	customEventName         = "custom-event-name"
	customEventEntityType   = "custom-event-entity-type"
	customEventQuery        = "custom-event-query"
	customEventDescription  = "custom-event-description"
	customEventSystemRuleID = "system-rule-id"
)

func TestSuccessfulGetOneCustomEventSpecification(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewCustomEventSpecificationResource(client)
	customEvent := makeTestCustomEventSpecification()
	serializedJSON, _ := json.Marshal(customEvent)

	client.EXPECT().GetOne(gomock.Eq(customEvent.ID), gomock.Eq(restapi.CustomEventSpecificationResourcePath)).Return(serializedJSON, nil)

	data, err := sut.GetOne(customEvent.ID)

	if err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}

	if !cmp.Equal(customEvent, data) {
		t.Fatalf(testutils.ExpectedUnmarshalledJSONWithStruct, customEvent, data, cmp.Diff(customEvent, data))
	}
}

func TestFailedGetOneCustomEventSpecificationBecauseOfErrorFromRestClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewCustomEventSpecificationResource(client)

	client.EXPECT().GetOne(gomock.Eq(customEventID), gomock.Eq(restapi.CustomEventSpecificationResourcePath)).Return(nil, errors.New("error during test"))

	_, err := sut.GetOne(customEventID)

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
}

func TestFailedGetOneCustomEventSpecificationBecauseOfInvalidJsonArray(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewCustomEventSpecificationResource(client)

	client.EXPECT().GetOne(gomock.Eq(customEventID), gomock.Eq(restapi.CustomEventSpecificationResourcePath)).Return([]byte("[{ \"invalid\" : \"data\" }]"), nil)

	_, err := sut.GetOne(customEventID)

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
}

func TestFailedGetOneCustomEventSpecificationBecauseOfInvalidJsonObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewCustomEventSpecificationResource(client)

	client.EXPECT().GetOne(gomock.Eq(customEventID), gomock.Eq(restapi.CustomEventSpecificationResourcePath)).Return([]byte("{ \"invalid\" : \"data\" }"), nil)

	_, err := sut.GetOne(customEventID)

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
}

func TestFailedGetOneCustomEventSpecificationBecauseOfNoJsonAsResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewCustomEventSpecificationResource(client)

	client.EXPECT().GetOne(gomock.Eq(customEventID), gomock.Eq(restapi.CustomEventSpecificationResourcePath)).Return([]byte("Invalid Data"), nil)

	_, err := sut.GetOne(customEventID)

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
}

func TestSuccessfulUpsertOfCustomEventSpecification(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewCustomEventSpecificationResource(client)
	customEvent := makeTestCustomEventSpecification()
	serializedJSON, _ := json.Marshal(customEvent)

	client.EXPECT().Put(gomock.Eq(customEvent), gomock.Eq(restapi.CustomEventSpecificationResourcePath)).Return(serializedJSON, nil)

	result, err := sut.Upsert(customEvent)

	if err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}

	if !cmp.Equal(customEvent, result) {
		t.Fatalf(testutils.ExpectedUnmarshalledJSONWithStruct, customEvent, result, cmp.Diff(customEvent, result))
	}
}

func TestFailedUpsertOfCustomEventSpecificationBecauseOfInvalidCustomEventSpecification(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewCustomEventSpecificationResource(client)
	systemRule := restapi.NewSystemRuleSpecification(customEventSystemRuleID, restapi.SeverityWarning.GetAPIRepresentation())
	customEvent := restapi.CustomEventSpecification{
		Name:       customEventName,
		EntityType: customEventEntityType,
		Rules:      []restapi.RuleSpecification{systemRule},
	}

	client.EXPECT().Put(gomock.Eq(customEvent), gomock.Eq(restapi.CustomEventSpecificationResourcePath)).Times(0)

	_, err := sut.Upsert(customEvent)

	if err == nil {
		t.Fatal(testutils.ExpectedError)
	}
}

func TestFailedUpsertOfCustomEventSpecificationBecauseOfInvalidResponseMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewCustomEventSpecificationResource(client)
	customEvent := makeTestCustomEventSpecification()

	client.EXPECT().Put(gomock.Eq(customEvent), gomock.Eq(restapi.CustomEventSpecificationResourcePath)).Return([]byte("invalid response"), nil)

	_, err := sut.Upsert(customEvent)

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
}

func TestFailedUpsertOfCustomEventSpecificationBecauseOfInvalidCustomEventSpecificationInResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewCustomEventSpecificationResource(client)
	customEvent := makeTestCustomEventSpecification()

	client.EXPECT().Put(gomock.Eq(customEvent), gomock.Eq(restapi.CustomEventSpecificationResourcePath)).Return([]byte("{ \"invalid\" : \"rule binding\" }"), nil)

	_, err := sut.Upsert(customEvent)

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
}

func TestFailedUpsertOfCustomEventSpecificationBecauseOfClientError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewCustomEventSpecificationResource(client)
	customEvent := makeTestCustomEventSpecification()

	client.EXPECT().Put(gomock.Eq(customEvent), gomock.Eq(restapi.CustomEventSpecificationResourcePath)).Return(nil, errors.New("Error during test"))

	_, err := sut.Upsert(customEvent)

	if err == nil {
		t.Fatal(testutils.ExpectedError)
	}
}

func TestSuccessfulDeleteOfCustomEventSpecificationByCustomEventSpecification(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewCustomEventSpecificationResource(client)
	customEvent := makeTestCustomEventSpecification()

	client.EXPECT().Delete(gomock.Eq(customEventID), gomock.Eq(restapi.CustomEventSpecificationResourcePath)).Return(nil)

	err := sut.Delete(customEvent)

	if err != nil {
		t.Fatalf("Expected no error got %s", err)
	}
}

func TestFailedDeleteOfCustomEventSpecificationByCustomEventSpecification(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mocks.NewMockRestClient(ctrl)
	sut := NewCustomEventSpecificationResource(client)
	customEvent := makeTestCustomEventSpecification()

	client.EXPECT().Delete(gomock.Eq(customEventID), gomock.Eq(restapi.CustomEventSpecificationResourcePath)).Return(errors.New("Error during test"))

	err := sut.Delete(customEvent)

	if err == nil {
		t.Fatal(testutils.ExpectedError)
	}
}

func makeTestCustomEventSpecification() restapi.CustomEventSpecification {
	description := customEventDescription
	query := customEventQuery
	expirationTime := 60000
	systemRule := restapi.NewSystemRuleSpecification(customEventSystemRuleID, restapi.SeverityWarning.GetAPIRepresentation())
	return restapi.CustomEventSpecification{
		ID:             customEventID,
		Name:           customEventName,
		EntityType:     customEventEntityType,
		Enabled:        true,
		Triggering:     true,
		Description:    &description,
		ExpirationTime: &expirationTime,
		Query:          &query,
		Rules:          []restapi.RuleSpecification{systemRule},
	}
}
