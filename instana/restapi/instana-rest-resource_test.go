package restapi_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	mocks "github.com/gessnerfl/terraform-provider-instana/mocks"
	"github.com/gessnerfl/terraform-provider-instana/testutils"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
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

	if err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}

	if !cmp.Equal(testObject, data) {
		t.Fatalf(testutils.ExpectedUnmarshalledJSONWithStruct, testObject, data, cmp.Diff(testObject, data))
	}
}

func TestShouldFailToGetOneTestObjectWhenErrorIsRetrievedFromRestClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := makeInstanaRestResourceSUT(client)

	client.EXPECT().GetOne(gomock.Eq(testObjectID), gomock.Eq(testObjectResourcePath)).Return(nil, errors.New("error during test"))

	_, err := sut.GetOne(testObjectID)

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
}

func TestShouldFailToGetOneTestObjectWhenResponseContainsAnInvalidJsonArray(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := makeInstanaRestResourceSUT(client)

	client.EXPECT().GetOne(gomock.Eq(testObjectID), gomock.Eq(testObjectResourcePath)).Return([]byte("[{ \"invalid\" : \"data\" }]"), nil)

	_, err := sut.GetOne(testObjectID)

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
}

func TestShouldFailToGetOneTestObjectWhenResponseContainsAnInvalidJsonObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := makeInstanaRestResourceSUT(client)

	client.EXPECT().GetOne(gomock.Eq(testObjectID), gomock.Eq(testObjectResourcePath)).Return([]byte("{ \"invalid\" : \"data\" }"), nil)

	_, err := sut.GetOne(testObjectID)

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
}

func TestShouldFailToGetOneTestObjectWhenResponseIsNotAJsonObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := makeInstanaRestResourceSUT(client)

	client.EXPECT().GetOne(gomock.Eq(testObjectID), gomock.Eq(testObjectResourcePath)).Return([]byte("Invalid Data"), nil)

	_, err := sut.GetOne(testObjectID)

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
}

func TestSuccessfulUpsertOfTestObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := makeInstanaRestResourceSUT(client)
	testObject := makeTestObject()
	serializedJSON, _ := json.Marshal(testObject)

	client.EXPECT().Put(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return(serializedJSON, nil)

	result, err := sut.Upsert(testObject)

	if err != nil {
		t.Fatalf(testutils.ExpectedNoErrorButGotMessage, err)
	}

	if !cmp.Equal(testObject, result) {
		t.Fatalf(testutils.ExpectedUnmarshalledJSONWithStruct, testObject, result, cmp.Diff(result, result))
	}
}

func TestShouldFailToUpsertTestObjectWhenErrorIsReturnedFromRestClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := makeInstanaRestResourceSUT(client)
	testObject := makeTestObject()

	client.EXPECT().Put(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return(nil, errors.New("Error during test"))

	_, err := sut.Upsert(testObject)

	if err == nil {
		t.Fatal(testutils.ExpectedError)
	}
}

func TestShouldFailToUpsertTestObjectWhenResponseMessageIsNotAValidJsonObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := makeInstanaRestResourceSUT(client)
	testObject := makeTestObject()

	client.EXPECT().Put(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return([]byte("invalid response"), nil)

	_, err := sut.Upsert(testObject)

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
}

func TestShouldFailToUpsertTestObjectWhenResponseMessageContainsAnInvalidTestObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := makeInstanaRestResourceSUT(client)
	testObject := makeTestObject()

	client.EXPECT().Put(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return([]byte("{ \"invalid\" : \"testObject\" }"), nil)

	_, err := sut.Upsert(testObject)

	if err == nil {
		t.Fatalf(testutils.ExpectedError)
	}
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

	_, err := sut.Upsert(testObject)

	if err == nil {
		t.Fatal(testutils.ExpectedError)
	}
}

func TestSuccessfulDeleteOfTestObjectByObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := makeInstanaRestResourceSUT(client)
	testObject := makeTestObject()

	client.EXPECT().Delete(gomock.Eq(testObjectID), gomock.Eq(testObjectResourcePath)).Return(nil)

	err := sut.Delete(testObject)

	if err != nil {
		t.Fatalf("Expected no error got %s", err)
	}
}

func TestShouldFailToDeleteTestObjectWhenErrorIsRetrunedFromRestClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)

	sut := makeInstanaRestResourceSUT(client)
	testObject := makeTestObject()

	client.EXPECT().Delete(gomock.Eq(testObjectID), gomock.Eq(testObjectResourcePath)).Return(errors.New("Error during test"))

	err := sut.Delete(testObject)

	if err == nil {
		t.Fatal(testutils.ExpectedError)
	}
}
