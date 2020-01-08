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
	alertingChannelID         = "test-alerting-channel-id"
	alertingChannelName       = "Test Alerting Channel"
	alertingChannelWebhookUrl = "https://webhook.example.com/test"
)

func TestSuccessfulGetOneAlertingChannel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewAlertingChannelResource(client)
	alertingChannel := makeTestAlertingChannel()
	serializedJSON, _ := json.Marshal(alertingChannel)

	client.EXPECT().GetOne(gomock.Eq(alertingChannel.ID), gomock.Eq(restapi.AlertingChannelsResourcePath)).Return(serializedJSON, nil)

	data, err := sut.GetOne(alertingChannel.ID)

	if err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}

	if !cmp.Equal(alertingChannel, data) {
		t.Fatalf(testutils.ExpectedUnmarshalledJSONWithStruct, alertingChannel, data, cmp.Diff(alertingChannel, data))
	}
}

func TestFailedGetOneAlertingChannelBecauseOfErrorFromRestClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewAlertingChannelResource(client)

	client.EXPECT().GetOne(gomock.Eq(alertingChannelID), gomock.Eq(restapi.AlertingChannelsResourcePath)).Return(nil, errors.New("error during test"))

	_, err := sut.GetOne(alertingChannelID)

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
}

func TestFailedGetOneAlertingChannelBecauseOfInvalidJsonArray(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewAlertingChannelResource(client)

	client.EXPECT().GetOne(gomock.Eq(alertingChannelID), gomock.Eq(restapi.AlertingChannelsResourcePath)).Return([]byte("[{ \"invalid\" : \"data\" }]"), nil)

	_, err := sut.GetOne(alertingChannelID)

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
}

func TestFailedGetOneAlertingChannelBecauseOfInvalidJsonObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewAlertingChannelResource(client)

	client.EXPECT().GetOne(gomock.Eq(alertingChannelID), gomock.Eq(restapi.AlertingChannelsResourcePath)).Return([]byte("{ \"invalid\" : \"data\" }"), nil)

	_, err := sut.GetOne(alertingChannelID)

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
}

func TestFailedGetOneAlertingChannelBecauseResponseIsNotAValidJsonDocument(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewAlertingChannelResource(client)

	client.EXPECT().GetOne(gomock.Eq(alertingChannelID), gomock.Eq(restapi.AlertingChannelsResourcePath)).Return([]byte("Invalid Data"), nil)

	_, err := sut.GetOne(alertingChannelID)

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
}

func TestSuccessfulUpsertOfAlertingChannel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewAlertingChannelResource(client)
	alertingChannel := makeTestAlertingChannel()
	serializedJSON, _ := json.Marshal(alertingChannel)

	client.EXPECT().Put(gomock.Eq(alertingChannel), gomock.Eq(restapi.AlertingChannelsResourcePath)).Return(serializedJSON, nil)

	result, err := sut.Upsert(alertingChannel)

	if err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}

	if !cmp.Equal(alertingChannel, result) {
		t.Fatalf(testutils.ExpectedUnmarshalledJSONWithStruct, alertingChannel, result, cmp.Diff(result, result))
	}
}

func TestFailedUpsertOfAlertingChannelBecauseOfClientError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewAlertingChannelResource(client)
	alertingChannel := makeTestAlertingChannel()

	client.EXPECT().Put(gomock.Eq(alertingChannel), gomock.Eq(restapi.AlertingChannelsResourcePath)).Return(nil, errors.New("Error during test"))

	_, err := sut.Upsert(alertingChannel)

	if err == nil {
		t.Fatal(testutils.ExpectedError)
	}
}

func TestFailedUpsertOfAlertingChannelBecauseOfInvalidResponseMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewAlertingChannelResource(client)
	alertingChannel := makeTestAlertingChannel()

	client.EXPECT().Put(gomock.Eq(alertingChannel), gomock.Eq(restapi.AlertingChannelsResourcePath)).Return([]byte("invalid response"), nil)

	_, err := sut.Upsert(alertingChannel)

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
}

func TestFailedUpsertOfAlertingChannelBecauseOfInvalidAlertingChannelInResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewAlertingChannelResource(client)
	alertingChannel := makeTestAlertingChannel()

	client.EXPECT().Put(gomock.Eq(alertingChannel), gomock.Eq(restapi.AlertingChannelsResourcePath)).Return([]byte("{ \"invalid\" : \"alertingChannel\" }"), nil)

	_, err := sut.Upsert(alertingChannel)

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
}

func TestFailedUpsertOfAlertingChannelBecauseOfInvalidAlertingChannelProvided(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewAlertingChannelResource(client)
	alertingChannel := restapi.AlertingChannel{
		Name: alertingChannelName,
	}

	client.EXPECT().Put(gomock.Eq(alertingChannel), gomock.Eq(restapi.AlertingChannelsResourcePath)).Times(0)

	_, err := sut.Upsert(alertingChannel)

	if err == nil {
		t.Fatal(testutils.ExpectedError)
	}
}

func TestSuccessfulDeleteOfAlertingChannelByAlertingChannel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewAlertingChannelResource(client)
	alertingChannel := makeTestAlertingChannel()

	client.EXPECT().Delete(gomock.Eq(alertingChannelID), gomock.Eq(restapi.AlertingChannelsResourcePath)).Return(nil)

	err := sut.Delete(alertingChannel)

	if err != nil {
		t.Fatalf("Expected no error got %s", err)
	}
}

func TestFailedDeleteOfAlertingChannelByAlertingChannel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := NewAlertingChannelResource(client)
	alertingChannel := makeTestAlertingChannel()

	client.EXPECT().Delete(gomock.Eq(alertingChannelID), gomock.Eq(restapi.AlertingChannelsResourcePath)).Return(errors.New("Error during test"))

	err := sut.Delete(alertingChannel)

	if err == nil {
		t.Fatal(testutils.ExpectedError)
	}
}

func makeTestAlertingChannel() restapi.AlertingChannel {
	webhookUrl := alertingChannelWebhookUrl
	return restapi.AlertingChannel{
		ID:         alertingChannelID,
		Name:       alertingChannelName,
		Kind:       restapi.Office365ChannelType,
		WebhookURL: &webhookUrl,
	}
}
