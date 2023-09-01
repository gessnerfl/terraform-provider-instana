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

func TestSuccessfulCreateOfTestObjectThroughCreatePOSTUpdateNotSupportedRestResource(t *testing.T) {
	executeCreateOrUpdateOperationThroughCreatePOSTUpdateNotSupportedRestResourceTest(t, func(t *testing.T, sut RestResource[*testObject], client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller[*testObject]) {
		testObject := makeTestObject()
		serializedJSON, _ := json.Marshal(testObject)

		client.EXPECT().Post(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return(serializedJSON, nil)
		unmarshaller.EXPECT().Unmarshal(serializedJSON).Times(1).Return(testObject, nil)

		result, err := sut.Create(testObject)

		assert.NoError(t, err)
		assert.Equal(t, testObject, result)
	})
}

func TestShouldFailToCreateTestObjectThroughCreatePOSTUpdateNotSupportedRestResourceWhenErrorIsReturnedFromRestClient(t *testing.T) {
	executeCreateOrUpdateOperationThroughCreatePOSTUpdateNotSupportedRestResourceTest(t, func(t *testing.T, sut RestResource[*testObject], client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller[*testObject]) {
		testObject := makeTestObject()

		client.EXPECT().Post(gomock.Eq(testObject), gomock.Eq(testObjectResourcePath)).Return(nil, errors.New("error during test"))
		client.EXPECT().Put(gomock.Any(), gomock.Eq(testObjectResourcePath)).Times(0)
		unmarshaller.EXPECT().Unmarshal(gomock.Any()).Times(0)

		_, err := sut.Create(testObject)

		assert.Error(t, err)
	})
}

func TestShouldFailToCreateTestObjectThroughCreatePOSTUpdateNotSupportedRestResourceWhenResponseCannotBeUnmarshalled(t *testing.T) {
	executeCreateOrUpdateOperationThroughCreatePOSTUpdateNotSupportedRestResourceTest(t, func(t *testing.T, sut RestResource[*testObject], client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller[*testObject]) {
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

func TestShouldFailToUpdateTestObjectThroughCreatePOSTUpdateNotSupportedRestResourceWhenEmptyObjectCanBeCreated(t *testing.T) {
	executeCreateOrUpdateOperationThroughCreatePOSTUpdateNotSupportedRestResourceTest(t, func(t *testing.T, sut RestResource[*testObject], client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller[*testObject]) {
		testData := makeTestObject()
		emptyObject := &testObject{}

		unmarshaller.EXPECT().Unmarshal(gomock.Any()).Times(1).Return(emptyObject, nil)

		_, err := sut.Update(testData)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "update is not supported for /test")
	})
}

func TestShouldFailToUpdateTestObjectThroughCreatePOSTUpdateNotSupportedRestResourceWhenEmptyObjectCannotBeCreated(t *testing.T) {
	executeCreateOrUpdateOperationThroughCreatePOSTUpdateNotSupportedRestResourceTest(t, func(t *testing.T, sut RestResource[*testObject], client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller[*testObject]) {
		testData := makeTestObject()
		emptyObject := &testObject{}
		unmarshallingError := errors.New("unmarshalling-error")

		unmarshaller.EXPECT().Unmarshal(gomock.Any()).Times(1).Return(emptyObject, unmarshallingError)

		_, err := sut.Update(testData)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "update is not supported for /test; unmarshalling-error")
	})
}

func executeCreateOrUpdateOperationThroughCreatePOSTUpdateNotSupportedRestResourceTest(t *testing.T, testFunction func(t *testing.T, sut RestResource[*testObject], client *mocks.MockRestClient, unmarshaller *mocks.MockJSONUnmarshaller[*testObject])) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mocks.NewMockRestClient(ctrl)
	unmarshaller := mocks.NewMockJSONUnmarshaller[*testObject](ctrl)

	sut := NewCreatePOSTUpdateNotSupportedRestResource[*testObject](testObjectResourcePath, unmarshaller, client)

	testFunction(t, sut, client, unmarshaller)
}
