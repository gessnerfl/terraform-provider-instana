package restapi_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	mocks "github.com/gessnerfl/terraform-provider-instana/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSuccessfulCreateOrUpdateOfTestObjectThroughCreatePUTUpdatePUTRestResource(t *testing.T) {
	executeCreateOrUpdateOperationThroughCreatePUTUpdatePUTRestResourceTest(t, func(t *testing.T, resourceFunc createUpdateFunc, client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller) {
		testObject := makeTestObject()
		serializedJSON, _ := json.Marshal(testObject)

		client.EXPECT().Put(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return(serializedJSON, nil)
		unmarshaller.EXPECT().Unmarshal(serializedJSON).Times(1).Return(testObject, nil)

		result, err := resourceFunc(testObject)

		assert.NoError(t, err)
		assert.Equal(t, testObject, result)
	})
}

func TestShouldFailToCreateOrUpdateTestObjectThroughCreatePUTUpdatePUTRestResourceWhenErrorIsReturnedFromRestClient(t *testing.T) {
	executeCreateOrUpdateOperationThroughCreatePUTUpdatePUTRestResourceTest(t, func(t *testing.T, resourceFunc createUpdateFunc, client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller) {
		testObject := makeTestObject()

		client.EXPECT().Put(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return(nil, errors.New("Error during test"))
		unmarshaller.EXPECT().Unmarshal(gomock.Any()).Times(0)

		_, err := resourceFunc(testObject)

		assert.Error(t, err)
	})
}

func TestShouldFailToCreateOrUpdateTestObjectThroughCreatePUTUpdatePUTRestResourceWhenResponseCannotBeUnmarshalled(t *testing.T) {
	executeCreateOrUpdateOperationThroughCreatePUTUpdatePUTRestResourceTest(t, func(t *testing.T, resourceFunc createUpdateFunc, client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller) {
		testObject := makeTestObject()
		response := []byte("invalid response")
		expectedError := errors.New("test")

		client.EXPECT().Put(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return(response, nil)
		unmarshaller.EXPECT().Unmarshal(response).Times(1).Return(nil, expectedError)

		_, err := resourceFunc(testObject)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}

func TestShouldFailToCreateOrUpdateTestObjectThroughCreatePUTUpdatePUTRestResourceWhenTheUnmarshalledResponseIsNotImplementingInstanaDataObject(t *testing.T) {
	executeCreateOrUpdateOperationThroughCreatePUTUpdatePUTRestResourceTest(t, func(t *testing.T, resourceFunc createUpdateFunc, client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller) {
		testObject := makeTestObject()
		response := []byte("{ \"invalid\" : \"testObject\" }")

		client.EXPECT().Put(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return(response, nil)
		unmarshaller.EXPECT().Unmarshal(response).Times(1).Return(&InvalidInstanaDataObject{}, nil)

		_, err := resourceFunc(testObject)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Unmarshalled object does not implement InstanaDataObject")
	})
}

func TestShouldFailedToCreateOrUpdateTestObjectThroughCreatePUTUpdatePUTRestResourceWhenAnInvalidTestObjectIsProvided(t *testing.T) {
	executeCreateOrUpdateOperationThroughCreatePUTUpdatePUTRestResourceTest(t, func(t *testing.T, resourceFunc createUpdateFunc, client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller) {
		testObject := &testObject{
			ID:   "some id",
			Name: "invalid name",
		}

		client.EXPECT().Put(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Times(0)
		unmarshaller.EXPECT().Unmarshal(gomock.Any()).Times(0)

		_, err := resourceFunc(testObject)

		assert.Error(t, err)
	})
}

func TestShouldFailedToCreateOrUpdateTestObjectThroughCreatePUTUpdatePUTRestResourceWhenAnInvalidTestObjectIsReceived(t *testing.T) {
	executeCreateOrUpdateOperationThroughCreatePUTUpdatePUTRestResourceTest(t, func(t *testing.T, resourceFunc createUpdateFunc, client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller) {
		object := makeTestObject()
		response := []byte("invalid response")

		client.EXPECT().Put(gomock.Eq(object), gomock.Eq(testObjectResourcePath)).Times(1).Return(response, nil)
		unmarshaller.EXPECT().Unmarshal(response).Times(1).Return(&testObject{ID: object.ID, Name: "invalid"}, nil)

		_, err := resourceFunc(object)

		assert.Error(t, err)
	})
}

type createUpdateFunc func(data InstanaDataObject) (InstanaDataObject, error)
type createPutUpdatePutContext struct {
	operation           string
	resourceFuncFactory func(RestResource) createUpdateFunc
}

func executeCreateOrUpdateOperationThroughCreatePUTUpdatePUTRestResourceTest(t *testing.T, testFunction func(t *testing.T, resourceFunc createUpdateFunc, client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller)) {
	contexts := []createPutUpdatePutContext{
		{operation: "Create", resourceFuncFactory: func(sut RestResource) createUpdateFunc { return sut.Create }},
		{operation: "Update", resourceFuncFactory: func(sut RestResource) createUpdateFunc { return sut.Update }},
	}

	caller := getCallerName()
	for _, context := range contexts {
		fullName := fmt.Sprintf("%s[%s]", caller, context.operation)
		t.Run(fullName, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			client := mocks.NewMockRestClient(ctrl)
			unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)

			sut := NewCreatePUTUpdatePUTRestResource(testObjectResourcePath, unmarshaller, client)

			client.EXPECT().Post(gomock.Any(), gomock.Eq(testObjectResourcePath)).Times(0)
			testFunction(t, context.resourceFuncFactory(sut), client, unmarshaller)
		})
	}
}
