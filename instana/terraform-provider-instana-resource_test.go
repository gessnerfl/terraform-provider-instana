package instana_test

import (
	"errors"
	"strings"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/mocks"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform/helper/schema"
)

const alertingChannelEmailID = "id"

func TestShouldSuccessfullyReadTestObjectFromInstanaAPIWhenBaseDataIsReturned(t *testing.T) {
	expectedModel := createTestAlertingChannelEmailObject()
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		resourceData := createEmptyAlertingChannelEmailResourceData(t)
		resourceData.SetId(alertingChannelEmailID)
		mockTestObjectApi := mocks.NewMockRestResource(ctrl)

		mockInstanaAPI.EXPECT().AlertingChannels().Return(mockTestObjectApi).Times(1)
		mockTestObjectApi.EXPECT().GetOne(gomock.Eq(alertingChannelEmailID)).Return(expectedModel, nil).Times(1)

		err := NewTerraformResource(NewAlertingChannelEmailResourceHandle()).Read(resourceData, providerMeta)

		if err != nil {
			t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
		}
		verifyTestObjectModelAppliedToResource(expectedModel, resourceData, t)
	})
}

func TestShouldFailToReadTestObjectFromInstanaAPIWhenIDIsMissing(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		resourceData := createEmptyAlertingChannelEmailResourceData(t)

		err := NewTerraformResource(NewAlertingChannelEmailResourceHandle()).Read(resourceData, providerMeta)

		if err == nil || !strings.HasPrefix(err.Error(), "ID of instana_alerting_channel_email") {
			t.Fatal("Expected error to occur because of missing id")
		}
	})
}

func TestShouldFailToReadTestObjectFromInstanaAPIAndDeleteResourceWhenRoleDoesNotExist(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		resourceData := createEmptyAlertingChannelEmailResourceData(t)
		resourceData.SetId(alertingChannelEmailID)
		mockTestObjectApi := mocks.NewMockRestResource(ctrl)

		mockInstanaAPI.EXPECT().AlertingChannels().Return(mockTestObjectApi).Times(1)
		mockTestObjectApi.EXPECT().GetOne(gomock.Eq(alertingChannelEmailID)).Return(restapi.AlertingChannel{}, restapi.ErrEntityNotFound).Times(1)

		err := NewTerraformResource(NewAlertingChannelEmailResourceHandle()).Read(resourceData, providerMeta)

		if err != nil {
			t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
		}
		if len(resourceData.Id()) > 0 {
			t.Fatal("Expected ID to be cleaned to destroy resource")
		}
	})
}

func TestShouldFailToReadTestObjectFromInstanaAPIAndReturnErrorWhenAPICallFails(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		resourceData := createEmptyAlertingChannelEmailResourceData(t)
		resourceData.SetId(alertingChannelEmailID)
		expectedError := errors.New("test")
		mockTestObjectApi := mocks.NewMockRestResource(ctrl)

		mockInstanaAPI.EXPECT().AlertingChannels().Return(mockTestObjectApi).Times(1)
		mockTestObjectApi.EXPECT().GetOne(gomock.Eq(alertingChannelEmailID)).Return(restapi.AlertingChannel{}, expectedError).Times(1)

		err := NewTerraformResource(NewAlertingChannelEmailResourceHandle()).Read(resourceData, providerMeta)

		if err == nil || err != expectedError {
			t.Fatal("Expected error should be returned")
		}
		if len(resourceData.Id()) == 0 {
			t.Fatal("Expected ID should still be set")
		}
	})
}

func TestShouldCreateTestObjectThroughInstanaAPI(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		data := createTestAlertingChannelEmailData()
		resourceData := createAlertingChannelEmailResourceData(data, t)
		expectedModel := createTestAlertingChannelEmailObject()
		mockTestObjectApi := mocks.NewMockRestResource(ctrl)

		mockInstanaAPI.EXPECT().AlertingChannels().Return(mockTestObjectApi).Times(1)
		mockResourceNameFormatter.EXPECT().Format(data[AlertingChannelFieldName]).Return(data[AlertingChannelFieldName]).Times(1)
		mockTestObjectApi.EXPECT().Upsert(gomock.AssignableToTypeOf(restapi.AlertingChannel{})).Return(expectedModel, nil).Times(1)

		err := NewTerraformResource(NewAlertingChannelEmailResourceHandle()).Create(resourceData, providerMeta)

		if err != nil {
			t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
		}
		verifyTestObjectModelAppliedToResource(expectedModel, resourceData, t)
	})
}

func TestShouldReturnErrorWhenCreateTestObjectFailsThroughInstanaAPI(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		data := createTestAlertingChannelEmailData()
		resourceData := createAlertingChannelEmailResourceData(data, t)
		expectedError := errors.New("test")
		mockTestObjectApi := mocks.NewMockRestResource(ctrl)

		mockInstanaAPI.EXPECT().AlertingChannels().Return(mockTestObjectApi).Times(1)
		mockResourceNameFormatter.EXPECT().Format(data[AlertingChannelFieldName]).Return(data[AlertingChannelFieldName]).Times(1)
		mockTestObjectApi.EXPECT().Upsert(gomock.AssignableToTypeOf(restapi.AlertingChannel{})).Return(restapi.AlertingChannel{}, expectedError).Times(1)

		err := NewTerraformResource(NewAlertingChannelEmailResourceHandle()).Create(resourceData, providerMeta)

		if err == nil || expectedError != err {
			t.Fatal("Expected definned error to be returned")
		}
	})
}

func TestShouldDeleteTestObjectThroughInstanaAPI(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		id := "test-id"
		data := createTestAlertingChannelEmailData()
		resourceData := createAlertingChannelEmailResourceData(data, t)
		resourceData.SetId(id)
		mockTestObjectApi := mocks.NewMockRestResource(ctrl)

		mockInstanaAPI.EXPECT().AlertingChannels().Return(mockTestObjectApi).Times(1)
		mockResourceNameFormatter.EXPECT().Format(data[AlertingChannelFieldName]).Return(data[AlertingChannelFieldName]).Times(1)
		mockTestObjectApi.EXPECT().DeleteByID(gomock.Eq(id)).Return(nil).Times(1)

		err := NewTerraformResource(NewAlertingChannelEmailResourceHandle()).Delete(resourceData, providerMeta)

		if err != nil {
			t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
		}
		if len(resourceData.Id()) > 0 {
			t.Fatal("Expected ID to be cleaned to destroy resource")
		}
	})
}

func TestShouldReturnErrorWhenDeleteTestObjectFailsThroughInstanaAPI(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		id := "test-id"
		data := createTestAlertingChannelEmailData()
		resourceData := createAlertingChannelEmailResourceData(data, t)
		resourceData.SetId(id)
		expectedError := errors.New("test")
		mockTestObjectApi := mocks.NewMockRestResource(ctrl)

		mockInstanaAPI.EXPECT().AlertingChannels().Return(mockTestObjectApi).Times(1)
		mockResourceNameFormatter.EXPECT().Format(data[AlertingChannelFieldName]).Return(data[AlertingChannelFieldName]).Times(1)
		mockTestObjectApi.EXPECT().DeleteByID(gomock.Eq(id)).Return(expectedError).Times(1)

		err := NewTerraformResource(NewAlertingChannelEmailResourceHandle()).Delete(resourceData, providerMeta)

		if err == nil || err != expectedError {
			t.Fatal("Expected error to be returned")
		}
		if len(resourceData.Id()) == 0 {
			t.Fatal("Expected ID not to be cleaned to avoid resource is destroy")
		}
	})
}

func verifyTestObjectModelAppliedToResource(model restapi.AlertingChannel, resourceData *schema.ResourceData, t *testing.T) {
	if model.ID != resourceData.Id() {
		t.Fatal("Expected ID to be identical")
	}
	if model.Name != resourceData.Get(AlertingChannelFieldFullName).(string) {
		t.Fatal("Expected Full Name to match Name of API")
	}
	mails := ReadStringArrayParameterFromResource(resourceData, AlertingChannelEmailFieldEmails)
	if !cmp.Equal(model.Emails, mails) {
		t.Fatalf("Expected Emails to be identical %s vs %s", model.Emails, mails)
	}
}

func createTestAlertingChannelEmailObject() restapi.AlertingChannel {
	return restapi.AlertingChannel{
		ID:     "id",
		Name:   "name",
		Emails: []string{"Email1", "Email2"},
	}
}

func createTestAlertingChannelEmailData() map[string]interface{} {
	emails := make([]interface{}, 2)
	emails[0] = "Email1"
	emails[1] = "Email2"

	data := make(map[string]interface{})
	data[AlertingChannelFieldName] = "name"
	data[AlertingChannelEmailFieldEmails] = emails
	return data
}

func createEmptyAlertingChannelEmailResourceData(t *testing.T) *schema.ResourceData {
	data := make(map[string]interface{})
	return createAlertingChannelEmailResourceData(data, t)
}

func createAlertingChannelEmailResourceData(data map[string]interface{}, t *testing.T) *schema.ResourceData {
	schemaMap := NewAlertingChannelEmailResourceHandle().GetSchema()
	return schema.TestResourceDataRaw(t, schemaMap, data)
}
