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

type testUnmarshaller struct{}

func (t *testUnmarshaller) Unmarshal(data []byte) (InstanaDataObject, error) {
	obj := testObject{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return &obj, fmt.Errorf("failed to parse json; %s", err)
	}
	return &obj, nil
}

func makeInstanaPutOnlyRestResourceSUT(client RestClient) RestResource {
	unmarshaller := &testUnmarshaller{}
	return NewPUTOnlyRestResource(testObjectResourcePath, unmarshaller, client)
}

func TestSuccessfulGetOneTestObjectThroughPutOnlyRestResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := makeInstanaPutOnlyRestResourceSUT(client)

	testObject := makeTestObject()
	serializedJSON, _ := json.Marshal(testObject)

	client.EXPECT().GetOne(gomock.Eq(testObject.ID), gomock.Eq(testObjectResourcePath)).Return(serializedJSON, nil)

	data, err := sut.GetOne(testObject.ID)

	assert.Nil(t, err)
	assert.Equal(t, testObject, data)
}

func TestShouldFailToGetOneTestObjectThroughPutOnlyRestResourceWhenErrorIsRetrievedFromRestClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := makeInstanaPutOnlyRestResourceSUT(client)

	client.EXPECT().GetOne(gomock.Eq(testObjectID), gomock.Eq(testObjectResourcePath)).Return(nil, errors.New("error during test"))

	_, err := sut.GetOne(testObjectID)

	assert.NotNil(t, err)
}

func TestShouldFailToGetOneTestObjectThroughPutOnlyRestResourceWhenResponseContainsAnInvalidJsonArray(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := makeInstanaPutOnlyRestResourceSUT(client)

	client.EXPECT().GetOne(gomock.Eq(testObjectID), gomock.Eq(testObjectResourcePath)).Return([]byte("[{ \"invalid\" : \"data\" }]"), nil)

	_, err := sut.GetOne(testObjectID)

	assert.NotNil(t, err)
}

func TestShouldFailToGetOneTestObjectThroughPutOnlyRestResourceWhenResponseContainsAnInvalidJsonObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := makeInstanaPutOnlyRestResourceSUT(client)

	client.EXPECT().GetOne(gomock.Eq(testObjectID), gomock.Eq(testObjectResourcePath)).Return([]byte("{ \"invalid\" : \"data\" }"), nil)

	_, err := sut.GetOne(testObjectID)

	assert.NotNil(t, err)
}

func TestShouldFailToGetOneTestObjectThroughPutOnlyRestResourceWhenResponseIsNotAJsonObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := makeInstanaPutOnlyRestResourceSUT(client)

	client.EXPECT().GetOne(gomock.Eq(testObjectID), gomock.Eq(testObjectResourcePath)).Return([]byte("Invalid Data"), nil)

	_, err := sut.GetOne(testObjectID)

	assert.NotNil(t, err)
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

		sut := makeInstanaPutOnlyRestResourceSUT(client)

		testObject := makeTestObject()
		serializedJSON, _ := json.Marshal(testObject)

		client.EXPECT().Put(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return(serializedJSON, nil)

		result, err := executeUpsertOperationThroughPutOnlyRestResource(operation, sut, testObject)

		assert.Nil(t, err)
		assert.Equal(t, testObject, result)
	})
}

func TestShouldFailToUpsertTestObjectThroughPutOnlyRestResourceWhenErrorIsReturnedFromRestClient(t *testing.T) {
	executeUpsertOperationThroughPutOnlyRestResourceTest(t, "TestShouldFailToUpsertTestObjectWhenErrorIsReturnedFromRestClient", func(operation string, t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		client := mocks.NewMockRestClient(ctrl)

		sut := makeInstanaPutOnlyRestResourceSUT(client)
		testObject := makeTestObject()

		client.EXPECT().Put(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return(nil, errors.New("Error during test"))

		_, err := executeUpsertOperationThroughPutOnlyRestResource(operation, sut, testObject)

		assert.NotNil(t, err)
	})
}

func TestShouldFailToUpsertTestObjectThroughPutOnlyRestResourceWhenResponseMessageIsNotAValidJsonObject(t *testing.T) {
	executeUpsertOperationThroughPutOnlyRestResourceTest(t, "TestShouldFailToUpsertTestObjectWhenResponseMessageIsNotAValidJsonObject", func(operation string, t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		client := mocks.NewMockRestClient(ctrl)

		sut := makeInstanaPutOnlyRestResourceSUT(client)
		testObject := makeTestObject()

		client.EXPECT().Put(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return([]byte("invalid response"), nil)

		_, err := executeUpsertOperationThroughPutOnlyRestResource(operation, sut, testObject)

		assert.NotNil(t, err)
	})
}

func TestShouldFailToUpsertTestObjectThroughPutOnlyRestResourceWhenResponseMessageContainsAnInvalidTestObject(t *testing.T) {
	executeUpsertOperationThroughPutOnlyRestResourceTest(t, "TestShouldFailToUpsertTestObjectWhenResponseMessageContainsAnInvalidTestObject", func(operation string, t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		client := mocks.NewMockRestClient(ctrl)

		sut := makeInstanaPutOnlyRestResourceSUT(client)
		testObject := makeTestObject()

		client.EXPECT().Put(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return([]byte("{ \"invalid\" : \"testObject\" }"), nil)

		_, err := executeUpsertOperationThroughPutOnlyRestResource(operation, sut, testObject)

		assert.NotNil(t, err)
	})
}

func TestShouldFailedToUpsertTestObjectThroughPutOnlyRestResourceWhenAnInvalidTestObjectIsProvided(t *testing.T) {
	executeUpsertOperationThroughPutOnlyRestResourceTest(t, "TestShouldFailedToUpsertTestObjectWhenAnInvalidTestObjectIsProvided", func(operation string, t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		client := mocks.NewMockRestClient(ctrl)

		sut := makeInstanaPutOnlyRestResourceSUT(client)
		testObject := &testObject{
			ID:   "some id",
			Name: "invalid name",
		}

		client.EXPECT().Put(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Times(0)

		_, err := executeUpsertOperationThroughPutOnlyRestResource(operation, sut, testObject)

		assert.NotNil(t, err)
	})
}

func TestSuccessfulDeleteOfTestObjectByObjectThroughPutOnlyRestResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := makeInstanaPutOnlyRestResourceSUT(client)
	testObject := makeTestObject()

	client.EXPECT().Delete(gomock.Eq(testObjectID), gomock.Eq(testObjectResourcePath)).Return(nil)

	err := sut.Delete(testObject)

	assert.Nil(t, err)
}

func TestShouldFailToDeleteTestObjectThroughPutOnlyRestResourceWhenErrorIsRetrunedFromRestClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := makeInstanaPutOnlyRestResourceSUT(client)
	testObject := makeTestObject()

	client.EXPECT().Delete(gomock.Eq(testObjectID), gomock.Eq(testObjectResourcePath)).Return(errors.New("Error during test"))

	err := sut.Delete(testObject)

	assert.NotNil(t, err)
}
