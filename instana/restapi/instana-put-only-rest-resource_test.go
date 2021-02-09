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

const (
	testObjectResourcePath = "/test"
	testObjectID           = "test-object-id"
	testObjectName         = "test-object-name"
)

type testObject struct {
	ID   string
	Name string
}

func (t *testObject) GetID() string {
	return t.ID
}

func (t *testObject) Validate() error {
	if t.Name != testObjectName {
		return errors.New("Name differs")
	}
	return nil
}

func makeTestObject() *testObject {
	return &testObject{ID: testObjectID, Name: testObjectName}
}

func TestSuccessfulGetOneTestObjectThroughPutOnlyRestResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)

	sut := NewPUTOnlyRestResource(testObjectResourcePath, unmarshaller, client)

	testObject := makeTestObject()
	serializedJSON, _ := json.Marshal(testObject)

	client.EXPECT().GetOne(gomock.Eq(testObject.ID), gomock.Eq(testObjectResourcePath)).Return(serializedJSON, nil)
	unmarshaller.EXPECT().Unmarshal(serializedJSON).Times(1).Return(testObject, nil)

	data, err := sut.GetOne(testObject.ID)

	assert.NoError(t, err)
	assert.Equal(t, testObject, data)
}

func TestShouldFailToGetOneTestObjectThroughPutOnlyRestResourceWhenErrorIsRetrievedFromRestClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)

	sut := NewPUTOnlyRestResource(testObjectResourcePath, unmarshaller, client)

	client.EXPECT().GetOne(gomock.Eq(testObjectID), gomock.Eq(testObjectResourcePath)).Return(nil, errors.New("error during test"))
	unmarshaller.EXPECT().Unmarshal(gomock.Any()).Times(0)

	_, err := sut.GetOne(testObjectID)

	assert.Error(t, err)
}

func TestShouldFailToGetOneTestObjectThroughPutOnlyRestResourceWhenResponseCannotBeUnmarshalled(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)

	expectedError := errors.New("test")
	response := []byte("[{ \"invalid\" : \"data\" }]")

	sut := NewPUTOnlyRestResource(testObjectResourcePath, unmarshaller, client)

	client.EXPECT().GetOne(gomock.Eq(testObjectID), gomock.Eq(testObjectResourcePath)).Return(response, nil)
	unmarshaller.EXPECT().Unmarshal(response).Times(1).Return(nil, expectedError)

	_, err := sut.GetOne(testObjectID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
}

type InvalidInstanaDataObject struct{}

func TestShouldFailToGetOneTestObjectThroughPutOnlyRestResourceWhenUnmarshalledObjectIsNotAInstanaDataObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)

	response := []byte("[{ \"some\" : \"data\" }]")

	sut := NewPUTOnlyRestResource(testObjectResourcePath, unmarshaller, client)

	client.EXPECT().GetOne(gomock.Eq(testObjectID), gomock.Eq(testObjectResourcePath)).Return(response, nil)
	unmarshaller.EXPECT().Unmarshal(response).Times(1).Return(&InvalidInstanaDataObject{}, nil)

	_, err := sut.GetOne(testObjectID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Unmarshalled object does not implement InstanaDataObject")
}

func TestShouldFailToGetOneTestObjectThroughPutOnlyRestResourceWhenUnmarshalledObjectIsNotValid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)

	response := []byte("[{ \"some\" : \"data\" }]")
	object := makeTestObject()
	object.Name = "invalid"

	sut := NewPUTOnlyRestResource(testObjectResourcePath, unmarshaller, client)

	client.EXPECT().GetOne(gomock.Eq(testObjectID), gomock.Eq(testObjectResourcePath)).Return(response, nil)
	unmarshaller.EXPECT().Unmarshal(response).Times(1).Return(object, nil)

	_, err := sut.GetOne(testObjectID)

	assert.Error(t, err)
}

var upsertOperations = []string{"Create", "Update"}

func executeUpsertOperationThroughPutOnlyRestResourceTest(t *testing.T, name string, testFunction func(operation string, t *testing.T)) {
	for _, operation := range upsertOperations {
		fullName := fmt.Sprintf("TestSuccessfulUpsertOfTestObject_%s", operation)
		t.Run(fullName, func(t *testing.T) {
			testFunction(operation, t)
		})
	}
}

func executeUpsertOperationThroughPutOnlyRestResource(operation string, sut RestResource, testObject InstanaDataObject) (InstanaDataObject, error) {
	if operation == "Create" {
		return sut.Create(testObject)
	}
	return sut.Update(testObject)
}

func TestSuccessfulUpsertOfTestObjectThroughPutOnlyRestResource(t *testing.T) {
	executeUpsertOperationThroughPutOnlyRestResourceTest(t, "TestSuccessfulUpsertOfTestObject", func(operation string, t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		client := mocks.NewMockRestClient(ctrl)
		unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)

		sut := NewPUTOnlyRestResource(testObjectResourcePath, unmarshaller, client)

		testObject := makeTestObject()
		serializedJSON, _ := json.Marshal(testObject)

		client.EXPECT().Put(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return(serializedJSON, nil)
		unmarshaller.EXPECT().Unmarshal(serializedJSON).Times(1).Return(testObject, nil)

		result, err := executeUpsertOperationThroughPutOnlyRestResource(operation, sut, testObject)

		assert.NoError(t, err)
		assert.Equal(t, testObject, result)
	})
}

func TestShouldFailToUpsertTestObjectThroughPutOnlyRestResourceWhenErrorIsReturnedFromRestClient(t *testing.T) {
	executeUpsertOperationThroughPutOnlyRestResourceTest(t, "TestShouldFailToUpsertTestObjectWhenErrorIsReturnedFromRestClient", func(operation string, t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		client := mocks.NewMockRestClient(ctrl)
		unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)

		sut := NewPUTOnlyRestResource(testObjectResourcePath, unmarshaller, client)
		testObject := makeTestObject()

		client.EXPECT().Put(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return(nil, errors.New("Error during test"))
		unmarshaller.EXPECT().Unmarshal(gomock.Any()).Times(0)

		_, err := executeUpsertOperationThroughPutOnlyRestResource(operation, sut, testObject)

		assert.Error(t, err)
	})
}

func TestShouldFailToUpsertTestObjectThroughPutOnlyRestResourceWhenResponseCannotBeUnmarshalled(t *testing.T) {
	executeUpsertOperationThroughPutOnlyRestResourceTest(t, "TestShouldFailToUpsertTestObjectWhenResponseMessageIsNotAValidJsonObject", func(operation string, t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		client := mocks.NewMockRestClient(ctrl)
		unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)

		sut := NewPUTOnlyRestResource(testObjectResourcePath, unmarshaller, client)
		testObject := makeTestObject()
		response := []byte("invalid response")
		expectedError := errors.New("test")

		client.EXPECT().Put(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return(response, nil)
		unmarshaller.EXPECT().Unmarshal(response).Times(1).Return(nil, expectedError)

		_, err := executeUpsertOperationThroughPutOnlyRestResource(operation, sut, testObject)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}

func TestShouldFailToUpsertTestObjectThroughPutOnlyRestResourceWhenTheUnmarshalledResponseIsNotImplementingInstanaDataObject(t *testing.T) {
	executeUpsertOperationThroughPutOnlyRestResourceTest(t, "TestShouldFailToUpsertTestObjectWhenResponseMessageContainsAnInvalidTestObject", func(operation string, t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		client := mocks.NewMockRestClient(ctrl)
		unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)

		sut := NewPUTOnlyRestResource(testObjectResourcePath, unmarshaller, client)
		testObject := makeTestObject()
		response := []byte("{ \"invalid\" : \"testObject\" }")

		client.EXPECT().Put(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return(response, nil)
		unmarshaller.EXPECT().Unmarshal(response).Times(1).Return(&InvalidInstanaDataObject{}, nil)

		_, err := executeUpsertOperationThroughPutOnlyRestResource(operation, sut, testObject)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Unmarshalled object does not implement InstanaDataObject")
	})
}

func TestShouldFailedToUpsertTestObjectThroughPutOnlyRestResourceWhenAnInvalidTestObjectIsProvided(t *testing.T) {
	executeUpsertOperationThroughPutOnlyRestResourceTest(t, "TestShouldFailedToUpsertTestObjectWhenAnInvalidTestObjectIsProvided", func(operation string, t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		client := mocks.NewMockRestClient(ctrl)
		unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)

		sut := NewPUTOnlyRestResource(testObjectResourcePath, unmarshaller, client)
		testObject := &testObject{
			ID:   "some id",
			Name: "invalid name",
		}

		client.EXPECT().Put(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Times(0)
		unmarshaller.EXPECT().Unmarshal(gomock.Any()).Times(0)

		_, err := executeUpsertOperationThroughPutOnlyRestResource(operation, sut, testObject)

		assert.Error(t, err)
	})
}

func TestShouldFailedToUpsertTestObjectThroughPutOnlyRestResourceWhenAnInvalidTestObjectIsReceived(t *testing.T) {
	executeUpsertOperationThroughPutOnlyRestResourceTest(t, "TestShouldFailedToUpsertTestObjectWhenAnInvalidTestObjectIsProvided", func(operation string, t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		client := mocks.NewMockRestClient(ctrl)
		unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)

		sut := NewPUTOnlyRestResource(testObjectResourcePath, unmarshaller, client)
		object := makeTestObject()
		response := []byte("invalid response")

		client.EXPECT().Put(gomock.Eq(object), gomock.Eq(testObjectResourcePath)).Times(1).Return(response, nil)
		unmarshaller.EXPECT().Unmarshal(response).Times(1).Return(&testObject{ID: object.ID, Name: "invalid"}, nil)

		_, err := executeUpsertOperationThroughPutOnlyRestResource(operation, sut, object)

		assert.Error(t, err)
	})
}

func TestSuccessfulDeleteOfTestObjectByObjectThroughPutOnlyRestResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)

	sut := NewPUTOnlyRestResource(testObjectResourcePath, unmarshaller, client)
	testObject := makeTestObject()

	client.EXPECT().Delete(gomock.Eq(testObjectID), gomock.Eq(testObjectResourcePath)).Return(nil)

	err := sut.Delete(testObject)

	assert.NoError(t, err)
}

func TestShouldFailToDeleteTestObjectThroughPutOnlyRestResourceWhenErrorIsRetrunedFromRestClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller(ctrl)

	sut := NewPUTOnlyRestResource(testObjectResourcePath, unmarshaller, client)
	testObject := makeTestObject()

	client.EXPECT().Delete(gomock.Eq(testObjectID), gomock.Eq(testObjectResourcePath)).Return(errors.New("Error during test"))

	err := sut.Delete(testObject)

	assert.Error(t, err)
}
