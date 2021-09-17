package restapi_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

const (
	testApplicationConfigId    = "test-application-config-id"
	testApplicationConfigLabel = "test-application-config-label"
)

func TestShouldSuccessfullyUnmarshalApplicationConfigWithMatchSpecification(t *testing.T) {
	id := testApplicationConfigId
	label := testApplicationConfigLabel
	applicationConfig := ApplicationConfig{
		ID:                 id,
		Label:              label,
		MatchSpecification: NewBinaryOperator(NewComparisonExpression("key", MatcherExpressionEntityDestination, EqualsOperator, "value"), LogicalAnd, NewUnaryOperationExpression("key", MatcherExpressionEntityDestination, NotBlankOperator)),
		Scope:              "scope",
		BoundaryScope:      "boundaryScope",
	}

	serializedJSON, _ := json.Marshal(applicationConfig)

	result, err := NewApplicationConfigUnmarshaller().Unmarshal(serializedJSON)

	require.NoError(t, err)
	require.Equal(t, &applicationConfig, result)
}

func TestShouldSuccessfullyUnmarshalApplicationConfigWithTagFilterExpressionContainingASingleTagFilter(t *testing.T) {
	value := "value"
	id := testApplicationConfigId
	label := testApplicationConfigLabel
	applicationConfig := ApplicationConfig{
		ID:                  id,
		Label:               label,
		TagFilterExpression: NewStringTagFilter(TagFilterEntityDestination, "entity.name", EqualsOperator, &value),
		Scope:               "scope",
		BoundaryScope:       "boundaryScope",
	}

	serializedJSON, _ := json.Marshal(applicationConfig)

	result, err := NewApplicationConfigUnmarshaller().Unmarshal(serializedJSON)

	require.NoError(t, err)
	require.Equal(t, &applicationConfig, result)
}

func TestShouldSuccessfullyUnmarshalApplicationConfigWithTagFilterExpressionContainingAnLogicalOr(t *testing.T) {
	value := "value"
	id := testApplicationConfigId
	label := testApplicationConfigLabel
	primaryExpression1 := NewStringTagFilter(TagFilterEntityDestination, "name1", EqualsOperator, &value)
	primaryExpression2 := NewStringTagFilter(TagFilterEntityDestination, "name2", EqualsOperator, &value)
	primaryExpression3 := NewStringTagFilter(TagFilterEntityDestination, "name3", EqualsOperator, &value)
	primaryExpression4 := NewStringTagFilter(TagFilterEntityDestination, "name4", EqualsOperator, &value)
	logicalOr := NewLogicalAndTagFilter([]TagFilterExpressionElement{primaryExpression1, primaryExpression2, NewLogicalAndTagFilter([]TagFilterExpressionElement{primaryExpression3, primaryExpression4})})
	applicationConfig := ApplicationConfig{
		ID:                  id,
		Label:               label,
		TagFilterExpression: logicalOr,
		Scope:               "scope",
		BoundaryScope:       "boundaryScope",
	}

	serializedJSON, _ := json.Marshal(applicationConfig)

	result, err := NewApplicationConfigUnmarshaller().Unmarshal(serializedJSON)

	require.NoError(t, err)
	require.Equal(t, &applicationConfig, result)
}

func TestShouldFailToUnmarshalApplicationConfigWhenResponseIsAJsonArray(t *testing.T) {
	response := `["foo","bar"]`

	_, err := NewApplicationConfigUnmarshaller().Unmarshal([]byte(response))

	require.Error(t, err)
}

func TestShouldReturnEmptyApplicationConfigWhenNoFieldOfResponseMatchesToModel(t *testing.T) {
	response := `{"foo" : "bar"}`
	config, err := NewApplicationConfigUnmarshaller().Unmarshal([]byte(response))

	require.NoError(t, err)
	require.Equal(t, &ApplicationConfig{}, config)
}

func TestShouldFailToUnmarshalApplicationConfigWhenResponseIsNotAValidJson(t *testing.T) {
	response := `Invalid Data`

	_, err := NewApplicationConfigUnmarshaller().Unmarshal([]byte(response))

	require.Error(t, err)
}

func TestShouldFailToUnmarshalApplicationConfigWhenExpressionTypeOfMatchSpecificationIsNotSupported(t *testing.T) {
	//config is invalid because there is no DType for the match specification.
	applicationConfig := ApplicationConfig{
		ID:    testApplicationConfigId,
		Label: testApplicationConfigLabel,
		MatchSpecification: TagMatcherExpression{
			Key:      "foo",
			Entity:   MatcherExpressionEntityDestination,
			Operator: NotEmptyOperator,
		},
		Scope:         "scope",
		BoundaryScope: "boundaryScope",
	}
	serializedJSON, _ := json.Marshal(applicationConfig)

	_, err := NewApplicationConfigUnmarshaller().Unmarshal(serializedJSON)

	require.Error(t, err)
}

func TestShouldFailToUnmarshalApplicationConfigWhenLeftSideOfBinaryExpressionTypeIsNotValid(t *testing.T) {
	left := &TagMatcherExpression{
		Key:      "foo",
		Operator: NotEmptyOperator,
	}
	right := NewUnaryOperationExpression("foo", MatcherExpressionEntityDestination, IsEmptyOperator)
	testShouldFailToUnmarshalApplicationConfigWhenOneSideOfBinaryExpressionIsNotValid(left, right, t)
}

func TestShouldFailToUnmarshalApplicationConfigWhenRightSideOfBinaryExpressionTypeIsNotValid(t *testing.T) {
	left := NewUnaryOperationExpression("foo", MatcherExpressionEntityDestination, IsEmptyOperator)
	right := &TagMatcherExpression{
		Key:      "foo",
		Entity:   MatcherExpressionEntityDestination,
		Operator: NotEmptyOperator,
	}
	testShouldFailToUnmarshalApplicationConfigWhenOneSideOfBinaryExpressionIsNotValid(left, right, t)
}

func testShouldFailToUnmarshalApplicationConfigWhenOneSideOfBinaryExpressionIsNotValid(left MatchExpression, right MatchExpression, t *testing.T) {
	applicationConfig := ApplicationConfig{
		ID:                 testApplicationConfigId,
		Label:              testApplicationConfigLabel,
		MatchSpecification: NewBinaryOperator(left, LogicalOr, right),
		Scope:              "scope",
		BoundaryScope:      "boundaryScope",
	}
	serializedJSON, _ := json.Marshal(applicationConfig)

	_, err := NewApplicationConfigUnmarshaller().Unmarshal(serializedJSON)

	require.Error(t, err)
}

func TestShouldFailToUnmarshalApplicationConfigWhenElementOfTagFilterExpressionIsNotValid(t *testing.T) {
	value := "value"
	primaryExpression := NewStringTagFilter(TagFilterEntityDestination, "name1", EqualsOperator, &value)
	invalidExpression := &TagFilterExpression{
		Type:            "INVALID",
		LogicalOperator: LogicalOr,
		Elements:        []TagFilterExpressionElement{},
	}
	applicationConfig := ApplicationConfig{
		ID:                  testApplicationConfigId,
		Label:               testApplicationConfigLabel,
		TagFilterExpression: NewLogicalOrTagFilter([]TagFilterExpressionElement{primaryExpression, invalidExpression}),
		Scope:               "scope",
		BoundaryScope:       "boundaryScope",
	}
	serializedJSON, _ := json.Marshal(applicationConfig)

	_, err := NewApplicationConfigUnmarshaller().Unmarshal(serializedJSON)

	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid tag filter element type INVALID")
}

func TestShouldFailToUnmarshalApplicationConfigWhenTagFilterIsNotAValidJsonObject(t *testing.T) {
	json := "{\"id\":\"test-application-config-id\",\"label\":\"test-application-config-label\",\"matchSpecification\":null,\"tagFilterExpression\":[\"foo\", \"bar\"],\"scope\":\"scope\",\"boundaryScope\":\"boundaryScope\"}"

	_, err := NewApplicationConfigUnmarshaller().Unmarshal([]byte(json))

	require.Error(t, err)
}
