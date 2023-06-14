package restapi_test

import (
	"errors"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	mocks "github.com/gessnerfl/terraform-provider-instana/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const (
	syntheticMonitorID                  = "id"
	syntheticMonitorLabel               = "label"
	syntheticMonitorAppName             = "appName"
	syntheticMonitorActive              = true
	syntheticMonitorConfigSyntheticType = "HTTPAction"
	syntheticMonitorLocation            = "location"
	syntheticMonitorConfigUrl           = "url"
	syntheticMonitorConfigOperation     = "operation"
)

var syntheticMonitorSerialized = []byte("serialized")
var IDQueryParameter = map[string]string{"id": syntheticMonitorID}

func makeSyntheticMonitor() *SyntheticMonitor {
	return &SyntheticMonitor{
		ID:        syntheticMonitorID,
		Label:     syntheticMonitorLabel,
		Active:    syntheticMonitorActive,
		Locations: []string{syntheticMonitorLocation},
		Configuration: SyntheticTestConfig{
			SyntheticType: syntheticMonitorConfigSyntheticType,
			URL:           syntheticMonitorConfigUrl,
			Operation:     syntheticMonitorConfigOperation,
		},
	}
}

// ########################################################
// GET Operation Tests
// ########################################################

func TestShouldSuccessfullyExecuteGetOperationOfSyntheticMonitorRestResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	syntheticMonitor := makeSyntheticMonitor()

	client.EXPECT().GetOne(syntheticMonitorID, SyntheticMonitorResourcePath).Times(1).Return(syntheticMonitorSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(syntheticMonitorSerialized).Times(1).Return(syntheticMonitor, nil)

	sut := NewSyntheticMonitorRestResource(unmarshaller, client)

	result, err := sut.GetOne(syntheticMonitorID)

	require.NoError(t, err)
	require.Equal(t, syntheticMonitor, result)
}

func TestShouldReturnErrorWhenExecutingGetOperationOfSyntheticMonitorRestResourceAndGetOperationFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	expectedError := errors.New("Error")

	client.EXPECT().GetOne(syntheticMonitorID, SyntheticMonitorResourcePath).Times(1).Return(syntheticMonitorSerialized, expectedError)
	unmarshaller.EXPECT().Unmarshal(gomock.Any()).Times(0)

	sut := NewSyntheticMonitorRestResource(unmarshaller, client)

	_, err := sut.GetOne(syntheticMonitorID)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

func TestShouldReturnErrorWhenExecutingGetOperationOfSyntheticMonitorRestResourceAndUnmarshallingFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	expectedError := errors.New("Error")

	client.EXPECT().GetOne(syntheticMonitorID, SyntheticMonitorResourcePath).Times(1).Return(syntheticMonitorSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(syntheticMonitorSerialized).Times(1).Return(SyntheticMonitor{}, expectedError)

	sut := NewSyntheticMonitorRestResource(unmarshaller, client)

	_, err := sut.GetOne(syntheticMonitorID)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

type InvalidSyntheticMonitor struct{}

func TestShouldReturnErrorWhenExecutingGetOperationOfSyntheticMonitorConfigRestResourceAndUnmarshallingDoesNotProvideAInstanaDataObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)

	client.EXPECT().GetOne(syntheticMonitorID, SyntheticMonitorResourcePath).Times(1).Return(syntheticMonitorSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(syntheticMonitorSerialized).Times(1).Return(&InvalidSyntheticMonitor{}, nil)

	sut := NewSyntheticMonitorRestResource(unmarshaller, client)

	_, err := sut.GetOne(syntheticMonitorID)

	require.Error(t, err)
	require.Contains(t, err.Error(), "unmarshalled object does not implement InstanaDataObject")
}

func TestShouldReturnErrorWhenExecutingGetOperationOfSyntheticMonitorRestResourceAndReceivedObjectIsNotValid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)

	client.EXPECT().GetOne(syntheticMonitorID, SyntheticMonitorResourcePath).Times(1).Return(syntheticMonitorSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(syntheticMonitorSerialized).Times(1).Return(&SyntheticMonitor{}, nil)

	sut := NewSyntheticMonitorRestResource(unmarshaller, client)

	_, err := sut.GetOne(syntheticMonitorID)

	require.Error(t, err)
	require.Contains(t, "id is missing", err.Error())
}

// ########################################################
// Create Operation Tests
// ########################################################

func TestShouldSuccessfullyExecuteCreateOperationOfSyntheticMonitorRestResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	syntheticMonitor := makeSyntheticMonitor()

	client.EXPECT().Post(gomock.Eq(syntheticMonitor), gomock.Eq(SyntheticMonitorResourcePath)).Times(1).Return(syntheticMonitorSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(syntheticMonitorSerialized).Times(1).Return(syntheticMonitor, nil)

	sut := NewSyntheticMonitorRestResource(unmarshaller, client)

	result, err := sut.Create(syntheticMonitor)

	require.NoError(t, err)
	require.Equal(t, syntheticMonitor, result)
}

func TestShouldReturnErrorWhenExecutingCreateOperationOfSyntheticMonitorRestResourceAndProvidedObjectIsNotValid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	syntheticMonitor := &SyntheticMonitor{}

	client.EXPECT().Post(gomock.Eq(syntheticMonitor), gomock.Eq(SyntheticMonitorResourcePath)).Times(0)
	unmarshaller.EXPECT().Unmarshal(syntheticMonitorSerialized).Times(0)

	sut := NewSyntheticMonitorRestResource(unmarshaller, client)

	_, err := sut.Create(syntheticMonitor)

	require.Error(t, err)
	require.Contains(t, "id is missing", err.Error())
}

func TestShouldReturnErrorWhenExecutingCreateOperationOfSyntheticMonitorRestResourceAndPostOperationFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	expectedError := errors.New("Error")
	syntheticMonitor := makeSyntheticMonitor()

	client.EXPECT().Post(gomock.Eq(syntheticMonitor), gomock.Eq(SyntheticMonitorResourcePath)).Times(1).Return(syntheticMonitorSerialized, expectedError)
	unmarshaller.EXPECT().Unmarshal(gomock.Any()).Times(0)

	sut := NewSyntheticMonitorRestResource(unmarshaller, client)

	_, err := sut.Create(syntheticMonitor)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

func TestShouldReturnErrorWhenExecutingCreateOperationOfSyntheticMonitorRestResourceAndUnmarshallingFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	expectedError := errors.New("Error")
	syntheticMonitor := makeSyntheticMonitor()

	client.EXPECT().Post(gomock.Eq(syntheticMonitor), gomock.Eq(SyntheticMonitorResourcePath)).Times(1).Return(syntheticMonitorSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(syntheticMonitorSerialized).Times(1).Return(&SyntheticMonitor{}, expectedError)

	sut := NewSyntheticMonitorRestResource(unmarshaller, client)

	_, err := sut.Create(syntheticMonitor)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

func TestShouldReturnErrorWhenExecutingCreateOperationOfSyntheticMonitorRestResourceAndReceivedObjectIsNotValid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	syntheticMonitor := makeSyntheticMonitor()

	client.EXPECT().Post(gomock.Eq(syntheticMonitor), gomock.Eq(SyntheticMonitorResourcePath)).Times(1).Return(syntheticMonitorSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(syntheticMonitorSerialized).Times(1).Return(&SyntheticMonitor{}, nil)

	sut := NewSyntheticMonitorRestResource(unmarshaller, client)

	_, err := sut.Create(syntheticMonitor)

	require.Error(t, err)
	require.Contains(t, "id is missing", err.Error())
}

// ########################################################
// Update Operation Tests
// ########################################################

func TestShouldSuccessfullyExecuteUpdateOperationOfSyntheticMonitorRestResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	syntheticMonitor := makeSyntheticMonitor()

	client.EXPECT().Put(gomock.Eq(syntheticMonitor), gomock.Eq(SyntheticMonitorResourcePath)).Times(1)
	client.EXPECT().GetOne(syntheticMonitorID, SyntheticMonitorResourcePath).Times(1).Return(syntheticMonitorSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(syntheticMonitorSerialized).Times(1).Return(syntheticMonitor, nil)

	sut := NewSyntheticMonitorRestResource(unmarshaller, client)

	result, err := sut.Update(syntheticMonitor)

	require.NoError(t, err)
	require.Equal(t, syntheticMonitor, result)
}

func TestShouldReturnErrorWhenExecutingUpdateOperationOfSyntheticMonitorRestResourceAndProvidedObjectIsNotValid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	syntheticMonitor := &SyntheticMonitor{}

	client.EXPECT().Put(gomock.Eq(syntheticMonitor), gomock.Eq(SyntheticMonitorResourcePath)).Times(0)
	unmarshaller.EXPECT().Unmarshal(syntheticMonitorSerialized).Times(0)

	sut := NewSyntheticMonitorRestResource(unmarshaller, client)

	_, err := sut.Update(syntheticMonitor)

	require.Error(t, err)
	require.Contains(t, "id is missing", err.Error())
}

func TestShouldReturnErrorWhenExecutingUpdateOperationOfSyntheticMonitorRestResourceAndPutOperationFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	expectedError := errors.New("Error")
	syntheticMonitor := makeSyntheticMonitor()

	client.EXPECT().Put(gomock.Eq(syntheticMonitor), gomock.Eq(SyntheticMonitorResourcePath)).Times(1).Return(syntheticMonitorSerialized, expectedError)
	unmarshaller.EXPECT().Unmarshal(gomock.Any()).Times(0)

	sut := NewSyntheticMonitorRestResource(unmarshaller, client)

	_, err := sut.Update(syntheticMonitor)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

func TestShouldReturnErrorWhenExecutingUpdateOperationOfSyntheticMonitorRestResourceAndUnmarshallingFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	expectedError := errors.New("Error")
	syntheticMonitor := makeSyntheticMonitor()

	client.EXPECT().Put(gomock.Eq(syntheticMonitor), gomock.Eq(SyntheticMonitorResourcePath)).Times(1)
	client.EXPECT().GetOne(syntheticMonitorID, SyntheticMonitorResourcePath).Times(1).Return(syntheticMonitorSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(syntheticMonitorSerialized).Times(1).Return(&SyntheticMonitor{}, expectedError)

	sut := NewSyntheticMonitorRestResource(unmarshaller, client)

	_, err := sut.Update(syntheticMonitor)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

func TestShouldReturnErrorWhenExecutingUpdateOperationOfSyntheticMonitorRestResourceAndReceivedObjectIsNotValid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	syntheticMonitor := makeSyntheticMonitor()

	client.EXPECT().Put(gomock.Eq(syntheticMonitor), gomock.Eq(SyntheticMonitorResourcePath)).Times(1)
	client.EXPECT().GetOne(syntheticMonitorID, SyntheticMonitorResourcePath).Times(1).Return(syntheticMonitorSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(syntheticMonitorSerialized).Times(1).Return(&SyntheticMonitor{}, nil)

	sut := NewSyntheticMonitorRestResource(unmarshaller, client)

	_, err := sut.Update(syntheticMonitor)

	require.Error(t, err)
	require.Contains(t, "id is missing", err.Error())
}

// ########################################################
// Delete Operation Tests
// ########################################################

func TestShouldSuccessfullyExecuteDeleteByObjectOperationOfSyntheticMonitorRestResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	syntheticMonitor := makeSyntheticMonitor()

	client.EXPECT().Delete(syntheticMonitorID, SyntheticMonitorResourcePath).Times(1).Return(nil)

	sut := NewSyntheticMonitorRestResource(unmarshaller, client)

	err := sut.Delete(syntheticMonitor)

	require.NoError(t, err)
}

func TestShouldReturnErrorWhenExecutingDeleteByObjectOperationOfSyntheticMonitorRestResourceAndDeleteRequestFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	expectedError := errors.New("Error")
	syntheticMonitor := makeSyntheticMonitor()

	client.EXPECT().Delete(syntheticMonitorID, SyntheticMonitorResourcePath).Times(1).Return(expectedError)

	sut := NewSyntheticMonitorRestResource(unmarshaller, client)

	err := sut.Delete(syntheticMonitor)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}
func TestShouldSuccessfullyExecuteDeleteByIdOperationOfSyntheticMonitorRestResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)

	client.EXPECT().Delete(syntheticMonitorID, SyntheticMonitorResourcePath).Times(1).Return(nil)

	sut := NewSyntheticMonitorRestResource(unmarshaller, client)

	err := sut.DeleteByID(syntheticMonitorID)

	require.NoError(t, err)
}

func TestShouldReturnErrorWhenExecutingDeleteByIdOperationOfSyntheticMonitorRestResourceAndDeleteRequestFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	expectedError := errors.New("Error")

	client.EXPECT().Delete(syntheticMonitorID, SyntheticMonitorResourcePath).Times(1).Return(expectedError)

	sut := NewSyntheticMonitorRestResource(unmarshaller, client)

	err := sut.DeleteByID(syntheticMonitorID)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}
