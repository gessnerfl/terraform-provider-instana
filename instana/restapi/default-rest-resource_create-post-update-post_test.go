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

func TestSuccessfulCreateOfTestObjectThroughCreatePOSTUpdatePOSTRestResource(t *testing.T) {
	executeCreateOrUpdateOperationThroughCreatePOSTUpdatePOSTRestResourceTest(t, func(t *testing.T, sut RestResource[*testObject], client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller[*testObject]) {
		testObject := makeTestObject()
		serializedJSON, _ := json.Marshal(testObject)

		client.EXPECT().Post(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return(serializedJSON, nil)
		client.EXPECT().PostWithID(gomock.Any(), gomock.Eq(testObjectResourcePath)).Times(0)
		unmarshaller.EXPECT().Unmarshal(serializedJSON).Times(1).Return(testObject, nil)

		result, err := sut.Create(testObject)

		assert.NoError(t, err)
		assert.Equal(t, testObject, result)
	})
}

func TestShouldFailToCreateTestObjectThroughCreatePOSTUpdatePOSTRestResourceWhenErrorIsReturnedFromRestClient(t *testing.T) {
	executeCreateOrUpdateOperationThroughCreatePOSTUpdatePOSTRestResourceTest(t, func(t *testing.T, sut RestResource[*testObject], client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller[*testObject]) {
		testObject := makeTestObject()

		client.EXPECT().Post(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return(nil, errors.New("Error during test"))
		client.EXPECT().PostWithID(gomock.Any(), gomock.Eq(testObjectResourcePath)).Times(0)
		client.EXPECT().Put(gomock.Any(), gomock.Eq(testObjectResourcePath)).Times(0)
		unmarshaller.EXPECT().Unmarshal(gomock.Any()).Times(0)

		_, err := sut.Create(testObject)

		assert.Error(t, err)
	})
}

func TestShouldFailToCreateTestObjectThroughCreatePOSTUpdatePOSTRestResourceWhenResponseCannotBeUnmarshalled(t *testing.T) {
	executeCreateOrUpdateOperationThroughCreatePOSTUpdatePOSTRestResourceTest(t, func(t *testing.T, sut RestResource[*testObject], client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller[*testObject]) {
		testObject := makeTestObject()
		expectedError := errors.New("test")

		client.EXPECT().Post(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return(invalidResponse, nil)
		client.EXPECT().PostWithID(gomock.Any(), gomock.Eq(testObjectResourcePath)).Times(0)
		client.EXPECT().Put(gomock.Any(), gomock.Eq(testObjectResourcePath)).Times(0)
		unmarshaller.EXPECT().Unmarshal(invalidResponse).Times(1).Return(nil, expectedError)

		_, err := sut.Create(testObject)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}

func TestSuccessfulUpdateOfTestObjectThroughCreatePOSTUpdatePOSTRestResource(t *testing.T) {
	executeCreateOrUpdateOperationThroughCreatePOSTUpdatePOSTRestResourceTest(t, func(t *testing.T, sut RestResource[*testObject], client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller[*testObject]) {
		testObject := makeTestObject()
		serializedJSON, _ := json.Marshal(testObject)

		client.EXPECT().PostWithID(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return(serializedJSON, nil)
		client.EXPECT().Post(gomock.Any(), gomock.Eq(testObjectResourcePath)).Times(0)
		client.EXPECT().Put(gomock.Any(), gomock.Eq(testObjectResourcePath)).Times(0)
		unmarshaller.EXPECT().Unmarshal(serializedJSON).Times(1).Return(testObject, nil)

		result, err := sut.Update(testObject)

		assert.NoError(t, err)
		assert.Equal(t, testObject, result)
	})
}

func TestShouldFailToUpdateTestObjectThroughCreatePOSTUpdatePOSTRestResourceWhenErrorIsReturnedFromRestClient(t *testing.T) {
	executeCreateOrUpdateOperationThroughCreatePOSTUpdatePOSTRestResourceTest(t, func(t *testing.T, sut RestResource[*testObject], client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller[*testObject]) {
		testObject := makeTestObject()

		client.EXPECT().PostWithID(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return(nil, errors.New("Error during test"))
		client.EXPECT().Post(gomock.Any(), gomock.Eq(testObjectResourcePath)).Times(0)
		client.EXPECT().Put(gomock.Any(), gomock.Eq(testObjectResourcePath)).Times(0)
		unmarshaller.EXPECT().Unmarshal(gomock.Any()).Times(0)

		_, err := sut.Update(testObject)

		assert.Error(t, err)
	})
}

func TestShouldFailToUpdateTestObjectThroughCreatePOSTUpdatePOSTRestResourceWhenResponseCannotBeUnmarshalled(t *testing.T) {
	executeCreateOrUpdateOperationThroughCreatePOSTUpdatePOSTRestResourceTest(t, func(t *testing.T, sut RestResource[*testObject], client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller[*testObject]) {
		testObject := makeTestObject()
		expectedError := errors.New("test")

		client.EXPECT().PostWithID(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return(invalidResponse, nil)
		client.EXPECT().Post(gomock.Any(), gomock.Eq(testObjectResourcePath)).Times(0)
		client.EXPECT().Put(gomock.Any(), gomock.Eq(testObjectResourcePath)).Times(0)
		unmarshaller.EXPECT().Unmarshal(invalidResponse).Times(1).Return(nil, expectedError)

		_, err := sut.Update(testObject)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}

func executeCreateOrUpdateOperationThroughCreatePOSTUpdatePOSTRestResourceTest(t *testing.T, testFunction func(t *testing.T, sut RestResource[*testObject], client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller[*testObject])) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller[*testObject](ctrl)

	sut := NewCreatePOSTUpdatePOSTRestResource[*testObject](testObjectResourcePath, unmarshaller, client)

	testFunction(t, sut, client, unmarshaller)
}
