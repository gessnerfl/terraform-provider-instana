package instana_test

import (
	"context"
	"errors"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana"
	"github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/mocks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

const alertingChannelEmailID = "id"
const resourceNameWithoutPrefixAndSuffix = "name"

func TestTerraformProviderInstanaResource(t *testing.T) {
	ut := &terraformProviderInstanaResourceUnitTest{}
	t.Run("should successfully read test object from instana API when base data is returned", ut.shouldSuccessfullyReadTestObjectFromInstanaAPIWhenBaseDataIsReturned)
	t.Run("should fail to read test object from instana API when id is missing", ut.shouldFailToReadTestObjectFromInstanaAPIWhenResourceIDIsMissing)
	t.Run("should fail to read test object from instana API and delete resource when role does not exist", ut.shouldFailToReadTestObjectFromInstanaAPIAndDeleteResourceWhenRoleDoesNotExist)
	t.Run("should fail to read test object from instana API and return error code when API call fails", ut.shouldFailToReadTestObjectFromInstanaAPIAndReturnErrorWhenAPICallFails)
	t.Run("should create test object through Instana API", ut.shouldCreateTestObjectThroughInstanaAPI)
	t.Run("should return error when create test object fails through Instana API", ut.shouldReturnErrorWhenCreateTestObjectFailsThroughInstanaAPI)
	t.Run("should update test object through Instana API", ut.shouldUpdateTestObjectThroughInstanaAPI)
	t.Run("should return error when update test object fails through Instana API", ut.shouldReturnErrorWhenUpdateTestObjectFailsThroughInstanaAPI)
	t.Run("should delete test object through Instana API", ut.shouldDeleteTestObjectThroughInstanaAPI)
	t.Run("should return error when delete test object fails through Instana API", ut.shouldReturnErrorWhenDeleteTestObjectFailsThroughInstanaAPI)
}

type terraformProviderInstanaResourceUnitTest struct{}

func (r *terraformProviderInstanaResourceUnitTest) shouldSuccessfullyReadTestObjectFromInstanaAPIWhenBaseDataIsReturned(t *testing.T) {
	expectedModel := r.createTestAlertingChannelEmailObject()
	testHelper := NewTestHelper[*restapi.AlertingChannel](t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI) {
		resourceData := r.createEmptyAlertingChannelResourceData(t)
		resourceData.SetId(alertingChannelEmailID)
		mockTestObjectApi := mocks.NewMockRestResource[*restapi.AlertingChannel](ctrl)

		mockInstanaAPI.EXPECT().AlertingChannels().Return(mockTestObjectApi).Times(1)
		mockTestObjectApi.EXPECT().GetOne(gomock.Eq(alertingChannelEmailID)).Return(expectedModel, nil).Times(1)

		resourceHandle := NewAlertingChannelResourceHandle()
		diag := NewTerraformResource(resourceHandle).Read(context.TODO(), resourceData, providerMeta)

		assert.Nil(t, diag)
		r.verifyTestObjectModelAppliedToResource(expectedModel, resourceData, t)
	})
}

func (r *terraformProviderInstanaResourceUnitTest) shouldFailToReadTestObjectFromInstanaAPIWhenResourceIDIsMissing(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AlertingChannel](t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI) {
		resourceData := r.createEmptyAlertingChannelResourceData(t)

		resourceHandle := NewAlertingChannelResourceHandle()
		diag := NewTerraformResource(resourceHandle).Read(context.TODO(), resourceData, providerMeta)

		assert.NotNil(t, diag)
		assert.True(t, diag.HasError())
		assert.Contains(t, diag[0].Summary, "ID of instana_alerting_channel")
	})
}

func (r *terraformProviderInstanaResourceUnitTest) shouldFailToReadTestObjectFromInstanaAPIAndDeleteResourceWhenRoleDoesNotExist(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AlertingChannel](t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI) {
		resourceData := r.createEmptyAlertingChannelResourceData(t)
		resourceData.SetId(alertingChannelEmailID)
		mockTestObjectApi := mocks.NewMockRestResource[*restapi.AlertingChannel](ctrl)

		mockInstanaAPI.EXPECT().AlertingChannels().Return(mockTestObjectApi).Times(1)
		mockTestObjectApi.EXPECT().GetOne(gomock.Eq(alertingChannelEmailID)).Return(&restapi.AlertingChannel{}, restapi.ErrEntityNotFound).Times(1)

		resourceHandle := NewAlertingChannelResourceHandle()
		diag := NewTerraformResource(resourceHandle).Read(context.TODO(), resourceData, providerMeta)

		assert.Nil(t, diag)
		assert.GreaterOrEqual(t, 0, len(resourceData.Id()))
	})
}

func (r *terraformProviderInstanaResourceUnitTest) shouldFailToReadTestObjectFromInstanaAPIAndReturnErrorWhenAPICallFails(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AlertingChannel](t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI) {
		resourceData := r.createEmptyAlertingChannelResourceData(t)
		resourceData.SetId(alertingChannelEmailID)
		expectedError := errors.New("test")
		mockTestObjectApi := mocks.NewMockRestResource[*restapi.AlertingChannel](ctrl)

		mockInstanaAPI.EXPECT().AlertingChannels().Return(mockTestObjectApi).Times(1)
		mockTestObjectApi.EXPECT().GetOne(gomock.Eq(alertingChannelEmailID)).Return(&restapi.AlertingChannel{}, expectedError).Times(1)

		resourceHandle := NewAlertingChannelResourceHandle()
		diag := NewTerraformResource(resourceHandle).Read(context.TODO(), resourceData, providerMeta)

		assert.NotNil(t, diag)
		assert.True(t, diag.HasError())
		assert.Equal(t, diag[0].Summary, expectedError.Error())
		assert.NotEqual(t, 0, len(resourceData.Id()))
	})
}

func (r *terraformProviderInstanaResourceUnitTest) shouldCreateTestObjectThroughInstanaAPI(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AlertingChannel](t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI) {
		data := r.createTestAlertingChannelEmailData()
		resourceData := r.createAlertingChannelResourceData(data, t)
		expectedModel := r.createTestAlertingChannelEmailObject()
		mockTestObjectApi := mocks.NewMockRestResource[*restapi.AlertingChannel](ctrl)

		mockInstanaAPI.EXPECT().AlertingChannels().Return(mockTestObjectApi).Times(1)
		mockTestObjectApi.EXPECT().Create(gomock.AssignableToTypeOf(&restapi.AlertingChannel{})).Return(expectedModel, nil).Times(1)

		resourceHandle := NewAlertingChannelResourceHandle()
		diag := NewTerraformResource(resourceHandle).Create(context.TODO(), resourceData, providerMeta)

		assert.Nil(t, diag)
		r.verifyTestObjectModelAppliedToResource(expectedModel, resourceData, t)
	})
}

func (r *terraformProviderInstanaResourceUnitTest) shouldReturnErrorWhenCreateTestObjectFailsThroughInstanaAPI(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AlertingChannel](t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI) {
		data := r.createTestAlertingChannelEmailData()
		resourceData := r.createAlertingChannelResourceData(data, t)
		expectedError := errors.New("test")
		mockTestObjectApi := mocks.NewMockRestResource[*restapi.AlertingChannel](ctrl)

		mockInstanaAPI.EXPECT().AlertingChannels().Return(mockTestObjectApi).Times(1)
		mockTestObjectApi.EXPECT().Create(gomock.AssignableToTypeOf(&restapi.AlertingChannel{})).Return(&restapi.AlertingChannel{}, expectedError).Times(1)

		resourceHandle := NewAlertingChannelResourceHandle()
		diag := NewTerraformResource(resourceHandle).Create(context.TODO(), resourceData, providerMeta)

		assert.NotNil(t, diag)
		assert.True(t, diag.HasError())
		assert.Equal(t, diag[0].Summary, expectedError.Error())
	})
}

func (r *terraformProviderInstanaResourceUnitTest) shouldUpdateTestObjectThroughInstanaAPI(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AlertingChannel](t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI) {
		data := r.createTestAlertingChannelEmailData()
		resourceData := r.createAlertingChannelResourceData(data, t)
		expectedModel := r.createTestAlertingChannelEmailObject()
		mockTestObjectApi := mocks.NewMockRestResource[*restapi.AlertingChannel](ctrl)

		mockInstanaAPI.EXPECT().AlertingChannels().Return(mockTestObjectApi).Times(1)
		mockTestObjectApi.EXPECT().Update(gomock.AssignableToTypeOf(&restapi.AlertingChannel{})).Return(expectedModel, nil).Times(1)

		resourceHandle := NewAlertingChannelResourceHandle()
		diag := NewTerraformResource(resourceHandle).Update(context.TODO(), resourceData, providerMeta)

		assert.Nil(t, diag)
		r.verifyTestObjectModelAppliedToResource(expectedModel, resourceData, t)
	})
}

func (r *terraformProviderInstanaResourceUnitTest) shouldReturnErrorWhenUpdateTestObjectFailsThroughInstanaAPI(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AlertingChannel](t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI) {
		data := r.createTestAlertingChannelEmailData()
		resourceData := r.createAlertingChannelResourceData(data, t)
		expectedError := errors.New("test")
		mockTestObjectApi := mocks.NewMockRestResource[*restapi.AlertingChannel](ctrl)

		mockInstanaAPI.EXPECT().AlertingChannels().Return(mockTestObjectApi).Times(1)
		mockTestObjectApi.EXPECT().Update(gomock.AssignableToTypeOf(&restapi.AlertingChannel{})).Return(&restapi.AlertingChannel{}, expectedError).Times(1)

		resourceHandle := NewAlertingChannelResourceHandle()
		diag := NewTerraformResource(resourceHandle).Update(context.TODO(), resourceData, providerMeta)

		assert.NotNil(t, diag)
		assert.True(t, diag.HasError())
		assert.Equal(t, diag[0].Summary, expectedError.Error())
	})
}

func (r *terraformProviderInstanaResourceUnitTest) shouldDeleteTestObjectThroughInstanaAPI(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AlertingChannel](t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI) {
		id := "test-id"
		data := r.createTestAlertingChannelEmailData()
		resourceData := r.createAlertingChannelResourceData(data, t)
		resourceData.SetId(id)
		mockTestObjectApi := mocks.NewMockRestResource[*restapi.AlertingChannel](ctrl)

		mockInstanaAPI.EXPECT().AlertingChannels().Return(mockTestObjectApi).Times(1)
		mockTestObjectApi.EXPECT().DeleteByID(gomock.Eq(id)).Return(nil).Times(1)

		resourceHandle := NewAlertingChannelResourceHandle()
		diag := NewTerraformResource(resourceHandle).Delete(context.TODO(), resourceData, providerMeta)

		assert.Nil(t, diag)
		assert.GreaterOrEqual(t, 0, len(resourceData.Id()))
	})
}

func (r *terraformProviderInstanaResourceUnitTest) shouldReturnErrorWhenDeleteTestObjectFailsThroughInstanaAPI(t *testing.T) {
	testHelper := NewTestHelper[*restapi.AlertingChannel](t)
	testHelper.WithMocking(t, func(ctrl *gomock.Controller, providerMeta *ProviderMeta, mockInstanaAPI *mocks.MockInstanaAPI) {
		id := "test-id"
		data := r.createTestAlertingChannelEmailData()
		resourceData := r.createAlertingChannelResourceData(data, t)
		resourceData.SetId(id)
		expectedError := errors.New("test")
		mockTestObjectApi := mocks.NewMockRestResource[*restapi.AlertingChannel](ctrl)

		mockInstanaAPI.EXPECT().AlertingChannels().Return(mockTestObjectApi).Times(1)
		mockTestObjectApi.EXPECT().DeleteByID(gomock.Eq(id)).Return(expectedError).Times(1)

		resourceHandle := NewAlertingChannelResourceHandle()
		diag := NewTerraformResource(resourceHandle).Delete(context.TODO(), resourceData, providerMeta)

		assert.NotNil(t, diag)
		assert.True(t, diag.HasError())
		assert.Equal(t, diag[0].Summary, expectedError.Error())
		assert.NotEqual(t, 0, len(resourceData.Id()))
	})
}

func (r *terraformProviderInstanaResourceUnitTest) verifyTestObjectModelAppliedToResource(model *restapi.AlertingChannel, resourceData *schema.ResourceData, t *testing.T) {
	assert.Equal(t, model.ID, resourceData.Id())
	assert.Equal(t, resourceNameWithoutPrefixAndSuffix, resourceData.Get(AlertingChannelFieldName))
	assert.Equal(t, model.Name, resourceData.Get(AlertingChannelFieldName))

	emailChannelConfig := resourceData.Get(AlertingChannelFieldChannelEmail)
	assert.IsType(t, []interface{}{}, emailChannelConfig)
	assert.Len(t, emailChannelConfig.([]interface{}), 1)
	assert.IsType(t, map[string]interface{}{}, emailChannelConfig.([]interface{})[0])

	emailConfig := emailChannelConfig.([]interface{})[0].(map[string]interface{})
	emails := ReadSetParameterFromMap[string](emailConfig, AlertingChannelEmailFieldEmails)
	assert.Equal(t, len(model.Emails), len(emails))
	for _, mail := range model.Emails {
		assert.Contains(t, emails, mail)
	}
}

func (r *terraformProviderInstanaResourceUnitTest) createTestAlertingChannelEmailObject() *restapi.AlertingChannel {
	return &restapi.AlertingChannel{
		ID:     "id",
		Name:   resourceName,
		Kind:   restapi.EmailChannelType,
		Emails: []string{"Email1", "Email2"},
	}
}

func (r *terraformProviderInstanaResourceUnitTest) createTestAlertingChannelEmailData() map[string]interface{} {
	emails := make([]interface{}, 2)
	emails[0] = "Email1"
	emails[1] = "Email2"

	data := make(map[string]interface{})
	data[AlertingChannelFieldName] = resourceNameWithoutPrefixAndSuffix
	data[AlertingChannelFieldChannelEmail] = []interface{}{
		map[string]interface{}{
			AlertingChannelEmailFieldEmails: emails,
		},
	}
	return data
}

func (r *terraformProviderInstanaResourceUnitTest) createEmptyAlertingChannelResourceData(t *testing.T) *schema.ResourceData {
	data := make(map[string]interface{})
	return r.createAlertingChannelResourceData(data, t)
}

func (r *terraformProviderInstanaResourceUnitTest) createAlertingChannelResourceData(data map[string]interface{}, t *testing.T) *schema.ResourceData {
	schemaMap := NewAlertingChannelResourceHandle().MetaData().Schema
	return schema.TestResourceDataRaw(t, schemaMap, data)
}
