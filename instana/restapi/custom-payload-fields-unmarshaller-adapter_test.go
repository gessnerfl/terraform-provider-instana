package restapi_test

import (
	"errors"
	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
	"github.com/gessnerfl/terraform-provider-instana/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestApplicationAlertConfigUnmarshaller(t *testing.T) {
	ut := &applicationAlertConfigUnmarshallerTest{}
	t.Run("should successfully map custom payload field aware object", ut.shouldSuccessfullyMapCustomPayloadFieldAwareObject)
	t.Run("should fail to map custom payload fields when unmarshalling fails", ut.shouldFailToMapCustomPayloadFieldAwareObjectWhenUnmarshallingFails)
	t.Run("should successfully map array of custom payload field aware objects", ut.shouldSuccessfullyMapArrayOfCustomPayloadFieldAwareObject)
	t.Run("should fail to unmarshal array of custom payload field aware object when unmarshalling fails", ut.shouldFailToMapArrayOfCustomPayloadFieldAwareObjectWhenUnmarshallingFails)
}

type applicationAlertConfigUnmarshallerTest struct {
}

func (ut *applicationAlertConfigUnmarshallerTest) shouldSuccessfullyMapCustomPayloadFieldAwareObject(t *testing.T) {
	tagKey := "tag-key"
	unmarshalled := &customFieldAwareTestObject{
		ID:   testObjectID,
		Name: testObjectName,
		CustomerPayloadFields: []CustomPayloadField[any]{
			{
				Type:  StaticStringCustomPayloadType,
				Key:   "key1",
				Value: "value1",
			},
			{
				Type:  DynamicCustomPayloadType,
				Key:   "key2",
				Value: map[string]interface{}{"key": tagKey, "tagName": "tag"},
			},
		},
	}
	serializedJSON := []byte("test-data")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	unmarshaler := mocks.NewMockJSONUnmarshaller[*customFieldAwareTestObject](ctrl)
	unmarshaler.EXPECT().Unmarshal(serializedJSON).Times(1).Return(unmarshalled, nil)

	result, err := NewCustomPayloadFieldsUnmarshallerAdapter[*customFieldAwareTestObject](unmarshaler).Unmarshal(serializedJSON)

	require.NoError(t, err)
	require.Equal(t, &customFieldAwareTestObject{
		ID:   testObjectID,
		Name: testObjectName,
		CustomerPayloadFields: []CustomPayloadField[any]{
			{
				Type:  StaticStringCustomPayloadType,
				Key:   "key1",
				Value: StaticStringCustomPayloadFieldValue("value1"),
			},
			{
				Type: DynamicCustomPayloadType,
				Key:  "key2",
				Value: DynamicCustomPayloadFieldValue{
					TagName: "tag",
					Key:     &tagKey,
				},
			},
		},
	}, result)
}

func (ut *applicationAlertConfigUnmarshallerTest) shouldFailToMapCustomPayloadFieldAwareObjectWhenUnmarshallingFails(t *testing.T) {
	expectedError := errors.New("test error")
	serializedJSON := []byte("test-data")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	unmarshaler := mocks.NewMockJSONUnmarshaller[*customFieldAwareTestObject](ctrl)
	unmarshaler.EXPECT().Unmarshal(serializedJSON).Times(1).Return(&customFieldAwareTestObject{}, expectedError)

	_, err := NewCustomPayloadFieldsUnmarshallerAdapter[*customFieldAwareTestObject](unmarshaler).Unmarshal(serializedJSON)

	require.Equal(t, expectedError, err)
}

func (ut *applicationAlertConfigUnmarshallerTest) shouldSuccessfullyMapArrayOfCustomPayloadFieldAwareObject(t *testing.T) {
	tagKey := "tag-key"
	unmarshalled := []*customFieldAwareTestObject{
		{
			ID:   testObjectID,
			Name: testObjectName,
			CustomerPayloadFields: []CustomPayloadField[any]{
				{
					Type:  StaticStringCustomPayloadType,
					Key:   "key1",
					Value: "value1",
				},
				{
					Type:  DynamicCustomPayloadType,
					Key:   "key2",
					Value: map[string]interface{}{"key": tagKey, "tagName": "tag"},
				},
			},
		},
	}
	serializedJSON := []byte("test-data")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	unmarshaler := mocks.NewMockJSONUnmarshaller[*customFieldAwareTestObject](ctrl)
	unmarshaler.EXPECT().UnmarshalArray(serializedJSON).Times(1).Return(&unmarshalled, nil)

	result, err := NewCustomPayloadFieldsUnmarshallerAdapter[*customFieldAwareTestObject](unmarshaler).UnmarshalArray(serializedJSON)

	require.NoError(t, err)
	require.Equal(t, &[]*customFieldAwareTestObject{
		{
			ID:   testObjectID,
			Name: testObjectName,
			CustomerPayloadFields: []CustomPayloadField[any]{
				{
					Type:  StaticStringCustomPayloadType,
					Key:   "key1",
					Value: StaticStringCustomPayloadFieldValue("value1"),
				},
				{
					Type: DynamicCustomPayloadType,
					Key:  "key2",
					Value: DynamicCustomPayloadFieldValue{
						TagName: "tag",
						Key:     &tagKey,
					},
				},
			},
		},
	}, result)
}

func (ut *applicationAlertConfigUnmarshallerTest) shouldFailToMapArrayOfCustomPayloadFieldAwareObjectWhenUnmarshallingFails(t *testing.T) {
	expectedError := errors.New("test error")
	serializedJSON := []byte("test-data")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	unmarshaler := mocks.NewMockJSONUnmarshaller[*customFieldAwareTestObject](ctrl)
	unmarshaler.EXPECT().UnmarshalArray(serializedJSON).Times(1).Return(&[]*customFieldAwareTestObject{}, expectedError)

	_, err := NewCustomPayloadFieldsUnmarshallerAdapter[*customFieldAwareTestObject](unmarshaler).UnmarshalArray(serializedJSON)

	require.Equal(t, expectedError, err)
}

type customFieldAwareTestObject struct {
	ID                    string
	Name                  string
	CustomerPayloadFields []CustomPayloadField[any]
}

func (t *customFieldAwareTestObject) GetIDForResourcePath() string {
	return t.ID
}

func (t *customFieldAwareTestObject) GetCustomerPayloadFields() []CustomPayloadField[any] {
	return t.CustomerPayloadFields
}

func (t *customFieldAwareTestObject) SetCustomerPayloadFields(fields []CustomPayloadField[any]) {
	t.CustomerPayloadFields = fields
}
