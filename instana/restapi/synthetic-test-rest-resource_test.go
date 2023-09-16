package restapi_test

import (
	"errors"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
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

func makeSyntheticTest() *SyntheticTest {
	url := syntheticTestConfigUrl
	operation := syntheticTestConfigOperation
	return &SyntheticTest{
		ID:        syntheticTestID,
		Label:     syntheticTestLabel,
		Active:    syntheticTestActive,
		Locations: []string{syntheticTestLocation},
		Configuration: SyntheticTestConfig{
			SyntheticType: syntheticTestConfigSyntheticType,
			URL:           &url,
			Operation:     &operation,
		},
	}
}

// ########################################################
// GET All Tests
// ########################################################
func TestShouldSuccessfullyGetAllSyntheticTests(t *testing.T) {
	testObject := makeSyntheticTest()
	expectedResult := []*SyntheticTest{testObject, testObject, testObject}
	restResponseData := []byte("server-response")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	restClient := mocks.NewMockRestClient(ctrl)
	restClient.EXPECT().Get(SyntheticTestResourcePath).Times(1).Return(restResponseData, nil)

	jsonUnmarshaller := mocks.NewMockJSONUnmarshaller[*SyntheticTest](ctrl)
	jsonUnmarshaller.EXPECT().UnmarshalArray(restResponseData).Times(1).Return(&expectedResult, nil)

	sut := NewSyntheticTestRestResource(jsonUnmarshaller, restClient)

	result, err := sut.GetAll()

	require.NoError(t, err)
	require.Equal(t, &expectedResult, result)
}

func TestShouldReturnEmptySliceWhenNoSyntheticTestsIsReturnedForGetAll(t *testing.T) {
	restResponseData := []byte("[]")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	restClient := mocks.NewMockRestClient(ctrl)
	restClient.EXPECT().Get(SyntheticTestResourcePath).Times(1).Return(restResponseData, nil)

	jsonUnmarshaller := mocks.NewMockJSONUnmarshaller[*SyntheticTest](ctrl)
	jsonUnmarshaller.EXPECT().UnmarshalArray(restResponseData).Times(1).Return(&[]*SyntheticTest{}, nil)

	sut := NewSyntheticTestRestResource(jsonUnmarshaller, restClient)

	result, err := sut.GetAll()

	require.NoError(t, err)
	require.Equal(t, &[]*SyntheticTest{}, result)
}

func TestShouldFailToGetAllSyntheticTestsWhenClientReturnsError(t *testing.T) {
	expectedError := errors.New("test")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	restClient := mocks.NewMockRestClient(ctrl)
	restClient.EXPECT().Get(SyntheticTestResourcePath).Times(1).Return(nil, expectedError)

	jsonUnmarshaller := mocks.NewMockJSONUnmarshaller[*SyntheticTest](ctrl)
	jsonUnmarshaller.EXPECT().UnmarshalArray(gomock.Any()).Times(0)

	sut := NewSyntheticTestRestResource(jsonUnmarshaller, restClient)

	_, err := sut.GetAll()

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

func TestShouldFailToGetAllSyntheticTestsWhenRestResultCannotBeUnmarshalled(t *testing.T) {
	restResponseData := []byte("invalidResponse")
	expectedError := errors.New("test")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	restClient := mocks.NewMockRestClient(ctrl)
	restClient.EXPECT().Get(SyntheticTestResourcePath).Times(1).Return(restResponseData, nil)

	jsonUnmarshaller := mocks.NewMockJSONUnmarshaller[*SyntheticTest](ctrl)
	jsonUnmarshaller.EXPECT().UnmarshalArray(restResponseData).Times(1).Return(nil, expectedError)

	sut := NewSyntheticTestRestResource(jsonUnmarshaller, restClient)

	_, err := sut.GetAll()

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

// ########################################################
// GET Operation Tests
// ########################################################

func TestShouldSuccessfullyExecuteGetOperationOfSyntheticTestRestResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller[*SyntheticTest](ctrl)
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
	unmarshaller := mocks.NewMockJSONUnmarshaller[*SyntheticTest](ctrl)
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
	unmarshaller := mocks.NewMockJSONUnmarshaller[*SyntheticTest](ctrl)
	expectedError := errors.New("Error")

	client.EXPECT().GetOne(syntheticTestID, SyntheticTestResourcePath).Times(1).Return(syntheticTestSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(syntheticTestSerialized).Times(1).Return(&SyntheticTest{}, expectedError)

	sut := NewSyntheticTestRestResource(unmarshaller, client)

	_, err := sut.GetOne(syntheticTestID)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

// ########################################################
// Create Operation Tests
// ########################################################

func TestShouldSuccessfullyExecuteCreateOperationOfSyntheticTestRestResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller[*SyntheticTest](ctrl)
	syntheticTest := makeSyntheticTest()

	client.EXPECT().Post(gomock.Eq(syntheticTest), gomock.Eq(SyntheticTestResourcePath)).Times(1).Return(syntheticTestSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(syntheticTestSerialized).Times(1).Return(syntheticTest, nil)

	sut := NewSyntheticTestRestResource(unmarshaller, client)

	result, err := sut.Create(syntheticTest)

	require.NoError(t, err)
	require.Equal(t, syntheticTest, result)
}

func TestShouldReturnErrorWhenExecutingCreateOperationOfSyntheticTestRestResourceAndPostOperationFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller[*SyntheticTest](ctrl)
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
	unmarshaller := mocks.NewMockJSONUnmarshaller[*SyntheticTest](ctrl)
	expectedError := errors.New("Error")
	syntheticTest := makeSyntheticTest()

	client.EXPECT().Post(gomock.Eq(syntheticTest), gomock.Eq(SyntheticTestResourcePath)).Times(1).Return(syntheticTestSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(syntheticTestSerialized).Times(1).Return(&SyntheticTest{}, expectedError)

	sut := NewSyntheticTestRestResource(unmarshaller, client)

	_, err := sut.Create(syntheticTest)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

// ########################################################
// Update Operation Tests
// ########################################################

func TestShouldSuccessfullyExecuteUpdateOperationOfSyntheticTestRestResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller[*SyntheticTest](ctrl)
	syntheticTest := makeSyntheticTest()

	client.EXPECT().Put(gomock.Eq(syntheticTest), gomock.Eq(SyntheticTestResourcePath)).Times(1)
	client.EXPECT().GetOne(syntheticTestID, SyntheticTestResourcePath).Times(1).Return(syntheticTestSerialized, nil)
	unmarshaller.EXPECT().Unmarshal(syntheticTestSerialized).Times(1).Return(syntheticTest, nil)

	sut := NewSyntheticTestRestResource(unmarshaller, client)

	result, err := sut.Update(syntheticTest)

	require.NoError(t, err)
	require.Equal(t, syntheticTest, result)
}

func TestShouldReturnErrorWhenExecutingUpdateOperationOfSyntheticTestRestResourceAndPutOperationFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller[*SyntheticTest](ctrl)
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
	unmarshaller := mocks.NewMockJSONUnmarshaller[*SyntheticTest](ctrl)
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

// ########################################################
// Delete Operation Tests
// ########################################################

func TestShouldSuccessfullyExecuteDeleteByObjectOperationOfSyntheticTestRestResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller[*SyntheticTest](ctrl)
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
	unmarshaller := mocks.NewMockJSONUnmarshaller[*SyntheticTest](ctrl)
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
	unmarshaller := mocks.NewMockJSONUnmarshaller[*SyntheticTest](ctrl)

	client.EXPECT().Delete(syntheticTestID, SyntheticTestResourcePath).Times(1).Return(nil)

	sut := NewSyntheticTestRestResource(unmarshaller, client)

	err := sut.DeleteByID(syntheticTestID)

	require.NoError(t, err)
}

func TestShouldReturnErrorWhenExecutingDeleteByIdOperationOfSyntheticTestRestResourceAndDeleteRequestFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller[*SyntheticTest](ctrl)
	expectedError := errors.New("Error")

	client.EXPECT().Delete(syntheticTestID, SyntheticTestResourcePath).Times(1).Return(expectedError)

	sut := NewSyntheticTestRestResource(unmarshaller, client)

	err := sut.DeleteByID(syntheticTestID)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}
