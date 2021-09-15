package restapi_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

func TestShouldSuccessfullyUnmarshalApplicationConfig(t *testing.T) {
	id := "test-application-config-id"
	label := "Test Application Config Label"
	applicationConfig := ApplicationConfig{
		ID:                 id,
		Label:              label,
		MatchSpecification: NewBinaryOperator(NewComparisionExpression("key", MatcherExpressionEntityDestination, EqualsOperator, "value"), LogicalAnd, NewUnaryOperationExpression("key", MatcherExpressionEntityDestination, NotBlankOperator)),
		Scope:              "scope",
		BoundaryScope:      "boundaryScope",
	}

	serializedJSON, _ := json.Marshal(applicationConfig)

	result, err := NewApplicationConfigUnmarshaller().Unmarshal(serializedJSON)

	require.NoError(t, err)
	require.Equal(t, &applicationConfig, result)
}

func TestShouldFailToUnmarashalApplicationConfigWhenResponseIsAJsonArray(t *testing.T) {
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

func TestShouldFailToUnmarshalApplicationConfigWhenExpressionTypeIsNotSupported(t *testing.T) {
	//config is invalid because there is no DType for the match specification.
	applicationConfig := ApplicationConfig{
		ID:    "id",
		Label: "label",
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

func TestShouldFailToUnmarashalApplicationConfigWhenLeftSideOfBinaryExpressionTypeIsNotValid(t *testing.T) {
	left := TagMatcherExpression{
		Key:      "foo",
		Operator: NotEmptyOperator,
	}
	right := NewUnaryOperationExpression("foo", MatcherExpressionEntityDestination, IsEmptyOperator)
	testShouldFailToUnmarashalApplicationConfigWhenOneSideOfBinaryExpressionIsNotValid(left, right, t)
}

func TestShouldFailToUnmarashalApplicationConfigWhenRightSideOfBinaryExpressionTypeIsNotValid(t *testing.T) {
	left := NewUnaryOperationExpression("foo", MatcherExpressionEntityDestination, IsEmptyOperator)
	right := TagMatcherExpression{
		Key:      "foo",
		Entity:   MatcherExpressionEntityDestination,
		Operator: NotEmptyOperator,
	}
	testShouldFailToUnmarashalApplicationConfigWhenOneSideOfBinaryExpressionIsNotValid(left, right, t)
}

func testShouldFailToUnmarashalApplicationConfigWhenOneSideOfBinaryExpressionIsNotValid(left MatchExpression, right MatchExpression, t *testing.T) {
	applicationConfig := ApplicationConfig{
		ID:                 "id",
		Label:              "label",
		MatchSpecification: NewBinaryOperator(left, LogicalOr, right),
		Scope:              "scope",
		BoundaryScope:      "boundaryScope",
	}
	serializedJSON, _ := json.Marshal(applicationConfig)

	_, err := NewApplicationConfigUnmarshaller().Unmarshal(serializedJSON)

	require.Error(t, err)
}
