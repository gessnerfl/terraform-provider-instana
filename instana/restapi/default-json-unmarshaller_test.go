package restapi_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

const (
	defaultObjectId   = "object-id"
	defaultObjectName = "object-name"
)

func TestShouldSuccessfullyUnmarshalSingleObject(t *testing.T) {
	testObject := TestObject{
		ID:   defaultObjectId,
		Name: defaultObjectName,
	}

	serializedJSON, _ := json.Marshal(testObject)

	sut := NewDefaultJSONUnmarshaller(&TestObject{})

	result, err := sut.Unmarshal(serializedJSON)

	require.NoError(t, err)
	require.Equal(t, &testObject, result)
}

func TestShouldSuccessfullyUnmarshalArrayOfObjects(t *testing.T) {
	testObject := TestObject{
		ID:   defaultObjectId,
		Name: defaultObjectName,
	}
	testObjects := []TestObject{testObject, testObject}
	arrayOfObjects := make([]TestObject, 0)

	serializedJSON, _ := json.Marshal(testObjects)

	sut := NewDefaultJSONUnmarshaller(&arrayOfObjects)

	result, err := sut.Unmarshal(serializedJSON)

	require.NoError(t, err)
	require.Equal(t, &testObjects, result)
}

func TestShouldFailToUnmarshalWhenObjectIsRequestedButResponseIsAJsonArray(t *testing.T) {
	testObject := TestObject{
		ID:   defaultObjectId,
		Name: defaultObjectName,
	}
	testObjects := []TestObject{testObject, testObject}

	serializedJSON, _ := json.Marshal(testObjects)

	sut := NewDefaultJSONUnmarshaller(&TestObject{})

	_, err := sut.Unmarshal(serializedJSON)

	require.Error(t, err)
}

func TestShouldFailToUnmarshalWhenResponseIsNotAJsonMessage(t *testing.T) {
	response := `foo bar`

	sut := NewDefaultJSONUnmarshaller(&TestObject{})

	_, err := sut.Unmarshal([]byte(response))

	require.Error(t, err)
}

func TestShouldReturnEmptyObjectWhenJsonObjectIsReturnWhereNoFiledMatches(t *testing.T) {
	response := `{"foo" : "bar" }`

	sut := NewDefaultJSONUnmarshaller(&TestObject{})

	result, err := sut.Unmarshal([]byte(response))

	require.NoError(t, err)
	require.Equal(t, &TestObject{}, result)
}

func TestShouldPanicToInitializeJSONUnmarshallerWhenObjectTypeIsNotAPointer(t *testing.T) {
	panicObject := recoverOnPanic(func() { NewDefaultJSONUnmarshaller(TestObject{}) })

	require.NotNil(t, panicObject)
	require.IsType(t, errors.New(""), panicObject)
	require.Contains(t, panicObject.(error).Error(), "objectType of defaultJSONUnmarshaller")
}

func recoverOnPanic(fn func()) (recovered interface{}) {
	defer func() {
		recovered = recover()
	}()
	fn()
	return
}

type TestObject struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
