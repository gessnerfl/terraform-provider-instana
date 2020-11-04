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

func makeInstanaRestResourceSUT(client RestClient) RestResource {
	unmarshaller := &testUnmarshaller{}
	return NewRestResource(testObjectResourcePath, unmarshaller, client)
}

func TestSuccessfulGetOneTestObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := makeInstanaRestResourceSUT(client)

	testObject := makeTestObject()
	serializedJSON, _ := json.Marshal(testObject)

	client.EXPECT().GetOne(gomock.Eq(testObject.ID), gomock.Eq(testObjectResourcePath)).Return(serializedJSON, nil)

	data, err := sut.GetOne(testObject.ID)

	assert.Nil(t, err)
	assert.Equal(t, testObject, data)
}

func TestShouldFailToGetOneTestObjectWhenErrorIsRetrievedFromRestClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := makeInstanaRestResourceSUT(client)

	client.EXPECT().GetOne(gomock.Eq(testObjectID), gomock.Eq(testObjectResourcePath)).Return(nil, errors.New("error during test"))

	_, err := sut.GetOne(testObjectID)

	assert.NotNil(t, err)
}

func TestShouldFailToGetOneTestObjectWhenResponseContainsAnInvalidJsonArray(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := makeInstanaRestResourceSUT(client)

	client.EXPECT().GetOne(gomock.Eq(testObjectID), gomock.Eq(testObjectResourcePath)).Return([]byte("[{ \"invalid\" : \"data\" }]"), nil)

	_, err := sut.GetOne(testObjectID)

	assert.NotNil(t, err)
}

func TestShouldFailToGetOneTestObjectWhenResponseContainsAnInvalidJsonObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := makeInstanaRestResourceSUT(client)

	client.EXPECT().GetOne(gomock.Eq(testObjectID), gomock.Eq(testObjectResourcePath)).Return([]byte("{ \"invalid\" : \"data\" }"), nil)

	_, err := sut.GetOne(testObjectID)

	assert.NotNil(t, err)
}

func TestShouldFailToGetOneTestObjectWhenResponseIsNotAJsonObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := makeInstanaRestResourceSUT(client)

	client.EXPECT().GetOne(gomock.Eq(testObjectID), gomock.Eq(testObjectResourcePath)).Return([]byte("Invalid Data"), nil)

	_, err := sut.GetOne(testObjectID)

	assert.NotNil(t, err)
}

func TestSuccessfulUpsertOfTestObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := makeInstanaRestResourceSUT(client)
	testObject := makeTestObject()
	serializedJSON, _ := json.Marshal(testObject)

	client.EXPECT().Put(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return(serializedJSON, nil)

	result, err := sut.Update(testObject)

	assert.Nil(t, err)
	assert.Equal(t, testObject, result)
}

func TestShouldFailToUpsertTestObjectWhenErrorIsReturnedFromRestClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := makeInstanaRestResourceSUT(client)
	testObject := makeTestObject()

	client.EXPECT().Put(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return(nil, errors.New("Error during test"))

	_, err := sut.Update(testObject)

	assert.NotNil(t, err)
}

func TestShouldFailToUpsertTestObjectWhenResponseMessageIsNotAValidJsonObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := makeInstanaRestResourceSUT(client)
	testObject := makeTestObject()

	client.EXPECT().Put(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return([]byte("invalid response"), nil)

	_, err := sut.Update(testObject)

	assert.NotNil(t, err)
}

func TestShouldFailToUpsertTestObjectWhenResponseMessageContainsAnInvalidTestObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := makeInstanaRestResourceSUT(client)
	testObject := makeTestObject()

	client.EXPECT().Put(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return([]byte("{ \"invalid\" : \"testObject\" }"), nil)

	_, err := sut.Update(testObject)

	assert.NotNil(t, err)
}

func TestShouldFailedToUpsertTestObjectWhenAnInvalidTestObjectIsProvided(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := makeInstanaRestResourceSUT(client)
	testObject := &testObject{
		ID:   "some id",
		Name: "invalid name",
	}

	client.EXPECT().Put(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Times(0)

	_, err := sut.Update(testObject)

	assert.NotNil(t, err)
}

func TestSuccessfulDeleteOfTestObjectByObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := makeInstanaRestResourceSUT(client)
	testObject := makeTestObject()

	client.EXPECT().Delete(gomock.Eq(testObjectID), gomock.Eq(testObjectResourcePath)).Return(nil)

	err := sut.Delete(testObject)

	assert.Nil(t, err)
}

func TestShouldFailToDeleteTestObjectWhenErrorIsRetrunedFromRestClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := makeInstanaRestResourceSUT(client)
	testObject := makeTestObject()

	client.EXPECT().Delete(gomock.Eq(testObjectID), gomock.Eq(testObjectResourcePath)).Return(errors.New("Error during test"))

	err := sut.Delete(testObject)

	assert.NotNil(t, err)
}
