package restapi_test

import (
	"errors"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const websiteMonitoringConfigID = "website-config-id"
const websiteMonitoringConfigName = "website-config-name"
const websiteMonitoringConfigAppName = "website-config-app-name"

var websiteMonitoringConfigSerialized = []byte("serialized")
var nameQueryParameter = map[string]string{"name": websiteMonitoringConfigName}

func makeTestWebsiteMonitoringConfig() *WebsiteMonitoringConfig {
	return &WebsiteMonitoringConfig{
		ID:      websiteMonitoringConfigID,
		Name:    websiteMonitoringConfigName,
		AppName: websiteMonitoringConfigAppName,
	}
}

// ########################################################
// GET All Tests
// ########################################################
func TestShouldSuccessfullyGetAllWebsiteMonitoringConfigs(t *testing.T) {
	websiteMonitoringConfig := makeTestWebsiteMonitoringConfig()
	expectedResult := []*WebsiteMonitoringConfig{websiteMonitoringConfig, websiteMonitoringConfig, websiteMonitoringConfig}
	restResponseData := []byte("server-response")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	restClient := mocks.NewMockRestClient(ctrl)
	restClient.EXPECT().Get(WebsiteMonitoringConfigResourcePath).Times(1).Return(restResponseData, nil)

	jsonUnmarshaller := mocks.NewMockJSONUnmarshaller[*WebsiteMonitoringConfig](ctrl)
	jsonUnmarshaller.EXPECT().UnmarshalArray(restResponseData).Times(1).Return(&expectedResult, nil)

	sut := NewWebsiteMonitoringConfigRestResource(jsonUnmarshaller, restClient)

	result, err := sut.GetAll()

	require.NoError(t, err)
	require.Equal(t, &expectedResult, result)
}

func TestShouldReturnEmptySliceWhenNoWebsiteMonitoringConfigsIsReturnedForGetAll(t *testing.T) {
	restResponseData := []byte("[]")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	restClient := mocks.NewMockRestClient(ctrl)
	restClient.EXPECT().Get(WebsiteMonitoringConfigResourcePath).Times(1).Return(restResponseData, nil)

	jsonUnmarshaller := mocks.NewMockJSONUnmarshaller[*WebsiteMonitoringConfig](ctrl)
	jsonUnmarshaller.EXPECT().UnmarshalArray(restResponseData).Times(1).Return(&[]*WebsiteMonitoringConfig{}, nil)

	sut := NewWebsiteMonitoringConfigRestResource(jsonUnmarshaller, restClient)

	result, err := sut.GetAll()

	require.NoError(t, err)
	require.Equal(t, &[]*WebsiteMonitoringConfig{}, result)
}

func TestShouldFailToGetAllWebsiteMonitoringConfigsWhenClientReturnsError(t *testing.T) {
	expectedError := errors.New("test")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	restClient := mocks.NewMockRestClient(ctrl)
	restClient.EXPECT().Get(WebsiteMonitoringConfigResourcePath).Times(1).Return(nil, expectedError)

	jsonUnmarshaller := mocks.NewMockJSONUnmarshaller[*WebsiteMonitoringConfig](ctrl)
	jsonUnmarshaller.EXPECT().UnmarshalArray(gomock.Any()).Times(0)

	sut := NewWebsiteMonitoringConfigRestResource(jsonUnmarshaller, restClient)

	_, err := sut.GetAll()

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

func TestShouldFailToGetAllWebsiteMonitoringConfigsWhenRestResultCannotBeUnmarshalled(t *testing.T) {
	restResponseData := []byte("invalidResponse")
	expectedError := errors.New("test")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	restClient := mocks.NewMockRestClient(ctrl)
	restClient.EXPECT().Get(WebsiteMonitoringConfigResourcePath).Times(1).Return(restResponseData, nil)

	jsonUnmarshaller := mocks.NewMockJSONUnmarshaller[*WebsiteMonitoringConfig](ctrl)
	jsonUnmarshaller.EXPECT().UnmarshalArray(restResponseData).Times(1).Return(nil, expectedError)

	sut := NewWebsiteMonitoringConfigRestResource(jsonUnmarshaller, restClient)

	_, err := sut.GetAll()

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

// ########################################################
// GET Operation Tests
// ########################################################

func TestShouldSuccessfullyExecuteGetOperationOfWebsiteMonitoringConfigRestResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller[*WebsiteMonitoringConfig](ctrl)
	websiteMonitoringConfig := makeTestWebsiteMonitoringConfig()

	client.EXPECT().GetOne(websiteMonitoringConfigID, WebsiteMonitoringConfigResourcePath).Times(1).Return(websiteMonitoringConfigSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(websiteMonitoringConfigSerialized).Times(1).Return(websiteMonitoringConfig, nil)

	sut := NewWebsiteMonitoringConfigRestResource(unmarshaller, client)

	result, err := sut.GetOne(websiteMonitoringConfigID)

	require.NoError(t, err)
	require.Equal(t, websiteMonitoringConfig, result)
}

func TestShouldReturnErrorWhenExecutingGetOperationOfWebsiteMonitoringConfigRestResourceAndGetOperationFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller[*WebsiteMonitoringConfig](ctrl)
	expectedError := errors.New("error")

	client.EXPECT().GetOne(websiteMonitoringConfigID, WebsiteMonitoringConfigResourcePath).Times(1).Return(websiteMonitoringConfigSerialized, expectedError)
	unmarshaller.EXPECT().Unmarshal(gomock.Any()).Times(0)

	sut := NewWebsiteMonitoringConfigRestResource(unmarshaller, client)

	_, err := sut.GetOne(websiteMonitoringConfigID)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

func TestShouldReturnErrorWhenExecutingGetOperationOfWebsiteMonitoringConfigRestResourceAndUnmarshallingFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller[*WebsiteMonitoringConfig](ctrl)
	expectedError := errors.New("error")

	client.EXPECT().GetOne(websiteMonitoringConfigID, WebsiteMonitoringConfigResourcePath).Times(1).Return(websiteMonitoringConfigSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(websiteMonitoringConfigSerialized).Times(1).Return(&WebsiteMonitoringConfig{}, expectedError)

	sut := NewWebsiteMonitoringConfigRestResource(unmarshaller, client)

	_, err := sut.GetOne(websiteMonitoringConfigID)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

// ########################################################
// Create Operation Tests
// ########################################################

func TestShouldSuccessfullyExecuteCreateOperationOfWebsiteMonitoringConfigRestResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller[*WebsiteMonitoringConfig](ctrl)
	websiteMonitoringConfig := makeTestWebsiteMonitoringConfig()

	client.EXPECT().PostByQuery(WebsiteMonitoringConfigResourcePath, nameQueryParameter).Times(1).Return(websiteMonitoringConfigSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(websiteMonitoringConfigSerialized).Times(1).Return(websiteMonitoringConfig, nil)

	sut := NewWebsiteMonitoringConfigRestResource(unmarshaller, client)

	result, err := sut.Create(websiteMonitoringConfig)

	require.NoError(t, err)
	require.Equal(t, websiteMonitoringConfig, result)
}

func TestShouldReturnErrorWhenExecutingCreateOperationOfWebsiteMonitoringConfigRestResourceAndPostOperationFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller[*WebsiteMonitoringConfig](ctrl)
	expectedError := errors.New("error")
	websiteMonitoringConfig := makeTestWebsiteMonitoringConfig()

	client.EXPECT().PostByQuery(WebsiteMonitoringConfigResourcePath, nameQueryParameter).Times(1).Return(websiteMonitoringConfigSerialized, expectedError)
	unmarshaller.EXPECT().Unmarshal(gomock.Any()).Times(0)

	sut := NewWebsiteMonitoringConfigRestResource(unmarshaller, client)

	_, err := sut.Create(websiteMonitoringConfig)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

func TestShouldReturnErrorWhenExecutingCreateOperationOfWebsiteMonitoringConfigRestResourceAndUnmarshallingFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller[*WebsiteMonitoringConfig](ctrl)
	expectedError := errors.New("error")
	websiteMonitoringConfig := makeTestWebsiteMonitoringConfig()

	client.EXPECT().PostByQuery(WebsiteMonitoringConfigResourcePath, nameQueryParameter).Times(1).Return(websiteMonitoringConfigSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(websiteMonitoringConfigSerialized).Times(1).Return(&WebsiteMonitoringConfig{}, expectedError)

	sut := NewWebsiteMonitoringConfigRestResource(unmarshaller, client)

	_, err := sut.Create(websiteMonitoringConfig)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

// ########################################################
// Update Operation Tests
// ########################################################

func TestShouldSuccessfullyExecuteUpdateOperationOfWebsiteMonitoringConfigRestResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller[*WebsiteMonitoringConfig](ctrl)
	websiteMonitoringConfig := makeTestWebsiteMonitoringConfig()

	client.EXPECT().PutByQuery(WebsiteMonitoringConfigResourcePath, websiteMonitoringConfigID, nameQueryParameter).Times(1).Return(websiteMonitoringConfigSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(websiteMonitoringConfigSerialized).Times(1).Return(websiteMonitoringConfig, nil)

	sut := NewWebsiteMonitoringConfigRestResource(unmarshaller, client)

	result, err := sut.Update(websiteMonitoringConfig)

	require.NoError(t, err)
	require.Equal(t, websiteMonitoringConfig, result)
}

func TestShouldReturnErrorWhenExecutingUpdateOperationOfWebsiteMonitoringConfigRestResourceAndPutOperationFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller[*WebsiteMonitoringConfig](ctrl)
	expectedError := errors.New("error")
	websiteMonitoringConfig := makeTestWebsiteMonitoringConfig()

	client.EXPECT().PutByQuery(WebsiteMonitoringConfigResourcePath, websiteMonitoringConfigID, nameQueryParameter).Times(1).Return(websiteMonitoringConfigSerialized, expectedError)
	unmarshaller.EXPECT().Unmarshal(gomock.Any()).Times(0)

	sut := NewWebsiteMonitoringConfigRestResource(unmarshaller, client)

	_, err := sut.Update(websiteMonitoringConfig)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

func TestShouldReturnErrorWhenExecutingUpdateOperationOfWebsiteMonitoringConfigRestResourceAndUnmarshallingFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller[*WebsiteMonitoringConfig](ctrl)
	expectedError := errors.New("error")
	websiteMonitoringConfig := makeTestWebsiteMonitoringConfig()

	client.EXPECT().PutByQuery(WebsiteMonitoringConfigResourcePath, websiteMonitoringConfigID, nameQueryParameter).Times(1).Return(websiteMonitoringConfigSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(websiteMonitoringConfigSerialized).Times(1).Return(&WebsiteMonitoringConfig{}, expectedError)

	sut := NewWebsiteMonitoringConfigRestResource(unmarshaller, client)

	_, err := sut.Update(websiteMonitoringConfig)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

// ########################################################
// Delete Operation Tests
// ########################################################

func TestShouldSuccessfullyExecuteDeleteByObjectOperationOfWebsiteMonitoringConfigRestResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller[*WebsiteMonitoringConfig](ctrl)
	websiteMonitoringConfig := makeTestWebsiteMonitoringConfig()

	client.EXPECT().Delete(websiteMonitoringConfigID, WebsiteMonitoringConfigResourcePath).Times(1).Return(nil)

	sut := NewWebsiteMonitoringConfigRestResource(unmarshaller, client)

	err := sut.Delete(websiteMonitoringConfig)

	require.NoError(t, err)
}

func TestShouldReturnErrorWhenExecutingDeleteByObjectOperationOfWebsiteMonitoringConfigRestResourceAndDeleteRequestFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller[*WebsiteMonitoringConfig](ctrl)
	expectedError := errors.New("error")
	websiteMonitoringConfig := makeTestWebsiteMonitoringConfig()

	client.EXPECT().Delete(websiteMonitoringConfigID, WebsiteMonitoringConfigResourcePath).Times(1).Return(expectedError)

	sut := NewWebsiteMonitoringConfigRestResource(unmarshaller, client)

	err := sut.Delete(websiteMonitoringConfig)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}
func TestShouldSuccessfullyExecuteDeleteByIdOperationOfWebsiteMonitoringConfigRestResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller[*WebsiteMonitoringConfig](ctrl)

	client.EXPECT().Delete(websiteMonitoringConfigID, WebsiteMonitoringConfigResourcePath).Times(1).Return(nil)

	sut := NewWebsiteMonitoringConfigRestResource(unmarshaller, client)

	err := sut.DeleteByID(websiteMonitoringConfigID)

	require.NoError(t, err)
}

func TestShouldReturnErrorWhenExecutingDeleteByIdOperationOfWebsiteMonitoringConfigRestResourceAndDeleteRequestFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller[*WebsiteMonitoringConfig](ctrl)
	expectedError := errors.New("error")

	client.EXPECT().Delete(websiteMonitoringConfigID, WebsiteMonitoringConfigResourcePath).Times(1).Return(expectedError)

	sut := NewWebsiteMonitoringConfigRestResource(unmarshaller, client)

	err := sut.DeleteByID(websiteMonitoringConfigID)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}
