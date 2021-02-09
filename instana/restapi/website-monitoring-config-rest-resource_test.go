package restapi_test

import (
	"errors"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	mocks "github.com/gessnerfl/terraform-provider-instana/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const nameIsMissingError = "Name is missing"

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
// GET Operation Tests
// ########################################################

func TestShouldSuccessfullyExecuteGetOperationOfWebsiteMonitoringConfigRestResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
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
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	expectedError := errors.New("Error")

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
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	expectedError := errors.New("Error")

	client.EXPECT().GetOne(websiteMonitoringConfigID, WebsiteMonitoringConfigResourcePath).Times(1).Return(websiteMonitoringConfigSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(websiteMonitoringConfigSerialized).Times(1).Return(WebsiteMonitoringConfig{}, expectedError)

	sut := NewWebsiteMonitoringConfigRestResource(unmarshaller, client)

	_, err := sut.GetOne(websiteMonitoringConfigID)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

type InvalidWebsiteMonitoringConfig struct{}

func TestShouldReturnErrorWhenExecutingGetOperationOfWebsiteMonitoringConfigRestResourceAndUnmarshallingDoesNotProvideAInstanaDataObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)

	client.EXPECT().GetOne(websiteMonitoringConfigID, WebsiteMonitoringConfigResourcePath).Times(1).Return(websiteMonitoringConfigSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(websiteMonitoringConfigSerialized).Times(1).Return(&InvalidWebsiteMonitoringConfig{}, nil)

	sut := NewWebsiteMonitoringConfigRestResource(unmarshaller, client)

	_, err := sut.GetOne(websiteMonitoringConfigID)

	require.Error(t, err)
	require.Contains(t, err.Error(), "Unmarshalled object does not implement InstanaDataObject")
}

func TestShouldReturnErrorWhenExecutingGetOperationOfWebsiteMonitoringConfigRestResourceAndReceivedObjectIsNotValid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)

	client.EXPECT().GetOne(websiteMonitoringConfigID, WebsiteMonitoringConfigResourcePath).Times(1).Return(websiteMonitoringConfigSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(websiteMonitoringConfigSerialized).Times(1).Return(&WebsiteMonitoringConfig{}, nil)

	sut := NewWebsiteMonitoringConfigRestResource(unmarshaller, client)

	_, err := sut.GetOne(websiteMonitoringConfigID)

	require.Error(t, err)
	require.Contains(t, nameIsMissingError, err.Error())
}

// ########################################################
// Create Operation Tests
// ########################################################

func TestShouldSuccessfullyExecuteCreateOperationOfWebsiteMonitoringConfigRestResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	websiteMonitoringConfig := makeTestWebsiteMonitoringConfig()

	client.EXPECT().PostByQuery(WebsiteMonitoringConfigResourcePath, nameQueryParameter).Times(1).Return(websiteMonitoringConfigSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(websiteMonitoringConfigSerialized).Times(1).Return(websiteMonitoringConfig, nil)

	sut := NewWebsiteMonitoringConfigRestResource(unmarshaller, client)

	result, err := sut.Create(websiteMonitoringConfig)

	require.NoError(t, err)
	require.Equal(t, websiteMonitoringConfig, result)
}

func TestShouldReturnErrorWhenExecutingCreateOperationOfWebsiteMonitoringConfigRestResourceAndProvidedObjectIsNotValid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	websiteMonitoringConfig := &WebsiteMonitoringConfig{}

	client.EXPECT().PostByQuery(WebsiteMonitoringConfigResourcePath, nameQueryParameter).Times(0)
	unmarshaller.EXPECT().Unmarshal(websiteMonitoringConfigSerialized).Times(0)

	sut := NewWebsiteMonitoringConfigRestResource(unmarshaller, client)

	_, err := sut.Create(websiteMonitoringConfig)

	require.Error(t, err)
	require.Contains(t, nameIsMissingError, err.Error())
}

func TestShouldReturnErrorWhenExecutingCreateOperationOfWebsiteMonitoringConfigRestResourceAndPostOperationFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	expectedError := errors.New("Error")
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
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	expectedError := errors.New("Error")
	websiteMonitoringConfig := makeTestWebsiteMonitoringConfig()

	client.EXPECT().PostByQuery(WebsiteMonitoringConfigResourcePath, nameQueryParameter).Times(1).Return(websiteMonitoringConfigSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(websiteMonitoringConfigSerialized).Times(1).Return(&WebsiteMonitoringConfig{}, expectedError)

	sut := NewWebsiteMonitoringConfigRestResource(unmarshaller, client)

	_, err := sut.Create(websiteMonitoringConfig)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

func TestShouldReturnErrorWhenExecutingCreateOperationOfWebsiteMonitoringConfigRestResourceAndReceivedObjectIsNotValid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	websiteMonitoringConfig := makeTestWebsiteMonitoringConfig()

	client.EXPECT().PostByQuery(WebsiteMonitoringConfigResourcePath, nameQueryParameter).Times(1).Return(websiteMonitoringConfigSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(websiteMonitoringConfigSerialized).Times(1).Return(&WebsiteMonitoringConfig{}, nil)

	sut := NewWebsiteMonitoringConfigRestResource(unmarshaller, client)

	_, err := sut.Create(websiteMonitoringConfig)

	require.Error(t, err)
	require.Contains(t, nameIsMissingError, err.Error())
}

// ########################################################
// Update Operation Tests
// ########################################################

func TestShouldSuccessfullyExecuteUpdateOperationOfWebsiteMonitoringConfigRestResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	websiteMonitoringConfig := makeTestWebsiteMonitoringConfig()

	client.EXPECT().PutByQuery(WebsiteMonitoringConfigResourcePath, websiteMonitoringConfigID, nameQueryParameter).Times(1).Return(websiteMonitoringConfigSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(websiteMonitoringConfigSerialized).Times(1).Return(websiteMonitoringConfig, nil)

	sut := NewWebsiteMonitoringConfigRestResource(unmarshaller, client)

	result, err := sut.Update(websiteMonitoringConfig)

	require.NoError(t, err)
	require.Equal(t, websiteMonitoringConfig, result)
}

func TestShouldReturnErrorWhenExecutingUpdateOperationOfWebsiteMonitoringConfigRestResourceAndProvidedObjectIsNotValid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	websiteMonitoringConfig := &WebsiteMonitoringConfig{}

	client.EXPECT().PutByQuery(WebsiteMonitoringConfigResourcePath, websiteMonitoringConfigID, nameQueryParameter).Times(0)
	unmarshaller.EXPECT().Unmarshal(websiteMonitoringConfigSerialized).Times(0)

	sut := NewWebsiteMonitoringConfigRestResource(unmarshaller, client)

	_, err := sut.Update(websiteMonitoringConfig)

	require.Error(t, err)
	require.Contains(t, nameIsMissingError, err.Error())
}

func TestShouldReturnErrorWhenExecutingUpdateOperationOfWebsiteMonitoringConfigRestResourceAndPutOperationFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	expectedError := errors.New("Error")
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
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	expectedError := errors.New("Error")
	websiteMonitoringConfig := makeTestWebsiteMonitoringConfig()

	client.EXPECT().PutByQuery(WebsiteMonitoringConfigResourcePath, websiteMonitoringConfigID, nameQueryParameter).Times(1).Return(websiteMonitoringConfigSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(websiteMonitoringConfigSerialized).Times(1).Return(&WebsiteMonitoringConfig{}, expectedError)

	sut := NewWebsiteMonitoringConfigRestResource(unmarshaller, client)

	_, err := sut.Update(websiteMonitoringConfig)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

func TestShouldReturnErrorWhenExecutingUpdateOperationOfWebsiteMonitoringConfigRestResourceAndReceivedObjectIsNotValid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	websiteMonitoringConfig := makeTestWebsiteMonitoringConfig()

	client.EXPECT().PutByQuery(WebsiteMonitoringConfigResourcePath, websiteMonitoringConfigID, nameQueryParameter).Times(1).Return(websiteMonitoringConfigSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(websiteMonitoringConfigSerialized).Times(1).Return(&WebsiteMonitoringConfig{}, nil)

	sut := NewWebsiteMonitoringConfigRestResource(unmarshaller, client)

	_, err := sut.Update(websiteMonitoringConfig)

	require.Error(t, err)
	require.Contains(t, nameIsMissingError, err.Error())
}

// ########################################################
// Delete Operation Tests
// ########################################################

func TestShouldSuccessfullyExecuteDeleteByObjectOperationOfWebsiteMonitoringConfigRestResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
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
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	expectedError := errors.New("Error")
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
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)

	client.EXPECT().Delete(websiteMonitoringConfigID, WebsiteMonitoringConfigResourcePath).Times(1).Return(nil)

	sut := NewWebsiteMonitoringConfigRestResource(unmarshaller, client)

	err := sut.DeleteByID(websiteMonitoringConfigID)

	require.NoError(t, err)
}

func TestShouldReturnErrorWhenExecutingDeleteByIdOperationOfWebsiteMonitoringConfigRestResourceAndDeleteRequestFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	expectedError := errors.New("Error")

	client.EXPECT().Delete(websiteMonitoringConfigID, WebsiteMonitoringConfigResourcePath).Times(1).Return(expectedError)

	sut := NewWebsiteMonitoringConfigRestResource(unmarshaller, client)

	err := sut.DeleteByID(websiteMonitoringConfigID)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}
