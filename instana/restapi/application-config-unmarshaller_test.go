package restapi_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/gessnerfl/terraform-provider-instana/instana/restapi"
)

func TestShouldSuccessfullyUnmarshalApplicationConfig(t *testing.T) {
	id := "test-application-config-id"
	label := "Test Application Config Label"
	applicationConfig := ApplicationConfig{
		ID:                 id,
		Label:              label,
		MatchSpecification: NewBinaryOperator(NewComparisionExpression("key", EqualsOperator, "value"), LogicalAnd, NewUnaryOperationExpression("key", NotBlankOperator)),
		Scope:              "scope",
	}

	serializedJSON, _ := json.Marshal(applicationConfig)

	result, err := NewApplicationConfigUnmarshaller().Unmarshal(serializedJSON)

	assert.Nil(t, err)
	assert.Equal(t, applicationConfig, result)
}

func TestShouldFailToUnmarashalApplicationConfigWhenResponseIsAJsonArray(t *testing.T) {
	response := `["foo","bar"]`

	_, err := NewApplicationConfigUnmarshaller().Unmarshal([]byte(response))

	assert.NotNil(t, err)
}

func TestShouldReturnEmptyApplicationConfigWhenNoFieldOfResponseMatchesToModel(t *testing.T) {
	response := `{"foo" : "bar"}`
	_, err := NewApplicationConfigUnmarshaller().Unmarshal([]byte(response))

	assert.NotNil(t, err)
}

func TestShouldFailToUnmarashalApplicationConfigWhenResponseIsNotAValidJson(t *testing.T) {
	response := `Invalid Data`

	_, err := NewApplicationConfigUnmarshaller().Unmarshal([]byte(response))

	assert.NotNil(t, err)
}

func TestShouldFailToUnmarashalApplicationConfigWhenExpressionTypeIsNotSupported(t *testing.T) {
	//config is invalid because there is no DType for the match specification.
	applicationConfig := ApplicationConfig{
		ID:    "id",
		Label: "label",
		MatchSpecification: TagMatcherExpression{
			Key:      "foo",
			Operator: NotEmptyOperator,
		},
		Scope: "scope",
	}
	serializedJSON, _ := json.Marshal(applicationConfig)

	_, err := NewApplicationConfigUnmarshaller().Unmarshal(serializedJSON)

	assert.NotNil(t, err)
}

func TestShouldFailToUnmarashalApplicationConfigWhenLeftSideOfBinaryExpressionTypeIsNotValid(t *testing.T) {
	left := TagMatcherExpression{
		Key:      "foo",
		Operator: NotEmptyOperator,
	}
	right := NewUnaryOperationExpression("foo", IsEmptyOperator)
	testShouldFailToUnmarashalApplicationConfigWhenOneSideOfBinaryExpressionIsNotValid(left, right, t)
}

func TestShouldFailToUnmarashalApplicationConfigWhenRightSideOfBinaryExpressionTypeIsNotValid(t *testing.T) {
	left := NewUnaryOperationExpression("foo", IsEmptyOperator)
	right := TagMatcherExpression{
		Key:      "foo",
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
	}
	serializedJSON, _ := json.Marshal(applicationConfig)

	_, err := NewApplicationConfigUnmarshaller().Unmarshal(serializedJSON)

	assert.NotNil(t, err)
}
