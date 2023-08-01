package restapi_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

const (
	defaultObjectId   = "object-id"
	defaultObjectName = "object-name"
)

func TestShouldSuccessfullyUnmarshalSingleObject(t *testing.T) {
	testData := &testObject{
		ID:   defaultObjectId,
		Name: defaultObjectName,
	}

	serializedJSON, _ := json.Marshal(testData)

	sut := NewDefaultJSONUnmarshaller(&testObject{})

	result, err := sut.Unmarshal(serializedJSON)

	require.NoError(t, err)
	require.Equal(t, testData, result)
}

func TestShouldSuccessfullyUnmarshalArrayOfObjects(t *testing.T) {
	testData := &testObject{
		ID:   defaultObjectId,
		Name: defaultObjectName,
	}
	testObjects := &[]*testObject{testData, testData}

	serializedJSON, _ := json.Marshal(testObjects)

	sut := NewDefaultJSONUnmarshaller(testData)

	result, err := sut.UnmarshalArray(serializedJSON)

	require.NoError(t, err)
	require.Equal(t, testObjects, result)
}

func TestShouldFailToUnmarshalArrayWhenNoValidJsonIsProvided(t *testing.T) {
	sut := NewDefaultJSONUnmarshaller(&testObject{})

	_, err := sut.UnmarshalArray([]byte("invalid json data"))

	require.Error(t, err)
}

func TestShouldFailToUnmarshalWhenObjectIsRequestedButResponseIsAJsonArray(t *testing.T) {
	testData := &testObject{
		ID:   defaultObjectId,
		Name: defaultObjectName,
	}
	testObjects := []*testObject{testData, testData}

	serializedJSON, _ := json.Marshal(testObjects)

	sut := NewDefaultJSONUnmarshaller(&testObject{})

	_, err := sut.Unmarshal(serializedJSON)

	require.Error(t, err)
}

func TestShouldFailToUnmarshalWhenResponseIsNotAJsonMessage(t *testing.T) {
	response := `foo bar`

	sut := NewDefaultJSONUnmarshaller(&testObject{})

	_, err := sut.Unmarshal([]byte(response))

	require.Error(t, err)
}

func TestShouldReturnEmptyObjectWhenJsonObjectIsReturnWhereNoFiledMatches(t *testing.T) {
	response := `{"foo" : "bar" }`

	sut := NewDefaultJSONUnmarshaller(&testObject{})

	result, err := sut.Unmarshal([]byte(response))

	require.NoError(t, err)
	require.Equal(t, &testObject{}, result)
}
