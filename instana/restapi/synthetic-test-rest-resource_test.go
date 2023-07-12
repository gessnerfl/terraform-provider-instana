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
	syntheticTestID                  = "id"
	syntheticTestLabel               = "label"
	syntheticTestActive              = true
	syntheticTestConfigSyntheticType = "HTTPAction"
	syntheticTestLocation            = "location"
	syntheticTestConfigUrl           = "url"
	syntheticTestConfigOperation     = "operation"
)

var syntheticTestSerialized = []byte("serialized")
var IDQueryParameter = map[string]string{"id": syntheticTestID}

func makeSyntheticTest() *SyntheticTest {
	return &SyntheticTest{
		ID:        syntheticTestID,
		Label:     syntheticTestLabel,
		Active:    syntheticTestActive,
		Locations: []string{syntheticTestLocation},
		Configuration: SyntheticTestConfig{
			SyntheticType: syntheticTestConfigSyntheticType,
			URL:           syntheticTestConfigUrl,
			Operation:     syntheticTestConfigOperation,
		},
	}
}

// ########################################################
// GET Operation Tests
// ########################################################

func TestShouldSuccessfullyExecuteGetOperationOfSyntheticTestRestResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	syntheticTest := makeSyntheticTest()

	client.EXPECT().GetOne(syntheticTestID, SyntheticTestResourcePath).Times(1).Return(syntheticTestSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(syntheticTestSerialized).Times(1).Return(syntheticTest, nil)

	sut := NewSyntheticTestRestResource(unmarshaller, client)

	result, err := sut.GetOne(syntheticTestID)

	require.NoError(t, err)
	require.Equal(t, syntheticTest, result)
}

func TestShouldReturnErrorWhenExecutingGetOperationOfSyntheticTestRestResourceAndGetOperationFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	expectedError := errors.New("Error")

	client.EXPECT().GetOne(syntheticTestID, SyntheticTestResourcePath).Times(1).Return(syntheticTestSerialized, expectedError)
	unmarshaller.EXPECT().Unmarshal(gomock.Any()).Times(0)

	sut := NewSyntheticTestRestResource(unmarshaller, client)

	_, err := sut.GetOne(syntheticTestID)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

func TestShouldReturnErrorWhenExecutingGetOperationOfSyntheticTestRestResourceAndUnmarshallingFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	expectedError := errors.New("Error")

	client.EXPECT().GetOne(syntheticTestID, SyntheticTestResourcePath).Times(1).Return(syntheticTestSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(syntheticTestSerialized).Times(1).Return(SyntheticTest{}, expectedError)

	sut := NewSyntheticTestRestResource(unmarshaller, client)

	_, err := sut.GetOne(syntheticTestID)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

type InvalidSyntheticTest struct{}

func TestShouldReturnErrorWhenExecutingGetOperationOfSyntheticTestConfigRestResourceAndUnmarshallingDoesNotProvideAInstanaDataObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)

	client.EXPECT().GetOne(syntheticTestID, SyntheticTestResourcePath).Times(1).Return(syntheticTestSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(syntheticTestSerialized).Times(1).Return(&InvalidSyntheticTest{}, nil)

	sut := NewSyntheticTestRestResource(unmarshaller, client)

	_, err := sut.GetOne(syntheticTestID)

	require.Error(t, err)
	require.Contains(t, err.Error(), "unmarshalled object does not implement InstanaDataObject")
}

func TestShouldReturnErrorWhenExecutingGetOperationOfSyntheticTestRestResourceAndReceivedObjectIsNotValid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)

	client.EXPECT().GetOne(syntheticTestID, SyntheticTestResourcePath).Times(1).Return(syntheticTestSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(syntheticTestSerialized).Times(1).Return(&SyntheticTest{}, nil)

	sut := NewSyntheticTestRestResource(unmarshaller, client)

	_, err := sut.GetOne(syntheticTestID)

	require.Error(t, err)
	require.Contains(t, "id is missing", err.Error())
}

// ########################################################
// Create Operation Tests
// ########################################################

func TestShouldSuccessfullyExecuteCreateOperationOfSyntheticTestRestResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	syntheticTest := makeSyntheticTest()

	client.EXPECT().Post(gomock.Eq(syntheticTest), gomock.Eq(SyntheticTestResourcePath)).Times(1).Return(syntheticTestSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(syntheticTestSerialized).Times(1).Return(syntheticTest, nil)

	sut := NewSyntheticTestRestResource(unmarshaller, client)

	result, err := sut.Create(syntheticTest)

	require.NoError(t, err)
	require.Equal(t, syntheticTest, result)
}

func TestShouldReturnErrorWhenExecutingCreateOperationOfSyntheticTestRestResourceAndProvidedObjectIsNotValid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	SyntheticTest := &SyntheticTest{}

	client.EXPECT().Post(gomock.Eq(SyntheticTest), gomock.Eq(SyntheticTestResourcePath)).Times(0)
	unmarshaller.EXPECT().Unmarshal(syntheticTestSerialized).Times(0)

	sut := NewSyntheticTestRestResource(unmarshaller, client)

	_, err := sut.Create(SyntheticTest)

	require.Error(t, err)
	require.Contains(t, "id is missing", err.Error())
}

func TestShouldReturnErrorWhenExecutingCreateOperationOfSyntheticTestRestResourceAndPostOperationFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	expectedError := errors.New("Error")
	syntheticTest := makeSyntheticTest()

	client.EXPECT().Post(gomock.Eq(syntheticTest), gomock.Eq(SyntheticTestResourcePath)).Times(1).Return(syntheticTestSerialized, expectedError)
	unmarshaller.EXPECT().Unmarshal(gomock.Any()).Times(0)

	sut := NewSyntheticTestRestResource(unmarshaller, client)

	_, err := sut.Create(syntheticTest)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

func TestShouldReturnErrorWhenExecutingCreateOperationOfSyntheticTestRestResourceAndUnmarshallingFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	expectedError := errors.New("Error")
	syntheticTest := makeSyntheticTest()

	client.EXPECT().Post(gomock.Eq(syntheticTest), gomock.Eq(SyntheticTestResourcePath)).Times(1).Return(syntheticTestSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(syntheticTestSerialized).Times(1).Return(&SyntheticTest{}, expectedError)

	sut := NewSyntheticTestRestResource(unmarshaller, client)

	_, err := sut.Create(syntheticTest)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

func TestShouldReturnErrorWhenExecutingCreateOperationOfSyntheticTestRestResourceAndReceivedObjectIsNotValid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	syntheticTest := makeSyntheticTest()

	client.EXPECT().Post(gomock.Eq(syntheticTest), gomock.Eq(SyntheticTestResourcePath)).Times(1).Return(syntheticTestSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(syntheticTestSerialized).Times(1).Return(&SyntheticTest{}, nil)

	sut := NewSyntheticTestRestResource(unmarshaller, client)

	_, err := sut.Create(syntheticTest)

	require.Error(t, err)
	require.Contains(t, "id is missing", err.Error())
}

// ########################################################
// Update Operation Tests
// ########################################################

func TestShouldSuccessfullyExecuteUpdateOperationOfSyntheticTestRestResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	syntheticTest := makeSyntheticTest()

	client.EXPECT().Put(gomock.Eq(syntheticTest), gomock.Eq(SyntheticTestResourcePath)).Times(1)
	client.EXPECT().GetOne(syntheticTestID, SyntheticTestResourcePath).Times(1).Return(syntheticTestSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(syntheticTestSerialized).Times(1).Return(syntheticTest, nil)

	sut := NewSyntheticTestRestResource(unmarshaller, client)

	result, err := sut.Update(syntheticTest)

	require.NoError(t, err)
	require.Equal(t, syntheticTest, result)
}

func TestShouldReturnErrorWhenExecutingUpdateOperationOfSyntheticTestRestResourceAndProvidedObjectIsNotValid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	syntheticTest := &SyntheticTest{}

	client.EXPECT().Put(gomock.Eq(syntheticTest), gomock.Eq(SyntheticTestResourcePath)).Times(0)
	unmarshaller.EXPECT().Unmarshal(syntheticTestSerialized).Times(0)

	sut := NewSyntheticTestRestResource(unmarshaller, client)

	_, err := sut.Update(syntheticTest)

	require.Error(t, err)
	require.Contains(t, "id is missing", err.Error())
}

func TestShouldReturnErrorWhenExecutingUpdateOperationOfSyntheticTestRestResourceAndPutOperationFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	expectedError := errors.New("Error")
	syntheticTest := makeSyntheticTest()

	client.EXPECT().Put(gomock.Eq(syntheticTest), gomock.Eq(SyntheticTestResourcePath)).Times(1).Return(syntheticTestSerialized, expectedError)
	unmarshaller.EXPECT().Unmarshal(gomock.Any()).Times(0)

	sut := NewSyntheticTestRestResource(unmarshaller, client)

	_, err := sut.Update(syntheticTest)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

func TestShouldReturnErrorWhenExecutingUpdateOperationOfSyntheticTestRestResourceAndUnmarshallingFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	expectedError := errors.New("Error")
	syntheticTest := makeSyntheticTest()

	client.EXPECT().Put(gomock.Eq(syntheticTest), gomock.Eq(SyntheticTestResourcePath)).Times(1)
	client.EXPECT().GetOne(syntheticTestID, SyntheticTestResourcePath).Times(1).Return(syntheticTestSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(syntheticTestSerialized).Times(1).Return(&SyntheticTest{}, expectedError)

	sut := NewSyntheticTestRestResource(unmarshaller, client)

	_, err := sut.Update(syntheticTest)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

func TestShouldReturnErrorWhenExecutingUpdateOperationOfSyntheticTestRestResourceAndReceivedObjectIsNotValid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	syntheticTest := makeSyntheticTest()

	client.EXPECT().Put(gomock.Eq(syntheticTest), gomock.Eq(SyntheticTestResourcePath)).Times(1)
	client.EXPECT().GetOne(syntheticTestID, SyntheticTestResourcePath).Times(1).Return(syntheticTestSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(syntheticTestSerialized).Times(1).Return(&SyntheticTest{}, nil)

	sut := NewSyntheticTestRestResource(unmarshaller, client)

	_, err := sut.Update(syntheticTest)

	require.Error(t, err)
	require.Contains(t, "id is missing", err.Error())
}

// ########################################################
// Delete Operation Tests
// ########################################################

func TestShouldSuccessfullyExecuteDeleteByObjectOperationOfSyntheticTestRestResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	syntheticTest := makeSyntheticTest()

	client.EXPECT().Delete(syntheticTestID, SyntheticTestResourcePath).Times(1).Return(nil)

	sut := NewSyntheticTestRestResource(unmarshaller, client)

	err := sut.Delete(syntheticTest)

	require.NoError(t, err)
}

func TestShouldReturnErrorWhenExecutingDeleteByObjectOperationOfSyntheticTestRestResourceAndDeleteRequestFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	expectedError := errors.New("Error")
	syntheticTest := makeSyntheticTest()

	client.EXPECT().Delete(syntheticTestID, SyntheticTestResourcePath).Times(1).Return(expectedError)

	sut := NewSyntheticTestRestResource(unmarshaller, client)

	err := sut.Delete(syntheticTest)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}
func TestShouldSuccessfullyExecuteDeleteByIdOperationOfSyntheticTestRestResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)

	client.EXPECT().Delete(syntheticTestID, SyntheticTestResourcePath).Times(1).Return(nil)

	sut := NewSyntheticTestRestResource(unmarshaller, client)

	err := sut.DeleteByID(syntheticTestID)

	require.NoError(t, err)
}

func TestShouldReturnErrorWhenExecutingDeleteByIdOperationOfSyntheticTestRestResourceAndDeleteRequestFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	expectedError := errors.New("Error")

	client.EXPECT().Delete(syntheticTestID, SyntheticTestResourcePath).Times(1).Return(expectedError)

	sut := NewSyntheticTestRestResource(unmarshaller, client)

	err := sut.DeleteByID(syntheticTestID)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}
