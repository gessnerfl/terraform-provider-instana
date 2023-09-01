package restapi_test

import (
	"encoding/json"
	"errors"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var (
	invalidResponse = []byte("invalid response")
)

func TestSuccessfulCreateOfTestObjectThroughCreatePOSTUpdatePUTRestResource(t *testing.T) {
	executeCreateOrUpdateOperationThroughCreatePOSTUpdatePUTRestResourceTest(t, func(t *testing.T, sut RestResource[*testObject], client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller[*testObject]) {
		testObject := makeTestObject()
		serializedJSON, _ := json.Marshal(testObject)

		client.EXPECT().Post(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return(serializedJSON, nil)
		unmarshaller.EXPECT().Unmarshal(serializedJSON).Times(1).Return(testObject, nil)

		result, err := sut.Create(testObject)

		assert.NoError(t, err)
		assert.Equal(t, testObject, result)
	})
}

func TestShouldFailToCreateTestObjectThroughCreatePOSTUpdatePUTRestResourceWhenErrorIsReturnedFromRestClient(t *testing.T) {
	executeCreateOrUpdateOperationThroughCreatePOSTUpdatePUTRestResourceTest(t, func(t *testing.T, sut RestResource[*testObject], client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller[*testObject]) {
		testObject := makeTestObject()

		client.EXPECT().Post(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return(nil, errors.New("error during test"))
		client.EXPECT().Put(gomock.Any(), gomock.Eq(testObjectResourcePath)).Times(0)
		unmarshaller.EXPECT().Unmarshal(gomock.Any()).Times(0)

		_, err := sut.Create(testObject)

		assert.Error(t, err)
	})
}

func TestShouldFailToCreateTestObjectThroughCreatePOSTUpdatePUTRestResourceWhenResponseCannotBeUnmarshalled(t *testing.T) {
	executeCreateOrUpdateOperationThroughCreatePOSTUpdatePUTRestResourceTest(t, func(t *testing.T, sut RestResource[*testObject], client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller[*testObject]) {
		testObject := makeTestObject()
		expectedError := errors.New("test")

		client.EXPECT().Post(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return(invalidResponse, nil)
		client.EXPECT().Put(gomock.Any(), gomock.Eq(testObjectResourcePath)).Times(0)
		unmarshaller.EXPECT().Unmarshal(invalidResponse).Times(1).Return(nil, expectedError)

		_, err := sut.Create(testObject)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}

func TestSuccessfulUpdateOfTestObjectThroughCreatePOSTUpdatePUTRestResource(t *testing.T) {
	executeCreateOrUpdateOperationThroughCreatePOSTUpdatePUTRestResourceTest(t, func(t *testing.T, sut RestResource[*testObject], client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller[*testObject]) {
		testObject := makeTestObject()
		serializedJSON, _ := json.Marshal(testObject)

		client.EXPECT().Put(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return(serializedJSON, nil)
		client.EXPECT().Post(gomock.Any(), gomock.Eq(testObjectResourcePath)).Times(0)
		unmarshaller.EXPECT().Unmarshal(serializedJSON).Times(1).Return(testObject, nil)

		result, err := sut.Update(testObject)

		assert.NoError(t, err)
		assert.Equal(t, testObject, result)
	})
}

func TestShouldFailToUpdateTestObjectThroughCreatePOSTUpdatePUTRestResourceWhenErrorIsReturnedFromRestClient(t *testing.T) {
	executeCreateOrUpdateOperationThroughCreatePOSTUpdatePUTRestResourceTest(t, func(t *testing.T, sut RestResource[*testObject], client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller[*testObject]) {
		testObject := makeTestObject()

		client.EXPECT().Put(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return(nil, errors.New("error during test"))
		client.EXPECT().Post(gomock.Any(), gomock.Eq(testObjectResourcePath)).Times(0)
		unmarshaller.EXPECT().Unmarshal(gomock.Any()).Times(0)

		_, err := sut.Update(testObject)

		assert.Error(t, err)
	})
}

func TestShouldFailToUpdateTestObjectThroughCreatePOSTUpdatePUTRestResourceWhenResponseCannotBeUnmarshalled(t *testing.T) {
	executeCreateOrUpdateOperationThroughCreatePOSTUpdatePUTRestResourceTest(t, func(t *testing.T, sut RestResource[*testObject], client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller[*testObject]) {
		testObject := makeTestObject()
		expectedError := errors.New("test")

		client.EXPECT().Put(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return(invalidResponse, nil)
		client.EXPECT().Post(gomock.Any(), gomock.Eq(testObjectResourcePath)).Times(0)
		unmarshaller.EXPECT().Unmarshal(invalidResponse).Times(1).Return(nil, expectedError)

		_, err := sut.Update(testObject)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}

func executeCreateOrUpdateOperationThroughCreatePOSTUpdatePUTRestResourceTest(t *testing.T, testFunction func(t *testing.T, sut RestResource[*testObject], client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller[*testObject])) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller[*testObject](ctrl)

	sut := NewCreatePOSTUpdatePUTRestResource[*testObject](testObjectResourcePath, unmarshaller, client)

	testFunction(t, sut, client, unmarshaller)
}
