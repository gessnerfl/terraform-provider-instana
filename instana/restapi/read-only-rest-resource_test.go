package restapi_test

import (
	"errors"
	"testing"

	"github.com/gessnerfl/terraform-provider-instana/mocks"
	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

const (
	testResourcePath = "/test"
)

func TestShouldSuccessfullyGetAllObjects(t *testing.T) {
	testObject1 := newTestObject("id1", "name1")
	testObject2 := newTestObject("id2", "name2")
	testObject3 := newTestObject("id3", "name3")
	expectedResult := []InstanaDataObject{testObject1, testObject2, testObject3}
	serverResponse := []TestInstanaDataObject{testObject1, testObject2, testObject3}
	restResponseData := []byte(`
	[
		{
			"id" : "id1",
			"name": "name1"
		},
		{
			"id" : "id2",
			"name": "name2"
		},
		{
			"id" : "id3",
			"name": "name2"
		}
	]
	`)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	restClient := mocks.NewMockRestClient(ctrl)
	restClient.EXPECT().Get(testResourcePath).Times(1).Return(restResponseData, nil)

	arrayJSONUnmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	arrayJSONUnmarshaller.EXPECT().Unmarshal(restResponseData).Times(1).Return(&serverResponse, nil)

	sut := NewReadOnlyRestResource(testResourcePath, nil, arrayJSONUnmarshaller, restClient)

	result, err := sut.GetAll()

	require.NoError(t, err)
	require.Equal(t, &expectedResult, result)
}

func TestShouldReturnEmptySliceWhenNoDataIsReturned(t *testing.T) {
	restResponseData := []byte("[]")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	restClient := mocks.NewMockRestClient(ctrl)
	restClient.EXPECT().Get(testResourcePath).Times(1).Return(restResponseData, nil)

	arrayJSONUnmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	arrayJSONUnmarshaller.EXPECT().Unmarshal(restResponseData).Times(1).Return(&[]TestInstanaDataObject{}, nil)

	sut := NewReadOnlyRestResource(testResourcePath, nil, arrayJSONUnmarshaller, restClient)

	result, err := sut.GetAll()

	require.NoError(t, err)
	require.Equal(t, &[]InstanaDataObject{}, result)
}

func TestShouldFailToGetAllWhenClientReturnsError(t *testing.T) {
	expectedError := errors.New("test")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	restClient := mocks.NewMockRestClient(ctrl)
	restClient.EXPECT().Get(testResourcePath).Times(1).Return(nil, expectedError)

	arrayJSONUnmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	arrayJSONUnmarshaller.EXPECT().Unmarshal(gomock.Any()).Times(0)

	sut := NewReadOnlyRestResource(testResourcePath, nil, arrayJSONUnmarshaller, restClient)

	_, err := sut.GetAll()

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

func TestShouldFailToGetAllWhenRestResultCannotBeUnmarshalled(t *testing.T) {
	restResponseData := []byte("invalidResponse")
	expectedError := errors.New("test")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	restClient := mocks.NewMockRestClient(ctrl)
	restClient.EXPECT().Get(testResourcePath).Times(1).Return(restResponseData, nil)

	arrayJSONUnmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	arrayJSONUnmarshaller.EXPECT().Unmarshal(restResponseData).Times(1).Return(nil, expectedError)

	sut := NewReadOnlyRestResource(testResourcePath, nil, arrayJSONUnmarshaller, restClient)

	_, err := sut.GetAll()

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

func TestShouldSuccessfullyGetObjectById(t *testing.T) {
	id := "id1"
	expectedResult := newTestObject(id, "name1")
	restResponseData := []byte(`
		{
			"id" : "id1",
			"name": "name1"
		}
	`)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	restClient := mocks.NewMockRestClient(ctrl)
	restClient.EXPECT().GetOne(id, testResourcePath).Times(1).Return(restResponseData, nil)

	objectJSONUnmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	objectJSONUnmarshaller.EXPECT().Unmarshal(restResponseData).Times(1).Return(&expectedResult, nil)

	sut := NewReadOnlyRestResource(testResourcePath, objectJSONUnmarshaller, nil, restClient)

	result, err := sut.GetOne(id)

	require.NoError(t, err)
	require.Equal(t, expectedResult, result)
}

func TestShouldFailToGetObjectByIdWhenRestClientResturnsError(t *testing.T) {
	id := "id1"
	expectedError := errors.New("test")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	restClient := mocks.NewMockRestClient(ctrl)
	restClient.EXPECT().GetOne(id, testResourcePath).Times(1).Return(nil, expectedError)

	objectJSONUnmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	objectJSONUnmarshaller.EXPECT().Unmarshal(gomock.Any()).Times(0)

	sut := NewReadOnlyRestResource(testResourcePath, objectJSONUnmarshaller, nil, restClient)

	_, err := sut.GetOne(id)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

func TestShouldFailToGetObjectByIdWhenRestResultCannotBeUnmarshalled(t *testing.T) {
	id := "id1"
	restResponseData := []byte("invalidResponse")
	expectedError := errors.New("test")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	restClient := mocks.NewMockRestClient(ctrl)
	restClient.EXPECT().GetOne(id, testResourcePath).Times(1).Return(restResponseData, nil)

	objectJSONUnmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)
	objectJSONUnmarshaller.EXPECT().Unmarshal(restResponseData).Times(1).Return(nil, expectedError)

	sut := NewReadOnlyRestResource(testResourcePath, objectJSONUnmarshaller, nil, restClient)

	_, err := sut.GetOne(id)

	require.Error(t, err)
	require.Equal(t, expectedError, err)
}

func newTestObject(id string, name string) TestInstanaDataObject {
	return TestInstanaDataObject{ID: id, Name: name}
}

type TestInstanaDataObject struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

//GetIDForResourcePath implemention of the interface InstanaDataObject
func (o TestInstanaDataObject) GetIDForResourcePath() string {
	return o.ID
}

//Validate implementation of the interface InstanaDataObject to verify if data object is correct. As this is read only datasource no validation is applied
func (o TestInstanaDataObject) Validate() error {
	return nil
}
