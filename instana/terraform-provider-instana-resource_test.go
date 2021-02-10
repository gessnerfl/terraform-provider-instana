package instana_test

import (
	"errors"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/mocks"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
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

		resourceHandle := NewAlertingChannelEmailResourceHandle()
		err := NewTerraformResource(resourceHandle).Read(resourceData, providerMeta)

		assert.Nil(t, err)
		verifyTestObjectModelAppliedToResource(expectedModel, resourceData, t)
	})
}

func TestShouldFailToReadTestObjectFromInstanaAPIWhenIDIsMissing(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		resourceData := createEmptyAlertingChannelEmailResourceData(t)

		resourceHandle := NewAlertingChannelEmailResourceHandle()
		err := NewTerraformResource(resourceHandle).Read(resourceData, providerMeta)

		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "ID of instana_alerting_channel_email")
	})
}

func TestShouldFailToReadTestObjectFromInstanaAPIAndDeleteResourceWhenRoleDoesNotExist(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		resourceData := createEmptyAlertingChannelEmailResourceData(t)
		resourceData.SetId(alertingChannelEmailID)
		mockTestObjectApi := mocks.NewMockRestResource(ctrl)

		mockInstanaAPI.EXPECT().AlertingChannels().Return(mockTestObjectApi).Times(1)
		mockTestObjectApi.EXPECT().GetOne(gomock.Eq(alertingChannelEmailID)).Return(&restapi.AlertingChannel{}, restapi.ErrEntityNotFound).Times(1)

		resourceHandle := NewAlertingChannelEmailResourceHandle()
		err := NewTerraformResource(resourceHandle).Read(resourceData, providerMeta)

		assert.Nil(t, err)
		assert.GreaterOrEqual(t, 0, len(resourceData.Id()))
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
		mockTestObjectApi.EXPECT().GetOne(gomock.Eq(alertingChannelEmailID)).Return(&restapi.AlertingChannel{}, expectedError).Times(1)

		resourceHandle := NewAlertingChannelEmailResourceHandle()
		err := NewTerraformResource(resourceHandle).Read(resourceData, providerMeta)

		assert.Equal(t, expectedError, err)
		assert.NotEqual(t, 0, len(resourceData.Id()))
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
		mockTestObjectApi.EXPECT().Create(gomock.AssignableToTypeOf(&restapi.AlertingChannel{})).Return(expectedModel, nil).Times(1)

		resourceHandle := NewAlertingChannelEmailResourceHandle()
		err := NewTerraformResource(resourceHandle).Create(resourceData, providerMeta)

		assert.Nil(t, err)
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
		mockTestObjectApi.EXPECT().Create(gomock.AssignableToTypeOf(&restapi.AlertingChannel{})).Return(&restapi.AlertingChannel{}, expectedError).Times(1)

		resourceHandle := NewAlertingChannelEmailResourceHandle()
		err := NewTerraformResource(resourceHandle).Create(resourceData, providerMeta)

		assert.Equal(t, expectedError, err)
	})
}

func TestShouldUpdateTestObjectThroughInstanaAPI(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		data := createTestAlertingChannelEmailData()
		resourceData := createAlertingChannelEmailResourceData(data, t)
		expectedModel := createTestAlertingChannelEmailObject()
		mockTestObjectApi := mocks.NewMockRestResource(ctrl)

		mockInstanaAPI.EXPECT().AlertingChannels().Return(mockTestObjectApi).Times(1)
		mockResourceNameFormatter.EXPECT().Format(data[AlertingChannelFieldName]).Return(data[AlertingChannelFieldName]).Times(1)
		mockTestObjectApi.EXPECT().Update(gomock.AssignableToTypeOf(&restapi.AlertingChannel{})).Return(expectedModel, nil).Times(1)

		resourceHandle := NewAlertingChannelEmailResourceHandle()
		err := NewTerraformResource(resourceHandle).Update(resourceData, providerMeta)

		assert.Nil(t, err)
		verifyTestObjectModelAppliedToResource(expectedModel, resourceData, t)
	})
}

func TestShouldReturnErrorWhenUpdateTestObjectFailsThroughInstanaAPI(t *testing.T) {
	testHelper := NewTestHelper(t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI, mockResourceNameFormatter *mocks.MockResourceNameFormatter) {
		data := createTestAlertingChannelEmailData()
		resourceData := createAlertingChannelEmailResourceData(data, t)
		expectedError := errors.New("test")
		mockTestObjectApi := mocks.NewMockRestResource(ctrl)

		mockInstanaAPI.EXPECT().AlertingChannels().Return(mockTestObjectApi).Times(1)
		mockResourceNameFormatter.EXPECT().Format(data[AlertingChannelFieldName]).Return(data[AlertingChannelFieldName]).Times(1)
		mockTestObjectApi.EXPECT().Update(gomock.AssignableToTypeOf(&restapi.AlertingChannel{})).Return(&restapi.AlertingChannel{}, expectedError).Times(1)

		resourceHandle := NewAlertingChannelEmailResourceHandle()
		err := NewTerraformResource(resourceHandle).Update(resourceData, providerMeta)

		assert.Equal(t, expectedError, err)
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

		resourceHandle := NewAlertingChannelEmailResourceHandle()
		err := NewTerraformResource(resourceHandle).Delete(resourceData, providerMeta)

		assert.Nil(t, err)
		assert.GreaterOrEqual(t, 0, len(resourceData.Id()))
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

		resourceHandle := NewAlertingChannelEmailResourceHandle()
		err := NewTerraformResource(resourceHandle).Delete(resourceData, providerMeta)

		assert.Equal(t, expectedError, err)
		assert.NotEqual(t, 0, len(resourceData.Id()))
	})
}

func verifyTestObjectModelAppliedToResource(model *restapi.AlertingChannel, resourceData *schema.ResourceData, t *testing.T) {
	assert.Equal(t, model.ID, resourceData.Id())
	assert.Equal(t, model.Name, resourceData.Get(AlertingChannelFieldFullName))

	emails := ReadStringSetParameterFromResource(resourceData, AlertingChannelEmailFieldEmails)
	assert.Equal(t, len(model.Emails), len(emails))
	for _, mail := range model.Emails {
		assert.Contains(t, emails, mail)
	}
}

func createTestAlertingChannelEmailObject() *restapi.AlertingChannel {
	return &restapi.AlertingChannel{
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
	schemaMap := NewAlertingChannelEmailResourceHandle().MetaData().Schema
	return schema.TestResourceDataRaw(t, schemaMap, data)
}
