package restapi_test

import (
	"encoding/json"
	"errors"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	mocks "github.com/gessnerfl/terraform-provider-instana/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSuccessfulCreateOfTestObjectThroughCreatePOSTUpdatePUTRestResource(t *testing.T) {
	executeCreateOrUpdateOperationThroughCreatePOSTUpdatePUTRestResourceTest(t, func(t *testing.T, sut RestResource, client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller) {
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
	executeCreateOrUpdateOperationThroughCreatePOSTUpdatePUTRestResourceTest(t, func(t *testing.T, sut RestResource, client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller) {
		testObject := makeTestObject()

		client.EXPECT().Post(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return(nil, errors.New("Error during test"))
		client.EXPECT().Put(gomock.Any(), gomock.Eq(testObjectResourcePath)).Times(0)
		unmarshaller.EXPECT().Unmarshal(gomock.Any()).Times(0)

		_, err := sut.Create(testObject)

		assert.Error(t, err)
	})
}

func TestShouldFailToCreateTestObjectThroughCreatePOSTUpdatePUTRestResourceWhenResponseCannotBeUnmarshalled(t *testing.T) {
	executeCreateOrUpdateOperationThroughCreatePOSTUpdatePUTRestResourceTest(t, func(t *testing.T, sut RestResource, client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller) {
		testObject := makeTestObject()
		response := []byte("invalid response")
		expectedError := errors.New("test")

		client.EXPECT().Post(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return(response, nil)
		client.EXPECT().Put(gomock.Any(), gomock.Eq(testObjectResourcePath)).Times(0)
		unmarshaller.EXPECT().Unmarshal(response).Times(1).Return(nil, expectedError)

		_, err := sut.Create(testObject)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}

func TestShouldFailToCreateTestObjectThroughCreatePOSTUpdatePUTRestResourceWhenTheUnmarshalledResponseIsNotImplementingInstanaDataObject(t *testing.T) {
	executeCreateOrUpdateOperationThroughCreatePOSTUpdatePUTRestResourceTest(t, func(t *testing.T, sut RestResource, client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller) {
		testObject := makeTestObject()
		response := []byte("{ \"invalid\" : \"testObject\" }")

		client.EXPECT().Post(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return(response, nil)
		client.EXPECT().Put(gomock.Any(), gomock.Eq(testObjectResourcePath)).Times(0)
		unmarshaller.EXPECT().Unmarshal(response).Times(1).Return(&InvalidInstanaDataObject{}, nil)

		_, err := sut.Create(testObject)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Unmarshalled object does not implement InstanaDataObject")
	})
}

func TestShouldFailedToCreateTestObjectThroughCreatePOSTUpdatePUTRestResourceWhenAnInvalidTestObjectIsProvided(t *testing.T) {
	executeCreateOrUpdateOperationThroughCreatePOSTUpdatePUTRestResourceTest(t, func(t *testing.T, sut RestResource, client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller) {
		testObject := &testObject{
			ID:   "some id",
			Name: "invalid name",
		}

		client.EXPECT().Post(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Times(0)
		client.EXPECT().Put(gomock.Any(), gomock.Eq(testObjectResourcePath)).Times(0)
		unmarshaller.EXPECT().Unmarshal(gomock.Any()).Times(0)

		_, err := sut.Create(testObject)

		assert.Error(t, err)
	})
}

func TestShouldFailedToCreateTestObjectThroughCreatePOSTUpdatePUTRestResourceWhenAnInvalidTestObjectIsReceived(t *testing.T) {
	executeCreateOrUpdateOperationThroughCreatePOSTUpdatePUTRestResourceTest(t, func(t *testing.T, sut RestResource, client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller) {
		object := makeTestObject()
		response := []byte("invalid response")

		client.EXPECT().Post(gomock.Eq(object), gomock.Eq(testObjectResourcePath)).Times(1).Return(response, nil)
		client.EXPECT().Put(gomock.Any(), gomock.Eq(testObjectResourcePath)).Times(0)
		unmarshaller.EXPECT().Unmarshal(response).Times(1).Return(&testObject{ID: object.ID, Name: "invalid"}, nil)

		_, err := sut.Create(object)

		assert.Error(t, err)
	})
}

func TestSuccessfulUpdateOfTestObjectThroughCreatePOSTUpdatePUTRestResource(t *testing.T) {
	executeCreateOrUpdateOperationThroughCreatePOSTUpdatePUTRestResourceTest(t, func(t *testing.T, sut RestResource, client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller) {
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
	executeCreateOrUpdateOperationThroughCreatePOSTUpdatePUTRestResourceTest(t, func(t *testing.T, sut RestResource, client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller) {
		testObject := makeTestObject()

		client.EXPECT().Put(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return(nil, errors.New("Error during test"))
		client.EXPECT().Post(gomock.Any(), gomock.Eq(testObjectResourcePath)).Times(0)
		unmarshaller.EXPECT().Unmarshal(gomock.Any()).Times(0)

		_, err := sut.Update(testObject)

		assert.Error(t, err)
	})
}

func TestShouldFailToUpdateTestObjectThroughCreatePOSTUpdatePUTRestResourceWhenResponseCannotBeUnmarshalled(t *testing.T) {
	executeCreateOrUpdateOperationThroughCreatePOSTUpdatePUTRestResourceTest(t, func(t *testing.T, sut RestResource, client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller) {
		testObject := makeTestObject()
		response := []byte("invalid response")
		expectedError := errors.New("test")

		client.EXPECT().Put(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return(response, nil)
		client.EXPECT().Post(gomock.Any(), gomock.Eq(testObjectResourcePath)).Times(0)
		unmarshaller.EXPECT().Unmarshal(response).Times(1).Return(nil, expectedError)

		_, err := sut.Update(testObject)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}

func TestShouldFailToUpdateTestObjectThroughCreatePOSTUpdatePUTRestResourceWhenTheUnmarshalledResponseIsNotImplementingInstanaDataObject(t *testing.T) {
	executeCreateOrUpdateOperationThroughCreatePOSTUpdatePUTRestResourceTest(t, func(t *testing.T, sut RestResource, client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller) {
		testObject := makeTestObject()
		response := []byte("{ \"invalid\" : \"testObject\" }")

		client.EXPECT().Put(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return(response, nil)
		client.EXPECT().Post(gomock.Any(), gomock.Eq(testObjectResourcePath)).Times(0)
		unmarshaller.EXPECT().Unmarshal(response).Times(1).Return(&InvalidInstanaDataObject{}, nil)

		_, err := sut.Update(testObject)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Unmarshalled object does not implement InstanaDataObject")
	})
}

func TestShouldFailedToUpdateTestObjectThroughCreatePOSTUpdatePUTRestResourceWhenAnInvalidTestObjectIsProvided(t *testing.T) {
	executeCreateOrUpdateOperationThroughCreatePOSTUpdatePUTRestResourceTest(t, func(t *testing.T, sut RestResource, client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller) {
		testObject := &testObject{
			ID:   "some id",
			Name: "invalid name",
		}

		client.EXPECT().Put(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Times(0)
		client.EXPECT().Post(gomock.Any(), gomock.Eq(testObjectResourcePath)).Times(0)
		unmarshaller.EXPECT().Unmarshal(gomock.Any()).Times(0)

		_, err := sut.Update(testObject)

		assert.Error(t, err)
	})
}

func TestShouldFailedToUpdateTestObjectThroughCreatePOSTUpdatePUTRestResourceWhenAnInvalidTestObjectIsReceived(t *testing.T) {
	executeCreateOrUpdateOperationThroughCreatePOSTUpdatePUTRestResourceTest(t, func(t *testing.T, sut RestResource, client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller) {
		object := makeTestObject()
		response := []byte("invalid response")

		client.EXPECT().Put(gomock.Eq(object), gomock.Eq(testObjectResourcePath)).Times(1).Return(response, nil)
		client.EXPECT().Post(gomock.Any(), gomock.Eq(testObjectResourcePath)).Times(0)
		unmarshaller.EXPECT().Unmarshal(response).Times(1).Return(&testObject{ID: object.ID, Name: "invalid"}, nil)

		_, err := sut.Update(object)

		assert.Error(t, err)
	})
}

func executeCreateOrUpdateOperationThroughCreatePOSTUpdatePUTRestResourceTest(t *testing.T, testFunction func(t *testing.T, sut RestResource, client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller)) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)

	sut := NewCreatePOSTUpdatePUTRestResource(testObjectResourcePath, unmarshaller, client)

	testFunction(t, sut, client, unmarshaller)
}
